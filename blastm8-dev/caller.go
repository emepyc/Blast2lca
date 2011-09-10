package main

import (
//	"fmt"
	"os"
	"bufio"
	"./blastm8_dev"
	"log"
	"runtime/pprof"
)

func main () {
	filename := "/home/mp/Documents/courses/PublicHealth_ValAgo2011/seqs/metagenome.blout.full"
	cpuprof := "cpuprof.bin"
	cpufh, err := os.Create(cpuprof)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpufh)
	defer pprof.StopCPUProfile()
	fh, oerr := os.Open(filename)
	if oerr != nil {
		log.Fatal(oerr)
	}
	defer fh.Close()
	buf := bufio.NewReader(fh)
	blastm8_dev.ProcFile(buf)
}
