package main

import (
	"encoding/json"
	"fmt"
	"fri-flowser-playground/internal/project"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
)

var currentProject *project.Project
var logger *zerolog.Logger
var logCache *CacheLogWriter
var port = 8080

func main() {
	logger, logCache = initLogger()

	mux := http.NewServeMux()
	mux.HandleFunc("/projects", projectsHandler)
	mux.HandleFunc("/projects/files", projectFilesHandler)
	mux.HandleFunc("/projects/logs", projectLogsHandler)
	mux.HandleFunc("/projects/blockchain-state", blockchainStateHandler)
	mux.HandleFunc("/projects/transactions", transactionsHandler)
	mux.HandleFunc("/projects/scripts", scriptsHandler)

	corsHandler := cors.Default().Handler(mux)
	logger.Info().Msgf("Server is running at http://localhost:%d", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), corsHandler)

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
	Location  string `json:"location"`
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

	result, err := currentProject.ExecuteTransaction([]byte(request.Source), request.Location, request.Arguments)

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
	Location  string `json:"location"`
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

	result, err := currentProject.ExecuteScript([]byte(request.Source), request.Location, request.Arguments)

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

func blockchainStateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getBlockchainState(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func getBlockchainState(w http.ResponseWriter, r *http.Request) {
	if currentProject == nil {
		http.Error(w, "Project not created", http.StatusBadRequest)
		return
	}

	jsonState, err := currentProject.BlockchainState()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(jsonState)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to write response")
	}
}

func projectLogsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listProjectLogs(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func listProjectLogs(w http.ResponseWriter, r *http.Request) {
	if currentProject == nil {
		http.Error(w, "Project not created", http.StatusBadRequest)
		return
	}

	jsonLogs, err := logCache.LogsJson()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(jsonLogs)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to write response")
	}
}

func projectFilesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listProjectFilesHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func listProjectFilesHandler(w http.ResponseWriter, r *http.Request) {
	if currentProject == nil {
		http.Error(w, "Project not created", http.StatusBadRequest)
		return
	}

	files, err := currentProject.Files()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonFiles, err := json.Marshal(files)

	if err != nil {
		logger.Fatal().Err(err)
	}

	_, err = w.Write(jsonFiles)

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

func initLogger() (*zerolog.Logger, *CacheLogWriter) {

	level := zerolog.InfoLevel
	zerolog.MessageFieldName = "msg"

	cacheWriter := NewCacheLogWriter()

	writer := zerolog.MultiLevelWriter(
		NewTextWriter(),
		cacheWriter,
	)

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(level)

	return &logger, cacheWriter
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

type CacheLogWriter struct {
	logs []string
}

func NewCacheLogWriter() *CacheLogWriter {
	return &CacheLogWriter{
		logs: make([]string, 0),
	}
}

var _ io.Writer = &CacheLogWriter{}

func (c *CacheLogWriter) Write(p []byte) (n int, err error) {
	c.logs = append(c.logs, string(p))
	return len(p), nil
}

func (c *CacheLogWriter) LogsJson() ([]byte, error) {
	res, err := json.Marshal(&c.logs)

	if err != nil {
		return nil, err
	}

	return res, nil
}
