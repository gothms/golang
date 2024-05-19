package pipe_filter

import (
	"errors"
)

var SumWrongFormatError = errors.New("input data should be []int")

type SumFilter struct{}

func NewSumFilter() *SumFilter {
	return &SumFilter{}
}
func (sf *SumFilter) Process(data Request) (Response, error) {
	arr, ok := data.([]int) // 检查数据格式/类型，是否可处理
	if !ok {
		return nil, SumWrongFormatError
	}
	ret := 0
	for _, v := range arr {
		ret += v
	}
	return ret, nil
}
