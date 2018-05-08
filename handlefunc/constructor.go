package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func showConstructor(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	processorID, _ := session.Values["processor"].(int)
	processor, _ := database.GetUnit(processorID)

	motherboardID, _ := session.Values["motherboard"].(int)
	motherboard, _ := database.GetUnit(motherboardID)

	videocardID, _ := session.Values["videocard"].(int)
	videocard, _ := database.GetUnit(videocardID)

	ramID, _ := session.Values["ram"].(int)
	ram, _ := database.GetUnit(ramID)

	nothingProcessorSelected := makeNothingData("../../assets/images/nothing_processor.jpg", "Процессор не выбран")
	nothingMotherboardSelected := makeNothingData("../../assets/images/nothing_motherboard.png", "Материнская плата не выбрана")
	nothingVideocardSelected := makeNothingData("../../assets/images/nothing_videocard.jpg", "Видеокарта не выбрана")
	nothingRamSelected := makeNothingData("../../assets/images/nothing_ram.jpg", "Оперативная память не выбрана")

	processorSelected := makeSingleData(processor.ID, "processors", "remove_constructor", "Убрать товар")
	motherboardSelected := makeSingleData(motherboard.ID, "motherboards", "remove_constructor", "Убрать товар")
	videocardSelected := makeSingleData(videocard.ID, "videocards", "remove_constructor", "Убрать товар")
	ramSelected := makeSingleData(ram.ID, "rams", "remove_constructor", "Убрать товар")

	processors := makeData(true, "Процессоры", "Процессор", "processors", "add_constructor", "Выбрать товар")
	motherboards := makeData(true, "Материнские платы", "Материнская плата", "motherboards", "add_constructor", "Выбрать товар")
	videocards := makeData(true, "Видеокарты", "Видеокарта", "videocards", "add_constructor", "Выбрать товар")
	rams := makeData(true, "Оперативная память", "Оперативная память", "rams", "add_constructor", "Выбрать товар")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", "Выбранная конфигурация")
		if processor.Name != "" {
			execute(w, "show_unit.html", processorSelected)
		} else {
			execute(w, "show_nothing.html", nothingProcessorSelected)
		}
		if motherboard.Name != "" {
			execute(w, "show_unit.html", motherboardSelected)
		} else {
			execute(w, "show_nothing.html", nothingMotherboardSelected)
		}
		if videocard.Name != "" {
			execute(w, "show_unit.html", videocardSelected)
		} else {
			execute(w, "show_nothing.html", nothingVideocardSelected)
		}
		if ram.Name != "" {
			execute(w, "show_unit.html", ramSelected)
		} else {
			execute(w, "show_nothing.html", nothingRamSelected)
		}
		execute(w, "two_buttons.html", nil)
		execute(w, "header.html", "Доступные компоненты")
		execute(w, "show_units.html", processors)
		execute(w, "show_units.html", motherboards)
		execute(w, "show_units.html", videocards)
		execute(w, "show_units.html", rams)
	}
}

func addConstructor(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	unitID, errString := strconv.Atoi(unitIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
	}
	unit, err := database.GetUnit(unitID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о товаре ", unitID, "!")
	} else {
		logger.Info("Данные о товаре ", unitID, " получены успешно.")
	}
	category, _ := database.GetCategory(unit.CategoryID)
	if category.Name == "Процессор" {
		session.Values["processor"] = unit.ID
	}
	if category.Name == "Материнская плата" {
		session.Values["motherboard"] = unit.ID
	}
	if category.Name == "Видеокарта" {
		session.Values["videocard"] = unit.ID
	}
	if category.Name == "Оперативная память" {
		session.Values["ram"] = unit.ID
	}
	session.Save(r, w)
	http.Redirect(w, r, "/constructor", 302)
}

func removeConstructor(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	unitID, errString := strconv.Atoi(unitIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
	}
	unit, err := database.GetUnit(unitID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о товаре ", unitID, "!")
	} else {
		logger.Info("Данные о товаре ", unitID, " получены успешно.")
	}
	category, _ := database.GetCategory(unit.CategoryID)
	if category.Name == "Процессор" {
		session.Values["processor"] = ""
	}
	if category.Name == "Материнская плата" {
		session.Values["motherboard"] = ""
	}
	if category.Name == "Видеокарта" {
		session.Values["videocard"] = ""
	}
	if category.Name == "Оперативная память" {
		session.Values["ram"] = ""
	}
	session.Save(r, w)
	http.Redirect(w, r, "/constructor", 302)
}

func clearConstructor(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["processor"] = ""
	session.Values["motherboard"] = ""
	session.Values["videocard"] = ""
	session.Values["ram"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/constructor", 302)
}

func orderConstructor(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/constructor"}
	execute(w, "error.html", errorMessage)
}
