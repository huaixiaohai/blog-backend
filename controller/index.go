package controller

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Index(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	resp.Write([]byte("index" + req.URL.Path))
}
