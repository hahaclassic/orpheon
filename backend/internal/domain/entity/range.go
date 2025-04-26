package entity

type Range struct {
	Start int
	End   int
}

func (r Range) Len() int {
	return r.End - r.Start
}
