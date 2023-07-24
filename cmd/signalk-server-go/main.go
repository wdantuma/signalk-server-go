package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/signalkserver"
	"github.com/wdantuma/signalk-server-go/source/cansource"
	"github.com/wdantuma/signalk-server-go/source/filesource"
)

type arrayFlag []string

func (s *arrayFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *arrayFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {

	ctx := context.Background()
	cfg := tls.Config{}

	var listenPort int = 3000

	enableTls := flag.Bool("tls", false, "Enable tls")
	tlsCertFile := flag.String("tlscert", "", "Tls certificate file")
	tlsKeyFile := flag.String("tlskey", "", "Tls key file")
	serveWebApps := flag.Bool("webapps", true, "Serve webapps")
	version := flag.Bool("version", false, "Show version")
	port := flag.Int("port", listenPort, "Listen port")
	debug := flag.Bool("debug", false, "Enable debugging")
	staticPath := flag.String("webapp-path", "./static", "Path to webapps")
	mmsi := flag.String("mmsi", "", "Vessel MMSI")
	var fileSources arrayFlag
	flag.Var(&fileSources, "file-source", "Path to candump file")
	var sources arrayFlag
	flag.Var(&sources, "source", "Source Can device")

	flag.Parse()

	listenPort = *port

	if *tlsCertFile != "" && *tlsKeyFile != "" && *enableTls {
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			log.Fatal(err)
		}

		cfg.Certificates = append(cfg.Certificates, cert)
	}

	router := mux.NewRouter()
	router.Use((loggingMiddleware))
	router.Use(handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"authorization", "content-type", "dpop"}),
		handlers.AllowedOriginValidator(func(_ string) bool {
			return true
		}),
	))
	signalkServer := signalkserver.NewSignalkServer()
	if *debug {
		signalkServer.SetDebug(true)
		router.Use(loggingMiddleware)
	}
	if *mmsi != "" {
		signalkServer.SetMMSI(*mmsi)
	}

	if *version {
		fmt.Printf("%s version : %s\n", signalkServer.GetName(), signalkServer.GetVersion())
		fmt.Printf("canboat version : %s\n", canboat.Version)
		fmt.Printf("signalk version : 2.0.0\n")
		return
	}

	if len(fileSources) > 0 {
		for _, fs := range fileSources {
			canSource, err := filesource.CreateFileSource(fs)
			if err != nil {
				log.Fatal(err)
			} else {
				signalkServer.AddSource(canSource)
			}
		}
	}

	if len(sources) > 0 {
		for _, s := range sources {
			canSource, err := cansource.NewCanSource(s)
			if err != nil {
				log.Fatal(err)
			} else {
				signalkServer.AddSource(canSource)
			}
		}
	}

	signalkServer.SetupServer(ctx, "", router)

	if *serveWebApps {
		fmt.Printf("Serving webapps from %s\n", *staticPath)
		// setup static file server at /@signalk
		fs := http.FileServer(http.Dir(*staticPath))
		router.PathPrefix("/@signalk").Handler(fs)
		router.Handle("/", http.RedirectHandler("/@signalk/freeboard-sk", http.StatusSeeOther))
	}

	// start listening
	fmt.Printf("Listening on :%d...\n", listenPort)
	server := http.Server{Addr: fmt.Sprintf(":%d", listenPort), Handler: router, TLSConfig: &cfg}
	if *enableTls {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}

	<-ctx.Done()
}
