include $(GOROOT)/src/Make.$(GOARCH)

OPTFILES=help.go optparse.go actions.go option.go types.go util.go

test: test.$(O)
	$(LD) -o test test.$(O)
optparse: optparse.$(O)
optparse.$(O): $(OPTFILES)
	$(GC) -o optparse.$(O) $(OPTFILES)
test.$(O): test.go optparse.$(O)
	$(GC) -I. -o test.$(O) test.go
.PHONY: clean
clean:
	rm -f test optparse.$(O) test.$(O)
