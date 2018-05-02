package errortemplate

import (
	"strings"
)

// Error хранит сообщение об ошибке
type Error struct {
	Message string
	Link    string
}

// GenerateMessage формирует сообщение об ошибке
func GenerateMessage(err error) (message string) {
	if strings.Contains(err.Error(), "pq: повторяющееся значение ключа нарушает ограничение уникальности") {
		stringArray := strings.Split(err.Error(), "\"")
		stringArray = strings.Split(stringArray[1], "_")
		if stringArray[0] == "features" {
			message = "Характеристика с таким названием уже существует!"
		}
		if stringArray[0] == "categories" {
			message = "Категория с таким названием уже существует!"
		}
		if stringArray[0] == "units" {
			message = "Товар с таким названием уже существует!"
		}
	}
	return message
}
