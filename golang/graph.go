package golang

type RefEdge struct {
	FQN string // chunk FQN
}

type RefGraph struct {
	// RefGraph represents FQN -> []FQN for *all* references of the chunk
	Edges map[string][]RefEdge
}

func NewRefGraph(chunks []Chunk) *RefGraph {
	g := &RefGraph{
		Edges: make(map[string][]RefEdge),
	}

	for _, chunk := range chunks {
		for _, ref := range chunk.References {
			g.Edges[chunk.FQN()] = append(g.Edges[chunk.FQN()], RefEdge{FQN: ref})
		}
	}

	return g
}
