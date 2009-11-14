package optparse

import "fmt"
import "reflect"
import "strings"

type _Help struct { x string; }
func Help(h string) *_Help { return &_Help{h} }
type _Nargs struct { x int; }
func Nargs(n int) *_Nargs { return &_Nargs{n} }
type _Argdesc struct { x string; }
func Argdesc(a string) *_Argdesc { return &_Argdesc{a} }
type _Default struct { x interface{}; }
func Default(d interface{}) *_Default { return &_Default{d} }

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
}

func createOption(args interface{}, dest interface{}, typ Type, action *Action) Option {
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
        case *_Help:
            opt.help = f.x;
        case *_Nargs:
            opt.nargs = f.x;
        case *_Argdesc:
            opt.argdesc = f.x;
        case *_Default:
            if false { fmt.Printf("%v\n", *f); }
            typ.storeDefault(dest, f.x);
        default:
            _, ok := field.(*reflect.FuncValue);
            if ok {
            }
        }
    }
    if opt.nargs == 0 && action.hasArgs {
        opt.nargs = 1;
    }
    opt.action = action;
    opt.setOpts(opts[0:max]);
    appendOpt(opt);
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
