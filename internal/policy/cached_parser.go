package policy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
)

const (
	DefaultCacheCapacity = 100
	invalidInnerParser   = "inner parser must not be nil"
	invalidCache         = "cache must not be nil"
)

type (
	GraphCache interface {
		Get(key string) (*Graph, bool)
		Set(key string, value *Graph)
	}

	CachedParser struct {
		inner Parser
		cache GraphCache
	}
)

func NewCachedParser(inner Parser, cache GraphCache) *CachedParser {
	if inner == nil {
		panic(invalidInnerParser)
	}
	if cache == nil {
		panic(invalidCache)
	}
	return &CachedParser{inner: inner, cache: cache}
}

func (p *CachedParser) Parse(ctx context.Context, dot string) (*Graph, error) {
	key := hashDOT(dot)
	if graph, ok := p.cache.Get(key); ok {
		slog.InfoContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:cache_hit] [dot_hash:%s]", key))
		return graph, nil
	}

	graph, err := p.inner.Parse(ctx, dot)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:cache_miss_parse_error] [dot_hash:%s] [err:%+v]", key, err))
		return nil, err
	}

	p.cache.Set(key, graph)
	slog.InfoContext(ctx, fmt.Sprintf("[feature:policy_inference] [msg:cache_miss_stored] [dot_hash:%s]", key))
	return graph, nil
}

func hashDOT(dot string) string {
	h := sha256.Sum256([]byte(dot))
	return hex.EncodeToString(h[:])
}
