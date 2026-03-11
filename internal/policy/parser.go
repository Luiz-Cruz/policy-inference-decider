package policy

import (
	"context"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/awalterschulze/gographviz/ast"
	"github.com/casbin/govaluate"
)

type DotParser struct{}

func NewDotParser() *DotParser {
	return &DotParser{}
}

func (DotParser) Parse(ctx context.Context, dot string) (*Graph, error) {
	astGraph, err := gographviz.ParseString(dot)
	if err != nil {
		return nil, ErrInvalidPolicyDot
	}

	nodes, edges := buildGraphFromAST(astGraph)
	if err = validateHasStart(nodes); err != nil {
		return nil, err
	}
	adjList := buildAdjList(edges)
	return &Graph{Nodes: nodes, Edges: edges, AdjList: adjList, Start: StartNodeID}, nil
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

func buildAdjList(edges []*Edge) map[string][]*Edge {
	adj := make(map[string][]*Edge)
	for _, e := range edges {
		adj[e.From] = append(adj[e.From], e)
	}
	return adj
}

func nodeFromStmt(stmt *ast.NodeStmt) *Node {
	id := string(stmt.NodeID.ID)
	result := extractResultFromNodeAttrs(stmt.Attrs)
	return &Node{ID: id, Result: result, ParsedResult: preParseResult(result)}
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

	valid := cond == "" || isValidCond(cond)
	var compiled *govaluate.EvaluableExpression
	if valid && cond != "" {
		compiled, _ = govaluate.NewEvaluableExpression(cond)
	}

	return &Edge{From: from, To: to, Cond: cond, ValidCond: valid, CompiledCond: compiled}, true
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
