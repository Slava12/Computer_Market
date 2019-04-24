package database

// User хранит данные о пользователе
type User struct {
	ID          int
	AccessLevel int
	Confirmed   bool
	Email       string
	Password    string
	FirstName   string
	SecondName  string
	Phone       string
}

// NewUser добавляет нового пользователя в базу данных
func NewUser(accessLevel int, confirmed bool, email string, password string, firstName string, secondName string, phone string) (id int, err error) {
	err = db.QueryRow("insert into users (access_level, confirmed, email, password, first_name, second_name, phone) values ($1, $2, $3, $4, $5, $6, $7) returning id",
		accessLevel, confirmed, email, password, firstName, secondName, phone).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DelUser удаляет пользователя из базы данных
func DelUser(ID int) error {
	_, err := db.Exec("delete from users where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// DelAllUsers удаляет всех пользователей из базы данных
func DelAllUsers() error {
	_, err := db.Exec("delete from users")
	if err != nil {
		return err
	}
	return nil
}

// GetUser возвращает данные о пользователе по его ID
func GetUser(ID int) (User, error) {
	row := db.QueryRow("select * from users where id=$1", ID)
	user := User{}
	err := row.Scan(&user.ID, &user.AccessLevel, &user.Confirmed, &user.Email, &user.Password, &user.FirstName, &user.SecondName, &user.Phone)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUserByEmail возвращает данные о пользователе по его e-mail
func GetUserByEmail(email string) (User, error) {
	row := db.QueryRow("select * from users where email=$1", email)
	user := User{}
	err := row.Scan(&user.ID, &user.AccessLevel, &user.Confirmed, &user.Email, &user.Password, &user.FirstName, &user.SecondName, &user.Phone)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetAllUsers возвращает данные обо всех пользователях
func GetAllUsers() ([]User, error) {
	rows, err := db.Query("select * from users order by id asc")
	if err != nil {
		return []User{}, err
	}
	users := []User{}
	user := User{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.AccessLevel, &user.Confirmed, &user.Email, &user.Password, &user.FirstName, &user.SecondName, &user.Phone)
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

// UpdateUserConfirmed обновляет значение логина пользователя
func UpdateUserConfirmed(ID int, confirmed bool) error {
	_, err := db.Exec("update users set confirmed = $1 where id = $2", confirmed, ID)
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

// UpdateUserPassword обновляет значение пароля пользователя
func UpdateUserPassword(ID int, password string) error {
	_, err := db.Exec("update users set password = $1 where id = $2", password, ID)
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

// UpdateUserPhone обновляет значение телефона
func UpdateUserPhone(ID int, phone string) error {
	_, err := db.Exec("update users set phone = $1 where id = $2", phone, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser обновляет значения полей пользователя
func UpdateUser(ID int, accessLevel int, confirmed bool, email string, password string, firstName string, secondName string, phone string) error {
	_, err := db.Exec("update users set access_level = $1, confirmed = $2, email = $3, password = $4, first_name = $5, second_name = $6, phone = $7 where id = $8",
		accessLevel, confirmed, email, password, firstName, secondName, phone, ID)
	if err != nil {
		return err
	}
	return nil
}
