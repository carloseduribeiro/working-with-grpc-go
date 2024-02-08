package database

import (
	"database/sql"
	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

func NewCourse(db *sql.DB) *Course {
	return &Course{db: db}
}

func (c *Course) Create(name, description, categoryID string) (*Course, error) {
	id := uuid.New().String()
	_, err := c.db.Exec(`INSERT INTO course (id, name, description, category_id) VALUES (?, ?, ?, ?)`, id, name, description, categoryID)
	if err != nil {
		return &Course{}, err
	}
	return &Course{
		ID:          id,
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
	}, nil
}

func (c *Course) FindAll() ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM course")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	courses := make([]Course, 0)
	for rows.Next() {
		var id, name, description, categoryID string
		if err := rows.Scan(&id, &name, &description, &categoryID); err != nil {
			return nil, err
		}
		courses = append(courses, Course{
			ID:          id,
			Name:        name,
			Description: description,
			CategoryID:  categoryID,
		})
	}
	return courses, nil
}

func (c *Course) FindAllByCategoryID(categoryID string) ([]Course, error) {
	stmt, err := c.db.Prepare(`SELECT id, name, description, category_id FROM course WHERE category_id = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(categoryID)
	if err != nil {
		return nil, err
	}
	courses := make([]Course, 0)
	for rows.Next() {
		var id, name, description, categoryID string
		if err := rows.Scan(&id, &name, &description, &categoryID); err != nil {
			return nil, err
		}
		courses = append(courses, Course{
			ID:          id,
			Name:        name,
			Description: description,
			CategoryID:  categoryID,
		})
	}
	return courses, nil
}
