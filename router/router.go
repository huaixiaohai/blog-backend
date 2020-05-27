package router

import (
	"github.com/julienschmidt/httprouter"
	"my/blog-backend/controller"
	"my/blog-backend/lib/log"
	"net/http"
)

var Router = httprouter.New()

func Init() {
	// test
	get("/dev-api", controller.Index)
	get("/dev-api/index", controller.Index)

	// user
	post("/api/user/login", controller.Login)
	get("/api/user/info", controller.GetUserInfo)

	// article_author
	get("/api/search/author/list", controller.GetAuthorList)

	//article
	post("/api/article/create", controller.AddArticle)
	post("/api/article/update", controller.UpdateArticle)
	delete("/api/article/delete", controller.DeleteArticle)
	get("/api/article/list", controller.ArticleList)
	get("/api/article/detail", controller.ArticleDetail)

	log.Info("路由初始化完成")
}

func get(path string, handles ...httprouter.Handle) {
	Router.GET(path, wrapper(handles...))
}

func post(path string, handles ...httprouter.Handle) {
	Router.POST(path, wrapper(handles...))
}

func put(path string, handles ...httprouter.Handle) {
	Router.PUT(path, wrapper(handles...))
}

func delete(path string, handles ...httprouter.Handle) {
	Router.DELETE(path, wrapper(handles...))
}

func wrapper(handles ...httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		for _, handle := range handles {
			handle(writer, request, params)
		}
	}
}
