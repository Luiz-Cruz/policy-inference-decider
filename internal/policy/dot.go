package policy

import (
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/awalterschulze/gographviz/ast"
)

func ParseDOT(dot string) (*Graph, error) {
	astGraph, err := gographviz.ParseString(dot)
	if err != nil {
		return nil, err
	}

	nodes, edges := buildGraphFromAST(astGraph)
	if err = validateHasStart(nodes); err != nil {
		return nil, err
	}
	return &Graph{Nodes: nodes, Edges: edges, Start: "start"}, nil
}

func buildGraphFromAST(astGraph *ast.Graph) (map[string]*Node, []*Edge) {
	nodes := make(map[string]*Node)
	var edges []*Edge
	for _, stmt := range astGraph.StmtList {
		if nodeStmt, ok := stmt.(*ast.NodeStmt); ok {
			n := nodeFromStmt(nodeStmt)
			nodes[n.ID] = n
		}
		if edgeStmt, ok := stmt.(*ast.EdgeStmt); ok {
			if e, ok := edgeFromStmt(edgeStmt); ok {
				edges = append(edges, e)
			}
		}
	}
	return nodes, edges
}

func nodeFromStmt(stmt *ast.NodeStmt) *Node {
	id := string(stmt.NodeID.ID)
	result := extractResultFromNodeAttrs(stmt.Attrs)
	return &Node{ID: id, Result: result}
}

func extractResultFromNodeAttrs(attrs ast.AttrList) string {
	for _, a := range attrs {
		s := a.String()
		if strings.HasPrefix(s, "result=") {
			return strings.Trim(strings.TrimPrefix(s, "result="), "\"")
		}
	}
	return ""
}

func edgeFromStmt(stmt *ast.EdgeStmt) (*Edge, bool) {
	if len(stmt.EdgeRHS) == 0 {
		return nil, false
	}
	from := string(stmt.Source.GetID())
	to := string(stmt.EdgeRHS[0].Destination.GetID())
	cond := extractCondFromEdgeAttrs(stmt.Attrs)
	return &Edge{From: from, To: to, Cond: cond}, true
}

func extractCondFromEdgeAttrs(attrs ast.AttrList) string {
	for _, attrList := range attrs {
		for _, a := range attrList {
			s := a.String()
			if strings.HasPrefix(s, "cond=") {
				return strings.Trim(strings.TrimPrefix(s, "cond="), "\"")
			}
		}
	}
	return ""
}

func validateHasStart(nodes map[string]*Node) error {
	if _, hasStart := nodes["start"]; !hasStart {
		return ErrNoStartNode
	}
	return nil
}
