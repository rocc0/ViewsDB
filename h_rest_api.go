package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

type saveRequest struct {
	TraceType int    `json:"type"`
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Data      string `json:"data"`
}
type deleteRequest struct {
	TraceId   int    `json:"item_id"`
	TableName string `json:"tbl_name"`
}

type editGovernName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

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
	var bTrace Trace
	basic := make(map[string]string)

	traceId, err := strconv.Atoi(c.Param("trk_id"))

	if err == nil {
		err := bTrace.getBasicData(traceId)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		pTrace, _ := getPeriodicData(traceId)

		for k, v := range bTrace.Fields {
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
			"pr": pTrace})
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

func postCreateItem(c *gin.Context) {
	var t Trace

	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &t.Fields)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		if id, err := t.createNewItem(); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"title": "Відстеження додано",
				"id":    id,
			})
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}

func postDeleteItem(c *gin.Context) {
	var delete deleteRequest

	x, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(x), &delete)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		if err := delete.deleteItem(); err == nil {
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
