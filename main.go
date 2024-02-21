package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	removeOverlap := flag.Bool("r", false, "Remove cue overlap if specified (ignore end time)")
	intOffset := flag.Int("o", 0, "Offset cue timestamps by N [ms]")
	intDefaultCueLength := flag.Int("l", 15000, "Max/default cue length if not specified [ms]")
	flag.Parse()
	args := flag.Args()

	offset, defaultCueLength := Timestamp(*intOffset), Timestamp(*intDefaultCueLength)

	var inFileName, outFileName string
	var rf *os.File
	switch {
	case len(args) == 0: // read stdin
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

	cues, err := readHeur(rf)
	if err != nil {
		panic("Heuristic file reading error")
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
	TransformCues(cues, offset, defaultCueLength, *removeOverlap)
	WriteVTT(w, cues)
}

func readHeur(rf io.Reader) ([]Cue, error) {
	fileScanner := bufio.NewScanner(rf)
	fileScanner.Split(bufio.ScanLines)
	// TODO ignore SRT cue numbers
	reIsTSLine := regexp.MustCompile("^\\d+:\\d+")

	cues := make([]Cue, 0)
	for fileScanner.Scan() {
		text := strings.TrimSpace(fileScanner.Text())
		if len(text) < 1 || (len(cues) == 0 && text == "WEBVTT") {
			continue
		}
		if reIsTSLine.MatchString(text) {
			tss := getTSs(text)
			fmt.Printf("TSS(%d): %q\n", len(tss), tss)

			if len(tss) > 2 || len(tss) == 0 {
				return nil, fmt.Errorf("Invalid timestamp format: %s", text)
			}

			var begin, end Timestamp = -1, -1
			var err, err2 error
			if len(tss) == 2 {
				end, err2 = getMs(tss[1])
			}

			begin, err = getMs(tss[0])

			if err != nil || err2 != nil {
				panic("Invalid timestamp format: " + text)
			}
			cues = append(cues, Cue{begin: begin, end: end, text: ""})

		} else {
			if len(cues) >= 1 {
				cues[len(cues)-1].text += text + "\n"
			}
		}
	}
	return cues, nil
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func WriteVTT(w *bufio.Writer, cues []Cue) {
	defer w.Flush()
	w.WriteString("WEBVTT\n\n")

	for _, c := range cues {
		fmt.Fprintf(w, "%s --> %s\n%s\n", c.begin, c.end, c.text)
	}
}

func WriteSRT(w *bufio.Writer, cues []Cue) {
	defer w.Flush()

	for i, c := range cues {
		fmt.Fprintf(w, "%d\n%s --> %s\n%s\n", i, c.begin, c.end, c.text)
	}
}
