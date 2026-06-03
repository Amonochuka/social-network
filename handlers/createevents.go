

package handlers

import (
	"fmt"
	"net/http"
)

func  CreateEvents(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, " Create Events \n")
}
