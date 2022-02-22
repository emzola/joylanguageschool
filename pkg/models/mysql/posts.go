package mysql

import (
	"database/sql"

	"joylanguageschool.ru/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

// Insert post to database
func (m *PostModel) Insert(title, content, image string) error {
	query := `INSERT INTO posts (title, content, image, created)
						VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err := m.DB.Exec(query, title, content, image)
	if err != nil {
		return err
	}

	return nil
}

// Gets a specific post from database
func (m *PostModel) Get(id int) (*models.Post, error) {
	query := `SELECT id, title, content, image, created FROM posts
						WHERE id = ?`

	row := m.DB.QueryRow(query, id)

	p := &models.Post{}

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

// Get all posts from database
func (m *PostModel) GetAllPosts() ([]*models.Post, error) {
	query := `SELECT id, title, content, image, created FROM posts
	ORDER BY created DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
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

// Get latest 5 posts from database
func (m *PostModel) GetLastFivePosts() ([]*models.Post, error) {
	query := `SELECT id, title, content, image, created FROM posts
	ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
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

// Get latest 3 posts from database
func (m *PostModel) GetLastThreePosts() ([]*models.Post, error) {
	query := `SELECT id, title, content, image, created FROM posts
	ORDER BY created DESC LIMIT 3`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		p := &models.Post{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
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

func (m *PostModel) Update(id int, title, content, image string) error {
	if id < 1 {
		return models.ErrNoRecord
	}

	query := `UPDATE posts SET title=?, content=?, image=? WHERE id=?`

	_, err := m.DB.Exec(query, title, content, image, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) Delete(id int) error {
	if id < 1 {
		return models.ErrNoRecord
	}

	query := `DELETE FROM posts WHERE id=?`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrNoRecord
	}

	return nil
}
