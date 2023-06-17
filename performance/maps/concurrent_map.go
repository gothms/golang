package maps

import cm "github.com/easierway/concurrent_map"

type ConcurrentMapBenchmarkAdapter struct {
	cm *cm.ConcurrentMap
}

func (m *ConcurrentMapBenchmarkAdapter) Set(key interface{}, value interface{}) {
	m.cm.Set(cm.StrKey(key.(string)), value)
}

func (m *ConcurrentMapBenchmarkAdapter) Get(key interface{}) (interface{}, bool) {
	return m.cm.Get(cm.StrKey(key.(string)))
}

func (m *ConcurrentMapBenchmarkAdapter) Del(key interface{}) {
	m.cm.Del(cm.StrKey(key.(string)))
}

func CreateConcurrentMapBenchmarkAdapter(numOfPartitions int) *ConcurrentMapBenchmarkAdapter {
	conMap := cm.CreateConcurrentMap(numOfPartitions)
	return &ConcurrentMapBenchmarkAdapter{conMap}
}
