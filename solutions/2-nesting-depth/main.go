package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type solution struct {
	caseNum int
	output  string
}

type ioStream struct {
	in  *bufio.Scanner
	out *bufio.Writer
}

func (i ioStream) read() (string, error) {
	if !i.in.Scan() {
		if err := i.in.Err(); err != nil {
			return "", err
		}
		return "", errors.New("end of stream")
	}

	return i.in.Text(), nil
}

func (i ioStream) readInt() (int, error) {
	inp, err := stream.read()
	if err != nil {
		return 0, err
	}

	int64, err := strconv.ParseInt(inp, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(int64), nil
}

func (i ioStream) write(s solution) error {
	fmt.Fprintf(i.out, "Case #%d: %s\n", s.caseNum, s.output)
	return i.out.Flush()
}

var stream = ioStream{in: bufio.NewScanner(os.Stdin), out: bufio.NewWriter(os.Stdout)}

func main() {
	numOfTestCases, err := stream.readInt()
	if err != nil {
		panic(err)
	}

	for i := 1; i <= numOfTestCases; i++ {
		if err = solve(i, stream); err != nil {
			panic(err)
		}
	}
}

func solve(caseNum int, stream ioStream) error {
	seq, err := readInts(stream)
	if err != nil {
		return err
	}

	var symbols []symbol
	for _, char := range seq {
		symbols = append(symbols, symbol{value: char})
	}

	populateRequiredParentheses(symbols)

	var output string
	for _, symbol := range symbols {
		output = output + symbol.toString()
	}

	stream.write(solution{caseNum: caseNum, output: output})
	return nil
}

type symbol struct {
	value              int
	openingParentheses int
	closingParentheses int
}

func (s symbol) toString() string {
	var prefix, suffix string

	for i := 0; i < s.openingParentheses; i++ {
		prefix = prefix + "("
	}

	for i := 0; i < s.closingParentheses; i++ {
		suffix = suffix + ")"
	}

	return fmt.Sprintf("%s%d%s", prefix, s.value, suffix)
}

func readInts(stream ioStream) ([]int, error) {
	inp, err := stream.read()
	if err != nil {
		return []int{}, err
	}

	var ints []int
	for _, char := range strings.Split(inp, "") {
		charAsInt, err := strconv.ParseInt(char, 10, 64)
		if err != nil {
			return []int{}, err
		}

		ints = append(ints, int(charAsInt))
	}

	return ints, nil
}

func populateRequiredParentheses(symbols []symbol) {
	var unresolved int
	for i := 0; i < len(symbols); i++ {
		// required number of opening parentheses is the difference between
		// the current value and the count of unresolved parentheses
		symbols[i].openingParentheses = symbols[i].value - unresolved

		// determine required number of closed parentheses
		if i+1 == len(symbols) {
			// no more symbols after this one, so closing parentheses should
			// wrap up the current value nicely
			symbols[i].closingParentheses = symbols[i].value
			continue
		}
		switch {
		case symbols[i+1].value < symbols[i].value:
			// next value is smaller, so we need to close off the difference
			// otherwise we will have too many unresolved parentheses for the next value!
			symbols[i].closingParentheses = symbols[i].value - symbols[i+1].value
		}

		// update the number of unresolved parentheses that we now have
		unresolved = unresolved + symbols[i].openingParentheses - symbols[i].closingParentheses
	}
}
