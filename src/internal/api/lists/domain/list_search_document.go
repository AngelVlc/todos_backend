package domain

type ListSearchDocument struct {
	ObjectID          string              `json:"objectID"`
	UserID            int32               `json:"userID"`
	Name              ListNameValueObject `json:"name"`
	ItemsTitles       []string            `json:"itemsTitles"`
	ItemsDescriptions []string            `json:"itemsDescriptions"`
}
