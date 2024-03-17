package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CrudItem struct {
	Id          int
	Name        string
	Description string
	internal    string
}

func main() {
	r := chi.NewRouter()
	currentId := 1
	storage := make(map[int]CrudItem)
	r.Get("/crud-items/", func(w http.ResponseWriter, r *http.Request) {
		result := make([]CrudItem, 0, len(storage))
		for _, item := range storage {
			result = append(result, item)
		}
		resultJson, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}
		w.Write(resultJson)
	})
	r.Get("/crud-items/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if _, ok := storage[id]; !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		resultJson, err := json.Marshal(storage[id])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}
		w.Write(resultJson)
	})
	r.Post("/crud-items/", func(w http.ResponseWriter, r *http.Request) {
		var item CrudItem
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		item.Id = currentId
		storage[currentId] = item
		jsonItem, err := json.Marshal(item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(jsonItem)
		currentId += 1
	})

	r.Put("/crud-items/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

			return
		}
		if _, ok := storage[id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var item CrudItem
		err = json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		item.Id = id
		storage[id] = item
		jsonItem, err := json.Marshal(item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

		}
		w.Write(jsonItem)
	})
	r.Delete("/crud-items/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))

			return
		}
		if _, ok := storage[id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		delete(storage, id)
	})
	http.ListenAndServe(":3000", r)
}
