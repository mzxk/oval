package oval

import "sync"

type Map struct {
	K sync.Map
	V sync.Map
}

func (t *Map) Store(k, v interface{}) {
	t.K.Store(k, v)
	t.V.Store(v, k)
}
func (t *Map) LoadK(k interface{}) (interface{}, bool) {
	return t.K.Load(k)
}
func (t *Map) LoadV(k interface{}) (interface{}, bool) {
	return t.V.Load(k)
}

//Delete 请保证k和v之间有一个不为nil
func (t *Map) Delete(k, v interface{}) {
	if k == nil {
		k, _ = t.V.Load(v)
	}
	if v == nil {
		v, _ = t.K.Load(k)
	}
	t.K.Delete(k)
	t.V.Delete(v)
}
