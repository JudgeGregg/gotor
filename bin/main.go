package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/JudgeGregg/gotor/globals"
	"github.com/JudgeGregg/gotor/parser"
	"github.com/JudgeGregg/gotor/raid"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {

	debug := flag.Bool("d", false, "debug")
	argFile := flag.String("f", "", "file")
	mainPlayer := flag.String("p", "", "main player")
	flag.Parse()
	globals.Debug = *debug
	if *argFile == "" || *mainPlayer == "" {
		flag.Usage()
		os.Exit(0)
	}
	_, filename := path.Split(*argFile)
	globals.RaidStartDate = raid.GetRaidStartDate(filename)
	globals.MainPlayerName = *mainPlayer
	lines := make(chan string)
	records := make(chan parser.Record)
	go tail(*argFile, lines)
	go parser.Parse(lines, records)
	raid_ := &parser.Raid{}
	for record := range records {
		raid.HandleRecord(raid_, record)
	}
}

func tail(filename string, lines chan string) {
	f, err := os.Open(filename)
	transformer := charmap.ISO8859_1.NewDecoder()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	oldSize := info.Size()
	tempLine := ""
	str := ""
	for {
		for line, err := r.ReadString('\n'); ; line, err = r.ReadString('\n') {
			fmt.Print("baseline: ", line)
			if err == io.EOF {
				// we have a "temp" line
				if line != "" {
					// not empty
					if !strings.Contains(line, "\n") {
						// partial line, wait
						tempLine += line
						fmt.Print("tempLine: ", tempLine)
						continue
					} else {
						fmt.Print("end of line!")
						// end of partial line, proceed
						tempLine += line
						fmt.Print("tempLine: ", tempLine)
					}
				} else {
					//empty temp line = "EOF"
					break
				}
			}
			if tempLine != "" {
				tempLine += line
				fmt.Print("tempLine: ", tempLine)
				str, _, _ = transform.String(transformer, string(tempLine))
				tempLine = ""
			} else {
				fmt.Print("line: ", line)
				str, _, _ = transform.String(transformer, string(line))
			}
			lines <- str
		}
		fmt.Print("DONE")
		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			panic(err)
		}
		for {
			time.Sleep(time.Second)
			newinfo, err := f.Stat()
			if err != nil {
				panic(err)
			}
			newSize := newinfo.Size()
			if newSize != oldSize {
				if newSize < oldSize {
					f.Seek(0, 0)
				} else {
					f.Seek(pos, io.SeekStart)
				}
				r = bufio.NewReader(f)
				oldSize = newSize
				break
			}
		}
	}
}
