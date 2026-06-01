package executor

import "context"

type Executor interface {
	// Execute runs the agent logic and returns a result map or an error
	Execute(ctx context.Context, agentID, userID string, config string, input interface{}) (map[string]interface{}, error)
}
