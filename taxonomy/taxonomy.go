// Package taxonomy implements NCBI taxonomy operations.
//
// These operations include loading the taxonomy database in memory 
// from names.dmp and nodes.dmp and basic operations over the database
// like LCA calculation
package taxonomy

import (
//	"log"
	"fmt"
	"os"
	"io"
	"bytes"
	"bufio"
	"strconv"
	"time"
	"math"
	"github.com/emepyc/Blast2lca/giTaxid"
	"github.com/emepyc/Blast2lca/wcl"
)
// TODO : Factor out LCA code in a different source file
const sep = "\t|\t"
// const maxNodes = 809800 // WARNING!! : Maximum number of nodes in nodes.dmp (from wc -l)
// TODO: This shouldn't be hardcoded
var no_rank []byte = []byte{'n', 'o', '_', 'r', 'a', 'n', 'k'}
var uc_ []byte = []byte{'u', 'c', '_'}

//Unknown is the taxon representation of any unknown (or lack of) taxon id
var Unknown []byte = []byte{'u', 'n', 'k', 'n', 'o', 'w', 'n'}

var sortedLevels map[string]int = map[string]int{ // TODO: Could we do this programmatically??
	"forma":            0,
	"varietas":         1,
	"subspecies":       2,
	"species":          3,
	"species subgroup": 4,
	"species group":    5,
	"subgenus":         6,
	"genus":            7,
	"subtribe":         8,
	"tribe":            9,
	"subfamily":        10,
	"family":           11,
	"superfamily":      12,
	"parvorder":        13,
	"infraorder":       14,
	"suborder":         15,
	"order":            16,
	"superorder":       17,
	"infraclass":       18,
	"subclass":         19,
	"class":            20,
	"superclass":       21,
	"subphylum":        22,
	"phylum":           23,
	"superphylum":      24,
	"subkingdom":       25,
	"kingdom":          26,
	"superkingdom":     27,
}

// InOpts collects the Input options for the constructor
type InOpts struct {
	Nodes, Names, Dict   string
	Savemem    bool
}

type pathnode struct {
	Name  []byte
	Taxon []byte
}

