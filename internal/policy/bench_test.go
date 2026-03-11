package policy

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
	"strconv"
	"testing"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})))
	if s := os.Getenv("BENCH_MEMORY_LIMIT_MB"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			debug.SetMemoryLimit(int64(n) * 1024 * 1024)
		}
	}
}

const heavyDOT = `digraph CreditPF { start [result=""]; tier1_check [result=""]; tier2_check [result=""]; approved_prime [result="approved=true,segment=prime,limit=50000"]; approved_standard [result="approved=true,segment=standard,limit=20000"]; approved_basic [result="approved=true,segment=basic,limit=5000"]; review_manual [result="approved=false,segment=manual_review,reason=inconclusive"]; review_docs [result="approved=false,segment=doc_review,reason=missing_docs"]; rejected_age [result="approved=false,reason=underage"]; rejected_score [result="approved=false,reason=low_score"]; rejected_income [result="approved=false,reason=insufficient_income"]; start -> rejected_age [cond="age<18"]; start -> tier1_check [cond="age>=18 && age<=25"]; start -> tier2_check [cond="age>25"]; tier1_check -> approved_basic [cond="score>=600 && income>=2000"]; tier1_check -> review_docs [cond="score>=400 && income>=2000"]; tier1_check -> rejected_score [cond="score<400"]; tier2_check -> approved_prime [cond="score>=750 && income>=10000"]; tier2_check -> approved_standard [cond="score>=600 && income>=5000"]; tier2_check -> approved_basic [cond="score>=600 && income>=2000"]; tier2_check -> review_manual [cond="score>=400 && income>=2000"]; tier2_check -> rejected_income [cond="income<2000"]; tier2_check -> rejected_score [cond="score<400"]; }`

var benchInput = map[string]any{"age": 22, "score": 650, "income": 5000}

func BenchmarkParse_NoCache(b *testing.B) {
	parser := NewDotParser()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(ctx, heavyDOT)
	}
}

func BenchmarkParse_CacheHit(b *testing.B) {
	parser := NewCachedParser(NewDotParser(), NewLFUGraphCache(DefaultCacheCapacity))
	ctx := context.Background()
	_, _ = parser.Parse(ctx, heavyDOT)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(ctx, heavyDOT)
	}
}

func BenchmarkParse_CacheMiss(b *testing.B) {
	parser := NewCachedParser(NewDotParser(), NewLFUGraphCache(DefaultCacheCapacity))
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dot := fmt.Sprintf(`digraph { start [result=""]; n%d [result="x=1"]; start -> n%d [cond="x==1"]; }`, i, i)
		_, _ = parser.Parse(ctx, dot)
	}
}

func BenchmarkInfer_CacheHit(b *testing.B) {
	parser := NewCachedParser(NewDotParser(), NewLFUGraphCache(DefaultCacheCapacity))
	executor := NewGraphExecutor()
	ctx := context.Background()
	_, _ = parser.Parse(ctx, heavyDOT)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph, _ := parser.Parse(ctx, heavyDOT)
		_, _ = executor.Process(ctx, graph, benchInput)
	}
}

func BenchmarkInfer_CacheMiss(b *testing.B) {
	parser := NewCachedParser(NewDotParser(), NewLFUGraphCache(DefaultCacheCapacity))
	executor := NewGraphExecutor()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dot := fmt.Sprintf(`digraph { start [result=""]; ok [result="approved=true"]; start -> ok [cond="age>=18"]; n%d [result=""]; }`, i)
		graph, _ := parser.Parse(ctx, dot)
		_, _ = executor.Process(ctx, graph, benchInput)
	}
}
