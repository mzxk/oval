package oval

import (
	"sync"
	"sync/atomic"
	"time"
)

type ExpireMap struct {
	m sync.Map
	c chan expChan
	i int64
}
type expStr struct {
	Data interface{}
	Tm   int64
}
type expChan struct {
	Key interface{}
	Tm  int64
}

func (t *ExpireMap) loop() {
	for {
		d := <-t.c
		//can delete
	RE:
		if time.Now().Unix() >= d.Tm {
			if date, ok := t.m.Load(d.Key); ok {
				if d.Tm == date.(*expStr).Tm {
					t.Delete(d.Key)
				}
			}
		} else {
			time.Sleep(1 * time.Second)
			goto RE
		}

	}
}

//Load justlike sync.Load
func (t *ExpireMap) Load(k interface{}) (interface{}, bool) {
	i, ok := t.m.Load(k)
	if ok {
		return i.(*expStr).Data, ok
	}
	return nil, ok
}

//Length justlike sync.Load
func (t *ExpireMap) Length() int64 {
	return t.i
}

//Store justlike sync.Store
func (t *ExpireMap) Store(k, v interface{}, expire int64) {
	ex := time.Now().Unix() + expire
	vv := &expStr{v, ex}
	t.m.Store(k, vv)
	t.c <- expChan{k, ex}
	atomic.AddInt64(&t.i, 1)
}

//Delete justlike sync.Delete
func (t *ExpireMap) Delete(k interface{}) {
	t.m.Delete(k)
	atomic.AddInt64(&t.i, -1)
}

//Expire justlike sync.Delete
func (t *ExpireMap) Expire(k string, expire int64) {
	if v, ok := t.m.Load(k); ok {
		ex := time.Now().Unix() + expire
		if vv, ok := v.(*expStr); ok {
			vv.Tm = ex
			t.c <- expChan{k, +ex}
		}
	}
}

func NewExpire() *ExpireMap {
	t := &ExpireMap{
		c: make(chan expChan, 99999999),
	}
	go t.loop()
	return t
}
