package database

// Order хранит данные о характеристике
type Order struct {
	ID     int
	State  string
	Units  string
	UserID int
	Cost   int
	Date   string
}

// NewOrder добавляет новый заказ в базу данных
func NewOrder(state string, units string, userID int, cost int, date string) (id int, err error) {
	err = db.QueryRow("insert into orders (state, units, user_id, cost, date) values ($1, $2, $3, $4, $5) returning id",
		state, units, userID, cost, date).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DelOrder удаляет заказ из базы данных
func DelOrder(ID int) error {
	_, err := db.Exec("delete from orders where id = $1", ID)
	if err != nil {
		return err
	}
	return nil
}

// DelAllOrders удаляет все заказы из базы данных
func DelAllOrders() error {
	_, err := db.Exec("delete from orders")
	if err != nil {
		return err
	}
	return nil
}

// GetOrder возвращает данные о заказе по его ID
func GetOrder(ID int) (Order, error) {
	row := db.QueryRow("select * from orders where id=$1", ID)
	order := Order{}
	err := row.Scan(&order.ID, &order.State, &order.Units, &order.UserID, &order.Cost, &order.Date)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

// GetAllOrders возвращает данные обо всех заказах
func GetAllOrders() ([]Order, error) {
	rows, err := db.Query("select * from orders order by id asc")
	if err != nil {
		return []Order{}, err
	}
	orders := []Order{}
	order := Order{}
	for rows.Next() {
		err = rows.Scan(&order.ID, &order.State, &order.Units, &order.UserID, &order.Cost, &order.Date)
		if err != nil {
			return []Order{}, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// UpdateOrder обновляет значения полей заказа
func UpdateOrder(ID int, state string, units string, userID int, cost int, date string) error {
	_, err := db.Exec("update orders set state = $1, units = $2, user_id = $3, cost = $4, date = $5 where id = $6",
		state, units, userID, cost, date, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrderState обновляет статус заказа
func UpdateOrderState(ID int, state string) error {
	_, err := db.Exec("update orders set state = $1 where id = $2",
		state, ID)
	if err != nil {
		return err
	}
	return nil
}
