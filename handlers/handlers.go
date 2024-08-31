package handlers

import (
	"net/http"

	"github.com/Mohamed-khattab/Message-Queue/messaging"
	"github.com/labstack/echo"
)

type GenericResponse struct {
	Error   string            `json:"error"`
	Data    map[string]string `json:"data"`
	Message string            `json:"message"`
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

	return c.String(http.StatusOK, "Message published successfully")
}

func (h *Handlers) Retrieve(c echo.Context) error {

	return c.JSON(http.StatusOK, []string{"message1", "message2"})
}
