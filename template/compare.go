package template

import (
	"fmt"
	"reflect"
	"strings"
)

// ResolveCondition evaluates a condition string (e.g., "x > 5").
func ResolveCondition(ctx *RenderCtx, condition string) (val bool, errs []error) {
	values := SplitRespectQuotes(condition)
	vlen := len(values)
	if vlen == 0 {
		return false, []error{fmt.Errorf("no conditions within %s statement", "if")}
	}

	// first value
	if len(errs) > 0 {
		return false, errs
	}
	//
	fv := values[0]
	if strings.HasPrefix(fv, ".") && vlen == 1 {
		value, verr := ResolveVariable(ctx, fv)
		if verr != nil {
			errs = append(errs, verr)
		}
		if value != nil && value != false {
			return true, errs
		} else {
			return false, errs
		}
	}
	//

	//
	if vlen == 3 {
		left := values[0]
		op := values[1]
		right := values[2]
		//
		lhsVal, lhsErr := ResolveValue(ctx, left)
		rhsVal, rhsErr := ResolveValue(ctx, right)
		if lhsErr != nil {
			errs = append(errs, fmt.Errorf("left operand error: %v", lhsErr))
		}
		if rhsErr != nil {
			errs = append(errs, fmt.Errorf("right operand error: %v", rhsErr))
		}

		result, compareErr := compareValues(op, lhsVal, rhsVal)
		if compareErr != nil {
			errs = append(errs, fmt.Errorf("right operand error: %v", rhsErr))
		}
		return result, errs
	}
	//
	// Handle truthy check
	value, err := ResolveTruthy(ctx, condition)
	if err != nil {
		errs = append(errs, err)
		return false, errs
	}
	return value, nil
}

func ResolveTruthy(ctx *RenderCtx, expr string) (bool, error) {
	value, err := ResolveValue(ctx, expr)
	if err != nil {
		return false, err
	}
	return isTruthy(value), nil
}

func isTruthy(val any) bool {
	if val == nil {
		return false
	}

	switch v := val.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0.0
	case string:
		return v != ""
	case []any:
		return len(v) > 0
	case map[string]any:
		return len(v) > 0
	default:
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
			return !rv.IsNil()
		case reflect.Struct:
			return true
		default:
			return !rv.IsZero()
		}
	}
}

func compareValues(op string, lhs, rhs any) (bool, error) {
	switch op {
	case "==":
		return equal(lhs, rhs), nil
	case "!=":
		return !equal(lhs, rhs), nil
	case ">", "<", ">=", "<=":
		return compareNumbers(op, lhs, rhs)
	default:
		return false, fmt.Errorf("unsupported operator %q", op)
	}
}

func equal(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func compareNumbers(op string, a, b any) (bool, error) {
	af, err := toFloat(a)
	if err != nil {
		return false, fmt.Errorf("left operand is not a number: %v", a)
	}
	bf, err := toFloat(b)
	if err != nil {
		return false, fmt.Errorf("right operand is not a number: %v", b)
	}

	switch op {
	case ">":
		return af > bf, nil
	case "<":
		return af < bf, nil
	case ">=":
		return af >= bf, nil
	case "<=":
		return af <= bf, nil
	default:
		return false, fmt.Errorf("unsupported operator %q for numeric comparison", op)
	}
}

func toFloat(val any) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("type %T is not a number", val)
	}
}
