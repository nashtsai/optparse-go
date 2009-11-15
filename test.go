package main

import "fmt"
import "strings"

import op "optparse";

var foo = op.String("--foo", "-f", op.Default("default"))
var i = op.Int("--int", "-i")
var bar = op.StringArray("--bar", "-b", op.Default([]string{"one,two"}))
var c = op.Int("--count", "-c", op.Count)
var baz = op.StringArray("--baz", op.Store, op.Nargs(3))
var list = op.StringArrayArray("--list", op.Nargs(3))

func main() {
    op.Parse();
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
