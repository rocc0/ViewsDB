package main

import (
	"net/http"
	"strings"

	"./httputil"
	"github.com/gin-gonic/gin"
)

// showIndexPage godoc
// @Summary Index page html endpoint
// @Description Show html main/search page
// @ID get-index-page
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "ok"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router / [get]
func showIndexPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Пошук відстежень",
	}, "index.html")
}

// showRegisterPage godoc
// @Summary Index page html endpoint
// @Description Show html register page
// @ID get-register-page
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "register"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /u/register [get]
func showRegisterPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Реєстрація",
	}, "index.html")
}

// showLoginPage godoc
// @Summary Index page html endpoint
// @Description Show html login page
// @ID get-login-page
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "login"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /u/login [get]
func showLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Вхід",
	}, "index.html")
}

// showRatings godoc
// @Summary Getting user data for displaying in user cabinet
// @Description Show html ratings page
// @ID get-string-by-string
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "ratings"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /ratings [get]
func showRatingsPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Якісні показники",
	}, "index.html")
}

// showTraceCreationPage godoc
// @Summary Getting user data for displaying in user cabinet
// @Description Show html trace creation page
// @ID get-creation-page
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "creation"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /track/create [get]
func showTraceCreationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Додати відстеження",
	}, "index.html")
}

// showUserPage godoc
// @Summary Getting user data for displaying in user cabinet
// @Description Show html user page
// @ID get-user-page
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "user page"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /u/cabinet [get]
func showUserPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Кабінет користувача",
	}, "index.html")
}

// showEditGovsNames godoc
// @Summary Getting user data for displaying in user cabinet
// @Description Show html govs page
// @ID get-govs
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "governments"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /govs/edit [get]
func showEditGovsNames(c *gin.Context) {
	render(c, gin.H{
		"title": "Редагувати назви державних органів",
	}, "index.html")
}

// showTracePage godoc
// @Summary Getting user data for displaying in user cabinet
// @Description Show html trace page
// @ID get-trac-by-string
// @Accept  text/html
// @Produce  text/html
// @Success 200 {string} string "trace"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /id/{trk_id} [get]
// @Router /id/{trk_id}/edit [get]
func showTracePage(c *gin.Context) {
	var (
		trace BasicTrace
		title string
	)
	traceID := c.Param("trk_id")
	url := c.Request.URL.Path
	if strings.Contains(url, "edit") {
		title = "Редагування | "
	} else {
		title = ""
	}

	if err := trace.getBasicData(traceID); err != nil {
		httputil.NewError(c, http.StatusNotFound, err)
		return
	}

	render(c, gin.H{
		"title": title + string(trace.Fields["reg_name"].([]uint8)),
	}, "index.html")

}
