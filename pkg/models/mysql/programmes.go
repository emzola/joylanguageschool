package mysql

import (
	"database/sql"

	"joylanguageschool.ru/pkg/models"
)

type ProgrammeModel struct {
	DB *sql.DB
}

// Insert post to database
func (m *ProgrammeModel) Insert(title, content, image string) error {
	query := `INSERT INTO programmes (title, content, image, created)
						VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err := m.DB.Exec(query, title, content, image)
	if err != nil {
		return err
	}

	return nil
}

// Gets a specific post from database
func (m *ProgrammeModel) Get(id int) (*models.Programme, error) {
	query := `SELECT id, title, content, image, created FROM programmes
						WHERE id = ?`

	row := m.DB.QueryRow(query, id)

	p := &models.Programme{}

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return p, nil
}

// Get all posts from database
func (m *ProgrammeModel) GetAllProgrammes() ([]*models.Programme, error) {
	query := `SELECT id, title, content, image, created FROM programmes
	ORDER BY created DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	programmes := []*models.Programme{}

	for rows.Next() {
		p := &models.Programme{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
		if err != nil {
			return nil, err
		}

		programmes = append(programmes, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return programmes, nil
}

// Get latest 5 posts from database
func (m *ProgrammeModel) GetLastFivePosts() ([]*models.Programme, error) {
	query := `SELECT id, title, content, image, created FROM programmes
	ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	programmes := []*models.Programme{}

	for rows.Next() {
		p := &models.Programme{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
		if err != nil {
			return nil, err
		}

		programmes = append(programmes, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return programmes, nil
}

// Get latest 3 posts from database
func (m *ProgrammeModel) GetLastThreePosts() ([]*models.Programme, error) {
	query := `SELECT id, title, content, image, created FROM programmes
	ORDER BY created DESC LIMIT 3`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	programmes := []*models.Programme{}

	for rows.Next() {
		p := &models.Programme{}

		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Created)
		if err != nil {
			return nil, err
		}

		programmes = append(programmes, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return programmes, nil
}

func (m *ProgrammeModel) Update(id int, title, content, image string) error {
	if id < 1 {
		return models.ErrNoRecord
	}

	query := `UPDATE programmes SET title=?, content=?, image=? WHERE id=?`

	_, err := m.DB.Exec(query, title, content, image, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *ProgrammeModel) Delete(id int) error {
	if id < 1 {
		return models.ErrNoRecord
	}

	query := `DELETE FROM programmes WHERE id=?`

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
