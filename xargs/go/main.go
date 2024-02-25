package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
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
var seperator = 0

var ScanNull = func(data []byte, atEOF bool) (advance int,
	token []byte,
	err error) {
	// Skip leading null bytes.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !(r != null) {
			break
		}
	}

	// Scan until null, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if !(r != null) {
			return i + width, data[start:i], nil
		}
	}

	// Request more data.
	return start, nil, nil
}

var num int
var procs int
var nullMode bool

func init() {
	seperator = seperatorSpace
	pflag.IntVarP(&num, "max-args", "n", 1,
		"Use at most max-args arguments per command line.")
	pflag.BoolVar(&nullMode, "0", false,
		"Input items are terminated by a null character instead of by whitespace, and the quotes and backslash are not special (every character is taken literally). Disables the end of file string, which is treated like any other argument. Useful when input items might contain white space, quote marks, or backslashes. The GNU find -print0 option produces input suitable for this mode.")
	pflag.IntVarP(&procs, "max-procs", "P", 1,
		"Run up to max-procs processes at a time; the default is 1. If max-procs is 0, xargs will run as many processes as possible at a time. Use the -n option with -P; otherwise chances are that only one exec will be done.")

	pflag.Parse()
}

func main() {
	//num = 1 // for testing
	if num == 0 {
		fmt.Println("cargs: -n 0: too small")
		os.Exit(1)
	}

	// init buffer
	stdin := bytes.Buffer{}
	_, err := stdin.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// check stdin args
	if nullMode {
		seperator = seperatorNull
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
	// using map to keep the order of results
	var results = make(map[int]string)
	var index = 0
	// ensure that more go routines than procs is not run at any time
	var currentProcsChan = make(chan int, procs)
	var wg sync.WaitGroup
	var m sync.Mutex

	// read into bytes until valid sep
	var batch = []string{}
	scanner := bufio.NewScanner(stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		argTxt := scanner.Text()

		// add arguments to the args from duplicated until num
		m.Lock()
		batch = append(batch, argTxt)
		m.Unlock()
		// execute when the additional commands added to the replicated args
		//equals the max-args
		if len(batch) == num {
			index++
			currentProcsChan <- 1 // should block until space is available
			wg.Add(1)
			go func(i int, batchArgs []string) {
				defer wg.Done()
				result, err := _proc(command, append(args, batchArgs...))
				if err != nil {
					fmt.Println("command failed with:", err.Error())
				}
				m.Lock()
				results[i] = result
				// reset concatenated args
				//concatArgs = args
				trash := <-currentProcsChan // remove from channel so make space
				// for new process
				_ = trash
				m.Unlock()
			}(index, batch)
			batch = []string{} // reset
		}
	}

	// handle split buffer length that isn't a mod of num.
	//when the num is lesser or greater than the split contents of stdin,
	//flushing the remaining contents
	if len(batch) > 0 {
		index++
		currentProcsChan <- 1
		wg.Add(1)
		go func(i int, batchArgs []string) {
			defer wg.Done()
			result, err := _proc(command, append(args, batchArgs...))
			if err != nil {
				fmt.Println("command failed with:", err.Error())
			}
			m.Lock()
			results[i] = result
			// reset concatenated args
			//concatArgs = args
			d := <-currentProcsChan
			_ = d
			m.Unlock()
		}(index, batch)
	}
	wg.Wait()

	return join(results)
}

func _proc(command string, concatArgs []string) (string, error) {
	var respBytes = []byte{}
	var err error
	_, respBytes, err = run(command, concatArgs)
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func run(command string, args []string) (int64, []byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		return 0, nil, err
	}

	var returnBytes []byte
	returnBytes = b
	return int64(len(returnBytes)), returnBytes, nil
}

func join(ss map[int]string) string {
	if len(ss) == 0 {
		return ""
	}

	keys := make([]int, 0, len(ss))
	for k := range ss {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var s strings.Builder
	for _, k := range keys {
		s.WriteString(ss[k])
	}

	return s.String()
}

// implement quick sort instead of sort.Int
func quicksort(s []string) []string { return s }

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
			//it is the same as the former sep, and they follow each other
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
