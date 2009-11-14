package optparse

import "fmt"
import "os"
import "strings"

//var optMap = make(map[string]Option)
var options = make([]Option, 0, 10)
func appendOpt(opt Option) {
    options = options[0:len(options)+1];
    options[len(options)-1] = opt;
    if len(options) == cap(options) {
        tmp := make([]Option, len(options), cap(options) * 2);
        for i, e := range options {
            tmp[i] = e;
        }
        options = tmp;
    }
}
func matches(s string) Option {
    for _, option := range options {
        if option.matches(s) {
            return option;
        }
    }
    return nil;
}

func Error(opt, msg string) {
    fmt.Fprintf(os.Stderr, "Error: %s: %s\n%s\n", opt, msg, Usage());
    os.Exit(1);
}

func Usage() string {
    return "";
}

func Parse() {
    ParseArgs(os.Args[1:len(os.Args)]);
}

var _args = make([]string, 0, 10);
func Args() []string {
    return _args;
}
func appendArg(arg string) {
    _args = _args[0:len(_args)+1];
    _args[len(_args)-1] = arg;
    if len(_args) == cap(_args) {
        tmp := make([]string, len(_args), cap(_args) * 2);
        for i, e := range _args {
            tmp[i] = e;
        }
        _args = tmp;
    }
}

func invalid(arg string) {
    Error(arg, "invalid option");
}

func ParseArgs(args []string) {
    var arg string;
    var option Option;
    var current []string;
    //fmt.Printf("%v\n", options);
    for i := 0; i < len(args); i++ {
        opt := args[i];
        if opt == "--" {
            i++;
            for ; i < len(args); i++ {
                appendArg(args[i]);
            }
        } else if strings.HasPrefix(opt, "--") {
            idx := strings.Index(opt, "=");
            if idx != -1 {
                arg = opt[idx + 1:len(opt)];
                opt = opt[0:idx];
            }
            //fmt.Printf("%v\n%v\n%v\n", idx, opt, arg);
            option = matches(opt);
            if option == nil {
                invalid(opt);
                continue
            }
            nargs := option.getNargs();
            if nargs > 0 {
                current = make([]string, nargs);
                j := 0;
                if idx != -1 {
                    current[0] = arg;
                    j = 1;
                }
                for ; j < len(current); j++{
                    i++;
                    if i >= len(args) {
                        Error(opt, "insufficient arguments for option");
                    }
                    current[j] = args[i];
                }
            } else {
                current = nil;
            }
            option.performAction(opt, current);
        } else if strings.HasPrefix(opt, "-") {
        } else {
            appendArg(opt);
        }
    }
}
