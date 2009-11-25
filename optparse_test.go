package optparse_test

import "os"
import "reflect"
import "testing"

import op "optparse"

func assertEqual(t *testing.T, a, b interface{}, msg string) {
    if !reflect.DeepEqual(a, b) {
        t.Errorf("%#v != %#v: %s", a, b, msg);
    }
}

func checkArgs(t *testing.T, args []string, err os.Error) {
    if len(args) > 0 {
        t.Error("Got args without passing any");
    }
    if err != nil {
        t.Fatalf("Received error from parsing: %s", err);
    }
}

func TestBool(t *testing.T) {
    p := op.NewParser("", 0);
    a := p.Bool("-a");
    assertEqual(t, *a, false, "Bool defaulted to true");
    b := p.Bool("-b", op.StoreFalse);
    assertEqual(t, *b, true, "Bool with StoreFalse defaulted to false");
    args, err := p.ParseArgs([]string{"-a", "-b"});
    checkArgs(t, args, err);
    assertEqual(t, *a, true, "Bool did not store true");
    assertEqual(t, *b, false, "Bool with StoreFalse did not store false");
}

func TestString(t *testing.T) {
    p := op.NewParser("", 0);
    a := p.String("-a");
    assertEqual(t, *a, "", "String did not default to empty string");
    b := p.String("-b", op.Const("foo"));
    assertEqual(t, *b, "", "String with Const did not default to empty string");
    args, err := p.ParseArgs([]string{"-a", "blah", "-b"});
    checkArgs(t, args, err);
    assertEqual(t, *a, "blah", "String did not get passed argument");
    assertEqual(t, *b, "foo", "String did not get Const value");
}

func TestInt(t *testing.T) {
    p := op.NewParser("", 0);
    a := p.Int("-a");
    assertEqual(t, *a, 0, "Int did not default to 0");
    b := p.Int("-b", op.Count);
    assertEqual(t, *b, 0, "Int with Count did not default to 0");
    c := p.Int("-c", op.Const(10));
    assertEqual(t, *c, 0, "Int with Const did not default to 0");
    args, err := p.ParseArgs([]string{"-a", "20", "-bbb", "-c"});
    checkArgs(t, args, err);
    assertEqual(t, *a, 20, "Int did not get passed argument");
    assertEqual(t, *b, 3, "Int did not count correctly");
    assertEqual(t, *c, 10, "Int did not get Const value");
}

func TestCallback(t *testing.T) {
    p := op.NewParser("", 0);
    a := 0;
    callbackOne := func() { a = 1; };
    p.Callback("-a", callbackOne);
    b := 0;
    callbackTwo := func(i int) { b = i; };
    p.Callback("-b", callbackTwo);
    args, err := p.ParseArgs([]string{"-a", "-b", "1"});
    checkArgs(t, args, err);
    assertEqual(t, a, 1, "Callback did not fire");
    assertEqual(t, b, 1, "Callback with argument did not fire");
}
