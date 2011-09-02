package blastm8


/*
The blastm8 package reads a tab-formatted blast output and passes each record to a channel
*/

import (
	"fmt"
	"os"
	"bufio"
	"sort" // For transform -- Schwartzian Transformation
	"bytes"
	"strconv"
//	"regexp" //HINT: To extract GIs
//	"getgi"
)

//Global parameters
const (
	lcalim = 0.1  // TODO: All these should be parameterized and imported from main
	blim   = 1000 // ""
	scLim  = 0.9  // ""
	ovl    = 0.8  // 80 % of the shortest sequence
)

//Global vars
// var (
// 	giRx *regexp.Regexp = regexp.MustCompile(`gi\|([0-9]+)\|`) // HINT: Compile globally
// )

//Query or Subject -- typical headers of nr/nt as they appear in blast-m8 results
type Header []byte

//Information of a single Blast hit
type Hit struct { // Was Blast
	subject                Header
	gi                     int // We may operate in GI space
	qfrom, qto, sfrom, sto int
	eval, bitsc            float32
}

type sortByPos Hits
type sortByBitsc Hits

type BySubj struct {
	Recs []*Hit
	Best *Hit // ??
}

//Information about the hits of a query
type QueryRes struct { // Was Rep2LCA
	Query  Header
	Best   *Hit
	Hits   []*Hit
	BySubj []*BySubj
	ByPos  []*Hit
}

//A collection of hits
type Hits []*Hit

//Stringify a hit
func (h Hit) String() string {
	return fmt.Sprintf("%s\t%d\t%d\t%.2f", h.subject, h.qfrom, h.qto, h.bitsc)
}

//Get gi
func (h *Hit) GetGI() int {
	return h.gi
}

//Get bitsc
func (h *Hit) GetBitsc() float32 {
	return h.bitsc
}

//Stringify a query result
func (t QueryRes) String() string {
	s := fmt.Sprintf("%s\nBEST:\n%v\n", t.Query, t.Best)
	s += fmt.Sprintf("HITS:\n")
	for _, v := range t.Hits {
		s += fmt.Sprintf("\t%v\n", v)
	}
	s += fmt.Sprintf("BY_SUBJ\n")
	for i, v := range t.BySubj {
		s += fmt.Sprintf("++++ %d ++++\n", i)
		s += fmt.Sprintf("%v\n", v)
	}
	return s
}

//Stringify a BySubj
func (b *BySubj) String() string {
	s := fmt.Sprintf("BEST:\n")
	s += fmt.Sprintf("\t%v\n", b.Best)
	s += fmt.Sprintf("ALL:\n")
	for _, v := range b.Recs {
		s += fmt.Sprintf("\t%v\n", v)
	}
	return s
}

// Sort interface for hits
func (h sortByPos) Len() int {
	return len(h)
}

func (h sortByPos) Less(i, j int) bool { // fields qFrom and qTo are already sorted
	if h[i].qfrom == h[j].qfrom {
		return h[i].qto < h[j].qto
	}
	return h[i].qfrom < h[j].qfrom
}

func (h sortByPos) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
// Sort interface -- End

// Sort interface by subject
func (h sortByBitsc) Len() int {
	return len(h)
}

func (h sortByBitsc) Less(i, j int) bool {
	return h[i].bitsc > h[j].bitsc
}

func (h sortByBitsc) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

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

