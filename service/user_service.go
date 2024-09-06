package service

import (
	"fmt"
	"os"
	"time"

	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
}

type implUserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &implUserService{
		db: db,
	}
}

type RegisterStruct struct {
	Username     string `json:"username" validate:"required,min=5,max=30"`
	Email        string `json:"email" validate:"required,max=50"`
	Phone_Number string `json:"phone_number" validate:"required,max=20"`
	Password     string `json:"password" validate:"required,min=5,max=20"`
}

func (s *implUserService) Register(c *fiber.Ctx) error {
	body := new(RegisterStruct)

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	user := &model.User{
		Username:    body.Username,
		Email:       body.Email,
		PhoneNumber: body.Phone_Number,
		Password:    string(hashedPassword),
	}

	if err := s.db.Create(user).Error; err != nil {
		fmt.Println("this is the error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to register",
			"error":   err.Error(),
		})
	}

	fmt.Println(user, hashedPassword)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "registration success",
	})
}

type LoginStruct struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *implUserService) Login(c *fiber.Ctx) error {
	body := new(LoginStruct)

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

	user := &model.User{}

	if err := s.db.Where("email = ?", body.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid username or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "database error",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid username or password",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Token expires in 24 hours
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to login",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "login successful",
	})
}

type PasswordStruct struct {
	Password string `json:"password" validate:"required,min=5,max=20"`
}

func (s *implUserService) UpdatePassword(c *fiber.Ctx) error {
	body := new(PasswordStruct)

	userID := c.Locals("user_id")

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	user := &model.User{}

	if err := s.db.Model(&user).Where("id = ?", userID).Update("password", string(hashedPassword)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}
