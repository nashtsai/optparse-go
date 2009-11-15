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
import "reflect"
import "strings"

type _Help struct { x string; }
func Help(h string) *_Help { return &_Help{h} }
type _Nargs struct { x int; }
func Nargs(n int) *_Nargs { return &_Nargs{n} }
type _Argdesc struct { x string; }
func Metavar(a string) *_Argdesc { return &_Argdesc{a} }
type _Default struct { x interface{}; }
func Default(d interface{}) *_Default { return &_Default{d} }
type _Const struct { x interface{}; }
func Const(c interface{}) *_Const { return &_Const{c} }

/*
const (
    Store = Action(iota);
    Append;
    Count;
)
*/

type Option interface {
    getNargs() int;
    hasArgs() bool;
    performAction(string, []string);
    //setType(Type);
    //setDest(interface{});
    getHelp() string;
    setOpts([]string);
    String() string;
    matches(string) bool;
}

type option struct {
    longOpts []string;
    shortOpts []string;
    typ Type;
    dest interface{};
    help string;
    argdesc string;
    nargs int;
    action *Action;
    const_ interface{};
}

func destTypecheck(dest, value interface{}) bool {
    return reflect.Typeof(dest).(*reflect.PtrType).Elem() == reflect.Typeof(value);
}

func (op *OptionParser)
createOption(args, dest interface{}, typ Type, action *Action)
Option
{
    v := reflect.NewValue(args).(*reflect.StructValue);
    opts := make([]string, v.NumField());
    max := 0;
    opt := new(option);
    opt.typ = typ;
    opt.dest = dest;
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i);
        switch f := field.Interface().(type) {
        case string:
            opts[max] = f;
            max++;
        case *Action:
            action = f;
            if action == StoreFalse {
                typ.storeDefault(dest, true);
            }
        case *_Help:
            opt.help = f.x;
        case *_Nargs:
            opt.nargs = f.x;
        case *_Argdesc:
            opt.argdesc = f.x;
        case *_Default:
            if false { fmt.Printf("%v\n", *f); }
            if !destTypecheck(dest, f.x) {
                ProgrammerError(fmt.Sprintf("%s: Type mismatch with default value.", opts[0]));
            }
            typ.storeDefault(dest, f.x);
        case *_Const:
            opt.const_ = f.x;
        default:
            fn, ok := field.(*reflect.FuncValue);
            if ok {
                fnType := fn.Type().(*reflect.FuncType);
                opt.nargs = fnType.NumIn();
                typ.(*CallbackType).fn = f;
                tmp := new(Action);
                tmp.name = callbackAction.name;
                tmp.fn = callbackAction.fn;
                tmp.hasArgs = opt.nargs > 0;
                action = tmp;
            }
        }
    }
    if opt.nargs == 0 && action.hasArgs {
        opt.nargs = 1;
    }
    if max == 0 {
        ProgrammerError("Option has no options!");
        return nil;
    }
    if !typ.validAction(action, opt.nargs) {
        ProgrammerError(fmt.Sprintf("Option '%s' is using invalid action '%s'.", opts[0], action.name));
        return nil;
    }
    if opt.const_ != nil && !destTypecheck(dest, opt.const_) {
        ProgrammerError(fmt.Sprintf("%s: Type mismatch with constant value.", opts[0]));
        return nil;
    }
    opt.action = action;
    opt.setOpts(opts[0:max]);
    op.appendOpt(opt);
    return opt;
}

func (o *option) performAction(optStr string, arg []string) {
    o.action.fn(o, optStr, arg);
}

func (o *option) hasArgs() bool {
    return o.action.hasArgs;
}

func (o *option) String() string {
    return strings.Join(o.longOpts, ",") + "," + strings.Join(o.shortOpts, ",")
}

func (o *option) setOpts(opts []string) {
    i := len(opts);
    longOpts := make([]string, 0, i);
    shortOpts := make([]string, 0, i);
    for _, opt := range opts {
        if strings.HasPrefix(opt, "--") {
            longOpts = longOpts[0:len(longOpts) + 1];
            longOpts[len(longOpts) - 1] = opt;
        } else if strings.HasPrefix(opt, "-") {
            shortOpts = shortOpts[0:len(shortOpts) + 1];
            shortOpts[len(shortOpts) - 1] = opt;
        }
    }
    o.longOpts = longOpts;
    o.shortOpts = shortOpts;
}

func (o *option) getNargs() int {
    return o.nargs;
}

func (o *option) getHelp() string {
    return o.help;
}

func (o *option) setType(t Type) {
    o.typ = t;
}
func (o *option) setDest(ptr interface{}) {
    o.dest = ptr;
}

func (o *option) matches(opt string) bool {
    if len(opt) < 2 ||
       len(opt) == 2 && (opt[0] != '-' || opt[1] == '-') ||
       len(opt) > 2 && opt[0:2] != "--" {
        return false;
    }
    if len(opt) == 2 {
        for _, s := range o.shortOpts {
            if s == opt {
                return true;
            }
        }
    } else {
        for _, s := range o.longOpts {
            if s == opt {
                return true;
            }
        }
    }
    return false;
}
