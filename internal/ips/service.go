/*
 * Copyright (c) 2023 shenjunzheng@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ips

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/sjzar/ips/internal/parser"
	"github.com/sjzar/ips/pkg/model"
)

// Service initializes and runs the main web service for the application.
// It sets up middlewares, routes and starts the HTTP server.
func (m *Manager) Service() {
	router := gin.New()

	// Handle error from SetTrustedProxies
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Error("Failed to set trusted proxies:", err)
	}

	// Middleware
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	m.router = router

	m.InitRouter()

	// Handle error from Run
	if err := m.router.Run(m.Conf.Addr); err != nil {
		log.Error("Failed to run server:", err)
	}
}

// EFS holds embedded file system data for static assets.
//
//go:embed static
var EFS embed.FS

// InitRouter sets up routes and static file servers for the web service.
// It defines endpoints for API as well as serving static content.
func (m *Manager) InitRouter() {
	staticDir, _ := fs.Sub(EFS, "static")
	m.router.StaticFS("/static", http.FS(staticDir))
	m.router.StaticFileFS("/favicon.ico", "./favicon.ico", http.FS(staticDir))
	m.router.StaticFileFS("/", "./index.htm", http.FS(staticDir))

	// API Router
	api := m.router.Group("/api")
	{
		api.GET("/v1/ip", m.GetIP)
		api.GET("/v1/query", m.GetQuery)
	}

	m.router.NoRoute(m.NoRoute)
}

// NoRoute handles 404 Not Found errors. If the request URL starts with "/api"
// or "/static", it responds with a JSON error. Otherwise, it redirects to the root path.
func (m *Manager) NoRoute(c *gin.Context) {
	path := c.Request.URL.Path
	switch {
	case strings.HasPrefix(path, "/api"), strings.HasPrefix(path, "/static"):
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	default:
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Redirect(http.StatusFound, "/")
	}
}

// GetIP handles the GET /v1/ip endpoint. It takes an IP as a query parameter
// and returns its associated information in JSON format.
// Example:
// GET /v1/ip?ip=<ip>
// Response:
// {}
func (m *Manager) GetIP(c *gin.Context) {
	info, err := m.parseIP(c.Query("ip"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info.Output(m.Conf.UseDBFields))
}

// GetQuery handles the GET /v1/query endpoint. It takes a text string as a query
// parameter, parses it into segments, and returns the information associated with each segment.
// Example:
// GET /v1/query?text=<text>
// Response:
// {"items": [{},{}]}
func (m *Manager) GetQuery(c *gin.Context) {

	text := c.Query("text")
	if len(text) == 0 && len(c.Request.URL.RawQuery) > 0 {
		text = c.Request.URL.RawQuery
	}

	ret := &model.DataList{}

	tp := parser.NewTextParser(text).Parse()

	for _, segment := range tp.Segments {
		info, err := m.parseSegment(segment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		switch v := info.(type) {
		case *model.IPInfo:
			ret.AddItem(v.Output(m.Conf.UseDBFields))
		case *model.DomainInfo:
			ret.AddDomain(v)
		}
	}

	c.JSON(http.StatusOK, ret)
}
