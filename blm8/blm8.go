package main

import (
	"fmt"
	"Blast2lca/blastm8"
	"os"
	"bufio"
	"time"
	"flag"
	"runtime/pprof"
	"log"
)

var cpuprof = flag.String("cpuprof", "", "write cpu profile to file")
var fname = flag.String("fname", "", "input blast file")

func main () {
//	fname := "/home/mp/Documents/courses/PublicHealth_ValAgo2011/seqs/metagenome.blout.full"
//	fname := "sample/metagenome.blout.127"
	flag.Parse()
	if *cpuprof != "" {
		f, err := os.Create(*cpuprof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	fh, err := os.Open(*fname)
	defer fh.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	buf := bufio.NewReader(fh)
	queryChan := make(chan *blastm8.QueryRes, 1000)
	go blastm8.Procfile(buf, queryChan, false)
	n := 0
	t1 := time.Nanoseconds()
	for k := range queryChan {
		n++
		k=k
		fmt.Printf("1\n")
	}
	t2 := time.Nanoseconds()
	secs := float32(t2-t1)/1e9
	fmt.Fprintf(os.Stderr, "%d sequences analyzed in %.3f seconds (%d sequences per second)\n", n, secs, int32(float32(n)/secs))
}

// func packN (chan *blastm8.QueryRes, n int) chan *[]blastm8.QueryRes {
// 	outchan := make(chan *[]blastm8.QueryRes)
	
// }
