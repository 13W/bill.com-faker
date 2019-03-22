package main

import (
	"encoding/json"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type customer struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	CompanyName string `json:"companyName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
}

func getMockCustomer(item customer) *JSON {
	return &JSON{
		"entity":                "Customer",
		"id":                    RandStringBytes(20),
		"isActive":              "1",
		"createdTime":           "2016-12-09T21:38:40.000+0000",
		"updatedTime":           "2016-12-09T21:38:40.000+0000",
		"name":                  item.Name,
		"shortName":             item.ShortName,
		"parentCustomerId":      "00000000000000000000",
		"companyName":           item.CompanyName,
		"contactFirstName":      "George",
		"contactLastName":       "Smith",
		"accNumber":             "212",
		"billAddress1":          "123 South North Street",
		"billAddress2":          "Suite 123",
		"billAddress3":          nil,
		"billAddress4":          nil,
		"billAddressCity":       "Santa Clara",
		"billAddressState":      "CA",
		"billAddressCountry":    "USA",
		"billAddressZip":        "95051",
		"shipAddress1":          "6868 First Second Avenue",
		"shipAddress2":          "Office 10",
		"shipAddress3":          nil,
		"shipAddress4":          nil,
		"shipAddressCity":       "Saint George",
		"shipAddressState":      "AL",
		"shipAddressCountry":    "USA",
		"shipAddressZip":        "35000",
		"email":                 item.Email,
		"phone":                 item.Phone,
		"altPhone":              "123-123-1256",
		"fax":                   "123-123-1235",
		"description":           item.Description,
		"printAs":               nil,
		"mergedIntoId":          "00000000000000000000",
		"hasAuthorizedToCharge": false,
		"accountType":           "1",
	}
}

var customersList []*JSON

// POST /List/Customer.json returns an empty search response
func (s *Rest) customerListCtrl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, createResponse(customersList))
}

// POST /Crud/Read/Customer.json returns a vendor attributes
func (s *Rest) customerReadCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var customerData JSON
	err := json.Unmarshal([]byte(r.FormValue("data")), &customerData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	customerItem := findInJSONList(customerData["id"].(string), customersList)

	render.JSON(w, r, createResponse(customerItem))
}

// POST /Crud/Update/Customer.json returns an updated vendor attributes
func (s *Rest) customerUpdateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Customer customer `json:"obj"`
	}

	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	customerItem := findInJSONList(postData.Customer.Id, customersList)
	*customerItem = *getMockCustomer(postData.Customer)

	render.JSON(w, r, createResponse(customerItem))
}

// POST /Crud/Create/Customer.json creates a vendor and returns a vendor attributes
func (s *Rest) customerCreateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Customer customer `json:"obj"`
	}
	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	postData.Customer.Id = RandStringBytes(20)

	customerItem := getMockCustomer(postData.Customer)
	customersList = append(customersList, customerItem)

	render.JSON(w, r, createResponse(customerItem))
}
