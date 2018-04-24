package database

// Feature хранит данные о характеристике
type Feature struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// NewFeature добавляет новую характеристику в базу данных
func NewFeature(name string) (id int, err error) {
	err = db.QueryRow("insert into features (name) values ($1) returning id", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DelFeature удаляет характеристику из базы данных
func DelFeature(ID int) error {
	_, err := db.Exec("delete from features where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetFeature возвращает данные о характеристике по её ID
func GetFeature(ID int) (Feature, error) {
	row := db.QueryRow("select * from features where id=$1", ID)
	feature := Feature{}
	err := row.Scan(&feature.ID, &feature.Name)
	if err != nil {
		return Feature{}, err
	}
	return feature, nil
}

// GetAllFeatures возвращает данные обо всех характеристиках
func GetAllFeatures() ([]Feature, error) {
	rows, err := db.Query("select * from features order by id asc")
	if err != nil {
		return []Feature{}, err
	}
	features := []Feature{}
	feature := Feature{}
	for rows.Next() {
		err = rows.Scan(&feature.ID, &feature.Name)
		if err != nil {
			return []Feature{}, err
		}
		features = append(features, feature)
	}
	return features, nil
}

// UpdateFeature обновляет значения полей характеристики
func UpdateFeature(ID int, name string) error {
	_, err := db.Exec("update features set name = $1 where id = $2", name, ID)
	if err != nil {
		return err
	}
	return nil
}
