package database

// User хранит данные о пользователе
type User struct {
	ID          int
	AccessLevel int
	Login       string
	Password    string
	Email       string
	FirstName   string
	SecondName  string
}

// NewUser добавляет нового пользователя в базу данных
func NewUser(accessLevel int, login string, password string, email string, firstName string, secondName string) error {
	_, err := db.Exec("insert into users (access_level, login, password, email, first_name, second_name) values ($1, $2, $3, $4, $5, $6)",
		accessLevel, login, password, email, firstName, secondName)
	if err != nil {
		return err
	}
	return nil
}

// DelUser удаляет пользователя из базы данных
func DelUser(ID int) error {
	_, err := db.Exec("delete from users where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetUser возвращает данные о пользователе по его ID
func GetUser(ID int) (User, error) {
	row := db.QueryRow("select * from users where id=$1", ID)
	user := User{}
	err := row.Scan(&user.ID, &user.AccessLevel, &user.Login, &user.Password, &user.Email, &user.FirstName, &user.SecondName)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetAllUsers возвращает данные обо всех пользователях
func GetAllUsers() ([]User, error) {
	rows, err := db.Query("select * from users")
	if err != nil {
		return []User{}, err
	}
	users := []User{}
	user := User{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.AccessLevel, &user.Login, &user.Password, &user.Email, &user.FirstName, &user.SecondName)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUserAccessLevel обновляет значение уровня доступа пользователя
func UpdateUserAccessLevel(ID int, accessLevel int) error {
	_, err := db.Exec("update users set access_level = $1 where id = $2", accessLevel, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserLogin обновляет значение логина пользователя
func UpdateUserLogin(ID int, login string) error {
	_, err := db.Exec("update users set login = $1 where id = $2", login, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserPassword обновляет значение пароля пользователя
func UpdateUserPassword(ID int, password string) error {
	_, err := db.Exec("update users set password = $1 where id = $2", password, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserEmail обновляет значение адреса почты пользователя
func UpdateUserEmail(ID int, email string) error {
	_, err := db.Exec("update users set email = $1 where id = $2", email, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserFirstName обновляет значение имени пользователя
func UpdateUserFirstName(ID int, firstName string) error {
	_, err := db.Exec("update users set first_name = $1 where id = $2", firstName, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUserSecondName обновляет значение фамилии пользователя
func UpdateUserSecondName(ID int, secondName string) error {
	_, err := db.Exec("update users set second_name = $1 where id = $2", secondName, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser обновляет значения полей пользователя
func UpdateUser(ID int, accessLevel int, login string, password string, email string, firstName string, secondName string) error {
	_, err := db.Exec("update users set access_level = $1, login = $2, password = $3, email = $4, first_name = $5, second_name = $6 where id = $7",
		accessLevel, login, password, email, firstName, secondName, ID)
	if err != nil {
		return err
	}
	return nil
}
