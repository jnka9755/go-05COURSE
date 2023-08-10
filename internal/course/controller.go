package course

import (
	"context"
	"errors"
	"fmt"

	"github.com/jnka9755/go-05META/meta"
	"github.com/jnka9755/go-05RESPONSE/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

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

	GetReq struct {
		ID string
	}

	GetAllReq struct {
		Name  string
		Limit int
		Page  int
	}

	UpdateReq struct {
		ID        string
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	DeleteReq struct {
		ID string
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateReq)

		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate == "" {
			return nil, response.BadRequest(ErrStarDateRequired.Error())
		}

		if req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		responseCourse, err := b.Create(ctx, &req)

		if err != nil {

			if err == ErrInvalidStartDate ||
				err == ErrInvalidEndtDate ||
				err == ErrEndDateHigherStart ||
				err == ErrEqualDates {
				return nil, response.BadRequest(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("Success create course", responseCourse, nil), nil
	}
}

func makeGetAllEndpoint(b Business, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			Name: req.Name,
		}

		count, err := b.Count(ctx, filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := b.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success get courses", courses, meta), nil
	}
}

func makeGetEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		course, err := b.Get(ctx, req.ID)

		if err != nil {
			return nil, response.NotFound(err.Error())
		}

		return response.OK("Success get course", course, nil), nil
	}
}

func makeDeleteEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		err := b.Delete(ctx, req.ID)

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK(fmt.Sprintf("Success delete course with ID -> '%s'", req.ID), nil, nil), nil
	}
}

func makeUpdateEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UpdateReq)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest(ErrStarDateRequired.Error())
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		if err := b.Update(ctx, &req); err != nil {

			if err == ErrInvalidStartDate ||
				err == ErrInvalidEndtDate ||
				err == ErrEndDateHigherStart ||
				err == ErrEqualDates {
				return nil, response.BadRequest(err.Error())
			}

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK(fmt.Sprintf("Success update course with ID -> '%s'", req.ID), nil, nil), nil
	}
}
