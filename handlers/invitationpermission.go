package handlers

import (
	"fmt"
	"net/http"
)

func InvitationPermission(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Invitation Permission \n")
}
