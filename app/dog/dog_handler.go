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

	htmlContent := `<tr>
        <td>` + newDog.Name + `</td>
        <td>` + newDog.Breed + `</td>
    </tr>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusCreated, htmlContent)
}
