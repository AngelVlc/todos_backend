package infrastructure

type ListResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	ItemsCount  int32  `json:"itemsCount"`
	IsQuickList bool   `json:"isQuickList"`
}
