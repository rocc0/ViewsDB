package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
<<<<<<< HEAD
	"reflect"
	"log"
	"io/ioutil"
	"encoding/json"
||||||| merged common ancestors
	"log"

=======
	"reflect"
>>>>>>> 219ac6dc41467e6dddf6f1dbfa707efdd8ecdf7b
)

type saveRequest struct {
	Name  string `json:"name"`
	Data string `json:"data"`
	Id int `json:"id,string"`
}

func viewData(c *gin.Context) {
	viewID, err := strconv.Atoi(c.Param("view_id"));
	if err == nil {
		view, err := getViewById(viewID)
<<<<<<< HEAD
		arr := make(map[string]string)
		selected := make(map[string]string)
		for k,v := range view.Fields {
			s := reflect.ValueOf(v)
			switch s.Interface().(type){
			case int64:
				arr[k] = string(strconv.Itoa(int(v.(int64))))
			default:
				arr[k] = string(v.([]uint8))
			}
		}
||||||| merged common ancestors
=======
		arr := make(map[string]string)
		for k,v := range view.Fields {
			s := reflect.ValueOf(v)
			switch s.Interface().(type){
			case int64:
				arr[k] = string(v.(int64))
			default:
				arr[k] = string(v.([]uint8))
			}
		}
>>>>>>> 219ac6dc41467e6dddf6f1dbfa707efdd8ecdf7b
		if err == nil {
<<<<<<< HEAD
			c.JSON(http.StatusOK, gin.H{
				"pl": arr,
				"selected": selected})
||||||| merged common ancestors
			// Call the render function with the title, article and the name of the
			// template
			c.SecureJSON(http.StatusOK, gin.H{
				"pl": view.Fields})
			log.Print(view.Fields)

=======
			// Call the render function with the title, article and the name of the
			// template
			c.JSON(http.StatusOK, gin.H{
				"pl": arr})
>>>>>>> 219ac6dc41467e6dddf6f1dbfa707efdd8ecdf7b
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		c.AbortWithStatus(http.StatusBadGateway)
	}
<<<<<<< HEAD
}

func viewRatings(c *gin.Context) {
	res, err := getReportData()
	if err != nil{
		log.Print(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"gov": res.Gov,
		"col_names": res.Header,
		})
}

func saveViewRow(c *gin.Context) {
	x, _ := ioutil.ReadAll(c.Request.Body)
	var srq saveRequest
	err := json.Unmarshal([]byte(x), &srq)
	if err != nil {
		log.Print(err.Error())
	} else {
		err := editView(srq.Name, srq.Data, srq.Id)
		if err != nil  {
			c.JSON(http.StatusBadRequest, gin.H{
				"ErrorTitle":   "Saving failed",
				"ErrorMessage": "Bad data"})
			log.Print(err.Error())

		} else {
			c.JSON(http.StatusOK, gin.H{
				"title": "Data saved"})
				log.Print("OK", srq.Name, srq.Data, srq.Id)
		}
	}
}

func  createView(c *gin.Context) {
	cols := make([]string, len(getColsNames()))
	formData := make([]interface{}, len(getColsNames()))
	colsNames := getColsNames()
	for i, _ := range cols{
		formData[i] = c.PostForm(colsNames[i])
	}
	if a, err := createNewView(formData); err == nil {
		render(c, gin.H{
			"title":   "Submission Successful",
			"pl": a}, "success.html")
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func allGovs(c *gin.Context) {
	res, err := getGovs()
	if err != nil{
		log.Print(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"govs": res,
	})
}
||||||| merged common ancestors
}
=======
}
>>>>>>> 219ac6dc41467e6dddf6f1dbfa707efdd8ecdf7b
