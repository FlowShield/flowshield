package verify

import "sync/atomic"

type Record struct {
	Total   uint64 `json:"total"`
	Success uint64 `json:"success"`
	Fail    uint64 `json:"fail"`
}

func (r *Record) AddSuccess(num uint64) {
	atomic.AddUint64(&r.Fail, num)
	atomic.AddUint64(&r.Total, num)
}

func (r *Record) AddFail(num uint64) {
	atomic.AddUint64(&r.Fail, num)
	atomic.AddUint64(&r.Total, num)
}

func (r *Record) GetFail() uint64 {
	return atomic.LoadUint64(&r.Fail)
}

func (r *Record) GetSuccess() uint64 {
	return atomic.LoadUint64(&r.Success)
}

func (r *Record) GetTotal() uint64 {
	return atomic.LoadUint64(&r.Total)
}
