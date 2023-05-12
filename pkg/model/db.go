package model

import (
	"encoding/json"
	"time"
)

type DBTreeNode struct {
	SourceID        int       `db:"source_id"`
	Depth           int       `db:"depth"` // 0 means 'root'
	CurrentFullKey  string    `db:"current_full_key"`
	CurrentNodeJSON string    `db:"current_node_json"` // here are serialized object of current depth type
	BusinessTime    time.Time `db:"current_node_json"`
}

func (n *DBTreeNode) ParentFullKey() string {
	parentFullKey, _ := ParseComplexKey(n.CurrentFullKey)
	return parentFullKey.ParentKey().FullKey()
}

type Source struct {
	ID                 int    `db:"id"`
	Description        string `db:"description"`
	CrawlerID          int    `db:"crawler_id"`
	CrawlerMeta        string `db:"crawler_meta"`
	Schedule           string `db:"schedule"` // https://en.wikipedia.org/wiki/Cron
	NumShouldBeMatched *int   `db:"num_should_be_matched"`
	HistoryState       string `db:"history_state"`
}

func (c *Source) ToJSON() string {
	serializedSource, _ := json.Marshal(c)
	return string(serializedSource)
}

type NumMatchedNotifier func(crawlerDescr string, expected *int, real int)
