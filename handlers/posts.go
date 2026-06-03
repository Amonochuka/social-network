

package handlers

import (
	"fmt"
	"net/http"
)

func  Posts(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, "Posts \n")
}
