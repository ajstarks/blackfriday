//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
//
// deck rendering backend
//
//

package blackfriday

import (
	"bytes"
	"fmt"
	"strings"
)

// Deck is a type that implements the Renderer interface for Deck output.
//
// Do not create this directly, instead use the DeckRenderer function.
type Deck struct {
	slidenumber int
	xp          float64
	yp          float64
	sp          float64
	cw          float64
	ch          float64
}

// DeckRenderer creates and configures a Deck object, which
// satisfies the Renderer interface.
//
// flags is a set of Deck_* options ORed together (currently no such options
// are defined).
func DeckRenderer(flags int) Renderer {
	return &Deck{xp: 10.0, yp: 90.0, sp: 2.0, cw: 1024, ch: 768}
}

func (options *Deck) renderText(ttype string) string {
	return fmt.Sprintf("<text xp=\"%.2f\" yp=\"%.2f\" sp=\"%.2f\" type=\"%s\">",
		options.xp, options.yp, options.sp, ttype)
}

// render code chunks using verbatim, or listings if we have a language
func (options *Deck) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	options.sp = 1.6
	out.WriteString(options.renderText("code"))
	attrEscape(out, text)
	out.WriteString("</text>\n")
}

func (options *Deck) BlockQuote(out *bytes.Buffer, text []byte) {
	options.xp += 5
	out.WriteString(options.renderText("block"))
	attrEscape(out, text)
	out.WriteString("</text>\n")
}

func (options *Deck) BlockHtml(out *bytes.Buffer, text []byte) {
}

func (options *Deck) Header(out *bytes.Buffer, text func() bool, level int) {
	marker := out.Len()
	options.xp = 10
	switch level {
	case 1:
		options.yp = 90
		options.sp = 4
		out.WriteString(options.renderText(""))
	case 2:
		options.yp = 40
		options.sp = 3
		out.WriteString(options.renderText(""))
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</text>\n")
	options.yp -= 10
}

func (options *Deck) HRule(out *bytes.Buffer) {
	options.slidenumber++
	options.xp = 10
	options.yp = 90
	options.sp = 2
	out.WriteString("</slide>\n<slide>\n")
}

func (options *Deck) renderList(ltype string) string {
	return fmt.Sprintf("<list xp=\"%.2f\" yp=\"%.2f\" sp=\"%.2f\" type=\"%s\">\n",
		options.xp, options.yp, options.sp, ltype)
}

func (options *Deck) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	options.sp = 2
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString(options.renderList("number"))
	} else {
		if flags&LIST_ITEM_PLAIN != 0 {
			out.WriteString(options.renderList("plain"))
		} else {
			out.WriteString(options.renderList("bullet"))
		}
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</list>\n")
}

func (options *Deck) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("<li>")
	out.Write(text) // attrEscape(out, text)
	out.WriteString("</li>\n")
}

func (options *Deck) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	options.sp = 2
	out.WriteString(options.renderText("block"))
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</text>\n")
	options.yp -= 20
}

func (options *Deck) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
}

func (options *Deck) TableRow(out *bytes.Buffer, text []byte) {
}

func (options *Deck) TableCell(out *bytes.Buffer, text []byte, align int) {
}

func (options *Deck) Footnotes(out *bytes.Buffer, text func() bool) {
}

func (options *Deck) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
}

func (options *Deck) AutoLink(out *bytes.Buffer, link []byte, kind int) {
}

func (options *Deck) CodeSpan(out *bytes.Buffer, text []byte) {
}

func (options *Deck) DoubleEmphasis(out *bytes.Buffer, text []byte) {
}

func (options *Deck) Emphasis(out *bytes.Buffer, text []byte) {
}

func (options *Deck) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	imgattr := strings.Split(string(alt), ",")

	out.WriteString("<image name=\"")
	attrEscape(out, link)
	if len(imgattr) == 4 {
		out.WriteString(fmt.Sprintf("\" xp=\"%s\" yp=\"%s\" width=\"%s\" height=\"%s\"",
			imgattr[0], imgattr[1], imgattr[2], imgattr[3]))
	} else {
		out.WriteString(" xp=\"50\" yp=\"50\" width=\"100\" height=\"100\"")
	}

	if len(title) > 0 {
		out.WriteString(` caption="`)
		attrEscape(out, title)
		out.WriteString(`"`)
	}
	out.WriteString(" />\n")
}

func (options *Deck) LineBreak(out *bytes.Buffer) {
}

func (options *Deck) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {

}

func (options *Deck) RawHtmlTag(out *bytes.Buffer, tag []byte) {
}

func (options *Deck) TripleEmphasis(out *bytes.Buffer, text []byte) {

}

func (options *Deck) StrikeThrough(out *bytes.Buffer, text []byte) {
}

func (options *Deck) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {

}

func (options *Deck) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (options *Deck) NormalText(out *bytes.Buffer, text []byte) {
	attrEscape(out, text)
}

// header and footer
func (options *Deck) DocumentHeader(out *bytes.Buffer) {
	out.WriteString("<deck>\n<slide>\n")
}

func (options *Deck) DocumentFooter(out *bytes.Buffer) {
	out.WriteString("</slide>\n</deck>\n")
}
