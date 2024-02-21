package platform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"integration_test/model/platform"
	"integration_test/util"
	"io/ioutil"
	"net/http"
)

type quoteClient struct {
	client *http.Client
}

func NewQuoteClient() QuoteInterface {
	client := &http.Client{}

	return &quoteClient{client: client}
}

func (q *quoteClient) Get(ctx context.Context, category string) (res []platform.QuoteResponse, errRes *platform.QuoteError, err error) {
	url := fmt.Sprintf("%s%s%s", util.Configuration.Platform.Host, "/v1/quotes", category)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, err
	}

	req.Header.Set("X-API-KEY", util.Configuration.Platform.APIKey)

	response, err := q.client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, nil, err
	}

	err = json.Unmarshal(body, &res)
	if err == nil && len(res) > 0 {
		return res, nil, nil
	}

	err = json.Unmarshal(body, &errRes)
	if err == nil {
		return nil, errRes, nil
	}

	return nil, nil, errors.New("invalid response")
}
