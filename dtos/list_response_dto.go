package dtos

type ListResponseDto struct {
	ID        int32                  `json:"id"`
	Name      string                 `json:"name"`
	ListItems []*ListItemResponseDto `json:"items"`
}
