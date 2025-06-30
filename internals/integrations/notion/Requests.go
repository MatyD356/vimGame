package notion

import (
	"net/http"

	"github.com/MatyD356/vimGame/internals/config"
)

func GetDatabaseReq(notionDbId string, cfg *config.Config) (*http.Request, error) {

	url := BaseURL + "databases/" + notionDbId + "/query"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Env.NotionSecret)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	return req, nil
}
