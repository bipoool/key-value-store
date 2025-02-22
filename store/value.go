package store

import (
	"time"
)

type Value struct {
	lastAccess int64
	isExpired  bool
	ttl        uint64
	value      any
}

func (v *Value) CheckIsExpired() bool {
	if v.isExpired {
		return v.isExpired
	}
	if v.ttl == 0 {
		return v.isExpired
	}
	currTime := time.Now().Unix()
	if v.ttl <= uint64(currTime) {
		v.isExpired = true
	}
	return v.isExpired
}

func (v *Value) SetIsExpired() {
	v.isExpired = true
}

func (v *Value) GetValue() any {
	return v.value
}

func (v *Value) SetValue(value any) {
	v.value = value
}

func (v *Value) SetTtl(ttl uint64) {
	v.ttl = ttl
}

func (v *Value) SetLastAccess() {
	v.lastAccess = time.Now().Unix()
}

func (v *Value) GetLastAccess() int64 {
	return v.lastAccess
}

func (v *Value) GetTtl() uint64 {
	return v.ttl
}

func NewValue() *Value {
	return &Value{
		lastAccess: time.Now().Unix(),
		isExpired:  false,
	}
}
