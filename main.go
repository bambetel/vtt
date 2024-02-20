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

// Convert timestamp string (hh:)mm:ss.ttt to milliseconds
func getMs(str string) Timestamp {
	reMST := regexp.MustCompile("^(\\d+):(\\d+)\\.?(\\d+)?")
	res := reMST.FindStringSubmatch(str)
	if res == nil {
		return -1
	}
	m, err := strconv.Atoi(res[1])
	s, err2 := strconv.Atoi(res[2])
	if err != nil || err2 != nil {
		return -1
	}
	ms := 0
	if len(res) > 3 {
		if res[3] != "" {
			mili, err := strconv.Atoi(res[3])
			if err == nil {
				ms = mili
			} else {
				return -1
			}
		}
	}
	return Timestamp((m*60+s)*1000 + ms)
}

type Timestamp int

func (ts Timestamp) String() string {
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

type Cue struct {
	begin Timestamp
	end   Timestamp
	text  string
}

func main() {
	var removeOverlap bool

	flag.BoolVar(&removeOverlap, "r", false, "Remove cue overlap if specified (ignore end time)")
	intOffset := flag.Int("o", 0, "Offset cue timestamps by N [ms]")
	intDefaultCueLength := flag.Int("l", 15000, "Max/default cue length if not specified [ms]")
	flag.Parse()
	args := flag.Args()

	offset, defaultCueLength := Timestamp(*intOffset), Timestamp(*intDefaultCueLength)

	var inFileName, outFileName string
	var rf *os.File
	switch {
	case len(args) == 0:
	case len(args) > 1:
		inFileName = args[0]
		outFileName = args[1]
	case len(args) == 1:
		rf = os.Stdin // TODO
		inFileName = args[0]

	}

	rf, err := os.Open(inFileName)
	defer rf.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fileScanner := bufio.NewScanner(rf)
	fileScanner.Split(bufio.ScanLines)

	// full cue from - to
	// TODO ignore and pass VTT cue settings after timestamps?

	// split line at `-->` (VTT) or `,` (sbv)
	reCueLine := regexp.MustCompile("^(\\d+:\\d+\\S*)\\s*-+>?\\s*(\\S*\\d+)$")
	// just beginning time
	reTSLine := regexp.MustCompile("^\\d+:[^\\s-]*\\d+$")

	cues := make([]Cue, 0)
	for fileScanner.Scan() {
		// sort cues here
		text := strings.TrimSpace(fileScanner.Text())
		if len(text) < 1 || (len(cues) == 0 && text == "WEBVTT") {
			continue
		}
		if reCueLine.MatchString(text) {
			// fmt.Println("Cue line: ", text)
			m := reCueLine.FindAllStringSubmatch(text, -1)[0]
			// fmt.Printf("%q\n", m)
			begin := getMs(m[1])
			end := getMs(m[2])
			cues = append(cues, Cue{begin: begin, end: end, text: ""})
		} else if reTSLine.MatchString(text) {
			sp := strings.Split(text, "-")
			begin := getMs(sp[0])
			end := Timestamp(-1)
			cues = append(cues, Cue{begin: begin, end: end, text: ""})

		} else {
			cues[len(cues)-1].text += text + "\n"
		}
	}
	// output
	if len(cues) < 1 {
		fmt.Fprint(os.Stderr, "No cues found in input, nothing to output.")
		return
	}
	var w *bufio.Writer
	if len(outFileName) > 0 {
		f, err := os.Create(outFileName)
		defer f.Close()
		if err != nil {
			fmt.Printf("Cannot open output file \"%s\" for writing.", outFileName)
			return
		}
		w = bufio.NewWriter(f)
	} else {
		w = bufio.NewWriter(os.Stdout)
	}

	TransformCues(cues, offset, defaultCueLength, removeOverlap)
	WriteVTT(w, cues)
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
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
		fmt.Printf("Transform %d %s %s %s -> %s %s", i, cues[i].begin, cues[i].end, offset, cues[i].begin+offset, cues[i].end+offset)
		cues[i].begin += offset
		cues[i].end += offset
	}
	for i, _ := range cues {
		fmt.Printf("%3d: %s --> %s %s\n", i, cues[i].begin, cues[i].end, cues[i].text)
	}
}

func WriteVTT(w *bufio.Writer, cues []Cue) {
	defer w.Flush()
	w.WriteString("WEBVTT\n\n")

	for _, c := range cues {
		fmt.Fprintf(w, "%s ---> %s\n%s\n", c.begin, c.end, c.text)
	}
}
