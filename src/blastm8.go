package blastm8

import (
	"fmt"
	"os"
	"bufio"
	//	"io";
	"bytes"
	"strconv"
//	"runtime"
	"regexp"
)

//func Init() { runtime.GOMAXPROCS(1) }

const (
	lcalim = 0.1
	blim   = 1000
	scLim  = 0.9
)

type Blast struct {
	query, subject         []byte
	qfrom, qto, sfrom, sto int
	eval, bitsc            float
}

type HType struct {
	GI int
	BS float
}

type Rep2LCA struct {
	Query  []byte
//	bscLim float
	Hits   []*HType
}

func (t Rep2LCA) String() string {
	s := fmt.Sprintf("%s\n", t.Query)
	for _, v := range t.Hits {
		if v.GI != 0 {
			s += fmt.Sprintf("%d %d", v.GI, v.BS);
		}
	}
	s += "\n"
	return s
}

/* Not used */
func bestOf(qC []*Blast) []*Blast {
	var n int
	var v *Blast
	for n, v = range qC {
		if v == nil {
			break
		}
	}
	nlim := float(n) * lcalim
	fmt.Println("..", nlim)
	return (qC[0:int(nlim)])
}

func (b *Blast) extractGI() (gi int) {
	giRx := regexp.MustCompile(`gi\|([0-9]+)\|`) // Sacar fuera?
	gis := giRx.FindSubmatch(b.subject)
	if gis == nil || len(gis) == 0 {
		fmt.Fprintf(os.Stderr, "No GI found in blast record\n")
		os.Exit(1)
	}
	if len(gis) > 2 {
		fmt.Fprintf(os.Stderr, "More than one GI found in blast record. Only the first will be used\n")
	}
	gi, converr := (strconv.Atoi(fmt.Sprintf("%s", gis[1])))
	if converr != nil {
		fmt.Fprintf(os.Stderr, "Error converting %s to number (GI)\n", gis[1])
		os.Exit(1)
	}
	return
}

func Procfile(iblast *bufio.Reader, queryChan chan<- *Rep2LCA, done chan<- bool) {
	hitCollect := new(Rep2LCA)
	var nextQ *Blast
LOOP : for i := 0; ; i++ {
		line, ierr := iblast.ReadString('\n')
		if ierr == os.EOF {
			queryChan <- hitCollect
			//			close(queryChan);
			done <- true
			return
		}
		if i>=blim {
		        i = 0
			continue LOOP
		}
		nextQ = parseblast([]byte(line[0 : len(line)-1]))
		if hitCollect.Query == nil {
			hitCollect.Hits = make([]*HType, blim)
			//			hitCollect.bscLim = nextQ.bitsc * scLim
			hitCollect.Query = nextQ.query
		}
		if (fmt.Sprintf("%s", hitCollect.Query)) != (fmt.Sprintf("%s", nextQ.query)) {
			queryChan <- hitCollect
			i = 0
			hitCollect = new(Rep2LCA)
			hitCollect.Query = nextQ.query
			//			hitCollect.bscLim = nextQ.bitsc * scLim
			hitCollect.Hits = make([]*HType, blim)
		}
		//		if nextQ.bitsc >= hitCollect.bscLim {
		hitCollect.Hits[i] = &HType{GI:nextQ.extractGI(), BS:nextQ.bitsc}
//		}
	}
	return
}

func parseblast(line []byte) *Blast {
	var newB *Blast
	parts := bytes.Split(line, []byte("\t"), -1) // For recent releases, n should be < 0
	//	fmt.Println(parts);
	/* When the bitscore is an integer, has a leading space -- weird, huh? */
	var bitscStr []byte
	if parts[11][0] == 32 {
		bitscStr = parts[11][1:]
	} else {
		bitscStr = parts[11]
	}

	qfrom, _ := strconv.Atoi(fmt.Sprintf("%s", parts[6]))
	qto, _ := strconv.Atoi(fmt.Sprintf("%s", parts[8]))
	sfrom, _ := strconv.Atoi(fmt.Sprintf("%s", parts[7]))
	sto, _ := strconv.Atoi(fmt.Sprintf("%s", parts[9]))
	eval, _ := strconv.Atof(fmt.Sprintf("%s", parts[10]))
	bitsc, bse := strconv.Atof(fmt.Sprintf("%s", bitscStr))
	if bse != nil {
		fmt.Fprintf(os.Stderr, "Error parsing bit score: %s as integer --- Aborting\n", parts[11])
		os.Exit(1)
	}

	newB = &Blast{
		query:   parts[0],
		subject: parts[1],
		qfrom:   qfrom,
		qto:     qto,
		sfrom:   sfrom,
		sto:     sto,
		bitsc:   bitsc,
		eval:    eval}
	return newB
}

// func main() {
// 	blastf,eopen := os.Open("test.blast",os.O_RDONLY,0644);
// 	if eopen != nil {
// 		fmt.Fprintf(os.Stderr,"file doesn't exist");
// 		os.Exit(1);
// 	}
// 	blastbuf := bufio.NewReader(io.Reader(blastf));
// 	queryChan := make(chan *Rep2LCA,1000000);
// 	doneChan := make (chan bool, 4);
// 	go nextblast (blastbuf,queryChan,doneChan);
// 	<- doneChan;
// 	for k := range queryChan {
// 		fmt.Printf("%s\n",k.query)
// 		fmt.Printf("%.3f\n",k.bscLim)
// 		for _,v := range k.hits {
// 			if v != 0 {
// 				fmt.Println(v);
// 			}
// 		}
// 		fmt.Println("")
// 	}

// }
