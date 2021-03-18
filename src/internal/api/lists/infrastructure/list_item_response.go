package infrastructure

type ListItemResponse struct {
	ID          int32  `json:"id"`
	ListID      int32  `json:"listId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
