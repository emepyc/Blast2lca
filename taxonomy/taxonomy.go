package taxonomy

import (
	"fmt"
	"os"
	"io"
	"bytes"
	"bufio"
	"strconv"
	"gob"
	"time"
	"math"
	"Blast2lca/giTaxid"
)
// TODO : Factor out LCA code? If not in a different package, at least in a different source file.
const sep = "\t|\t"
const maxNodes = 573200  // HINT: Maximum number of nodes in nodes.dmp (from wc -l)
var no_rank []byte = []byte{'n','o','_','r','a','n','k'}
var uc_ []byte = []byte{'u','c','_'}
var unknown []byte = []byte{'u','n','k','n','o','w','n'}

var sortedLevels map[string]int = map[string]int {  // TODO: Could we do this programmatically??
	"forma"              : 0,
	"varietas"           : 1,  
	"subspecies"         : 2,
	"species"            : 3,
	"species subgroup"   : 4,
	"species group"      : 5,
	"subgenus"           : 6,
	"genus"              : 7,
	"subtribe"           : 8,
	"tribe"              : 9,
	"subfamily"          : 10,
	"family"             : 11,
	"superfamily"        : 12,
	"parvorder"          : 13,
	"infraorder"         : 14,
	"suborder"           : 15,
	"order"              : 16,
	"superorder"         : 17,
	"infraclass"         : 18,
	"subclass"           : 19,
	"class"              : 20,
	"superclass"         : 21,
	"subphylum"          : 22,
	"phylum"             : 23,
	"superphylum"        : 24,
	"subkingdom"         : 25,
	"kingdom"            : 26,
	"superkingdom"       : 27,
}

var (
	namesfile, nodesfile, outfile, dictfile string
	helpflag                                bool
)

// Input options for the constructor
type InOpts struct {
	Nodes, Names, Dict, Bintax string
	TaxIsBin, Savemem bool
}

type Pathnode struct {
	Name []byte
	Taxon []byte
}

type Taxnode struct {
	id     int
	Taxid  int
	Parent int
	Childs []int
	Name   []byte
	Taxon  []byte
}

type auxNode struct {
	_id    int
	id     int
	childs []int
	parent int
	taxon  []byte
}

type taxTree map[int]*Taxnode
type auxTree map[int]*auxNode

type taxonomy struct {
	T taxTree
	G giTaxid.GiMapper
	D map[int]int  // from values to indexes
	E, L, H []int
	M [][]int
}

func (n *Pathnode) String () string {
	return fmt.Sprintf("{taxon:%s, name:%s}",n.Taxon, n.Name)
}

func (n *Taxnode) String () string {
 	retStr := fmt.Sprintf ("\n\t{\n\t\tid:%d\n\t\tTaxid:%d\n\t\tParent:%d\n\t\tName:%s\n\t\tTaxon:%s\n\t\tnext:",n.id,n.Taxid,n.Parent,n.Name,n.Taxon)
	for _,v := range n.Childs {
		retStr += fmt.Sprintf("%d,", v)
	}
	retStr += "\n\t}\n"
	return retStr
 }

func (t taxTree) String () string {
	retStr := ""
	for k,v := range t {
		retStr += fmt.Sprintf("%d%s\n",k,v)
	}
	return retStr
}

func (n *auxNode) String () string {
 	retStr := fmt.Sprintf ("\n\t{\n\t\t_id:%d\n\t\tid:%d\n\t\tprev:%d\n\t\ttaxon:%s\n\t\tnext:",n._id,n.id,n.parent, n.taxon)
	for _,v := range n.childs {
		retStr += fmt.Sprintf("%d,", v)
	}
	retStr += "\n\t}\n"
	return retStr
 }

func (t auxTree) String () string {
	retStr := ""
	for k,v := range t {
		retStr += fmt.Sprintf("%d%s\n",k,v)
	}
	return retStr
}

func log2 (x int) int {
	return int (math.Log(float64(x)) / math.Log(float64(2)))
}

