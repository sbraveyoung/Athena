package ast

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"

	"github.com/SmartBrave/utils/easyerrors"
)

type MODE int

const (
	STRICT MODE = iota //BUG: do not support functions
	COMPATIBLE
)

type AST struct {
	mode MODE
}

func NewAST(m MODE) (ast *AST) {
	return &AST{mode: m}
}

func (a *AST) Judge(args map[string]interface{}, ops map[string]func(interface{}) bool, rule string) (bool, error) {
	exprAst, err := parser.ParseExpr(rule)
	if err != nil {
		return false, err
	}

	ret, err := a.do(args, ops, exprAst)
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

func (a *AST) do(args map[string]interface{}, ops map[string]func(interface{}) bool, node ast.Node) (interface{}, error) {
	switch node.(type) {
	case *ast.BadExpr:
		return nil, badError
	case *ast.Ident:
		//XXX: could do better!
		ident := node.(*ast.Ident)
		switch ident.Name {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			//!!!NOTE!!!: ensure that the keys in args and ops do not overlap, and don't conflict with predefined functions.
			if val, ok := args[ident.Name]; ok {
				return val, nil
			}
			if val, ok := ops[ident.Name]; ok {
				return val, nil
			}
			switch a.mode {
			case STRICT:
				return nil, badError
			case COMPATIBLE:
				//compatible with app=\"media_std\", app=media_std and predefined functions.
				return ident.Name, nil
			default:
				//XXX: do nothing?
			}
		}
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
		return a.do(args, ops, parenExpr.X)
	case *ast.IndexExpr: //NOTE: index with slice is not supported!
		indexExpr := node.(*ast.IndexExpr)
		x, err1 := a.do(args, ops, indexExpr.X)
		index, err2 := a.do(args, ops, indexExpr.Index)
		if err := easyerrors.HandleMultiError(easyerrors.Simple(), err1, err2); err != nil {
			return nil, badError
		}

		xv, okx := x.(func(interface{}) bool)
		if !okx {
			return nil, badError
		}
		return xv(index), nil
	case *ast.UnaryExpr:
		unaryExpr := node.(*ast.UnaryExpr)
		x, err := a.do(args, ops, unaryExpr.X)
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
		x, err1 := a.do(args, ops, binaryExpr.X)
		y, err2 := a.do(args, ops, binaryExpr.Y)
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
			xv, okx := x.(int)
			yv, oky := y.(int)
			if !okx || !oky {
				return nil, badError
			}
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
			case token.AND: //&, TODO
			case token.OR: //|, TODO
			default:
				return nil, badError
			}
		case bool:
			xv, okx := x.(bool)
			yv, oky := y.(bool)
			if !okx || !oky {
				return nil, badError
			}
			switch binaryExpr.Op {
			case token.LAND: //&&
				return xv && yv, nil
			case token.LOR: //||
				return xv || yv, nil
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

		x, err := a.do(args, ops, sliceExpr.X)
		if err != nil {
			return nil, badError
		}
		xv, okx := x.(string) //NOTE: only support string, eg: cv[8:]=="Iphone" while cv=="IK7.8.9_Iphone"
		if !okx {
			return nil, badError
		}

		lowv := 0
		if sliceExpr.Low != nil {
			low, err := a.do(args, ops, sliceExpr.Low)
			if err != nil {
				return nil, badError
			}
			oklow := false
			lowv, oklow = low.(int)
			if !oklow {
				return nil, badError
			}
		}

		highv := len(xv)
		if sliceExpr.High != nil {
			high, err := a.do(args, ops, sliceExpr.High)
			if err != nil {
				return nil, badError
			}
			okhigh := false
			highv, okhigh = high.(int)
			if !okhigh {
				return nil, badError
			}
		}

		if lowv >= highv || lowv < 0 || highv > len(xv) {
			return nil, badError
		}
		return xv[lowv:highv], nil
	case *ast.CallExpr:
		callExpr := node.(*ast.CallExpr)
		fun, err := a.do(args, ops, callExpr.Fun)
		if err != nil {
			return nil, err
		}
		funv, ok := fun.(string)
		if !ok {
			return nil, badError
		}

		switch strings.ToLower(funv) {
		case "contains":
			if len(callExpr.Args) != 2 {
				return nil, badError
			}
			arg0, err0 := a.do(args, ops, callExpr.Args[0])
			arg1, err1 := a.do(args, ops, callExpr.Args[1])
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
			arg0, err0 := a.do(args, ops, callExpr.Args[0])
			arg1, err1 := a.do(args, ops, callExpr.Args[1])
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
	case *ast.AssignStmt:
		assignStmt := node.(*ast.AssignStmt)
		switch assignStmt.Tok {
		case token.ASSIGN:
			switch a.mode {
			case STRICT:
				//TODO
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
		// case *ast.SelectorExpr: //not support
		// case *ast.CompositeLit: //not support
		// case *ast.Ellipsis: //not support
		// case *ast.TypeAssertExpr: //not support
		// case *ast.FuncLit: //not support
		// case *ast.StarExpr: //not support
		// case *ast.KeyValueExpr://not support
	}

	return nil, badError
}
