package main

import (
	"fmt"
	"log"

	"github.com/Mohamed-khattab/Message-Queue/handlers"
	"github.com/Mohamed-khattab/Message-Queue/messaging"
	"github.com/labstack/echo"
)

func main() {

	e := echo.New()

	var broker = messaging.NewBroker()
	var handler = handlers.NewHandlers(broker)

	e.POST("/v1/subescribe", handler.Subscribe)
	e.POST("/v1/unsubescribe", handler.Unsubscribe)

	e.POST("/v1/publish", handler.Subscribe)
	e.GET("/v1/retrieve", handler.Unsubscribe)

	// START THE SERVER
	fmt.Println("Server is running at port 3000")
	log.Fatal(e.Start(":3000"))
}
