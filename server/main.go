package main

import (
	"github-crawler/database"
	"fmt"
	"net/http"
	"encoding/json"
	"os"
)

var redisClient *database.Redis

type Response struct {
	State		int `json:"state"`
	ErrorMsg	string `json:"error_msg"`
	Result 		interface{} `json:"result"`
}

type Server struct {

}

func main() {

	redisHost := os.Getenv("RedisHost")
	redisPassword := os.Getenv("RedisPassword")
	if redisHost == "" {
		redisHost = "127.0.0.1:6379"
	}
	redisClient = database.NewRedis(redisHost, redisPassword, 0)

	if err := redisClient.Connection(); err != nil {
		fmt.Println("Redis connect error: ", err)
		return
	}

	defer redisClient.Disconnect()

	fmt.Println("start server http://localhost:8080 ")
	err := http.ListenAndServe(":8080", &Server{})
	if err != nil {
		fmt.Println("start server error: ", err)
	}
}

func (_ *Server) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	switch r.URL.Path {
	case "/language", "/language/":
		AllLanguageHandler(w, r)
	case "/repo", "/repo/":
		TrendingRepoHandler(w, r)
	case "/developer", "/developer/":
		TrendingDeveloperHandler(w, r)
	case "/":
		response := Response{
			State: 1,
			Result: map[string]string {
				"language" : "/language",
				"repo" : "/repo",
				"developer" : "/developer",
			},
		}
		WriteResponse(w, response)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func AllLanguageHandler(w http.ResponseWriter, r *http.Request)  {
	cacheLanguages := redisClient.GetCacheLanguages()
	response := Response{
		State: 1,
	}
	if len(cacheLanguages) >= 0 {
		response.Result = cacheLanguages
	}
	WriteResponse(w, response)
}

func TrendingRepoHandler(w http.ResponseWriter, r *http.Request)  {
	lan := r.Form.Get("lan")
	since := r.Form.Get("since")
	if lan == "all-language" || lan == "" {
		lan = "all-language"
	}
	if since == "" {
		since = "daily"
	}

	key := fmt.Sprintf("trending-repo-%v-%v", lan, since)
	str, err := redisClient.Get(key)
	response := Response{
		State: 1,
	}
	if err != nil {
		fmt.Println("repo error: ", err)
		response.ErrorMsg = err.Error()
	} else {
		var val interface{}
		err = json.Unmarshal([]byte(str), &val)
		if err != nil {
			response.ErrorMsg = err.Error()
		} else {
			response.Result = val
		}
	}

	WriteResponse(w, response)
}

func TrendingDeveloperHandler(w http.ResponseWriter, r *http.Request)  {
	lan := r.Form.Get("lan")
	since := r.Form.Get("since")
	if lan == "all-language" || lan == "" {
		lan = "all-language"
	}
	if since == "" {
		since = "daily"
	}

	key := fmt.Sprintf("trending-developer-%v-%v", lan, since)
	str, err := redisClient.Get(key)
	response := Response{
		State: 1,
	}
	if err != nil {
		fmt.Println("repo error: ", err)
		response.ErrorMsg = err.Error()
	} else {
		var val interface{}
		err = json.Unmarshal([]byte(str), &val)
		if err != nil {
			response.ErrorMsg = err.Error()
		} else {
			response.Result = val
		}
	}

	WriteResponse(w, response)
}

func WriteResponse(w http.ResponseWriter, response Response)  {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(`{"state" : 1, "error_msg" : "", "result" : null}`))
	} else {
		w.Write(data)
	}

}