func makeMatrix () [][]int {
	dim := maxNodes * 2
	mat := make([][]int, dim)
//	fmt.Printf("Creating matrix of size: %d x %d", maxNodes, dim2)
	for i:=0; i < dim; i++ {
		mat[i] = make([]int, log2(dim)+1)
	}
	return mat
}

func parsename(l []byte) (int, []byte, os.Error) {
	parts := bytes.Split(l, []byte(sep), -1)
	tid, e := strconv.Atoi(string(parts[0]))
	if e != nil {
		return 0, nil, e
	}
	return tid, parts[1], nil
}

func (t taxTree) loadNames (fname string, dict map[int]int) os.Error {
	namesf, eopen := os.OpenFile(fname, os.O_RDONLY, 0644)
	defer namesf.Close()
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file doesn't exist %s\n", fname)
		return eopen
	}
	b := bufio.NewReader(io.Reader(namesf))
	for {
		line, _,  err := b.ReadLine()
		if err == os.EOF {
			return nil
		}
		if pos := bytes.Index(line, []byte("scientific name")); pos == -1 {
			continue
		}
		gi, name, e := parsename([]byte(line)[0 : len(line)-2]) // HINT: ends in "\t|"
		if e != nil {
			return e
		}
		t[dict[gi]].Name = make([]byte, len(name))
		copy (t[dict[gi]].Name, name)
	}
	return nil
}

func (t auxTree) addIdx (node int, idx *int, corrs map[int]int) {
	t[node]._id = *idx
	corrs[node] = *idx
	*idx = *idx+1
	for _,nextNode := range t[node].childs {
		t.addIdx(nextNode, idx, corrs)
	}
}

func (t auxTree) restoreRels(node int, corrs map[int]int) taxTree {
	taxtree := make(taxTree, len(t))
	for _,v := range(t) {
		newchilds := make([]int, 0, len(v.childs))
		for _, c := range v.childs {
			newchilds = append(newchilds, corrs[c])
		}
		taxtree[v._id] = &Taxnode{
		id : v._id,
		Taxid : v.id,
		Parent : corrs[v.parent],
		Childs : newchilds,
		Taxon : v.taxon,
		}
	}
	return taxtree
}

func newTaxonomy (b *bufio.Reader) (taxTree, map[int]int, os.Error) {
	auxtree :=  make(auxTree, maxNodes)
	i := 0
	for {
		l, _, err := b.ReadLine()
		i++
		if err == os.EOF {
			corrs := make(map[int]int, maxNodes)
			var idx int = 1
			auxtree.addIdx(1, &idx, corrs) // HINT: Args : first node in tree, first idx to assign, correspondences between old ids and new ids
			taxtree := auxtree.restoreRels(1, corrs)
			return taxtree, corrs, nil
		}
		if err != nil {
			return nil, nil, err
		}
		line := l[:len(l)-2]   // HINT : Ends with \t|
//		newNode := parsenode([]byte(line)[0 : len(line)-3]) // ends with \t\|\n
//		taxmap.T[newNode.This.Tids.Taxid] = newNode
		parts := bytes.Split(line, []byte(sep),4)
		this, ae1 := strconv.Atoi(string(parts[0]))
		if ae1 != nil {
			return nil, nil, ae1
		}
		that, ae2 := strconv.Atoi(string(parts[1]))
		if ae2 != nil {
			return nil, nil, ae2
		}
		if this == that {
			continue   // To avoid circular references in the tree
		}
		if _, ok := auxtree[this]; ok {
			auxtree[this].parent = that
			auxtree[this].taxon = make([]byte, len(parts[2]))
			copy (auxtree[this].taxon, parts[2])
		} else {
			auxNode := &auxNode {
			id     : this,
			parent : that,
			childs : []int{},
			}
			auxNode.taxon = make([]byte,len(parts[2]))
			copy (auxNode.taxon, parts[2])
			auxtree[this] = auxNode
		}
		if _, ok := auxtree[that]; ok {
			auxtree[that].childs = append(auxtree[that].childs, this)
		} else {
			auxNode := &auxNode {
			id : that,
			childs : []int{this},
			}
			auxtree[that] = auxNode
		}
	}
	return nil, nil, nil  // never used
}

