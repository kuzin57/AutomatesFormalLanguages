package display

import (
	"fmt"
	"workspace/internal/entities/automate"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func DisplayGraph(states []*automate.State, name string) (err error) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return
	}

	defer func() {
		if err = graph.Close(); err != nil {
			return
		}
		g.Close()
	}()

	if err = createGraph(states, graph); err != nil {
		return
	}

	filename := "./" + name + ".png"
	if err = g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		return
	}

	return
}

func createGraph(states []*automate.State, graph *cgraph.Graph) (err error) {
	nodes := make([]*cgraph.Node, len(states))
	for i, state := range states {
		nodes[i], err = graph.CreateNode(fmt.Sprint(state.Number))
		if err != nil {
			return err
		}

		if state.IsTerminal {
			nodes[i].SetFontColor("blue")
		}
	}

	return createEdges(states, graph, nodes)
}

func createEdges(states []*automate.State, graph *cgraph.Graph, nodes []*cgraph.Node) error {
	mapNodes := make(map[int]*cgraph.Node)
	for i, node := range nodes {
		mapNodes[states[i].Number] = node
	}

	for i, state := range states {
		for key, to := range state.Transitions {
			for _, v := range to {
				edge, err := graph.CreateEdge(string(key), nodes[i], mapNodes[v])
				if err != nil {
					return err
				}

				edge.SetLabel(string(key))
			}
		}
	}
	return nil
}
