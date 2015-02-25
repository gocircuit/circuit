package render

import (
	"sync"
)

// Unknown is a computable value, which is determined when requested by a method call.
type Unknown interface {
	Value() interface{}
}

// Variable is an Unknown whose value can be set once.
type Variable chan interface{}

func NewVariable() Variable {
	return make(Variable, 1)
}

func (u Variable) Fix(v interface{}) {
	u <- v
}

func (u Variable) Value() interface{} {
	r := <-u
	u <- r
	return r
}

// Compute is an Unknown which invokes a function if its value is needed.
type Compute struct {
	logic  func() interface{}
	once   sync.Once
	result interface{}
}

func NewCompute(logic func() interface{}) Unknown {
	return &Compute{logic: logic}
}

func (c *Compute) Value() interface{} {
	c.once.Do(c.compute)
	return c.result
}

func (c *Compute) compute() {
	c.result = c.logic()
}

// // Page
// type Page struct {
// 	Url    Unknown
// 	Render Unknown
// }

// type PageRender func(*Page) string

// func NewPage(url Unknown, renderer PageRender) *Page {
// 	p := &Page{Url: url}
// 	p.Render = NewCompute(func() interface{} { return renderer(p) })
// 	return p
// }

// type Index struct {
// 	page   []*Page
// 	Render Unknown
// }

// func NewIndex(page ...*Page) *Index {
// }
