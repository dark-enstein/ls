package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	stdin := bytes.Buffer{}
	_, err := stdin.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	out, err := CoreOp(os.Args[1:], &stdin)
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
		res = parseStdin("echo", []string{}, stdin)
	case lenArgs == 1:
		res = parseStdin(args[0], []string{}, stdin)
	case lenArgs > 1:
		res = parseStdin(args[0], args[1:], stdin)
	}

	return res, nil
}

func parseStdin(command string, args []string, stdin io.Reader) string {
	var res = ""
	var lenBytesResp int64
	// read into bytes
	scanner := bufio.NewScanner(stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		var resp = []byte{}
		//var err error
		argTxt := scanner.Text()
		//log.Println("Processing:", scanner.Text())
		// trim extraeneous bytes, extra space, empty lines, etc
		//strings.Trim
		// process words
		var concatArgs = []string{}
		concatArgs = append(append(concatArgs, args...), argTxt)
		lenBytesResp, resp, _ = run(command, concatArgs)
		//if err != nil {
		//	fmt.Println("Command failed with:", err.Error())
		//}
		res = join(res, string(resp))
	}

	// process output
	if lenBytesResp > 0 {
		res += "\n"
	}

	return res
}

func run(command string, args []string) (int64, []byte, error) {
	//log.Println("Command:", command)
	//log.Println("Args:", args)
	cmd := exec.Command(command, args...)
	//cmd.Stdout = os.Stdout
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
