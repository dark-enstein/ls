package main

import (
	"fmt"
	"strings"
)

var (
	message = "MD5 is one in a series of message digest algorithms designed by Professor Ronald Rivest of MIT (Rivest, 1992). When analytic work indicated that MD5's predecessor MD4 was likely to be insecure, Rivest designed MD5 in 1991 as a secure replacement. (Hans Dobbertin did indeed la "
)

func main() {
	fmt.Printf("message: %s\n\n", message)
	fmt.Println(_main(message))
}

func _main(s string) string {
	// convert input string to binary
	bin := _convToBinary(s)

	// pad message as defined in RFC
	padded_bin := _pad(bin)

	// init constants
	cA, cB, cC, cD := _initConstants()

	return padded_bin
}

func F(X, Y, Z uint32) uint32 {}
func G(X, Y, Z uint32) uint32 {}
func H(X, Y, Z uint32) uint32 {}
func I(X, Y, Z uint32) uint32 {}

func _pad(bin string) string {
	// pad message as defined in RFC
	lengthMessage := len(bin)
	//fmt.Println("message length:", lengthMessage)

	rem := 512 - (lengthMessage % 512)
	//fmt.Println("needed padding length:", rem)
	padding := ""
	if rem != 0 {
		for i := 0; rem > 64; rem-- {
			// pad with one 1
			if i == 0 {
				padding += "1"
				i++
				continue
			}

			// then pad with zeros
			padding += "0"
			i++
		}
		//fmt.Println("after padding:", padding)
		//fmt.Println("len padding:", len(padding))
		// length remaining should be exactly 64
		if rem != 64 {
			fmt.Println("error calculating padding stuff")
		}

		// pad the original length
		validLength := fmt.Sprintf("%b", lengthMessage)
		leftPaddedTo64 := strings.Repeat("0", 64-len(validLength))
		totalMessageLengthPadding := leftPaddedTo64 + validLength
		//fmt.Println("length of totalMessageLengthPadding:", len(totalMessageLengthPadding))

		bin += padding + totalMessageLengthPadding
	}

	fmt.Println("length after padding: ", len(bin))
	if len(bin)%512 != 0 {
		fmt.Println("error calculating padding stuff")
	}

	return bin
}

func _convToBinary(s string) string {
	bins := ""
	for i := 0; i < len(s); i++ {
		bins += fmt.Sprintf("%b", s[i])
	}
	return bins
}

func _initConstants() ([][]int, [][]int, [][]int, [][]int) {
	a := [][]int{
		{
			0, 1,
		},
		{
			2, 3,
		},
		{
			4, 5,
		},
		{
			6, 7,
		},
	}
	b := [][]int{
		{
			8, 9,
		},
		{
			10, 11,
		},
		{
			12, 13,
		},
		{
			14, 15,
		},
	}
	c := [][]int{
		{
			15, 14,
		},
		{
			13, 12,
		},
		{
			11, 10,
		},
		{
			9, 8,
		},
	}
	d := [][]int{
		{
			7, 6,
		},
		{
			5, 4,
		},
		{
			3, 2,
		},
		{
			1, 0,
		},
	}
	return a, b, c, d
}
