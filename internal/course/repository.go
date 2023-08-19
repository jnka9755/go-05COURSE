package course

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jnka9755/go-05DOMAIN/domain"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, course *domain.Course) error
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, course *UpdateCourse) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repository struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepository(l *log.Logger, db *gorm.DB) Repository {

	return &repository{
		db:  db,
		log: l,
	}
}

func (r *repository) Create(ctx context.Context, course *domain.Course) error {

	if err := r.db.WithContext(ctx).Create(course).Error; err != nil {
		r.log.Println("Error-Repository CreateCourse->", err)
		return err
	}

	r.log.Println("Repository -> Create course with id: ", course.ID)

	return nil
}

func (r *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {

	var courses []domain.Course

	tx := r.db.WithContext(ctx).Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)

	result := tx.Order("created_at desc").Find(&courses)

	if result.Error != nil {
		r.log.Println("Error-Repository GetAllCourses ->", result.Error)
		return nil, result.Error
	}

	return courses, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Course, error) {

	course := domain.Course{ID: id}

	if err := r.db.WithContext(ctx).First(&course).Error; err != nil {
		r.log.Println("Error-Repository GetCourse ->", err)

		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}

		return nil, err
	}

	return &course, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {

	course := domain.Course{ID: id}

	result := r.db.WithContext(ctx).Delete(&course)

	if result.Error != nil {
		r.log.Println("Error-Repository DeleteCourse ->", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("Course with ID -> '%s' doesn't exist", course.ID)
		return ErrNotFound{id}
	}

	r.log.Println("Repository -> Delete course with id: ", course.ID)

	return nil
}

func (r *repository) Update(ctx context.Context, course *UpdateCourse) error {

	r.log.Println("Udate course Repository")

	values := make(map[string]interface{})

	if course.Name != nil {
		values["name"] = *course.Name
	}

	if course.StartDate != nil {
		values["start_date"] = *course.StartDate
	}

	if course.EndDate != nil {
		values["end_date"] = *course.EndDate
	}

	result := r.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", course.ID).Updates(values)

	if result.Error != nil {
		r.log.Println("Error-Repository UdateCourse ->", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("Course with ID -> '%s' doesn't exist", course.ID)
		return ErrNotFound{course.ID}
	}

	return nil
}

func (r *repository) Count(ctx context.Context, filters Filters) (int, error) {

	var count int64
	tx := r.db.WithContext(ctx).Model(domain.Course{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		r.log.Println("Error-Repository CountCourse ->", err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}

	return tx
}
