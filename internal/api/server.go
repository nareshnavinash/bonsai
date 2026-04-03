package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/nareshnavinash/bonsai/internal/llm"
	"github.com/nareshnavinash/bonsai/internal/registry"
)

type Server struct {
	Client       *llm.Client
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

	models := []llm.Model{}

	// Add models from the backend server
	if resp, err := s.Client.Models(r.Context()); err == nil {
		models = append(models, resp.Data...)
	}

	// Add bonsai registry models (if not already listed)
	installed := map[string]bool{}
	for _, m := range models {
		installed[m.ID] = true
	}
	for _, m := range registry.Models {
		if !installed[m.Name] {
			models = append(models, llm.Model{
				ID:      m.Name,
				Object:  "model",
				OwnedBy: "prism-ml",
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(llm.ModelList{
		Object: "list",
		Data:   models,
	})
}

func (s *Server) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "invalid_request_error")
		return
	}

	var req llm.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error(), "invalid_request_error")
		return
	}

	if len(req.Messages) == 0 {
		writeError(w, http.StatusBadRequest, "messages is required and must not be empty", "invalid_request_error")
		return
	}

	// Resolve model name
	if req.Model == "" {
		req.Model = s.DefaultModel
	} else {
		req.Model = registry.Resolve(req.Model)
	}

	requestID := "chatcmpl-" + uuid.New().String()[:8]

	if req.Stream {
		s.handleStream(w, r.Context(), &req, requestID)
	} else {
		s.handleNonStream(w, r.Context(), &req, requestID)
	}
}

func (s *Server) handleNonStream(w http.ResponseWriter, ctx context.Context, req *llm.ChatCompletionRequest, id string) {
	req.Stream = false
	resp, err := s.Client.ChatCompletion(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "not reachable") {
			writeError(w, http.StatusServiceUnavailable, "cannot connect to inference server", "server_error")
		} else if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "model not found: "+req.Model, "invalid_request_error")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error(), "server_error")
		}
		return
	}

	// Filter artifacts from response
	if len(resp.Choices) > 0 && resp.Choices[0].Message != nil {
		content := resp.Choices[0].Message.Content
		content = strings.ReplaceAll(content, "<tool_call>", "")
		content = strings.ReplaceAll(content, "</tool_call>", "")
		resp.Choices[0].Message.Content = strings.TrimSpace(content)
	}

	resp.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleStream(w http.ResponseWriter, ctx context.Context, req *llm.ChatCompletionRequest, id string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported", "server_error")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	req.Stream = true
	err := s.Client.ChatCompletionStream(ctx, req, func(resp llm.ChatCompletionResponse) error {
		if len(resp.Choices) == 0 {
			return nil
		}
		delta := resp.Choices[0].Delta
		if delta == nil || delta.Content == "" {
			return nil
		}

		// Filter artifacts
		content := delta.Content
		content = strings.ReplaceAll(content, "<tool_call>", "")
		content = strings.ReplaceAll(content, "</tool_call>", "")
		if content == "" {
			return nil
		}

		chunk := llm.ChatCompletionResponse{
			ID:     id,
			Object: "chat.completion.chunk",
			Model:  req.Model,
			Choices: []llm.Choice{{
				Index:        0,
				Delta:        &llm.Message{Role: "assistant", Content: content},
				FinishReason: nil,
			}},
		}

		data, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return nil
	})

	if err != nil {
		errChunk := llm.ErrorResponse{Error: llm.ErrorDetail{Message: err.Error(), Type: "server_error"}}
		data, _ := json.Marshal(errChunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return
	}

	// Send finish chunk
	stop := "stop"
	finishChunk := llm.ChatCompletionResponse{
		ID:     id,
		Object: "chat.completion.chunk",
		Model:  req.Model,
		Choices: []llm.Choice{{
			Index:        0,
			Delta:        &llm.Message{},
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
	json.NewEncoder(w).Encode(llm.ErrorResponse{
		Error: llm.ErrorDetail{Message: message, Type: errType},
	})
}
