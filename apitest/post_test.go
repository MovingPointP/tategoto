package apitest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"tategoto/config/msg"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	//router取得 table初期化
	r := NewRouter()
	signup(t, r)
	beforeLoginPost_303(t, r)
	login(t, r)
	successPost_200(t, r)
}

// サインアップ
func signup(t *testing.T, r *gin.Engine) {
	requestJson := `{ "name": "hogeman", "mail": "hoge@mail.com", "password": "hogehoge"}`
	body := bytes.NewBuffer([]byte(requestJson))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/signup", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
}

// ログインなしポスト
func beforeLoginPost_303(t *testing.T, r *gin.Engine) {
	requestJson := `{ "content": "hello", "user_id": "1"}`
	body := bytes.NewBuffer([]byte(requestJson))

	//ログインが必要な処理
	responseJson := `{ "message":"` + msg.ShouldLoginErr + `", "path":"/api/posts"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/posts", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.JSONEq(t, w.Body.String(), responseJson)
	assert.Equal(t, w.Code, 303)
}

// ログイン
func login(t *testing.T, r *gin.Engine) {
	requestJson := `{ "mail": "hoge@mail.com", "password": "hogehoge"}`
	body := bytes.NewBuffer([]byte(requestJson))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/login", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
}

// 正常なポスト
func successPost_200(t *testing.T, r *gin.Engine) {
	requestJson := `{ "content": "hello", "user_id": "1"}`
	body := bytes.NewBuffer([]byte(requestJson))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/posts", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "test",
	})
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
}