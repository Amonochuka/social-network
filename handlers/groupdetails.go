package handlers

import (
	"fmt"
	"net/http"
)

func GroupDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "group details here\n")
}
