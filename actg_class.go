package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

var actgClassList []*JSON

type actgClass struct {
	Id       string `json:"id"`
	ParentId string `json:"parentId"`
	Name     string `json:"name"`
}

func getMockActgClass(item actgClass) *JSON {
	id := fmt.Sprintf("cls%s", RandStringBytes(17))
	parentId := fmt.Sprintf("cls%s", RandStringBytes(17))
	if item.Id != "" {
		id = item.Id
	}
	if item.ParentId != "" {
		parentId = item.ParentId
	}

	return &JSON{
		"updatedTime":       "2019-01-30T08:05:20.000+0000",
		"parentActgClassId": parentId,
		"name":              item.Name,
		"mergedIntoId":      "00000000000000000000",
		"entity":            "ActgClass",
		"createdTime":       "2019-01-30T08:05:20.000+0000",
		"shortName":         "",
		"id":                id,
		"isActive":          "1",
		"description":       "",
	}
}

// POST /List/ActgClass.json returns an empty search response
func (s *Rest) actgClassListCtrl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, createResponse(actgClassList))
}

// POST /Crud/Read/ActgClass.json returns a vendor attributes
func (s *Rest) actgClassReadCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var actgClassData JSON
	err := json.Unmarshal([]byte(r.FormValue("data")), &actgClassData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	actgClassItem := findInJSONList(actgClassData["id"].(string), actgClassList)

	render.JSON(w, r, createResponse(actgClassItem))
}

// POST /Crud/Update/ActgClass.json returns an updated vendor attributes
func (s *Rest) actgClassUpdateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		ActgClass actgClass `json:"obj"`
	}

	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	actgClassItem := findInJSONList(postData.ActgClass.Id, actgClassList)
	*actgClassItem = *getMockActgClass(postData.ActgClass)

	render.JSON(w, r, createResponse(actgClassItem))
}

// POST /Crud/Create/ActgClass.json creates a vendor and returns a vendor attributes
func (s *Rest) actgClassCreateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		ActgClass actgClass `json:"obj"`
	}
	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	postData.ActgClass.Id = RandStringBytes(20)

	actgClassItem := getMockActgClass(postData.ActgClass)
	actgClassList = append(actgClassList, actgClassItem)

	render.JSON(w, r, createResponse(actgClassItem))
}
