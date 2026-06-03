package handlers

import (
	"fmt"
	"net/http"
)

func RejectInvitation(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Reject Invitation\n")
}