func (b Header) extractGI () (int, os.Error) {
	gib := make([]byte, 0, 10)
	for i,v := range b {
		if v == 'g' && b[i+1] == 'i' && b[i+2] == '|' {
			for j:=i+3; j<len(b); j++ {
				if b[j] == '|' {
					gi, err := strconv.Atoi(string(gib))
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

func (queryRes *QueryRes) fillBySubj() {
	queryRes.ByPos = make([]*Hit, len(queryRes.Hits))
	copy(queryRes.ByPos, queryRes.Hits) // Assume Hits are sorted by BS, is this really happening?
	sort.Sort((sortByPos)(queryRes.ByPos))
	queryRes.BySubj = make([]*BySubj, 0, 10)
	currGroup := &BySubj{ Recs : make([]*Hit, 0, 10), Best : &Hit{} }
	currGroup.Best = queryRes.ByPos[0]
	var minTo int = 0
	var maxFrom int = 0
	for _, nextHit := range queryRes.ByPos {
		nextFrom, nextTo := nextHit.qfrom, nextHit.qto
		bestFrom, bestTo := currGroup.Best.qfrom, currGroup.Best.qto
		if nextTo < bestTo {
			minTo = nextTo
		} else {
			minTo = bestTo
		}
		if nextFrom > bestFrom {
			maxFrom = nextFrom
		} else {
			maxFrom = bestFrom
		}
		lenBest := bestTo - bestFrom //HINT: Both are already sorted (To<From)
		lenNext := nextTo - nextFrom

		percBest := float32(minTo - maxFrom) / float32(lenBest)
		percNext := float32(minTo - maxFrom) / float32(lenNext)

		if percBest < ovl && percNext < ovl {
			sort.Sort((sortByBitsc)(currGroup.Recs))  // sorted by bitsc descending order
			queryRes.BySubj = append(queryRes.BySubj, currGroup)
			currGroup = &BySubj{Recs : make([]*Hit, 0, 10), Best : &Hit{}}
			currGroup.Recs = append(currGroup.Recs, nextHit)
			currGroup.Best = nextHit
			continue
		}
		currGroup.Recs = append(currGroup.Recs, nextHit)
		if currGroup.Best.bitsc < nextHit.bitsc {
			currGroup.Best = nextHit
		}
	}
	queryRes.BySubj = append(queryRes.BySubj, currGroup)
}


//Read the query results from a blast m8-formatted file and pass the results
//to a channel
//func Procfile(iblast *bufio.Reader, queryChan chan<- *QueryRes, done chan<- bool, byFunc bool) {
func Procfile(iblast *bufio.Reader, queryChan chan<- *QueryRes, byFunc bool) {
	hitCollect := new(QueryRes)
	hitCollect.Best = &Hit{bitsc : 0}
	var nextHit *Hit
	var query []byte
LOOP:
	for i := 0; ; i++ {
		line, ierr := iblast.ReadBytes('\n')
//		fmt.Fprintf(os.Stderr, "++LINE: %s, ierr: %s\n", line, ierr);
//		if len(line) == 0 {   // HINT: For blank lines
//			continue
//		}
		if ierr == os.EOF {
			if byFunc {
				hitCollect.fillBySubj()
			}
			sort.Sort((sortByBitsc)(hitCollect.Hits))
			queryChan <- hitCollect
			close(queryChan)
			return
		}
		line = line[0:len(line)-1]
		query, nextHit = parseblast(line)
//		fmt.Fprintf(os.Stderr, "Q:%s\n", query)
		if hitCollect.Query == nil {
			hitCollect.Hits = make([]*Hit, 0, blim)
			hitCollect.Query = make([]byte, len(query))
			copy(hitCollect.Query, query)
			fmt.Fprintf (os.Stderr, "Query => %s\n", hitCollect.Query)
		}
		if ! bytes.Equal(hitCollect.Query, query) {
			if byFunc {
				hitCollect.fillBySubj()
			}
			sort.Sort((sortByBitsc)(hitCollect.Hits))
			queryChan <- hitCollect
			hitCollect = new(QueryRes)
			hitCollect.Best = nextHit
			hitCollect.Query = make([]byte, len(query))
			copy (hitCollect.Query, query)
			hitCollect.Hits = make([]*Hit, 0, blim)
			i = 0
		}
		if i >= blim {
			continue LOOP
		}
		if hitCollect.Best.bitsc < nextHit.bitsc {
			hitCollect.Best = nextHit
		}
		hitCollect.Hits = append(hitCollect.Hits, nextHit)
	}
	return
}

//Parses a single blast hit
func parseblast(line []byte) ([]byte, *Hit) {
	var newB *Hit

	parts := bytes.Split(line, []byte("\t")) // For recent releases, n should be < 0
	var bitscStr []byte
	if parts[11][0] == 32 {
		bitscStr = parts[11][1:]
	} else {
		bitscStr = parts[11]
	}

	// TODO: we don't give strand information, should we??
	qfrom, _ := strconv.Atoi(string(parts[6]))
	qto, _ := strconv.Atoi(string(parts[7]))   // Weird ... it was 8!!??
	sfrom, _ := strconv.Atoi(string(parts[8])) // Weird ... it was 7!!??
	sto, _ := strconv.Atoi(string(parts[9]))
	eval, _ := strconv.Atof32(string(parts[10]))  // Used? If not... remove
	bitsc, bse := strconv.Atof32(string(bitscStr))
	qfromto := []int{qfrom, qto}
	sort.Ints(qfromto)
	qfrom, qto = qfromto[0], qfromto[1] // qfrom and qto are now sorted
	sfromto := []int{sfrom, sto}
	sort.Ints(sfromto)
	sfrom, sto = sfromto[0], sfromto[1] // sfrom and sto are now sorted
	if bse != nil {
		fmt.Fprintf(os.Stderr, "Error parsing bit score: %s as integer --- Aborting\n", parts[11])
		os.Exit(1)
	}

	query := parts[0]
	gi, gierr := Header(parts[1]).extractGI()
	if gierr != nil {
		fmt.Fprintf(os.Stderr, "%s", gierr)
		os.Exit(1)
	}
	newB = &Hit{
//		subject: parts[1],
//	gi:      Header(parts[1]).extractGI(),
	gi: gi,
	qfrom:   qfrom,
	qto:     qto,
	sfrom:   sfrom,
	sto:     sto,
	bitsc:   bitsc,
	eval:    eval}
	newB.subject = make([]byte, len(parts[1]))
	copy(newB.subject, parts[1])
	return query, newB
}

