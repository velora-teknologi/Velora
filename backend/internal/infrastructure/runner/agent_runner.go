package runner

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/infrastructure/executor"
)

// StartAgentWorker subscribes to the `agent.execute` subject and processes incoming events.
// This implementation performs a minimal, deterministic processing step and publishes
// either `agent.executed` or `agent.failed` with a JSON payload containing the outcome.
func StartAgentWorker(ctx context.Context, nc *nats.Conn, exec executor.Executor, logger *zap.SugaredLogger) (*nats.Subscription, error) {
	if nc == nil {
		logger.Warn("NATS connection is nil; agent worker will not start")
		return nil, nil
	}

	// Install lightweight connection handlers for observability and reconnection logging.
	nc.SetDisconnectHandler(func(nc *nats.Conn) {
		logger.Warnf("NATS disconnected: %v", nc.LastError())
	})
	nc.SetReconnectHandler(func(nc *nats.Conn) {
		logger.Infof("NATS reconnected to %s", nc.ConnectedUrl())
	})
	nc.SetClosedHandler(func(nc *nats.Conn) {
		logger.Infof("NATS connection closed")
	})

	// Try to subscribe with retries (exponential backoff) using a queue group so multiple
	// workers can share the load.
	var sub *nats.Subscription
	var err error
	maxAttempts := 5
	backoff := 100 * time.Millisecond
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		sub, err = nc.QueueSubscribe("agent.execute", "agent_workers", func(m *nats.Msg) {
			logger.Debugw("agent.execute received", "subject", m.Subject)

			var evt map[string]interface{}
			if err := json.Unmarshal(m.Data, &evt); err != nil {
				logger.Errorf("failed to unmarshal agent.execute payload: %v", err)
				// publish failure
				fail := map[string]interface{}{"status": "failed", "error": "invalid payload"}
				if b, err := json.Marshal(fail); err == nil {
					_ = nc.Publish("agent.failed", b)
				}
				return
			}

			// In the real runner, use the provided executor to run the agent.
			agentID, _ := evt["agent_id"].(string)
			userID, _ := evt["user_id"].(string)

			// extract config (may be string or object)
			var configStr string
			switch v := evt["config"].(type) {
			case string:
				configStr = v
			default:
				if b, err := json.Marshal(v); err == nil {
					configStr = string(b)
				}
			}

			if exec == nil {
				logger.Error("no executor configured; cannot run agent")
				fail := map[string]interface{}{"status": "failed", "error": "executor not configured"}
				if b, err := json.Marshal(fail); err == nil {
					_ = nc.Publish("agent.failed", b)
				}
				return
			}

			// execute and capture result
			res, err := exec.Execute(context.Background(), agentID, userID, configStr, evt["input"])
			if err != nil {
				logger.Errorf("executor failed: %v", err)
				fail := map[string]interface{}{"status": "failed", "error": err.Error()}
				if b, err := json.Marshal(fail); err == nil {
					_ = nc.Publish("agent.failed", b)
				}
				return
			}

			// attach timestamp and publish executed
			res["timestamp"] = time.Now().UTC().Format(time.RFC3339)
			payload, err := json.Marshal(res)
			if err != nil {
				logger.Errorf("failed to marshal execution result: %v", err)
				fail := map[string]interface{}{"status": "failed", "error": "marshal error"}
				if b, err := json.Marshal(fail); err == nil {
					_ = nc.Publish("agent.failed", b)
				}
				return
			}

			if err := nc.Publish("agent.executed", payload); err != nil {
				logger.Errorf("failed to publish agent.executed: %v", err)
			} else {
				logger.Infow("agent execution published", "agent_id", agentID, "user_id", userID)
			}
		})
		if err == nil {
			break
		}
		logger.Warnf("subscribe attempt %d failed: %v", attempt, err)
		time.Sleep(backoff)
		backoff *= 2
	}
	if err != nil {
		logger.Errorf("failed to subscribe to agent.execute after retries: %v", err)
		return nil, err
	}

	// ensure subscription is flushed
	_ = nc.Flush()

	logger.Info("agent worker started and subscribed to agent.execute (queue: agent_workers)")

	// Stop subscriber when context is done
	go func() {
		<-ctx.Done()
		logger.Info("context canceled — draining and unsubscribing agent worker")
		// drain allows in-flight messages to be processed
		if err := sub.Drain(); err != nil {
			logger.Warnf("error draining subscription: %v", err)
			_ = sub.Unsubscribe()
		}
		// flush to ensure pending protocol messages are sent
		_ = nc.Flush()
	}()

	return sub, nil
}
