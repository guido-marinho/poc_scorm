package scormrt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RuntimeRequest represents a request to the runtime API.
type RuntimeRequest struct {
	Session string `json:"session"`
	Method  string `json:"method"`
	Element string `json:"element,omitempty"`
	Value   string `json:"value,omitempty"`
}

// RuntimeHandler dispatches runtime API calls. The external LMS can POST a JSON
// payload describing the method to invoke.
func RuntimeHandler(c *gin.Context) {
	var req RuntimeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var result interface{}
	switch req.Method {
	case "Initialize":
		result = Initialize(req.Session)
	case "Terminate":
		result = Terminate(req.Session)
	case "GetValue":
		result = GetValue(req.Session, req.Element)
	case "SetValue":
		result = SetValue(req.Session, req.Element, req.Value)
	case "Commit":
		result = Commit(req.Session)
	case "GetLastError":
		result = GetLastError(req.Session)
	case "GetErrorString":
		result = GetErrorString(req.Value)
	case "GetDiagnostic":
		result = GetDiagnostic(req.Value)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
