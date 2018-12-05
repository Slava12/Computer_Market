package database

// Pair хранит данные о пользователе
type Pair struct {
	ID    int
	One   int
	Two   int
	Count int
}

// NewPair добавляет новую пару товаров в базу данных
func NewPair(one int, two int, count int) (id int, err error) {
	err = db.QueryRow("insert into pairs (one, two, count) values ($1, $2, $3) returning id",
		one, two, count).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DelPair удаляет пару из базы данных
func DelPair(ID int) error {
	_, err := db.Exec("delete from pairs where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// DelAllPairs удаляет все пары из базы данных
func DelAllPairs() error {
	_, err := db.Exec("delete from pairs")
	if err != nil {
		return err
	}
	return nil
}

// GetPair возвращает данные о паре по её ID
func GetPair(ID int) (Pair, error) {
	row := db.QueryRow("select * from pairs where id=$1", ID)
	pair := Pair{}
	err := row.Scan(&pair.ID, &pair.One, &pair.Two, &pair.Count)
	if err != nil {
		return Pair{}, err
	}
	return pair, nil
}

// GetPairByUnitsID возвращает данные о паре по ID её товаров
func GetPairByUnitsID(one int, two int) (Pair, error) {
	row := db.QueryRow("select * from pairs where one=$1 and two=$2", one, two)
	pair := Pair{}
	err := row.Scan(&pair.ID, &pair.One, &pair.Two, &pair.Count)
	if err != nil {
		return Pair{}, err
	}
	return pair, nil
}

// GetPairsByUnitID возвращает данные о парах по ID одного из товаров
func GetPairsByUnitID(ID int) ([]Pair, error) {
	rows, err := db.Query("select * from pairs where one=$1 or two=$1 order by count desc", ID)
	if err != nil {
		return []Pair{}, err
	}
	pairs := []Pair{}
	pair := Pair{}
	for rows.Next() {
		err = rows.Scan(&pair.ID, &pair.One, &pair.Two, &pair.Count)
		if err != nil {
			return []Pair{}, err
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}

// GetAllPairs возвращает данные обо всех парах
func GetAllPairs() ([]Pair, error) {
	rows, err := db.Query("select * from pairs order by id asc")
	if err != nil {
		return []Pair{}, err
	}
	pairs := []Pair{}
	pair := Pair{}
	for rows.Next() {
		err = rows.Scan(&pair.ID, &pair.One, &pair.Two, &pair.Count)
		if err != nil {
			return []Pair{}, err
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}

// UpdatePairCount обновляет поле Count в паре товаров
func UpdatePairCount(ID int, count int) error {
	_, err := db.Exec("update pairs set count = $1 where id = $2", count, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePair обновляет пару товаров
func UpdatePair(ID int, one int, two int, count int) error {
	_, err := db.Exec("update pairs set one = $1, two = $2, count = $3 where id = $4", one, two, count, ID)
	if err != nil {
		return err
	}
	return nil
}
