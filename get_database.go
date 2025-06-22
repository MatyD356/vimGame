package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/MatyD356/vimGame/internals/cache"
	"github.com/MatyD356/vimGame/internals/notionApi"
)

func parseNumber(s string) (int, error) {
	// Trim spaces first
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	// Build a cleaned string keeping digits and one leading '-'
	var builder strings.Builder
	for i, r := range s {
		if unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else if r == '-' && i == 0 {
			// keep leading minus
			builder.WriteRune(r)
		}
		// ignore other chars like '+', '?', etc.
	}

	cleaned := builder.String()
	if cleaned == "" || cleaned == "-" {
		return 0, fmt.Errorf("no digits found")
	}

	num, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (cfg *Config) GetDatabase(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetDatabase called")
	notionDbId := r.PathValue("databaseId")
	fmt.Println("Notion Database ID:", notionDbId)
	if notionDbId == "" {
		http.Error(w, "missing database ID", http.StatusBadRequest)
		return
	}
	url := notionApi.BaseURL + "databases/" + notionDbId + "/query"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		http.Error(w, "failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Env.NotionSecret)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Request URL:", req.URL.String())
	resp, err := cfg.NotionApi.httpClient.Do(req)

	if err != nil {
		http.Error(w, "failed to make request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("unexpected status code: %d, response: %s", resp.StatusCode, string(dat)), resp.StatusCode)
		return
	}
	var database NotionDatabase
	err = json.Unmarshal(dat, &database)
	if err != nil {
		http.Error(w, "failed to unmarshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Get sub-database IDs from the results
	var pages []cache.PageCache
	for _, result := range database.Results {
		if result.Object == "page" {
			//on a full run this should be done in parallel
			cachePage, ok := cfg.Cache.GetPage(result.ID)
			if !ok {
				fmt.Println("Cache miss for page ID:", result.ID)
				cachePage, err = cfg.GetPageChildrenDatabaseId(result.ID, result.Properties.Name.Title[0].PlainText)
				cfg.Cache.SetPage(result.ID, cachePage)
				if err != nil {
					fmt.Println("Error getting page children database ID:", err)
				}
			} else {
				fmt.Println("Cache hit for page ID:", result.ID)
			}
			pages = append(pages, cachePage)
		}
	}
	//sort pages by created time
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].CreatedTime.Before(pages[j].CreatedTime)
	})
	// fetch each sub-database
	for idx, page := range pages {
		cacheChildDb, ok := cfg.Cache.GetParsedChildDatabase(page.ChildDatabaseId)
		if ok {
			fmt.Println("Cache hit for child database ID:", page.ChildDatabaseId)
			pages[idx].ChildDatabase = cacheChildDb
			continue
		}
		if page.ChildDatabaseId == "" {
			continue // skip empty IDs
		}
		subDbUrl := notionApi.BaseURL + "databases/" + page.ChildDatabaseId + "/query"
		subReq, err := http.NewRequest("POST", subDbUrl, nil)
		if err != nil {
			http.Error(w, "failed to create sub-database request: "+err.Error(), http.StatusInternalServerError)
			return
		}
		subReq.Header.Set("Authorization", "Bearer "+cfg.Env.NotionSecret)
		subReq.Header.Set("Notion-Version", "2022-06-28")
		subReq.Header.Set("Content-Type", "application/json")
		subResp, err := cfg.NotionApi.httpClient.Do(subReq)
		fmt.Println("Request URL:", req.URL.String())

		if err != nil {
			http.Error(w, "failed to make sub-database request: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer subResp.Body.Close()
		if subResp.StatusCode != http.StatusOK {
			subDat, _ := io.ReadAll(subResp.Body)
			http.Error(w, fmt.Sprintf("unexpected status code for sub-database: %d, response: %s", subResp.StatusCode, string(subDat)), subResp.StatusCode)
			return
		}
		var subDatabase NotionDatabase
		subDat, err := io.ReadAll(subResp.Body)
		if err != nil {
			http.Error(w, "failed to read sub-database response body: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(subDat, &subDatabase)
		if err != nil {
			fmt.Println("Error unmarshalling sub-database response:", err)
			return
		}
		for _, subResult := range subDatabase.Results {
			if subResult.Object == "page" {
				title := ""
				if len(subResult.Properties.Title.Title) > 0 {
					title = subResult.Properties.Title.Title[0].PlainText
					if title == "Ćwiczenie" {
						continue
					}
				}
				value := ""
				if len(subResult.Properties.ColumnOne.RichText) > 0 {
					parsedValue, err := parseNumber(subResult.Properties.ColumnOne.RichText[0].PlainText)
					if err == nil {
						value = strconv.Itoa(parsedValue)
					} else {
						value = subResult.Properties.ColumnOne.RichText[0].PlainText
					}
				}
				cacheChildDbEntry := cache.ChildDatabaseCache{
					ID:    subResult.ID,
					Title: title,
					Value: value,
				}
				pages[idx].ChildDatabase = append(pages[idx].ChildDatabase, cacheChildDbEntry)
				cfg.Cache.SetParsedChildDatabase(page.ChildDatabaseId, pages[idx].ChildDatabase)
			}
		}

	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if err := json.NewEncoder(w).Encode(pages); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
