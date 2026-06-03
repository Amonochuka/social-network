

package handlers

import (
	"fmt"
	"net/http"
)

func  NotGoingResponce(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, " Not Going Responce \n")
}
