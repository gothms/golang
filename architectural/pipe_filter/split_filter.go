package pipe_filter

import (
	"errors"
	"strings"
)

var SplitFilterWrongFormatError = errors.New("input data should be string")

type SplitFilter struct {
	delimiter string
}

func NewSplitFilter(delimiter string) *SplitFilter {
	return &SplitFilter{delimiter}
}
func (sf *SplitFilter) Process(data Request) (Response, error) {
	str, ok := data.(string) // 检查数据格式/类型，是否可处理
	if !ok {
		return nil, SplitFilterWrongFormatError
	}
	return strings.Split(str, sf.delimiter), nil
}
