package notion

import "time"

type Block struct {
	Object  string `json:"object"`
	Results []struct {
		Object      string    `json:"object"`
		ID          string    `json:"id"`
		Type        string    `json:"type"`
		CreatedTime time.Time `json:"created_time"`
		Parent      struct {
			Type   string `json:"type"`
			PageID string `json:"page_id"`
		} `json:"parent"`
		ChildDatabase struct {
			Title string `json:"title"`
		} `json:"child_database"`
	} `json:"results"`
	Type string `json:"type"`
}
