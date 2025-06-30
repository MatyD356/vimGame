package notionservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/MatyD356/vimGame/internals/cache"
	"github.com/MatyD356/vimGame/internals/config"
	"github.com/MatyD356/vimGame/internals/integrations/notion"
)

func parseNumber(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	var builder strings.Builder
	for i, r := range s {
		if unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else if r == '-' && i == 0 {
			builder.WriteRune(r)
		}
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

func GetDatabase(cfg *config.Config) error {
	fmt.Println("GetDatabase called")
	notionDbId := cfg.Env.NotionDbId
	fmt.Println("Notion Database ID:", notionDbId)
	if notionDbId == "" {
		return errors.New("missing database ID")
	}
	req, err := notion.GetDatabaseReq(notionDbId, cfg)
	if err != nil {
		return err
	}
	resp, err := cfg.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode) + ", response: " + string(dat))
	}
	var database notion.Database
	err = json.Unmarshal(dat, &database)
	if err != nil {
		return err
	}
	// Get sub-database IDs from the results
	var pages []cache.PageCache
	for _, result := range database.Results {
		if result.Object == "page" {
			//on a full run this should be done in parallel
			cachePage, ok := cfg.Cache.GetPage(result.ID)
			if !ok {
				fmt.Println("Cache miss for page ID:", result.ID)
				cachePage, err = GetPageChildrenDatabaseId(result.ID, result.Properties.Name.Title[0].PlainText, cfg)
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
		subDbUrl := notion.BaseURL + "databases/" + page.ChildDatabaseId + "/query"
		subReq, err := notion.GetDatabaseReq(subDbUrl, cfg)
		if err != nil {
			return err
		}
		subResp, err := cfg.HttpClient.Do(subReq)
		fmt.Println("Request URL:", req.URL.String())

		if err != nil {
			return err
		}
		defer subResp.Body.Close()
		if subResp.StatusCode != http.StatusOK {
			return errors.New("unexpected status code for sub-database:" + page.ChildDatabaseId)
		}
		var subDatabase notion.Database
		subDat, err := io.ReadAll(subResp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(subDat, &subDatabase)
		if err != nil {
			return err
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
	return nil
}
