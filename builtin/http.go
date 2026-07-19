package builtin

import (
	"encoding/json"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/mykytaserdiuk/serc/ast"
	"github.com/mykytaserdiuk/serc/runtime"
)

var (
	RT runtime.Runtime
)

func LoadHttp(rt runtime.Runtime) map[string]ast.BuiltinFunc {
	funcs := make(map[string]ast.BuiltinFunc)
	funcs["serve"] = HttpServe
	funcs["get"] = HttpGet
	funcs["post"] = HttpPost
	funcs["handle"] = HttpHandle
	RT = rt
	return funcs
}

func HttpHandle(args []ast.Value) ast.FuncResult {
	endpoint := args[0]
	if endpoint.Type != ast.StringValue {
		panic("RUNTIME ERROR: expected string value as first argument")
	}
	callback := args[1]
	if callback.Type != ast.FuncValue {
		panic("RUNTIME ERROR: expected func value as second argument")
	}

	http.HandleFunc(endpoint.Data.(string), func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		ret := RT.Call(callback, []ast.Value{
			ast.GetStringValue(r.Method),
			ast.GetNativeObjectValue("body", body),
			ast.GetNativeObjectValue("request", r),
		})

		obj := ret.Data.(ast.Object)

		var returnMap map[string]ast.Value
		err := mapstructure.Decode(obj.Fields, &returnMap)
		if err != nil {
			// todo write error
		}
		b, err := json.Marshal(returnMap["body"].Data)
		if err != nil {
			// todo write error
		}
		w.WriteHeader(returnMap["status"].Data.(int))
		w.Write(b)
	})

	return ast.FuncResult{
		Value: ast.NumberLiteral{Value: 0},
	}
}

func HttpPost(args []ast.Value) ast.FuncResult {
	endpoint := args[0]
	if endpoint.Type != ast.StringValue {
		panic("RUNTIME ERROR: expected string value as first argument")
	}
	response := args[1]
	if response.Type != ast.FuncValue {
		panic("RUNTIME ERROR: expected func value as second argument")
	}

	http.HandleFunc(endpoint.Data.(string), func(w http.ResponseWriter, r *http.Request) {
		ret := RT.Call(response, []ast.Value{
			ast.GetNativeObjectValue("request", r),
		})
		w.Write([]byte(ret.Data.(string)))
	})

	return ast.FuncResult{
		Value: ast.NumberLiteral{Value: 0},
	}
}

func HttpGet(args []ast.Value) ast.FuncResult {
	endpoint := args[0]
	if endpoint.Type != ast.StringValue {
		panic("RUNTIME ERROR: expected string value as first argument")
	}
	response := args[1]
	if response.Type != ast.FuncValue {
		panic("RUNTIME ERROR: expected func value as second argument")
	}

	http.HandleFunc(endpoint.Data.(string), func(w http.ResponseWriter, r *http.Request) {
		ret := RT.Call(response, []ast.Value{
			ast.GetStringValue(r.URL.Query().Get("q")),
		})
		w.Write([]byte(ret.Data.(string)))
	})

	return ast.FuncResult{
		Value: ast.NumberLiteral{Value: 0},
	}
}

func HttpServe(args []ast.Value) ast.FuncResult {
	port := args[0]
	if port.Type != ast.StringValue {
		panic("RUNTIME ERROR: unsoportet 'http.serve' type")
	}
	err := http.ListenAndServe(port.Data.(string), nil)
	if err != nil {
		// TODO: change to error
		return ast.FuncResult{
			Value: ast.StringLiteral{
				Value: err.Error(),
			},
		}
	}
	return ast.FuncResult{
		Value: ast.NumberLiteral{Value: 0},
	}
}
