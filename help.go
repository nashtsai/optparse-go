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
import "path"
import "strings"
import "utf8"
import "strconv"

// Splits the string s over a number of lines width characters wide.
func linewrap(s string, width int, firstblank bool) []string {
    start := 0;
    length := -1;
    words := splitWords(s);
    lines := make([]string, 0, 5);
    if firstblank {
        appendString(&lines, "");
    }
    for i, word := range words {
        wordLen := utf8.RuneCountInString(word);
        if length + wordLen + 1 > width {
            appendString(&lines, strings.Join(words[start:i], " "));
            start = i;
            length = wordLen;
        } else {
            length += wordLen + 1;
        }
    }
    appendString(&lines, strings.Join(words[start:len(words)], " "));
    return lines;
}

func helpLines(opts []Option) []string {
    return nil;
}


func maxOptionColsize(op *OptionParser, max, limit int) int {
    for _, opt := range op.options {
        length := utf8.RuneCountInString(opt.String());
        if length > max && length < limit {
            max = length
        }
    }
    return max
}


func optionUsage(lines *[]string, opt *Option, format string, width int, max_argcol int) {
    optstr := opt.String();
    help := linewrap(opt.getHelp(), width, len(optstr) > max_argcol);
    appendString(lines, fmt.Sprintf(format, optstr, help[0]))
    for _, line := range help[1:len(help)] {
        appendString(lines, fmt.Sprintf(format, "", line));
    }
}

func (op *OptionParser) Usage() string {

	//   |  -a, --arg    help information for arg      |
    //   |  -b           help information for b        |
    //   |  -l, --long-arg                             |
    //   |               help information for a long   |
    //   |               argument spanning multiple    |
    //   |               lines                         |
    //   ^--^---------^--^--------------------------^--^
    //    |     |     ||            |                | |
    //    |    max    |`-colsep   width              | `-COLUMNS
    //    indent      max_argcol                     `-gutter
    //
    //    COLUMNS    = read from env var of same name, defaulting to 80
    //                 if env var not set or setting is too small (min_width)
    //    min_argcol = minimum value of max

	const (
		indent = 2;
		colsep = 2;
		gutter = 1;
        min_argcol = 4;
        min_width = 5;
	)
    filler := indent + colsep + gutter;
	COLUMNS,enverr := strconv.Atoi(os.Getenv("COLUMNS")); 
	if enverr != nil || COLUMNS < min_width {
		COLUMNS = 80;
	}
	max_argcol := COLUMNS / 3 - 2;
	if max_argcol < min_argcol { max_argcol = min_argcol; }

    //helps := make([][]string, len(op.options));
    //lines := make([]string, 0, 10);
    _, binName := path.Split(os.Args[0]);
    lines := []string {
        fmt.Sprintf("Usage: %s %s", binName, op.usage),
        "",
        "Options:",
    };
    max := maxOptionColsize(op, 0, max_argcol);
    width := COLUMNS - max - filler;
    if width < min_width { width = 55; }
    format := fmt.Sprintf(fmt.Sprintf("%%%ds%%%%-%%ds%%%ds%%%%s", indent, colsep),
                           " ", max, " ");

    for _, opt := range op.options {
        optionUsage(&lines, &opt, format, width, max_argcol);
    }
    return strings.Join(lines, "\n");
}
