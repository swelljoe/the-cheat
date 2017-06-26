package main

import (
	"bytes"
	"fmt"
	bf "gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
)

func main() {
	input, err := ioutil.ReadFile("sheet.md")
	if err != nil {
		fmt.Print(err)
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
	//html := r.Render(ast)
	fmt.Printf("%s\n", buff.Bytes())
}
