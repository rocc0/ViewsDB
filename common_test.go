package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"github.com/gin-gonic/gin"
)

var userList = []User{
	User{Email: "user1", Password: "pass1"},
	User{Email: "user2", Password: "pass2"},
	User{Email: "user3", Password: "pass3"},
}

var tmpUserList []User


func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

// Helper function to create a router during testing
func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	if withTemplates {
		r.LoadHTMLGlob("templates/*")
	}
	return r
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func testMiddlewareRequest(t *testing.T, r *gin.Engine, expectedHTTPCode int) {
	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/", nil)

	// Process the request and test the response
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == expectedHTTPCode
	})
}

// This function is used to store the main lists into the temporary one
// for testing
func saveLists() {
	tmpUserList = userList
}

// This function is used to restore the main lists from the temporary one
func restoreLists() {
	userList = tmpUserList
}