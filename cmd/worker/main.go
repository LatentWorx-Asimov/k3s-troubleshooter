package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/LatentWorx-Asimov/k3s-troubleshooter/internal/agent"
)

// A2AEnvelope defines the standard Agent-to-Agent messaging structure
type A2AEnvelope struct {
	Sender  string `json:"sender"`
	Intent  string `json:"intent"`
	ReplyTo string `json:"reply_to,omitempty"`
	Payload struct {
		Context string `json:"context"`
	} `json:"payload"`
}

// A2AResponse defines the standard callback format
type A2AResponse struct {
	Sender  string `json:"sender"`
	Status  string `json:"status"`
	Intent  string `json:"intent"`
	Payload struct {
		Result string `json:"result"`
	} `json:"payload"`
}

type AgentSpec struct {
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	AllowedTools        []string `json:"allowed_tools"`
	PreferredModel      string   `json:"preferred_model"`
	FallbackModel       string   `json:"fallback_model"`
	ConfidenceThreshold float64  `json:"confidence_threshold"`
	Version             string   `json:"version"`
	CostPerToken        float64  `json:"cost_per_token"`
}

func main() {
	log.Println("Starting A2A Agent Webhook Listener...")

	agentSpec := AgentSpec{
		Name:                "k3s-troubleshooter",
		Description:         "Diagnoses k3s infrastructure, CrashLoopBackOffs, Ingress misconfigurations, and networking issues.",
		AllowedTools:        []string{"kubectl_get_pods", "kubectl_describe", "kubectl_logs", "kubectl_get_events", "kubectl_get_services", "kubectl_get_ingress"},
		PreferredModel:      "gemini-2.5-pro",
		FallbackModel:       "gemini-2.5-flash",
		ConfidenceThreshold: 0.8,
		Version:             "v1.0.0",
		CostPerToken:        0.0000025,
	}

	// 1. Self-Registration (Optional, can be done via CI/CD DB migrations)
	// registryClient.Register(context.Background(), agentSpec)
	log.Printf("Agent [%s] initializing...", agentSpec.Name)

	// 2. Start A2A Webhook Listener
	http.HandleFunc("/inbox", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var envelope A2AEnvelope
		if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
			http.Error(w, "Invalid A2A Envelope", http.StatusBadRequest)
			return
		}

		// Acknowledge receipt of the task immediately (Asynchronous pattern)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"status":"accepted"}`))

		// 3. Spin up the agent's internal cognitive loop asynchronously
		go func(env A2AEnvelope) {
			result, err := agent.ProcessTask(context.Background(), env.Payload.Context)
			
			// 4. Emit callback if a reply endpoint was provided
			if env.ReplyTo != "" {
				sendA2ACallback(agentSpec.Name, env.ReplyTo, result, err)
			}
		}(envelope)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Agent [%s] listening for A2A messages on port %s", agentSpec.Name, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func sendA2ACallback(agentName, replyTo, result string, err error) {
	status := "completed"
	if err != nil {
		status = "failed"
		result = err.Error()
	}

	response := A2AResponse{
		Sender: agentName,
		Status: status,
		Intent: "task_result",
	}
	response.Payload.Result = result

	body, _ := json.Marshal(response)
	resp, postErr := http.Post(replyTo, "application/json", bytes.NewBuffer(body))
	if postErr != nil {
		log.Printf("Failed to send A2A callback to %s: %v", replyTo, postErr)
		return
	}
	defer resp.Body.Close()
	log.Printf("Successfully sent A2A callback to %s (Status: %d)", replyTo, resp.StatusCode)
}
