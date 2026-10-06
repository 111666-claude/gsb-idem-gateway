// Package idemgateway 按幂等键去重并缓存结果。
package idemgateway

// Counter 记录台账比较次数。
type Counter struct {
	Scanned int
}

// Record 是一条幂等记录。
type Record struct {
	Key     string
	State   string
	Value   string
	Expires int
}

// State 是幂等台账。
type State struct {
	TTL  int
	Recs []Record
}

// NewState 建一份台账。
func NewState(ttl int) *State {
	return &State{TTL: ttl}
}

func find(st *State, key string, c *Counter) *Record {
	c.Scanned++
	for i := range st.Recs {
		c.Scanned++
		if st.Recs[i].Key == key {
			return &st.Recs[i]
		}
	}
	return nil
}

// Begin 申请处理一个幂等键。
func Begin(st *State, key string, now int, c *Counter) string {
	rec := find(st, key, c)
	if rec != nil && rec.Value != "" {
		return "replay"
	}
	st.Recs = append(st.Recs, Record{Key: key, State: "pending", Expires: now + st.TTL})
	return "new"
}

// Complete 记录一次处理结果。
func Complete(st *State, key, value string, now int, c *Counter) {
	rec := find(st, key, c)
	if rec == nil {
		return
	}
	rec.State = "done"
	rec.Value = value
}

// Result 读取已缓存的结果。
func Result(st *State, key string, c *Counter) string {
	rec := find(st, key, c)
	if rec == nil {
		return ""
	}
	return rec.Value
}