type taxnode struct {
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

type taxTree map[int]*taxnode
type auxTree map[int]*auxNode

// Taxonomy is the internal representation of the NCBI taxonomy database.
// It includes high level structures to search efficiently for LCAs (in constant time)
type Taxonomy struct {
	// TODO: Unexport all the fields that don't require to be exported
	T       taxTree
	G       giTaxid.GiMapper
	D       map[int]int // from values to indexes
	E, L, H []int
	M       [][]int
}

func (n *pathnode) String() string {
	return fmt.Sprintf("{taxon:%s, name:%s}", n.Taxon, n.Name)
}

func (n *taxnode) String() string {
	retStr := fmt.Sprintf("\n\t{\n\t\tid:%d\n\t\tTaxid:%d\n\t\tParent:%d\n\t\tName:%s\n\t\tTaxon:%s\n\t\tnext:", n.id, n.Taxid, n.Parent, n.Name, n.Taxon)
	for _, v := range n.Childs {
		retStr += fmt.Sprintf("%d,", v)
	}
	retStr += "\n\t}\n"
	return retStr
}

func (t taxTree) String() string {
	retStr := ""
	for k, v := range t {
		retStr += fmt.Sprintf("%d%s\n", k, v)
	}
	return retStr
}

func (n *auxNode) String() string {
	retStr := fmt.Sprintf("\n\t{\n\t\t_id:%d\n\t\tid:%d\n\t\tprev:%d\n\t\ttaxon:%s\n\t\tnext:", n._id, n.id, n.parent, n.taxon)
	for _, v := range n.childs {
		retStr += fmt.Sprintf("%d,", v)
	}
	retStr += "\n\t}\n"
	return retStr
}

func (t auxTree) String() string {
	retStr := ""
	for k, v := range t {
		retStr += fmt.Sprintf("%d%s\n", k, v)
	}
	return retStr
}

func log2 (x int) int {
	return int (math.Log(float64(x)) / math.Log(float64(2)))
}

func makeMatrix (maxNodes int) [][]int {
	dim := maxNodes * 2
	mat := make([][]int, dim)
//	fmt.Printf("Creating matrix of size: %d x %d", maxNodes, dim2)
	for i:=0; i < dim; i++ {
		mat[i] = make([]int, log2(dim)+1)
	}
	return mat
}

func parsename(l []byte) (int, []byte, error) {
	parts := bytes.Split(l, []byte(sep))
	tid, e := strconv.Atoi(string(parts[0]))
	if e != nil {
		return 0, nil, e
	}
	return tid, parts[1], nil
}

func (t taxTree) loadNames (fname string, dict map[int]int) error {
	namesf, eopen := os.OpenFile(fname, os.O_RDONLY, 0644)
	defer namesf.Close()
	if eopen != nil {
		fmt.Fprintf(os.Stderr, "file doesn't exist %s\n", fname)
		return eopen
	}
	b := bufio.NewReader(io.Reader(namesf))
	for {
		line, _,  err := b.ReadLine()
		if err == io.EOF {
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
	for _, v := range t {
		newchilds := make([]int, 0, len(v.childs))
		for _, c := range v.childs {
			newchilds = append(newchilds, corrs[c])
		}
		taxtree[v._id] = &taxnode{
			id:     v._id,
			Taxid:  v.id,
			Parent: corrs[v.parent],
			Childs: newchilds,
			Taxon:  v.taxon,
		}
	}
	return taxtree
}

func newTaxonomy(b *bufio.Reader, maxNodes int) (taxTree, map[int]int, error) {
	auxtree := make(auxTree, maxNodes)
	i := 0
	for {
		l, _, err := b.ReadLine()
		i++
		if err == io.EOF {
			corrs := make(map[int]int, maxNodes)
			var idx int = 1
			auxtree.addIdx(1, &idx, corrs) // HINT: Args : first node in tree, first idx to assign, correspondences between old ids and new ids
			taxtree := auxtree.restoreRels(1, corrs)
			return taxtree, corrs, nil
		}
		if err != nil {
			return nil, nil, err
		}
		line := l[:len(l)-2] // HINT : Ends with \t|
		//		newNode := parsenode([]byte(line)[0 : len(line)-3]) // ends with \t\|\n
		//		taxmap.T[newNode.This.Tids.Taxid] = newNode
		parts := bytes.SplitN(line, []byte(sep), 4)
		this, ae1 := strconv.Atoi(string(parts[0]))
		if ae1 != nil {
			return nil, nil, ae1
		}
		that, ae2 := strconv.Atoi(string(parts[1]))
		if ae2 != nil {
			return nil, nil, ae2
		}
		if this == that {
			continue // To avoid circular references in the tree
		}
		if _, ok := auxtree[this]; ok {
			auxtree[this].parent = that
			auxtree[this].taxon = make([]byte, len(parts[2]))
			copy(auxtree[this].taxon, parts[2])
		} else {
			auxNode := &auxNode{
				id:     this,
				parent: that,
				childs: []int{},
			}
			auxNode.taxon = make([]byte, len(parts[2]))
			copy(auxNode.taxon, parts[2])
			auxtree[this] = auxNode
		}
		if _, ok := auxtree[that]; ok {
			auxtree[that].childs = append(auxtree[that].childs, this)
		} else {
			auxNode := &auxNode{
				id:     that,
				childs: []int{this},
			}
			auxtree[that] = auxNode
		}
	}
	return nil, nil, nil // never used
}


// New creates a new NCBI taxonomy representation from the options given in opts
// returns the newly created taxonomy or any error it may encounter in the process
func New(nodesfn, namesfn, dictfn string, savemem bool) (*Taxonomy, error) {
	// T : taxtree => tax
	t := &Taxonomy{}
	fmt.Fprintf(os.Stderr, "Creating new taxonomy tree ... ")
	s1 := time.Now()
	maxNodes, ewcl := wcl.FromFile(nodesfn)
	if ewcl != nil {
		return nil, ewcl
	}
	nodesf, eopen := os.OpenFile(nodesfn, os.O_RDONLY, 0644)
	if eopen != nil {
		return nil, eopen
	}
	nodesbuf := bufio.NewReader(io.Reader(nodesf))

	tax, dict, ee := newTaxonomy(nodesbuf, maxNodes)
	if ee != nil {
		return nil, ee
	}
	s2 := time.Now()
	dur := s2.Sub(s1)
	fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", dur.Seconds())

	// T.Name(s) 
	fmt.Fprintf(os.Stderr, "Filling names in taxonomy tree ... ")
	s1 = time.Now()
	err := tax.loadNames(namesfn, dict)
	if err != nil {
		return nil, err
	}
	s2 = time.Now()
	dur = s2.Sub(s1)
	fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", dur.Seconds())

	// E, L, H
	fmt.Fprintf(os.Stderr, "Creating indexes ... ")
	s1 = time.Now()
	E := make([]int, 0, maxNodes*2)
	L := make([]int, 0, maxNodes*2)
	H := make([]int, maxNodes)
	tax.elh(1, 0, &E, &L, &H)
	s2 = time.Now()
	dur = s2.Sub(s1)
	fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", dur.Seconds())

	// M
	fmt.Fprintf(os.Stderr, "Preprocessing RMQ ... ")
	s1 = time.Now()
	M := makeMatrix(maxNodes)
	rmqPrep(&M, L[:len(L)-1], maxNodes * 2)
	s2 = time.Now()
	dur = s2.Sub(s1)
	fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", dur.Seconds())

	fmt.Fprintf(os.Stderr, "Combining data ... ")
	s1 = time.Now()
	t.T = tax
	t.D = dict
	//		t.E = make([]int, len(E))
	//		copy(t.E, E)
	t.E = E
	//		t.L = make([]int, len(L))
	//		copy(t.L, L)
	t.L = L
	//		t.H = make([]int, len(H))
	//		copy(t.H, H)
	t.H = H
	// 		t.M = make([][]int, len(M))
	// 		for i := 0; i < len(M); i++ {
	// 			t.M[i] = make([]int, len(M[0]))
	// 		}
	// 		copy(t.M, M)
	t.M = M
	s2 = time.Now()
	dur = s2.Sub(s1)
	fmt.Fprintf(os.Stderr, "Done (%.3f sec)\n", dur.Seconds())

	t.G, err = giTaxid.Load(dictfn, savemem)
	if err != nil {
		return nil, err
	}
	return t, nil
}

//TODO: Unexport this function?
func (t Taxonomy) Node(taxid int) *taxnode {
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

//TODO: Unexport this function?
func (t Taxonomy) Path(taxid int) []*pathnode {
	node := t.Node(taxid)
	path := make([]*pathnode, 0, 10)
	for {
		if node.Taxid == 1 {
			return path
		}
		newNode := &pathnode{}
		newNode.Name = make([]byte, len(node.Name))
		newNode.Taxon = make([]byte, len(node.Taxon))
		copy(newNode.Name, node.Name)
		copy(newNode.Taxon, node.Taxon)
		path = append(path, newNode)
		node = t.Parent(node)
	}
	return nil
}

//TODO: Unexport this function?
func (t *Taxonomy) PathFromGi(gi int) ([]*pathnode, error) {
	taxid, err := t.TaxidFromGi(gi)
	if err != nil {
		return nil, err
	}
	return t.Path(taxid), nil
}

// TaxidFromGi returns the Taxid associated with a given GI
func (t *Taxonomy) TaxidFromGi(gi int) (int, error) {
	taxid, err := t.G.GiTaxid(gi)
	if err != nil {
		return -1, err
	}
	return taxid, nil
}

//TODO: Unexport this function?
func (t *Taxonomy) Parent(node *taxnode) *taxnode {
	return t.T[node.Parent]
}

// AtLevels returns a slice of slices having the taxons at the specified
// taxonomic levels
func (t *Taxonomy) AtLevels(node *taxnode, levs ...[]byte) [][]byte {
	taxAtLevels := make([][]byte, 0, len(levs))
	//	fmt.Println("NODE : ", node)
	//	os.Exit(0)
	allLevs := t.AllLevels(node)                    // HINT: Taxonomy levels to names
	baseLevN, _ := sortedLevels[string(node.Taxon)] // HINT: Levels below this are "uc_"
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

		taxAtLevels = append(taxAtLevels, Unknown)
	}
	return taxAtLevels
}

func (t *Taxonomy) AllLevels(node *taxnode) map[string][]byte {
	taxons := make(map[string][]byte, 15)
	for {
		if node.Taxid == 1 {
			return taxons
		}
		if bytes.Compare(node.Taxon, no_rank) != 0 {
			taxons[string(node.Taxon)] = node.Name
		}
		node = t.Parent(node)
	}
	return nil // Should be never used
}

func (t *Taxonomy) AtLevel(node *taxnode, lev []byte) []byte {
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


