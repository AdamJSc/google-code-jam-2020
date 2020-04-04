package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
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
	numOfActivities, err := stream.readInt()
	if err != nil {
		return err
	}

	for i := 0; i < numOfActivities; i++ {
		// pairs of activity start and end time on their way
		stream.read()
	}

	stream.write(solution{caseNum: caseNum, output: fmt.Sprintf("%d", numOfActivities)})
	return nil
}
