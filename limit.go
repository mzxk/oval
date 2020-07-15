/*
只适合不重要的限制，如api接口等，这将直接使用内存，并且不会落盘
*/
package oval

import (
	"sync"
	"sync/atomic"
)

//只启动一次limit
var onceLimit sync.Once
var expLimit *ExpireMap

//Limited 输入唯一key，单位时间,单位时间内限制的次数，
//当返回false时，代表没有限制，当返回true时，代表已经限制
func Limited(key string, second, times int64) bool {
	onceLimit.Do(func() {
		expLimit = NewExpire()
	})
	var tm int64 = 1
	i, loaded := expLimit.LoadOrStore(key, &tm, second)
	if !loaded {
		return false
	}
	if ii := atomic.AddInt64(i.(*int64), 1); ii > times {
		return true
	}
	return false
}

//UnLimited 取消限制
func UnLimited(key string) {
	onceLimit.Do(func() {
		expLimit = NewExpire()
	})
	expLimit.Delete(key)
}
