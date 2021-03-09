package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type entry struct {
	path string
	offset uint32
	size uint32
}

func main() {
	flag.String("o", "", "Output file path.")
	flag.Parse()

	if flag.Arg(0) == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	entries := gatherEntries()
	outputFile := getOutputPath()

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Unable to open file for writing: ", outputFile)
		os.Exit(1)
	}

	writeHeader(file, entries)
	appendFiles(file, entries)

	file.Close()
}

func getOutputPath() string {
	outPathValue := flag.CommandLine.Lookup("o").Value.String()

	if outPathValue != "" {
		return outPathValue
	}

	firstArg := flag.Arg(0)
	outDir := filepath.Dir(firstArg)
	name := strings.TrimSuffix(filepath.Base(firstArg), filepath.Ext(firstArg)) + ".sbk"

	return filepath.Join(outDir, name)
}

func gatherEntries() []entry {
	inputFiles := flag.Args()
	var entries []entry

	offset :=	uint32(4 /* magic num & count */ + len(inputFiles) * 8 /* item offset/size */)

	// apply padding to the header if needed
	if offset % 8 != 0 {
		offset += 8 - offset % 8
	}

	for _, inputFile := range  inputFiles {
		stats, err := os.Stat(inputFile)

		if err != nil {
			fmt.Println("Unable to stat file: ", inputFile)
			os.Exit(1)
		}

		// sbc will pad file sizes to be multiples of 8
		size := stats.Size()

		if size % 8 != 0 {
			size += 8 - size % 8
		}

		e := entry{path: inputFile, offset: offset, size: uint32(size)}

		offset += e.size

		entries = append(entries, e)
	}

	return entries
}

func padFile(file *os.File, value byte, count uint32) {
	fill := make([]byte, count)

	for i := uint32(0); i < count; i++ {
		fill[i] = value
	}

	file.Write(fill)
}

func writeHeader(file *os.File, entries []entry) bool {
	var magicNum uint16 = 21297

	// header intro
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, magicNum)
	binary.Write(&buf, binary.BigEndian, uint16(len(entries)))

	// entries table
	for _, entry := range entries {
		binary.Write(&buf, binary.BigEndian, entry.offset)
		binary.Write(&buf, binary.BigEndian, entry.size)
	}

	file.Write(buf.Bytes())

	// pad if necessary
	paddingNeeded := entries[0].offset - uint32(buf.Len())
	if paddingNeeded > 0  {
		padFile(file, 0xFF, paddingNeeded)
	}

	return true
}

func appendFiles(file *os.File, entries []entry) {
	for _, entry := range entries {
		buffer, _ := os.ReadFile(entry.path)
		file.Write(buffer)

		bufferLen := uint32(len(buffer))
		if bufferLen != entry.size {
			padFile(file, 0, entry.size - bufferLen)
		}
	}
}