package service

import (
	"fmt"
	"strconv"

	"app2/config"
	"app2/middleware"
	"app2/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProductGroup(router fiber.Router) {
	ProductGroup := router.Group("/product")

	ProductGroup.Post("/", middleware.Authenticate, middleware.GetCredential, middleware.Authorize(0, 1), CreateProduct)
	ProductGroup.Get("/", GetAllProducts)
	ProductGroup.Get("/page", PaginatedProduct)
	ProductGroup.Get("/:id", GetProductById)
	ProductGroup.Put("/:id", middleware.Authenticate, middleware.GetCredential, UpdateProduct)
	ProductGroup.Delete("/:id", middleware.Authenticate, middleware.GetCredential, DeleteProduct)
}

/*
*
*
*
*
*
 */

type ProductStruct struct {
	Name       string  `json:"name" validate:"required,min=1,max=100"`
	Price      float64 `json:"price" validate:"required"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

// create product

func CreateProduct(c *fiber.Ctx) error {
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

	db := config.GetDB()

	if err := db.Create(product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "new product created",
	})
}

/*
*
*
*
*
*
 */

// get all product

func GetAllProducts(c *fiber.Ctx) error {
	products := make([]model.Product, 0)

	db := config.GetDB()

	if err := db.Preload("Category").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "unable to get all products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

/*
*
*
*
*
*
 */

// paginate product

func PaginatedProduct(c *fiber.Ctx) error {
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

	db := config.GetDB()

	query := db.Preload("Category").Model(&model.Product{}).Where("name ILIKE ?", "%"+search+"%")

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

/*
*
*
*
*
*
 */

// get product by id

func GetProductById(c *fiber.Ctx) error {
	product := &model.Product{}

	id := c.Params("id")

	db := config.GetDB()

	if err := db.Preload("Category").First(product, id).Error; err != nil {
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

/*
*
*
*
*
*
 */

// update product

func UpdateProduct(c *fiber.Ctx) error {
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

	db := config.GetDB()

	if err := db.Model(&model.Product{}).Where("id = ?", id).Updates(product).Error; err != nil {
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

/*
*
*
*
*
*
 */

// delete product

func DeleteProduct(c *fiber.Ctx) error {
	product := &model.Product{}

	id := c.Params("id")

	db := config.GetDB()

	if err := db.Delete(product, id).Error; err != nil {
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
