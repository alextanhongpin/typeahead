package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	ahead "github.com/alextanhongpin/go-ahead"
)

func main() {
	var (
		cpuprofile  = flag.String("cpu", "", "The file to output the cpu profiling, e.g. cpu.out")
		memprofile  = flag.String("mem", "", "The file to output the memory profiling, e.g. mem.out")
		interactive = flag.Bool("i", false, "Whether to allow interactive mode or not")
		source      = flag.String("source", "", "the default dictionary to load")
		in          = flag.String("in", "", "the file that stores the struct")
		out         = flag.String("out", "", "the destination to store the file to")
	)
	flag.Parse()
	if *cpuprofile != "" {
		cpufile, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpufile)
		defer pprof.StopCPUProfile()
	}

	root := ahead.New()

	if *in != "" {
		var f *os.File
		_, err := os.Stat(*in)
		if os.IsNotExist(err) {
			f, err = os.Create(*in)
		} else {
			f, err = os.Open(*in)
		}
		if err != nil {
			log.Fatal(err)
		}
		dec := gob.NewDecoder(f)
		err = dec.Decode(root)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("read from", *in)
	}

	if *source != "" {
		f, err := os.Open(*source)
		// f, err := os.Open("/usr/share/dict/words")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		var count, words int
		for scanner.Scan() {
			b := bytes.ToLower(scanner.Bytes())
			words++
			count += len(b)
			root.Insert(b, nil)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("inserted", words, "words", count, "characters")
	}

	if *out != "" {
		var f *os.File
		_, err := os.Stat(*out)
		if os.IsNotExist(err) {
			f, err = os.Create(*out)
		} else {
			f, err = os.Open(*out)
		}
		if err != nil {
			log.Fatal(err)
		}
		enc := gob.NewEncoder(f)
		err = enc.Encode(root)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("store to", *out)
	}

	if *memprofile != "" {
		memfile, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(memfile)
		defer memfile.Close()
	}

	if *interactive {
		fmt.Println("Enter a search keyword:")
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			b := bytes.TrimSpace(reader.Bytes())
			if len(b) == 0 {
				continue
			}
			fmt.Printf("searching for %s:\n", b)
			start := time.Now()
			// result := root.Find(b)
			result := root.FindRecursive(b)
			var count int
			fmt.Printf("found %d results in %s\n", len(result), time.Since(start))
			// for r, _ := range result {
			//         fmt.Println(r)
			//         count++
			// }
			for _, b := range result {
				fmt.Println(string(b))
				count++
			}
			fmt.Printf("found %d results in %s\n", count, time.Since(start))
			fmt.Println()
		}
		if err := reader.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
