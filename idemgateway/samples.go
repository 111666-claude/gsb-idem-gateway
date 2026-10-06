package idemgateway

import "fmt"

// Sample 是一个固定场景，用来复现题面里的锚点。
type Sample struct {
	Name string
	Run  func() map[string]any
}

// Samples 返回全部场景。
func Samples() []Sample {
	return []Sample{
		{Name: "inflight", Run: func() map[string]any {
			st := NewState(50)
			c := &Counter{}
			Begin(st, "k", 0, c)
			second := Begin(st, "k", 0, c)
			return map[string]any{"second": second}
		}},
		{Name: "expiry", Run: func() map[string]any {
			st := NewState(50)
			c := &Counter{}
			Begin(st, "k", 0, c)
			Complete(st, "k", "v1", 0, c)
			return map[string]any{"at100": Begin(st, "k", 100, c)}
		}},
		{Name: "empty", Run: func() map[string]any {
			st := NewState(50)
			c := &Counter{}
			Begin(st, "k", 0, c)
			Complete(st, "k", "", 0, c)
			return map[string]any{"second": Begin(st, "k", 0, c)}
		}},
		{Name: "overwrite", Run: func() map[string]any {
			st := NewState(50)
			c := &Counter{}
			Begin(st, "k", 0, c)
			Complete(st, "k", "v1", 0, c)
			Complete(st, "k", "v2", 0, c)
			return map[string]any{"value": Result(st, "k", c)}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	st := NewState(1000000)
	c := &Counter{}
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("k%04d", i)
		Begin(st, key, i, c)
		Complete(st, key, "v", i, c)
	}
	for i := 0; i < n; i++ {
		Begin(st, fmt.Sprintf("k%04d", i), 0, c)
	}
	return map[string]any{"keys": n, "scanned": c.Scanned}
}
