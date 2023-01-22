package main

import (
	"context"
	"fmt"
	"github.com/pyropy/decentralised-lambda/lambda"
)

type ExampleInput struct {
	Name string
}

func Handler(ctx context.Context, input ExampleInput) (string, error) {
	return fmt.Sprintf("Hello %s!", input.Name), nil
}

func main() {
	lambda.Start(Handler)
}
