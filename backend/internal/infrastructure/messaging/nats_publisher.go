package messaging

import (
	"github.com/nats-io/nats.go"
	"github.com/velora-teknologi/velora/internal/application/services"
	"go.uber.org/zap"
)

type NatsPublisher struct {
	nc     *nats.Conn
	logger *zap.SugaredLogger
}

func NewNatsPublisher(nc *nats.Conn, logger *zap.SugaredLogger) services.Publisher {
	return &NatsPublisher{nc: nc, logger: logger}
}

func (p *NatsPublisher) Publish(subject string, data []byte) error {
	if p == nil || p.nc == nil {
		if p != nil && p.logger != nil {
			p.logger.Warnw("NATS connection not configured; dropping publish", "subject", subject)
		}
		return nil
	}
	return p.nc.Publish(subject, data)
}
