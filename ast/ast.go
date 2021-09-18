package ast

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/SmartBrave/Athena/ast/go/ast"
	"github.com/SmartBrave/Athena/ast/go/parser"
	"github.com/SmartBrave/Athena/ast/go/token"
	"github.com/SmartBrave/Athena/easyerrors"
)

type MODE int

const (
	STRICT MODE = iota
	COMPATIBLE
)

type AST struct {
	mode MODE
}

func NewAST(m MODE) (ast *AST) {
	return &AST{mode: m}
}

func (a *AST) Judge(variable map[string]interface{}, functions map[string]interface{}, rule string) (bool, error) {
	exprAst, err := parser.ParseExpr(rule)
	if err != nil {
		return false, err
	}

	ret, err := a.do(variable, functions, exprAst)
	if flag, ok := ret.(bool); err == nil && ok {
		return flag, nil
	}
	if err != nil {
		return false, err
	}
	return false, errors.New("the result of expression is not boolean")
}

var (
	badError = errors.New("bad expression")
)

func (a *AST) do(variable map[string]interface{}, functions map[string]interface{}, node ast.Node) (interface{}, error) {
	switch node.(type) {
	case *ast.BadExpr:
		return nil, badError
	case *ast.Ident:
		ident := node.(*ast.Ident)
		iName := ident.Name

		switch iName {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return iName, nil
		}
	case *ast.Variable:
		v := node.(*ast.Variable)
		if val, ok := variable[v.Name]; ok {
			return val, nil
		}
		return nil, badError
	case *ast.BasicLit:
		basicLit := node.(*ast.BasicLit)
		switch basicLit.Kind {
		case token.INT: //using int instead int64
			return strconv.Atoi(basicLit.Value)
		// case token.FLOAT:
		// return strconv.ParseFloat(basicLit.Value, 64)
		// case token.IMAG:
		// case token.CHAR:
		case token.STRING:
			return strconv.Unquote(basicLit.Value)
		default:
			return nil, badError
		}
	case *ast.ParenExpr:
		parenExpr := node.(*ast.ParenExpr)
		return a.do(variable, functions, parenExpr.X)
	case *ast.IndexExpr:
		indexExpr := node.(*ast.IndexExpr)
		x, err1 := a.do(variable, functions, indexExpr.X)
		index, err2 := a.do(variable, functions, indexExpr.Index)
		if err := easyerrors.HandleMultiError(easyerrors.Simple(), err1, err2); err != nil {
			return nil, badError
		}

		xv := reflect.ValueOf(x)
		switch reflect.TypeOf(x).Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			if iIndex, ok := index.(int); !ok {
				return nil, badError
			} else {
				return xv.Index(iIndex).Interface(), nil
			}
		case reflect.Map:
			return xv.MapIndex(reflect.ValueOf(index)).Interface(), nil
		default:
			return nil, badError
		}
	case *ast.UnaryExpr:
		unaryExpr := node.(*ast.UnaryExpr)
		x, err := a.do(variable, functions, unaryExpr.X)
		if err != nil {
			return nil, badError
		}
		switch unaryExpr.Op {
		case token.NOT: //!
			xv, okx := x.(bool)
			if !okx {
				return nil, badError
			}
			return !xv, nil
		}
	case *ast.BinaryExpr:
		binaryExpr := node.(*ast.BinaryExpr)
		x, err1 := a.do(variable, functions, binaryExpr.X)
		y, err2 := a.do(variable, functions, binaryExpr.Y)
		if err := easyerrors.HandleMultiError(easyerrors.Simple(), err1, err2); err != nil {
			return nil, badError
		}
		switch binaryExpr.Op {
		case token.EQL: //==
			return reflect.DeepEqual(x, y), nil
		case token.NEQ: //!=
			return !reflect.DeepEqual(x, y), nil
		default:
			//continue
		}

		if reflect.TypeOf(x).Kind() != reflect.TypeOf(y).Kind() {
			return nil, badError
		}

		switch x.(type) {
		case int:
			xv, _ := x.(int)
			yv, _ := y.(int)
			switch binaryExpr.Op {
			case token.REM: //%
				return xv % yv, nil
			case token.LSS: //<
				return xv < yv, nil
			case token.GTR: //>
				return xv > yv, nil
			case token.LEQ: //<=
				return xv <= yv, nil
			case token.GEQ: //>=
				return xv >= yv, nil
			case token.AND: //&
				switch a.mode {
				case STRICT:
					return xv & yv, nil
				case COMPATIBLE:
					return xv != 0 && yv != 0, nil
				default:
					return nil, badError
				}
			case token.OR: //|
				switch a.mode {
				case STRICT:
					return xv | yv, nil
				case COMPATIBLE:
					return xv != 0 || yv != 0, nil
				default:
					return nil, badError
				}
			default:
				return nil, badError
			}
		case bool:
			xv, _ := x.(bool)
			yv, _ := y.(bool)
			switch binaryExpr.Op {
			case token.LAND: //&&
				return xv && yv, nil
			case token.LOR: //||
				return xv || yv, nil
			case token.AND: //&
				switch a.mode {
				case STRICT:
					return nil, badError
				case COMPATIBLE:
					return xv && yv, nil
				default:
					return nil, badError
				}
			case token.OR: //|
				switch a.mode {
				case STRICT:
					return nil, badError
				case COMPATIBLE:
					return xv || yv, nil
				default:
					return nil, badError
				}
			default:
				return nil, badError
			}
		default:
			return nil, badError
		}
	case *ast.SliceExpr:
		sliceExpr := node.(*ast.SliceExpr)
		if sliceExpr.Slice3 {
			return nil, badError
		}
		x, err := a.do(variable, functions, sliceExpr.X)
		if err != nil {
			return nil, badError
		}
		if t := reflect.TypeOf(x).Kind(); t != reflect.Array && t != reflect.Slice && t != reflect.String {
			return nil, badError
		}
		xLen := reflect.ValueOf(x).Len()

		lowv := 0
		if sliceExpr.Low != nil {
			low, err := a.do(variable, functions, sliceExpr.Low)
			if err != nil {
				return nil, badError
			}
			oklow := false
			lowv, oklow = low.(int)
			if !oklow {
				return nil, badError
			}
		}

		highv := xLen
		if sliceExpr.High != nil {
			high, err := a.do(variable, functions, sliceExpr.High)
			if err != nil {
				return nil, badError
			}
			okhigh := false
			highv, okhigh = high.(int)
			if !okhigh {
				return nil, badError
			}
		}

		if lowv >= highv || lowv < 0 || highv > xLen {
			return nil, badError
		}

		return reflect.ValueOf(x).Slice(lowv, highv).Interface(), nil
	case *ast.CallExpr:
		callExpr := node.(*ast.CallExpr)
		fun, err := a.do(variable, functions, callExpr.Fun)
		if err != nil {
			return nil, err
		}
		funv, ok := fun.(string)
		if !ok {
			return nil, badError
		}
		if f, ok := functions[funv]; ok {
			in := []reflect.Value{}
			for i := 0; i < len(callExpr.Args); i++ {
				arg, err := a.do(variable, functions, callExpr.Args[i])
				if err != nil {
					return nil, badError
				}
				in = append(in, reflect.ValueOf(arg))
			}
			out := reflect.ValueOf(f).Call(in)
			if len(out) == 0 {
				return nil, nil
			}
			return out[0].Interface(), nil //NOTE: ignore other values
		} else {
			switch strings.ToLower(funv) {
			case "contains":
				if len(callExpr.Args) != 2 {
					return nil, badError
				}
				arg0, err0 := a.do(variable, functions, callExpr.Args[0])
				arg1, err1 := a.do(variable, functions, callExpr.Args[1])
				if err := easyerrors.HandleMultiError(easyerrors.Simple(), err0, err1); err != nil {
					return nil, badError
				}

				arg0v, ok0 := arg0.(string)
				arg1v, ok1 := arg1.(string)
				if !ok0 || !ok1 {
					return nil, badError
				}

				return strings.Contains(arg0v, arg1v), nil
			case "mod": //mod with string
				if len(callExpr.Args) != 2 {
					return nil, badError
				}
				arg0, err0 := a.do(variable, functions, callExpr.Args[0])
				arg1, err1 := a.do(variable, functions, callExpr.Args[1])
				if err := easyerrors.HandleMultiError(easyerrors.Simple(), err0, err1); err != nil {
					return nil, badError
				}

				arg0v, ok0 := arg0.(string)
				arg1v, ok1 := arg1.(int)
				if !ok0 || !ok1 {
					return nil, badError
				}
				arg0vi, err := strconv.Atoi(arg0v)
				if err != nil {
					return nil, badError
				}

				return arg0vi % arg1v, nil
			default:
				return nil, badError
			}
		}
	case *ast.AssignStmt:
		assignStmt := node.(*ast.AssignStmt)
		switch assignStmt.Tok {
		case token.ASSIGN:
			switch a.mode {
			case STRICT:
				//XXX: support assign operator?
				return nil, badError
			case COMPATIBLE:
				if len(assignStmt.Lhs) != len(assignStmt.Rhs) {
					return nil, badError
				}
				return reflect.DeepEqual(assignStmt.Lhs, assignStmt.Rhs), nil
			default:
				return nil, badError
			}
		}
	default:
		return nil, badError
	}

	return nil, badError
}
