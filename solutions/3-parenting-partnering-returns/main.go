package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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

	var schedule activitySchedule

	for i := 0; i < numOfActivities; i++ {
		input, err := stream.read()
		if err != nil {
			return err
		}

		pairOfMinutesAfterMidnight := strings.Split(input, " ")
		if len(pairOfMinutesAfterMidnight) != 2 {
			return fmt.Errorf("cannot split '%s' into a pair of integers", input)
		}

		startMinute, err := strconv.ParseInt(pairOfMinutesAfterMidnight[0], 10, 64)
		if err != nil {
			return err
		}

		endMinute, err := strconv.ParseInt(pairOfMinutesAfterMidnight[1], 10, 64)
		if err != nil {
			return err
		}

		activity, err := newActivityFromMinutes(startMinute, endMinute)
		if err != nil {
			return err
		}

		schedule.activities = append(schedule.activities, activity)
	}

	stream.write(solution{caseNum: caseNum, output: fmt.Sprintf("%+v", schedule)})
	return nil
}

type timespan struct {
	start time.Time
	end   time.Time
}

type parent struct {
	occupied []timespan
}

type activity struct {
	timespan      timespan
	parentInitial string
}

func newActivityFromMinutes(start int64, end int64) (activity, error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return activity{}, err
	}

	midnight := time.Date(1970, 1, 1, 0, 0, 0, 0, loc)
	startTime := midnight.Add(time.Duration(start) * time.Minute)
	endTime := midnight.Add(time.Duration(end) * time.Minute)

	return activity{timespan: timespan{start: startTime, end: endTime}}, nil
}

type activitySchedule struct {
	activities []activity
	parents    map[string]parent
}
