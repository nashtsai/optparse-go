package optparse

import "fmt"
import "reflect"
import "strconv"

type Action struct {
    fn func (*option, string, []string);
    hasArgs bool;
}

// StoreConst
type StoreConst struct {
    option;
}

// Count
/*
type Count struct {
    option;
}
*/

var Count = &Action{
    fn: func (c *option, opt string, arg []string) {
        c.typ.(incrementable).increment(c.dest);
    },
    hasArgs: false
}

// Store
var Store = &Action{
    fn: func (s *option, optStr string, arg []string) {
        val := reflect.NewValue(s.typ.parseArg(optStr, arg));
        reflect.NewValue(s.dest).(*reflect.PtrValue).Elem().SetValue(val);
    },
    hasArgs: true
}

// Append
var Append = &Action{
    fn: func (a *option, opt string, arg []string) {
        val := a.typ.parseArg(opt, arg);
        a.typ.(array).append(a.dest, val);
    },
    hasArgs: true
}

var converters = map[reflect.Type]func (a string) interface{} {
    reflect.Typeof(""): func (a string) interface{} {
        return a;
    }
}

var callbackAction = &Action{
    fn: func (c *option, opt string, arg []string) {
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
