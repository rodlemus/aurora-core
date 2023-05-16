package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

type ServerHttp struct {
	Router      *chi.Mux
	l           *log.Logger
	Server      *http.Server
	serverPort  string
	apiVersion  string
	routerExist bool
}

func NewServerHttp(l *log.Logger, port string, v string) *ServerHttp {

	sh := &ServerHttp{serverPort: port, apiVersion: v, routerExist: false}
	sh.Router = chi.NewRouter()
	sh.Server = &http.Server{
		Addr:    sh.serverPort,
		Handler: sh.Router,
	}
	sh.l = l

	return sh
}

type RouterMap map[string]func(chi.Router)

// Server Middlewares goes previous to public and private routers
func (sh *ServerHttp) MountHandlers(rootRouter RouterMap, serverMiddlewares ...func(http.Handler) http.Handler) {

	if len(rootRouter) == 0 {
		sh.l.Panicln("No routers provided")
		os.Exit(1)
	}

	sh.routerExist = true

	for _, m := range serverMiddlewares {
		sh.Router.Use(m)
	}

	sh.Router.Route(sh.apiVersion, func(r chi.Router) {

		for k, v := range rootRouter {
			r.Route(k, v)
		}
	})

	if err := chi.Walk(sh.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		sh.l.Printf("%s %s\n", method, route)
		return nil
	}); err != nil {
		sh.l.Printf("Error al listar las rutas: %s\n", err.Error())
	}
}

func (sh *ServerHttp) Run() {
	if !sh.routerExist {
		sh.l.Println("No routers provided")
		os.Exit(1)
	}
	sh.l.Println("running on port ", sh.serverPort)
	err := sh.Server.ListenAndServe()

	if err != nil {
		sh.l.Println(err)
		os.Exit(1)
	}
}

func (sh *ServerHttp) Shutdown() {
	sh.l.Println("Stopping server")
	const timeout = 30 * time.Second

	ctxTimeout, _ := context.WithTimeout(context.Background(), timeout)
	if err := sh.Server.Shutdown(ctxTimeout); err != nil {
		sh.l.Printf("No se pudo cerrar el servidor correctamente: %v", err)
	} else {
		sh.l.Println("El servidor se cerró correctamente.")
	}

}
