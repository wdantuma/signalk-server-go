package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wdantuma/signalk-server-go/signalkserver"
)

func main() {

	ctx := context.Background()
	cfg := tls.Config{}

	// cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// cfg.Certificates = append(cfg.Certificates, cert)

	router := mux.NewRouter()
	signalkServer := signalkserver.NewSignalkServer()
	signalkServer.SetupServer(ctx, "", router)

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/@signalk").Handler(fs)
	router.Handle("/", http.RedirectHandler("/@signalk/instrumentpanel", http.StatusSeeOther))

	server := http.Server{Addr: ":3000", Handler: router, TLSConfig: &cfg}
	err := server.ListenAndServe()
	//err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Listening on :3000...")
	<-ctx.Done()
}
