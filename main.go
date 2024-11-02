package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/lib/pq"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db, _ = gorm.Open(sqlite.Open("bloggingDatabase.db"), &gorm.Config{})
)
type Post struct {
	ID uint				`json:"id"`
	Title string 		`json:"title"`
	Content string 		`json:"content"`
	Category string		`json:"category"`
	Tags pq.StringArray	`gorm:"type:text[]" json:"tags"`
	CreatedAt time.Time	`gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time	`gorm:"autoUpdateTime" json:"updatedAt"`
}

func main() {
	e := echo.New()

	db.AutoMigrate(&Post{})

	fmt.Println("Blog running on port 8080...")

	e.POST("/posts", CreateBlogPost)
	e.PUT("/posts/:id", UpdateBlogPost)
	e.DELETE("/posts/:id", DeleteBlogPost)
	e.GET("/posts/:id", GetBlogPost)
	e.GET("/posts", GetAllBlogPosts)

	e.Logger.Fatal(e.Start(":8080"))
}

func CreateBlogPost(c echo.Context) error {
	var post Post

	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	if err := db.Create(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, post)
}

func UpdateBlogPost(c echo.Context) error {
	var post Post

	id := c.Param("id")

	if err := db.Where("ID = ?", id).First(&post).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	db.Save(&post)
	return c.JSON(http.StatusOK, post)
}

func DeleteBlogPost(c echo.Context) error {
	var post Post
	id := c.Param("id")
	if err := db.Where("ID = ?", id).Delete(&post).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}
	return c.JSON(http.StatusNoContent, echo.Map{})
}

func GetBlogPost(c echo.Context) error {
	var post Post
	id := c.Param("id")
	if err := db.Where("ID = ?", id).First(&post).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}
	return c.JSON(http.StatusOK, post)
}

func GetAllBlogPosts(c echo.Context) error {
	post := []Post{}
	term := c.QueryParam("term")
	if term == "" {
		if err := db.Find(&post).Error; err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
		}
		return c.JSON(http.StatusOK, post)
	}
	if err := db.Where("Title LIKE ?", "%"+term+"%").Find(&post).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}
	return c.JSON(http.StatusOK, post)
}