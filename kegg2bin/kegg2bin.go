package main

import (
	. "Blast2lca/kegg"
	"fmt"
	"os"
	"flag"
	"time"
)

const VERSION = 0.01

var (
	gi2kegg, outfile string
	helpflag bool
)

func init () {
	flag.StringVar(&gi2kegg, "gi2kegg", "genes_ncbi-gi.list", "genes_ncbi-gi.list file")
//	flag.StringVar(&kegg2pw, "kegg2pw", "genes_pathway.list", "genes_pathway.list")
	flag.StringVar(&outfile, "outfile", "kegg.bin", "Output file")
	flag.BoolVar(&helpflag, "help", false, "Print USAGE and exits")
	flag.Parse()
	if helpflag {
		fmt.Printf("kegg2bin\n")
		flag.Usage()
		os.Exit(0)
	}
}

func main () {
	fmt.Fprintf(os.Stderr, "Generating KEGG binary format on file %s ... ", outfile)
	t1 := time.Nanoseconds()
	err := Store( gi2kegg, outfile )
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: Can't Build KEGG Dictionary: %s\n", err)
		os.Exit(1)
	}
	t2 := time.Nanoseconds()
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2-t1)/1e9)

// 	fmt.Fprintf(os.Stderr, "Loading KEGG binary database ... ")
// 	t1 := time.Nanoseconds()
// 	keggDB, err := kegg.Load("kegg.bin", .... more things and "false")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error loading DB")
// 		os.Exit(1)
// 	}
// 	t2 := time.Nanoseconds()
// 	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2 - t1)/1e9)
// 	if keggDB != nil {}
}

