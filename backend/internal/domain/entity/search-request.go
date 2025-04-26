package entity

type SearchRequest struct {
	Query   string
	Filters Filters
	Limit   int
	Offset  int
}

type Filters struct {
	Genre   string
	Country string
}
