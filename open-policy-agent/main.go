package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

const (
	opaEndpoint = "http://localhost:8181/v1/data/app/rbac/allow"
)

func main() {
	r := chi.NewRouter()
	r.Use(opaMiddleware)
	r.Route("/articles", func(r chi.Router) {
		r.Get("/", listArticles)
		r.Post("/", createArticle)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getArticle)
			r.Put("/", updateArticle)
			r.Delete("/", deleteArticle)
		})
	})
	http.ListenAndServe(":3000", r)
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("got %s\n", id))
}

func listArticles(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "listed\n")
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "created\n")
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("updated %s\n", id))
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("deleted %s\n", id))
}

type data struct {
	Input input `json:"input"`
}
type input struct {
	Method string   `json:"method"`
	Path   []string `json:"path"`
	Roles  []string `json:"roles"`
}

type result struct {
	Result bool `json:"result"`
}

func opaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		trimed := strings.Trim(r.URL.Path, "/")
		p := strings.Split(trimed, "/")
		// NOTE: コンテキストからユーザ情報の取得→ロールの設定
		userRoles := []string{"article.editor"}
		input := data{
			Input: input{
				Method: r.Method,
				Path:   p,
				Roles:  userRoles,
			},
		}
		body, _ := json.Marshal(input)
		resp, err := http.Post(opaEndpoint, "application/json", bytes.NewBuffer(body))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		res := result{}
		err = json.NewDecoder(resp.Body).Decode(&res)
		if !res.Result {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
