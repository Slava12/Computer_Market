package handlefunc

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

// Record хранит данные о товаре
type Record struct {
	ID    int
	Count int
}

func splitBasket(basket string) []Record {
	units := strings.Split(basket, ";")
	records := []Record{}
	record := Record{}
	for i := range units {
		if i == len(units)-1 { // Последняя запись всегда пустая
			break
		}
		unitInfo := strings.Split(units[i], ":")
		id, err := strconv.Atoi(unitInfo[0])
		if err != nil {
			logger.Warn(err, "Не удалось конвертировать строку в число!")
			return []Record{}
		}
		count, err := strconv.Atoi(unitInfo[1])
		if err != nil {
			logger.Warn(err, "Не удалось конвертировать строку в число!")
			return []Record{}
		}
		record.ID = id
		record.Count = count
		records = append(records, record)
	}
	return records
}

func showBasket(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	basket, _ := session.Values["basket"].(string)
	cost := 0
	records := splitBasket(basket)
	data := make([]Data, len(records))
	for i, record := range records {
		unit, err := database.GetUnit(record.ID)
		if err != nil {
			logger.Warn(err, "Не удалось получить запись о товаре ", record.ID, "!")
			return
		}
		if len(unit.Pictures) > 0 {
			data[i].Picture = unit.Pictures[0]
		}
		data[i].LinkUnit = "/units/" + strconv.Itoa(unit.ID)
		data[i].Name = unit.Name
		data[i].Price = unit.Price
		data[i].Link = "/remove_from_basket/" + strconv.Itoa(unit.ID)
		data[i].Text = "Убрать из корзины"
		cost += unit.Price * record.Count
	}
	dataFull := DataFull{}
	dataFull.ShowCategory = false
	dataFull.CategoryNames = ""
	dataFull.CategoryLink = ""
	dataFull.Data = data
	if r.Method == "GET" {
		menu(w, r)
		if len(records) == 0 {
			execute(w, "header.html", "Корзина пуста")
		} else {
			execute(w, "header.html", "Корзина")
		}
		execute(w, "show_units.html", dataFull)
		if len(records) > 0 {
			execute(w, "header.html", "Общая стоимость")
			execute(w, "basket.html", cost)
		}
	}
	if r.Method == "POST" {
		email, _ := session.Values["login"].(string)
		user, err := database.GetUserByEmail(email)
		if err != nil {
			logger.Warn(err, "Не удалось получить запись о пользователе ", email, "!")
			return
		}
		temp := time.Now().String()
		tempStrings := strings.Split(temp, ".")
		tempTime := tempStrings[0]
		tempStrings = strings.Split(tempStrings[1], " ")
		tempTime += " " + tempStrings[1] + " " + tempStrings[2]
		date := tempTime
		id, err := database.NewOrder("Выполняется", basket, user.ID, cost, date)
		if err != nil {
			logger.Warn(err, "Не удалось добавить новый заказ!")
			message := "Приносим свои извинения. Не удалось добавить новый заказ!"
			errorMessage := errortemplate.Error{Message: message, Link: "/basket"}
			execute(w, "error.html", errorMessage)
			return
		}
		logger.Info("Добавление заказа ", id, " прошло успешно.")
		if user.Phone != "00000000000" {
			adress := "https://gateway.api.sc/get/?user=79208788087&pwd=soBH5UdzEu&sadr=dupelmarket&dadr=" + user.Phone + "&text=%D0%97%D0%B0%D0%BA%D0%B0%D0%B7%20%D1%83%D1%81%D0%BF%D0%B5%D1%88%D0%BD%D0%BE%20%D0%BE%D1%84%D0%BE%D1%80%D0%BC%D0%BB%D0%B5%D0%BD.%20%D0%9D%D0%BE%D0%BC%D0%B5%D1%80%20%D0%B7%D0%B0%D0%BA%D0%B0%D0%B7%D0%B0:%20" + strconv.Itoa(id)
			http.Get(adress)
		}
		records := splitBasket(basket)
		if len(records) > 1 {
			idList := []int{}
			for _, record := range records {
				idList = append(idList, record.ID)
			}
			createPairs(idList)
		}
		session.Values["basket"] = ""
		session.Values["cost"] = cost
		session.Save(r, w)
		http.Redirect(w, r, "/payment", 302)
	}
}

func addToBasket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	session, _ := store.Get(r, "cookie-name")
	basket, _ := session.Values["basket"].(string)
	units := strings.Split(basket, ";")
	compare := false
	for i := range units {
		unitInfo := strings.Split(units[i], ":")
		if unitInfo[0] == unitIDstring {
			number, errString := strconv.Atoi(unitInfo[1])
			if errString != nil {
				logger.Warn(errString, "Не удалось конвертировать строку в число!")
				return
			}
			number++
			units[i] = unitInfo[0] + ":" + strconv.Itoa(number)
			compare = true
			break
		}
	}
	if compare {
		basket = ""
		for i := range units {
			if i == len(units)-1 { // Последняя запись всегда пустая
				break
			}
			basket += units[i] + ";"
		}
		session.Values["basket"] = basket
	} else {
		session.Values["basket"] = basket + unitIDstring + ":1;"
	}
	session.Save(r, w)
	http.Redirect(w, r, currentPage, 302)
}

func removeOneFromBasket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	session, _ := store.Get(r, "cookie-name")
	basket, _ := session.Values["basket"].(string)
	units := strings.Split(basket, ";")
	for i := range units {
		unitInfo := strings.Split(units[i], ":")
		if unitInfo[0] == unitIDstring {
			number, errString := strconv.Atoi(unitInfo[1])
			if errString != nil {
				logger.Warn(errString, "Не удалось конвертировать строку в число!")
				return
			}
			number--
			if number == 0 {
				units[i] = "removed"
				break
			}
			units[i] = unitInfo[0] + ":" + strconv.Itoa(number)
			break
		}
	}
	basket = ""
	for i := range units {
		if i == len(units)-1 { // Последняя запись всегда пустая
			break
		}
		if units[i] == "removed" {
			continue
		}
		basket += units[i] + ";"
	}
	session.Values["basket"] = basket
	session.Save(r, w)
	http.Redirect(w, r, "/basket", 302)
}

func removeFromBasket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	session, _ := store.Get(r, "cookie-name")
	basket, _ := session.Values["basket"].(string)
	units := strings.Split(basket, ";")
	for i := range units {
		unitInfo := strings.Split(units[i], ":")
		if unitInfo[0] == unitIDstring {
			units[i] = "removed"
			break
		}
	}
	basket = ""
	for i := range units {
		if i == len(units)-1 { // Последняя запись всегда пустая
			break
		}
		if units[i] == "removed" {
			continue
		}
		basket += units[i] + ";"
	}
	session.Values["basket"] = basket
	session.Save(r, w)
	http.Redirect(w, r, "/basket", 302)
}

func clearBasket(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["basket"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/basket", 302)
}
