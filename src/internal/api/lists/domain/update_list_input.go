package domain

type UpdateListInput struct {
	Name          ListNameValueObject `json:"name"`
	IDsByPosition []int32             `json:"idsByPosition"`
}
