package runtime

import "net/http"

type ResultMetadata struct {
	StatusCode int         `json:"statusCode"`
	Header     http.Header `json:"header"`
}

// TODO: Replace argument types with custom types
type LambdaFunction = func(w http.ResponseWriter, r *http.Request)

type LambdaRuntime struct {
	Handler LambdaFunction
}
