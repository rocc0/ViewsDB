package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"strings"
)

type (
	saveRequest struct {
		ID        int    `json:"id"`
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

func getRatings(c *gin.Context) {
	columns, ratings, err := getReportData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"columns": columns,
			"ratings": ratings,
		})
	}
}

func getGoverns(c *gin.Context) {
	res, err := getGovernsList()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"govs": res,
		})
	}
}

func getTrace(c *gin.Context) {
	var b BasicTrace
	basic := make(map[string]string)
	traceID, err := strconv.Atoi(c.Param("trk_id"))

	if err == nil {
		err := b.getBasicData(traceID)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		period, err := getPeriodicData(traceID)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		for k, v := range b.Fields {
			s := reflect.ValueOf(v)
			switch s.Interface().(type) {
			case int64:
				basic[k] = string(strconv.Itoa(int(v.(int64))))
			default:
				basic[k] = string(v.([]uint8))
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"pl": basic,
			"pr": period})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func postTrackField(c *gin.Context) {
	var save saveRequest

	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &save)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		err := save.saveTraceField()
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"title": "Зміни збережено"})
		}
	}
}

//TODO: needs refactoring and decoupling
func postCreateItem(c *gin.Context) {
	var trace NewTrace
	var period BasicTrace

	x, _ := ioutil.ReadAll(c.Request.Body)

	if ok := strings.Contains(c.Request.URL.Path, "create-period"); ok {
		err := json.Unmarshal([]byte(x), &period.Fields)
		if err != nil {

			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			if id, err := period.createNewPeriod(); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"title": "Відстеження додано",
					"id":    id,
				})
			} else {

				c.AbortWithStatus(http.StatusBadRequest)
			}
		}
	} else {
		err := json.Unmarshal([]byte(x), &trace)
		if err != nil {
			log.Print(1, err)
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			if id, err := trace.createNewTrace(); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"title": "Відстеження додано",
					"id":    id,
				})
			} else {
				log.Print(2, err)
				c.AbortWithStatus(http.StatusBadRequest)
			}
		}
	}
}

func postDeleteItem(c *gin.Context) {
	var del deleteRequest

	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &del)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		if err := del.deleteItem(); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"title": "Відстеження видалено",
			})
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}

func postEditGovernments(c *gin.Context) {
	var edit editGovernName

	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &edit)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		if err := edit.editGovName(); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"title": "Назву змінено",
			})
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}
