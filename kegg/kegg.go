package kegg

/* This package reads the  */
import (
	"fmt"
	"os"
	"bufio"
	"io"
	"bytes"
	"strconv"
	"time"
<<<<<<< HEAD
//	"gob"
	"sync"
)

const (
	maxKeggs    = 7000000 //7e6 
	maxPathways = 4000000 //4e6  We don't use pathways for now
	recs4report = 1000   //1e3
)

type Mapper interface {
	Gi2Kegg ( int ) ( []byte, os.Error )
=======
	//	"gob"
	"sync"
	"Blast2lca/wcl"
)

const (
//	maxKeggs    = 7000000 //7e6 WARNING!! Hardcoded!
	maxPathways = 4000000 //4e6  We don't use pathways for now
	recs4report = 1000    //1e3
)

type Mapper interface {
	Gi2Kegg(int) ([]byte, os.Error)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
}

type data2store struct {
	gi, lpos int
<<<<<<< HEAD
	off byte
=======
	off      byte
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
}

type Info struct {
	Gene    []byte // TODO: Unexport fields & make accessors ?
	Pathway [][]byte
}

type PathwMap map[string][][]byte
type OnMemory map[int][]byte

type MutFile struct {
<<<<<<< HEAD
	Fh *os.File
=======
	Fh   *os.File
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	lock *sync.Mutex
}

type OnFile struct {
	Bin, Gi2gene *MutFile
}

<<<<<<< HEAD

