package main

import (
	"fmt"
	"os"
	"io"
	"bytes"
	"bufio"
	"strings"
	"strconv"
	"gob"
	"time"
	"flag"
)

const VERSION = 0.01
const sep = "\t|\t"

var (
	namesfile, nodesfile, outfile string
	helpflag                      bool
)

type tids struct {
	Taxid int
	Name  string
}

type taxon struct {
	Tids  *tids
	Taxon string
}

type Taxnode struct {
	This     taxon
	ParentId int
}

type taxonomy map[int]*Taxnode

func init() {
	flag.StringVar(&nodesfile, "nodes", "nodes.dmp", "nodes.dmp file of taxonomy")
	flag.StringVar(&namesfile, "names", "names.dmp", "names.dmp file of taxonomy")
<<<<<<< HEAD
	flag.StringVar(&outfile,    "out",   "taxonomy.bin", "output file -- defaults to taxonomy.bin")
=======
	flag.StringVar(&outfile, "out", "taxonomy.bin", "output file -- defaults to taxonomy.bin")
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73

	flag.BoolVar(&helpflag, "help", false, "Print usage and exits")
	flag.Parse()
	if helpflag {
		fmt.Printf("tax2bin\nVERSiON: %.3f\n\n", VERSION)
		flag.Usage()
		os.Exit(0)
	}
}

func readnames(b *bufio.Reader, namesch chan<- *tids) {
	for {
		line, err := b.ReadString('\n')
		if err == os.EOF {
<<<<<<< HEAD
//			bch <- true
=======
			//			bch <- true
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
			close(namesch)
			return
		}
		if pos := strings.Index(line, "scientific name"); pos == -1 {
			continue
		}
		newName := parsename([]byte(line)[0 : len(line)-3]) // ???
		//		fmt.Println(*newName);
		namesch <- newName
	}
<<<<<<< HEAD
//	bch <- true
=======
	//	bch <- true
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	close(namesch)
	return
}

func parsename(l []byte) *tids {
<<<<<<< HEAD
	parts := bytes.Split(l, []byte(sep), -1)
=======
	parts := bytes.Split(l, []byte(sep))
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	tid, _ := strconv.Atoi(fmt.Sprintf("%s", parts[0]))
	newName := &tids{
		tid,
		fmt.Sprintf("%s", parts[1])}
	return newName
}

func parsenode(l []byte) *Taxnode {
<<<<<<< HEAD
	parts := bytes.Split(l, []byte(sep), 4)
=======
	parts := bytes.SplitN(l, []byte(sep), 4)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	tid, _ := strconv.Atoi(fmt.Sprintf("%s", parts[0]))
	pid, _ := strconv.Atoi(fmt.Sprintf("%s", parts[1]))
	//	fmt.Println("KKKK:",tid);
	newRec := Taxnode{
		taxon{
			&tids{
				tid,
				""},
			fmt.Sprintf("%s", parts[2])},
		pid,
	}
	return &newRec
}

func newTaxonomy(b *bufio.Reader) taxonomy {
<<<<<<< HEAD
	taxmap :=  make(map[int]*Taxnode)
=======
	taxmap := make(map[int]*Taxnode)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	for {
		line, err := b.ReadString('\n')
		if err == os.EOF {
			return taxmap
		}
		newNode := parsenode([]byte(line)[0 : len(line)-3]) // ends with \t\|\n
		//		fmt.Println("KK:",newNode.this.taxid);
		taxmap[newNode.This.Tids.Taxid] = newNode
	}
	return taxmap
}

func New(nodes, names string) taxonomy {
	nodesf, eopen := os.OpenFile(nodes, os.O_RDONLY, 0644)
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file doesn't exist %s\n", nodes)
		os.Exit(1)
	}
	nodesbuf := bufio.NewReader(io.Reader(nodesf))

	namesf, eopen := os.OpenFile(names, os.O_RDONLY, 0644)
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file doesn't exist %s\n", names)
		os.Exit(1)
	}
	namesbuf := bufio.NewReader(io.Reader(namesf))

<<<<<<< HEAD
//	endch := make(chan bool)
=======
	//	endch := make(chan bool)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	namesch := make(chan *tids, 1000)
	go readnames(namesbuf, namesch)

	tax := newTaxonomy(nodesbuf)
<<<<<<< HEAD
//	<-endch

	for {
		n,ok := <-namesch
=======
	//	<-endch

	for {
		n, ok := <-namesch
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		if !ok {
			break
		}
		tax[n.Taxid].This.Tids = n
	}

<<<<<<< HEAD
// 	for n := range namesch {
// 		tax[n.Taxid].This.Tids = n
// 	}
=======
	// 	for n := range namesch {
	// 		tax[n.Taxid].This.Tids = n
	// 	}
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73

	return tax
}

// Not used in this program -- We use it only to store the new binary version of taxonomy
<<<<<<< HEAD
func (t taxonomy) Store (fname string) os.Error {
=======
func (t taxonomy) Store(fname string) os.Error {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(t)
	if err != nil {
		return err
	}

	fh, eopen := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
	defer fh.Close()
	if eopen != nil {
		return eopen
	}
<<<<<<< HEAD
	_,e := fh.Write(b.Bytes())
	if e != nil {
		return e
	}
//	fmt.Fprintf(os.Stderr, "%d bytes successfully written to file\n", n)
	return nil
}

func Load (fname string) (taxonomy, os.Error) {
=======
	_, e := fh.Write(b.Bytes())
	if e != nil {
		return e
	}
	//	fmt.Fprintf(os.Stderr, "%d bytes successfully written to file\n", n)
	return nil
}

func Load(fname string) (taxonomy, os.Error) {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	fh, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	t := make(taxonomy)
	dec := gob.NewDecoder(fh)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	return t, nil
<<<<<<< HEAD
} 
=======
}
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73

func main() {
	fmt.Fprintf(os.Stderr, "Building the taxonomy tree ... ")
	s1 := time.Nanoseconds()
	newtax := New(nodesfile, namesfile)
	s2 := time.Nanoseconds()
	if newtax == nil {
<<<<<<< HEAD
		fmt.Fprintf(os.Stderr,"\n ERROR: the map is empty\n")
=======
		fmt.Fprintf(os.Stderr, "\n ERROR: the map is empty\n")
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Done (%.5f)\n", (float32(s2-s1))/1e9)
	fmt.Fprintf(os.Stderr, "Writing tree into binary file ... ")
	s1 = time.Nanoseconds()
	e := newtax.Store(outfile)
	s2 = time.Nanoseconds()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Done (%.5f)\n", (float32(s2-s1))/1e9)

<<<<<<< HEAD
// 	s1 = time.Nanoseconds()
// 	_,err := Load("file.bin")
// 	s2 = time.Nanoseconds()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	fmt.Printf("tree read from file in %.5f seconds\n", float32(s2-s1)/1e9)

}

=======
	// 	s1 = time.Nanoseconds()
	// 	_,err := Load("file.bin")
	// 	s2 = time.Nanoseconds()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Printf("tree read from file in %.5f seconds\n", float32(s2-s1)/1e9)

}
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
