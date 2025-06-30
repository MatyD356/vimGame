package notionservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MatyD356/vimGame/internals/cache"
	"github.com/MatyD356/vimGame/internals/config"
	"github.com/MatyD356/vimGame/internals/integrations/notion"
)

func GetPageChildrenDatabaseId(pageId string, pageTitle string, cfg *config.Config) (cache.PageCache, error) {
	fmt.Println("GetPageBlocks page ID:", pageId)
	if pageId == "" {
		return cache.PageCache{}, fmt.Errorf("missing page ID")
	}
	url := notion.BaseURL + "blocks/" + pageId + "/children"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return cache.PageCache{}, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Env.NotionSecret)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfg.HttpClient.Do(req)
	fmt.Println("Request URL:", req.URL.String())
	if err != nil {
		return cache.PageCache{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return cache.PageCache{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var blocks notion.Block
	if err := json.NewDecoder(resp.Body).Decode(&blocks); err != nil {
		return cache.PageCache{}, fmt.Errorf("failed to decode response: %w", err)
	}
	for _, result := range blocks.Results {
		if result.Type == "child_database" {
			fmt.Println("Found child database with ID:", result.ID)
			return cache.PageCache{
				ChildDatabaseId: result.ID,
				CreatedTime:     result.CreatedTime,
				Title:           pageTitle,
			}, nil
		}
	}
	return cache.PageCache{}, fmt.Errorf("no child database found in page blocks")
}
