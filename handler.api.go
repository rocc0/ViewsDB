package main

import (
	"strconv"
	"net/http"
	"reflect"
	"io/ioutil"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type saveRequest struct {
	TypeOf 	int 	`json:"type"`
	Name  	string 	`json:"name"`
	Data 	string 	`json:"data"`
	Id 		int 	`json:"id"`
}
type deleteRequest struct {
	ItemId	int		`json:"item_id"`
	TblName	string	`json:"tbl_name"`
}

type editGov struct {
	Id int
	Name string
}

func getTrack(c *gin.Context) {
	viewID, err := strconv.Atoi(c.Param("trk_id"))
	basic := make(map[string]string)
	selected := make(map[string]string)
	if err == nil {
		btrace, err := getBasicData(viewID)
		ptrace, err := getPeriodicData(viewID)
		for k,v := range btrace.Fields {
			s := reflect.ValueOf(v)
			switch s.Interface().(type){
			case int64:
				basic[k] = string(strconv.Itoa(int(v.(int64))))
			default:
				basic[k] = string(v.([]uint8))
			}
		}
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"pl": basic,
				"pr": ptrace,
				"selected": selected})
		} else {
			log.Print(err)
		}
	} else {
		log.Print(err)
	}
}

func getRatings(c *gin.Context) {
	columns, ratings, err := getReportData()
	check(err)

	c.JSON(http.StatusOK, gin.H{
		"columns": columns,
		"ratings": ratings,
		})
}

func getGovernments(c *gin.Context) {
	res, err := getGovs()
	check(err)
	c.JSON(http.StatusOK, gin.H{
		"govs": res,
	})
}

func postTrackField(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var srq saveRequest
	err := json.Unmarshal([]byte(x), &srq)
	if err != nil {
		log.Print(err.Error())
	} else {
		err := editView(srq.Name, srq.Data, srq.TypeOf, srq.Id)
		if err != nil  {
			c.JSON(http.StatusBadRequest, gin.H{
				"ErrorTitle":   "Saving failed",
				"ErrorMessage": "Bad data"})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"title": "Data saved"})
			log.Print("OK", srq.TypeOf, srq.Name, srq.Data, srq.Id)
		}
	}
}

func postCreateItem(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var form map[string]interface{}
	json.Unmarshal([]byte(x), &form)
	if id, err := createNewItem(form); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"title": "Item added",
			"id": id,
		})

	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func postDeleteItem(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var delrq deleteRequest
	err := json.Unmarshal([]byte(x), &delrq)
	check(err)

	if err := deleteItem(delrq.ItemId, delrq.TblName); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"title": "Item removed",
		})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}


func postEditGovernments(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var editgov editGov
	err := json.Unmarshal([]byte(x), &editgov)
	check(err)
	if err := editGovName(editgov.Id, editgov.Name); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"title": "Gov name changed",
		})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
