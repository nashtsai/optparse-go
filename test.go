package main

import "fmt"
import "strings"

import op "optparse";

var foo = op.String("--foo", "-f", op.Default("default"))
var i = op.Int("--int", "-i")
var bar = op.StringArray("--bar", "-b", op.Append, op.Default([]string{"one,two"}))
var c = op.Int("--count", "-c", op.Count)
var baz = op.StringArray("--baz", op.Nargs(3))

func main() {
    op.Parse();
    fmt.Printf("--foo=%s\n", *foo);
    fmt.Printf("--int=%d\n", *i);
    fmt.Printf("--bar=[%s]\n", strings.Join(*bar, ","));
    fmt.Printf("--count=%d\n", *c);
    fmt.Printf("--baz=[%s]\n", strings.Join(*baz, ","));
    fmt.Printf("%v\n", op.Args());
}
