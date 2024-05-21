package main

import (
	"encoding/json"
	"fmt"
	"fri-flowser-playground/internal/project"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
)

var currentProject *project.Project
var logger *zerolog.Logger
var port = 8080

func main() {
	logger := initLogger()

	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/transactions", transactionsHandler)
	http.HandleFunc("/scripts", scriptsHandler)

	logger.Info().Msgf("Server is running at http://localhost:%d", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}

}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createTransactionHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

type CreateTransactionRequest struct {
	Source    string `json:"source"`
	Arguments string `json:"arguments"`
}

func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if currentProject == nil {
		http.Error(w, "Project not created", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var request CreateTransactionRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	result, err := currentProject.ExecuteTransaction([]byte(request.Source), request.Arguments)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(result)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to write response")
	}
}

func scriptsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createScriptHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

type CreateScriptRequest struct {
	Source    string `json:"source"`
	Arguments string `json:"arguments"`
}

func createScriptHandler(w http.ResponseWriter, r *http.Request) {
	if currentProject == nil {
		http.Error(w, "Project not created", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var request CreateScriptRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	result, err := currentProject.ExecuteScript([]byte(request.Source), request.Arguments)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(result)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to write response")
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createProjectHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

type CreateProjectRequest struct {
	ProjectUrl string `json:"projectUrl"`
}

func createProjectHandler(w http.ResponseWriter, r *http.Request) {
	logger = initLogger()
	currentProject = project.New(logger)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var request CreateProjectRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	err = currentProject.Open(request.ProjectUrl)

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func initLogger() *zerolog.Logger {

	level := zerolog.InfoLevel
	zerolog.MessageFieldName = "msg"

	writer := zerolog.MultiLevelWriter(
		NewTextWriter(),
	)

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(level)

	return &logger
}

func NewTextWriter() zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{Out: os.Stdout}
	writer.FormatMessage = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("%-44s", i)
	}

	return writer
}
