package buildin

import (
	"net/http"

	"github.com/mykytaserdiuk/serc/ast"
)

func LoadHttp() map[string]ast.BuiltinFunc {
	funcs := make(map[string]ast.BuiltinFunc)
	funcs["serve"] = HttpServe
	funcs["get"] = HttpGet
	return funcs
}

func HttpGet(args []ast.Value) ast.FuncResult {
	endpoint := args[0]
	if endpoint.Type != ast.StringValue {
		panic("RUNTIME ERROR: expected string value as first argument")
	}
	response := args[1]
	if response.Type != ast.StringValue {
		panic("RUNTIME ERROR: expected string value as second argument")
	}

	http.HandleFunc(endpoint.Data.(string), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response.Data.(string)))
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
