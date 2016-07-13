package taxonomy

import(
	"math"
	"errors"
)

func (t taxTree) elh(node, level int, E, L, H *[]int) {
	n, _ := t[node]
	*E = append(*E, n.id)
	*L = append(*L, level)
	if (*H)[n.id-1] == 0 {
		(*H)[n.id-1] = len(*E)
	}
	for _, nextNode := range n.Childs {
		t.elh(nextNode, level+1, E, L, H)
	}
	*E = append(*E, n.Parent)
	*L = append(*L, level-1)
}

func rmqPrep(M *[][]int, A []int, N int) {
	for i := 0; i < N; i++ {
		(*M)[i][0] = i
	}
	for j := 1; 1<<uint(j) <= N; j++ {
		for i := 0; i+(1<<uint(j))-1 < N; i++ {
			if A[(*M)[i][j-1]] < A[(*M)[i+(1<<(uint(j)-1))][j-1]] {
				(*M)[i][j] = (*M)[i][j-1]
			} else {
				(*M)[i][j] = (*M)[i+(1<<(uint(j)-1))][j-1]
			}
		}
	}
}

func rmq(M [][]int, A []int, i, j int) (rmq int) {
	k := int(log2(j - i + 1))
	//	k := int(math.Log(float64(j - i + 1)))
	if A[M[i][k]] <= A[M[j-int(math.Pow(float64(2), float64(k)))+1][k]] {
		rmq = M[i][k]
	} else {
		rmq = M[j-int(math.Pow(float64(2), float64(k)))+1][k]
	}
	return
}


// LCA calculates the lowest common ancestor of a list of taxon ids
func (t Taxonomy) LCA(values ...int) (*taxnode, error) {
	indexes := make([]int, 0, len(values))
	for _, v := range values { // from values to indexes
		if _, ok := t.D[v]; ok { // HINT -- There may be taxids not in taxonomy
			indexes = append(indexes, t.D[v])
		}
	}
	if len(indexes) == 0 {
		return &taxnode{}, errors.New("EMPTY")
	}
	red := indexes[0]
	indexes = indexes[1:]
	for _, ind := range indexes {
		//fmt.Fprintf(os.Stderr, "LCA so far:%v VS %v\n", red, ind)
		red = lcaHelper(t.E, t.L, t.H, t.M, red, ind)
	}
	return t.T[red], nil
}

func lcaHelper(E, L, H []int, M [][]int, i, j int) int {
	if i == j {
		return i
	}
	v1 := H[i-1]
	v2 := H[j-1]
	if v1 > v2 {
		v1, v2 = v2, v1
	}
	rmq := rmq(M, L, v1, v2)
	lca := E[rmq]
	return lca
}
