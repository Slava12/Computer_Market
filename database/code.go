package database

// Code хранит данные о кодах подтверждения пользователей
type Code struct {
	ID     int
	Code   int
	UserID int
}

// NewCode добавляет новый код подтверждения в базу данных
func NewCode(code int, userID int) (id int, err error) {
	err = db.QueryRow("insert into codes (code, user_id) values ($1, $2) returning id", code, userID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DelCode удаляет код подтверждения из базы данных
func DelCode(ID int) error {
	_, err := db.Exec("delete from codes where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// DelAllCodes удаляет все коды подтверждения из базы данных
func DelAllCodes() error {
	_, err := db.Exec("delete from codes")
	if err != nil {
		return err
	}
	return nil
}

// GetCode возвращает данные о коде подтверждения по его ID
func GetCode(ID int) (Code, error) {
	row := db.QueryRow("select * from codes where id=$1", ID)
	code := Code{}
	err := row.Scan(&code.ID, &code.Code, &code.UserID)
	if err != nil {
		return Code{}, err
	}
	return code, nil
}

// GetCodeByUserID возвращает данные о коде подтверждения по ID пользователя
func GetCodeByUserID(userID int) (Code, error) {
	row := db.QueryRow("select * from codes where user_id=$1", userID)
	code := Code{}
	err := row.Scan(&code.ID, &code.Code, &code.UserID)
	if err != nil {
		return Code{}, err
	}
	return code, nil
}

// GetAllCodes возвращает данные обо всех кодах подтверждения
func GetAllCodes() ([]Code, error) {
	rows, err := db.Query("select * from codes order by id asc")
	if err != nil {
		return []Code{}, err
	}
	codes := []Code{}
	code := Code{}
	for rows.Next() {
		err = rows.Scan(&code.ID, &code.Code, &code.UserID)
		if err != nil {
			return []Code{}, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

// UpdateCode обновляет значения полей кода подтверждения
func UpdateCode(ID int, code int, userID int) error {
	_, err := db.Exec("update codes set code = $1, user_id = $2 where id = $3",
		code, userID, ID)
	if err != nil {
		return err
	}
	return nil
}
