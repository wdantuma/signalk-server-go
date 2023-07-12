package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wdantuma/signalk-server-go/canboat"
	"github.com/wdantuma/signalk-server-go/signalkserver"
)

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
	staticPath := flag.String("webapppath", "./static", "Path to webapps")
	mmsi := flag.String("mmsi", "", "Vessel MMSI")
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
	signalkServer := signalkserver.NewSignalkServer()
	if *debug {
		signalkServer.EnableDebug()
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

	signalkServer.SetupServer(ctx, "", router)

	if *serveWebApps {
		fmt.Printf("Serving webapps from %s\n", *staticPath)
		// setup static file server at /@signalk
		fs := http.FileServer(http.Dir(*staticPath))
		router.PathPrefix("/@signalk").Handler(fs)
		router.Handle("/", http.RedirectHandler("/@signalk/feeboard-sk", http.StatusSeeOther))
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
