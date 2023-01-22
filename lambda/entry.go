package lambda

import (
	"context"
)

// TODO: User protobufs instead of JSON
// TODO: Add some kinda settings for handler func
func Start(handler interface{}) error {
	return start(handler)
}

func start(handler interface{}) error {
	input, err := parseInput()
	if err != nil {
		return err
	}

	handlerFunc := newHandler(handler)
	ctx := context.Background()

	response, err := handlerFunc(ctx, input)
	if err != nil {
		return err
	}

	err = writeOutput(response)
	if err != nil {
		return err
	}
	return nil
}
