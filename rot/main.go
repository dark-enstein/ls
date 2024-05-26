package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	Alpha = map[string]string{
		"A": "N",
		"B": "O",
		"C": "P",
		"D": "Q",
		"E": "R",
		"F": "S",
		"G": "T",
		"H": "U",
		"I": "V",
		"J": "W",
		"K": "X",
		"L": "Y",
		"M": "Z",
		"N": "A",
		"O": "B",
		"P": "C",
		"Q": "D",
		"R": "E",
		"S": "F",
		"T": "G",
		"U": "H",
		"V": "I",
		"W": "J",
		"X": "K",
		"Y": "L",
		"Z": "M",
		"a": "n",
		"b": "o",
		"c": "p",
		"d": "q",
		"e": "r",
		"f": "s",
		"g": "t",
		"h": "u",
		"i": "v",
		"j": "w",
		"k": "x",
		"l": "y",
		"m": "z",
		"n": "a",
		"o": "b",
		"p": "c",
		"q": "d",
		"r": "e",
		"s": "f",
		"t": "g",
		"u": "h",
		"v": "i",
		"w": "j",
		"x": "k",
		"y": "l",
		"z": "m",
	}
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <filename>")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(cipherline(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func cipherline(s string) string {
	var result string
	for _, c := range s {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			result += Alpha[string(c)]
		} else {
			result += string(c)
		}
	}
	return result
}
