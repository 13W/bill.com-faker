package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jessevdk/go-flags"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type JSON map[string]interface{}

type Rest struct {
	Version    string
	lock       sync.Mutex
	httpServer *http.Server
}

type Opts struct {
	Port int `long:"port" env:"PORT" default:"8080" description:"port"`
}

const VERSION = "v2"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!-1234567890+_"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Run the lister and request's router, activate rest server
func (s *Rest) Run(hostUrl string, port int) {
	s.lock.Lock()
	router := s.routes()
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", hostUrl, port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	s.lock.Unlock()
	_ = s.httpServer.ListenAndServe()
}

// Shutdown the rest server
func (s *Rest) Shutdown() {
	s.httpServer.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[DEBUG] rest shutdown error, %s", err)
		}
	}
	log.Print("[INFO] shutdown rest server completed")
	s.lock.Unlock()
}

func (s *Rest) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/api/v2", func(rapi chi.Router) {
		// Login & Logout
		rapi.Group(func(r chi.Router) {
			r.Post("/Login.json", s.loginCtrl)
			r.Post("/Logout.json", s.logoutCtrl)
		})
		// Actg
		rapi.Group(func(r chi.Router) {
			r.Post("/List/ActgClass.json", s.actgClassListCtrl)
			r.Post("/Crud/Create/ActgClass.json", s.actgClassCreateCtrl)
			r.Post("/Crud/Read/ActgClass.json", s.actgClassReadCtrl)
			r.Post("/Crud/Update/ActgClass.json", s.actgClassUpdateCtrl)
		})
		// Vendor
		rapi.Group(func(r chi.Router) {
			r.Post("/List/Vendor.json", s.vendorListCtrl)
			r.Post("/Crud/Create/Vendor.json", s.vendorCreateCtrl)
			r.Post("/Crud/Read/Vendor.json", s.vendorReadCtrl)
			r.Post("/Crud/Update/Vendor.json", s.vendorUpdateCtrl)
		})
		// Customer
		rapi.Group(func(r chi.Router) {
			r.Post("/List/Customer.json", s.customerListCtrl)
			r.Post("/Crud/Create/Customer.json", s.customerCreateCtrl)
			r.Post("/Crud/Read/Customer.json", s.customerReadCtrl)
			r.Post("/Crud/Update/Customer.json", s.customerUpdateCtrl)
		})
		// Invoice
		rapi.Group(func(r chi.Router) {
			r.Post("/List/Invoice.json", s.invoiceListCtrl)
			r.Post("/Crud/Create/Invoice.json", s.invoiceCreateCtrl)
			r.Post("/Crud/Read/Invoice.json", s.invoiceReadCtrl)
			r.Post("/Crud/Update/Invoice.json", s.invoiceUpdateCtrl)
		})
		// Bill
		rapi.Group(func(r chi.Router) {
			r.Post("/List/Bill.json", s.billListCtrl)
			r.Post("/Crud/Read/Bill.json", s.billReadCtrl)
			r.Post("/Crud/Create/Bill.json", s.billCreateCtrl)
			r.Post("/Crud/Delete/Bill.json", s.billReadCtrl)
			r.Post("/Crud/Update/Bill.json", s.billUpdateCtrl)
		})
		// RecurringBill
		rapi.Group(func(r chi.Router) {
			r.Post("/List/RecurringBill.json", s.recurringBillListCtrl)
			r.Post("/Crud/Read/RecurringBill.json", s.recurringBillReadCtrl)
			r.Post("/Crud/Create/RecurringBill.json", s.recurringBillCreateCtrl)
			r.Post("/Crud/Delete/RecurringBill.json", s.recurringBillReadCtrl)
			r.Post("/Crud/Update/RecurringBill.json", s.recurringBillUpdateCtrl)
		})
	})

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		render.HTML(w, r, "pong")
	})
	return router
}

func findInJSONList(id string, list []*JSON) *JSON {
	for _, item := range list {
		if (*item)["id"].(string) == id {
			return item
		}
	}

	return nil
}

