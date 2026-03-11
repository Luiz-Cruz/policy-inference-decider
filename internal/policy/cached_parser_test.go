package policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDOT = `digraph { start [result=""]; ok [result="approved=true"]; no [result="approved=false"]; start -> ok [cond="age>=18"]; start -> no [cond="age<18"]; }`

func newTestCachedParser(capacity int) *CachedParser {
	return NewCachedParser(NewDotParser(), NewLFUGraphCache(capacity))
}

func TestCachedParser(t *testing.T) {
	t.Run("first call parses and caches", func(t *testing.T) {
		parser := newTestCachedParser(10)

		graph, err := parser.Parse(context.Background(), testDOT)

		require.NoError(t, err)
		assert.NotNil(t, graph)
		assert.Equal(t, "start", graph.Start)
	})

	t.Run("second call returns same graph from cache", func(t *testing.T) {
		parser := newTestCachedParser(10)

		graph1, err := parser.Parse(context.Background(), testDOT)
		require.NoError(t, err)

		graph2, err := parser.Parse(context.Background(), testDOT)
		require.NoError(t, err)

		assert.Same(t, graph1, graph2)
	})

	t.Run("different DOT strings produce different graphs", func(t *testing.T) {
		parser := newTestCachedParser(10)
		otherDOT := `digraph { start [result=""]; end [result="x=1"]; start -> end [cond="x==1"]; }`

		graph1, err := parser.Parse(context.Background(), testDOT)
		require.NoError(t, err)

		graph2, err := parser.Parse(context.Background(), otherDOT)
		require.NoError(t, err)

		assert.NotSame(t, graph1, graph2)
	})

	t.Run("LFU keeps most frequently used when full", func(t *testing.T) {
		parser := newTestCachedParser(2)
		dot1 := `digraph { start [result=""]; a [result="x=1"]; start -> a [cond="x==1"]; }`
		dot2 := `digraph { start [result=""]; b [result="x=2"]; start -> b [cond="x==2"]; }`
		dot3 := `digraph { start [result=""]; c [result="x=3"]; start -> c [cond="x==3"]; }`

		_, _ = parser.Parse(context.Background(), dot1)
		_, _ = parser.Parse(context.Background(), dot1)
		_, _ = parser.Parse(context.Background(), dot1)
		_, _ = parser.Parse(context.Background(), dot2)

		_, _ = parser.Parse(context.Background(), dot3)

		graph1, err := parser.Parse(context.Background(), dot1)
		require.NoError(t, err)
		assert.NotNil(t, graph1)
	})

	t.Run("invalid DOT is not cached", func(t *testing.T) {
		parser := newTestCachedParser(10)
		invalidDOT := `digraph { start [result=]; }`

		_, err := parser.Parse(context.Background(), invalidDOT)

		assert.Error(t, err)
	})

	t.Run("DOT without start node is not cached", func(t *testing.T) {
		parser := newTestCachedParser(10)
		noStartDOT := `digraph { foo [result=""]; bar [result="x=1"]; foo -> bar [cond="true"]; }`

		_, err := parser.Parse(context.Background(), noStartDOT)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNoStartNode)
	})
}

func TestCachedParserWithMockCache(t *testing.T) {
	t.Run("uses injected cache implementation", func(t *testing.T) {
		mock := &mockGraphCache{store: make(map[string]*Graph)}
		parser := NewCachedParser(NewDotParser(), mock)

		_, err := parser.Parse(context.Background(), testDOT)
		require.NoError(t, err)
		assert.Equal(t, 1, mock.setCount)

		_, err = parser.Parse(context.Background(), testDOT)
		require.NoError(t, err)
		assert.Equal(t, 1, mock.getHitCount)
		assert.Equal(t, 1, mock.setCount)
	})
}

func TestPanicsOnInvalidConfiguration(t *testing.T) {
	t.Run("panics when inner parser is nil", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidInnerParser, func() {
			NewCachedParser(nil, NewLFUGraphCache(10))
		})
	})

	t.Run("panics when cache is nil", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidCache, func() {
			NewCachedParser(NewDotParser(), nil)
		})
	})

	t.Run("panics when LFU capacity is zero", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidCacheCapacity, func() {
			NewLFUGraphCache(0)
		})
	})

	t.Run("panics when LFU capacity is negative", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidCacheCapacity, func() {
			NewLFUGraphCache(-1)
		})
	})
}

func TestHashDOT(t *testing.T) {
	t.Run("same input produces same hash", func(t *testing.T) {
		h1 := hashDOT("digraph { start -> ok }")
		h2 := hashDOT("digraph { start -> ok }")

		assert.Equal(t, h1, h2)
	})

	t.Run("different input produces different hash", func(t *testing.T) {
		h1 := hashDOT("digraph { start -> ok }")
		h2 := hashDOT("digraph { start -> no }")

		assert.NotEqual(t, h1, h2)
	})
}

type mockGraphCache struct {
	store       map[string]*Graph
	getHitCount int
	setCount    int
}

func (m *mockGraphCache) Get(key string) (*Graph, bool) {
	graph, ok := m.store[key]
	if ok {
		m.getHitCount++
	}
	return graph, ok
}

func (m *mockGraphCache) Set(key string, value *Graph) {
	m.store[key] = value
	m.setCount++
}