=======
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
// for now, pathway mapping is not used
func MapPathways(fname []byte) PathwMap {
	mapPw := make(map[string][][]byte, maxPathways)
	dictf, err := os.Open(fmt.Sprintf("%s", fname))
	if err != nil {
		fmt.Fprintf(os.Stderr, "file %s doesn't exist\n", fname)
		os.Exit(1)
	}
	defer dictf.Close()

	dictBuff := bufio.NewReader(io.Reader(dictf))
	for {
		line, err := dictBuff.ReadString('\n')
		if err == os.EOF {
			return mapPw
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "file %s can't be read\n", fname)
			os.Exit(1)
		}
<<<<<<< HEAD
		parts := bytes.Split([]byte(line), []byte("\t"), -1)
		gene := string(parts[0])
		pathw_bytes := (bytes.Split(parts[1], []byte(":"), -1))[1] // HINT: example: path:hsa00232
=======
		parts := bytes.Split([]byte(line), []byte("\t"))
		gene := string(parts[0])
		pathw_bytes := (bytes.Split(parts[1], []byte(":")))[1] // HINT: example: path:hsa00232
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		mapPw[gene] = append(mapPw[gene], pathw_bytes)
		//mapPw[gene] = pathw_bytes
	}
	return mapPw
}
<<<<<<< HEAD
func BuildDB (gene2giFn []byte) (OnMemory, os.Error) { // HINT: OnMemory is a map, it is passed by reference
=======
func BuildDB(gene2giFn []byte) (OnMemory, os.Error) { // HINT: OnMemory is a map, it is passed by reference
	maxKeggs, wclerr := wcl.FromFile(string(gene2giFn))
	if wclerr != nil {
		return nil, wclerr
	}
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	keggDict := make(OnMemory, maxKeggs) // <7,000,000 of records currently in genes_ncbi-gi.list.
	dictf, err := os.Open(fmt.Sprintf("%s", gene2giFn))
	if err != nil {
		return nil, err
	}
	defer dictf.Close()

	dictBuff := bufio.NewReader(io.Reader(dictf))
	for {
		line, err := dictBuff.ReadString('\n')
		if err == os.EOF {
<<<<<<< HEAD
			return keggDict, nil     // Normal return
=======
			return keggDict, nil // Normal return
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		}
		if err != nil {
			return nil, err
		}
<<<<<<< HEAD
		parts := bytes.Split([]byte(line), []byte("\t"), -1) // HINT: -1 means "split as much as you can"
		gene := parts[0]
		gi_bytes := (bytes.Split(parts[1], []byte(":"), -1))[1] // HINT: ncbi-gi:21071030
=======
		parts := bytes.Split([]byte(line), []byte("\t"))
		gene := parts[0]
		gi_bytes := (bytes.Split(parts[1], []byte(":")))[1] // HINT: ncbi-gi:21071030
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		gi, serr := strconv.Atoi(string(gi_bytes))
		if serr != nil {
			return nil, serr
		}
		keggDict[gi] = gene
	}
<<<<<<< HEAD
	return keggDict, nil   // never in use!
=======
	return keggDict, nil // never in use!
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
}

// Not used for now... we have disabled pathway lookups
// func BuildDB2(gene2giFn, gene2PwFn []byte) (OnMemory, os.Error) { // HINT: DB is a map, it is passed by reference
// 	fmt.Fprintf(os.Stderr, "Building Kegg DB ... ")
// 	s1 := time.Nanoseconds()
// 	pathwMap := MapPathways(gene2PwFn)
// 	keggDict := make(OnMemory, maxKeggs) // <7,000,000 of records currently in genes_ncbi-gi.list
// 	dictf, err := os.Open(fmt.Sprintf("%s", gene2giFn))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer dictf.Close()

// 	dictBuff := bufio.NewReader(io.Reader(dictf))
// 	for {
// 		line, err := dictBuff.ReadString('\n')
// 		if err == os.EOF {
// 			s2 := time.Nanoseconds()
// 			fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(s2 - s1)/1e9)
// 			return keggDict, nil     // Normal return
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		parts := bytes.Split([]byte(line), []byte("\t"), -1) // HINT: -1 means "split as much as you can"
// 		gene := parts[0]
// 		gi_bytes := (bytes.Split(parts[1], []byte(":"), -1))[1] // HINT: ncbi-gi:21071030
// 		gi, serr := strconv.Atoi(fmt.Sprintf("%s", gi_bytes))
// 		if serr != nil {
// 			return nil, serr
// 		}
// 		pathw, ok := pathwMap[(fmt.Sprintf("%s", gene))]
// 		if ok {
// 			keggDict[gi] = &Info{Gene: gene, Pathway: pathw}
// 		} else {
// 			keggDict[gi] = &Info{Gene: gene}
// 		}
// 	}
// 	return keggDict, nil   // never in use!
// }

<<<<<<< HEAD
func (m *OnMemory) Gi2Kegg ( gi int ) ( []byte, os.Error ) {
=======
func (m *OnMemory) Gi2Kegg(gi int) ([]byte, os.Error) {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	gene, ok := (*m)[gi]
	if ok {
		return gene, nil
	}
	return nil, nil
}

<<<<<<< HEAD
func (m *OnFile) Gi2Kegg ( gi int ) ( []byte, os.Error ) {
=======
func (m *OnFile) Gi2Kegg(gi int) ([]byte, os.Error) {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	pos := int64(gi * 5)
	m.Bin.lock.Lock()
	_, err := m.Bin.Fh.Seek(pos, 0)
	if err != nil {
		return nil, err
	}
	data := make([]byte, 5)
	n, err := m.Bin.Fh.Read(data)
	m.Bin.lock.Unlock()
	if n != 5 {
<<<<<<< HEAD
=======
		fmt.Fprintf(os.Stderr,"READ only %d byes: %s\n", n, data)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		return nil, os.NewError("Too few bytes read")
	}
	if err != nil {
		return nil, err
	}
	if data[0] == 0 && data[1] == 0 && data[2] == 0 && data[3] == 0 && data[4] == 0 {
		// It is not present in the database
		return nil, nil
	}
<<<<<<< HEAD
	pos2 := int64(uint32(data[3]) | uint32(data[2]) << 8 | uint32(data[1]) << 16 | uint32(data[0]) << 24)
=======
	pos2 := int64(uint32(data[3]) | uint32(data[2])<<8 | uint32(data[1])<<16 | uint32(data[0])<<24)
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	offset := int(data[4])
	m.Gi2gene.lock.Lock()
	_, err = m.Gi2gene.Fh.Seek(pos2, 0)
	geneBytes := make([]byte, offset)
	n, err = m.Gi2gene.Fh.Read(geneBytes)
	m.Gi2gene.lock.Unlock()
	if n != offset {
		return nil, os.NewError("Too few bytes read")
	}
	return geneBytes, nil
}

// For now, not used (pathways disallowed)
// func (m *OnMemory) Get(gi int) *Info {
// 	info, ok := (*m)[gi]
// 	if ok {
// 		return info
// 	}
// 	return nil
// }

<<<<<<< HEAD
func parseLine (line []byte) (*data2store, os.Error) { // gi, pos, offset, err
	parts := bytes.Split(line, []byte("\t"), -1)
=======
func parseLine(line []byte) (*data2store, os.Error) { // gi, pos, offset, err
	parts := bytes.Split(line, []byte("\t"))
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	gene := parts[0]
	preN := bytes.IndexByte(gene, ':')
	if preN == -1 {
		return &data2store{}, os.NewError("Kegg Gene doesn't have ':' in name")
	}
<<<<<<< HEAD
	gi_bytes := (bytes.Split(parts[1], []byte(":"), -1))[1]
=======
	gi_bytes := (bytes.Split(parts[1], []byte(":")))[1]
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73

	gi, err := strconv.Atoi(string(gi_bytes))
	if err != nil {
		return &data2store{}, err
	}
	return &data2store{
<<<<<<< HEAD
	gi   : gi,
	lpos : preN+1,
	off  : byte(len(gene)-preN-1),
	}, nil
}

func readFull (keggfn string) ([]byte, os.Error) {
=======
		gi:   gi,
		lpos: preN + 1,
		off:  byte(len(gene) - preN - 1),
	}, nil
}

func readFull(keggfn string) ([]byte, os.Error) {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	fmt.Fprintf(os.Stderr, "Reading kegg2gi file ... ")
	t1 := time.Nanoseconds()
	fh, err := os.Open(keggfn)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	d, err := fh.Stat()
	if err != nil {
		return nil, err
	}
	fullkegg := make([]byte, d.Size)
	_, err = fh.Read(fullkegg)
	if err != nil {
		return nil, err
	}
	t2 := time.Nanoseconds()
<<<<<<< HEAD
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2 - t1)/1e9)
	return fullkegg, nil
}

func Store (keggfn, binout string) os.Error {
=======
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2-t1)/1e9)
	return fullkegg, nil
}

