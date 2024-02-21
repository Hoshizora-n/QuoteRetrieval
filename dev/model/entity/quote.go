package entity

type Quote struct {
	Quote    string `bson:"quote"`
	Author   string `bson:"author"`
	Category string `bson:"category"`
}
