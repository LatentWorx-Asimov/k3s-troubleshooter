package agent

import (
	"context"
	"log"
	"time"
)

// SystemPrompt enforces the architectural constraints and capabilities of this agent.
const SystemPrompt = `You are a specialized agent.
Your primary responsibility is: [Insert Responsibility].
Your strict constraints are: [Insert Constraints].`

// ProcessTask represents the Agent's internal cognitive loop.
// It receives a task context, does the work asynchronously, and returns a result string.
func ProcessTask(ctx context.Context, taskContext string) (string, error) {
	log.Printf("Agent starting task. Context: %s", taskContext)

	// Here you would instantiate your LLM client, inject the SystemPrompt,
	// and execute the agent's Tool/Thinking loop.
	// For demonstration, we simulate work with a sleep.
	time.Sleep(2 * time.Second)

	log.Println("Agent task processing completed.")
	return "Agent successfully completed the task based on constraints.", nil
}
