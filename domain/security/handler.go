package security

import "net/http"

type SecurityHandler struct {
}

func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{}
}

func (h *SecurityHandler) LoginCustomer(w http.ResponseWriter, r *http.Request) {
	handler.LoginCustomer(a.DB, w, r)
}
