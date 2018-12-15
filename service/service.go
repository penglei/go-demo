package service

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// NewServer initializes the service with the given Database, and sets up appropriate routes.
func NewServer(db *Database) *Server {
	router := echo.New()
	server := &Server{
		router: router,
		db:     db,
	}

	server.setupRoutes()
	return server
}

// Server contains all that is needed to respond to incoming requests, like a database. Other services like a mail,
// redis, or S3 server could also be added.
type Server struct {
	router *echo.Echo
	db     *Database
}

// The ServerError type allows errors to provide an appropriate HTTP status code and message. The Server checks for
// this interface when recovering from a panic inside a handler.
type ServerError interface {
	HttpStatusCode() int
	HttpStatusMessage() string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) setupRoutes() {
	type Data struct {
		Message string
	}
	t, err := template.New("index_template").Parse(`
<html>
<head>
<title>go-demo home page</title>
</head>
<body style="background-color:red">{{.Message}}</body>
</html>
`)

	if err != nil {
		panic(err)
	}

	s.router.POST("/contacts", s.AddContact)
	s.router.GET("/contacts/:email", s.GetContactByEmail)
	s.router.GET("/", func(c echo.Context) error {
		err := t.Execute(c.Response().Writer, Data{Message: "hello, this is index of go-demo"})
		if err == nil {
			c.Response().Committed = true
		}
		return err
	})

	// By default the router will handle errors. But the service should always return JSON if possible, so these
	// custom handlers are added.

	/*
		s.router.NotFound = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				writeJSONError(w, http.StatusNotFound, "")
			},
		)

		s.router.HandleMethodNotAllowed = true
		s.router.MethodNotAllowed = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				writeJSONError(w, http.StatusMethodNotAllowed, "")
			},
		)

		s.router.PanicHandler = func(w http.ResponseWriter, r *http.Request, e interface{}) {
			serverError, ok := e.(ServerError)
			if ok {
				writeJSONError(w, serverError.HttpStatusCode(), serverError.HttpStatusMessage())
			} else {
				log.Printf("Panic during request: %v", e)
				writeJSONError(w, http.StatusInternalServerError, "")
			}
		}
	*/
}

// AddContact handles HTTP requests to add a Contact.
func (s *Server) AddContact(c echo.Context) error {
	var contact Contact
	//w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	var r, w = c.Request(), c.Response().Writer

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&contact); err != nil {
		// writeJSONError(w, http.StatusBadRequest, "Error decoding JSON")
		return errors.New("error decoding JSON")
	}

	contactID, err := s.db.AddContact(contact)
	if err != nil {
		return err
	}
	contact.ID = contactID

	writeJSON(
		w,
		http.StatusCreated,
		&ContactResponse{
			Contact: &contact,
		},
	)
	return nil
}

//GetContactByEmail handles HTTP requests to GET a Contact by an email address.
func (s *Server) GetContactByEmail(c echo.Context) error {
	// (w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	var email = strings.TrimSpace(c.Param("email"))
	if email == "" {
		return errors.New("expected a single email")
	}

	var w = c.Response().Writer
	contact, err := s.db.GetContactByEmail(email)
	if err != nil {
		writeUnexpectedError(w, err)
	} else if contact == nil {
		writeJSONNotFound(w)
	} else {
		writeJSON(
			w,
			http.StatusOK,
			&ContactResponse{
				Contact: contact,
			},
		)
	}
	return nil
}

// ===== JSON HELPERS ==================================================================================================

func writeJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	var err = encoder.Encode(response)
	if err != nil {
		panic(err)
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	if message == "" {
		message = http.StatusText(statusCode)
	}

	writeJSON(
		w,
		statusCode,
		&ErrorResponse{
			StatusCode: statusCode,
			Message:    message,
		},
	)
}

func writeJSONNotFound(w http.ResponseWriter) {
	writeJSONError(w, http.StatusNotFound, "")
}

func writeUnexpectedError(w http.ResponseWriter, err error) {
	writeJSONError(w, http.StatusInternalServerError, err.Error())
}
