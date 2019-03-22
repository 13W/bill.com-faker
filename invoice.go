package main

import (
	"encoding/json"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"time"
)

type invoice struct {
	Id               string            `json:"id"`
	CustomerId       string            `json:"customerId"`
	InvoiceNumber    string            `json:"invoiceNumber"`
	InvoiceDate      string            `json:"invoiceDate"`
	DueDate          string            `json:"dueDate"`
	Description      string            `json:"description"`
	ActgClassId      string            `json:"actgClassId"`
	InvoiceLineItems []invoiceLineItem `json:"invoiceLineItem"`
}

type invoiceLineItem struct {
	ItemId           string `json:"itemId"`
	Quantity         int64  `json:"quantity"`
	Amount           int64  `json:"amount"`
	Price            int64  `json:"price"`
	RatePercent      int64  `json:"ratePercent"`
	ChartOfAccountId string `json:"chartOfAccountId"`
	DepartmentId     string `json:"departmentId"`
	LocationId       string `json:"locationId"`
	ActgClassId      string `json:"actgClassId"`
	JobId            string `json:"jobId"`
	Description      string `json:"description"`
	Taxable          bool   `json:"taxable"`
	TaxCode          string `json:"taxCode"`
}

func getMockInvoiceLineItems(item invoiceLineItem) JSON {
	return JSON{
		"entity":           "InvoiceLineItem",
		"itemId":           RandStringBytes(20),
		"quantity":         item.Quantity,
		"amount":           item.Amount,
		"price":            item.Price,
		"ratePercent":      1,
		"chartOfAccountId": RandStringBytes(20),
		"departmentId":     RandStringBytes(20),
		"locationId":       RandStringBytes(20),
		"actgClassId":      item.ActgClassId,
		"jobId":            RandStringBytes(20),
		"description":      item.Description,
		"taxable":          false,
		"taxCode":          "",
	}
}

func getMockInvoice(item invoice) *JSON {
	id := RandStringBytes(20)
	if item.Id != "" {
		id = item.Id
	}

	actgClassId := "00000000000000000000"
	if item.ActgClassId != "" {
		actgClassId = item.ActgClassId
	}

	modifyTime := time.Now().UTC().Format("2006-01-02T15:04:05.000-0700")
	invoiceLineItems := make([]JSON, 0)
	if len(item.InvoiceLineItems) > 0 {
		for _, invoice := range item.InvoiceLineItems {
			invoiceLineItems = append(invoiceLineItems, getMockInvoiceLineItems(invoice))
		}
	} else {
		invoiceLineItems = append(invoiceLineItems, getMockInvoiceLineItems(invoiceLineItem{
			Quantity:    200,
			Amount:      500,
			Price:       1234,
			ActgClassId: actgClassId,
			Description: item.Description,
		}))
	}

	amount := int64(0)
	for _, i := range invoiceLineItems {
		amount += i["amount"].(int64)
	}

	return &JSON{
		"entity":                "Invoice",
		"id":                    id,
		"isActive":              "1",
		"createdTime":           modifyTime,
		"updatedTime":           modifyTime,
		"customerId":            item.CustomerId,
		"invoiceNumber":         item.InvoiceNumber,
		"invoiceDate":           item.InvoiceDate,
		"dueDate":               item.DueDate,
		"glPostingDate":         "2016-12-12",
		"amount":                amount,
		"amountDue":             amount,
		"paymentStatus":         "1",
		"description":           item.Description,
		"poNumber":              "PO1818",
		"isToBePrinted":         true,
		"isToBeEmailed":         true,
		"lastSentTime":          nil,
		"itemSalesTax":          "00000000000000000000",
		"salesTaxPercentage":    0,
		"salesTaxTotal":         0.00,
		"terms":                 "Net 15",
		"salesRep":              "George",
		"FOB":                   nil,
		"shipDate":              "2016-12-20",
		"shipMethod":            nil,
		"departmentId":          RandStringBytes(20),
		"locationId":            RandStringBytes(20),
		"actgClassId":           actgClassId,
		"jobId":                 RandStringBytes(20),
		"payToBankAccountId":    RandStringBytes(20),
		"payToChartOfAccountId": RandStringBytes(20),
	}
}

var invoicesList []*JSON

// POST /List/Invoice.json returns an empty search response
func (s *Rest) invoiceListCtrl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, createResponse(invoicesList))
}

// POST /Crud/Read/Invoice.json returns a vendor attributes
func (s *Rest) invoiceReadCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var invoiceData JSON
	err := json.Unmarshal([]byte(r.FormValue("data")), &invoiceData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	invoiceItem := findInJSONList(invoiceData["id"].(string), invoicesList)

	render.JSON(w, r, createResponse(invoiceItem))
}

// POST /Crud/Update/Invoice.json returns an updated vendor attributes
func (s *Rest) invoiceUpdateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Invoice invoice `json:"obj"`
	}

	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	invoiceItem := findInJSONList(postData.Invoice.Id, invoicesList)
	*invoiceItem = *getMockInvoice(postData.Invoice)

	render.JSON(w, r, createResponse(invoiceItem))
}

// POST /Crud/Create/Invoice.json creates a vendor and returns a vendor attributes
func (s *Rest) invoiceCreateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Invoice invoice `json:"obj"`
	}
	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	postData.Invoice.Id = RandStringBytes(20)

	invoiceItem := getMockInvoice(postData.Invoice)
	invoicesList = append(invoicesList, invoiceItem)

	render.JSON(w, r, createResponse(invoiceItem))
}
