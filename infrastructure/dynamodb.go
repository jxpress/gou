package infrastructure

// Karma is attribute struct of dynamodb
type Karma struct {
	Identifier        string `json:"identifier" dynamo:"identifier,hash"`
	Score             int    `json:"score" dynamo:"score"`
}
