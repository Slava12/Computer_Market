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

func showBasket(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	basket, _ := session.Values["basket"].(string)
	units := strings.Split(basket, ";")
	data := make([]Data, len(units)-1)
	cost := 0
	for i := range units {
		if i == len(units)-1 { // Последняя запись всегда пустая
			break
		}
		unitInfo := strings.Split(units[i], ":")
		unitID, errString := strconv.Atoi(unitInfo[0])
		if errString != nil {
			logger.Warn(errString, "Не удалось конвертировать строку в число!")
			return
		}
		unit, err := database.GetUnit(unitID)
		if err != nil {
			logger.Warn(err, "Не удалось получить запись о товаре ", unitID, "!")
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
		unitNumber, errString := strconv.Atoi(unitInfo[1])
		if errString != nil {
			logger.Warn(errString, "Не удалось конвертировать строку в число!")
			return
		}
		cost += unit.Price * unitNumber
	}
	dataFull := DataFull{}
	dataFull.ShowCategory = false
	dataFull.CategoryNames = ""
	dataFull.CategoryLink = ""
	dataFull.Data = data
	if r.Method == "GET" {
		menu(w, r)
		if len(units) == 1 {
			execute(w, "header.html", "Корзина пуста")
		} else {
			execute(w, "header.html", "Корзина")
		}
		err := tpl.ExecuteTemplate(w, "show_units.html", dataFull)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		if len(units) != 1 {
			execute(w, "header.html", "Общая стоимость")
			err = tpl.ExecuteTemplate(w, "basket.html", cost)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
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
		session.Values["basket"] = ""
		session.Save(r, w)
		http.Redirect(w, r, "/orders", 302)
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
	http.Redirect(w, r, "/index", 302)
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
