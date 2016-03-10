# stopwatch

Stopwatch is a command line tool to track the start and duration of named timers.  The data is stored in your home directory as `.stopwatch.json`.

![](stopwatch.gif)

### Install

    go get github.com/jweir/stopwatch

### Usage

````
Usage for stopwatch version 0.2:
  stopwatch
         Prints all existing stopwatches.
  stopwatch label...
         Starts a new stopwatch with the given label.
         Or stops a existing stopwatch with that label.

Flags:
  -prompt
    	Prints the label of the first of the stopwatch.
      Prints an empty string if there are no stopwatches.
  -stopall
    	Issues a stop command all stopwatches

````

This was a quick 45 minute project to build a simple stop watch and scratch that itch.

## LICENSE

The MIT License (MIT)
Copyright (c) 2016 John Weir

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
