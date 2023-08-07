package infrastructure

type MoveListItemInput struct {
	OriginListItemID  int32 `json:"originListItemId"`
	DestinationListID int32 `json:"destinationListItemId"`
}
