package main

import (
	"encoding/json"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"time"
)

type recurringBill struct {
	Entity                 string                  `json:"entity"`
	Id                     string                  `json:"id"`
	IsActive               string                  `json:"isActive"`
	VendorId               string                  `json:"vendorId"`
	TimePeriod             string                  `json:"timePeriod"`
	FrequencyPerTimePeriod int64                   `json:"frequencyPerTimePeriod"`
	NextDueDate            string                  `json:"nextDueDate"`
	EndDate                string                  `json:"endDate"`
	DaysInAdvance          int64                   `json:"daysInAdvance"`
	Description            string                  `json:"description"`
	CreatedTime            string                  `json:"createdTime"`
	UpdatedTime            string                  `json:"updatedTime"`
	RecurringBillLineItems []recurringBillLineItem `json:"recurringBillLineItems"`
}

type recurringBillLineItem struct {
	Entity           string  `json:"entity"`
	Id               string  `json:"id"`
	RecurringBillId  string  `json:"recurringBillId"`
	Amount           float64 `json:"amount"`
	ChartOfAccountId string  `json:"chartOfAccountId"`
	DepartmentId     string  `json:"departmentId"`
	LocationId       string  `json:"locationId"`
	Description      string  `json:"description"`
	CreatedTime      string  `json:"createdTime"`
	UpdatedTime      string  `json:"updatedTime"`
}

func getMockRecurringBill(item recurringBill) *JSON {
	modifyTime := time.Now().UTC().Format("2006-01-02T15:04:05.000-0700")
	recurringBillId := RandStringBytes(20)
	if item.Id != "" {
		recurringBillId = item.Id
	}
	recurringBillLineItems := make([]JSON, 0)
	recurringBillLineItems = append(recurringBillLineItems, JSON{
		"entity":           "RecurringBillLineItem",
		"id":               RandStringBytes(20),
		"recurringBillId":  recurringBillId,
		"amount":           200.00,
		"chartOfAccountId": RandStringBytes(20),
		"departmentId":     RandStringBytes(20),
		"locationId":       RandStringBytes(20),
		"description":      nil,
		"createdTime":      modifyTime,
		"updatedTime":      modifyTime,
	})
	return &JSON{
		"entity":                 "RecurringBill",
		"id":                     recurringBillId,
		"isActive":               "1",
		"vendorId":               item.VendorId,
		"timePeriod":             "2",
		"frequencyPerTimePeriod": item.FrequencyPerTimePeriod,
		"nextDueDate":            item.NextDueDate,
		"endDate":                "2017-12-05",
		"daysInAdvance":          item.DaysInAdvance,
		"description":            "Services and deliveries.",
		"createdTime":            modifyTime,
		"updatedTime":            modifyTime,
		"recurringBillLineItems": recurringBillLineItems,
	}

}

var recurringBillsList []*JSON

// POST /List/RecurringBill.json returns an empty search response
func (s *Rest) recurringBillListCtrl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, createResponse(recurringBillsList))
}

// POST /Crud/Read/RecurringBill.json returns a vendor attributes
func (s *Rest) recurringBillReadCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var recurringBillData JSON
	err := json.Unmarshal([]byte(r.FormValue("data")), &recurringBillData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	recurringBillItem := *findInJSONList(recurringBillData["id"].(string), recurringBillsList)

	render.JSON(w, r, createResponse(recurringBillItem))
}

// POST /Crud/Update/RecurringBill.json returns an updated vendor attributes
func (s *Rest) recurringBillUpdateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		RecurringBill recurringBill `json:"obj"`
	}

	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	recurringBillItem := findInJSONList(postData.RecurringBill.Id, recurringBillsList)
	*recurringBillItem = *getMockRecurringBill(postData.RecurringBill)

	render.JSON(w, r, createResponse(recurringBillItem))
}

// POST /Crud/Create/RecurringBill.json creates a vendor and returns a vendor attributes
func (s *Rest) recurringBillCreateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		RecurringBill recurringBill `json:"obj"`
	}
	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	postData.RecurringBill.Id = RandStringBytes(20)

	recurringBillItem := getMockRecurringBill(postData.RecurringBill)
	recurringBillsList = append(recurringBillsList, recurringBillItem)

	render.JSON(w, r, createResponse(recurringBillItem))
}
