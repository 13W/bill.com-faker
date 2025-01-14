package main

import (
	"encoding/json"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"time"
)

type billLineItem struct {
	ID               string  `json:"id"`
	ChartOfAccountId string  `json:"chartOfAccountId"`
	BillId           string  `json:"billId"`
	Amount           float64 `json:"amount"`
	ActgClassId      string  `json:"actgClassId"`
	Quantity         int32   `json:"quantity"`
}

type bill struct {
	ID            string         `json:"id"`
	VendorId      string         `json:"vendorId"`
	Description   string         `json:"description"`
	InvoiceDate   string         `json:"invoiceDate"`
	InvoiceNumber string         `json:"invoiceNumber"`
	DueDate       string         `json:"dueDate"`
	BillLineItems []billLineItem `json:"billLineItems"`
	CustomerId    string         `json:"customerId"`
}

var billList []*JSON

func getMockBill(item bill) *JSON {
	modifyTime := time.Now().UTC().Format("2006-01-02T15:04:05.000-0700")
	billLineItems := make([]JSON, 0)
	billLineItems = append(billLineItems, JSON{
		"itemId":           "00000000000000000000",
		"updatedTime":      modifyTime,
		"description":      "",
		"chartOfAccountId": item.BillLineItems[0].ChartOfAccountId,
		"billId":           item.ID,
		"entity":           "BillLineItem",
		"customerId":       item.CustomerId,
		"employeeId":       "00000000000000000000",
		"amount":           item.BillLineItems[0].Amount,
		"locationId":       "00000000000000000000",
		"departmentId":     "00000000000000000000",
		"lineType":         "1",
		"jobBillable":      false,
		"createdTime":      modifyTime,
		"actgClassId":      item.BillLineItems[0].ActgClassId,
		"jobId":            "00000000000000000000",
		"unitPrice":        "",
		"id":               item.BillLineItems[0].ID,
		"quantity":         item.BillLineItems[0].Quantity,
	})
	return &JSON{
		"vendorId":                item.VendorId,
		"invoiceDate":             item.InvoiceDate,
		"dueAmount":               item.BillLineItems[0].Amount,
		"entity":                  "Bill",
		"paymentTermId":           "00000000000000000000",
		"hasAutoPay":              false,
		"paidAmount":              "",
		"dueDate":                 item.DueDate,
		"localAmount":             "",
		"glPostingDate":           "",
		"approvalStatus":          "0",
		"id":                      item.ID,
		"poNumber":                "",
		"billLineItems":           billLineItems,
		"payFromBankAccountId":    "00000000000000000000",
		"description":             item.Description,
		"exchangeRate":            "",
		"invoiceNumber":           item.InvoiceNumber,
		"isActive":                "1",
		"updatedTime":             modifyTime,
		"eBillCreated":            false,
		"scheduledAmount":         0.0,
		"paymentStatus":           "1",
		"amount":                  item.BillLineItems[0].Amount,
		"createdTime":             modifyTime,
		"payFromChartOfAccountId": "00000000000000000000",
	}
}

// POST /List/Bill.json returns an empty search response
func (s *Rest) billListCtrl(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, createResponse(billList))
}

// POST /Crud/Read/Bill.json returns a bill attributes as well as bill line item attributes
func (s *Rest) billReadCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		ID string `json:"id"`
	}
	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}

	billItem := findInJSONList(postData.ID, billList)

	render.JSON(w, r, createResponse(billItem))
}

// POST /Crud/Create/Bill.json creates a bill
func (s *Rest) billCreateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Bill bill `json:"obj"`
	}

	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	postData.Bill.ID = RandStringBytes(20)
	postData.Bill.BillLineItems[0].ID = RandStringBytes(20)

	billItem := getMockBill(postData.Bill)
	billList = append(billList, billItem)

	render.JSON(w, r, createResponse(billItem))
}

// POST /Crud/Update/Bill.json or /Crud/Delete/Bill.json returns a bill attributes
func (s *Rest) billUpdateCtrl(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}
	var postData struct {
		Bill bill `json:"obj"`
	}

	err := json.Unmarshal([]byte(r.FormValue("data")), &postData)
	if err != nil {
		log.Println("[ERROR] Can't parse request data")
		return
	}

	postData.Bill.BillLineItems[0].ID = RandStringBytes(20)

	billItem := findInJSONList(postData.Bill.ID, billList)
	*billItem = *getMockBill(postData.Bill)

	render.JSON(w, r, createResponse(billItem))
}
