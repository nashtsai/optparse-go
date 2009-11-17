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

    parseArg(opt string, arg []string) interface{};
    storeDefault(dest, def interface{});
    validAction(action *Action, nargs int) bool;
    getOption() *option;
    getDest() interface{};
    getConst() interface{};
}

type option struct {
    longOpts []string;
    shortOpts []string;
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
createOption(args, dest interface{}, typ Option, action *Action)
Option
{
    v := reflect.NewValue(args).(*reflect.StructValue);
    opts := make([]string, v.NumField());
    max := 0;
    opt := typ.getOption();
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
            opt.help = strings.TrimSpace(f.x);
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
            }
        }
    }
    opt.action = action;
    if action == helpAction && opt.help == "" {
        opt.help = "Print this help message and exit.";
    }
    if opt.nargs == 0 && typ.hasArgs() {
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
    typ.setOpts(opts[0:max]);
    if opt.argdesc == "" && opt.nargs > 0 {
        if len(opt.longOpts) > 0 {
            tmp := opt.longOpts[0];
            opt.argdesc = tmp[2:len(tmp)];
        } else {
            tmp := opt.shortOpts[0];
            opt.argdesc = tmp[1:len(tmp)];
        }
        opt.argdesc = strings.ToUpper(opt.argdesc);
        // strings.Replace would be nice...
        opt.argdesc = strings.Map(func(x int)int {
            if x == '-' {
                return '_';
            }
            return x;
        }, opt.argdesc);
    }
    op.appendOpt(typ);
    return typ;
}

func (o *option) getOption() *option {
    return o;
}

func (o *option) getDest() interface{} {
    return o.dest;
}

func (o *option) getConst() interface{} {
    return o.const_;
}

func (o *option) parseArg(opt string, arg []string) interface{} {
    return nil;
}

func (o *option) storeDefault(dest, def interface{}) {}
func (o *option) validAction(action *Action, nargs int) bool {
    return false;
}

func (o *option) performAction(optStr string, arg []string) {
    o.action.fn(o, optStr, arg);
}

func (o *option) hasArgs() bool {
    return o.action.hasArgs;
}

func (o *option) String() string {
    var ret string;
    if o.nargs == 0 {
        short := strings.Join(o.shortOpts, ", ");
        long := strings.Join(o.longOpts, ", ");
        if short != "" && long != "" {
            ret = short + ", " + long;
        } else if short != "" {
            ret = short;
        } else {
            ret = long;
        }
    } else {
        parts := make([]string, len(o.shortOpts) + len(o.longOpts));
        for i, opt := range o.shortOpts {
            parts[i] = opt + " " + o.argdesc;
        }
        for i, opt := range o.longOpts {
            parts[i + len(o.shortOpts)] = opt + " " + o.argdesc;
        }
        ret = strings.Join(parts, ", ");
    }
    return ret;
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
