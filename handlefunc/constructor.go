package handlefunc

import (
	"net/http"
)

func showConstructor(w http.ResponseWriter, r *http.Request) {
	processors := makeData("Процессор", "processors", "add_constructor", "Выбрать товар")
	motherboards := makeData("Материнская плата", "motherboards", "add_constructor", "Выбрать товар")
	videocards := makeData("Видеокарта", "videocards", "add_constructor", "Выбрать товар")
	rams := makeData("Оперативная память", "rams", "add_constructor", "Выбрать товар")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", processors)
		execute(w, "show_units.html", motherboards)
		execute(w, "show_units.html", videocards)
		execute(w, "show_units.html", rams)
	}
}
