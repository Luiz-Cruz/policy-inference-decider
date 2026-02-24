package policy

import (
	"context"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/awalterschulze/gographviz/ast"
)

type DotParser struct{}

func (DotParser) Parse(ctx context.Context, dot string) (*Graph, error) {
	astGraph, err := gographviz.ParseString(dot)
	if err != nil {
		return nil, ErrInvalidPolicyDot
	}

	nodes, edges := buildGraphFromAST(astGraph)
	if err = validateHasStart(nodes); err != nil {
		return nil, err
	}
	return &Graph{Nodes: nodes, Edges: edges, Start: StartNodeID}, nil
}

func buildGraphFromAST(astGraph *ast.Graph) (map[string]*Node, []*Edge) {
	nodes := make(map[string]*Node)
	var edges []*Edge
	for _, stmt := range astGraph.StmtList {
		if nodeStmt, ok := stmt.(*ast.NodeStmt); ok {
			node := nodeFromStmt(nodeStmt)
			nodes[node.ID] = node
		}
		if edgeStmt, ok := stmt.(*ast.EdgeStmt); ok {
			if edge, ok := edgeFromStmt(edgeStmt); ok {
				edges = append(edges, edge)
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
	if _, hasStart := nodes[StartNodeID]; !hasStart {
		return ErrNoStartNode
	}
	return nil
}
