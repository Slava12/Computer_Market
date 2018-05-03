package database

import "encoding/json"

// FeatureUnit хранит данные о характеристиках товара
type FeatureUnit struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

// DelAllUnits удаляет все товары из базы данных
func DelAllUnits() error {
	_, err := db.Exec("delete from units")
	if err != nil {
		return err
	}
	return nil
}

// GetUnit возвращает данные о товаре по его ID
func GetUnit(ID int) (Unit, error) {
	row := db.QueryRow("select * from units where id=$1", ID)
	unit := Unit{}
	var features string
	var pictures string
	err := row.Scan(&unit.ID, &unit.Name, &unit.CategoryID, &unit.Quantity, &unit.Price, &unit.Discount, &features, &pictures)
	if err != nil {
		return Unit{}, err
	}
	errU := json.Unmarshal([]byte(features), &unit.Features)
	if errU != nil {
		return Unit{}, err
	}
	errU = json.Unmarshal([]byte(pictures), &unit.Pictures)
	if errU != nil {
		return Unit{}, err
	}
	return unit, nil
}

// GetUnitsByCategoryID возвращает данные обо всех товарах
func GetUnitsByCategoryID(id int) ([]Unit, error) {
	rows, err := db.Query("select * from units where category_id=$1 order by id asc", id)
	if err != nil {
		return []Unit{}, err
	}
	units := []Unit{}
	unit := Unit{}
	features := []string{}
	feature := ""
	pictures := []string{}
	picture := ""
	for rows.Next() {
		err = rows.Scan(&unit.ID, &unit.Name, &unit.CategoryID, &unit.Quantity, &unit.Price, &unit.Discount, &feature, &picture)
		if err != nil {
			return []Unit{}, err
		}
		features = append(features, feature)
		pictures = append(pictures, picture)
		units = append(units, unit)
	}
	for i := 0; i < len(units); i++ {
		errU := json.Unmarshal([]byte(features[i]), &units[i].Features)
		if errU != nil {
			return []Unit{}, err
		}
		errU = json.Unmarshal([]byte(pictures[i]), &units[i].Pictures)
		if errU != nil {
			return []Unit{}, err
		}
	}
	return units, nil
}

// GetAllUnits возвращает данные обо всех товарах
func GetAllUnits() ([]Unit, error) {
	rows, err := db.Query("select * from units order by id asc")
	if err != nil {
		return []Unit{}, err
	}
	units := []Unit{}
	unit := Unit{}
	features := []string{}
	feature := ""
	pictures := []string{}
	picture := ""
	for rows.Next() {
		err = rows.Scan(&unit.ID, &unit.Name, &unit.CategoryID, &unit.Quantity, &unit.Price, &unit.Discount, &feature, &picture)
		if err != nil {
			return []Unit{}, err
		}
		features = append(features, feature)
		pictures = append(pictures, picture)
		units = append(units, unit)
	}
	for i := 0; i < len(units); i++ {
		errU := json.Unmarshal([]byte(features[i]), &units[i].Features)
		if errU != nil {
			return []Unit{}, err
		}
		errU = json.Unmarshal([]byte(pictures[i]), &units[i].Pictures)
		if errU != nil {
			return []Unit{}, err
		}
	}
	return units, nil
}

// UpdateUnit обновляет значения полей товара
func UpdateUnit(ID int, name string, categoryID int, quantity int, price int, discount int, features []FeatureUnit, pictures []string) error {
	featuresJSON, errMarshal := json.Marshal(features)
	if errMarshal != nil {
		return errMarshal
	}
	picturesJSON, errMarshal := json.Marshal(pictures)
	if errMarshal != nil {
		return errMarshal
	}
	_, err := db.Exec("update units set name = $1, category_id = $2, quantity = $3, price = $4, discount = $5, features = $6, pictures = $7 where id = $8",
		name, categoryID, quantity, price, discount, featuresJSON, picturesJSON, ID)
	if err != nil {
		return err
	}
	return nil
}
