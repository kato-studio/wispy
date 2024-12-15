package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

// Raw HTML input
var rawHTMLSmall = `
<TitleThing title="I Am A Title!">
	lorem ipsum dolor sit amet consectetur adipiscing elit
</TitleThing>
<Card title="Dynamic Card Title 0">
	This is the card body content.
</Card>
<TitleThing title="I Am A Title! [Two]">
	lorem ipsum dolor sit amet consectetur adipiscing elit
</TitleThing>
<Card title="Dynamic Card Title 2">
	This is the card body content.
</Card>
<TitleThing title="I Am A Title! [Three]">
	lorem ipsum dolor sit amet consectetur adipiscing elit
</TitleThing>
<section class="section section--resources section--faq" id="">
	<div class="container  ">
		<h2 class="section__title  " tabindex="0"></h2>
		<p class="section__desc" tabindex="0"></p>
		<div class="section__content" tabindex="0">
	</div>
</section>
`

func init() {
	rawHTMLSmall = strings.Repeat(rawHTMLSmall, 3)
}

func isCapitalByte(b byte) bool {
	return 'A' <= b && b <= 'Z'
}

func newScannerV1(w http.ResponseWriter, req *http.Request) {
	var result strings.Builder

	// var scanOne_time = time.Now()
	scanner := bufio.NewScanner(strings.NewReader(rawHTMLSmall))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		length := len(bytes)
		if length > 3 {
			fmt.Println(string(bytes[:3]))
			//
			fir_byte := bytes[0]
			sec_byte := bytes[1]
			//
			if fir_byte == '<' && isCapitalByte(sec_byte) {
				result.WriteString("<Component:")
				continue
			}
		}
		// isCapitalByte
		result.WriteRune(' ')
		result.Write(bytes)
	}
	// scanOne_Duration := time.Since(scanOne_time)
	// w.Write([]byte(result.String()))
	io.WriteString(w, result.String())
}

// isSpace reports whether the character is a Unicode white space character.
// We avoid dependency on the unicode package, but check validity of the implementation
// in the tests.
func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

func ScanV2(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r == '<' {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		// if isSpace(r) {
		// 	return i + width, data[start:i], nil
		// }
		if r == '>' {
			return i + width, data[start : i+1], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func newScannerV2(w http.ResponseWriter, req *http.Request) {
	var result strings.Builder
	var scanOne_time = time.Now()
	//
	scanner := bufio.NewScanner(strings.NewReader(rawHTMLSmall))
	fmt.Println("NewReader Duration: ", time.Since(scanOne_time))
	scanner.Split(ScanV2)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		length := len(bytes)
		if length > 3 {
			fmt.Println("=-  ", string(bytes))
			//
			fir_byte := bytes[0]
			sec_byte := bytes[1]
			//
			if fir_byte == '<' && isCapitalByte(sec_byte) {
				result.WriteString("<Component:")
				continue
			}
		}

		// isCapitalByte
		result.WriteRune(' ')
		result.Write(bytes)
	}
	scanOne_Duration := time.Since(scanOne_time)
	fmt.Println("ScanV2 Duration: ", scanOne_Duration)
	io.WriteString(w, result.String())
}
