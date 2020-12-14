package feed

type Item struct {
	ID        int    `json:"id" db:"id"`
	ItemTitle string `json:"item_title" db:"item_title"`
	Link      string `json:"link" db:"link"`
	FullItem  string `json:"description" db:"description"`
	SourceID  string `json:"source_id" db:"source_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type Source struct {
	ID  int    `db:"id"`
	URL string `db:"url"`
}
