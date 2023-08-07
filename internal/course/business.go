package course

import (
	"log"
	"time"

	"github.com/jnka9755/go-05DOMAIN/domain"
)

type (
	Business interface {
		Create(course *CreateReq) (*domain.Course, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Delete(id string) error
		Update(id string, course *UpdateReq) error
		Count(filters Filters) (int, error)
	}

	business struct {
		log        *log.Logger
		repository Repository
	}

	Filters struct {
		Name string
	}

	UpdateCourse struct {
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

func (b business) Create(request *CreateReq) (*domain.Course, error) {

	var startDateParsed, endDateParsed time.Time

	if request.StartDate != "" {
		date, err := time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			b.log.Println(err)
			return nil, err
		}
		startDateParsed = date
	}

	if request.EndDate != "" {
		date, err := time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			b.log.Println(err)
			return nil, err
		}
		endDateParsed = date
	}

	course := domain.Course{
		Name:      request.Name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := b.repository.Create(&course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (b business) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {

	b.log.Println("GetAll course Business")
	courses, err := b.repository.GetAll(filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (b business) Get(id string) (*domain.Course, error) {

	b.log.Println("Get course Business")
	course, err := b.repository.Get(id)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (b business) Delete(id string) error {

	b.log.Println("Delete course Business")
	return b.repository.Delete(id)
}

func (b business) Update(id string, course *UpdateReq) error {

	b.log.Println("Update course Business")

	var startDateParsed, endDateParsed *time.Time

	if course.StartDate != nil {
		date, err := time.Parse("2006-01-02", *course.StartDate)
		if err != nil {
			b.log.Println(err)
			return err
		}
		startDateParsed = &date
	}

	if course.EndDate != nil {
		date, err := time.Parse("2006-01-02", *course.StartDate)
		if err != nil {
			b.log.Println(err)
			return err
		}
		endDateParsed = &date
	}

	courseUpdate := UpdateCourse{
		Name:      course.Name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	return b.repository.Update(id, &courseUpdate)
}

func (b business) Count(filters Filters) (int, error) {
	return b.repository.Count(filters)
}
