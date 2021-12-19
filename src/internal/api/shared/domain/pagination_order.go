package domain

type PaginationOrder int

const (
	OrderAsc PaginationOrder = iota
	OrderDesc
)

func (o PaginationOrder) String() string {
	if o == OrderAsc {
		return "asc"
	}

	return "desc"
}
