package httpgin

import (
	"advertisingService/internal/ports/httpgin/httpAd"
	"advertisingService/internal/ports/httpgin/httpUser"
	"advertisingService/internal/userApp"
	"net/http"

	"github.com/gin-gonic/gin"

	"advertisingService/internal/adApp"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a adApp.App, ua userApp.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}
	api := s.app.Group("/api/v1")
	httpAd.AppRouter(api, a)
	httpUser.AppRouter(api, ua)
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
