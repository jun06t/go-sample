package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/openfga/go-sdk/client"
)

const (
	apiURL = "http://127.0.0.1:8080"
	storeID = "01JJ7TD9J1T5H03V6V4H0PH37P"
)

func main() {
	r := chi.NewRouter()

	// 簡単のため、本来リクエストから取得すべきJWTオブジェクトを埋め込む
	r.Use(embedJWT)

	r.Route("/articles", func(r chi.Router) {
		r.Use(preauthorize("article"))
		r.Use(checkAuthorization)

		r.Get("/", list)
		r.Post("/", create)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", get)
			r.Put("/", update)
			r.Delete("/", del)
		})
	})

	r.Route("/items", func(r chi.Router) {
		r.Use(preauthorize("item"))
		r.Use(checkAuthorization)

		r.Get("/", list)
		r.Post("/", create)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", get)
			r.Put("/", update)
			r.Delete("/", del)
		})
	})

	err := http.ListenAndServe(":9000", r)
	if err != nil {
		fmt.Println(err)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "got %s\n", id)
}

func list(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "listed\n")
}

func create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "created\n")
}

func update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "updated %s\n", id)
}

func del(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "deleted %s\n", id)
}

func embedJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userName := r.Header.Get("X-User")
		claims := jwt.MapClaims{
			"sub": userName,
			"iss": "example.com",
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		ctx := context.WithValue(r.Context(), "jwt", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func preauthorize(resource string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value("jwt").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			user := claims["sub"].(string)
			ctx := context.WithValue(r.Context(), "user", user)

			var relation string
			switch r.Method {
			case "GET":
				relation = "viewer"
			case "POST", "PUT":
				relation = "editor"
			case "DELETE":
				relation = "owner"
			default:
				relation = "owner"
			}

			ctx = context.WithValue(ctx, "relation", relation)
			object := "api:" + resource
			ctx = context.WithValue(ctx, "object", object)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func checkAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fgaClient, err := NewSdkClient(&ClientConfiguration{
			ApiUrl:  apiURL,
			StoreId: storeID,
		})
		if err != nil {
			log.Printf("%#v\n", err)
			http.Error(w, "Unable to build OpenFGA client", http.StatusServiceUnavailable)
			return
		}

		user := r.Context().Value("user").(string)
		relation := r.Context().Value("relation").(string)
		object := r.Context().Value("object").(string)

		body := ClientCheckRequest{
			User:     "user:" + user,
			Relation: relation,
			Object:   object,
		}

		data, err := fgaClient.Check(context.Background()).Body(body).Execute()
		if err != nil {
			log.Printf("%#v\n", err)
			http.Error(w, "Unable to check for authorization", http.StatusServiceUnavailable)
			return
		}

		if !(*data.Allowed) {
			log.Printf("%#v\n", data)
			http.Error(w, "Forbidden to access document", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
