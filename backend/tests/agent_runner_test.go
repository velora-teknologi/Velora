package services

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/infrastructure/runner"
)

// mock executor that records calls
type mockExecutor struct {
	called chan map[string]interface{}
}

func (m *mockExecutor) Execute(ctx context.Context, agentID, userID string, config string, input interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{"agent_id": agentID, "user_id": userID, "status": "executed", "result": map[string]interface{}{"echo": input}}
	m.called <- res
	return res, nil
}

func runNATSServer(t *testing.T) *server.Server {
	opts := &server.Options{Host: "127.0.0.1", Port: -1}
	s, err := server.NewServer(opts)
	if err != nil {
		t.Fatalf("failed to create nats server: %v", err)
	}
	go s.Start()
	if !s.ReadyForConnections(5 * time.Second) {
		t.Fatalf("nats server not ready")
	}
	return s
}

func TestAgentRunner_Execute(t *testing.T) {
	// start embedded nats server
	s := runNATSServer(t)
	defer s.Shutdown()

	nc, err := nats.Connect(s.ClientURL())
	if err != nil {
		t.Fatalf("nats connect: %v", err)
	}
	defer nc.Close()

	z, _ := zap.NewDevelopment()
	logger := z.Sugar()

	m := &mockExecutor{called: make(chan map[string]interface{}, 1)}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = runner.StartAgentWorker(ctx, nc, m, logger)
	if err != nil {
		t.Fatalf("start agent worker: %v", err)
	}

	// subscribe to agent.executed capture
	sub, err := nc.SubscribeSync("agent.executed")
	if err != nil {
		t.Fatalf("subscribe executed: %v", err)
	}

	// publish an agent.execute event
	evt := map[string]interface{}{"agent_id": "a1", "user_id": "u1", "input": map[string]interface{}{"q": "hello"}}
	b, _ := json.Marshal(evt)
	if err := nc.Publish("agent.execute", b); err != nil {
		t.Fatalf("publish execute: %v", err)
	}

	// wait for executor called
	select {
	case res := <-m.called:
		if res["agent_id"] != "a1" {
			t.Fatalf("unexpected agent id: %v", res["agent_id"])
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("executor was not called in time")
	}

	// ensure agent.executed published
	msg, err := sub.NextMsg(3 * time.Second)
	if err != nil {
		t.Fatalf("did not receive executed message: %v", err)
	}
	var got map[string]interface{}
	_ = json.Unmarshal(msg.Data, &got)
	if got["agent_id"] != "a1" {
		t.Fatalf("unexpected executed payload agent_id: %v", got["agent_id"])
	}
}
