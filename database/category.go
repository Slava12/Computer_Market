package database

// Category хранит данные о категории
type Category struct {
	ID       int
	ParentID int
	Name     string
	Features string // заменить!!!
}
