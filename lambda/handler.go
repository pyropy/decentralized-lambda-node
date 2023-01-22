package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type handlerFunc = func(ctx context.Context, input []byte) (io.Reader, error)

func newHandler(handler interface{}) handlerFunc {
	return reflectHandler(handler)
}

type jsonOutBuffer struct {
	*bytes.Buffer
}

func errorHandler(err error) handlerFunc {
	return func(_ context.Context, _ []byte) (io.Reader, error) {
		return nil, err
	}
}

// handlerTakesContext returns whether the handler takes a context.Context as its first argument.
func handlerTakesContext(handler reflect.Type) (bool, error) {
	switch handler.NumIn() {
	case 0:
		return false, nil
	case 1:
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
		argumentType := handler.In(0)
		if argumentType.Kind() != reflect.Interface {
			return false, nil
		}

		// handlers like func(event any) are valid.
		if argumentType.NumMethod() == 0 {
			return false, nil
		}

		if !contextType.Implements(argumentType) || !argumentType.Implements(contextType) {
			return false, fmt.Errorf("handler takes an interface, but it is not context.Context: %q", argumentType.Name())
		}
		return true, nil
	case 2:
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
		argumentType := handler.In(0)
		if argumentType.Kind() != reflect.Interface || !contextType.Implements(argumentType) || !argumentType.Implements(contextType) {
			return false, fmt.Errorf("handler takes two arguments, but the first is not Context. got %s", argumentType.Kind())
		}
		return true, nil
	}
	return false, fmt.Errorf("handlers may not take more than two arguments, but handler takes %d", handler.NumIn())
}

func validateReturns(handler reflect.Type) error {
	errorType := reflect.TypeOf((*error)(nil)).Elem()

	switch n := handler.NumOut(); {
	case n > 2:
		return fmt.Errorf("handler may not return more than two values")
	case n > 1:
		if !handler.Out(1).Implements(errorType) {
			return fmt.Errorf("handler returns two values, but the second does not implement error")
		}
	case n == 1:
		if !handler.Out(0).Implements(errorType) {
			return fmt.Errorf("handler returns a single value, but it does not implement error")
		}
	}

	return nil
}

func reflectHandler(f interface{}) handlerFunc {
	if f == nil {
		return errorHandler(errors.New("handler is nil"))
	}

	handler := reflect.ValueOf(f)
	handlerType := reflect.TypeOf(f)
	if handlerType.Kind() != reflect.Func {
		return errorHandler(fmt.Errorf("handler kind %s is not %s", handlerType.Kind(), reflect.Func))
	}

	takesContext, err := handlerTakesContext(handlerType)
	if err != nil {
		return errorHandler(err)
	}

	if err := validateReturns(handlerType); err != nil {
		return errorHandler(err)
	}

	out := &jsonOutBuffer{bytes.NewBuffer(nil)}
	return func(ctx context.Context, payload []byte) (io.Reader, error) {
		out.Reset()
		in := bytes.NewBuffer(payload)
		decoder := json.NewDecoder(in)
		encoder := json.NewEncoder(out)

		// construct arguments
		var args []reflect.Value
		if takesContext {
			args = append(args, reflect.ValueOf(ctx))
		}
		if (handlerType.NumIn() == 1 && !takesContext) || handlerType.NumIn() == 2 {
			eventType := handlerType.In(handlerType.NumIn() - 1)
			event := reflect.New(eventType)
			if err := decoder.Decode(event.Interface()); err != nil {
				return nil, err
			}
			args = append(args, event.Elem())
		}

		response := handler.Call(args)

		// return the error, if any
		if len(response) > 0 {
			if errVal, ok := response[len(response)-1].Interface().(error); ok && errVal != nil {
				return nil, errVal
			}
		}
		// set the response value, if any
		var val interface{}
		if len(response) > 1 {
			val = response[0].Interface()
		}

		// encode to JSON
		if err := encoder.Encode(val); err != nil {
			// if response is not JSON serializable, but the response type is a reader, return it as-is
			if reader, ok := val.(io.Reader); ok {
				return reader, nil
			}
			return nil, err
		}

		// if response value is an io.Reader, return it as-is
		if reader, ok := val.(io.Reader); ok {
			// back-compat, don't return the reader if the value serialized to a non-empty json
			if strings.HasPrefix(out.String(), "{}") {
				return reader, nil
			}
		}

		return out, nil
	}

}
