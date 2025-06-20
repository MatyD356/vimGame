package env

import "errors"

var ErrMissingNotionSecret = errors.New("NOTION_SECRET is not set in the environment")
