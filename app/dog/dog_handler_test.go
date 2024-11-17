package dog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type DogHandlerTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockDogDatabase *MockDogDatabase
	router          *gin.Engine
}

func (suite *DogHandlerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockDogDatabase = NewMockDogDatabase(suite.ctrl)
	suite.router = gin.Default()
	RegisterHandler(suite.router, suite.mockDogDatabase)
}

func (suite *DogHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *DogHandlerTestSuite) TestList_Success() {
	// Arrange
	expectedDogs := []Dog{
		{ID: "1", Name: "Rover", Breed: "Labrador"},
		{ID: "2", Name: "Max", Breed: "Beagle"},
	}
	suite.mockDogDatabase.EXPECT().SelectAll().Return(expectedDogs, nil)

	req := httptest.NewRequest(http.MethodGet, "/components/dog-list", nil)
	rec := httptest.NewRecorder()

	// Act
	suite.router.ServeHTTP(rec, req)

	// Assert
	suite.Equal(http.StatusOK, rec.Code)
	suite.Contains(rec.Body.String(), `<td>Rover</td>`)
	suite.Contains(rec.Body.String(), `<td>Max</td>`)
}

func (suite *DogHandlerTestSuite) TestList_Error() {
	// Arrange
	suite.mockDogDatabase.EXPECT().SelectAll().Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/components/dog-list", nil)
	rec := httptest.NewRecorder()

	// Act
	suite.router.ServeHTTP(rec, req)

	// Assert
	suite.Equal(http.StatusInternalServerError, rec.Code)
	suite.Contains(rec.Body.String(), "Failed to retrieve dogs")
}

func (suite *DogHandlerTestSuite) TestCreate_Success() {
	// Arrange
	suite.mockDogDatabase.EXPECT().Insert(gomock.Any()).Return(nil)

	formData := "name=Charlie&breed=Poodle"
	req := httptest.NewRequest(http.MethodPost, "/dogs", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	// Act
	suite.router.ServeHTTP(rec, req)

	// Assert
	suite.Equal(http.StatusCreated, rec.Code)
	suite.Contains(rec.Body.String(), `<td>Charlie</td>`)
	suite.Contains(rec.Body.String(), `<td>Poodle</td>`)
}

func (suite *DogHandlerTestSuite) TestCreate_ValidationError() {
	// Arrange
	formData := "name=&breed="
	req := httptest.NewRequest(http.MethodPost, "/dogs", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	// Act
	suite.router.ServeHTTP(rec, req)

	// Assert
	suite.Equal(http.StatusBadRequest, rec.Code)
	suite.Contains(rec.Body.String(), "Name and Breed are required")
}

func (suite *DogHandlerTestSuite) TestCreate_DatabaseError() {
	// Arrange
	suite.mockDogDatabase.EXPECT().Insert(gomock.Any()).Return(assert.AnError)

	formData := "name=Charlie&breed=Poodle"
	req := httptest.NewRequest(http.MethodPost, "/dogs", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	// Act
	suite.router.ServeHTTP(rec, req)

	// Assert
	suite.Equal(http.StatusInternalServerError, rec.Code)
	suite.Contains(rec.Body.String(), "Failed to add dog")
}

func TestDogHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DogHandlerTestSuite))
}
