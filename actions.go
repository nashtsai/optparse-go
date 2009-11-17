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

import "reflect"

type Action struct {
    name string;
    fn func (Option, string, []string);
    hasArgs bool;
}

// StoreConst
var StoreConst = &Action{
    name: "StoreConst",
    fn: func (c Option, opt string, arg []string) {
        val := reflect.NewValue(c.getConst());
        elem := reflect.NewValue(c.getDest()).(*reflect.PtrValue).Elem();
        elem.SetValue(val);
    },
    hasArgs: false
}

// StoreTrue
var StoreTrue = &Action{
    name: "StoreTrue",
    fn: func (c Option, opt string, arg []string) {
        *c.getDest().(*bool) = true;
    },
    hasArgs: false
}

// StoreFalse
var StoreFalse = &Action{
    name: "StoreFalse",
    fn: func (c Option, opt string, arg []string) {
        *c.getDest().(*bool) = false;
    },
    hasArgs: false
}

// Count
var Count = &Action{
    name: "Count",
    fn: func (c Option, opt string, arg []string) {
        c.(incrementable).increment(c.getDest());
    },
    hasArgs: false
}

// Store
var Store = &Action{
    name: "Store",
    fn: func (s Option, optStr string, arg []string) {
        val := reflect.NewValue(s.parseArg(optStr, arg));
        reflect.NewValue(s.getDest()).(*reflect.PtrValue).Elem().SetValue(val);
    },
    hasArgs: true
}

// Append
var Append = &Action{
    name: "Append",
    fn: func (a Option, opt string, arg []string) {
        val := a.parseArg(opt, arg);
        a.(array).append(a.getDest(), val);
    },
    hasArgs: true
}

var callbackAction = &Action{
    name: "callbackAction",
    fn: nil,
    // this is ignored
    hasArgs: true
}

var helpAction = &Action{
    name: "helpAction",
    fn: nil,
    hasArgs: false
}
