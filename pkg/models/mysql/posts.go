package mysql

import (
	"database/sql"

	"joylanguageschool.ru/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

// Insert post to database
func (m *PostModel) Insert(title, content string) (int, error) {
	query := `INSERT INTO posts (title, content, created)
						VALUES(?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(query, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Gets a specific post from database
func (m *PostModel) Get(id int) (*models.Post, error) {
	query := `SELECT id, title, content, created FROM posts
						WHERE id = ?`

row := m.DB.QueryRow(query, id)

p := &models.Post{}

err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Created)
if err == sql.ErrNoRows {
	return nil, models.ErrNoRecord
} else if err != nil {
	return nil, err
}

return p, nil
}

// Get all posts from database
func (m *PostModel) GetAll() ([]*models.Post, error) {
	query := `SELECT id, title, content, created FROM posts
	ORDER BY created DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// Get latest posts from database
func (m *PostModel) Latest() ([]*models.Post, error) {
	query := `SELECT id, title, content, created FROM posts
	ORDER BY created DESC LIMIT 3`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

