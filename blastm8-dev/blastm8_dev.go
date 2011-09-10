package blastm8_dev

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"bytes"
	"log"
	"strconv"
)

type Line []byte
type Header []byte
type Hit struct {
	subject Header
	gi int
	qfrom, qto, sfrom, sto int
	bitsc int // float??
}
type QueryRes struct {
	Query Header
	Best *Hit
	Hits []*Hit
}

func ProcFile (iblast *bufio.Reader) {
	lines := 0
	recs := 0
	t1 := time.Nanoseconds()
	readQuery := make([]byte, 30)
	readBlock := make([][]byte, 100) // TODO: set length to blim
	procBlocChan := make(chan *QueryRes, 4)
	for {
		line, ierr := iblast.ReadBytes('\n')
		if ierr == os.EOF {
			fmt.Fprintf(os.Stdout, "Query! %s\n", readQuery)
			go procBloc(readQuery, readBlock, procBlocChan)
			break
		}
		lines++
		currQuery,err := Line(line).extractQuery()
		if err != nil {
			log.Fatal(err)  // TODO: Convert it to non-fatal
		}
		if readQuery == nil {
			readQuery = currQuery
			continue
		}
		if ! bytes.Equal(readQuery,currQuery) {
			recs++
//			_ = procBloc(readBlock)
			fmt.Fprintf(os.Stdout, "Query! %s\n", readQuery)
			readBlock = make([][]byte, 100)
			readQuery = currQuery
			continue
		}
		readBlock = append(readBlock, line)
//		fmt.Printf("%s\n", line)
	}
	t2 := time.Nanoseconds()
	elapse := float32(t2 - t1) / 1e9
	linesXsec := int(float32(lines) / elapse)
	fmt.Fprintf(os.Stderr, "%d recs (%d lines) in %.3f secs (%d lines per sec)\n",
		recs, lines, elapse, linesXsec)
	
}

func parseLine (line Line) *Hit {
	parts := bytes.Split(line, []byte("\t"))
	var bitscStr []byte
	if parts[11][0] == ' ' {
		bitscStr = parts[11][1:]
	} else {
		bitscStr = parts[11]
	}

	qfrom, _ := strconv.Atoi(string(parts[6]))
	qto  , _ := strconv.Atoi(string(parts[7]))
	sfrom, _ := strconv.Atoi(string(parts[8]))
	sto  , _ := strconv.Atoi(string(parts[9]))
	bitsc, _ := strconv.Atof32(string(parts[10]))
	if qfrom > qto {
		qfrom, qto = qto, qfrom
	}
	if sfrom > sto {
		sfrom, sto = sto, sfrom
	}
	gi, gierr := Header(parts[1]).extractGI()
	if gierr != nil {
		log.Fatal(gierr)
	}
	return &Hit {
	gi: gi,
        qfrom: qfrom,
        qto: qto,
        sfrom: sfrom,
        sto: sto,
        bitsc: bitsc }
	}
}

func (b Header) extractGI () (int, os.Error) {
	gib := make([]byte, 0, 10)
	for i,v := range b {
		if v == 'g' && b[i+1] == 'i' && b[i+2] == '|' {
			for j:=i+3; j<len(b); j++ {
				if b[j] == '|' {
					gi, err := strconv.Atoi(string(gib))
//					gi := atoi(gib)
					if err != nil {
						return -1, err
					}
					return gi, nil
				}
				gib = append(gib, b[j])
			}
			nerr := os.NewError(fmt.Sprintf("No | found after GI in %s", b))
			return -1, nerr
		}
	}
	lerr := os.NewError(fmt.Sprintf("No gi| found in: %s", b))
	return -1, lerr
}

func (l Line) extractQuery () ([]byte, os.Error) {
	pos := bytes.IndexByte(l, '\t')
	if pos < 0  { // -1 if it is not present
		return nil, os.NewError("No tab found")
	}
	return l[:pos], nil
}

func ProcBloc (query []byte, block [][]byte, sendChan chan<- *QueryRes ) {
	newQuery := &QueryRes{}
	newQuery.Best = &Hit{bitsc : 0}
	newQuery.Query = make([]byte, len(query)) // IS THIS ALLOCATION NEEDED??
	copy(newQuery.Query, query)
	for _,hitline := range block {
		nextHit := parseLine(Line(hitline));
		
		if nextHit.bitsc > newQuery.Best {
			newQuery.Best = nextHit
		}
	}
}
