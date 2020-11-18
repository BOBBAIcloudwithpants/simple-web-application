package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/bobbaicloudwithpants/simple-web-application/models"
	"html/template"
	"net/http"
	m "github.com/gorilla/mux"
	"sync"
)
var mutex sync.Mutex

func InitHandlers() *m.Router{
	mux := m.NewRouter()

	// 处理所有的对于页面，css, js等GET操作的处理
	mux.PathPrefix("/views/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			path := "." + request.URL.Path
			fmt.Println(path)
			http.ServeFile(writer, request, path)
		}
	})

	// 根目录，返回用户登录的页面
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("method: ", request.Method)
		request.ParseForm()
		if request.Method == "GET" {

			// 渲染模版
			t, _ := template.ParseFiles("./views/login.html")
			t.Execute(writer, nil)

		} else {
			// 将用户通过表单发送过来的字段存放到内存当中
			// 加锁，防止多个go程同时修改
			mutex.Lock()
			u := models.User{UserId: request.Form["user_id"][0], Username: request.Form["username"][0], Password: request.Form["password"][0]}
			models.Users[u.UserId] = u
			mutex.Unlock()

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)
			json.NewEncoder(writer).Encode(u)
		}
	})

	// 用户页面的处理
	mux.HandleFunc("/users/{user_id}", func(writer http.ResponseWriter, request *http.Request) {
		// 解析url中的参数
		vars := m.Vars(request)

		// 渲染模版
		t, _ := template.ParseFiles("./views/info.html")
		u := models.Users[vars["user_id"]]
		user := models.User{Username: u.Username, Password: u.Password, UserId: u.UserId}
		t.Execute(writer, user)

	})
	return mux
}

