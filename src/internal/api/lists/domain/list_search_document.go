package domain

type ListSearchDocument struct {
	ObjectID          string              `json:"objectID"`
	Name              ListNameValueObject `json:"name"`
	ItemsTitles       []string            `json:"itemsTitles"`
	ItemsDescriptions []string            `json:"itemsDescriptions"`
}
