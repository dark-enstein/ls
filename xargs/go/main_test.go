package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/mitchellh/go-homedir"
)

var tableTest = []struct {
	buf         *bytes.Buffer
	args        []string
	expectedOut string
}{{
	bytes.NewBuffer([]byte("gamma epsilon thundero")),
	[]string{"echo"},
	"gamma epsilon thundero\n",
}}

func TestCoreOp(t *testing.T) {
	for i := 0; i < len(tableTest); i++ {
		out, err := CoreOp(tableTest[i].args, tableTest[i].buf)
		if err != nil {
			t.Fatal(err)
		}

		expectedOut := tableTest[i].expectedOut
		if out != expectedOut {
			t.Errorf("expected %s, got %s", expectedOut, out)
		}
	}
}

func TestCargs(t *testing.T) {
	runMake()
	bin, _ := homedir.Expand("~/.adhoc/bin/cargs")

	for i := 0; i < len(tableTest); i++ {
		cmd := exec.Command(bin, tableTest[i].args[0])
		cmd.Stdin = tableTest[i].buf
		b, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}

		out := string(b)
		fmt.Println("output:", out)
		expectedOut := tableTest[i].expectedOut
		if out != expectedOut {
			t.Errorf("expected %s, got %s", expectedOut, out)
		}
	}
}

func BenchmarkCargs(b *testing.B) {
	bin, err := homedir.Expand("~/.adhoc/bin/cargs")
	if err != nil {
		log.Fatalln("error determining homedir:", err)
	}
	var cmds = []string{
		bin, "xargs",
	}

	runMake()

	b.ResetTimer()
	// run benchmark
	for i := 0; i < len(tableTest); i++ {
		for j := 0; j < len(cmds); j++ {
			b.Run(cmds[j], func(b *testing.B) {
				cmd := exec.Command(cmds[j], tableTest[i].args[0])
				cmd.Stdin = tableTest[i].buf
				_, err := cmd.Output()
				if err != nil {
					log.Println(err)
				}
			})
		}
	}
}

func runMake() {
	_, err := exec.Command("make").Output()
	if err != nil {
		log.Fatalln(err)
	}
}
