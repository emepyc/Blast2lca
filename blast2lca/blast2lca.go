package main

import (
	"Blast2lca/blastm8"
	"Blast2lca/taxonomy"
	"Blast2lca/kegg"
	"fmt"
	"os"
	"bufio"
	"io"
	"runtime"
	"flag"
	"sync"
	"bytes"
<<<<<<< HEAD
)

const VERSION = 0.1

const (
	lcalim      = 0.1
	blim        = 1000
=======
	"time"
	"runtime/pprof"
	"log"
)

const VERSION = 0.2

const (
	lcalim       = 0.1
	blim         = 1000
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	bscLimFactor = 0.9
)

var (
<<<<<<< HEAD
	cpuflag int
	dictflag, nodesflag, namesflag, blastfile, taxlevel, gi2kegg, bintaxflag, binkegg string
	savememflag, verflag, helpflag, keggflag, taxIsBin, keggIsBin bool // HINT: keggflag is not set by the user, will be true if gi2kegg and kegg2pw are defined / The same with taxIsBin
	printfLock sync.Mutex
=======
	cpuprofile, memprofile string
	cpuflag                                                                           int
	dictflag, nodesflag, namesflag, blastfile, taxlevel, gi2kegg, binkegg, bintaxflag string
	savememflag, verflag, helpflag, keggflag, taxIsBin, keggIsBin                     bool // HINT: keggflag is not set by the user, will be true if gi2kegg and kegg2pw are defined / The same with taxIsBin
	printfLock                                                                        sync.Mutex
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
)

// HINT : Use this function for sensible output! it is threadsafe
func printf(format string, args ...interface{}) {
	printfLock.Lock()
	fmt.Printf(format, args...)
	printfLock.Unlock()
}

func init() {
	flag.IntVar(&cpuflag, "ncpus", 4, "Number of cpus for multithreading [optional]")
	flag.StringVar(&nodesflag, "nodes", "nodes.dmp", "nodes.dmp file of taxonomy")
	flag.StringVar(&namesflag, "names", "names.dmp", "names.dmp file of taxonomy")
	flag.StringVar(&bintaxflag, "bintax", "", "taxonomy in binary format (from tax2bin) [optional]")
	flag.StringVar(&dictflag, "dict", "", "dict file of taxonomy")
	flag.StringVar(&taxlevel, "levels", "family", "Desired taxonomical levels")
	flag.StringVar(&gi2kegg, "gi2kegg", "", "genes_ncbi-gi.list file from kegg [optional]")
<<<<<<< HEAD
//	flag.StringVar(&kegg2pw, "kegg2pw", "", "genes_pathway.list file from kegg [optional]")  // For now, we leave out pathways... can be recovered later
=======
	//	flag.StringVar(&kegg2pw, "kegg2pw", "", "genes_pathway.list file from kegg [optional]")  // For now, we leave out pathways... can be recovered later
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	flag.StringVar(&binkegg, "binkegg", "", "kegg in binary format (from kegg2bin) [optional]")
	flag.BoolVar(&savememflag, "savemem", false, "save memory by keeping files in disk [optional]")
	flag.BoolVar(&verflag, "version", false, "Print VERSION and exits")
	flag.BoolVar(&helpflag, "help", false, "Print USAGE and exits")
<<<<<<< HEAD
	flag.Parse()
//	if (kegg2gi != "" && kegg2pw != "") { // TODO: Treat kegg2gi and kegg2pw independently?
=======
	flag.StringVar(&cpuprofile, "cpuprof", "", "write cpu profile to file")
	flag.StringVar(&memprofile, "memprof", "", "write mem profile to file")
	flag.Parse()
	//	if (kegg2gi != "" && kegg2pw != "") { // TODO: Treat kegg2gi and kegg2pw independently?
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	if gi2kegg != "" {
		keggflag = true
	} else {
		keggflag = false
	}
<<<<<<< HEAD
//	if (kegg2gi != "" && kegg2pw == "") || (kegg2gi == "" && kegg2pw != "") {
//		fmt.Fprintf(os.Stderr, "If KEGG analysis is required, both genes_ncbi-gi.list and genes_pathway.list files are required. Only one provided\n")
//		os.Exit(1)
//	}
=======
	//	if (kegg2gi != "" && kegg2pw == "") || (kegg2gi == "" && kegg2pw != "") {
	//		fmt.Fprintf(os.Stderr, "If KEGG analysis is required, both genes_ncbi-gi.list and genes_pathway.list files are required. Only one provided\n")
	//		os.Exit(1)
	//	}
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	if keggflag && savememflag && binkegg == "" {
		fmt.Fprintf(os.Stderr, "In savemem mode you need to specify the kegg binary file. See the docs for details\n")
		os.Exit(1)
	}
<<<<<<< HEAD
	if (binkegg != "") {
=======
	if binkegg != "" {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		fmt.Fprintf(os.Stderr, "Binary kegg file %s will be used for kegg assignment\n", binkegg)
		keggIsBin = true
	} else {
		keggIsBin = false
	}
<<<<<<< HEAD
	if (bintaxflag != "") {
		fmt.Fprintf(os.Stderr,"Binary taxonomy file %s will be used for taxonomy tree construction\n", bintaxflag)
=======
	if bintaxflag != "" {
		fmt.Fprintf(os.Stderr, "Binary taxonomy file %s will be used for taxonomy tree construction\n", bintaxflag)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		taxIsBin = true
	} else {
		taxIsBin = false
	}
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
	runtime.GOMAXPROCS(cpuflag)
}

func main() {
	opts := &taxonomy.InOpts{
<<<<<<< HEAD
	Nodes : nodesflag,
	Names : namesflag,
	Dict  : dictflag,
	Bintax : bintaxflag,
	TaxIsBin : taxIsBin,
	Savemem : savememflag,
	}
	levs := bytes.Split([]byte(taxlevel), []byte{':'}, -1)
=======
		Nodes:    nodesflag,
		Names:    namesflag,
		Dict:     dictflag,
		Bintax:   bintaxflag,
		TaxIsBin: taxIsBin,
		Savemem:  savememflag,
	}
	levs := bytes.Split([]byte(taxlevel), []byte{':'})
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	taxDB, err := taxonomy.New(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : Impossible to get a valid Taxonomy: %s\n", err)
		os.Exit(1)
	}

	// BLAST
	blastf, eopen := os.OpenFile(blastfile, os.O_RDONLY, 0644)
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file %s doesn't exist\n", blastfile)
		os.Exit(1)
	}
	defer blastf.Close()

	blastbuf := bufio.NewReader(io.Reader(blastf))
	var keggDB kegg.Mapper
<<<<<<< HEAD
	keggflag = false
=======
//	keggflag = false
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	if keggflag == true {
		keggDB, err = kegg.Load(gi2kegg, binkegg, savememflag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "I can't load the Kegg database: %s", err)
			os.Exit(1)
		}
	}
	queryChan := make(chan *blastm8.QueryRes) // TODO: Should we do this buffered?
