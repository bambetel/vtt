package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Convert timestamp string (hh:)mm:ss.ttt to milliseconds
func getMs(str string) (Timestamp, error) {
	spl := strings.Split(str, ".")
	if len(spl) > 2 {
		return -1, fmt.Errorf("Too much dots")
	}
	var ms int
	var err error
	if len(spl) == 2 {
		ms, err = strconv.Atoi(spl[1])
		if err != nil {
			return -1, fmt.Errorf("Invalid milliseconds format")
		}
	}
	spl = strings.Split(spl[0], ":")
	if l := len(spl); l < 2 || l > 3 {
		return -1, fmt.Errorf("Invalid timestamp format")
	}
	var t [3]int

	for i := range spl {
		t[3-len(spl)+i], err = strconv.Atoi(spl[i])
		if err != nil {
			return -1, fmt.Errorf("Invalid timestamp format")
		}
	}
	return Timestamp(((t[0]*60+t[1])*60+t[2])*1000 + ms), nil
}

type Timestamp int

func (ts Timestamp) String() string {
	if ts == -1 {
		return "-1"
	}
	ms := int64(ts)
	t := ms % 1000
	ms /= 1000
	h := ms / 3600
	ms %= 3600
	m := ms / 60
	s := ms % 60

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, t)
	} else {
		return fmt.Sprintf("%02d:%02d.%03d", m, s, t)
	}
}

func (ts Timestamp) Offset(o int) Timestamp {
	return max(0, ts+Timestamp(o))
}

type Cue struct {
	begin Timestamp
	end   Timestamp
	text  string
}

// split line to timestamps (1 or more)
func getTSs(line string) []string {
	if strings.ContainsRune(line, ',') {
		return strings.Split(line, ",")
	}
	if strings.Contains(line, "-") {
		reTSSep := regexp.MustCompile("\\s*-+>?\\s*")
		return reTSSep.Split(line, -1)
	}
	return []string{line}
}

// Modifies cues slice
func TransformCues(cues []Cue, offset Timestamp, defaultCueLength Timestamp, removeOverlap bool) {
	for i := range cues {
		if removeOverlap && i < len(cues)-1 {
			cues[i].end = min(cues[i].end, cues[i+1].begin)
		}
		// extrapolate cue end, no overlap
		if cues[i].end == -1 {
			if i < len(cues)-1 {
				cues[i].end = min(cues[i+1].begin, cues[i].begin+defaultCueLength)
			} else {
				cues[i].end = cues[i].begin + defaultCueLength
			}
		}
		cues[i].begin += offset
		cues[i].end += offset
	}
}
