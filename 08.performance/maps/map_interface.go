package maps

type Map interface {
	Get(interface{}) (interface{}, bool)
	Set(interface{}, interface{})
	Del(interface{})
}
