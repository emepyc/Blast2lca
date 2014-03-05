//Package blastm8 allows reading and parsing of m8-formatted blast files
package blastm8

import (
	"fmt"
	"io"
	"bufio"
	"sort"
	"bytes"
	"strconv"
	"regexp" //HINT: To extract GIs -- TODO: Profile the 2 alternatives given below
	"log"
	"errors"
//	"os"
)

//giRx is a regexp for GI extraction
var giRx *regexp.Regexp = regexp.MustCompile(`gi\|([0-9]+)\|`) // HINT: Compile globally

//Header represents a Query or Subject -- typical headers of nr/nt as they appear in blast-m8 results
type Header []byte

//Hit gives single Blast hit information
type Hit struct { // Was Blast
	gi               int // We may operate in GI space
	bitsc            float64
}

//Hits represent a  collection of hits
type Hits []*Hit

//BlastBlock represents a block of Blast hits
type BlastBlock struct {
	header Header
	block  []byte
}

//QueryRes has the needed information about the hits of a query
type QueryRes struct {
	Query  Header
	Hits   Hits
}

//findIndex returns the index in the Hits slice with the last significant Hit.
func (hits Hits) findIndex (bsLim float64) int {
	for i,hit := range hits {
		if hit.bitsc < bsLim {
			return i
		}
	}
	return len(hits)
}

//String stringify a BlastBlock
func (b BlastBlock) String() string {
	s := fmt.Sprintf("HEADER: %s\n", b.header)
	s += fmt.Sprintf("BLOCK:\n%s\n", b.block)
	return s
}

//String stringify a hit
func (h Hit) String() string {
	return fmt.Sprintf("GI:%d\t%.2f", h.gi, h.bitsc)
}

//GI returns the GI of the corresponding Hit
func (h *Hit) GI() int {
	return h.gi
}

//Bitsc returns the bit score of the corresponding Hit
func (h *Hit) Bitsc() float64 {
	return h.bitsc
}

//String Stringifies a query result
func (t QueryRes) String() string {
	s := fmt.Sprintf("%s\n", t.Query)
	s += fmt.Sprintf("HITS:\n")
	for _, v := range t.Hits {
		s += fmt.Sprintf("\t%v\n", v)
	}
	return s
}

// Sort interface for Hits by bitscore
func (h Hits) Len() int {
	return len(h)
}

// Sort interface for Hits by bitscore
func (h Hits) Less(i, j int) bool {
	return h[i].bitsc > h[j].bitsc
}

// Sort interface for Hits by bitscore
func (h Hits) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}


// 2 OPTIONS HERE! SEE WHAT IS FASTER / MORE CONSISTENT
//Extract the GI of a Subject/Query
// func (b Header) extractGI() (gi int) {
// 	gis := giRx.FindSubmatch(b)
// 	if gis == nil || len(gis) == 0 {
// 		fmt.Fprintf(os.Stderr, "No GI found in Header: %s => %s\n", b)
// 		gi = -1
// 		return
// 	}
// 	if len(gis) > 2 {
// 		fmt.Fprintf(os.Stderr, "More than one GI found in blast record. Only the first will be used\n")
// 	}
// 	gi, converr := (strconv.Atoi(string(gis[1])))
// 	if converr != nil {
// 		fmt.Fprintf(os.Stderr, "Error converting %s to number (GI)\n", gis[1])
// 		os.Exit(1)
// 	}
// 	return
// }

func (b Header) extractGI () (int, error) {
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
			nerr := errors.New(fmt.Sprintf("No | found after GI in %s", b))
			return -1, nerr
		}
	}
	lerr := errors.New(fmt.Sprintf("No gi| found in: %s", b))
	return -1, lerr
}

