package helpers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ParseInt32UrlVar returns an int32 from an url variable
func ParseInt32UrlVar(r *http.Request, varName string) int32 {
	vars := mux.Vars(r)
	value := vars[varName]
	res, _ := strconv.ParseInt(value, 10, 32)

	return int32(res)
}
