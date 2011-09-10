package main

import (
	. "Blast2lca/giTaxid"
	"fmt"
	"os"
	"time"
	"flag"
)

const VERSION = 0.01

var (
	giTaxid, outfile string
	helpflag bool
)

func init () {
	flag.StringVar(&giTaxid, "gitaxid", "gi_taxid_prot.dmp", "gi_taxid_prot.dmp [default] or gi_taxid_nucl.dmp file from NCBI's taxonomy repository")
	flag.StringVar(&outfile, "outbin", "gi_taxid_prot.bin", "binary converted version [defaults to gi_taxid_prot.bin]")
	flag.BoolVar(&helpflag, "help", false, "Print USAGE and exits")
	flag.Parse()
	if helpflag {
		fmt.Printf("gitaxid2bin\n")
		flag.Usage()
		os.Exit(0)
	}
}

func main () {
	fmt.Fprintf(os.Stderr, "Encoding dictionary ... ")
	t1 := time.Nanoseconds()
	m, err := New(giTaxid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", err)
		os.Exit(1)
	}
	t2 := time.Nanoseconds()
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2 - t1)/1e9)
	fmt.Fprintf(os.Stderr, "Storing dictionary ... ")
	t1 = time.Nanoseconds()
	err = m.Store(outfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", err)
		os.Exit(1)
	}
	t2 = time.Nanoseconds()
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2 - t1)/1e9)
}
