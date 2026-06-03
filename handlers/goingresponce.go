

package handlers

import (
	"fmt"
	"net/http"
)

func  GoingResponce(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, " Going Responce \n")
}
