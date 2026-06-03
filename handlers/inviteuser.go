

package handlers

import (
	"fmt"
	"net/http"
)

func  InviteUsers(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, "Invite Users\n")
}
