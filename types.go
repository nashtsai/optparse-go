/*
Copyright 2009 Kirk McDonald

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/
package optparse

import "fmt"
import "os"
import "reflect"
import "strconv"

/*
type Type interface {
    parseArg(opt string, arg []string) interface{};
    storeDefault(dest, def interface{});
    validAction(action *Action, nargs int) bool;
    getOption() *option;
}
*/

// BoolType
type BoolType struct {
    option;
}

func (b *BoolType) storeDefault(dest, def interface{}) {
    *dest.(*bool) = def.(bool);
}

func (b *BoolType) validAction(action *Action, nargs int) bool {
    switch action {
    case StoreTrue, StoreFalse:
        return true;
    }
    return false;
}

func (op *OptionParser) Bool(a...) *bool {
    dest := new(bool);
    op.BoolVar(dest, a);
    return dest;
}

func (op *OptionParser) BoolVar(dest *bool, a...) {
    typ := new(BoolType);
    op.createOption(a, dest, typ, StoreTrue);
}

// StringType
type StringType struct {
    option;
}

func (s *StringType) parseArg(opt string, arg []string) interface{} {
    return arg[0];
}

func (s *StringType) storeDefault(dest, def interface{}) {
    *dest.(*string) = def.(string);
}

func (s *StringType) validAction(action *Action, nargs int) bool {
    switch action {
    case StoreConst:
        return true;
    case Store, Append:
        if nargs == 1 {
            return true;
        }
    }
    return false
}

func (op *OptionParser) String(a...) *string {
    dest := new(string);
    op.StringVar(dest, a);
    return dest;
}

func (op *OptionParser) StringVar(dest *string, a...) {
    typ := new(StringType);
    op.createOption(a, dest, typ, Store);
}

// incrementable
type incrementable interface {
    increment(interface{});
}

// IntType
type IntType struct {
    option;
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

func (it *IntType) validAction(action *Action, nargs int) bool {
    switch action {
    case StoreConst, Count:
        return true;
    case Store, Append:
        if nargs == 1 {
            return true;
        }
    }
    return false;
}
func (it *IntType) increment(dest interface{}) {
    ptr := dest.(*int);
    *ptr++;
}

func (op *OptionParser) Int(a...) *int {
    dest := new(int);
    op.IntVar(dest, a);
    return dest;
}

func (op *OptionParser) IntVar(dest *int, a...) {
    typ := new(IntType);
    op.createOption(a, dest, typ, Store);
}

// CallbackType
type CallbackType struct {
    option;
    fn interface{};
}

// Ditto.
func (cb *CallbackType) storeDefault(dest, def interface{}) {
}

func (cb *CallbackType) validAction(action *Action, nargs int) bool {
    // Each callback uses a different Action object, but they all have the
    // same function.
    return action == nil;
}

func (cb *CallbackType) hasArgs() bool {
    return cb.nargs > 0;
}

func (cb *CallbackType) performAction(optStr string, arg []string) {
    fn := reflect.NewValue(cb.fn).(*reflect.FuncValue);
    fnType := fn.Type().(*reflect.FuncType);
    values := make([]reflect.Value, len(arg));
    for i := 0; i < len(arg); i++ {
        switch v := fnType.In(i).(type) {
        case *reflect.StringType:
            values[i] = reflect.NewValue(arg[i]);
        case *reflect.IntType:
            x, ok := strconv.Atoi(arg[i]);
            if ok != nil {
                Error(optStr, fmt.Sprintf("'%s' is not an integer", arg[i]));
            }
            values[i] = reflect.NewValue(x);
        }
    }
    fn.Call(values);
}

func (op *OptionParser) Callback(a...) {
    typ := new(CallbackType);
    op.createOption(a, nil, typ, nil);
}

// array
type array interface {
    append(dest interface{}, val interface{});
}

// StringArrayType
type StringArrayType struct {
    option;
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

func (sa *StringArrayType) validAction(action *Action, nargs int) bool {
    switch action {
    case StoreConst:
        return true;
    case Store:
        if nargs > 1 {
            return true;
        }
    case Append:
        if nargs == 1 {
            return true;
        }
    }
    return false;
}

func (sa *StringArrayType) append(dest interface{}, val interface{}) {
    a := dest.(*[]string);
    *a = appendString(*a, val.(string));
}

func (op *OptionParser) StringArray(a...) *[]string {
    dest := new([]string);
    *dest = make([]string, 0, 5);
    op.StringArrayVar(dest, a);
    return dest;
}

func (op *OptionParser) StringArrayVar(dest *[]string, a...) {
    typ := new(StringArrayType);
    op.createOption(a, dest, typ, Append);
}

// IntArrayType
type IntArrayType struct {
    option;
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

func (ia *IntArrayType) validAction(action *Action, nargs int) bool {
    switch action {
    case StoreConst:
        return true;
    case Store:
        if nargs > 1 {
            return true;
        }
    case Append:
        if nargs == 1 {
            return true;
        }
    }
    return false;
}

func (ia *IntArrayType) append(dest interface{}, val interface{}) {
    i := dest.(*[]int);
    *i = appendInt(*i, val.(int));
}

func (op *OptionParser) IntArray(a...) *[]int {
    dest := new([]int);
    *dest = make([]int, 0, 5);
    op.IntArrayVar(dest, a);
    return dest;
}

func (op *OptionParser) IntArrayVar(dest *[]int, a...) {
    typ := new(IntArrayType);
    op.createOption(a, dest, typ, Append);
}

// StringArrayArray
type StringArrayArrayType struct {
    option;
}

func (sa *StringArrayArrayType) parseArg(opt string, arg []string) interface{} {
    return arg;
}

func (sa *StringArrayArrayType) storeDefault(dest, def interface{}) {
    *dest.(*[][]string) = def.([][]string);
}

func (sa *StringArrayArrayType) validAction(action *Action, nargs int) bool {
    switch action {
    case StoreConst:
        return true;
    case Append:
        if nargs > 1 {
            return true;
        }
    }
    return false;
}

func (sa *StringArrayArrayType) append(dest, val interface{}) {
    a := dest.(*[][]string);
    *a = appendStringArray(*a, val.([]string));
}

func (op *OptionParser) StringArrayArray(a...) *[][]string {
    dest := new([][]string);
    *dest = make([][]string, 0, 5);
    op.StringArrayArrayVar(dest, a);
    return dest;
}

func (op *OptionParser) StringArrayArrayVar(dest *[][]string, a...) {
    typ := new(StringArrayArrayType);
    op.createOption(a, dest, typ, Append);
}

// HelpType
type HelpType struct {
    option;
    op *OptionParser;
}

func (h *HelpType) storeDefault(dest, def interface{}) {
}

func (h *HelpType) validAction(action *Action, nargs int) bool {
    return action == helpAction;
}

func (h *HelpType) performAction(optStr string, arg []string) {
    usage := h.op.Usage();
    fmt.Println(usage);
    os.Exit(0);
}

func (op *OptionParser) Help(a...) {
    typ := new(HelpType);
    typ.op = op;
    op.createOption(a, nil, typ, helpAction);
}
