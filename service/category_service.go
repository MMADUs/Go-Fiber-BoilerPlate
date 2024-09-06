package service

import (
	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(c *fiber.Ctx) error
	GetAllCategory(c *fiber.Ctx) error
	GetCategoryById(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type implCategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &implCategoryService{
		db: db,
	}
}

type CategoryStruct struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

func (s *implCategoryService) CreateCategory(c *fiber.Ctx) error {
	body := new(CategoryStruct)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "invalid data",
		})
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation error",
		})
	}

	category := &model.Category{
		Name: body.Name,
	}

	if err := s.db.Create(category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "new category created",
	})
}

func (s *implCategoryService) GetAllCategory(c *fiber.Ctx) error {
	category := make([]model.Category, 0)

	if err := s.db.Preload("Products").Find(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to get all category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

func (s *implCategoryService) GetCategoryById(c *fiber.Ctx) error {
	category := &model.Category{}

	id := c.Params("id")

	if err := s.db.First(category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

func (s *implCategoryService) UpdateCategory(c *fiber.Ctx) error {
	body := new(CategoryStruct)

	id := c.Params("id")

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "invalid data",
		})
	}

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation error",
		})
	}

	category := &model.Category{
		Name: body.Name,
	}

	if err := s.db.Model(&model.Category{}).Where("id = ?", id).Updates(category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "category updated",
	})
}

func (s *implCategoryService) DeleteCategory(c *fiber.Ctx) error {
	category := &model.Category{}

	id := c.Params("id")

	if err := s.db.Delete(category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "category deleted",
	})
}
