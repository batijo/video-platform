package controllers

import (
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("admin panel"))
}
