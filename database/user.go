package database

// User хранит данные о пользователе
type User struct {
	ID          int
	AccessLevel int
	Login       string
	Password    string
	Email       string
	FullName    string
}

// NewUser добавляет нового пользователя в базу данных
func NewUser(accessLevel int, login string, password string, email string, fullName string) (ID int, err error) {
	err = db.QueryRow(`INSERT INTO users (access_level, login, password, email, full_name) VALUES 
	($1, $2, $3, $4, $5) RETURNING id;`, accessLevel, login, password, email, fullName).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}
