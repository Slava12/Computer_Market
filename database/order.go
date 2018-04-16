package database

// Order хранит данные о характеристике
type Order struct {
	ID     int
	Units  []string
	UserID int
	Cost   string
	Date   string
}

// NewOrder добавляет новый заказ в базу данных
func NewOrder() error {
	return nil
}
