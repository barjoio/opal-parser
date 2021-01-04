package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

var markupBytes []byte
var err error

var lineBuffer string
var lines []string
var lineNum = 1
var linePos = 0

var elements [][]string

// current rune and index
var r rune
var i int

var currentTag string

var textBufferStack []string
var contextStack []string
var detectedPosStack []int

var contextMap = map[string]func(){
	"start": func() {
		if r == '!' {
			contextSet("tagName", true)
		} else if r == '\\' {
			textBufferNew()
			contextStackPush("escapeChar", true)
		} else if !isSpace(r) {
			textBufferAppend(r)
			contextSet("text", true)
		}
	},
	"escapeChar": func() {
		if !isSpace(r) {
			textBufferAppend(r)
		} else {
			elementAdd("\\", textBufferGet())
			textBufferStackPop()
			contextStackPop()
			if contextGet() == "text" {
				textBufferAppend(r)
			}
		}
	},
	"tagName": func() {
		if isLetter(r) || isNumber(r) {
			textBufferAppend(r)
		} else if r == '\n' {
			if textBufferGet() == "" {
				printErrExpectedTextAfterTag()
			} else {
				elementAdd(textBufferGet(), "")
				textBufferStackPop()
				contextSet("start", false)
			}
		} else if isSpace(r) {
			currentTagSet(textBufferGet())
			textBufferStackPop()
			contextSet("tagContent", false)
		} else {
			printErrUnexpectedChar(r)
		}

	},
	"tagContent": func() {
		if r == '\n' {
			elementAdd(currentTag, textBufferGet())
			textBufferStackPop()
			contextSet("start", false)
		} else {
			textBufferAppend(r)
		}
	},
	"text": func() {
		if r == '\n' {
			elementAdd("", textBufferGet())
			textBufferStackPop()
			contextSet("start", false)
		} else {
			textBufferAppend(r)
		}
	},
	// "readingInlineTag": func() {
	// 	if r != ')' {
	// 		textBufferAdd(r)
	// 	} else {
	// 		elementAdd(currentTag, textBuffer)
	// 		contextSet("readingText")
	// 	}
	// },
}

func main() {
	markupBytes, err = ioutil.ReadFile("file.opal")
	if err != nil {
		panic(err)
	}

	// add new line at the end of markup
	markupBytes = append(markupBytes, byte('\n'))

	// parse markup
	contextSet("start", false)
	for index, char := range markupBytes {
		r = rune(char)
		i = index

		linePos++
		lineBuffer += string(r)

		if contextAction, ok := contextMap[contextGet()]; ok {
			contextAction()
		} else {
			fmt.Println("Unknown context")
		}

		if r == '\n' {
			lineNum++
			linePos = 0
			lines = append(lines, lineBuffer)
			lineBuffer = ""
		}
	}

	for _, e := range elements {
		fmt.Println("Context:", "#", e[0], "#")
		fmt.Println("Line num:", "#", e[1], "#")
		fmt.Println("Line pos:", "#", e[2], "#")
		fmt.Println("Tag name:", "#", e[3], "#")
		fmt.Println("Content:", "#", strings.Join(e[4:], "|"), "#")
		fmt.Println()
	}

	b, err := json.Marshal(elements)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

// prints

func printRemedy(format string, args ...interface{}) {
	fmt.Printf("%s %s.\n", color.GreenString("Try:"), fmt.Sprintf(format, args...))
}

func printErrLocation() {
	errChar := string(lineBuffer[len(lineBuffer)-1])
	if errChar == "\n" { // edge-case for newline being the errChar
		errChar = ""
	}
	beforeErrText := lineBuffer[:len(lineBuffer)-1]
	if len(beforeErrText) >= 5 {
		beforeErrText = beforeErrText[len(beforeErrText)-5:]
	}
	lineBuffer = ""
	for _, c := range markupBytes[i:] {
		if rune(c) != '\n' {
			lineBuffer += string(c)
		} else {
			break
		}
	}
	var afterErrText string
	if len(lineBuffer) > 1 {
		afterErrText = lineBuffer[1:]
	}
	if len(afterErrText) > 5 {
		afterErrText = afterErrText[:5]
	}
	lineBuffer = ""
	lineNumString := strconv.Itoa(lineNum)
	fmt.Printf("%s  %s%s%s\n", color.HiBlackString(lineNumString), beforeErrText, color.RedString(errChar), afterErrText)
	s := fmt.Sprintf("%s%s^%s", strings.Repeat(" ", len(lineNumString+"  ")), strings.Repeat("-", len(beforeErrText)), strings.Repeat("-", len(afterErrText)))
	fmt.Println(color.HiRedString(s))
}

func printErr(format string, args ...interface{}) {
	fmt.Printf("%s %s at %d,%d.\n", color.RedString("Error:"), fmt.Sprintf(format, args...), lineNum, linePos)
	printErrLocation()
}

func printErrFatal(format string, args ...interface{}) {
	printErr(format, args...)
	die()
}

func printErrUnexpectedChar(r rune) {
	printErrFatal("Unexpected character (%s)", string(r))
}

func printErrExpectedTextAfterTag() {
	printErr("Expected a valid tag name after tag declaration")
	printRemedy("Exclamation mark followed by a word or number, such as %s", color.HiGreenString("!Doc"))
	die()
}

func die() {
	os.Exit(1)
}

//

func prevRune() rune {
	return rune(markupBytes[i-1])
}

func nextRune() rune {
	return rune(markupBytes[i+1])
}

//

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}

