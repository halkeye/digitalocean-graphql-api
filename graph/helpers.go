package graph

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
)

func mustStringPtr(s string) *string {
	return &s
}

type doContextKey string

const DOContextKey = doContextKey("DoContextKey")

func DoClientFromContext(ctx context.Context) (*godo.Client, error) {
	doClient := ctx.Value(DOContextKey)
	if doClient == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	do, ok := doClient.(*godo.Client)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return do, nil
}
