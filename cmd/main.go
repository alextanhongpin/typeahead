package main

import (
	"bufio"
	"bytes"
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
	f, err := os.Open("/usr/share/dict/words")
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
