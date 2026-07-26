# Agent Template

This is a template repository for creating new specialized agents in the LatentWorx-Asimov organization.

## Repository Layout
- `cmd/worker/main.go`: The main entrypoint for the lightweight A2A Agent server. Handles self-registration (if needed) and listens for A2A webhooks on `/inbox`.
- `internal/agent/`: Contains the agent's core cognitive loop (`ProcessTask`), prompts, and tool implementations.
- `Dockerfile`: Multi-stage Docker build, outputting a scratch/distroless image to run the agent.
- `deploy/manifests/`: Agnostic Kubernetes manifests (Deployment, ServiceAccount, RBAC if necessary) to deploy this agent into any standard Kubernetes environment (like `k3s`).
- `Makefile`: Standard local dev targets (`make build`, `make run`, `make test`).

## Setup Instructions
1. Use this repository as a template (`gh repo create LatentWorx-Asimov/my-new-agent --template LatentWorx-Asimov/agent-template`).
2. Clone your new repository.
3. Update `go.mod` with your new module name (`go mod edit -module github.com/LatentWorx-Asimov/my-new-agent`).
4. Update `cmd/worker/main.go` to define the new agent's specific metadata (Name, Description, AllowedTools) in the `AgentSpec`.
5. Implement your agent's LLM logic, constraints, and tool execution in `internal/agent/agent.go`.
6. Update the Deployment manifest in `deploy/manifests/deployment.yaml` to reference your agent's image name.

## Deployment Agnosticism
The actual orchestration/deployment of these agents should ideally be managed from a central GitOps repository (e.g. `lwx-infra`). This repository merely provides a standard container image and baseline Kubernetes objects that are agnostic to the specific environment.
