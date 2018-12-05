package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func pairs(w http.ResponseWriter, r *http.Request) {
	pairs, err := database.GetAllPairs()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список пар товаров!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "pairs.html", pairs)
	}
}

func editPair(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pairIDstring := vars["id"]
	pairID, err := strconv.Atoi(pairIDstring)
	if err != nil {
		logger.Warn(err, "Не удалось конвертировать строку в число!")
		return
	}
	pair, err := database.GetPair(pairID)
	if err != nil {
		logger.Warn(err, "Не удалось получить данные о паре товаров ", pairID, "!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "pair.html", pair)
	}
}

func updatePair(w http.ResponseWriter, r *http.Request) {
	pairID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		logger.Warn(err, "Не удалось конвертировать строку в число!")
		return
	}
	pair, err := database.GetPair(pairID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о паре товаров ", pairID, "!")
		return
	}
	if r.Method == "POST" {
		pair.One, _ = strconv.Atoi(r.FormValue("one"))
		pair.Two, _ = strconv.Atoi(r.FormValue("two"))
		pair.Count, _ = strconv.Atoi(r.FormValue("count"))
		err := database.UpdatePair(pair.ID, pair.One, pair.Two, pair.Count)
		if err != nil {
			logger.Warn(err, "Не удалось обновить пару товаров ", pair.ID, "!")
			return
		}
		logger.Info("Пара товаров ", pair.ID, " обновлёна успешно.")
		http.Redirect(w, r, "/edit/pairs", 302)
	}
}

func addPair(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "add_pair.html", nil)
	}
	if r.Method == "POST" {
		result := database.Pair{}
		result.One, _ = strconv.Atoi(r.FormValue("one"))
		result.Two, _ = strconv.Atoi(r.FormValue("two"))
		result.Count, _ = strconv.Atoi(r.FormValue("count"))
		id, err := database.NewPair(result.One, result.Two, result.Count)
		if err != nil {
			logger.Warn(err, "Не удалось добавить новую пару товаров!")
			return
		}
		logger.Info("Добавление пары товаров ", id, " прошло успешно.")
		http.Redirect(w, r, "/edit/pairs", 302)
	}
}

func delPair(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pairID, _ := strconv.Atoi(r.FormValue("id"))
		err := database.DelPair(pairID)
		if err != nil {
			logger.Warn(err, "Не удалось удалить запись о паре товаров ", pairID, "!")
			return
		}
		logger.Info("Удаление записи о паре товаров ", pairID, " прошло успешно.")
		http.Redirect(w, r, "/edit/pairs", 302)
	}
}

func delAllPairs(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := database.DelAllPairs()
		if err != nil {
			logger.Warn(err, "Не удалось удалить все пары товаров!")
			return
		}
		logger.Info("Удаление всех пар товаров прошло успешно.")
		http.Redirect(w, r, "/edit/pairs", 302)
	}
}
