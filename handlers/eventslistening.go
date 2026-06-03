

package handlers

import (
	"fmt"
	"net/http"
)

func  EventsListening(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, " Events  Listening \n")
}
