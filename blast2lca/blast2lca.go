//blast2lca calculates the lca from blast results tabular (-m8 format)
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/emepyc/Blast2lca/blastm8"
	"github.com/emepyc/Blast2lca/taxonomy"
)

const VERSION = 0.6

const (
	DEFAULT_BLAST_BUFFER_SIZE = 200 * 1024 * 1024 // Default size for blast reader is 200Mb. TODO -- Try other sizes and profile
)

var (
	cpuprofile, memprofile                              string
	procsflag                                           int
	dictflag, nodesflag, namesflag, blastfile, taxlevel string
	savememflag, verflag, helpflag                      bool
	// order                                               bool
	bscLimFactor float64
	printfLock   sync.Mutex // TODO: Try to avoid this mutex -- Print from a channel -- Make it optionally ordered
	totalQueries int
)

// printf performs threadsafe prints to os.Stdout
func printf(format string, args ...interface{}) {
	printfLock.Lock()
	fmt.Printf(format, args...)
	printfLock.Unlock()
}

func init() {
	flag.IntVar(&procsflag, "nprocs", 4, "Number of cpus for multithreading [optional]")
	flag.StringVar(&nodesflag, "nodes", "nodes.dmp", "nodes.dmp file of taxonomy")
	flag.StringVar(&namesflag, "names", "names.dmp", "names.dmp file of taxonomy")
	flag.StringVar(&dictflag, "dict", "", "Dict file of taxonomy")
	flag.StringVar(&taxlevel, "levels", "", "Desired LCA taxonomical levels [optional]")
	flag.BoolVar(&savememflag, "savemem", false, "Save memory by keeping the gi2taxid mapping file in disk [optional]")
	flag.BoolVar(&verflag, "version", false, "Print VERSION and exits")
	flag.BoolVar(&helpflag, "help", false, "Print USAGE and exits")
	flag.StringVar(&cpuprofile, "cpuprof", "", "Write cpu profile to file")
	flag.StringVar(&memprofile, "memprof", "", "Write mem profile to file")
	flag.Float64Var(&bscLimFactor, "bsfactor", 0.9, "Limit factor for bit score significance")
	// flag.BoolVar(&order, "order", false, "Keep the sequences output in the same order as in the input blast file")
	flag.Parse()

	// The blast file is the first unparsed argument
	blastfile = flag.Arg(0)
	if verflag {
		fmt.Printf("blast2lca\nVERSION: %.3f\n\n", VERSION)
		os.Exit(0)
	}
	if helpflag {
		fmt.Printf("blast2lca\n")
		flag.Usage()
		os.Exit(0)
	}

	if blastfile == "" {
		fmt.Printf("blast2lca\n")
		flag.Usage()
		fmt.Printf("\nA blast file is mandatory\n\n")
		os.Exit(1)
	}
	runtime.GOMAXPROCS(procsflag)
}

func bl2lca(BlastChan <-chan *blastm8.BlastBlock, taxDB *taxonomy.Taxonomy, levs [][]byte, done chan<- struct{}) {
	for {
		select {
		case queryBlock, ok := <-BlastChan:
			if ok {
				totalQueries++
				queryRec := blastm8.ParseRecord(*queryBlock, bscLimFactor)
				taxids := make([]int, 0, len(queryRec.Hits))
				for _, gibs := range queryRec.Hits {
					taxid, err := taxDB.TaxidFromGi(gibs.GI())
					if err != nil {
						log.Printf("WARNING: Taxid can't be retrieved from %s -- Ignoring this record\n", gibs.GI())
						continue
					} else {
						taxids = append(taxids, taxid)
					}
				}
				var atLevs [][]byte
				var allLevs []byte
				lcaNode, err := taxDB.LCA(taxids...)
				if err != nil {
					atLevs = make([][]byte, 1)
					atLevs[0] = taxonomy.Unknown
				} else {
					if len(levs[0]) != 0 {
						atLevs = taxDB.AtLevels(lcaNode, levs...)
					}
				}
				allLevs = bytes.Join(atLevs, []byte{';'})
				msg := fmt.Sprintf("%s\t%s\t%s\t%s\n", queryRec.Query, lcaNode.Name, lcaNode.Taxon, allLevs)
				fmt.Print(msg)
			} else {
				done <- struct{}{}
				return
			}
		default:
		}

	}
}

func main() {
	levs := bytes.Split([]byte(taxlevel), []byte{':'})
	taxDB, err := taxonomy.New(nodesflag, namesflag, dictflag, savememflag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : Impossible to get a valid Taxonomy: %s\n", err)
		os.Exit(1)
	}

	// BLAST
	blastf, eopen := os.OpenFile(blastfile, os.O_RDONLY, 0644) // Use os.Open instead?
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open file %s: %s\n", blastfile, eopen)
		os.Exit(1)
	}
	defer blastf.Close()

	blastbuf := bufio.NewReaderSize(io.Reader(blastf), DEFAULT_BLAST_BUFFER_SIZE)

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}

	blastBlockChan := make(chan *blastm8.BlastBlock, 200)
	done := make(chan struct{})
	t1 := time.Now()

	go blastm8.Procfile(blastbuf, blastBlockChan)
	go bl2lca(blastBlockChan, taxDB, levs, done)
	<-done

	t2 := time.Now()
	dur := t2.Sub(t1)
	secs := dur.Seconds()
	log.Printf("%d sequences analyzed in %.3f seconds (%d sequences per second)\n", totalQueries, secs, int32(float64(totalQueries)/secs))

}
