package main

import "fmt"
//import "strings"

import op "optparse";

var foo = op.String("--foo", "-f", op.Default("default"))
var i = op.Int("--int", "-i")
//var bar = optparse.StringArray("--bar", "-b", optparse.Append)
//var c = optparse.Int("--count", &optparse.Count{})
//var baz = optparse.StringArray("--baz", &optparse.Store{Nargs: 3})

func main() {
    //optparse.AddOption(optparse.StoreString{S: "foo"}, "bar", 123);
    op.Parse();
    fmt.Printf("--foo=%s\n", *foo);
    fmt.Printf("--int=%d\n", *i);
//    fmt.Printf("--bar=[%s]\n", strings.Join(*bar, ","));
//    fmt.Printf("--count=%d\n", *c);
    fmt.Printf("%v\n", op.Args());
}
