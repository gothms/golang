package maps

import "github.com/orcaman/concurrent-map/v2"

type OrcamanConcurrentMapBenchmarkAdapter struct {
	cmap *cmap.ConcurrentMap[string, any]
}

func (o *OrcamanConcurrentMapBenchmarkAdapter) Get(key interface{}) (interface{}, bool) {
	return o.cmap.Get(key.(string))
}

func (o *OrcamanConcurrentMapBenchmarkAdapter) Set(key interface{}, value interface{}) {
	o.cmap.Set(key.(string), value)
}

func (o *OrcamanConcurrentMapBenchmarkAdapter) Del(key interface{}) {
	o.cmap.Pop(key.(string))
}
func CreateOrcamanConcurrentMapBenchmarkAdapter() *OrcamanConcurrentMapBenchmarkAdapter {
	c := cmap.New[any]()
	return &OrcamanConcurrentMapBenchmarkAdapter{&c}
}
