package taxonomy

import (
	"fmt"
	"os"
	"io"
	"bytes"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
//	"encoding/binary"
//	"runtime"
)

//func Init() { runtime.GOMAXPROCS(4) }

const sep = "\t|\t"

type tids struct {
	taxid int
	name  string
}

type taxon struct {
	*tids
	taxon string
}

type Taxnode struct {
	this     taxon
	parentId int
}

type taxonomy struct {
	tree  map[int]*Taxnode
	gimap []uint8
	//	overLevels map[[]byte]map[[]byte]int
}

func (t *taxonomy) OverLevel(level []byte) *map[string]int {
	lstr := fmt.Sprintf("%s", level)
	mm := make(map[string]int)
	for k, v := range t.tree {
		if t.tree[k].this.taxon == lstr {
			mm[t.tree[v.parentId].this.taxon] = 1
		}
	}
	mm["no rank"] = 0, false
	return &mm
}

func (t *taxonomy) AtLevel(id int, lev []byte) []byte {
	if id == 1 {
		return ([]byte("uc_no_rank"))
	}
	overL := *t.OverLevel(lev)
	//	fmt.Println(overL)
	lstr := fmt.Sprintf("%s", lev)
	var node *Taxnode
	for {
		if id == 1 || id == 131567 || id==0 {
			return []byte("uc_no_rank")
		}
		node = t.tree[id]
		if node.this.taxon == lstr {
			return []byte(node.this.tids.name)
		}
		if _, ok := overL[node.this.taxon]; ok {
			return []byte("uc_" + node.this.tids.name)
		}
		id = node.parentId
	}
	return ([]byte("no rank"))
}

func OutPath(path []*Taxnode) {
	for _, n := range path {
		if n.parentId == -1 {
			break
		}
		if n.parentId == 1 {
			break
		}
		fmt.Print(n.this.tids.name, "(", n.this.taxid, ") ")
	}
	fmt.Println()
}

func usedSlots(t [][]*Taxnode) int {
	for p, v := range t {
		if v == nil {
			return (p)
		}
	}
	return len(t)
}

func LCA(taxs [][]*Taxnode) int {
	levels := make(map[int]int)
	nslots := usedSlots(taxs)
	for _, t := range taxs {
		for _, n := range t {
			if n == nil {
				break
			}
			if n.parentId == -1 { nslots--; continue }
			if _, present := levels[n.this.taxid]; present {
				levels[n.this.taxid] += 1
			} else {
				levels[n.this.taxid] = 1
			}
		}
	}
	for _, n := range taxs[0] {
		if n.parentId == -1 {
			return 1
		}
		if n.parentId == 1 {
			return n.this.taxid
		}
		if n.this.taxid == 1 {
			return 1
		}
		if levels[n.this.taxid] == nslots {
			return n.this.taxid
		}
	}
	return 0
}

func (t *taxonomy) Path(id int) []*Taxnode {
	path := make([]*Taxnode, 200)
	var node *Taxnode
	var ok bool
	pos := 0
	for {
		node, ok = t.tree[id]
		if !ok {
			fmt.Fprintf(os.Stderr, "%d is not found\n", id)
			path[pos] = &Taxnode{parentId:-1}
			return path
//			id = 2
//			node, _ = t.tree[2]
			//			os.Exit(1)
		}
		path[pos] = node
		if node.parentId == 1 {
			return path
		}
		id = node.parentId
		pos++
	}
	return path // superfluo
}

func (t *taxonomy) PathFromGi(gi int) []*Taxnode {
	taxid := t.TaxidFromGi(gi)
	tax := t.Path(int(taxid))
	return tax
}

func (t *taxonomy) TaxidFromGi(gi int) int {
	pos := gi*4
	b := t.gimap
	if len(b) < pos+3 {
		fmt.Fprintf(os.Stderr, "GI not found: %d\n", gi)
		return -1
	}
	taxid := int (uint32(b[pos+3]) | uint32(b[pos+2])<<8 | uint32(b[pos+1])<<16 | uint32(b[pos])<<24)
	if taxid == 0 {
		fmt.Fprintf(os.Stderr, "GI not found: %d\n", gi)
		return -1	
	}
	return taxid
}

// func (t *taxonomy) TaxidFromGi(gi int) int {
// 	t.gimap.Seek(int64(gi)*4, 0)
// 	var taxid uint32
// 	err := binary.Read(t.gimap, binary.BigEndian, &taxid)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "GI not found: %i", gi)
// 		return 2
// 		//		os.Exit(1)
// 	}
// 	return int(taxid)
// }

