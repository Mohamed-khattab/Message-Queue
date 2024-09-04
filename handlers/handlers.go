package handlers

import (
	"fmt"
	"net/http"

	"github.com/Mohamed-khattab/Message-Queue/messaging"
	"github.com/labstack/echo"
)

type GenericResponse struct {
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type Handlers struct {
	Broker *messaging.Broker
}

func NewHandlers(broker *messaging.Broker) *Handlers {
	return &Handlers{
		Broker: broker,
	}
}

func (h *Handlers) Subscribe(c echo.Context) error {

	var req messaging.SubscribeRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid request Format",
		})
	}

	subscriberID, err := h.Broker.Subscribe(req.Endpoint, req.Topics)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &GenericResponse{
		Data: map[string]string{
			"subscriber_id": subscriberID,
		},
		Message: "Subscribed successfully",
	})
}
func (h *Handlers) Unsubscribe(c echo.Context) error {

	var req messaging.UnsubscribeRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid Request Format",
		})
	}

	err := h.Broker.Unsubscribe(req.SubscriberID, req.Topics)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &GenericResponse{
		Message: "Unsubscribed successfully to the Specified Topics",
	})

}

func (h *Handlers) Publish(c echo.Context) error {

	var req messaging.PublishRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid Request Format",
		})
	}

	err := h.Broker.Publish(req.Topic, req.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &GenericResponse{
		Message: "Message published successfully to all Topic subscribers",
	})
}

func (h *Handlers) Retrieve(c echo.Context) error {

	var req messaging.RetrieveRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &GenericResponse{
			Error: "Invalid Request Format",
		})
	}

	messages, err := h.Broker.Retrieve(req.SubscriberId, req.Topic, req.StartDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &GenericResponse{
			Error: fmt.Sprintf("Failed to retrieve messages with error %v", err),
		})
	}
	return c.JSON(http.StatusOK, &GenericResponse{
		Data: messages,
	})
}
