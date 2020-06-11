package controller

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"my/blog-backend/dao"
	"my/blog-backend/lib/log"
	"my/blog-backend/lib/redis"
	"my/blog-backend/model"
	"net/http"
	"strconv"
)

const RKBlogPV = "redis_key_blog_pv_"

// 请求参数
type ArticleDataReq struct {
	ID       int64
	Title    string
	Summary  string
	Content  string
	Status   model.ArticleStatus
	AuthorID int64
	Image    string
}

func AddArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s := &ArticleDataReq{}
	err = json.Unmarshal(buf, s)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dao.ArticleInfo.Insert(s.Title, s.Summary, s.Image, s.Content, s.AuthorID, s.Status)
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

func UpdateArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info(string(buf))
	s := &ArticleDataReq{}
	err = json.Unmarshal(buf, s)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dao.ArticleInfo.Update(s.ID, s.Title, s.Summary, s.Image, s.Content, s.AuthorID, s.Status)
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

func DeleteArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

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
		items = append(items, map[string]interface{}{
			"id":        one.ID,
			"title":     one.Title,
			"status":    one.Status,
			"image":     one.Image,
			"summary":   one.Summary,
			"authorId":  one.AuthorID,
			"timestamp": one.CreateTime.Unix()})
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
		"title":     info.Title,
		"status":    info.Status,
		"image":     info.Image,
		"summary":   info.Summary,
		"authorId":  info.AuthorID,
		"timestamp": info.CreateTime.Unix(),
		"content":   content.Content,
	}, "code": 20000})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
	redisKey := RKBlogPV + strconv.FormatInt(id, 10)
	//if !redis.EXISTS(redisKey) {
	//	redis.SET(redisKey, 0)
	//}
	redis.INCR(redisKey)
}
