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
		err := fmt.Errorf("godo.Client has wrong type")
		return nil, err
	}
	return do, nil
}

func WithContext(ctx context.Context, bearerToken string) context.Context {
	client := godo.NewFromToken(bearerToken)
	return context.WithValue(ctx, DOContextKey, client)
}
