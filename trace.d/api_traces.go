package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	saveRequest struct {
		ID        string `json:"id"`
		Name      string `json:"column"`
		Data      string `json:"data"`
		TraceType string `json:"type"`
	}

	deleteRequest struct {
		TraceID int    `json:"trace_id"`
		Table   string `json:"table"`
	}

	editGovernName struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

//TODO: needs refactoring and decoupling

// postCreateItem godoc
// @Summary Creation new item or period of item
// @Description endpoint for two operations
// @ID post-create-new-trace-period
// @Accept  json
// @Produce  json
// @Success 200 {string} string "ID"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/create [post]
// @Router /api/create-period [post]
func postCreateItem(c *gin.Context) {
	var (
		trace  NewTrace
		period BasicTrace
	)
	x, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}

	if ok := strings.Contains(c.Request.URL.Path, "create-period"); ok {
		if err := json.Unmarshal([]byte(x), &period.Fields); err != nil {
			NewError(c, http.StatusBadRequest, err)
			return
		}
		if id, err := period.createNewPeriod(); err != nil {
			NewError(c, http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"title": "Відстеження додано",
				"id":    id,
			})
		}
	}
	if err := json.Unmarshal([]byte(x), &trace); err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	if id, err := trace.createNewTrace(); err != nil {
		NewError(c, http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"title": "Відстеження додано",
			"id":    id,
		})
	}

}

// getRatings godoc
// @Summary Getting user data for displaying in user cabinet
// @Description get string by ID
// @ID get-ratings
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Ratings"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/ratings [get]
func getRatings(c *gin.Context) {
	columns, ratings, err := getReportData()
	if err != nil {
		NewError(c, http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"columns": columns,
			"ratings": ratings,
		})
	}
}

// getGoverns godoc
// @Summary Getting user data for displaying in user cabinet
// @Description get string by ID
// @ID get-governments
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Ratings"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/govs [get]
func getGoverns(c *gin.Context) {
	res, err := getGovernsList()
	if err != nil {
		NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"govs": res,
	})

}

// getTrace godoc
// @Summary Getting user data for displaying in user cabinet
// @Description get string by ID
// @ID get-trace
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Trace"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/v/{trk_id} [get]
func getTrace(c *gin.Context) {
	var (
		b  BasicTrace
		b2 BasicTrace
	)
	traceID := c.Param("trk_id")

	basic, err := b.getBasicData(traceID)
	if err != nil {
		NewError(c, http.StatusNotFound, err)
		return
	}
	period, err := b2.getPeriodicData(traceID)
	if err != nil {
		NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pl": basic.Fields,
		"pr": period})
}

// postEditTrackField godoc
// @Summary Getting user data for displaying in user cabinet
// @Description get string by ID
// @ID post-edit-field
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Trace field edited"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/v/{trk_id} [post]
func postEditTrackField(c *gin.Context) {
	var save saveRequest
	x, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal([]byte(x), &save); err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	if err := save.saveTraceChanges(); err != nil {
		NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"title": "Зміни збережено",
	})

}

// postEditGovernments godoc
// @Summary Getting user data for displaying in user cabinet
// @Description get string by ID
// @ID post-edit-govs
// @Accept  json
// @Produce  json
// @Success 200 {string} string "gov edited"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/govs/edit [post]
func postEditGovernments(c *gin.Context) {
	var edit editGovernName
	x, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal([]byte(x), &edit); err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	if err := edit.editGovName(); err != nil {
		NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"title": "Назву змінено",
	})

}

// postDeleteItem godoc
// @Summary Getting user data for displaying in user cabinet
// @Description get string by ID
// @ID post-delete-item
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Deleted"
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /api/delete [post]
func postDeleteItem(c *gin.Context) {
	var del deleteRequest
	x, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal([]byte(x), &del); err != nil {
		NewError(c, http.StatusBadRequest, err)
		return
	}
	if err := del.deleteItem(); err != nil {
		NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"title": "Відстеження видалено",
	})

}
