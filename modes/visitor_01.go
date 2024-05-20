package modes

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Visitor01 simple demo
type Visitor01 func(shape Shape)

type Shape interface {
	accept(Visitor01)
}
type Circle struct {
	Radius int
}

func (c Circle) accept(v Visitor01) {
	v(c)
}

type Rectangle struct {
	Width, Heigh int
}

func (r Rectangle) accept(v Visitor01) {
	v(r)
}
func JsonVisitor(shape Shape) {
	bytes, err := json.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
func XmlVisitor(shape Shape) {
	bytes, err := xml.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
func Visitor01Test() {
	c := Circle{10}
	r := Rectangle{100, 200}
	shapes := []Shape{c, r}

	for _, s := range shapes {
		s.accept(JsonVisitor)
		s.accept(XmlVisitor)
	}
}
