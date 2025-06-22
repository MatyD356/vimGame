package notionApi

import "errors"

var ErrMissingNotionDatabaseId = errors.New("NOTION_DATABASE_ID is not set in the environment")