func readnames(b *bufio.Reader, namesch chan<- *tids, bch chan<- bool) {
	for {
		line, err := b.ReadString('\n')
		if err == os.EOF {
			bch <- true
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
	bch <- true
	close(namesch)
	return
}

func parsename(l []byte) *tids {
	parts := bytes.Split(l, []byte(sep), -1)
	tid, _ := strconv.Atoi(fmt.Sprintf("%s", parts[0]))
	newName := &tids{
		tid,
		fmt.Sprintf("%s", parts[1])}
	return newName
}

func parsenode(l []byte) *Taxnode {
	parts := bytes.Split(l, []byte(sep), 4)
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

func LoadDict(fname string) []uint8 {
	b,rerr := ioutil.ReadFile(fname)
	if rerr != nil {
		fmt.Fprintf(os.Stderr, "File %s not found\n",fname)
		os.Exit(1)
	}
	return b
}

func newTaxonomy(b *bufio.Reader, g string) *taxonomy {
	
	taxmap := &taxonomy{
	tree:  make(map[int]*Taxnode),
	gimap: LoadDict(g)}
	for {
		line, err := b.ReadString('\n')
		if err == os.EOF {
			return taxmap
		}
		newNode := parsenode([]byte(line)[0 : len(line)-3]) // ends with \t\|\n
		//		fmt.Println("KK:",newNode.this.taxid);
		taxmap.tree[newNode.this.taxid] = newNode
	}
	return taxmap
}

func New(nodes, names, dict string) *taxonomy {
	nodesf, eopen := os.Open(nodes, os.O_RDONLY, 0644)
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file doesn't exist %s\n", nodes)
		os.Exit(1)
	}
	nodesbuf := bufio.NewReader(io.Reader(nodesf))

	namesf, eopen := os.Open(names, os.O_RDONLY, 0644)
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file doesn't exist %s\n", names)
		os.Exit(1)
	}
	namesbuf := bufio.NewReader(io.Reader(namesf))

// 	dictfh, edict := os.Open(dict, os.O_RDONLY, 0644)
// 	if edict != nil {
// 		fmt.Fprintf(os.Stderr, "file doesn't exist %s\n", dict)
// 		os.Exit(1)
// 	}

	endch := make(chan bool, 100)
	namesch := make(chan *tids, 1000000)
	go readnames(namesbuf, namesch, endch)

//	tax := newTaxonomy(nodesbuf, io.ReadSeeker(dictfh))
	tax := newTaxonomy(nodesbuf, dict)
	<-endch

	for n := range namesch {
		tax.tree[n.taxid].this.tids = n
	}

	return tax
}

// func main() {
// 	taxDB := New(
// 		"nodes.dmp",
// 		"names.dmp",
// 		"/home/pignatelli/lib/taxbuild/taxonomy-16-03-2010/gi_taxid_prot.bin")
	
// //	fmt.Println(tax);

// // 	fmt.Println("xx:",tax.tree[201239]);
// // 	for n := range namesch {
// // 		tax.tree[n.taxid].this.tids = n;
// // 	}

//  	path1 := taxDB.PathFromGi(227546195)
//  	path2 := taxDB.PathFromGi(239622198)
//  	path3 := taxDB.PathFromGi(213692709)
// 	path4 := taxDB.PathFromGi(23465596)
// 	path5 := taxDB.PathFromGi(23335262)
// 	path6 := taxDB.PathFromGi(268619926)


// 	OutPath(path1)
// 	OutPath(path2)
// 	OutPath(path3)
// 	OutPath(path4)
// 	OutPath(path5)
// 	OutPath(path6)
// // // 	for _,v := range path1 {
// // // 		if v!=nil {
// // // 			fmt.Print(v," ");
// // // 		}
// // // 	}
// // // 	fmt.Println("\n")
// // // 	for _,v := range path2 {
// // // 		if v!=nil {
// // // 			fmt.Print(v," ");
// // // 		}
// // // 	}
// // // 	fmt.Println("\n")
// // // 	for _,v := range path3 {
// // // 		if v!=nil {
// // // 			fmt.Print(v," ");
// // // 		}
// // // 	}
//  	fmt.Println("\n")
//  	fmt.Println("LCA: ")
//  	lca := LCA([][]*Taxnode{path1,path2,path6,path4,path5,path3})
//  	fmt.Println("lca => ",lca)
// // // 	path := tax.Path(2);
// // // 	for _,v := range path {
// // // 		if v!=nil {
// // // 			fmt.Print(v,"  ");
// // // 		}
// // // 	}

// // 	atL := tax.AtLevel(227961,strings.Bytes("family"));
// // 	fmt.Printf("%s\n",atL);

// }
