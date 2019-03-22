package main

import (
	"encoding/json"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"time"
)

func getMockVendor(item vendor) *JSON {
	modifyTime := time.Now().UTC().Format("2006-01-02T15:04:05.000-0700")
	return &JSON{
		"sendNotifications":      true,
		"mergedIntoId":           "00000000000000000000",
		"taxId":                  "",
		"entity":                 "Vendor",
		"paymentEmail":           "",
		"paymentTermId":          "00000000000000000000",
		"hasBankAccountAutoPay":  false,
		"paymentCurrency":        "",
		"bankCountry":            "",
		"id":                     item.ID,
		"accNumber":              "",
		"paymentPhone":           item.Phone,
		"paymentPurpose":         "",
		"since":                  "",
		"payDaysBefore":          "",
		"addressCountry":         item.AddressCountry,
		"billSyncPref":           "0",
		"isActive":               "1",
		"track1099":              true,
		"nameOnCheck":            "",
		"fax":                    "",
		"description":            "",
		"billCurrency":           "",
		"companyName":            "",
		"payBy":                  "0",
		"address1":               item.Address1,
		"address2":               "",
		"address3":               "",
		"address4":               "",
		"phone":                  item.Phone,
		"contactLastName":        "",
		"accountType":            "0",
		"enabledCombinePayments": true,
		"shortName":              "",
		"externalBillPayIn12m":   "",
		"addressCity":            item.AddressCity,
		"updatedTime":            modifyTime,
		"contactFirstName":       "",
		"name":                   item.Name,
		"addressState":           item.AddressState,
		"email":                  item.Email,
		"lastBalanceUpdate":      modifyTime,
		"createdTime":            modifyTime,
		"balance":                100000.0,
		"prefPmtMethod":          "1",
		"addressZip":             "11111",
	}
}

type vendor struct {
	ID             string `json:"id"`
	AddressCountry string `json:"addressCountry"`
	AddressCity    string `json:"addressCity"`
	AddressState   string `json:"addressState"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Name           string `json:"name"`
	Address1       string `json:"address1"`
}

var vendorsList []*JSON

// POST /List/Vendor.json returns an empty search response
func (s *Rest) vendorListCtrl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, createResponse(vendorsList))
}

// POST /Crud/Read/Vendor.json returns a vendor attributes
func (s *Rest) vendorReadCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var vendorData JSON
	err := json.Unmarshal([]byte(r.FormValue("data")), &vendorData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}

	vendorItem := findInJSONList(vendorData["id"].(string), vendorsList)

	render.JSON(w, r, createResponse(vendorItem))
}

// POST /Crud/Update/Vendor.json returns an updated vendor attributes
func (s *Rest) vendorUpdateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Vendor vendor `json:"obj"`
	}

	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}

	vendorItem := findInJSONList(postData.Vendor.ID, vendorsList)
	*vendorItem = *getMockVendor(postData.Vendor)

	render.JSON(w, r, createResponse(vendorItem))
}

// POST /Crud/Create/Vendor.json creates a vendor and returns a vendor attributes
func (s *Rest) vendorCreateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Vendor vendor `json:"obj"`
	}
	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	postData.Vendor.ID = RandStringBytes(20)

	vendorItem := getMockVendor(postData.Vendor)
	vendorsList = append(vendorsList, vendorItem)
	vendorResponse := JSON{
		"response_status":  0,
		"response_message": "Success",
		"response_data":    vendorItem,
	}

	render.JSON(w, r, vendorResponse)
}
