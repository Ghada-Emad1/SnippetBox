package models

import (
	"database/sql"
	
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	var pk int
	expiresTime := time.Now().AddDate(0, 0, expires)
	err := m.DB.QueryRow(`INSERT INTO snippets (title, content, expires) VALUES ($1,$2,$3) RETURNING id`, title, content, expiresTime).Scan(&pk)
	if err != nil {
		return 0, err
	}
	

	return pk, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
