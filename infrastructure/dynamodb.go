package infrastructure

// Karma is attribute struct of dynamodb
type Karma struct {
	Identifier        string `json:"Identifier" dynamo:"Identifier,hash"`
	Score             int    `json:"Score" dynamo:"Score"`
}
