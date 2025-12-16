//go:build js && wasm
// +build js,wasm

package wasm

import (
	"context"
	"errors"
	"syscall/js"

	"github.com/nkechr1s/gofetch/infrastructure"
)

var defaultClient *infrastructure.Client

func init() {
	defaultClient = infrastructure.NewClient()
}

// ExposeFunctions exposes GoFetch functions to JavaScript.
func ExposeFunctions() {
	js.Global().Set("gofetch", js.ValueOf(map[string]interface{}{
		"newClient":  js.FuncOf(newClient),
		"get":        js.FuncOf(get),
		"post":       js.FuncOf(post),
		"put":        js.FuncOf(put),
		"patch":      js.FuncOf(patch),
		"delete":     js.FuncOf(delete_),
		"setBaseURL": js.FuncOf(setBaseURL),
		"setTimeout": js.FuncOf(setTimeout),
		"setHeader":  js.FuncOf(setHeader),
	}))
}

// newClient creates a new client instance.
func newClient(this js.Value, args []js.Value) interface{} {
	client := infrastructure.NewClient()

	// Return a JavaScript object with methods
	return map[string]interface{}{
		"get":         js.FuncOf(makeGetFunc(client)),
		"post":        js.FuncOf(makePostFunc(client)),
		"put":         js.FuncOf(makePutFunc(client)),
		"patch":       js.FuncOf(makePatchFunc(client)),
		"delete":      js.FuncOf(makeDeleteFunc(client)),
		"setBaseURL":  js.FuncOf(makeSetBaseURLFunc(client)),
		"setTimeout":  js.FuncOf(makeSetTimeoutFunc(client)),
		"setHeader":   js.FuncOf(makeSetHeaderFunc(client)),
		"newInstance": js.FuncOf(makeNewInstanceFunc(client)),
	}
}

// get performs a GET request using the default client.
func get(this js.Value, args []js.Value) interface{} {
	return makeGetFunc(defaultClient)(this, args)
}

// post performs a POST request using the default client.
func post(this js.Value, args []js.Value) interface{} {
	return makePostFunc(defaultClient)(this, args)
}

// put performs a PUT request using the default client.
func put(this js.Value, args []js.Value) interface{} {
	return makePutFunc(defaultClient)(this, args)
}

// patch performs a PATCH request using the default client.
func patch(this js.Value, args []js.Value) interface{} {
	return makePatchFunc(defaultClient)(this, args)
}

// delete_ performs a DELETE request using the default client.
func delete_(this js.Value, args []js.Value) interface{} {
	return makeDeleteFunc(defaultClient)(this, args)
}

// setBaseURL sets the base URL on the default client.
func setBaseURL(this js.Value, args []js.Value) interface{} {
	return makeSetBaseURLFunc(defaultClient)(this, args)
}

// setTimeout sets the timeout on the default client.
func setTimeout(this js.Value, args []js.Value) interface{} {
	return makeSetTimeoutFunc(defaultClient)(this, args)
}

// setHeader sets a header on the default client.
func setHeader(this js.Value, args []js.Value) interface{} {
	return makeSetHeaderFunc(defaultClient)(this, args)
}

// Helper functions to create closures for specific client instances

func makeGetFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		return promiseWrapper(func() (interface{}, error) {
			if len(args) < 1 {
				return nil, errors.New("path is required")
			}

			path := args[0].String()
			var params map[string]interface{}

			if len(args) >= 2 && !args[1].IsUndefined() && !args[1].IsNull() {
				params = jsObjectToMap(args[1])
			}

			var target interface{}
			resp, err := client.Get(context.Background(), path, params, &target)
			if err != nil {
				return nil, err
			}

			return responseToJS(resp), nil
		})
	}
}

func makePostFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		return promiseWrapper(func() (interface{}, error) {
			if len(args) < 1 {
				return nil, errors.New("path is required")
			}

			path := args[0].String()
			var params map[string]interface{}
			var body interface{}

			if len(args) >= 2 && !args[1].IsUndefined() && !args[1].IsNull() {
				params = jsObjectToMap(args[1])
			}

			if len(args) >= 3 && !args[2].IsUndefined() && !args[2].IsNull() {
				body = jsValueToGo(args[2])
			}

			var target interface{}
			resp, err := client.Post(context.Background(), path, params, body, &target)
			if err != nil {
				return nil, err
			}

			return responseToJS(resp), nil
		})
	}
}

func makePutFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		return promiseWrapper(func() (interface{}, error) {
			if len(args) < 1 {
				return nil, errors.New("path is required")
			}

			path := args[0].String()
			var params map[string]interface{}
			var body interface{}

			if len(args) >= 2 && !args[1].IsUndefined() && !args[1].IsNull() {
				params = jsObjectToMap(args[1])
			}

			if len(args) >= 3 && !args[2].IsUndefined() && !args[2].IsNull() {
				body = jsValueToGo(args[2])
			}

			var target interface{}
			resp, err := client.Put(context.Background(), path, params, body, &target)
			if err != nil {
				return nil, err
			}

			return responseToJS(resp), nil
		})
	}
}

func makePatchFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		return promiseWrapper(func() (interface{}, error) {
			if len(args) < 1 {
				return nil, errors.New("path is required")
			}

			path := args[0].String()
			var params map[string]interface{}
			var body interface{}

			if len(args) >= 2 && !args[1].IsUndefined() && !args[1].IsNull() {
				params = jsObjectToMap(args[1])
			}

			if len(args) >= 3 && !args[2].IsUndefined() && !args[2].IsNull() {
				body = jsValueToGo(args[2])
			}

			var target interface{}
			resp, err := client.Patch(context.Background(), path, params, body, &target)
			if err != nil {
				return nil, err
			}

			return responseToJS(resp), nil
		})
	}
}

func makeDeleteFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		return promiseWrapper(func() (interface{}, error) {
			if len(args) < 1 {
				return nil, errors.New("path is required")
			}

			path := args[0].String()
			var params map[string]interface{}

			if len(args) >= 2 && !args[1].IsUndefined() && !args[1].IsNull() {
				params = jsObjectToMap(args[1])
			}

			var target interface{}
			resp, err := client.Delete(context.Background(), path, params, &target)
			if err != nil {
				return nil, err
			}

			return responseToJS(resp), nil
		})
	}
}

func makeSetBaseURLFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return nil
		}
		client.SetBaseURL(args[0].String())
		return this
	}
}

func makeSetTimeoutFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return this
		}

		// JavaScript numbers can be passed as different types
		// Try to get as float first, then convert to int
		var timeout int
		jsVal := args[0]

		// Check the type and convert appropriately
		switch jsVal.Type() {
		case js.TypeNumber:
			// Get as float and convert to int (milliseconds)
			timeout = int(jsVal.Float())
		default:
			// If it's not a number, return without setting
			return this
		}

		client.SetTimeout(durationFromMillis(timeout))
		return this
	}
}

func makeSetHeaderFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			return nil
		}
		client.SetHeader(args[0].String(), args[1].String())
		return this
	}
}

func makeNewInstanceFunc(client *infrastructure.Client) func(js.Value, []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		newClient := client.NewInstance()
		return map[string]interface{}{
			"get":         js.FuncOf(makeGetFunc(newClient)),
			"post":        js.FuncOf(makePostFunc(newClient)),
			"put":         js.FuncOf(makePutFunc(newClient)),
			"patch":       js.FuncOf(makePatchFunc(newClient)),
			"delete":      js.FuncOf(makeDeleteFunc(newClient)),
			"setBaseURL":  js.FuncOf(makeSetBaseURLFunc(newClient)),
			"setTimeout":  js.FuncOf(makeSetTimeoutFunc(newClient)),
			"setHeader":   js.FuncOf(makeSetHeaderFunc(newClient)),
			"newInstance": js.FuncOf(makeNewInstanceFunc(newClient)),
		}
	}
}
