package main

import (
	"fmt"
	"go-pg-gorm/models"
	"go-pg-gorm/storage"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

// Repository is a struct that contains the database connection
type Repository struct {
	Db *gorm.DB
}

// gorm gives us the ability to interact with the database and DB is the connection to the database that we will use to interact with the database while Db is a pointer to the connection to the database

// book functions
func (r *Repository) CreateBooks(context *fiber.Ctx) error {
	book := Book{}

	if err := context.BodyParser(&book); err != nil {
		return context.Status(400).JSON(&fiber.Map{
			"message": "unable to parse JSON",
		})
	}

	err := r.Db.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to create book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book created successfully",
	})

	return nil
}

//  createbook is a method that takes a context and returns an error. The context is a pointer to the fiber context that we will use to interact with the request and response objects. The error is a pointer to the error object that we will use to handle errors in our application. book is a struct that we will use to store the data that we will receive from the request body.

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	books := &[]models.Books{}
	// slice of models of books
	err := r.Db.Find(&books).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to fetch books",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    books,
	})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	books := models.Books{}

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id is required",
		})
		return nil
	}

	err := r.Db.Delete(books, id)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book deleted successfully",
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	books := &models.Books{}

	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id is required",
		})
		return nil
	}
	fmt.Println("the id is: ", id)

	err := r.Db.Where("id = ?", id).First(books).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to fetch book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book fetched successfully",
		"data":    books,
	})
	return nil
}

func (r *Repository) UpdateBooks(context *fiber.Ctx) error{
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id is required",
		})
		return nil
	}

	books := models.Books{}
	if err := context.BodyParser(&books); err != nil {
		return context.Status(400).JSON(&fiber.Map{
			"message": "unable to parse JSON",
		})
	}

	err := r.Db.Model(&books).Where("id = ?", id).Updates(books).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to update book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book updated successfully",
	})
	return nil
}

// user functions
func(r *Repository) CreateUser(c *fiber.Ctx) error{
	user := models.Users{}

	if err := c.BodyParser(&user); err != nil{
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to parse JSON",
		})
	}

	err := r.Db.Create(&user).Error
	if err != nil{
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "user can't be created",
		})
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user created successfully",
		"data": user,
	})

	return nil
}

func (r *Repository) GetUsers(c *fiber.Ctx) error{
	users := &[]models.Users{}

	err := r.Db.Preload("Book").Find(&users).Error
	if err != nil{
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "can't fetch all the users",
		})
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "users fetched successfully",
		"data": users,
	})
	return nil
}

func (r *Repository) GetUserByID(c *fiber.Ctx) error{
	users := &models.Users{}

	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id is required",
		})
	}

	err := r.Db.Preload("Book").Where("id = ?", id).First(users).Error
	if err != nil{
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "can't fetch the user",
		})
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "user fetched successfully",
		"data": users,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// book routes
	// all these createbooks, etc are functions used here as methods
	api.Post("/books", r.CreateBooks)
	api.Put("/books/:id", r.UpdateBooks)
	api.Get("/books", r.GetBooks)
	api.Get("/books/:id", r.GetBookByID)
	api.Delete("/books/:id", r.DeleteBook)

	// user routes
	api.Post("/users", r.CreateUser)
	api.Get("/users", r.GetUsers)
	api.Get("/users/:id", r.GetUserByID)
}

func main() {
	fmt.Println("Hello, World!")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		log.Fatal("CLIENT_URL is not set in the environment")
	}

	db, err := storage.NewConnection(dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Error migrating the database", err)
	}

	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("Error migrating the database", err)
	}

	r := Repository{
		Db: db,
	}
	// fiber has almost similar syntax than express.js but way faster
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: clientURL,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	r.SetupRoutes(app)
	app.Listen(":8000")
}
