package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/errortemplate"
)

func orders(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/edit"}
	execute(w, "error.html", errorMessage)
}