func New(opts *InOpts) (*taxonomy, os.Error) {
	// T : taxtree => tax
	t := &taxonomy{}
	if opts.TaxIsBin {
		s1 := time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Loading taxonomy tree ... ")
		var e os.Error
		t,e = Load(opts.Bintax)
		if e != nil {
			return nil, e
		}
		s2 := time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", (float32(s2-s1))/1e9) 
	} else {
		fmt.Fprintf(os.Stderr, "Creating new taxonomy tree ... ")
		s1 := time.Nanoseconds()
		nodesf, eopen := os.OpenFile(opts.Nodes, os.O_RDONLY, 0644)
		if eopen != nil {
			return nil, eopen
		}
		nodesbuf := bufio.NewReader(io.Reader(nodesf))

		tax, dict, ee := newTaxonomy(nodesbuf)
		if ee != nil {
			return nil, ee
		}
		s2 := time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", (float32(s2-s1))/1e9)

		// T.Name(s) 
		fmt.Fprintf(os.Stderr, "Filling names in taxonomy tree ... ")
		s1 = time.Nanoseconds()
		err := tax.loadNames(opts.Names, dict)
		if err != nil {
			return nil, err
		}
		s2 = time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", (float32(s2-s1))/1e9)

		// E, L, H
		fmt.Fprintf(os.Stderr, "Creating indexes ... ")
		s1 = time.Nanoseconds()
		E := make([]int, 0, maxNodes*2)
		L := make([]int, 0, maxNodes*2)
		H := make([]int, maxNodes)
		tax.ELH (1,0, &E, &L, &H)
//		fmt.Fprintf(os.Stderr, "Selection of H:\n%v\n", H[301502]); os.Exit(0)

		fmt.Fprintf(os.Stderr, "H has grown to : %d\n", len(H))
		fmt.Fprintf(os.Stderr, "For pos %d, value in H is %d\n", 301501, H[301501])
//		os.Exit(0)
		s2 = time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", (float32(s2-s1))/1e9)

		// M
		fmt.Fprintf(os.Stderr,"Preprocessing RMQ ... ")
		s1 = time.Nanoseconds()
		M := makeMatrix()
		RMQprep(&M, L[:len(L)-1], maxNodes)
		s2 = time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", (float32(s2-s1))/1e9)

		fmt.Fprintf(os.Stderr, "Combining data ... ")
		s1 = time.Nanoseconds()
		t.T = tax
		t.D = dict
		t.E = make([]int, len(E))
		copy (t.E, E)
		t.L = make([]int, len(L))
		copy (t.L, L)
		t.H = make([]int, len(H))
		copy (t.H, H)
		t.M = make([][]int, len(M))
		for i:=0; i<len(M); i++ {
			t.M[i] = make([]int, len(M[0]))
		}
		copy (t.M, M)
		s2 = time.Nanoseconds()
		fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", (float32(s2-s1))/1e9)
	}
	fmt.Fprintf(os.Stderr, "Incorporating Gi Taxid mappings ... ")
	s1 := time.Nanoseconds()
	var err os.Error
	t.G, err = giTaxid.Load(opts.Dict, opts.Savemem)
	if err != nil {
		return nil, err
	}
	s2 := time.Nanoseconds()
	fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", (float32(s2-s1))/1e9)
	return t, nil
}

func (t taxTree) ELH (node, level int, E, L, H *[]int) {
	n,_ := t[node]
	if n.id == 301502 {
		fmt.Fprintf(os.Stderr, "n.id == %d -- len(E) == %d curr value in H[%d] == %d\n", 301502, len(*E), 301502, (*H)[n.id-1])
	}
	*E = append(*E, n.id)
	*L = append(*L, level)
	if (*H)[n.id-1] == 0 {
		(*H)[n.id-1] = len(*E)
	}
	for _,nextNode := range n.Childs {
		t.ELH(nextNode, level+1, E, L, H)
	}
	*E = append(*E, n.Parent)
	*L = append(*L, level-1)
}

