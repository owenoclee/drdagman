package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

type Node struct {
	ID             string `yaml:"id"`
	Implementation string `yaml:"implementation"`
}

type Transition struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Dag struct {
	Nodes       []Node       `yaml:"nodes"`
	Transitions []Transition `yaml:"transitions"`
}

type DagGraph struct {
	Nodes        map[string]*Node
	Nexts        map[string][]*Node
	Dependencies map[string][]*Node
}

func buildDagGraph(d *Dag) *DagGraph {
	graph := &DagGraph{
		Nodes:        make(map[string]*Node),
		Nexts:        make(map[string][]*Node),
		Dependencies: make(map[string][]*Node),
	}
	for i := range d.Nodes {
		n := &d.Nodes[i]
		graph.Nodes[n.ID] = n
		graph.Dependencies[n.ID] = []*Node{}
	}
	for _, t := range d.Transitions {
		parent, okFrom := graph.Nodes[t.From]
		child, okChild := graph.Nodes[t.To]
		if !okFrom || !okChild {
			log.Printf("warning: transition from %s to %s contains unknown nodes", t.From, t.To)
			continue
		}
		graph.Nexts[t.From] = append(graph.Nexts[t.From], child)
		graph.Dependencies[t.To] = append(graph.Dependencies[t.To], parent)
	}
	return graph
}

func (g *DagGraph) TopologicalSort() []*Node {
	visited := make(map[string]bool)
	var sorted []*Node
	var visit func(n *Node)
	visit = func(n *Node) {
		if visited[n.ID] {
			return
		}
		visited[n.ID] = true
		for _, dep := range g.Dependencies[n.ID] {
			visit(dep)
		}
		sorted = append(sorted, n)
	}
	for _, n := range g.Nodes {
		visit(n)
	}
	return sorted
}

func (g *DagGraph) Root() (*Node, error) {
	roots := []*Node{}
	for _, n := range g.Nodes {
		if len(g.Dependencies[n.ID]) == 0 {
			roots = append(roots, n)
		}
	}
	if len(roots) != 1 {
		return nil, fmt.Errorf("expected exactly one root, got %d", len(roots))
	}
	return roots[0], nil
}

func (g *DagGraph) Leaf() (*Node, error) {
	leaves := []*Node{}
	for _, n := range g.Nodes {
		if len(g.Nexts[n.ID]) == 0 {
			leaves = append(leaves, n)
		}
	}
	if len(leaves) != 1 {
		return nil, fmt.Errorf("expected exactly one leaf, got %d", len(leaves))
	}
	return leaves[0], nil
}

func main() {
	data, err := os.ReadFile("dag.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var d Dag
	if err := yaml.Unmarshal(data, &d); err != nil {
		log.Fatal(err)
	}
	graph := buildDagGraph(&d)

	// debug printing
	func() {
		fmt.Println("nodes:")
		for id, node := range graph.Nodes {
			fmt.Printf("  %s: { id: %s, implementation: %s }\n", id, node.ID, node.Implementation)
		}

		fmt.Println("nexts:")
		for id, children := range graph.Nexts {
			fmt.Printf("  %s: [", id)
			for i, child := range children {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", child.ID)
			}
			fmt.Println("]")
		}

		fmt.Println("dependencies:")
		for id, deps := range graph.Dependencies {
			fmt.Printf("  %s: [", id)
			for i, dep := range deps {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", dep.ID)
			}
			fmt.Println("]")
		}
	}()

	root, err := graph.Root()
	if err != nil {
		log.Fatal(err)
	}
	leaf, err := graph.Leaf()
	if err != nil {
		log.Fatal(err)
	}

	atoi := func(s string) int {
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("expected integer, got %s", s)
		}
		return n
	}
	atof := func(s string) float64 {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatalf("expected float, got %s", s)
		}
		return n
	}
	execute := func(n *Node, inputs []int) (int, error) {
		actionAndArgs := strings.Fields(n.Implementation)
		action := actionAndArgs[0]
		args := actionAndArgs[1:]
		fmt.Printf("executing node %s (impl: %s) with inputs %v\n", n.ID, n.Implementation, inputs)
		switch action {
		case "add":
			return inputs[0] + atoi(args[0]), nil
		case "multiply":
			return int(float64(inputs[0]) * atof(args[0])), nil
		case "sum":
			sum := 0
			for _, input := range inputs {
				sum += input
			}
			return sum, nil
		}
		return 0, fmt.Errorf("unknown action %s", action)
	}

	startingValue := 13
	outputForNode := make(map[string]int)
	for _, n := range graph.TopologicalSort() {
		output, err := func() (int, error) {
			if n == root {
				return execute(n, []int{startingValue})
			}
			inputs := make([]int, 0, len(graph.Dependencies[n.ID]))
			for _, dep := range graph.Dependencies[n.ID] {
				input, ok := outputForNode[dep.ID]
				if !ok {
					return 0, fmt.Errorf("missing output for node %s", dep.ID)
				}
				inputs = append(inputs, input)
			}
			return execute(n, inputs)
		}()
		if err != nil {
			log.Fatalf("error executing node %s: %v", n.ID, err)
		}
		outputForNode[n.ID] = output
	}
	fmt.Println(outputForNode[leaf.ID])
}
