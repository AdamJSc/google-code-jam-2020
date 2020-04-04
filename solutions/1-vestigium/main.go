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

type inOut struct {
	in  *bufio.Scanner
	out *bufio.Writer
}

func (i inOut) read() (string, error) {
	if !i.in.Scan() {
		if err := i.in.Err(); err != nil {
			return "", err
		}
		return "", errors.New("end of stream")
	}

	return i.in.Text(), nil
}

func (i inOut) write(s solution) error {
	fmt.Fprintf(i.out, "Case #%d: %s\n", s.caseNum, s.output)
	return i.out.Flush()
}

var stream = inOut{in: bufio.NewScanner(os.Stdin), out: bufio.NewWriter(os.Stdout)}

func readNumOfTestCases(stream inOut) (int, error) {
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

func main() {
	tc, err := readNumOfTestCases(stream)
	if err != nil {
		panic(err)
	}

	for i := 1; i <= tc; i++ {
		if err = solve(i, stream); err != nil {
			panic(err)
		}
	}
}

func solve(caseNum int, stream inOut) error {
	matrixSize, err := readMatrixSize(stream)
	if err != nil {
		return err
	}

	var matrix matrix
	for i := 1; i <= matrixSize; i++ {
		inputRow, err := stream.read()
		if err != nil {
			return err
		}

		rowAsInts, err := parseRowAsInts(inputRow)
		if err != nil {
			return err
		}

		matrix.appendRow(rowAsInts)
	}

	stream.write(solution{caseNum: caseNum, output: fmt.Sprintf("%d 0 0", matrix.getTrace())})
	return nil
}

func readMatrixSize(stream inOut) (int, error) {
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

func parseRowAsInts(row string) ([]int, error) {
	var rowAsInts []int

	for _, strNum := range strings.Split(row, " ") {
		intNum, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			return []int{}, err
		}

		rowAsInts = append(rowAsInts, int(intNum))
	}

	return rowAsInts, nil
}

type matrix struct {
	rows [][]int
}

func (m *matrix) appendRow(row []int) {
	m.rows = append(m.rows, row)
}

func (m *matrix) getTrace() int {
	var trace int
	for i, row := range m.rows {
		trace += row[i]
	}
	return trace
}
