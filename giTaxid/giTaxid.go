
// Package giTaxid provides functionality to work with GI to Taxid mapping files
package giTaxid
import (
	"fmt"
	"os"
	"time"
	"bufio"
	"bytes"
	"strconv"
	"sync"
	"errors"
	"io"
	"log"
)

const posJump = -100 // We will only read the last 100bp

// GiMapper is the interface that wraps the GiTaxid method
// GiTaxid maps Gi to Taxids
type GiMapper interface {
	GiTaxid ( int ) ( int, error )
}

// OnMemory is the type of the GiTaxid mapper loaded in memory
type OnMemory []uint8

// OnFile is the type of the GiTaxid mapper kept in file
type OnFile struct {
	FileMap *os.File
	FileLen int64
	lock *sync.Mutex
}

// GiTaxid maps Gi to Taxids
func (m OnMemory) GiTaxid (gi int) ( int, error ) {
	pos := gi * 3

	if (len(m)) < pos+3 {
		return -1, errors.New(fmt.Sprintf("GI too high: %d\n", gi))
	}

	taxid := uint32(m[pos+2]) | uint32(m[pos+1]) << 8 | uint32(m[pos]) << 16 | uint32(0) << 24
	return int(taxid), nil
}

// GiTaxid maps Gi to Taxids
func (r *OnFile) GiTaxid ( gi int ) ( int, error ) {
//	fmt.Fprintf(os.Stderr, "----> GI : %d\n", gi)

	r.lock.Lock()
	f := r.FileMap
	pos := int64(gi * 3)

	if pos >= r.FileLen {   // 3 bytes
		return -1, errors.New(fmt.Sprintf("GI too high: %d\n", gi))
	}

	newpos, err := f.Seek(pos, os.SEEK_SET)
	// Seeking beyond EOF doesn't report error
	if newpos != pos {
		return -1, errors.New(fmt.Sprintf("Can't seek to pos %d (maybe too high?)\n", pos))
	}
	if err != nil {
		return -1, err
	}

	bts := make([]byte, 3)
	n, err := f.Read(bts)
	r.lock.Unlock()
	if n != 3 {
		return -1, errors.New(fmt.Sprintf("Can't read for GI %d\n", gi))
	}
	if err != nil {
		return -1, err
	}

	taxid := int(uint32(bts[2]) | uint32(bts[1]) << 8 | uint32(bts[0]) << 16 | uint32(0) << 24)
	return taxid, nil
}

func readLastGI (fname string) (gi int, ok bool) {
	fh, err := os.Open(fname)
	if err != nil {
		log.Printf("WARNING: Unable to open file: %s\n", err)
		return
	}
	defer fh.Close()
	_, e := fh.Seek(posJump, 2)  // We read only the last part of the file
	if e != nil {
		log.Printf("WARNING: Unable to seek file: %s\n", e)
		return
	}
	buf := bufio.NewReader(fh)
	prevLine, isIndex , err := buf.ReadLine()		
	if err != nil {
		log.Printf("WARNING: Unable to read file: %s\n", err)
		return
	}
	if isIndex {
		log.Print("WARNING: Suspicious file (line too long)\n")
		return
	}
	for {
		line, isIndex , err := buf.ReadLine()
		if err == io.EOF { 
			parts := bytes.SplitN(prevLine, []byte("\t"), 2)
			gi, err = strconv.Atoi(string(parts[0]))
			if err != nil {
				log.Printf("WARNING: %s doesn't seem a valid GI number: %s\n", parts[0], err)
				return
			}
			ok = true
			return
		}
		if err != nil {
			log.Printf("WARNING: Unable to read file: %s\n", err)
			return
		}
		if isIndex {
			log.Printf("WARNING: Suspicious long file in file\n")
			return
		}
		prevLine = line
	}
	log.Fatal("Unreachable code")
	return // Never used
}

func readLatestGI(files []string) (gi int) {
	for _, fname := range files {
		log.Printf("Processing dict file: %s ... ", fname)
		t1 := time.Now()
		if lastGi, ok := readLastGI(fname); ok {
			if lastGi > gi {
				gi = lastGi
			}
			t2 := time.Now()
			dur := t2.Sub(t1)
			log.Printf("Done (%.3f secs)\n", dur.Seconds())
		} else {
			log.Fatal("ERROR: File %s can't be processed\n", fname)
		}
	}
	return gi
}

// encode encodes the gi taxid mappings
func (m OnMemory) encode (gi, taxid int) {
	pos := gi * 3
	m[pos+2] = byte(taxid >> 0)
	m[pos+1] = byte(taxid >> 8)
	m[pos] = byte(taxid >> 16)
}


// loadTextMapper incorporates the fh file mapper to m
func (m OnMemory) loadTextMapper (fh *os.File) error {
	buf := bufio.NewReader(fh)
	for {
		line, _ , err := buf.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		parts := bytes.Split(line, []byte("\t"))
		gi, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			return err
		}
		taxid, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			return err
		}
		m.encode(gi, taxid)
	}
	log.Fatal("Unreachable code")
	return nil // never reached
}

// New creates a new binary dict file from the input text dict file/s
// Returns a the binary structure or any error it may encounter in the process
func New (files []string) (OnMemory, error) {
	log.Print("Creating new Gi2taxid binary structure ")
	t1 := time.Now()
	var m OnMemory
	lastGi := readLatestGI(files)
	m = make(OnMemory, (lastGi*3)+3)

	for _, file := range files {
		log.Printf("(%s) ... ", file)
		fh, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer fh.Close()

		m.loadTextMapper(fh)

		t2 := time.Now()
		dur := t2.Sub(t1)
		log.Printf("Done (%.3f secs)\n", dur.Seconds())
		return m, nil
	}
	return m, nil // Never used
} 

// Store stores the binary representation of the Gi => Taxid mapper into the fname file
// Returns nil or any error it may encounter in the process
func (m OnMemory) Store (fname string) error {
	fmt.Fprintf(os.Stderr, "Storing binary structure to file ... ")
	t1 := time.Now()
	fh, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = fh.Write(m)
	t2 := time.Now()
	dur := t2.Sub(t1)
	log.Printf("Done (%.3f secs)\n", dur.Seconds())
	return err
}

// Load loads the binary representation of the Gi => Taxid mapper from the fname file
// Returns a GiMapper or any error it may encounter in the process
// If savemem is true, the mapper is not loaded into memory and will be kept in the file
func Load (fname string, savemem bool) ( GiMapper, error ) {
	fmt.Fprintf(os.Stderr, "Loading Gi2taxid binary file ... ")
	t1 := time.Now()
	if savemem {
		fh, err := os.Open(fname)
		if err != nil {
			return nil, err
		}
		d, err := os.Stat(fname)
		if err != nil {
			return nil, err
		}
		t2 := time.Now()
		dur := t2.Sub(t1)
		fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", dur.Seconds())
		return &OnFile {
		FileMap : fh,
		FileLen : d.Size(),
		lock : new(sync.Mutex),
		}, nil
	}
	fh, err := os.OpenFile(fname, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	d, err := fh.Stat()
	if err != nil {
		return nil, err
	}
	b := make([]byte, d.Size())
	_, err = io.ReadFull(fh, b)
	if err != nil {
		return nil, err
	}
	t2 := time.Now()
	dur := t2.Sub(t1)
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", dur.Seconds())
	return OnMemory(b), nil
}

