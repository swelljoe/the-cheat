// The Cheat
// Produces beautiful, simple, HTML cheat sheets from Markdown tables
// Copyright 2017 Joe Cooper <swelljoe@gmail.com>
// Distributed under the Apache license

package main

import (
	"bytes"
	"flag"
	"fmt"
	bf "gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"os"
)

const version = "0.1.0"

func main() {
	// Command line args
	var page bool
	var css string
	flag.BoolVar(&page, "page", true,
		"Generate a standalone HTML page")
	flag.StringVar(&css, "css", "",
		"Link to a custom CSS style sheet (implies -page)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "The Cheat Cheat Sheet Maker v"+version+
			"Copyright 2017 Joe Cooper <swelljoe@gmail.com>\n\n"+
			"Usage:\n"+
			"	%s [options] [inputfile [outputfile]]\n\n"+
			"Options:\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Implied Options
	if css != "" {
		page = true
	} else {
		css = "css/cheat.css"
	}

	var input []byte
	var err error
	args := flag.Args()
	switch len(args) {
	case 0:
		if input, err = ioutil.ReadAll(os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
			os.Exit(-1)
		}
	case 1, 2:
		if input, err = ioutil.ReadFile(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from", args[0], ":", err)
			os.Exit(-1)
		}
	default:
		flag.Usage()
		os.Exit(-1)
	}

	md := bf.New(bf.WithExtensions(bf.CommonExtensions))
	ast := md.Parse(input)
	var buff bytes.Buffer
	r := bf.NewHTMLRenderer(bf.HTMLRendererParameters{})
	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == bf.Table {
			if entering {
				buff.WriteString("<div>\n")
				r.RenderNode(&buff, node, entering)
			} else {
				r.RenderNode(&buff, node, entering)
				buff.WriteString("</div>\n")
			}
		} else {
			r.RenderNode(&buff, node, entering)
		}
		return bf.GoToNext
	})
	fmt.Printf("%s\n", buff.Bytes())
}
