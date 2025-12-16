//go:build js && wasm
// +build js,wasm

package wasm

import (
	"encoding/json"
	"syscall/js"
	"time"

	"github.com/nikoskechris/gofetch/domain/models"
)

// promiseWrapper wraps a Go function in a JavaScript Promise.
func promiseWrapper(fn func() (interface{}, error)) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			result, err := fn()
			if err != nil {
				reject.Invoke(err.Error())
			} else {
				resolve.Invoke(result)
			}
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// jsObjectToMap converts a JavaScript object to a Go map.
func jsObjectToMap(obj js.Value) map[string]interface{} {
	result := make(map[string]interface{})

	if obj.Type() != js.TypeObject {
		return result
	}

	keys := js.Global().Get("Object").Call("keys", obj)
	for i := 0; i < keys.Length(); i++ {
		key := keys.Index(i).String()
		value := obj.Get(key)
		result[key] = jsValueToGo(value)
	}

	return result
}

// jsValueToGo converts a JavaScript value to a Go value.
func jsValueToGo(val js.Value) interface{} {
	switch val.Type() {
	case js.TypeBoolean:
		return val.Bool()
	case js.TypeNumber:
		return val.Float()
	case js.TypeString:
		return val.String()
	case js.TypeObject:
		if val.Get("constructor").Get("name").String() == "Array" {
			length := val.Length()
			result := make([]interface{}, length)
			for i := 0; i < length; i++ {
				result[i] = jsValueToGo(val.Index(i))
			}
			return result
		}
		return jsObjectToMap(val)
	case js.TypeNull, js.TypeUndefined:
		return nil
	default:
		return nil
	}
}

// responseToJS converts a Response to a JavaScript object.
func responseToJS(resp *models.Response) interface{} {
	// Convert headers to JS object
	headers := make(map[string]interface{})
	for key, values := range resp.Headers {
		if len(values) == 1 {
			headers[key] = values[0]
		} else {
			headers[key] = values
		}
	}

	// Convert data to JSON-compatible format
	var data interface{}
	if resp.Data != nil {
		// Marshal and unmarshal to ensure JSON compatibility
		jsonData, _ := json.Marshal(resp.Data)
		json.Unmarshal(jsonData, &data)
	}

	return map[string]interface{}{
		"statusCode": resp.StatusCode,
		"headers":    headers,
		"data":       data,
		"rawBody":    string(resp.RawBody),
	}
}

// durationFromMillis converts milliseconds to time.Duration.
func durationFromMillis(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}
