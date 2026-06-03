package handlers

import (
	"fmt"
	"net/http"
)

func Groups(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GROUP\n")
}
