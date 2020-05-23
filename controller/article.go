package controller

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"my/blog-backend/dao"
	"my/blog-backend/lib/log"
	"my/blog-backend/model"
	"net/http"
	"strconv"
	"time"
)

func AddArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	type data struct {
		Title    string
		Summary  string
		Content  string
		Status   int32
		AuthorID int64
	}
	s := &data{}
	err = json.Unmarshal(buf, s)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dao.ArticleInfo.Insert(s.Title, s.Summary, "", s.Content, s.AuthorID, model.ArticleStatus(s.Status), time.Now())
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//var buf []byte
	buf, err = json.Marshal(map[string]interface{}{"code": 20000})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	w.Write(buf)
	return
}

func ArticleList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	args := r.URL.Query()
	var list []*model.ArticleInfo
	var total int64
	var err error
	var page int64
	var limit int64
	page, err = strconv.ParseInt(args.Get("page"), 10, 64)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	limit, err = strconv.ParseInt(args.Get("limit"), 10, 64)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	total, list, err = dao.ArticleInfo.List(page, limit)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := make([]interface{}, 0)
	for _, one := range list {
		items = append(items, map[string]interface{}{"id": one.ID, "timestamp": one.PublishedTime.Unix(), "author_id": one.AuthorID, "status": one.Status, "title": one.Title, "summary": one.Summary})
	}
	var buf []byte
	buf, err = json.Marshal(map[string]interface{}{"data": map[string]interface{}{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	}, "code": 20000})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)

}

func ArticleDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	args := r.URL.Query()
	id, err := strconv.ParseInt(args.Get("id"), 10, 64)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var info *model.ArticleInfo
	var content *model.ArticleContent
	info, content, err = dao.ArticleInfo.Detail(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var buf []byte
	buf, err = json.Marshal(map[string]interface{}{"data": map[string]interface{}{
		"id":        info.ID,
		"author_id": info.AuthorID,
		"title":     info.Title,
		"summary":   info.Summary,
		"content":   content.Content,
		//"published_time": content.pulished_ti
	}, "code": 20000})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
