package model

type DBTreeNode struct {
	SourceID        int    `db:"source_id"`
	Depth           int    `db:"depth"` // 0 means 'root'
	ParentFullKey   string `db:"parent_full_key"`
	CurrentNodeJSON string `db:"current_node_json"`
}

type Source struct {
	ID          int    `db:"id"`
	Description string `db:"description"`
	CrawlerID   int    `db:"crawler_id"`
	CrawlerMeta string `db:"crawler_meta"`
	Schedule    string `db:"schedule"`
}
