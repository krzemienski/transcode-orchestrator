package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/gziphandler"
	"github.com/Sirupsen/logrus"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/provider"
)

// JSONService will implement server.JSONService and
// handle all requests to the server.
type JSONService struct {
	providers map[string]provider.ProviderFactory
}

// NewJSONService will instantiate a JSONService
// with the given configuration.
func NewJSONService(cfg *config.Config) *JSONService {
	return &JSONService{
		providers: map[string]provider.ProviderFactory{
			"encoding.com": provider.EncodingComProvider,
		},
	}
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *JSONService) Prefix() string {
	return ""
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *JSONService) Middleware(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
// In this implementation, we're using it to provide application logging and to check errors
// and provide generic responses.
func (s *JSONService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {

		status, res, err := j(r)
		if err != nil {
			server.LogWithFields(r).WithFields(logrus.Fields{
				"error": err,
			}).Error("problems with serving request")
			return http.StatusServiceUnavailable, nil, &jsonErr{"sorry, this service is unavailable"}
		}

		server.LogWithFields(r).Info("success!")
		return status, res, nil
	}
}

// JSONEndpoints is a listing of all endpoints available in the JSONService.
func (s *JSONService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/jobs": map[string]server.JSONEndpoint{
			"POST": s.newTranscodeJob,
		},
		"/jobs/{jobID}": map[string]server.JSONEndpoint{
			"GET": s.getTranscodeJob,
		},
	}
}

type jsonErr struct {
	Err string `json:"error"`
}

func (e *jsonErr) Error() string {
	return e.Err
}
