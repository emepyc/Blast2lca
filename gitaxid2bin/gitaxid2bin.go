
// gitaxid2bin creates a new binary dictionary from the input files specified in the command line
package main

import (
	"Blast2lca/giTaxid"
	"fmt"
	"os"
	"time"
	"flag"
	"log"
)

const VERSION = 0.01

var (
	outfile string
	helpflag bool
)

func init () {
	flag.StringVar(&outfile, "outbin", "gi_taxid.bin", "binary converted version [defaults to gi_taxid.bin]")
	flag.BoolVar(&helpflag, "help", false, "Print this message and exits")
	flag.Usage = usage
	flag.Parse()
	if helpflag || len(os.Args) == 1 {
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "\n%s converts a list of gi => taxid mapping files to binary format\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "The resulting binary mapper is written on the file specified by the -outfile parameter\n")
	fmt.Fprintf(os.Stderr, "Usage: %s gi_taxid_[nucl|prot].dmp...\n\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(2)
}

func main () {
	fmt.Fprintf(os.Stderr, "Encoding dictionary ... ")
	t1 := time.Now()

	m, err := giTaxid.New(flag.Args())
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	t2 := time.Now()
	dur := t2.Sub(t1)
	log.Printf("Done (%.3f secs)\n", dur.Seconds())
	fmt.Fprintf(os.Stderr, "Storing dictionary ... ")
	t1 = time.Now()
	err = m.Store(outfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", err)
		os.Exit(1)
	}
	t2 = time.Now()
	dur = t2.Sub(t1)
	log.Printf("Done (%.3f secs)\n", dur.Seconds())
}
