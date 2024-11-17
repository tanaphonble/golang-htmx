package dog

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	dogDatabase DogDatabase
}

func RegisterHandler(e *gin.Engine, dogDatabase DogDatabase) {
	h := newHandler(dogDatabase)

	e.GET("/components/dog-list", h.List)
	e.POST("/dogs", h.Create)
	e.DELETE("/dogs/:id", h.Delete)
}

func newHandler(dogDatabase DogDatabase) *handler {
	if dogDatabase == nil {
		log.Fatal("dog database dependency cannot be nil")
	}

	return &handler{dogDatabase: dogDatabase}
}

func (h *handler) List(c *gin.Context) {
	dogs, err := h.dogDatabase.SelectAll()
	if err != nil {
		slog.Error("Failed to retrieve dogs", slog.String("error", err.Error()))
		c.String(http.StatusInternalServerError, "Failed to retrieve dogs")
		return
	}

	c.Header("Content-Type", "text/html")
	if err := dogListTableTemplate.Execute(c.Writer, dogs); err != nil {
		slog.Error("Failed to render template", slog.String("error", err.Error()))
		c.String(http.StatusInternalServerError, "Failed to render template")
	}
}

func (h *handler) Create(c *gin.Context) {
	name := c.PostForm("name")
	breed := c.PostForm("breed")

	if name == "" || breed == "" {
		c.String(http.StatusBadRequest, "Name and Breed are required")
		return
	}

	newDog := Dog{
		ID:    uuid.NewString(),
		Name:  name,
		Breed: breed,
	}

	err := h.dogDatabase.Insert(newDog)
	if err != nil {
		slog.Error("Failed to add dog", slog.String("error", err.Error()))
		c.String(http.StatusInternalServerError, "Failed to add dog")
		return
	}

	c.Status(http.StatusCreated)
}

func (h *handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.dogDatabase.Delete(id)
	if err != nil {
		slog.Error("Failed to delete dog", slog.String("error", err.Error()), slog.String("id", id))
		c.String(http.StatusInternalServerError, "Failed to delete dog")
		return
	}

	c.Status(http.StatusNoContent)
}
