package runtime

import "github.com/mykytaserdiuk/serc/ast"

type Runtime interface {
	Call(fn ast.Value, args []ast.Value) ast.Value
}
