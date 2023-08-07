package course

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jnka9755/go-05META/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Get    Controller
		Delete Controller
		Update Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	UpdateReq struct {
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(b Business, config Config) Endpoints {

	return Endpoints{
		Create: makeCreateEndpoint(b),
		GetAll: makeGetAllEndpoint(b, config),
		Get:    makeGetEndpoint(b),
		Delete: makeDeleteEndpoint(b),
		Update: makeUpdateEndpoint(b),
	}
}

func makeCreateEndpoint(b Business) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var request CreateReq

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		if request.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "name is required"})
			return
		}

		if request.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "start_date is required"})
			return
		}

		if request.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "end_date is required"})
			return
		}

		responseCourse, err := b.Create(&request)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 201, Data: responseCourse})
	}
}

func makeGetAllEndpoint(b Business, config Config) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		value := r.URL.Query()

		filters := Filters{
			Name: value.Get("name"),
		}

		limit, _ := strconv.Atoi(value.Get("limit"))
		page, _ := strconv.Atoi(value.Get("page"))

		count, err := b.Count(filters)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		meta, err := meta.New(page, limit, count, config.LimPageDef)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		courses, err := b.GetAll(filters, meta.Offset(), meta.Limit())

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: courses, Meta: meta})
	}
}

func makeGetEndpoint(b Business) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		course, err := b.Get(id)

		if err != nil {

			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: course})
	}
}

func makeDeleteEndpoint(b Business) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		if err := b.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Err: "Course doesn't exist"})
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "Successful delete"})
	}
}

func makeUpdateEndpoint(b Business) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var request UpdateReq

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request formar"})
			return
		}

		if request.Name != nil && *request.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "name is required"})
			return
		}

		if request.StartDate != nil && *request.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "start_date is required"})
			return
		}

		if request.EndDate != nil && *request.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "end_date is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := b.Update(id, &request); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Err: "Course doesn't exist"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "Successful update"})
	}
}
