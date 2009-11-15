package main

import "fmt"
import "strings"

import op "optparse";

var flag = op.Bool("--flag", "-t")
var invert = op.Bool("--invert", "-T", op.StoreFalse)
var foo = op.String("--foo", "-f", op.Default("default"))
var i = op.Int("--int", "-i")
var bar = op.StringArray("--bar", "-b", op.Default([]string{"one,two"}))
var c = op.Int("--count", "-c", op.Count)
var baz = op.StringArray("--baz", op.Store, op.Nargs(3))
var list = op.StringArrayArray("--list", op.Nargs(3))


func main() {
    op.Callback("--callback", func() { fmt.Println("Callback"); });
    op.Callback("--callback-arg", "-a", func(i int, s string) {
        fmt.Printf("Callback: %d %s\n", i, s);
    });
    op.Parse();
    fmt.Printf("--flag=%t\n", *flag);
    fmt.Printf("--invert=%t\n", *invert);
    fmt.Printf("--foo=%s\n", *foo);
    fmt.Printf("--int=%d\n", *i);
    fmt.Printf("--bar=[%s]\n", strings.Join(*bar, ","));
    fmt.Printf("--count=%d\n", *c);
    fmt.Printf("--baz=[%s]\n", strings.Join(*baz, ","));
    fmt.Printf("--list=[\n");
    for i := 0; i < len(*list); i++ {
        fmt.Printf("  %v\n", (*list)[i]);
    }
    fmt.Printf("]\n");
    fmt.Printf("%v\n", op.Args());
}
