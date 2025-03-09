package HttpServer

import (
	"fmt"
	HTTP "net/http"
)

type Response = HTTP.ResponseWriter
type Request = *HTTP.Request

var ADD_NUMBER int = 0

func h_NotFound(response Response, request Request) {
	response.WriteHeader(404)

	fmt.Fprintln(response, "Not Found")
}

func h_Add(response Response, request Request) {
	defer fmt.Fprintln(response, "new add is : ", ADD_NUMBER)

	ADD_NUMBER++
}