<<<<<<< HEAD
//	doneChan := make(chan bool)               // TODO: Should we do this buffered?
	launched := 0
	finished := make(chan bool) // TODO: Needed that much?? Profile... any performance impact?
//	go blastm8.Procfile(blastbuf, queryChan, doneChan, keggflag)
	go blastm8.Procfile(blastbuf, queryChan, keggflag)

LOOP:	for {
		if k, ok := <-queryChan; ok {
//			fmt.Fprintf(os.Stderr, "%v", k)
			launched++
			go func() {
				taxids := make([]int, 0, 100) 
=======
	//	doneChan := make(chan bool)                    // TODO: Should we do this buffered?
	launched := 0
	finished := make(chan bool) // TODO: Needed that much?? Profile... any performance impact?
	//	go blastm8.Procfile(blastbuf, queryChan, doneChan, keggflag)
	go blastm8.Procfile(blastbuf, queryChan, keggflag)

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

	t1 := time.Nanoseconds()
LOOP:
	for {
		if k, ok := <-queryChan; ok {
			launched++
			go func() {
				taxids := make([]int, 0, 100)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
				bscLim := k.Best.GetBitsc() * bscLimFactor
				for _, gibs := range k.Hits {
					taxid, err := taxDB.TaxidFromGi(gibs.GetGI())
					if err != nil {
						continue
					}
					if gibs.GetBitsc() >= bscLim {
						taxids = append(taxids, taxid)
					} else {
						break
					}
				}

<<<<<<< HEAD
//				fmt.Fprintf(os.Stderr, "Query: %s => TAXIDS: %v\n", k.Query, taxids)
				lcaNode, err := taxDB.LCA(taxids...)
//				fmt.Fprintf(os.Stderr, "LCANODE:\n", lcaNode)
//				os.Exit(0)
				var atLevs [][]byte
				if err != nil { 
=======
				lcaNode, err := taxDB.LCA(taxids...)
				var atLevs [][]byte
				if err != nil {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
					atLevs = make([][]byte, 1)
					atLevs[0] = []byte("unknown")
				} else {
					atLevs = taxDB.AtLevels(lcaNode, levs...)
				}

				allLevs := bytes.Join(atLevs, []byte(";"))

				keggIDs := ""
				for _, hits := range k.BySubj {
					hitBsLimit := hits.Best.GetBitsc() * bscLimFactor
					for _, hit := range hits.Recs {
						if hit.GetBitsc() < hitBsLimit {
							break
						}
						kegg, err := keggDB.Gi2Kegg(hit.GetGI())
						if err != nil {
							fmt.Fprintf(os.Stderr, "ERROR: Got an error from kegg.Gi2Kegg: %s\n", err)
						}
						if kegg != nil {
<<<<<<< HEAD
							keggIDs += fmt.Sprintf("%s", kegg)
=======
							keggIDs += fmt.Sprintf("%s;", kegg)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
							break
						}
					}
				}
				printf("%s\t%s\t%s\n", k.Query, allLevs, keggIDs)
				finished <- true
			}()
		} else {
			for i := 0; i < launched; i++ {
				<-finished
			}
			break LOOP
		}
	}
<<<<<<< HEAD
=======
	t2 := time.Nanoseconds()
	secs := float32(t2-t1)/1e9;
	fmt.Fprintf(os.Stderr, "%d sequences analyzed in %.3f seconds (%d sequences per second)\n", launched, secs, int32(float32(launched) / secs))
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
}
