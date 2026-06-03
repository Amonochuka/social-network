

package handlers

import (
	"fmt"
	"net/http"
)

func  Comments(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, "Comments\n")
}