func RMQprep (M *[][]int, A []int, N int) {
	for i := 0; i < N; i++ {
		(*M)[i][0] = i
	}
	for j := 1; 1 << uint(j) <= N; j++ {
		for i := 0; i + (1 << uint(j)) - 1 < N; i++ {
			if (A[(*M)[i][j - 1]] < A[(*M)[i + (1 << (uint(j) - 1))][j - 1]]) {
				(*M)[i][j] = (*M)[i][j - 1]
			} else {
				(*M)[i][j] = (*M)[i + (1 << (uint(j) - 1))][j - 1]
			}
		}
	}
}

func RMQ (M [][]int, A []int, i, j int) (rmq int) {
	k := int(log2(j-i+1))
//	k := int(math.Log(float64(j - i + 1)))
	if A[M[i][k]] <= A[M[j-int(math.Pow(float64(2),float64(k))) + 1][k]] {
		rmq = M[i][k]
	} else {
		rmq = M[j - int(math.Pow(float64(2),float64(k))) + 1][k]
	}
	return
}

func (t taxonomy) LCA (values ...int) (*Taxnode, os.Error) {
	indexes := make([]int, 0, len(values))
	for _, v := range values {   // from values to indexes
		if _,ok := t.D[v]; ok { // HINT -- There may be taxids not in taxonomy
			indexes = append(indexes,t.D[v])
		}
	}
	if len(indexes) == 0 {
		return &Taxnode{}, os.NewError("EMPTY")
	}
	red := indexes[0]
	indexes = indexes[1:]
	for _,ind := range indexes {
//		fmt.Fprintf(os.Stderr, "LCA so far:%v VS %v\n", red, ind)
		red = LCAhelper(t.E,t.L,t.H,t.M,red,ind)
	}
	return t.T[red], nil
}

func LCAhelper (E,L,H []int, M [][]int, i, j int) int {
	v1 := H[i-1]
	v2 := H[j-1]
	if v1 > v2 {
		v1, v2 = v2, v1
	}
	rmq := RMQ(M, L, v1, v2)
//	fmt.Printf("RMQ: %d\n", rmq)
	lca := E[rmq]
	return lca
}

func (t taxonomy) Node (taxid int) *Taxnode {
	id, ok := t.D[taxid]
	if !ok {
		return nil 
	}
	node, ok := t.T[id]
	if !ok {
		return nil
	}
	return node
}

func (t taxonomy) Path (taxid int) []*Pathnode {
	node := t.Node(taxid)
	path := make([]*Pathnode, 0, 10)
	for {
		if node.Taxid == 1 {
			return path
		}
		newNode := &Pathnode{}
		newNode.Name = make([]byte, len(node.Name))
		newNode.Taxon = make([]byte, len(node.Taxon))
		copy(newNode.Name, node.Name)
		copy(newNode.Taxon, node.Taxon)
		path = append(path, newNode)
		node = t.Parent(node)
	}
	return nil
}

func (t *taxonomy) PathFromGi(gi int) ([]*Pathnode, os.Error) {
	taxid, err := t.TaxidFromGi(gi)
	if err != nil {
		return nil, err
	}
	return t.Path(taxid), nil
}

func (t *taxonomy) TaxidFromGi(gi int) (int, os.Error) {
	taxid, err := t.G.GiTaxid(gi)
	if err != nil {
		return -1, err
	}
	return taxid, nil
}

func (t *taxonomy) Parent (node *Taxnode) *Taxnode {
	return t.T[node.Parent]
}

func (t *taxonomy) AtLevels (node *Taxnode, levs... []byte) [][]byte {
	taxAtLevels := make ([][]byte, 0, len(levs))
//	fmt.Println("NODE : ", node)
//	os.Exit(0)
	allLevs := t.AllLevels(node)           // HINT: Taxonomy levels to names
	baseLevN,_ := sortedLevels[string(node.Taxon)]  // HINT: Levels below this are "uc_"
	baseLev := append(uc_, node.Name...)
//	fmt.Fprintf(os.Stderr, "baseLevN : %s (%d)\n", baseLev, baseLevN)
	for _, lev := range levs {
//		fmt.Fprintf(os.Stderr, "LEV %s\n", lev)
		if atLev, ok := allLevs[string(lev)]; ok {
			taxAtLevels = append(taxAtLevels, atLev)
			continue
		} 
		// If not ... 2 possible causes: i) too low level ("uc_")
		if sortedLevels[string(lev)] < baseLevN {
			taxAtLevels = append(taxAtLevels, baseLev)
			continue
		}
		// ii) No taxon for the LCA at the required level -- give the first known upstream
		
		
		taxAtLevels = append(taxAtLevels, unknown)
	}
//	s2 := time.Nanoseconds()
//	fmt.Fprintf(os.Stderr, "Done (%.3f secs)\n", (float32(s2-s1)/1e9))
	return taxAtLevels
}

