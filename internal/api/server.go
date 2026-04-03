package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	ollamaapi "github.com/ollama/ollama/api"

	"github.com/nareshnavinash/bonsai/internal/registry"
)

type Server struct {
	Client       *ollamaapi.Client
	DefaultModel string
	Host         string
	Port         int
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/v1/chat/completions", s.handleChatCompletions)
	mux.HandleFunc("/v1/models", s.handleModels)

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	fmt.Printf("Bonsai API server listening on http://%s\n", addr)
	fmt.Printf("OpenAI-compatible endpoint: http://%s/v1/chat/completions\n", addr)
	return http.ListenAndServe(addr, cors(mux))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "invalid_request_error")
		return
	}

	models := []Model{}

	// Add installed Ollama models
	if resp, err := s.Client.List(context.Background()); err == nil {
		for _, m := range resp.Models {
			models = append(models, Model{
				ID:      m.Name,
				Object:  "model",
				OwnedBy: "local",
			})
		}
	}

	// Add bonsai registry models (if not already installed)
	installed := map[string]bool{}
	for _, m := range models {
		installed[m.ID] = true
	}
	for _, m := range registry.Models {
		if !installed[m.HFRepo] && !installed[m.Name] {
			models = append(models, Model{
				ID:      m.Name,
				Object:  "model",
				OwnedBy: "prism-ml",
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ModelList{
		Object: "list",
		Data:   models,
	})
}

func (s *Server) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "invalid_request_error")
		return
	}

	var req ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error(), "invalid_request_error")
		return
	}

	if len(req.Messages) == 0 {
		writeError(w, http.StatusBadRequest, "messages is required and must not be empty", "invalid_request_error")
		return
	}

	// Resolve model
	model := s.DefaultModel
	if req.Model != "" {
		model = registry.Resolve(req.Model)
	}

	// Build Ollama messages
	ollamaMessages := make([]ollamaapi.Message, len(req.Messages))
	for i, m := range req.Messages {
		ollamaMessages[i] = ollamaapi.Message{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	// Build options
	options := map[string]interface{}{}
	if req.Temperature != nil {
		options["temperature"] = *req.Temperature
	}
	if req.TopP != nil {
		options["top_p"] = *req.TopP
	}

	chatReq := &ollamaapi.ChatRequest{
		Model:    model,
		Messages: ollamaMessages,
		Options:  options,
	}

	requestID := "chatcmpl-" + uuid.New().String()[:8]

	if req.Stream {
		s.handleStream(w, chatReq, requestID, model)
	} else {
		s.handleNonStream(w, chatReq, requestID, model)
	}
}

func (s *Server) handleNonStream(w http.ResponseWriter, req *ollamaapi.ChatRequest, id, model string) {
	var fullResponse string

	err := s.Client.Chat(context.Background(), req, func(resp ollamaapi.ChatResponse) error {
		if resp.Message.Content != "" {
			fullResponse += resp.Message.Content
		}
		// If model only produces thinking tokens, capture those as fallback
		if resp.Message.Thinking != "" && fullResponse == "" {
			// Will be used only if no content tokens arrive
		}
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			writeError(w, http.StatusServiceUnavailable, "cannot connect to Ollama", "server_error")
		} else if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "model not found: "+model, "invalid_request_error")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error(), "server_error")
		}
		return
	}

	// Filter artifacts
	fullResponse = strings.ReplaceAll(fullResponse, "<tool_call>", "")
	fullResponse = strings.ReplaceAll(fullResponse, "</tool_call>", "")
	fullResponse = strings.TrimSpace(fullResponse)

	stop := "stop"
	resp := ChatCompletionResponse{
		ID:     id,
		Object: "chat.completion",
		Model:  model,
		Choices: []Choice{{
			Index:        0,
			Message:      &Message{Role: "assistant", Content: fullResponse},
			FinishReason: &stop,
		}},
		Usage: Usage{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleStream(w http.ResponseWriter, req *ollamaapi.ChatRequest, id, model string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported", "server_error")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	err := s.Client.Chat(context.Background(), req, func(resp ollamaapi.ChatResponse) error {
		content := resp.Message.Content
		if content == "" {
			return nil
		}

		// Filter artifacts
		content = strings.ReplaceAll(content, "<tool_call>", "")
		content = strings.ReplaceAll(content, "</tool_call>", "")
		if content == "" {
			return nil
		}

		chunk := ChatCompletionResponse{
			ID:     id,
			Object: "chat.completion.chunk",
			Model:  model,
			Choices: []Choice{{
				Index:        0,
				Delta:        &Message{Role: "assistant", Content: content},
				FinishReason: nil,
			}},
		}

		data, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return nil
	})

	if err != nil {
		errChunk := ErrorResponse{Error: ErrorDetail{Message: err.Error(), Type: "server_error"}}
		data, _ := json.Marshal(errChunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return
	}

	// Send finish chunk
	stop := "stop"
	finishChunk := ChatCompletionResponse{
		ID:     id,
		Object: "chat.completion.chunk",
		Model:  model,
		Choices: []Choice{{
			Index:        0,
			Delta:        &Message{},
			FinishReason: &stop,
		}},
	}
	data, _ := json.Marshal(finishChunk)
	fmt.Fprintf(w, "data: %s\n\n", data)
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func writeError(w http.ResponseWriter, status int, message, errType string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: ErrorDetail{Message: message, Type: errType},
	})
}
