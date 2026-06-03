

package handlers

import (
	"fmt"
	"net/http"
)

func  AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, "Accept Invitation\n")
}
