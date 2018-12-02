package handlefunc

import (
	"net/http"
	"time"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/logger"
)

func showOrders(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	email, _ := session.Values["login"].(string)
	user, err := database.GetUserByEmail(email)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о пользователе ", email, "!")
		message := "Произошла ошибка, информация о пользователе не найдена."
		errorMessage := errortemplate.Error{Message: message, Link: "/index"}
		execute(w, "error.html", errorMessage)
		return
	}
	orders, err := database.GetOrdersByUserID(user.ID)
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список заказов пользователя ", user.ID, "!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", "Ваши заказы")
		execute(w, "show_orders.html", orders)
	}
}

func showOrder(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/orders"}
	execute(w, "error.html", errorMessage)
}

func makeOrder(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	email, _ := session.Values["login"].(string)
	user, err := database.GetUserByEmail(email)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о пользователе ", email, "!")
		return
	}
	basket, _ := session.Values["basket"].(string)
	id, err := database.NewOrder("Выполняется", basket, user.ID, 0, time.Now().String())
	if err != nil {
		logger.Warn(err, "Не удалось добавить новый заказ!")
		return
	}
	logger.Info("Добавление заказа ", id, " прошло успешно.")
	http.Redirect(w, r, "/orders", 302)
}
