package testing

import (
	"encoding/json"
	"fmt"
	"integration_test/model/delivery"
	"integration_test/model/entity"
	"integration_test/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

type test struct {
	name  string
	param string
	want  want
}

type want struct {
	httpCode int
	res      *delivery.Response
	success  bool
	data     entity.Quote
}

func (suite *testSvc) TestAPI() {
	t := suite.T()

	tests := []test{
		{
			name:  "success",
			param: "?category=freedom",
			want: want{
				httpCode: 200,
				res: &delivery.Response{
					Quote:    "if you don't fight, you won't win",
					Author:   "eren yeager",
					Category: "freedom",
				},
				success: true,
				data: entity.Quote{
					Quote:    "if you don't fight, you won't win",
					Author:   "eren yeager",
					Category: "freedom",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpCode, res, err := suite.requestSend(t, test.param)
			assert.NoError(t, err)

			assert.Equal(t, test.want.httpCode, httpCode)
			assert.Equal(t, test.want.res, res)

			if test.want.success {
				suite.validateData(t, test.want.data)
			}
		})
	}
}

func (s *testSvc) validateData(t *testing.T, data entity.Quote) {
	filter := bson.M{
		"quote": data.Quote,
	}

	var resData entity.Quote
	coll := s.mongoClient.Database(util.Configuration.MongoDB.Database).Collection(util.Configuration.MongoDB.Collection)
	coll.FindOne(s.ctx, filter).Decode(&resData)

	t.Log(resData)

	assert.Equal(t, data, resData)
}

func (s *testSvc) requestSend(t *testing.T, param string) (int, *delivery.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:30103%s%s", util.Configuration.Server.Endpoint, param), nil)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	respBody := json.NewDecoder(resp.Body)

	var responseMap map[string]interface{}
	err = respBody.Decode(&responseMap)
	if err != nil {
		return 0, nil, err
	}

	respBytes, _ := json.Marshal(responseMap)

	var response delivery.Response
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, &response, nil
}
