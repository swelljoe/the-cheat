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
	var css, title, header string
	var cols int
	flag.BoolVar(&page, "page", true,
		"Generate a standalone HTML page")
	flag.StringVar(&title, "title", "",
		"Title for the page (implies -page)")
	flag.StringVar(&header, "header", "",
		"Header for the page (implies -page)")
	flag.StringVar(&css, "css", "",
		"Link to a custom CSS style sheet (implies -page)")
	flag.IntVar(&cols, "cols", 3,
		"Maximum number of columns on page (implies -page)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "The Cheat - A Cheat Sheet Maker v"+version+
			"\nCopyright 2017 Joe Cooper <swelljoe@gmail.com>\n\n"+
			"Usage:\n"+
			"%s [options] [inputfile [outputfile]]\n\n"+
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

	if title != "" {
		page = true
	}

	if title != "" {
		page = true
	}

	// This might make it impossible to leave page off
	if cols > 0 {
		page = true
	}

	var input []byte
	var err error
	args := flag.Args()
	switch len(args) {
	case 1, 2:
		if input, err = ioutil.ReadFile(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from", args[0], ":", err)
			os.Exit(-1)
		}
	default:
		flag.Usage()
		os.Exit(-1)
	}

	//var htmlFlags bf.HTMLFlags

	//if page {
	//	htmlFlags = bf.CompletePage
	//}

	md := bf.New(bf.WithExtensions(bf.CommonExtensions))
	ast := md.Parse(input)
	var buff bytes.Buffer
	if page {
		writeHeader(&buff, title, header)
	}
	buff.WriteString("<div class='cheat flex three'>\n")
	r := bf.NewHTMLRenderer(bf.HTMLRendererParameters{
		Flags: bf.CompletePage})
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
	buff.WriteString("</div>\n")
	if page {
		writeFooter(&buff)
	}
	fmt.Printf("%s\n", buff.Bytes())
}

func writeHeader(w *bytes.Buffer, title string, header string) {
	w.WriteString("<!DOCTYPE html>\n")
	w.WriteString("<html lang='en-us'>\n")
	w.WriteString("<head>\n")
	w.WriteString("<meta charset='utf-8'>\n")
	if title != "" {
		w.WriteString("<title>")
		w.WriteString(title)
		w.WriteString("</title>\n")
	}
	w.WriteString("<link rel='stylesheet' href='css/cheat.min.css' />\n")
	w.WriteString("</head>\n\n")
	w.WriteString("<body>\n")
	w.WriteString("<main>\n")
	if header != "" {
		w.WriteString("<h1>")
		w.WriteString(header)
		w.WriteString("</h1>\n")
	}
}

func writeFooter(w *bytes.Buffer) {
	w.WriteString("</main>\n")
	w.WriteString("</body>\n")
	w.WriteString("</html>\n")
}