//

func elementAdd(tagName, tagContent string, additionalAttrs ...string) {
	startPos := detectedPosGet()
	elements = append(elements, append([]string{contextGet(), strconv.Itoa(lineNum), strconv.Itoa(startPos), tagName, tagContent}, additionalAttrs...))
}

// textBufferStack

func textBufferStackPush(s string) {
	textBufferStack = append(textBufferStack, s)
}

func textBufferStackPop() {
	l := len(textBufferStack)
	if l > 0 {
		textBufferStack = textBufferStack[:l-1]
	}
}

func textBufferGet() string {
	l := len(textBufferStack)
	if l > 0 {
		return textBufferStack[l-1]
	}
	return ""
}

func textBufferAppend(r rune) {
	l := len(textBufferStack)
	if l == 0 {
		textBufferStackPush("")
		l++
	}
	textBufferStack[l-1] = textBufferGet() + string(r)
}

func textBufferNew() {
	textBufferStackPush("")
}

// contextStack

func contextStackPush(s string, detectPos bool) {
	contextStack = append(contextStack, s)
	if detectPos {
		detectedPosSet(linePos)
	}
}

func contextStackPop() {
	l := len(contextStack)
	if l > 0 {
		contextStack = contextStack[:l-1]
	}
}

func contextSet(s string, detectPos bool) {
	l := len(contextStack)
	if l == 0 {
		contextStackPush("", false)
		l++
	}
	contextStack[l-1] = s
	if detectPos {
		detectedPosSet(linePos)
	}
}

func contextGet() string {
	l := len(contextStack)
	if l > 0 {
		return contextStack[l-1]
	}
	return ""
}

// detectedPosStack

func detectedPosStackPush(p int) {
	detectedPosStack = append(detectedPosStack, p)
}

func detectedPosStackPop() {
	l := len(detectedPosStack)
	if l > 0 {
		detectedPosStack = detectedPosStack[:l-1]
	}
}

func detectedPosGet() int {
	l := len(detectedPosStack)
	if l > 0 {
		return detectedPosStack[l-1]
	}
	return -1
}

func detectedPosSet(p int) {
	l := len(detectedPosStack)
	if l == 0 {
		detectedPosStackPush(0)
		l++
	}
	detectedPosStack[l-1] = p
}

//

func currentTagSet(s string) {
	currentTag = s
}

func currentTagReset() {
	currentTag = ""
}
