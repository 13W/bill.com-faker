package main

import (
	"github.com/go-chi/render"
	"net/http"
)

// POST /Login.json returns the login response
func (s *Rest) loginCtrl(w http.ResponseWriter, r *http.Request) {
	loginResponse := JSON{
		"response_status":  0,
		"response_message": "Success",
		"response_data": JSON{
			"apiEndPoint": "https://api-mock.bill.com/api/v2",
			"usersId":     RandStringBytes(20),
			"sessionId":   RandStringBytes(45),
			"orgId":       RandStringBytes(20),
		},
	}
	render.JSON(w, r, loginResponse)
}

// POST /Logout.json returns the logout response
func (s *Rest) logoutCtrl(w http.ResponseWriter, r *http.Request) {
	logoutResponse := JSON{
		"response_status": 0, "response_message": "Success", "response_data": JSON{},
	}
	render.JSON(w, r, logoutResponse)
}
