package main

import "fmt"
import "strings"

import op "optparse";

var s = "This is some sample text. Watermelon. This is some sample text. Watermelon. This is some sample text.\n Watermelon.     This is some sample text. Watermelon. This is some sample text. Watermelon.";
var p = op.Parser("[options] args...")
var flag = p.Bool("--flag", "-t", op.Help(s))
var invert = p.Bool("--invert", "-T", op.StoreFalse)
var foo = p.String("--foo", "-f", op.Default("default"))
var i = p.Int("--int", "-i", op.Default(78))
var bar = p.StringArray("--bar", "-b", op.Default([]string{"one,two"}))
var c = p.Int("--count", "-c", op.Count)
var baz = p.StringArray("--baz", op.Store, op.Nargs(3), op.Metavar("A B C"))
var list = p.StringArrayArray("--list", op.Nargs(3), op.Metavar("X Y Z"))

func main() {
    p.Callback("--callback", func() { fmt.Println("Callback"); });
    p.Callback("--callback-arg", "-a", func(i int, s string) {
        fmt.Printf("Callback: %d %s\n", i, s);
    }, op.Help(s));
    p.Help("-h", "--help");
    args := p.Parse();
    fmt.Printf("--flag=%t\n", *flag);
    fmt.Printf("--invert=%t\n", *invert);
    fmt.Printf("--foo=%s\n", *foo);
    fmt.Printf("--int=%d\n", *i);
    fmt.Printf("--bar=[%s]\n", strings.Join(*bar, ","));
    fmt.Printf("--count=%d\n", *c);
    fmt.Printf("--baz=[%s]\n", strings.Join(*baz, ","));
    if len(*list) > 0 {
        fmt.Printf("--list=[\n");
        for i := 0; i < len(*list); i++ {
            fmt.Printf("  %v\n", (*list)[i]);
        }
        fmt.Printf("]\n");
    } else {
        fmt.Printf("--list=[]\n");
    }
    fmt.Printf("%v\n", args);
}
