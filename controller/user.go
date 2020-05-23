package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"my/blog-backend/dao"
	"my/blog-backend/lib/log"
	"my/blog-backend/model"
	"net/http"
	"time"
)

const (
	SecretKey = "1885df74d00dbbe19274c6d955feeb5b"
)

func Login(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
	}
	user := &model.User{}
	err = json.Unmarshal(bytes, user)
	if err != nil {
		log.Error(err)
	}

	var dbUser *model.User
	dbUser, err = dao.User.One(user.UserName)
	if err != nil {
		log.Error(err)
	}
	if dbUser.Password != user.Password {
		log.Info("登录失败")
		resp.WriteHeader(http.StatusForbidden)
		return
	}
	var tokenStr string
	tokenStr, err = genTokenStr()
	if err != nil {
		resp.WriteHeader(http.StatusForbidden)
		return
	}
	log.Info("登录成功")
	var buf []byte
	buf, err = json.Marshal(map[string]interface{}{"data": map[string]interface{}{"token": tokenStr}, "code": 20000})
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.WriteHeader(http.StatusOK)

	resp.Write(buf)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	rToken := query["token"][0]
	token, err := genTokenStr()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rToken != token {
		log.Error("验证不通过")
		//w.WriteHeader(http.StatusInternalServerError)
		//return
	}

	var buf []byte
	buf, err = json.Marshal(map[string]interface{}{"data": map[string]interface{}{
		"roles":        []string{"admin"},
		"introduction": "I am a super administrator",
		"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		"name":         "Super Admin",
	}, "code": 20000})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func genTokenStr() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		log.ErrorWithFields("Error while signing the token", log.Fields{
			"err": err,
		})
		return "", err
	}
	return tokenString, nil
}
