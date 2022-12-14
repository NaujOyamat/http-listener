package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	PORT = flag.Int("port", 8080, "Server port")
)

func main() {
	flag.Parse()

	gin.SetMode(gin.DebugMode)

	server := gin.Default()

	server.Any("/*path", handler())

	if err := server.Run(fmt.Sprintf(":%d", *PORT)); err != nil {
		log.Fatal(err.Error())
	}
}

func handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var response response

		response.Method = ctx.Request.Method

		for k, v := range ctx.Request.Header {
			response.Headers = append(response.Headers, fmt.Sprintf("%s: %s", k, v))
		}

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Error(fmt.Errorf("can't read body"))
			return
		}

		response.Body = string(body)

		data, _ := json.Marshal(response)

		log.Printf("[RESPONSE] %s\n", data)

		ctx.JSON(http.StatusOK, response)
	}
}

type response struct {
	Method  string
	Headers []string
	Body    string
}
