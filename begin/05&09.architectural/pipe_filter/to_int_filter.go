package pipe_filter

import (
	"errors"
	"strconv"
)

var ToIntFilterWrongFormatError = errors.New("input data should be []string")

type ToIntFilter struct{}

func NewToIntFilter() *ToIntFilter {
	return &ToIntFilter{}
}
func (tif *ToIntFilter) Process(data Request) (Response, error) {
	str, ok := data.([]string) // 检查数据格式/类型，是否可处理
	if !ok {
		return nil, ToIntFilterWrongFormatError
	}
	ret := make([]int, len(str))
	for i := 0; i < len(str); i++ {
		v, err := strconv.Atoi(str[i])
		if err != nil {
			return nil, err
		}
		ret[i] = v
	}
	return ret, nil
}
