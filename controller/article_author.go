package controller

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"my/blog-backend/dao"
	"my/blog-backend/lib/log"
	"net/http"
)

func GetAuthorList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	authors, err := dao.ArticleAuthor.List()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := make([]interface{}, 0)
	for _, author := range authors {
		items = append(items, map[string]interface{}{"id": author.ID, "name": author.Name})
	}

	var buf []byte
	buf, err = json.Marshal(map[string]interface{}{"data": map[string]interface{}{
		"items": items,
	}, "code": 20000})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)

}
