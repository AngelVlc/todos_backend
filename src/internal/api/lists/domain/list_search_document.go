package domain

type ListSearchDocument struct {
	ObjectID          string   `json:"objectID"`
	UserID            int32    `json:"userID"`
	Name              string   `json:"name"`
	ItemsTitles       []string `json:"itemsTitles"`
	ItemsDescriptions []string `json:"itemsDescriptions"`
}
