package main

import (
	"./blastm8"
	"./taxonomy"
	"fmt"
	"os"
	"bufio"
	"io"
	"runtime"
	"flag"
	"sync"
//	"bytes"
)

const VERSION = 0.01

const (
	lcalim = 0.1
	blim   = 1000
	scLim  = 0.9
)

var (
	cpuflag int
	dictflag,nodesflag,namesflag,blastfile,taxlevel string
	verflag,helpflag bool
	printfLock sync.Mutex
)

func printf(format string, args ...interface{}) {
	printfLock.Lock()
	defer printfLock.Unlock()
	fmt.Printf(format, args...)
}

func init() { 
	flag.IntVar(&cpuflag,"ncpus", 4, "Number of cpus for multithreading")
	flag.StringVar(&nodesflag,"nodes","nodes.dmp","nodes.dmp file of taxonomy")
	flag.StringVar(&namesflag,"names","names.dmp","names.dmp file of taxonomy")
	flag.StringVar(&dictflag,"dict","", "dict file of taxonomy")
	flag.StringVar(&taxlevel,"level","family","Desired taxonomical level")
	flag.BoolVar(&verflag,"version",false,"Print VERSION and exits")
	flag.BoolVar(&helpflag,"help",false,"Print USAGE and exits")
	flag.Parse()
	blastfile = flag.Arg(0)
	if verflag {
		fmt.Printf("blast2lca\nVERSION: %.3f\n\n",VERSION)
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
	runtime.GOMAXPROCS(cpuflag)
}

func main() {
	taxDB := taxonomy.New(
		nodesflag,
		namesflag,
		dictflag)

	// BLAST
	blastf, eopen := os.Open(blastfile, os.O_RDONLY, 0644)
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file %s doesn't exist\n",blastfile)
		os.Exit(1)
	}

	blastbuf := bufio.NewReader(io.Reader(blastf))
	queryChan := make(chan *blastm8.Rep2LCA)
	doneChan := make(chan bool)
	launched := 0
	busy := make(chan bool, 1000000)
	go blastm8.Procfile(blastbuf, queryChan, doneChan)
LOOP: for {
		select {
		case k := <-queryChan:
//			fmt.Println(k.Hits)
			launched++
			go func() {
				paths := make([][]*taxonomy.Taxnode, 10000)
				fPos := 0
				rpos := 0
				var bscLim float
				for _, gibs := range k.Hits {
//					fmt.Println(gibs,gi)
					if gibs == nil {
						break
					}
					taxid := taxDB.TaxidFromGi(gibs.GI)
					if taxid != -1 && fPos == 0 {
						bscLim = gibs.BS * scLim
						fPos = 1
				        }
					if taxid == -1 {
						continue
					}
					if gibs.BS >= bscLim {
						path := taxDB.PathFromGi(gibs.GI)
//						taxonomy.OutPath(path
						paths[rpos] = path // TO DO: Check that rpos can not be greater than 10000. This would cause an "index out of bounds" crash
						rpos++
					} else {
						break
					}
				}
//				func() {
				lcaId := taxonomy.LCA(paths)
				fam := taxDB.AtLevel(lcaId, []byte(taxlevel))
				printf("%s\t%s\n",k.Query,fam)
//				fmt.Printf("%s\t%s\n", k.Query, fam)
//				fmt.Println("\n\n")
				busy <- true
			}()
		case <-doneChan:
			for i := 0; i < launched; i++ {
				<-busy
			}
			break LOOP
		}
	}

// 		for k := range queryChan {
// 			go func () {
// 				paths := make([][]*taxonomy.Taxnode,250)
// 				for pos,gi := range k.Hits {
// 					if gi == 0 {
// 						break
// 					}
// 					path := taxDB.PathFromGi(gi)
// 					//			taxonomy.OutPath(path)
// 					paths[pos] = path
// 				}
// 				lcaId := taxonomy.LCA(paths)
// 				//		fmt.Println("LCA: ",lcaId)
// 				fam := taxDB.AtLevel(lcaId,[]byte("family"))
// 				fmt.Printf("%s\t%s\n",k.Query,fam)
// 			}()
// 		}

}
