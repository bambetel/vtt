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
	i := 0

	// reMST := regexp.MustCompile("^(\\d+):(\\d+)\\.?(\\d+)?")
	reTSLine := regexp.MustCompile("^\\d+.*\\d+$")

	w.WriteString("WEBVTT\n")
	for fileScanner.Scan() {
		text := strings.TrimSpace(fileScanner.Text())
		if reTSLine.MatchString(text) {
			fmt.Println("cue line: ", text)
			// single number - extrapolate
			w.WriteString("\n" + formatTimestamp(getMs(text)+offset) + " --> " + formatTimestamp(getMs(text)+offset+3000) + "\n")
		} else {
			// res := reMST.FindStringSubmatch(text)
			// fmt.Println("re find: ", res)

			fmt.Println("sub line: ", text)
			w.WriteString(text + "\n")
		}
		i++
	}
	w.Flush()
}
