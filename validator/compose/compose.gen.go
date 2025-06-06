// Code generated by compose-gen. DO NOT EDIT.
package compose

import (
	"github.com/katydid/validator-go-jsonschema/validator/ast"
	"github.com/katydid/validator-go-jsonschema/validator/funcs"
	"github.com/katydid/validator-go-jsonschema/validator/types"
)

func composeDouble(expr *ast.Expr) (funcs.Double, error) {
	f, err := prep(expr, types.SINGLE_DOUBLE)
	if err != nil {
		return nil, err
	}
	if expr.Terminal != nil {
		if expr.GetTerminal().Variable != nil {
			return funcs.DoubleVar(), nil
		} else {
			return funcs.DoubleConst(expr.GetTerminal().GetDoubleValue()), nil
		}
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewDouble(values...)
}

func composeDoubles(expr *ast.Expr) (funcs.Doubles, error) {
	f, err := prep(expr, types.LIST_DOUBLE)
	if err != nil {
		return nil, err
	}
	if expr.List != nil {
		vs, err := newValues(expr.GetList().GetElems())
		if err != nil {
			return nil, err
		}
		bs := make([]funcs.Double, len(vs))
		var ok bool
		for i := range vs {
			bs[i], ok = vs[i].(funcs.Double)
			if !ok {
				return nil, &errExpected{types.SINGLE_DOUBLE.String(), expr.String()}
			}
		}
		return funcs.NewListOfDouble(bs), nil
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewDoubles(values...)
}

func composeInt(expr *ast.Expr) (funcs.Int, error) {
	f, err := prep(expr, types.SINGLE_INT)
	if err != nil {
		return nil, err
	}
	if expr.Terminal != nil {
		if expr.GetTerminal().Variable != nil {
			return funcs.IntVar(), nil
		} else {
			return funcs.IntConst(expr.GetTerminal().GetIntValue()), nil
		}
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewInt(values...)
}

func composeInts(expr *ast.Expr) (funcs.Ints, error) {
	f, err := prep(expr, types.LIST_INT)
	if err != nil {
		return nil, err
	}
	if expr.List != nil {
		vs, err := newValues(expr.GetList().GetElems())
		if err != nil {
			return nil, err
		}
		bs := make([]funcs.Int, len(vs))
		var ok bool
		for i := range vs {
			bs[i], ok = vs[i].(funcs.Int)
			if !ok {
				return nil, &errExpected{types.SINGLE_INT.String(), expr.String()}
			}
		}
		return funcs.NewListOfInt(bs), nil
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewInts(values...)
}

func composeUint(expr *ast.Expr) (funcs.Uint, error) {
	f, err := prep(expr, types.SINGLE_UINT)
	if err != nil {
		return nil, err
	}
	if expr.Terminal != nil {
		if expr.GetTerminal().Variable != nil {
			return funcs.UintVar(), nil
		} else {
			return funcs.UintConst(expr.GetTerminal().GetUintValue()), nil
		}
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewUint(values...)
}

func composeUints(expr *ast.Expr) (funcs.Uints, error) {
	f, err := prep(expr, types.LIST_UINT)
	if err != nil {
		return nil, err
	}
	if expr.List != nil {
		vs, err := newValues(expr.GetList().GetElems())
		if err != nil {
			return nil, err
		}
		bs := make([]funcs.Uint, len(vs))
		var ok bool
		for i := range vs {
			bs[i], ok = vs[i].(funcs.Uint)
			if !ok {
				return nil, &errExpected{types.SINGLE_UINT.String(), expr.String()}
			}
		}
		return funcs.NewListOfUint(bs), nil
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewUints(values...)
}

func composeBool(expr *ast.Expr) (funcs.Bool, error) {
	f, err := prep(expr, types.SINGLE_BOOL)
	if err != nil {
		return nil, err
	}
	if expr.Terminal != nil {
		if expr.GetTerminal().Variable != nil {
			return funcs.BoolVar(), nil
		} else {
			return funcs.BoolConst(expr.GetTerminal().GetBoolValue()), nil
		}
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewBool(values...)
}

func composeBools(expr *ast.Expr) (funcs.Bools, error) {
	f, err := prep(expr, types.LIST_BOOL)
	if err != nil {
		return nil, err
	}
	if expr.List != nil {
		vs, err := newValues(expr.GetList().GetElems())
		if err != nil {
			return nil, err
		}
		bs := make([]funcs.Bool, len(vs))
		var ok bool
		for i := range vs {
			bs[i], ok = vs[i].(funcs.Bool)
			if !ok {
				return nil, &errExpected{types.SINGLE_BOOL.String(), expr.String()}
			}
		}
		return funcs.NewListOfBool(bs), nil
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewBools(values...)
}

func composeString(expr *ast.Expr) (funcs.String, error) {
	f, err := prep(expr, types.SINGLE_STRING)
	if err != nil {
		return nil, err
	}
	if expr.Terminal != nil {
		if expr.GetTerminal().Variable != nil {
			return funcs.StringVar(), nil
		} else {
			return funcs.StringConst(expr.GetTerminal().GetStringValue()), nil
		}
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewString(values...)
}

func composeStrings(expr *ast.Expr) (funcs.Strings, error) {
	f, err := prep(expr, types.LIST_STRING)
	if err != nil {
		return nil, err
	}
	if expr.List != nil {
		vs, err := newValues(expr.GetList().GetElems())
		if err != nil {
			return nil, err
		}
		bs := make([]funcs.String, len(vs))
		var ok bool
		for i := range vs {
			bs[i], ok = vs[i].(funcs.String)
			if !ok {
				return nil, &errExpected{types.SINGLE_STRING.String(), expr.String()}
			}
		}
		return funcs.NewListOfString(bs), nil
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewStrings(values...)
}

func composeBytes(expr *ast.Expr) (funcs.Bytes, error) {
	f, err := prep(expr, types.SINGLE_BYTES)
	if err != nil {
		return nil, err
	}
	if expr.Terminal != nil {
		if expr.GetTerminal().Variable != nil {
			return funcs.BytesVar(), nil
		} else {
			return funcs.BytesConst(expr.GetTerminal().GetBytesValue()), nil
		}
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewBytes(values...)
}

func composeListOfBytes(expr *ast.Expr) (funcs.ListOfBytes, error) {
	f, err := prep(expr, types.LIST_BYTES)
	if err != nil {
		return nil, err
	}
	if expr.List != nil {
		vs, err := newValues(expr.GetList().GetElems())
		if err != nil {
			return nil, err
		}
		bs := make([]funcs.Bytes, len(vs))
		var ok bool
		for i := range vs {
			bs[i], ok = vs[i].(funcs.Bytes)
			if !ok {
				return nil, &errExpected{types.SINGLE_BYTES.String(), expr.String()}
			}
		}
		return funcs.NewListOfBytes(bs), nil
	}
	values, err := newValues(expr.GetFunction().GetParams())
	if err != nil {
		return nil, err
	}
	return f.NewListOfBytes(values...)
}