func (t *taxonomy) AllLevels (node *Taxnode) map[string][]byte {
//	fmt.Println("NODE2: ",node)
	taxons := make(map[string][]byte, 15)
	for {
//		fmt.Fprintf(os.Stderr, "TAXID: %d\n", node.Taxid)
		if node.Taxid == 1 {
			return taxons
		}
		if bytes.Compare (node.Taxon, no_rank) != 0 {
			taxons[string(node.Taxon)] = node.Name
		}
		node = t.Parent(node)
	}
	return nil // Should be never used
}

func (t *taxonomy) AtLevel (node *Taxnode, lev []byte) []byte {
	for {
		if bytes.Compare(node.Taxon, lev) == 0 {
			return node.Name
		}

		if l, ok := sortedLevels[string(node.Taxon)]; ok && l > sortedLevels[string(lev)] { 
			return append(uc_, node.Name...)
		}
		
		node = t.Parent(node)
	}
	return []byte("no rank++")
}


func Load (fname string) (*taxonomy, os.Error) {
	fh, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	t := &taxonomy{}
	dec := gob.NewDecoder(fh)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t taxonomy) Store (fname string) (int, os.Error) {
	t.G = nil   // We Store without GiTaxid mappings
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(t)
	if err != nil {
		return 0, err
	}
	fh, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
	defer fh.Close()
	if err != nil {
		return 0, err
	}

	n,err := fh.Write(b.Bytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}


// func main() {
// 	opts := &InOpts{
// 	Nodes : nodesfile,
// 	Names : namesfile,
// 	Dict : dictfile,
// 	TaxIsBin : false,
// 	Savemem : true,
// 	}
// 	newtax, err := New(opts)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "\n ERROR: the taxonomy is empty\n")
// 		os.Exit(1)
// 	}

// //	fmt.Println(newtax.Path(67971))
// //	os.Exit(0)
// //	fmt.Println(newtax)
// 	fmt.Println("Calculating Lca")
// 	lca := newtax.LCA(542295, 316120, 258110)
// 	fmt.Printf("LCA : %v\n", lca)
// 	fmt.Printf("PATH: %v\n", newtax.Path(lca.Taxid))
// 	lca_atlevel := newtax.AtLevel(lca, []byte("species"))
// 	fmt.Printf("At Level Phylum : %s\n", lca_atlevel)
// //	fmt.Printf("Parent: %v\n", newtax.GetParent(lca))
// 	os.Exit(0)

// 	fmt.Fprintf(os.Stderr, "Storing all files into binary file ... ")
// 	s1 := time.Nanoseconds()
// 	n, e := newtax.Store(outfile)
// 	s2 := time.Nanoseconds()
// 	if e != nil {
// 		fmt.Fprintf(os.Stderr, "\n ERROR: %s\n", e)
// 		os.Exit(1)
// 	}
// 	fmt.Fprintf(os.Stderr, "Done (%.3f sec) -- %d bytes written\n", (float32(s2-s1))/1e9, n)

// 	fmt.Fprintf(os.Stderr, "Loading the binary file back ... ")
//  	s1 = time.Nanoseconds()
//  	_,err = Load(outfile)
//  	s2 = time.Nanoseconds()
//  	if err != nil {
//  		fmt.Fprintf(os.Stderr, "\n ERROR: %s\n",err)
//  		os.Exit(1)
//  	}
//  	fmt.Printf("Done (%.3f sec)\n", float32(s2-s1)/1e9)
// }

