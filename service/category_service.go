package service

import (
	"app2/config"
	"app2/middleware"
	"app2/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CategoryGroup(router fiber.Router) {
	categoryRoutes := router.Group("/category")

	categoryRoutes.Post("/", middleware.Authenticate, middleware.GetCredential, CreateCategory)
	categoryRoutes.Get("/", GetAllCategory)
	categoryRoutes.Get("/:id", GetCategoryById)
	categoryRoutes.Put("/:id", middleware.Authenticate, middleware.GetCredential, UpdateCategory)
	categoryRoutes.Delete("/:id", middleware.Authenticate, middleware.GetCredential, DeleteCategory)
}

/*
*
*
*
*
*
 */

type CategoryStruct struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

// create category

func CreateCategory(c *fiber.Ctx) error {
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

	db := config.GetDB()

	category := &model.Category{
		Name: body.Name,
	}

	if err := db.Create(category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":"new category created",
	})
}

/*
*
*
*
*
*
 */

// get all category

func GetAllCategory(c *fiber.Ctx) error {
	category := make([]model.Category, 0)

	db := config.GetDB()

	if err := db.Preload("Products").Find(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to get all category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

/*
*
*
*
*
*
 */

// get category by id

func GetCategoryById(c *fiber.Ctx) error {
	category := &model.Category{}

	id := c.Params("id")

	db := config.GetDB()

	if err := db.First(category, id).Error; err != nil {
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

/*
*
*
*
*
*
 */

// update category

func UpdateCategory(c *fiber.Ctx) error {
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

	db := config.GetDB()

	if err := db.Model(&model.Category{}).Where("id = ?", id).Updates(category).Error; err != nil {
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

/*
*
*
*
*
*
 */

// delete category

func DeleteCategory(c *fiber.Ctx) error {
	category := &model.Category{}

	id := c.Params("id")

	db := config.GetDB()

	if err := db.Delete(category, id).Error; err != nil {
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
