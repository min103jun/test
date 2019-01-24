package controllers

import (
	"net/http"
)

func (app *Application) RequestHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
	}
	if r.FormValue("submitted") == "true" {
		//txid, err := app.Fabric.InsertUser("min103jun", "kimminjun", "alswns1031", "941031-1111111", "Daegu")
		txid, err := app.Fabric.InsertVote("vote02", "2017/12/31", "2018/12/31", "question01", "question02", "question03", "question04")
		//txid, err := app.Fabric.InsertVoteResult("vote01", "warhand9", "1")
		if err != nil {
			http.Error(w, "Unable to invoke hello in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "request.html", data)
}
