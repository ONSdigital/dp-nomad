package client

import (
	"context"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
)

// ServiceName
const ServiceName = "Nomad"

var (
	// StatusMessage contains a map of messages to service response statuses
	StatusMessage = map[string]string{
		health.StatusOK:       " is ok",
		health.StatusWarning:  " is degraded, but at least partially functioning",
		health.StatusCritical: " functionality is unavailable or non-functioning",
	}
)

func (c *Client) Checker(ctx context.Context, state *health.CheckState) error {
	service := ServiceName
	logData := log.Data{
		"service": service,
	}

	code, err := c.Get(ctx, "/v1/agent/health?type=client")
	if err != nil {
		log.Event(ctx, "failed to request service health", log.ERROR, log.Error(err), logData)
		return err
	}

	switch code {
	case 0: // When there is a problem with the client return error in message
		state.Update(health.StatusCritical, err.Error(), 0)
	case 200:
		message := generateMessage(service, health.StatusOK)
		state.Update(health.StatusOK, message, code)
	default:
		message := generateMessage(service, health.StatusCritical)
		state.Update(health.StatusCritical, message, code)
	}

	return nil
}

func generateMessage(service string, state string) string {
	return service + StatusMessage[state]
}