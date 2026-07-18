package buildin

import (
	"net/http"

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
	RT = rt
	return funcs
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

	ret := RT.Call(response, []ast.Value{
		ast.GetStringValue("hello"),
	})

	http.HandleFunc(endpoint.Data.(string), func(w http.ResponseWriter, r *http.Request) {
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
