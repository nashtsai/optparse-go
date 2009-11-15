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
import "strings"

type OptionParser struct {
    options []Option;
    _args []string;
}

func Parser() *OptionParser {
    ret := new(OptionParser);
    ret.options = make([]Option, 0, 10);
    ret._args = make([]string, 0, 10);
    return ret;
}

//var options = make([]Option, 0, 10)
func (op *OptionParser) appendOpt(opt Option) {
    op.options = op.options[0:len(op.options)+1];
    op.options[len(op.options)-1] = opt;
    if len(op.options) == cap(op.options) {
        tmp := make([]Option, len(op.options), cap(op.options) * 2);
        for i, e := range op.options {
            tmp[i] = e;
        }
        op.options = tmp;
    }
}
func (op *OptionParser) matches(s string) Option {
    for _, option := range op.options {
        if option.matches(s) {
            return option;
        }
    }
    return nil;
}

func Error(opt, msg string) {
    fmt.Fprintf(os.Stderr, "Error: %s: %s\n%s\n", opt, msg/*, Usage()*/);
    os.Exit(1);
}
func ProgrammerError(msg string) {
    fmt.Fprintf(os.Stderr, "Programmer error: %s\n", msg);
    os.Exit(2);
}

func (op *OptionParser) Usage() string {
    return "";
}

func (op *OptionParser) Parse() {
    op.ParseArgs(os.Args[1:len(os.Args)]);
}

//var _args = make([]string, 0, 10);
func (op *OptionParser) Args() []string {
    return op._args;
}
func (op *OptionParser) appendArg(arg string) {
    op._args = op._args[0:len(op._args)+1];
    op._args[len(op._args)-1] = arg;
    if len(op._args) == cap(op._args) {
        tmp := make([]string, len(op._args), cap(op._args) * 2);
        for i, e := range op._args {
            tmp[i] = e;
        }
        op._args = tmp;
    }
}

func (op *OptionParser) invalid(arg string) {
    Error(arg, "invalid option");
}

func (op *OptionParser)
doAction(opt, arg string, hasArg bool, args []string, i int)
(int, bool)
{
    var current []string;
    usedArg := false;
    option := op.matches(opt);
    if option == nil {
        op.invalid(opt);
        return i, false
    }
    nargs := option.getNargs();
    if nargs > 0 {
        current = make([]string, nargs);
        j := 0;
        if hasArg {
            current[0] = arg;
            j = 1;
            usedArg = true;
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
    return i, usedArg
}

func (op *OptionParser) ParseArgs(args []string) {
    var arg string;
    var hasArg bool;
    for i := 0; i < len(args); i++ {
        opt := args[i];
        if opt == "--" {
            i++;
            for ; i < len(args); i++ {
                op.appendArg(args[i]);
            }
        } else if strings.HasPrefix(opt, "--") {
            idx := strings.Index(opt, "=");
            if idx != -1 {
                arg = opt[idx + 1:len(opt)];
                hasArg = true;
                opt = opt[0:idx];
            } else {
                hasArg = false;
            }
            i, _ = op.doAction(opt, arg, hasArg, args, i);
        } else if strings.HasPrefix(opt, "-") {
            for j, c := range opt[1:len(opt)] {
                s := "-" + string(c);
                if j == len(opt) - 2 {
                    hasArg = false;
                } else {
                    arg = opt[j + len(s):len(opt)];
                    hasArg = true;
                }
                var usedArg bool;
                i, usedArg = op.doAction(s, arg, hasArg, args, i);
                if usedArg {
                    break;
                }
            }
        } else {
            op.appendArg(opt);
        }
    }
}
