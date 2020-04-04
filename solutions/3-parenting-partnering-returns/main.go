package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
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

	sched := schedule{
		parents: map[string]parent{
			"cameron": parent{initial: "C"},
			"jamie":   parent{initial: "J"},
		},
	}

	for i := 0; i < numOfActivities; i++ {
		// parse start and end minutes from each forthcoming input row
		// and inflate our activity schedule
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
		sched.activities = append(sched.activities, activity)
	}

	if err = assignParentsToActivities(&sched); err != nil {
		return err
	}

	output := sched.toString()
	stream.write(solution{caseNum: caseNum, output: output})
	return nil
}

type timespan struct {
	start time.Time
	end   time.Time
}

func (ts timespan) overlapsWith(t timespan) bool {
	// check if end of t is equal to or before start of ts
	if t.end.Equal(ts.start) || t.end.Before(ts.start) {
		return false
	}

	// check if start of t is equal to or after end of ts
	if t.start.Equal(ts.end) || t.start.After(ts.end) {
		return false
	}

	return true
}

type parent struct {
	initial   string
	timetable []timespan
}

func (p parent) isAvailableFor(t timespan) bool {
	for _, timespan := range p.timetable {
		if timespan.overlapsWith(t) {
			return false
		}
	}

	return true
}

type activity struct {
	timespan      timespan
	parentInitial string
}

func (a activity) getRef() string {
	start := a.timespan.start.Format("150405")
	end := a.timespan.end.Format("150405")
	return fmt.Sprintf("%s:%s", start, end)
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

type schedule struct {
	activities []activity
	parents    map[string]parent
}

func (s schedule) toString() string {
	var output string
	for _, activity := range s.activities {
		if activity.parentInitial == "" {
			return "IMPOSSIBLE"
		}
		output = output + activity.parentInitial
	}
	return output
}

func assignParentsToActivities(s *schedule) error {
	// let's map each activity's ref to its slice index
	// this way we can iterate over the activity refs in alphabetical order
	activityKeyMap := make(map[string]int)
	var activityKeyMapKeys []string
	for idx, activity := range s.activities {
		activityKeyMap[activity.getRef()] = idx
		activityKeyMapKeys = append(activityKeyMapKeys, activity.getRef())
	}

	// sort map keys alphabetically (which is chronological by start date because of the format)
	sort.Strings(activityKeyMapKeys)

	// iterate over the activities in chronological order
	// in order to identify an available parent to fulfil each one
	for _, mapkey := range activityKeyMapKeys {
		idx := activityKeyMap[mapkey]
		activity := s.activities[idx]
		activity.parentInitial = getInitialOfAvailableParentForActivity(activity, s.parents)
		s.activities[idx] = activity
	}

	return nil
}

func getInitialOfAvailableParentForActivity(a activity, parents map[string]parent) string {
	for key, parent := range parents {
		if parent.isAvailableFor(a.timespan) {
			parent.timetable = append(parent.timetable, a.timespan)
			parents[key] = parent
			return parent.initial
		}
	}

	return ""
}
