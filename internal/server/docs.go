package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type docsTranscriptionForm struct {
	File           huma.FormFile `form:"file"`
	Language       string        `form:"language"`
	Prompt         string        `form:"prompt"`
	DetectLanguage bool          `form:"detect_language"`
	EnhanceAudio   bool          `form:"enhance_audio"`
	ResponseFormat string        `form:"response_format"`
	Stream         bool          `form:"stream"`
	Model          string        `form:"model"`
	BeamSize       int           `form:"beam_size"`
	BestOf         int           `form:"best_of"`
	DiarizeModel   string        `form:"diarize_model"`
	MaxSegLen      int           `form:"max_segment_len"`
	MaxTextCtx     int           `form:"max_text_ctx"`
	NThreads       int           `form:"n_threads"`
	SamplingStrat  string        `form:"sampling_strategy"`
	StabTimestamps bool          `form:"stable_timestamps"`
	Temperature    float32       `form:"temperature"`
	Translate      bool          `form:"translate"`
	VadModel       string        `form:"vad_model"`
	WordTimestamps bool          `form:"word_timestamps"`
}

type docsTranscriptionInput struct {
	RawBody huma.MultipartFormFiles[docsTranscriptionForm]
}

type docsTranscriptionOutput struct {
	Body struct {
		Text string `json:"text"`
	}
}

type docsModelsOutput struct {
	Body map[string]any
}

type docsModelLoadInput struct {
	Body struct {
		Path string `json:"path"`
	}
}

type docsModelLoadOutput struct {
	Body struct {
		Status string `json:"status"`
		Model  string `json:"model"`
	}
}

type docsStatusOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

type docsReadyOutput struct {
	Body struct {
		Status string `json:"status"`
		Model  string `json:"model,omitempty"`
	}
}

func (s *Server) registerDocsRoutes(mux *http.ServeMux) {
	docsMux := http.NewServeMux()
	config := huma.DefaultConfig("Sona API", "dev")
	config.DocsPath = ""
	api := humago.New(docsMux, config)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/v1/audio/transcriptions",
		OperationID: "createTranscription",
		Summary:     "Create transcription",
	}, func(context.Context, *docsTranscriptionInput) (*docsTranscriptionOutput, error) {
		return nil, huma.Error501NotImplemented("spec-only operation")
	})

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/v1/models",
		OperationID: "listModels",
		Summary:     "List loaded models",
	}, func(context.Context, *struct{}) (*docsModelsOutput, error) {
		return nil, huma.Error501NotImplemented("spec-only operation")
	})

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/v1/models/load",
		OperationID: "loadModel",
		Summary:     "Load a model",
	}, func(context.Context, *docsModelLoadInput) (*docsModelLoadOutput, error) {
		return nil, huma.Error501NotImplemented("spec-only operation")
	})

	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/v1/models",
		OperationID: "unloadModel",
		Summary:     "Unload the current model",
	}, func(context.Context, *struct{}) (*docsStatusOutput, error) {
		return nil, huma.Error501NotImplemented("spec-only operation")
	})

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/health",
		OperationID: "healthCheck",
		Summary:     "Health check (always 200)",
	}, func(context.Context, *struct{}) (*docsStatusOutput, error) {
		return nil, huma.Error501NotImplemented("spec-only operation")
	})

	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/ready",
		OperationID: "readyCheck",
		Summary:     "Ready check (200 if model loaded, 503 otherwise)",
	}, func(context.Context, *struct{}) (*docsReadyOutput, error) {
		return nil, huma.Error501NotImplemented("spec-only operation")
	})

	mux.HandleFunc("GET /docs", serveSwaggerUI)
	mux.HandleFunc("GET /docs/", redirectDocs)
	mux.Handle("/openapi.json", docsMux)
	mux.Handle("/openapi.yaml", docsMux)
	mux.Handle("/openapi-3.0.json", docsMux)
	mux.Handle("/openapi-3.0.yaml", docsMux)
	mux.Handle("/schemas/", docsMux)
}

func redirectDocs(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/docs", http.StatusTemporaryRedirect)
}

func serveSwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Sona API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
      html, body { margin: 0; padding: 0; }
      #swagger-ui { max-width: 1200px; margin: 0 auto; }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
    <script>
      SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        displayRequestDuration: true
      });
    </script>
  </body>
</html>`)
}
