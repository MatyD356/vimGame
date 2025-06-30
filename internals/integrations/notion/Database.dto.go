package notion

import "time"

type DatabaseProperty struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title []struct {
		Type      string `json:"type"`
		PlainText string `json:"plain_text"`
	} `json:"title"`
	RichText  []RichText `json:"rich_text"`
	PlainText string     `json:"plain_text"`
}

type RichText struct {
	Type string `json:"type"`
	Text struct {
		Content string `json:"content"`
		Link    any    `json:"link"`
	} `json:"text"`
	PlainText string `json:"plain_text"`
}

type EditEvent struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type Result struct {
	Object      string    `json:"object"`
	ID          string    `json:"id"`
	CreatedTime time.Time `json:"created_time"`
	Parent      struct {
		Type       string `json:"type"`
		DatabaseID string `json:"database_id"`
	} `json:"parent"`
	Properties struct {
		ColumnOne DatabaseProperty `json:"Column 1"`
		Title     DatabaseProperty `json:"Title"`
		Name      DatabaseProperty `json:"Name"`
	} `json:"properties"`
}

type Database struct {
	Object  string   `json:"object"`
	Results []Result `json:"results"`
}
