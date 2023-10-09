package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getMs(str string) int {
	reMST := regexp.MustCompile("^(\\d+):(\\d+)\\.?(\\d+)?")
	res := reMST.FindStringSubmatch(str)
	if res == nil {
		fmt.Println("No match (1)!")
		return -1
	}
	m, err := strconv.Atoi(res[1])
	s, err2 := strconv.Atoi(res[2])
	if err != nil || err2 != nil {

		fmt.Println("Atoi error (1)!")
		return -1
	}
	ms := 0
	if len(res) > 3 {
		if res[3] != "" {

			mili, err := strconv.Atoi(res[3])
			if err == nil {
				ms = mili
			} else {
				fmt.Println("ms present but no int")
			}
		}
	}
	return (m*60+s)*1000 + ms
}

func formatTimestamp(ms int) string {
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

type Cue struct {
	begin int
	end   int
	text  string
}

func main() {
	var offset int
	var defaultCueLength int
	var removeOverlap bool
	flag.IntVar(&offset, "o", 0, "Offset cue timestamps by N [ms]")
	flag.IntVar(&defaultCueLength, "l", 15000, "Max/default cue length if not specified [ms]")
	flag.BoolVar(&removeOverlap, "r", false, "Remove cue overlap if specified (ignore end time)")
	flag.Parse()
	fmt.Printf("flags o=%d\n", offset)
	rf, err := os.Open("in.txt")
	defer rf.Close()
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(rf)

	fileScanner.Split(bufio.ScanLines)

	reTSLine := regexp.MustCompile("^\\d+:.*\\d+$")

	cues := make([]Cue, 0)
	for fileScanner.Scan() {
		// sort cues here
		text := strings.TrimSpace(fileScanner.Text())
		if len(text) < 1 || (len(cues) == 0 && text == "WEBVTT") {
			continue
		}
		if reTSLine.MatchString(text) {
			sp := strings.Split(text, "-")
			begin := getMs(sp[0])
			end := -1
			if len(sp) > 1 {
				end = getMs(sp[1])
			}
			cues = append(cues, Cue{begin: begin, end: end, text: ""})

		} else {
			cues[len(cues)-1].text += text + "\n"
		}
	}
	// output
	if len(cues) < 1 {
		fmt.Println("No cues found in input, nothing to output.")
		return
	}

	f, err := os.Create("out.txt")
	defer f.Close()
	w := bufio.NewWriter(f)

	w.WriteString("WEBVTT\n\n")

	for i, c := range cues {
		if removeOverlap && i < len(cues)-1 {
			c.end = min(c.end, cues[i+1].begin)
		}
		// extrapolate cue end, no overlap
		if c.end == -1 {
			if i < len(cues)-1 {
				c.end = min(cues[i+1].begin, c.begin+defaultCueLength)
			} else {
				c.end = c.begin + defaultCueLength
			}
		}

		w.WriteString(
			formatTimestamp(c.begin+offset) + " --> " +
				formatTimestamp(c.end+offset) + "\n" +
				c.text + "\n")
	}
	w.Flush()
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
