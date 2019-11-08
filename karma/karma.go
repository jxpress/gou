package karma

import (
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/guregu/dynamo"
	"github.com/jxpress/gou/infrastructure"
	moduledynamodb "github.com/jxpress/gou/module/dynamodb"
)

func GetAll(config *aws.Config, identifier string) (karma []*infrastructure.Karma, err error) {
	c, err := moduledynamodb.NewClient(config)
	if err != nil {
		return
	}
	err = c.GetAllItems("Karma", &moduledynamodb.HashKey{
		Name:  "Identifier",
		Value: identifier,
	}, &karma)
	if err != nil {
		return
	}
	return
}

func Put(config *aws.Config, karma *infrastructure.Karma) (err error) {
	c, err := moduledynamodb.NewClient(config)
	if err != nil {
		return
	}
	if err := c.Put("Karma", karma); err != nil {
		return
	}
	return
}
