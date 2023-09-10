package generic

type DefineGeneric[K comparable, V any] struct {
	key K
	val V
	m   map[K]V
}
