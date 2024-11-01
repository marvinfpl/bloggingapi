package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"time"
	"strconv"
)

type Post struct {
	ID int				`json:"id"`
	Title string 		`json:"title"`
	Content string 		`json:"content"`
	Category string		`json:"category"`
	Tags []string 		`json:"tags"`
	CreatedAt time.Time	`json:"createdAt"`
	UpdatedAt time.Time	`json:"updatedAt"`
}

var posts = []Post{
	{ID: 1, Title: "first post", Content: "none", Category: "coding", Tags: []string{"coding with marvin"},CreatedAt: time.Now() },
}

func main() {
	e := echo.New()

	fmt.Println("Blog running on port 8080...")

	e.POST("/posts", createPost)
	e.PUT("/posts/:id", updatePost)
	e.DELETE("/posts/:id", deletePost)
	e.GET("/posts/:id", getPost)
	e.GET("/posts", getAllPosts)

	e.Logger.Fatal(e.Start(":8080"))
}

func createPost(c echo.Context) error {
	post := new(Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	post.CreatedAt = time.Now()
	posts = append(posts, *post)
	return c.JSON(http.StatusCreated, post)
}

func updatePost(c echo.Context) error {
	idStr := c.Param("/posts/:id")
	if idStr == "" {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}
	id, _ := strconv.Atoi(idStr)
	for _, post := range posts{
		if post.ID == id {
			if err := c.Bind(&post); err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
			}
			post.UpdatedAt = time.Now()
			return c.JSON(http.StatusOK, post)
		}
	}
	return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
}

func deletePost(c echo.Context) error {
	idStr := c.Param("/posts/:id")
	if idStr == "" {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})		
	}
	id, _ := strconv.Atoi(idStr)
	for index, post := range posts {
		if post.ID == id {
			posts = append(posts[:index], posts[index+1:]... )
			return c.JSON(http.StatusNoContent, echo.Map{"message": "post deleted"})
		}
	}
	return c.JSON(http.StatusBadRequest, echo.Map{"error": "post not found"})
}

func getPost(c echo.Context) error {
	idStr := c.Param("/posts/:id")
	if idStr == "" {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})		
	}
	id, _ := strconv.Atoi(idStr)
	for _, post := range posts {
		if post.ID == id {
			return c.JSON(http.StatusOK, post)
		}
	}
	return c.JSON(http.StatusBadRequest, echo.Map{"error": "post not found"})
}

func getAllPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, posts)
}