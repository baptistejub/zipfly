package zipfly

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ServerOptions struct {
	ValidateSignature bool
	SigningSecret     string
	PublicUrl         string
}

type Server struct {
	environment string
	options     ServerOptions
	router      *mux.Router
}

type zipPayload struct {
	Filename  string `json:"filename"`
	Files     []File `json:"files"`
	Signature string `json:"signature,omitempty"`
}

type File struct {
	Url      string `json:"url"`
	Filename string `json:"filename"`
	Compress bool   `json:"compress,omitempty"`
}

func fetch(sourceUrl string) (*zipPayload, error) {
	fmt.Println("Fetching files to zip from", sourceUrl)
	resp, err := http.Get(sourceUrl)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("couldn't fetch from URL")
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return UnmarshalPayload(bodyBytes)
}

func UnmarshalPayload(payload []byte) (*zipPayload, error) {
	var parsed zipPayload
	err := json.Unmarshal(payload, &parsed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func NewServer(env string, options ServerOptions) *Server {
	r := mux.NewRouter()

	server := Server{environment: env, options: options, router: r}

	r.HandleFunc("/zip", server.HandleGetStreamZip).Methods("GET")
	r.HandleFunc("/zip", server.HandlePostStreamZip).Methods("POST")
	r.HandleFunc("/healthz", server.HealthCheck).Methods("GET")

	return &server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	handlers.CORS(originsOk, headersOk, methodsOk)(s.router).ServeHTTP(w, r)
}

func (s *Server) HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) HandleGetStreamZip(w http.ResponseWriter, req *http.Request) {
	if !s.validateGetRequestSignature(req) {
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	payload, err := s.extractZipPayloadFromQueryString(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.streamZip(w, payload)
}

func (s *Server) HandlePostStreamZip(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !s.validatePostRequestSignature(req, body) {
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	payload, err := s.zipPayloadFromBody(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.streamZip(w, payload)
}

func (s *Server) extractZipPayloadFromQueryString(req *http.Request) (*zipPayload, error) {
	query := req.URL.Query()
	if query["source"] == nil || query["source"][0] == "" {
		return nil, errors.New("missing source url")
	}

	decodedSourceUrl, err := base64.StdEncoding.DecodeString(query["source"][0])
	if err != nil {
		return nil, errors.New("invalid source url")
	}

	payload, err := fetch(string(decodedSourceUrl))
	if err != nil {
		return nil, err
	}

	if query["filename"] != nil {
		payload.Filename = query["filename"][0]
	}

	return payload, nil
}

func (s *Server) zipPayloadFromBody(body []byte) (*zipPayload, error) {
	payload, err := UnmarshalPayload(body)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (s *Server) streamZip(w http.ResponseWriter, payload *zipPayload) {
	if payload.Filename == "" {
		payload.Filename = "archive.zip"
	}

	fmt.Println("Creating zip:", payload.Filename)

	zipStreamer, err := NewZipStreamer(payload.Files)
	if err != nil {
		fmt.Println("Error while parsing source files for", payload.Filename, ":", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// need to write the header before bytes
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", payload.Filename))
	w.WriteHeader(http.StatusOK)
	err = zipStreamer.StreamFiles(w)

	fmt.Println("Done streaming zip:", payload.Filename)

	if err != nil {
		fmt.Println("Streaming error for", payload.Filename, ":", err.Error())
		closeForError(w)
	}
}

// Close the connection so the client gets an error instead of 200 with an invalid file
func closeForError(w http.ResponseWriter) {
	hj, ok := w.(http.Hijacker)

	if !ok {
		return
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		return
	}

	conn.Close()
}
