package optparse

import "fmt"
import "strconv"

type Type interface {
    parseArg(opt string, arg []string) interface{};
    storeDefault(dest, def interface{});
}

// StringType
type StringType struct {
}

func (s *StringType) parseArg(opt string, arg []string) interface{} {
    return arg[0];
}

func (s *StringType) storeDefault(dest, def interface{}) {
    *dest.(*string) = def.(string);
}

func String(a...) *string {
    dest := new(string);
    StringVar(dest, a);
    return dest;
}

func StringVar(dest *string, a...) {
    typ := new(StringType);
    createOption(a, dest, typ, Store);
}

// incrementable
type incrementable interface {
    increment(interface{});
}

// IntType
type IntType struct {
}

func (it *IntType) parseArg(opt string, arg []string) interface{} {
    i, ok := strconv.Atoi(arg[0]);
    if ok != nil {
        Error(opt, fmt.Sprintf("'%s' is not an integer", arg[0]));
    }
    return i;
}

func (it *IntType) storeDefault(dest, def interface{}) {
    ptr, ok := dest.(*int);
    if !ok {
        Error("..", "blargh");
    }
    *ptr = def.(int);
}

func (it *IntType) increment(dest interface{}) {
    ptr := dest.(*int);
    *ptr++;
}

func Int(a...) *int {
    dest := new(int);
    IntVar(dest, a);
    return dest;
}

func IntVar(dest *int, a...) {
    typ := new(IntType);
    createOption(a, dest, typ, Store);
}

// CallbackType
type CallbackType struct {
    fn interface{};
}

// This is not actually used. The CallbackType can only be used with the
// CallbackAction.
func (cb *CallbackType) parseArg(opt string, arg []string) interface{} {
    return nil;
}

// Ditto.
func (cb *CallbackType) storeDefault(dest, def interface{}) {
}

func Callback(a...) {
    typ := new(CallbackType);
    createOption(a, nil, typ, nil);
}

// array
type array interface {
    append(dest interface{}, val interface{});
}

// StringArrayType
type StringArrayType struct {
}

func (sa *StringArrayType) parseArg(opt string, arg []string) interface{} {
    if len(arg) == 1 {
        return arg[0];
    } else {
        return arg;
    }
    return nil;
}

func (s *StringArrayType) storeDefault(dest, def interface{}) {
    *dest.(*[]string) = def.([]string);
}

func (sa *StringArrayType) append(dest interface{}, val interface{}) {
    a := dest.(*[]string);
    *a = appendString(*a, val.(string));
}

func StringArray(a...) *[]string {
    dest := new([]string);
    *dest = make([]string, 0, 5);
    StringArrayVar(dest, a);
    return dest;
}

func StringArrayVar(dest *[]string, a...) {
    typ := new(StringArrayType);
    createOption(a, dest, typ, Append);
}

// IntArrayType
type IntArrayType struct {
}

func (ia *IntArrayType) parseArg(opt string, arg []string) interface{} {
    if len(arg) == 1 {
        i, ok := strconv.Atoi(arg[0]);
        if ok != nil {
            Error(opt, fmt.Sprintf("'%s' is not an integer", arg[0]));
        }
        return i;
    } else {
        ret := make([]int, len(arg));
        for i, str := range arg {
            x, ok := strconv.Atoi(str);
            if ok != nil {
                Error(opt, fmt.Sprintf("'%s' is not an integer", str));
            }
            ret[i] = x
        }
        return ret;
    }
    return nil;
}

func (s *IntArrayType) storeDefault(dest, def interface{}) {
    *dest.(*[]int) = def.([]int);
}

func (ia *IntArrayType) append(dest interface{}, val interface{}) {
    i := dest.(*[]int);
    *i = appendInt(*i, val.(int));
}

func IntArray(a...) *[]int {
    dest := new([]int);
    *dest = make([]int, 0, 5);
    IntArrayVar(dest, a);
    return dest;
}

func IntArrayVar(dest *[]int, a...) {
    typ := new(IntArrayType);
    createOption(a, dest, typ, Append);
}

// StringArrayArray
type StringArrayArrayType struct {
}

func (sa *StringArrayArrayType) parseArg(opt string, arg []string) interface{} {
    return arg;
}

func (sa *StringArrayArrayType) storeDefault(dest, def interface{}) {
    *dest.(*[][]string) = def.([][]string);
}

func (sa *StringArrayArrayType) append(dest, val interface{}) {
    a := dest.(*[][]string);
    *a = appendStringArray(*a, val.([]string));
}

func StringArrayArray(a...) *[][]string {
    dest := new([][]string);
    *dest = make([][]string, 0, 5);
    StringArrayArrayVar(dest, a);
    return dest;
}

func StringArrayArrayVar(dest *[][]string, a...) {
    typ := new(StringArrayArrayType);
    createOption(a, dest, typ, Append);
}
