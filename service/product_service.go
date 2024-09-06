package service

import (
	"fmt"
	"strconv"

	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
	GetProductById(c *fiber.Ctx) error
	PaginatedProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type implProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductService {
	return &implProductService{
		db: db,
	}
}

type ProductStruct struct {
	Name       string  `json:"name" validate:"required,min=1,max=100"`
	Price      float64 `json:"price" validate:"required"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

func (s *implProductService) CreateProduct(c *fiber.Ctx) error {
	body := new(ProductStruct)

	// user := c.Locals("user").(*model.User)
	// fmt.Println("this is user", user.ID)

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

	product := &model.Product{
		Name:       body.Name,
		Price:      body.Price,
		CategoryID: body.CategoryID,
	}

	if err := s.db.Create(product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "new product created",
	})
}

func (s *implProductService) GetAllProducts(c *fiber.Ctx) error {
	products := make([]model.Product, 0)

	if err := s.db.Preload("Category").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to get all products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (s *implProductService) PaginatedProduct(c *fiber.Ctx) error {
	sortOrder := c.Query("sort", "asc")

	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}

	search := c.Query("search", "")
	skip := limit * page

	products := make([]model.Product, 0)

	query := s.db.Preload("Category").Model(&model.Product{}).Where("name ILIKE ?", "%"+search+"%")

	var totalRows int64
	if err := query.Count(&totalRows).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count products",
		})
	}

	if err := query.Offset(skip).Limit(limit).Order("id " + sortOrder).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve products",
		})
	}

	totalPages := int(totalRows) / limit
	if int(totalRows)%limit != 0 {
		totalPages++
	}

	response := map[string]interface{}{
		"data":         products,
		"current_page": page,
		"data_limit":   limit,
		"total_rows":   totalRows,
		"total_pages":  totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (s *implProductService) GetProductById(c *fiber.Ctx) error {
	product := &model.Product{}

	id := c.Params("id")

	if err := s.db.Preload("Category").First(product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (s *implProductService) UpdateProduct(c *fiber.Ctx) error {
	body := new(ProductStruct)

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

	product := &model.Product{
		Name:       body.Name,
		Price:      body.Price,
		CategoryID: body.CategoryID,
	}

	if err := s.db.Model(&model.Product{}).Where("id = ?", id).Updates(product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database error",
		})
	}

	fmt.Println(product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "product updated",
	})
}

func (s *implProductService) DeleteProduct(c *fiber.Ctx) error {
	product := &model.Product{}

	id := c.Params("id")

	if err := s.db.Delete(product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "product deleted",
	})
}
