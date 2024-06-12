# Subtitle conversion utility

## Installation

Requires GO 1.21+ to run.

## Usage

    vtt <options> <input-file> <output-file>

Files:
- If no file is specified, standard input is read.
- If no output file is specified, the result is printed to standard output.

CLI options:
- `-o=N` Offset all timestamps by N [ms]
- `-r`   Remove cue overlap if specified (ignore end time and set to next cue start time)
- `-l=N`    Set max/default cue length if not specified to N [ms], default: 15s
- `-t=srt` Set output format: vtt or srt (default: vtt)

## Input data

This program reads simple text files with subtitles and converts them to VTT. 

Subtitles in the following format specifying start of the caption and end in format minutes:seconds-minutes:seconds.

    0:00-0:05
    First subtitle line.
    0:06-0:10 
    Second subtitle line.
    0:12-0:20
    Third subtitle part.

If any end time will be ommited, the subtitle duration will be set to the begining of the next line.

     0:00                    -->  0:00-0:06
     First subtitle line.    -->  First subtitle line.
     0:06                    -->  0:06-0:12
     Second subtitle line.   -->  Second subtitle line.
     0:12-...                -->  0:12-...
     Third subtitle part.    -->  Third subtitle part.
                                            
