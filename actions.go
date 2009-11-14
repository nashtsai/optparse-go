package optparse

import "reflect"

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