func Store(keggfn, binout string) os.Error {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	fout, err := os.OpenFile(binout, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
<<<<<<< HEAD
=======
	maxKeggs, wclerr := wcl.FromFile(keggfn)
	if wclerr != nil {
		return wclerr
	}
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	fullkegg, err := readFull(keggfn)
	if err != nil {
		return err
	}
	fullbuff := bytes.NewBuffer(fullkegg)
	nrecs := 0
	posLine := 0
	for {
<<<<<<< HEAD
		newline, err := fullbuff.ReadBytes('\n')    // WARNING... Non portable
		if err == os.EOF {
			return nil   // We are done
=======
		newline, err := fullbuff.ReadBytes('\n') // WARNING... Non portable
		if err == os.EOF {
			return nil // We are done
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		}
		if err != nil {
			return err
		}
		nrecs++
		if (nrecs % recs4report) == 0 {
<<<<<<< HEAD
			fmt.Printf("\r%d lines read (%d %%)   ", nrecs, nrecs * 100 / maxKeggs)
		}
	//	ch := make(chan *data2store, 100)
//		done := make(chan bool)
		d, err := parseLine(newline)  // newline has '\n' still at the end... no problem with that, right?
=======
			fmt.Printf("\r%d lines read (%d %%)   ", nrecs, (nrecs*100/maxKeggs)+1)
		}
		//	ch := make(chan *data2store, 100)
		//		done := make(chan bool)
		d, err := parseLine(newline) // newline has '\n' still at the end... no problem with that, right?
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		if err != nil {
			return err
		}
		_, err = fout.Seek(int64(d.gi*5), 0)
		if err != nil {
			return err
		}
		fpos := posLine + d.lpos
		data := make([]byte, 4)
		data[0] = byte(fpos >> 24)
		data[1] = byte(fpos >> 16)
		data[2] = byte(fpos >> 8)
		data[3] = byte(fpos >> 0)
		_, err = fout.Write(data)
		if err != nil {
			return err
		}
		_, err = fout.Write([]byte{d.off})
		if err != nil {
			return err
		}
		posLine += len(newline)
	}
	return nil // never used
}
// func Store (keggfn, binout string) os.Error {
// 	fh, err := os.Open(keggfn)
// 	if err != nil {
// 		return err
// 	}
// 	defer fh.Close()
// 	fout, err := os.OpenFile(binout, os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer fout.Close()
// 	var line []byte
// 	pos := 0
// 	posLine := 0
// 	nrecs := 0
// 	fmt.Printf("\n")
// 	for {
// 		var nextByte []byte = []byte{' '}
// 		_, e := fh.Read(nextByte)
// 		if e == os.EOF {
// 			return nil
// 		}
// 		if nextByte[0] == '\n' {
// 			nrecs++
// 			if (nrecs % recs4report) == 0 {
// 				fmt.Printf("\rLines: %d (%d %%)  ", nrecs, nrecs * 100 / maxKeggs)
// 			}
// 			gi, lpos, off, err := parseLine(line)
// 			if err != nil {
// 				return err
// 			}
// 			_, err = fout.Seek(int64(gi*5), 0)
// 			if err != nil {
// 				return err
// 			}
// 			fpos := posLine+lpos
// 			data := make([]byte, 4)
// 			data[0] = byte(fpos >> 24)
// 			data[1] = byte(fpos >> 16)
// 			data[2] = byte(fpos >> 8)
// 			data[3] = byte(fpos >> 0)  // Is this doing something?
// 			_, err = fout.Write(data)
// 			if err != nil {
// 				return err
// 			}
// 			_, err = fout.Write([]byte{off})
// 			if err != nil {
// 				return err
// 			}
// 			line = []byte{}
// 			posLine = pos
// 		}
// 		line = append(line, nextByte[0])
// 		pos++
// 	}
// 	return nil // never used
// }

<<<<<<< HEAD
func Load ( gi2gene, binfile string, savemem bool ) ( Mapper, os.Error ) {
=======
func Load(gi2gene, binfile string, savemem bool) (Mapper, os.Error) {
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
	fmt.Fprintf(os.Stderr, "Loading Kegg Database ... ")
	t1 := time.Nanoseconds()
	if savemem {
		keggFh, err := os.Open(gi2gene)
		if err != nil {
			return nil, err
		}
<<<<<<< HEAD
		
=======

>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		binFh, err := os.Open(binfile)
		if err != nil {
			return nil, err
		}
		t2 := time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2-t1)/1e9)
<<<<<<< HEAD
		return &OnFile {
		Bin : &MutFile{Fh : binFh, lock : new(sync.Mutex)},
		Gi2gene : &MutFile{Fh : keggFh, lock : new(sync.Mutex)},
=======
		return &OnFile{
			Bin:     &MutFile{Fh: binFh, lock: new(sync.Mutex)},
			Gi2gene: &MutFile{Fh: keggFh, lock: new(sync.Mutex)},
>>>>>>> e89fd5ad3651902c534e4c995e3d1f4805443d73
		}, nil
	}
	newdb, err := BuildDB([]byte(gi2gene))
	if err != nil {
		return nil, err
	}
	return &newdb, nil
}

// This Store is for "gob"bed versions of the DB (old versions -- with pathways)
// func (m *DB) Store (fname string) os.Error{
// 	fmt.Fprintf(os.Stderr, "Storing the KEGG DB ... ")
// 	t1 := time.Nanoseconds()
// 	b := new(bytes.Buffer)
// 	enc := gob.NewEncoder(b)
// 	err := enc.Encode(m)
// 	if err != nil {
// 		return err
// 	}
// 	fh, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644)
// 	defer fh.Close()
// 	if err != nil {
// 		return err
// 	}

// 	_, err = fh.Write(b.Bytes())
// 	if err != nil {
// 		return err
// 	}
// 	t2 := time.Nanoseconds()
// 	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2-t1)/1e9)
// 	return nil
// }

// This Load is for "gob"bed versions of the DB (old versions -- with pathways)
// func Load (fname string) (DB, os.Error) {
// 	fmt.Fprintf("Loading KEGG DB ... ")
// 	t1 := time.Nanoseconds()
// 	fh, err := os.Open(fname)
// 	if err != nil {
// 		return nil, err
// 	}
// 	m := make(DB, maxKeggs)
// 	dec := gob.NewDecoder(fh)
// 	err = dec.Decode(&m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	t2 := time.Nanoseconds()
// 	fmt.Fprintf("Done (%.3f secs)\n", float32(t2 - t1)/1e9)
// 	return m, nil
// }

func (m *Info) String() string {
	str := fmt.Sprintf("Kegg_Gene: %s\n", m.Gene)
	for _, pw := range m.Pathway {
		str += fmt.Sprintf("\tPathway:%s\n", pw)
	}
	return str
}

// func main () {
// 	kegg2gi_fn := "/Users/pignatelli/Desktop/genes_ncbi-gi.list"
// 	kegg2pw_fn := "/Users/pignatelli/Desktop/genes_pathway.list"
// 	myMap := BuildDB([]byte(kegg2gi_fn),[]byte(kegg2pw_fn))
// 	k := myMap.Get(14589889)
// 	fmt.Println(k)
// }
