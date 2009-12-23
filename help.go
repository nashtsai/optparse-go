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
func linewrap(s string, width int) []string {
    start := 0;
    length := -1;
    words := splitWords(s);
    lines := make([]string, 0, 5);
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

    optStrs := make([]string, len(op.options));
    optLong := make([]bool, len(op.options));
    //helps := make([][]string, len(op.options));
    //lines := make([]string, 0, 10);
    _, binName := path.Split(os.Args[0]);
    lines := []string {
        fmt.Sprintf("Usage: %s %s", binName, op.usage),
        "",
        "Options:"
    };
    max := 0;
    for i, opt := range op.options {
        optStr := opt.String();
        optStrs[i] = optStr;
        length := utf8.RuneCountInString(optStr);
        if length > max && length < max_argcol {
            max = length;
        }
        optLong[i] = length >= max_argcol;
    }
    width := COLUMNS - max - filler;
	if width < min_width { width = 55; }
	format := fmt.Sprintf(fmt.Sprintf("%%%ds%%%%-%%ds%%%ds%%%%s", indent, colsep),
		                   " ", max, " ");

    for i, opt := range op.options {
        help := linewrap(opt.getHelp(), width);
        if optLong[i] {
            appendString(&lines, fmt.Sprintf(format, optStrs[i], ""));
            if opt.getHelp() == "" {
                continue;
            }
            for _, line := range help {
                appendString(&lines, fmt.Sprintf(format, "", line));
            }
        } else {
            if opt.getHelp() == "" {
                appendString(&lines, fmt.Sprintf(format, optStrs[i], ""));
                continue;
            }
            firstLine := fmt.Sprintf(format, optStrs[i], help[0]);
            appendString(&lines, firstLine);
            for _, line := range help[1:len(help)] {
                appendString(&lines, fmt.Sprintf(format, "", line));
            }
        }
    }
    return strings.Join(lines, "\n");
}
