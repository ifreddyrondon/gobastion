package todo

import (
	"net/http"
	"strconv"

	"fmt"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
)

type Handler struct {
	gobastion.Reader
	gobastion.Responder
}

// Routes creates a REST router for the todos resource
func (h *Handler) Routes() chi.Router {
	r := gobastion.NewRouter()

	r.Get("/", h.List)    // GET /todos - read a list of todos
	r.Post("/", h.Create) // POST /todos - create a new todo and persist it
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", h.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", h.Delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	todo1 := todo{Description: "do something 1"}
	todo2 := todo{Description: "do something 2"}

	h.Send(w, []todo{todo1, todo2})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	todo1 := todo{Description: "do something 1"}
	h.Created(w, todo1)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	todo1 := todo{Id: i, Description: fmt.Sprintf("do something %v", id)}
	h.Send(w, todo1)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(id) // the error should be handle
	todo1 := todo{Id: i, Description: fmt.Sprintf("do something %v", id)}
	h.Send(w, todo1)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.NoContent(w)
}