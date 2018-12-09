package handlefunc

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func orders(w http.ResponseWriter, r *http.Request) {
	orders, err := database.GetAllOrders()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список заказов!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "orders.html", orders)
	}
}

func editOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDstring := vars["id"]
	orderID, errString := strconv.Atoi(orderIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
		return
	}
	order, err := database.GetOrder(orderID)
	if err != nil {
		logger.Warn(err, "Не удалось получить данные о заказе ", orderID, "!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "order.html", order)
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	orderID, errString := strconv.Atoi(r.FormValue("id"))
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
		return
	}
	order, err := database.GetOrder(orderID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о заказе ", orderID, "!")
		return
	}
	if r.Method == "POST" {
		if r.FormValue("state_1") != "" {
			order.State = "Выполняется"
		}
		if r.FormValue("state_2") != "" {
			order.State = "Исполнен"
		}
		if r.FormValue("state_3") != "" {
			order.State = "Отменён"
		}
		order.Units = r.FormValue("units")
		order.UserID, _ = strconv.Atoi(r.FormValue("user_id"))
		order.Cost, _ = strconv.Atoi(r.FormValue("cost"))
		order.Date = r.FormValue("date")
		err := database.UpdateOrder(order.ID, order.State, order.Units, order.UserID, order.Cost, order.Date)
		if err != nil {
			logger.Warn(err, "Не удалось обновить заказ ", order.ID, "!")
			return
		}
		logger.Info("Заказ ", order.ID, " обновлён успешно.")
		http.Redirect(w, r, "/edit/orders", 302)
	}
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "add_order.html", nil)
	}
	if r.Method == "POST" {
		result := database.Order{}
		if r.FormValue("state_1") != "" {
			result.State = "Выполняется"
		}
		if r.FormValue("state_2") != "" {
			result.State = "Исполнен"
		}
		if r.FormValue("state_3") != "" {
			result.State = "Отменён"
		}
		result.Units = r.FormValue("units")
		result.UserID, _ = strconv.Atoi(r.FormValue("user_id"))
		result.Cost, _ = strconv.Atoi(r.FormValue("cost"))
		temp := time.Now().String()
		tempStrings := strings.Split(temp, ".")
		tempTime := tempStrings[0]
		tempStrings = strings.Split(tempStrings[1], " ")
		tempTime += " " + tempStrings[1] + " " + tempStrings[2]
		result.Date = tempTime
		id, err := database.NewOrder(result.State, result.Units, result.UserID, result.Cost, result.Date)
		if err != nil {
			logger.Warn(err, "Не удалось добавить новый заказ!")
			return
		}
		logger.Info("Добавление заказа ", id, " прошло успешно.")
		http.Redirect(w, r, "/edit/orders", 302)
	}
}

func delOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		orderID, _ := strconv.Atoi(r.FormValue("id"))
		err := database.DelOrder(orderID)
		if err != nil {
			logger.Warn(err, "Не удалось удалить запись о заказе ", orderID, "!")
			return
		}
		logger.Info("Удаление записи о заказе ", orderID, " прошло успешно.")
		http.Redirect(w, r, "/edit/orders", 302)
	}
}

func delAllOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := database.DelAllOrders()
		if err != nil {
			logger.Warn(err, "Не удалось удалить все заказы!")
			return
		}
		logger.Info("Удаление всех заказов прошло успешно.")
		http.Redirect(w, r, "/edit/orders", 302)
	}
}
