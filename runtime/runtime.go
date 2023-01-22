package runtime

import (
	"net/http"
	"net/http/httptest"
)

// TODO: User protobufs instead of JSON

func NewResultMetadata(r *http.Response) *ResultMetadata {
	return &ResultMetadata{
		StatusCode: r.StatusCode,
		Header:     r.Header,
	}
}

// TODO: Add some kinda settings for runtime
func NewLambdaRuntime() *LambdaRuntime {
	return &LambdaRuntime{}
}

func (r *LambdaRuntime) RegisterHandlerFunc(lambdaFunc LambdaFunction) {
	r.Handler = lambdaFunc
}

func (r *LambdaRuntime) Run() error {
	input, err := r.parseInput()
	if err != nil {
		return err
	}

	w := httptest.NewRecorder()
	r.Handler(w, input)
	result := w.Result()
	resultMetadata := NewResultMetadata(result)

	r.writeOutput(result, resultMetadata)

	return nil
}
