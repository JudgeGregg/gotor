package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"path"
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
	if *argFile == "" {
		argFile = &os.Args[1]
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
	for {
		for line, prefix, err := r.ReadLine(); err != io.EOF; line, prefix, err = r.ReadLine() {
			if prefix {
				str, _, _ := transform.String(transformer, string(line))
				lines <- str
			} else {
				str, _, _ := transform.String(transformer, string(line))
				lines <- str
			}
		}
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
