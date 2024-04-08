package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
)

type doContextKey string

const DOContextKey = doContextKey("DoContextKey")

func For(ctx context.Context) (*godo.Client, error) {
	doClient := ctx.Value(DOContextKey)
	if doClient == nil {
		err := fmt.Errorf("could not retrieve godo.Client")
		return nil, err
	}

	do, ok := doClient.(*godo.Client)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return do, nil
}

func WithContext(ctx context.Context, client *godo.Client) context.Context {
	return context.WithValue(ctx, DOContextKey, client)
}
