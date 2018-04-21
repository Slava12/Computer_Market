package database

import (
	"encoding/json"
)

// Category хранит данные о категории
type Category struct {
	ID       int
	ParentID int
	Name     string
	Features []Feature
}

// NewCategory добавляет новую категорию в базу данных
func NewCategory(parentID int, name string, features []Feature) (id int, err error) {
	featuresJSON, errMarshal := json.Marshal(features)
	if errMarshal != nil {
		return 0, errMarshal
	}
	err = db.QueryRow("insert into categories (parent_id, name, features) values ($1, $2, $3) returning id", parentID, name, featuresJSON).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DelCategory удаляет категорию из базы данных
func DelCategory(ID int) error {
	_, err := db.Exec("delete from categories where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetCategory возвращает данные о категории по её ID
func GetCategory(ID int) (Category, error) {
	row := db.QueryRow("select * from categories where id=$1", ID)
	category := Category{}
	var features string
	err := row.Scan(&category.ID, &category.ParentID, &category.Name, &features)
	if err != nil {
		return Category{}, err
	}
	errU := json.Unmarshal([]byte(features), &category.Features)
	if errU != nil {
		return Category{}, err
	}
	return category, nil
}

// GetAllCategories возвращает данные обо всех категориях
func GetAllCategories() ([]Category, error) {
	rows, err := db.Query("select * from categories")
	if err != nil {
		return []Category{}, err
	}
	categories := []Category{}
	category := Category{}
	var features string
	for rows.Next() {
		err = rows.Scan(&category.ID, &category.ParentID, &category.Name, &features)
		if err != nil {
			return []Category{}, err
		}
		errU := json.Unmarshal([]byte(features), &category.Features)
		if errU != nil {
			return []Category{}, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// UpdateCategory обновляет значения полей категории
func UpdateCategory(ID int, parentID int, name string, features []Feature) error {
	_, err := db.Exec("update categories set parent_id = $1, name = $2, features = $3 where id = $4", parentID, name, features, ID)
	if err != nil {
		return err
	}
	return nil
}
