package handlefunc

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"
)

type paymentData struct {
	OutSum         int
	SignatureValue string
}

func pay(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	cost, _ := session.Values["cost"].(int)
	str := "compmarket:" + strconv.Itoa(cost) + ":0:AQhJe1jw2GSepy9nX04e"
	hash := md5.Sum([]byte(str))
	data := paymentData{cost, hex.EncodeToString(hash[:])}
	if r.Method == "GET" {
		execute(w, "payment.html", data)
	}
}

func success(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		execute(w, "success.html", nil)
	}
}

func fail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		execute(w, "fail.html", nil)
	}
}
