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

func ProcFile (iblast *bufio.Reader) {
	lines := 0
	recs := 0
	t1 := time.Nanoseconds()
	readQuery := make([]byte, 30)
	readBlock := make([][]byte, 100) // TODO: set length to blim
	procBloc := make(chan [], 4)
	for {
		line, ierr := iblast.ReadBytes('\n')
		if ierr == os.EOF {
			fmt.Fprintf(os.Stdout, "Query! %s\n", readQuery)
//			_ = procBloc(readBlock
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

func parseblast (line Line) [][]byte {
	parts := bytes.Split(line, []byte("\t"))
	return parts
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

func ProcBloc (block [][]byte) 