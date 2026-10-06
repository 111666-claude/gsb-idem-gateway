package idemgateway

import "testing"

func TestBeginOnEmptyIsNew(t *testing.T) {
	if got := Begin(NewState(50), "k", 0, &Counter{}); got != "new" {
		t.Fatalf("首次申请应是 new，得到 %q", got)
	}
}

func TestCompleteStoresValue(t *testing.T) {
	st := NewState(50)
	c := &Counter{}
	Begin(st, "k", 0, c)
	Complete(st, "k", "v", 0, c)
	if got := Result(st, "k", c); got != "v" {
		t.Fatalf("结果应是 v，得到 %q", got)
	}
}

func TestResultUnknownIsEmpty(t *testing.T) {
	if got := Result(NewState(50), "k", &Counter{}); got != "" {
		t.Fatalf("未知键结果应为空，得到 %q", got)
	}
}

func TestReplayAfterComplete(t *testing.T) {
	st := NewState(50)
	c := &Counter{}
	Begin(st, "k", 0, c)
	Complete(st, "k", "v", 0, c)
	if got := Begin(st, "k", 0, c); got != "replay" {
		t.Fatalf("完成后再次申请应是 replay，得到 %q", got)
	}
}

func TestTTLIsStored(t *testing.T) {
	if got := NewState(50).TTL; got != 50 {
		t.Fatalf("窗口应是 50，得到 %d", got)
	}
}

func TestCounterCounts(t *testing.T) {
	c := &Counter{}
	Begin(NewState(50), "k", 0, c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}
