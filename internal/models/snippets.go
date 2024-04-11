package models

import (
	"database/sql"
	"errors"

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
	stat := `SELECT id,title,content,created,expires FROM snippets WHERE expires > CURRENT_TIMESTAMP AND id=$1`
	row := m.DB.QueryRow(stat, id)

	s := &Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stat := `SELECT id,title,content,created,expires FROM snippets WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
