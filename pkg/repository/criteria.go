package repository

type Criteria struct {
	Filter any
	Sort   any
	Index  *int64
	Size   *int64
}

func (c *Criteria) SetIndex(index int64) *Criteria {
	if index < 0 {
		index = 0
	}
	c.Index = &index
	return c
}

func (c *Criteria) SetSize(size int64) *Criteria {
	if size < 1 {
		size = 1
	}
	if size > 100 {
		size = 100
	}
	c.Size = &size
	return c
}
