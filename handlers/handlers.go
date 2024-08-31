package handlers

import (
	"net/http"

	"github.com/Mohamed-khattab/Message-Queue/messaging"
	"github.com/labstack/echo"
)

type Handlers struct {
	Broker *messaging.Broker
}

func NewHandlers(broker *messaging.Broker) *Handlers {
	return &Handlers{
		Broker: broker,
	}
}

func (h *Handlers) Subscribe(c echo.Context) error {

	var req messaging.SubscriptionRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request Format")
	}

	success , err  := h.Broker.Subscribe(req.Endpoint, req.Topics)
	if err != nil {
    	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !success {
		return c.JSON(http.StatusInternalServerError,"Failed to subscribe to specified topics")
	}

	return c.String(http.StatusOK, "Subscribed successfully")
}

func (h *Handlers) Unsubscribe(c echo.Context) error {

	var req messaging.SubscriptionRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request Format")
	}

	success := h.Broker.Unsubscribe(req.Topics)
	if !success {
		return c.JSON(http.StatusInternalServerError, "Failed to unsubscribe from specified topics")
	}
	
	return c.String(http.StatusOK, "Unsubscribed successfully")
}

func (h *Handlers) Publish(c echo.Context) error {
	

	return c.String(http.StatusOK, "Message published successfully")
}

func (h *Handlers) Retrieve(c echo.Context) error {

	return c.JSON(http.StatusOK, []string{"message1", "message2"})
}
