package main

import (
	"too-simple-twitter/server/controller"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

// SPAのパスプレフィックス
const APP_PREFIX = "/app"

// APIのパスプレフィックス
const API_PREFIX = "/api/v1"

// SPAの静的ファイルが出力されているディレクトリ
const STATIC_PATH = "/workspace/client/out"

// SPAのindex.html
const INDEX_PATH = "index.html"

func main() {
	port := getPort()
	fmt.Println("HTTP Server is running on localhost" + port)

	// HTML
	router := mux.NewRouter()
	router.PathPrefix(APP_PREFIX).HandlerFunc(servePages)

	// API
	apiRouter := router.PathPrefix(API_PREFIX).Subrouter()
	apiRouter.Use(controller.CORSMiddleware) // 開発用CORS許可
	controller.UseAuthRouter(apiRouter)
	controller.UseUserRouter(apiRouter)
	controller.UseTweetRouter(apiRouter)

	// CORSプリフライトリクエスト
	router.PathPrefix("/").HandlerFunc(controller.CORSPreflight).Methods(http.MethodOptions)

	// Static files
	router.PathPrefix("/").HandlerFunc(serveStaticFiles)

	http.Handle("/", router)
	http.ListenAndServe(port, nil)
}

/*
  SSGしたHTMLをホスティングする
*/
func servePages(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch path {
	case "/app":
		http.ServeFile(w, r, filepath.Join(STATIC_PATH, "index.html"))
	case "/app/home":
		http.ServeFile(w, r, filepath.Join(STATIC_PATH, "home.html"))
	case "/app/signup":
		http.ServeFile(w, r, filepath.Join(STATIC_PATH, "signup.html"))

	}
}

/*
  静的ファイルをホスティングする
  参考：https://github.com/gorilla/mux#serving-single-page-applications
*/
func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// "/"のときはアプリのパスにリダイレクトする
	if path == "/" {
		http.Redirect(w, r, APP_PREFIX, http.StatusMovedPermanently)
	}

	path = filepath.Join(STATIC_PATH, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(STATIC_PATH)).ServeHTTP(w, r)
}

// 環境変数を見てサーバ稼働ポートを取得
func getPort() string {
	if string(os.Getenv("DEV")) == "true" {
		return ":80"
	} else {
		return ":5000"
	}
}
