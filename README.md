optigo
=========

[![Build Status](https://travis-ci.org/coryb/optigo.svg?branch=master)](https://travis-ci.org/coryb/optigo)
[![Coverage Status](https://coveralls.io/repos/coryb/optigo/badge.svg?branch=master)](https://coveralls.io/r/coryb/optigo?branch=master)
[![GoDoc](https://godoc.org/github.com/coryb/optigo?status.png)](https://godoc.org/github.com/coryb/optigo)

**optigo** is a simple command line options parser.  

```go
package main

import (
	"fmt"
	"github.com/coryb/optigo"
	"os"
)

func main() {
	usage := func() {
		fmt.Println(`
Usage: optigo-test [-h] [-i INT] [-f FLOAT] [-s STRING] [-m STRING]... [-c]...
`)
		os.Exit(0)
	}

	myint := 1
	myfloat := 1.0
	mystr := "default"
	many := make([]string, 0)
	options := make(map[string]string)
	count := 0

	op := optigo.NewDirectAssignParser(map[string]interface{}{
		"h|help":       usage,
		"i|int=i":      &myint,
		"f|float=f":    &myfloat,
		"s|str=s":      &mystr,
		"m|many=s@":    &many,
		"c|count+":     &count,
		"o|options=s%": &options,
	})

	op.ProcessAll(os.Args[1:])

	fmt.Printf("myint: %d\n", myint)
	fmt.Printf("myfloat: %f\n", myfloat)
	fmt.Printf("mystr: %s\n", mystr)
	fmt.Printf("many: %v\n", many)
	fmt.Printf("count: %d\n", count)
	fmt.Printf("options: %v\n", options)
}
```

**optigo** parses command-line arguments either into reference values you pass into the
NewDirectAssingParser constructor or if you use NewParser it will populate the Results
object in the OptionParser object created via NewParser.

## Installation

```go
import "github.com/coryb/optigo"
```

To install optigo according to your `$GOPATH`:

```console
$ go get github.com/coryb/optigo
```

## Usage

#### type OptionParser

```go
type OptionParser struct {
     Results map[string]interface{}
     Args    []string
}
```

OptionParser struct will contain the `Results` and `Args` after one of the
Process routines is called. A OptionParser object is created with either
NewParser or NewDirectAssignParser

#### func  NewDirectAssignParser

```go
func NewDirectAssignParser(opts map[string]interface{}) OptionParser
```
NewDirectAssignParser generates an OptionParser object from the `opts` passed
in. After calling OptionParser.Parser([]string) the options will be assigned
directly to the references passed in `opts`.

#### func  NewParser

```go
func NewParser(opts []string) OptionParser
```
NewParser generates an OptionParser object from the opts passed in. After
calling OptionParser.Parser([]string) the option results will be stored in
OptionParser.Results

#### func (*OptionParser) ProcessAll

```go
func (o *OptionParser) ProcessAll(args []string) error
```
ProcessAll will parse all arguments in args. If there are any arguments in args
that start with '-' and are not known options then an error will be returned.
Any non-options will be available in OptionParser.Args.

#### func (*OptionParser) ProcessSome

```go
func (o *OptionParser) ProcessSome(args []string) error
```
ProcessSome will parse all known arguments in args. Any non-options and unknown
options will be available in OPtionParser.Args. This can be used to implement
multple pass options parsing, for example perhaps sub-commands options are
parsed seperately from global options.

Documentation and examples for optigo are available at
[GoDoc.org](https://godoc.org/github.com/coryb/optigo).

## Notes

The the option declaration syntax was inspired by the Getopt::Long Perl package, but **optigo** makes no attempt to have complete
parity with the Getopt::Long functionality.   