package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/content"
)

func addFeatures(w http.ResponseWriter, r *http.Request) {
	content.AddFeaturesFromFile("./lists/features.txt")
	http.Redirect(w, r, "/edit/features", 302)
}

func addProcessor(w http.ResponseWriter, r *http.Request) {
	content.AddProcessorsFromFile("./lists/processors.txt", filesFolder)
	http.Redirect(w, r, "/edit/units", 302)
}

func addMotherboard(w http.ResponseWriter, r *http.Request) {
	content.AddMotherboardsFromFile("./lists/motherboards.txt", filesFolder)
	http.Redirect(w, r, "/edit/units", 302)
}

func addVideocard(w http.ResponseWriter, r *http.Request) {
	content.AddVideocardsFromFile("./lists/videocards.txt", filesFolder)
	http.Redirect(w, r, "/edit/units", 302)
}

func addRAM(w http.ResponseWriter, r *http.Request) {
	access := hidePage(w, r)
	if access == false {
		return
	}
	content.AddRAMFromFile("./lists/rams.txt", filesFolder)
	http.Redirect(w, r, "/edit/units", 302)
}
