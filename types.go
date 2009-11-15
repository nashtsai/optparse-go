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
import "strconv"

type Type interface {
    parseArg(opt string, arg []string) interface{};
    storeDefault(dest, def interface{});
    validAction(action *Action, nargs int) bool;
}

// BoolType
type BoolType struct {
}

func (b *BoolType) parseArg(opt string, arg []string) interface{} {
    return nil;
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

func Bool(a...) *bool {
    dest := new(bool);
    BoolVar(dest, a);
    return dest;
}

func BoolVar(dest *bool, a...) {
    typ := new(BoolType);
    createOption(a, dest, typ, StoreTrue);
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

func (cb *CallbackType) validAction(action *Action, nargs int) bool {
    // Each callback uses a different Action object, but they all have the
    // same function.
    return action.fn == callbackAction.fn;
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
