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

import "strings"
import "utf8"

// Splits the string s over a number of lines width characters wide.
func Linewrap(s string, width int) []string {
    start := 0;
    length := -1;
    words := splitWords(s);
    lines := make([]string, 0, 5);
    for i, word := range words {
        wordLen := utf8.RuneCountInString(word);
        if length + wordLen + 1 > width {
            lines = appendString(lines, strings.Join(words[start:i], " "));
            start = i;
            length = wordLen;
        } else {
            length += wordLen + 1;
        }
    }
    lines = appendString(lines, strings.Join(words[start:len(words)], " "));
    return lines;
}