//ProcFile reads the query results from a blast m8-formatted file and passes the results
//to the queryChan channel. What is passed is the raw block of lines corresponding to a single query in the blast file.
func Procfile(iblast *bufio.Reader, queryChan chan<- *BlastBlock) {
	subjCollect := bytes.NewBuffer(make([]byte, 0, 100000))
	var query []byte
	for ;; {
		line, _, ierr := iblast.ReadLine() // TODO: Check isIndex
		if ierr == io.EOF {
			block := subjCollect.Bytes()
			queryChan <- &BlastBlock{ header : Header(query), block : block[:len(block)-1] }
			close(queryChan)
			return
		}
		currQuery, qerr := extractQuery(line)
		if qerr != nil {
			log.Printf("WARNING: I can't extract the query field from this line: %s\n%s\n", line, qerr)
			continue // offending line is not passed
		}
		if (query == nil) {
			query = currQuery
		}
		if (bytes.Equal(currQuery, query)) {
			_, err := subjCollect.Write(append(line, '\n'))
			if err != nil {
				log.Println("WARNING: Error collecting line from blast:\nLINE:\n%s\nERROR: %s\n", line, err)
			}
		} else {
			block := subjCollect.Bytes()
			passQuery := make([]byte, len(query))
			copy (passQuery, query)
			queryChan <- &BlastBlock{ header: Header(passQuery), block : block[:len(block)-1] }
			subjCollect = bytes.NewBuffer(make([]byte, 0, 100000))
			subjCollect.Write(append(line, '\n'))
			query = currQuery
		}
	}
}

// extractQuery extracts and returns the query field of a line of m8-formatted blast hit.
func extractQuery (line []byte) ([]byte, error) {
	pos := bytes.IndexByte(line, '\t')
	if pos < 0 {
		return nil, errors.New("Line is not tab separated. This line can't be processed")
	}
	if pos == 0 {
		return nil, errors.New("Line with a blank query field (starts with a <tab> character")
	}
	return line[0:pos], nil
}

// ParseRecord parses the lines for a query (blast m8-formatted) and write the information in a QueryRes
// Only the lines with bit score greater than the best score * scLim are processed
func ParseRecord (bb BlastBlock, scLim float64) *QueryRes {
	qRes := &QueryRes{}
	qRes.Query = bb.header
	bestBs := float64(0)
	recs := bytes.Split(bb.block, []byte{'\n'})
	for _, blastLine := range recs {
		nextHit, err := parseblast(blastLine)
		if err != nil {
			log.Printf("WARNING: Ignoring blast line: %s\n", blastLine)
			continue
		}
		qRes.Hits = append(qRes.Hits, nextHit)
		if bestBs < nextHit.bitsc {
			bestBs = nextHit.bitsc
		}
	}
	bsLim := bestBs * scLim
	sort.Sort(qRes.Hits)
	index := qRes.Hits.findIndex(bsLim)
	qRes.Hits = qRes.Hits[:index]
	return qRes
}


func parseblast(line []byte) (*Hit, error) {
	var newB *Hit
	
	parts := bytes.Split(line, []byte("\t"))
//	parts := getFields(line)
	var bitscStr []byte
	bitscStr := bytes.TrimSpace(parts[11])
	// if parts[11][0] == 32 {
	// 	bitscStr = parts[11][1:]
	// } else {
	// 	bitscStr = parts[11]
	// }

	bitsc, bse := strconv.ParseFloat(string(bitscStr), 64)
	if bse != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing bit score %s as integer: %s\n", parts[11], bse))
	}
	gi, gierr := Header(parts[1]).extractGI()
	if gierr != nil {
		return nil, errors.New(fmt.Sprintf("Error extracting GI from %s: %s", parts[1], gierr))
	}
	newB = &Hit{
	gi: gi,
	bitsc:   bitsc}
	return newB, nil
}

func getFields (line []byte) [][]byte {
	fields := make ([][]byte, 0, 12)
	this := make([]byte, 0, 20)
	for _, ch := range line {
		if ch == 9 { // \t
			fields = append(fields, this)
			this = make([]byte, 0, 20)
			continue
		}
		if ch == 10 { // \n
			fields = append(fields, this)
			return fields
		}
		this = append(this, ch)
	}
	fields = append(fields, this)
	return (fields)
}

//To profile a custom implementation of atoi (-- atoi is one of the most expensive operations in the program
//This may have now less impact, since the conversions have been moved to parallelized parts of the program
func atoi (dbytes []byte) int {
	i := 0
	is_negative := false

	for _, ch := range dbytes {
		if ch < 15 && ch > 26 && ch != 13 { // !is_digit(ch) && ch != '-'
			continue
		}
		if ch == 13 { // -
			is_negative = true
		}
		i = (i*10) + (int(ch) - 16)%10
	}
	if is_negative {
		return i
	}
	return -i
}
