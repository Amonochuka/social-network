

package handlers

import (
	"fmt"
	"net/http"
)

func  Membership(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, "Membership \n")
}
