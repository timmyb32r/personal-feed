package model

import "encoding/json"

type DBTreeNode struct {
	SourceID        int    `db:"source_id"`
	Depth           int    `db:"depth"` // 0 means 'root'
	ParentFullKey   string `db:"parent_full_key"`
	CurrentNodeJSON string `db:"current_node_json"`
}

type Source struct {
	ID                 int    `db:"id"`
	Description        string `db:"description"`
	CrawlerID          int    `db:"crawler_id"`
	CrawlerMeta        string `db:"crawler_meta"`
	Schedule           string `db:"schedule"` // https://en.wikipedia.org/wiki/Cron
	NumShouldBeMatched *int   `db:"num_should_be_matched"`
}

func (c *Source) ToJSON() string {
	serializedSource, _ := json.Marshal(c)
	return string(serializedSource)
}

type NumMatchedNotifier func(crawlerDescr string, expected *int, real int)
