package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"unicode/utf8"

	"github.com/spf13/pflag"
)

const (
	seperatorSpace = iota
	seperatorLine
	seperatorNull
)

var seperatorMap = map[int][]rune{
	seperatorNull: {null},
	seperatorLine: {newline},
	seperatorSpace: {'\t', '\v', '\f', '\r', newline, space, '\u00FF',
		'\u0085',
		'\u00A0', '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000'},
}

var null = '\u0000'
var space = ' '
var newline = '\n'

var separator int

var ScanNull = func(data []byte, atEOF bool) (advance int,
	token []byte,
	err error) {
	// Skip leading null bytes.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isNull(r) {
			break
		}
	}

	// Scan until null, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if !isNull(r) {
			return i + width, data[start:i], nil
		}
	}

	// Request more data.
	return start, nil, nil
}

var num int
var nullMode bool

func init() {
	separator = seperatorSpace
	pflag.IntVarP(&num, "max-args", "n", 0,
		"Use at most max-args arguments per command line.")
	pflag.BoolVar(&nullMode, "0", false,
		"Input items are terminated by a null character instead of by whitespace, and the quotes and backslash are not special (every character is taken literally). Disables the end of file string, which is treated like any other argument. Useful when input items might contain white space, quote marks, or backslashes. The GNU find -print0 option produces input suitable for this mode.")

	pflag.Parse()
}

func main() {
	// init buffer
	stdin := bytes.Buffer{}
	_, err := stdin.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// check stdin args
	if nullMode {
		separator = seperatorNull
	}

	// check argument length
	stdinSS := split(stdin.String(), separator)
	if len(stdinSS) > num {
		log.Fatal("too many arguments")
	}

	out, err := CoreOp(pflag.Args(), &stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(out)
	// for every entry from stdin, concat with xargs arg and execute
}

func CoreOp(args []string, stdin io.Reader) (string, error) {
	lenArgs := len(args)

	var res string
	switch {
	case lenArgs == 0:
		res = process("echo", []string{}, stdin)
	case lenArgs == 1:
		res = process(args[0], []string{}, stdin)
	case lenArgs > 1:
		res = process(args[0], args[1:], stdin)
	}

	return res, nil
}

func process(command string, args []string, stdin io.Reader) string {
	var res = ""
	var lenBytesResp int64

	// read into bytes
	scanner := bufio.NewScanner(stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		var resp = []byte{}
		var err error
		argTxt := scanner.Text()

		// process words
		var concatArgs = []string{}
		concatArgs = append(append(concatArgs, args...), argTxt)
		lenBytesResp, resp, err = run(command, concatArgs)
		if err != nil {
			fmt.Println("command failed with:", err.Error())
		}
		res = join(res, string(resp))
	}

	// process output
	if lenBytesResp > 0 {
		res += "\n"
	}

	return res
}

func run(command string, args []string) (int64, []byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		return 0, nil, err
	}

	var returnBytes []byte
	if len(b) > 1 && b[len(b)-1] == '\n' {
		// b[:len(b)-1] removes newlines from return bytes
		returnBytes = b[:len(b)-1]
	} else {
		returnBytes = b
	}
	return int64(len(returnBytes)), returnBytes, nil
}

func join(s1, s2 string) string {
	if len(s1) == 0 {
		return s2
	} else if len(s2) == 0 {
		return s1
	}
	//fmt.Println("bytes:", []byte(s2))
	return s1 + " " + s2
}

func isNull(r rune) bool {
	if r == null {
		return true
	}
	return false
}

func split(s string, separator int) []string {
	seps := seperatorMap[separator]
	// unwrap seps into map for easy retrieval
	mapSeps := make(map[rune]bool, len(seps))
	for i := 0; i < len(seps); i++ {
		mapSeps[seps[i]] = true
	}

	// init variables
	mutS := []rune(s)
	var res = []string{}

	lastSepIndex := -1 // ensure it doesn't conflate with the actual index 0
	for i := 0; i < len(mutS); i++ {
		if _, ok := mapSeps[mutS[i]]; ok {
			// return early if index equals zero,
			//it is the same as the former sep and they follow each other
			if !(i == 0 || i == lastSepIndex+1) {
				res = append(res, string(mutS[lastSepIndex+1:i]))
			}
			lastSepIndex = i
		}
	}
	if lastSepIndex < len(mutS)-1 {
		res = append(res, string(mutS[lastSepIndex+1:]))
	}
	return res
}

//func isUniformUpToIndex(rr []rune, start int, target int) bool {
//	if len(rr) == 0 {
//		return false
//	}
//
//	if start >= len(rr)-1 || target >= len(rr)-1 {
//		return false
//	}
//
//	if start == target {
//		return true
//	}
//
//	var flag = rr[0]
//	for i := start; i <= target; i++ {
//		if flag != rr[i] {
//			return false
//		}
//	}
//	return true
//}
