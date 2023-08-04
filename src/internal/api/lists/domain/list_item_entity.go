package domain

type ListItemEntity struct {
	ID          int32                      `json:"id"`
	ListID      int32                      `json:"-"`
	UserID      int32                      `json:"-"`
	Title       ItemTitleValueObject       `json:"title"`
	Description ItemDescriptionValueObject `json:"description"`
	Position    int32                      `json:"position"`
}
