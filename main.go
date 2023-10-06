package main

import (
	"bufio"
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
	bT := ms % 1000
	ms /= 1000
	bM := ms / 60
	bS := ms % 60

	return fmt.Sprintf("%02d:%02d.%03d", bM, bS, bT)
}

type Cue struct {
	begin int
	end   int
	text  string
}

func main() {
	offset := 8000 // offset timestamps
	fmt.Println("vim-go")
	rf, err := os.Open("in.txt")
	defer rf.Close()
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("out.txt")
	defer f.Close()
	w := bufio.NewWriter(f)
	fileScanner := bufio.NewScanner(rf)

	fileScanner.Split(bufio.ScanLines)

	// reMST := regexp.MustCompile("^(\\d+):(\\d+)\\.?(\\d+)?")
	reTSLine := regexp.MustCompile("^\\d+.*\\d+$")

	w.WriteString("WEBVTT\n")
	lines := make([]string, 0)
	cues := make([]Cue, 0)
	for fileScanner.Scan() {
		// sort lines and cues here
		text := strings.TrimSpace(fileScanner.Text())
		lines = append(lines, text)
		if reTSLine.MatchString(text) {
			fmt.Println("cue line: ", text)
			sp := strings.Split(text, "-")
			fmt.Printf("Split: %q\n", sp)
			begin := getMs(sp[0])
			end := -1
			if len(sp) > 1 {
				end = getMs(sp[1])
			}
			cues = append(cues, Cue{begin: begin, end: end, text: ""})

		} else {
			fmt.Println("sub line: ", text)
			cues[len(cues)-1].text += text
		}
	}
	fmt.Println("------")
	fmt.Println(cues)
	fmt.Println("------")
	for i, c := range cues {

		if c.end == -1 {
			if i < len(cues)-1 {
				c.end = cues[i+1].begin
			} else {
				c.end = c.begin + 15000
			}
		}
		w.WriteString("\n" + formatTimestamp(c.begin+offset) + " --> " + formatTimestamp(c.end+offset) + "\n" + c.text + "\n")
	}
	w.Flush()
}
