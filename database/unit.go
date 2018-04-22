package database

import "encoding/json"

// FeatureUnit хранит данные о характеристиках товара
type FeatureUnit struct {
	FeatureName string `json:"featurename"`
	Value       string `json:"value"`
}

// Unit хранит данные о товаре
type Unit struct {
	ID         int
	Name       string
	CategoryID int
	Quantity   int
	Price      int
	Discount   int
	Features   []FeatureUnit
	Pictures   []string `json:"pictures"`
}

// NewUnit добавляет новый товар в базу данных
func NewUnit(name string, categoryID int, quantity int, price int, discount int, features []FeatureUnit, pictures []string) (id int, err error) {
	featuresJSON, errMarshal := json.Marshal(features)
	if errMarshal != nil {
		return 0, errMarshal
	}
	picturesJSON, errMarshal := json.Marshal(pictures)
	if errMarshal != nil {
		return 0, errMarshal
	}
	err = db.QueryRow("insert into units (name, category_id, quantity, price, discount, features, pictures) values ($1, $2, $3, $4, $5, $6, $7) returning id",
		name, categoryID, quantity, price, discount, featuresJSON, picturesJSON).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DelUnit удаляет товар из базы данных
func DelUnit(ID int) error {
	_, err := db.Exec("delete from units where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetUnit возвращает данные о товаре по его ID
func GetUnit(ID int) (Unit, error) {
	row := db.QueryRow("select * from units where id=$1", ID)
	unit := Unit{}
	err := row.Scan(&unit.ID, &unit.Name, &unit.CategoryID, &unit.Quantity, &unit.Price, &unit.Discount, &unit.Features, &unit.Pictures)
	if err != nil {
		return Unit{}, err
	}
	return unit, nil
}

// GetAllUnits возвращает данные обо всех товарах
func GetAllUnits() ([]Unit, error) {
	rows, err := db.Query("select * from units")
	if err != nil {
		return []Unit{}, err
	}
	units := []Unit{}
	unit := Unit{}
	for rows.Next() {
		err = rows.Scan(&unit.ID, &unit.Name, &unit.CategoryID, &unit.Quantity, &unit.Price, &unit.Discount, &unit.Features, &unit.Pictures)
		if err != nil {
			return []Unit{}, err
		}
		units = append(units, unit)
	}
	return units, nil
}

// UpdateUnit обновляет значения полей товара
func UpdateUnit(ID int, name string, categoryID int, quantity int, price int, discount int, features []FeatureUnit, pictures []string) error {
	_, err := db.Exec("update categories set name = $1, category_id = $2, quantity = $3, price = $4, discount = $5, features = $6, pictures = $7 where id = $8",
		name, categoryID, quantity, price, discount, features, pictures, ID)
	if err != nil {
		return err
	}
	return nil
}
