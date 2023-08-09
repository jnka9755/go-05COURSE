package course

import (
	"context"
	"log"
	"time"

	"github.com/jnka9755/go-05DOMAIN/domain"
)

type (
	Business interface {
		Create(ctx context.Context, course *CreateReq) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, course *UpdateReq) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	business struct {
		log        *log.Logger
		repository Repository
	}

	Filters struct {
		Name string
	}

	UpdateCourse struct {
		ID        string
		Name      *string
		StartDate *time.Time
		EndDate   *time.Time
	}
)

func NewBusiness(log *log.Logger, repository Repository) Business {
	return &business{
		log:        log,
		repository: repository,
	}
}

func (b business) Create(ctx context.Context, request *CreateReq) (*domain.Course, error) {

	var startDateParsed, endDateParsed time.Time

	if request.StartDate != "" {
		date, err := time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			b.log.Println(err)
			return nil, ErrInvalidStartDate
		}
		startDateParsed = date
	}

	if request.EndDate != "" {
		date, err := time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			b.log.Println(err)
			return nil, ErrInvalidEndtDate
		}
		endDateParsed = date
	}

	course := domain.Course{
		Name:      request.Name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := b.repository.Create(ctx, &course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (b business) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {

	courses, err := b.repository.GetAll(ctx, filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (b business) Get(ctx context.Context, id string) (*domain.Course, error) {

	course, err := b.repository.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (b business) Delete(ctx context.Context, id string) error {

	return b.repository.Delete(ctx, id)
}

func (b business) Update(ctx context.Context, request *UpdateReq) error {

	var startDateParsed, endDateParsed *time.Time

	if request.StartDate != nil {
		date, err := time.Parse("2006-01-02", *request.StartDate)
		if err != nil {
			b.log.Println(err)
			return ErrInvalidStartDate
		}
		startDateParsed = &date
	}

	if request.EndDate != nil {
		date, err := time.Parse("2006-01-02", *request.StartDate)
		if err != nil {
			b.log.Println(err)
			return ErrInvalidEndtDate
		}
		endDateParsed = &date
	}

	courseUpdate := UpdateCourse{
		ID:        request.ID,
		Name:      request.Name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	return b.repository.Update(ctx, &courseUpdate)
}

func (b business) Count(ctx context.Context, filters Filters) (int, error) {
	return b.repository.Count(ctx, filters)
}
