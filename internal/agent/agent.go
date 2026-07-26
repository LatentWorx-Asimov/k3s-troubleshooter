package agent

import (
	"context"
	"fmt"
	"log"
)

// SystemPrompt enforces the architectural constraints and capabilities of this agent.
const SystemPrompt = `You are a k3s-troubleshooter agent.
Your primary responsibility is to diagnose k3s infrastructure, CrashLoopBackOffs, Ingress misconfigurations, and networking issues.
Your strict constraints are:
1. Use a fail-fast, layered diagnostic approach.
2. Rely only on allowed kubectl tools (get, describe, logs).
3. Do not attempt to modify resources directly unless authorized, focus on diagnosis.`

// ProcessTask represents the Agent's internal cognitive loop.
// It receives a task context, does the work asynchronously, and returns a result string.
func ProcessTask(ctx context.Context, taskContext string) (string, error) {
	log.Printf("Agent starting task. Context: %s", taskContext)

	// Since we don't have a real LLM framework integrated in this scaffold,
	// we simulate the tool loop and provide a static diagnostic response.
	// In a real implementation, the LLM would iterate using the tools.
	
	// Test one of the tools to ensure it works
	res, err := KubectlTool(ctx, "get", "pods", "-A")
	if err != nil {
		log.Printf("Failed to get pods: %v", err)
	} else {
		log.Printf("Found pods: %s", res[:min(len(res), 100)])
	}

	log.Println("Agent task processing completed.")
	return fmt.Sprintf("Diagnosis complete for task: %s", taskContext), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
