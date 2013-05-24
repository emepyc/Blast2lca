// readGi reads the taxid associated with the Gi given as argument using the binary index created by gitaxid2bin
package main

import (
	"github.com/emepyc/Blast2lca/giTaxid"
	"fmt"
	"os"
	"flag"
	"log"
)

const VERSION = 0.01

var (
	dictfile string
	gi       int
	helpflag bool
)

func init () {
	flag.StringVar(&dictfile, "dict", "gi_taxid.bin", "binary version of the gi2taxid file")
	flag.IntVar(&gi, "gi", 3, "GI to look for")
	flag.BoolVar(&helpflag, "help", false, "Print this message and exists")
	flag.Usage = usage
	flag.Parse()
	if helpflag || len(os.Args) == 1 {
		usage()
	}
}

func usage () {
	fmt.Fprintf(os.Stderr, "%s -- Retrieval of Taxids associated with GI using the binary index created by gitaxid2bin\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s -dict=<git_taxid_[nucl|prot].bin> -gi=<GI>\n\n", os.Args[0])
	os.Exit(1);
}

func main () {
	giMapper, err := giTaxid.Load(dictfile, true);
	if err != nil {
		log.Fatalf("Problem reading dict file: %s\n", err)
	}

	taxid, err := giMapper.GiTaxid(gi)
	if err != nil {
		log.Fatalf("Problem retrieving taxid from dict: %s\n", err);
	}
	fmt.Printf("%d\n", taxid)
}