func getJSONElement(id int64, list []*JSON) JSON {
	return *list[id]
}

func createResponse(data interface{}) JSON {
	status := 0
	message := "Success"

	if data == nil {
		status = 1
		message = "Error"
	}

	return JSON{
		"response_status":  status,
		"response_message": message,
		"response_data":    data,
	}
}

func fillLists() {
	actgClassList = append(actgClassList, getMockActgClass(actgClass{
		Name: "MockActgClasName",
	}))

	customersList = append(customersList, getMockCustomer(customer{
		Id:          RandStringBytes(20),
		Description: "First Customer",
		Phone:       "+1234567890",
		Email:       "customer+1@example.com",
		CompanyName: "ASD Inc,",
		Name:        "First Customer",
		ShortName:   "First",
	}))

	customersList = append(customersList, getMockCustomer(customer{
		Id:          RandStringBytes(20),
		Description: "Second Customer",
		Phone:       "+1234567890",
		Email:       "customer+2@example.com",
		CompanyName: "ASD Inc,",
		Name:        "Second Customer",
		ShortName:   "Second",
	}))

	vendorsList = append(vendorsList, getMockVendor(vendor{
		Name:           "Customer 1",
		Email:          "test@vemdor.com",
		Phone:          "+1234567890",
		Address1:       "In the middle of the world",
		AddressCity:    "Neverlands",
		AddressCountry: "Neveros",
		AddressState:   "Neverland",
		ID:             RandStringBytes(20),
	}))

	vendorsList = append(vendorsList, getMockVendor(vendor{
		Name:           "Customer 2",
		Email:          "test+2@vemdor.com",
		Phone:          "+1234567890",
		Address1:       "In the middle of the world",
		AddressCity:    "Neverlands",
		AddressCountry: "Neveros",
		AddressState:   "Neverland",
		ID:             RandStringBytes(20),
	}))

	billId := RandStringBytes(20)
	billList = append(billList, getMockBill(bill{
		ID:            billId,
		Description:   "Bill some description",
		InvoiceDate:   "2010-01-01",
		InvoiceNumber: "Test_Jan_01",
		DueDate:       "2010-01-01",
		VendorId:      getJSONElement(0, vendorsList)["id"].(string),
		CustomerId:    getJSONElement(0, customersList)["id"].(string),
		BillLineItems: []billLineItem{
			{
				ID:               RandStringBytes(20),
				ChartOfAccountId: RandStringBytes(20),
				BillId:           billId,
				Amount:           float64(rand.Intn(1000)),
				ActgClassId:      getJSONElement(0, actgClassList)["id"].(string),
				Quantity:         1,
			},
		},
	}))
	billId = RandStringBytes(20)
	billList = append(billList, getMockBill(bill{
		ID:            billId,
		Description:   "Bill second description",
		InvoiceDate:   "2010-01-01",
		InvoiceNumber: "Test_Jan_01",
		DueDate:       "2010-01-01",
		VendorId:      getJSONElement(1, vendorsList)["id"].(string),
		CustomerId:    getJSONElement(1, customersList)["id"].(string),
		BillLineItems: []billLineItem{
			{
				ID:               RandStringBytes(20),
				ChartOfAccountId: RandStringBytes(20),
				BillId:           billId,
				Amount:           float64(rand.Intn(1000)),
				ActgClassId:      getJSONElement(0, actgClassList)["id"].(string),
				Quantity:         1,
			},
		},
	}))

	invoicesList = append(invoicesList, getMockInvoice(invoice{
		ActgClassId: getJSONElement(0, actgClassList)["id"].(string),
		Description: "Test Invoice Description",
	}))
}

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		os.Exit(1)
	}

	fillLists()

	restSrv := &Rest{
		Version: VERSION,
	}
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		restSrv.Shutdown()
	}()
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	log.Printf("[INFO] Start the mock.bill.com server http://0.0.0.0:%d", opts.Port)
	restSrv.Run("0.0.0.0", opts.Port)
}
