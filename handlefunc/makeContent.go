package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/content"
)

func addFeatures(w http.ResponseWriter, r *http.Request) {
	content.AddFeaturesFromFile("./lists/features.txt")
	http.Redirect(w, r, "/edit/features", 302)
}
