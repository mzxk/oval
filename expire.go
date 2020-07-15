/*
ExpireMap提供一个自动过期的sync.map对象
这样可以避免内存泄漏的问题
没有提供遍历的函数，同时增加了一个重置过期时间的函数 Expire
*/
package oval

import (
	"sync"
	"time"
)

//ExpireMap 主要结构
//m 主要的map
//c 在每次写入时会写入这个chan，来处理过期问题
//i sync.map的长度
type ExpireMap struct {
	m sync.Map
	c chan expireChanStru
}

//保存在sync.map里的结构
type expireMapStru struct {
	Data interface{}
	Tm   int64
}

//保存在过期chan里的结构
type expireChanStru struct {
	Key interface{}
	Tm  int64
}

//这是一个循环 ， 首先从过期时间chan里读取最近的一条
//如果这一条的时间已经过期，那么删除此条对应的map并立刻处理下一条，如果这条没有过期，那么等待一秒后重新处理
func (t *ExpireMap) loop() {
	for {
		d := <-t.c
		//can delete
	RE:
		if time.Now().Unix() >= d.Tm {
			if date, ok := t.m.Load(d.Key); ok {
				if d.Tm == date.(*expireMapStru).Tm {
					t.Delete(d.Key)
				}
			}
		} else {
			time.Sleep(500 * time.Millisecond)
			goto RE
		}

	}
}

//Load 同sync.map.Load()
func (t *ExpireMap) Load(k interface{}) (interface{}, bool) {
	i, ok := t.m.Load(k)
	if ok {
		return i.(*expireMapStru).Data, ok
	}
	return nil, ok
}

//Store 同 sync.map.store(),但是需要输入一个过期时间（秒）
func (t *ExpireMap) Store(k, v interface{}, expire int64) {
	ex := time.Now().Unix() + expire
	vv := &expireMapStru{v, ex} //保存map
	t.m.Store(k, vv)
	t.c <- expireChanStru{k, ex} //保存过期chan
}

//LoadOrStore 同 sync.map.LoadOrStore(),但是需要输入一个过期时间（秒）
func (t *ExpireMap) LoadOrStore(k, v interface{}, expire int64) (interface{}, bool) {
	ex := time.Now().Unix() + expire
	vv := &expireMapStru{v, ex} //保存map
	i, loaded := t.m.LoadOrStore(k, vv)
	if loaded {
		return i.(*expireMapStru).Data, loaded
	}
	t.c <- expireChanStru{k, ex} //保存过期chan
	return nil, false
}

//Delete j同 sync.map.Delete()
func (t *ExpireMap) Delete(k interface{}) {
	t.m.Delete(k)
}

//Expire 这个函数会修改当前map的过期时间
func (t *ExpireMap) Expire(k string, expire int64) {
	if v, ok := t.m.Load(k); ok {
		ex := time.Now().Unix() + expire
		if vv, ok := v.(*expireMapStru); ok {
			vv.Tm = ex
			t.c <- expireChanStru{k, +ex}
		}
	}
}

//NewExpire 这个将返回一个带自动过期时间带sync。map
func NewExpire() *ExpireMap {
	t := &ExpireMap{
		c: make(chan expireChanStru, 99999999),
	}
	go t.loop()
	return t
}
