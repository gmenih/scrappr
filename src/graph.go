package p

import (
	"fmt"
	"net/http"
)

func DrawGraphs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, World!")
}
