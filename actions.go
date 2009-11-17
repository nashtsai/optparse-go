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

type Action struct {
    name string;
    fn func (*OptionParser, *option, string, []string);
    hasArgs bool;
}

// StoreConst
var StoreConst = &Action{
    name: "StoreConst",
    fn: func (op *OptionParser, c *option, opt string, arg []string) {
        val := reflect.NewValue(c.const_);
        elem := reflect.NewValue(c.dest).(*reflect.PtrValue).Elem();
        elem.SetValue(val);
    },
    hasArgs: false
}

// StoreTrue
var StoreTrue = &Action{
    name: "StoreTrue",
    fn: func (op *OptionParser, c *option, opt string, arg []string) {
        *c.dest.(*bool) = true;
    },
    hasArgs: false
}

// StoreFalse
var StoreFalse = &Action{
    name: "StoreFalse",
    fn: func (op *OptionParser, c *option, opt string, arg []string) {
        *c.dest.(*bool) = false;
    },
    hasArgs: false
}

// Count
var Count = &Action{
    name: "Count",
    fn: func (op *OptionParser, c *option, opt string, arg []string) {
        c.typ.(incrementable).increment(c.dest);
    },
    hasArgs: false
}

// Store
var Store = &Action{
    name: "Store",
    fn: func (op *OptionParser, s *option, optStr string, arg []string) {
        val := reflect.NewValue(s.typ.parseArg(optStr, arg));
        reflect.NewValue(s.dest).(*reflect.PtrValue).Elem().SetValue(val);
    },
    hasArgs: true
}

// Append
var Append = &Action{
    name: "Append",
    fn: func (op *OptionParser, a *option, opt string, arg []string) {
        val := a.typ.parseArg(opt, arg);
        a.typ.(array).append(a.dest, val);
    },
    hasArgs: true
}

var callbackAction = &Action{
    name: "callbackAction",
    fn: func (op *OptionParser, c *option, opt string, arg []string) {
        fn := reflect.NewValue(c.typ.(*CallbackType).fn).(*reflect.FuncValue);
        fnType := fn.Type().(*reflect.FuncType);
        values := make([]reflect.Value, len(arg));
        for i := 0; i < len(arg); i++ {
            switch v := fnType.In(i).(type) {
            case *reflect.StringType:
                values[i] = reflect.NewValue(arg[i]);
            case *reflect.IntType:
                x, ok := strconv.Atoi(arg[i]);
                if ok != nil {
                    Error(opt, fmt.Sprintf("'%s' is not an integer", arg[i]));
                }
                values[i] = reflect.NewValue(x);
            }
        }
        fn.Call(values);
    },
    // this is ignored
    hasArgs: true
}

var helpAction = &Action{
    name: "helpAction",
    fn: func (op *OptionParser, c *option, opt string, arg []string) {
        usage := op.Usage();
        fmt.Println(usage);
        os.Exit(0);
    },
    hasArgs: false
}
