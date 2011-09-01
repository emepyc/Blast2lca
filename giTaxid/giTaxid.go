package giTaxid
import (
	"fmt"
	"os"
	"time"
	"bufio"
	"bytes"
	"strconv"
	"sync"
//	"io/ioutil"
	"io"
//	"runtime/pprof"
)

const posJump = -100 // We will only read the last 100bp

type GiMapper interface {
	GiTaxid ( int ) ( int, os.Error )
}

type OnMemory []uint8

type OnFile struct {
	FileMap *os.File
	FileLen int64
	lock *sync.Mutex
}

func (m OnMemory) GiTaxid (gi int) ( int, os.Error ) {
	pos := gi * 3

	if (len(m)) < pos+3 {
		return -1, os.NewError(fmt.Sprintf("GI too high: %d\n", gi))
	}

	taxid := uint32(m[pos+2]) | uint32(m[pos+1]) << 8 | uint32(m[pos]) << 16 | uint32(0) << 24
	return int(taxid), nil
}

func (r *OnFile) GiTaxid ( gi int ) ( int, os.Error ) {
//	fmt.Fprintf(os.Stderr, "----> GI : %d\n", gi)

	r.lock.Lock()
	f := r.FileMap
	pos := int64(gi * 3)

	if pos > r.FileLen - 24 {   // 3 bytes
		return -1, os.NewError(fmt.Sprintf("GI too high: %d\n", gi))
	}

	newpos, err := f.Seek(pos, os.SEEK_SET)
	// Seeking beyond EOF doesn't report error
	if newpos != pos {
		return -1, os.NewError(fmt.Sprintf("Can't seek to pos %d (maybe too high?)\n", pos))
	}
	if err != nil {
		return -1, err
	}

	bts := make([]byte, 3)
	n, err := f.Read(bts)
	r.lock.Unlock()
	if n != 3 {
		return -1, os.NewError(fmt.Sprintf("Can't read for GI %d\n", gi))
	}
	if err != nil {
		return -1, err
	}

	taxid := int(uint32(bts[2]) | uint32(bts[1]) << 8 | uint32(bts[0]) << 16 | uint32(0) << 24)
	return taxid, nil
}

func (m OnMemory) Encode (gi, taxid int) {
	pos := gi * 3
	m[pos+2] = byte(taxid >> 0)
	m[pos+1] = byte(taxid >> 8)
	m[pos] = byte(taxid >> 16)
}

func getLastGi (fh *os.File) ([]byte, os.Error) {
	_, e := fh.Seek(posJump, 2)  // We read only the last part of the file
	if e != nil {
		return nil, e
	}
	buf := bufio.NewReader(fh)
	prevLine, _ , err := buf.ReadLine()
	if err != nil {
		return nil, err
	}
	for {
		line, _ , err := buf.ReadLine()
		if err == os.EOF { 
			parts := bytes.SplitN(prevLine, []byte("\t"), 2)
			return parts[0], nil
		}
		if err != nil {
			return nil, err
		}
		prevLine = line
	}
	return nil, nil // Never used
}

func New (fileIn string) (OnMemory, os.Error) {
	fmt.Fprintf(os.Stderr, "Creating new Gi2taxid binary structure ... ")
	t1 := time.Nanoseconds()
	fh, err := os.Open(fileIn)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	lastGiStr,err := getLastGi(fh)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Atoi error: %s\n", err)
		return nil, err
	}

	lastGi, err := strconv.Atoi(string(lastGiStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Atoi error: %s\n", err)
	}

	m := make(OnMemory, (lastGi*3)+3)
	if err != nil {
		return nil, err
	}

	_, err = fh.Seek(0,0)
	buf := bufio.NewReader(fh)

	for {
		line, _ , err := buf.ReadLine()
		if err == os.EOF {
			t2 := time.Nanoseconds()
			fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2 - t1)/1e9)
			return m, nil
		}
		if err != nil {
			return nil, err
		}
		parts := bytes.Split(line, []byte("\t"))
		gi, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			return nil, err
		}
		taxid, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			return nil, err
		}
		m.Encode(gi, taxid)
	}
	return m, nil // Never used
} 

func (m *OnMemory) Store (fname string) os.Error {
	fmt.Fprintf(os.Stderr, "Storing binary structure to file ... ")
	t1 := time.Nanoseconds()
	fh, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = fh.Write(*m)
	t2 := time.Nanoseconds()
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2-t1)/1e9)
	return err
}

func Load (fname string, savemem bool) ( GiMapper, os.Error ) {
	fmt.Fprintf(os.Stderr, "Loading Gi2taxid binary file ... ")
	t1 := time.Nanoseconds()
	if savemem {
		fh, err := os.Open(fname)
		if err != nil {
			return nil, err
		}
		d, err := os.Stat(fname)
		if err != nil {
			return nil, err
		}
		t2 := time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2-t1)/1e9)
		return &OnFile {
		FileMap : fh,
		FileLen : d.Size,
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
	b := make([]byte, d.Size)
	_, err = io.ReadFull(fh, b)
//	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	t2 := time.Nanoseconds()
	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2-t1)/1e9)
	return OnMemory(b), nil
}

// func main() {

// 	eh := http.ListenAndServe(":6060", nil) 
// 	if eh != nil { 
// 		panic("ListenAndServe: " + eh.String()) 
// 	} 

// 	e := pprof.StartCPUProfile(os.Stdout)
// 	if e != nil {
// 		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", e)
// 		os.Exit(1)
// 	}
// 	ifile := "/Users/pignatelli/src/tests/gi_taxid_prot.dmp"
// 	ofile := "/Users/pignatelli/src/tests/gi_taxid_prot.24bits.bin"

// 	fmt.Fprintf(os.Stderr, "Creating the dictionary ... ")
// 	t1 := time.Nanoseconds()
// 	m, err := New(ifile)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", err)
// 		os.Exit(1)
// 	}
// 	t2 := time.Nanoseconds()
// 	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2 - t1)/1e9)

// 	fmt.Fprintf(os.Stderr, "Storing the dictionary ... ")
// 	t2 = time.Nanoseconds()
// 	err = m.Store(ofile)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", err)
// 		os.Exit(1)
// 	}	
// 	t2 = time.Nanoseconds()
// 	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", float32(t2 - t1)/1e9)
// 	pprof.StopCPUProfile()

//	m.Encode(3, 1234)
//	m.Encode(4, 4321)
//	fmt.Printf("%v\n", m)
// 	v, _ := m.Decode(46)
// 	fmt.Printf("%v -- %d\n", v, v)
// 	v, _ = m.Decode(1909)
// 	fmt.Printf("%v -- %d\n", v, v)
//}
