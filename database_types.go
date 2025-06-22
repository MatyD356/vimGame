package main

import "time"

type DatabaseProperty struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title []struct {
		Type      string `json:"type"`
		PlainText string `json:"plain_text"`
	} `json:"title"`
	RichText  []NotionRichText `json:"rich_text"`
	PlainText string           `json:"plain_text"`
}

type NotionRichText struct {
	Type string `json:"type"`
	Text struct {
		Content string `json:"content"`
		Link    any    `json:"link"`
	} `json:"text"`
	PlainText string `json:"plain_text"`
}

type NotionEditEvent struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type NotionResult struct {
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

type NotionDatabase struct {
	Object  string         `json:"object"`
	Results []NotionResult `json:"results"`
}
