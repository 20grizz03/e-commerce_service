package tests

import (
	"e-com/app/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAllProducts(t *testing.T) {
	mockMethod := new(mocks.Methods)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/products", nil)
	mockMethod.On("GetAllProducts", w, req).Return().Once()

	mockMethod.GetAllProducts(w, req)

	mockMethod.AssertCalled(t, "GetAllProducts", w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestGetProductByID(t *testing.T) {
	mockMethod := new(mocks.Methods)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/products/1", nil)
	mockMethod.On("GetProductsById", w, req).Return().Once()

	mockMethod.GetProductsById(w, req)

	mockMethod.AssertCalled(t, "GetProductsById", w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestPostAndPutProducts_POST(t *testing.T) {
	mockMethods := new(mocks.Methods)

	w := httptest.NewRecorder()
	reqBody := `{"id": 123, "name": "Шапка"}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	mockMethods.On("PostAndPutProducts", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
		resp := args.Get(0).(http.ResponseWriter)
		request := args.Get(1).(*http.Request)

		// Проверяем, что метод POST, и записываем ответ
		if request.Method == http.MethodPost {
			resp.WriteHeader(http.StatusCreated)
			_, _ = resp.Write([]byte(`{"message": "Product created successfully"}`))
		}
	}).Return()

	mockMethods.PostAndPutProducts(w, req)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Equal(t, `{"message": "Product created successfully"}`, w.Body.String())

	mockMethods.AssertCalled(t, "PostAndPutProducts", w, req)
}
