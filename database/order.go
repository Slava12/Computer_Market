package database

import "encoding/json"

// Order хранит данные о характеристике
type Order struct {
	ID     int
	Units  []string
	UserID int
	Cost   int
	Date   string
}

// NewOrder добавляет новый заказ в базу данных
func NewOrder(units []string, userID int, cost int, date string) (id int, err error) {
	unitsJSON, errMarshal := json.Marshal(units)
	if errMarshal != nil {
		return 0, errMarshal
	}
	err = db.QueryRow("insert into orders (units, user_id, cost, date) values ($1, $2, $3, $4) returning id",
		unitsJSON, userID, cost, date).Scan(&id)
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
	var units string
	err := row.Scan(&order.ID, &units, &order.UserID, &order.Cost, &order.Date)
	if err != nil {
		return Order{}, err
	}
	errU := json.Unmarshal([]byte(units), &order.Units)
	if errU != nil {
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
	units := []string{}
	unit := ""
	for rows.Next() {
		err = rows.Scan(&order.ID, &units, &order.UserID, &order.Cost, &order.Date)
		if err != nil {
			return []Order{}, err
		}
		orders = append(orders, order)
		units = append(units, unit)
	}
	for i := 0; i < len(orders); i++ {
		errU := json.Unmarshal([]byte(units[i]), &orders[i].Units)
		if errU != nil {
			return []Order{}, err
		}
	}
	return orders, nil
}

// UpdateOrder обновляет значения полей заказа
func UpdateOrder(ID int, units []string, userID int, cost int, date string) error {
	unitsJSON, errMarshal := json.Marshal(units)
	if errMarshal != nil {
		return errMarshal
	}
	_, err := db.Exec("update orders set units = $1, user_id = $2, cost = $3, date = $4 where id = $5",
		unitsJSON, userID, cost, date, ID)
	if err != nil {
		return err
	}
	return nil
}
