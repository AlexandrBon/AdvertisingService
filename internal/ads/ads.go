package ads

type Ad struct {
	ID           int64
	Title        string `validate:"minmax:1,100"`
	Text         string `validate:"minmax:1,500"`
	AuthorID     int64
	Published    bool
	LastUpdate   string
	CreationDate string
}
