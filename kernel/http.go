package kernel

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func CreateSessionStore() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(os.Getenv(config.Key)))
}

func GetHttpServer(router *mux.Router, cfg ServerConf) *http.Server {
	serverString := fmt.Sprintf("%s:%d", cfg.Name, cfg.Port)

	var httpServerConf = http.Server{}
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	if cfg.SSL {
		sslCfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		httpServerConf = http.Server{
			Addr:         serverString,
			Handler:      loggedRouter,
			TLSConfig:    sslCfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

	} else {
		httpServerConf = http.Server{
			Addr:    serverString,
			Handler: loggedRouter,
		}
	}

	return &httpServerConf
}
