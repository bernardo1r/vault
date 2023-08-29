package main

import (
	"api/database"
	"api/handler"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func checkerr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func readClientCerts(certsFile string) (*x509.CertPool, error) {
	file, err := os.ReadFile(certsFile)
	if err != nil {
		return nil, fmt.Errorf("reading certs files: %w", err)
	}
	var certs []string
	err = json.Unmarshal(file, &certs)
	if err != nil {
		return nil, fmt.Errorf("parsing certs files: %w", err)
	}

	certPool := x509.NewCertPool()
	for key, certFilename := range certs {
		cert, err := os.ReadFile(certFilename)
		if err != nil {
			return nil, fmt.Errorf("reading client cert file %d: %w", key, err)
		}
		ok := certPool.AppendCertsFromPEM([]byte(cert))
		if !ok {
			return nil, fmt.Errorf("parsing cert %d: %w", key, err)
		}
	}

	return certPool, nil
}

func main() {
	if len(os.Args) < 5 {
		log.Fatalln("Usage: api [SERVER_ADDR] [SERVER_CERT] [SERVER_KEY] [CLIENTS_CERT_FILE]")
	}
	serverAddr := os.Args[1]
	serverCert, err := tls.LoadX509KeyPair(os.Args[2], os.Args[3])
	if err != nil {
		err = fmt.Errorf("loading server key/certificate: %w", err)
		checkerr(err)
	}

	clientCerts, err := readClientCerts(os.Args[4])
	checkerr(err)

	db, err := database.New(os.Getenv(("DATABASE_URL")))
	checkerr(err)

	router := handler.New(db)
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handler.Log(router.Register))
	mux.HandleFunc("/login", handler.Log(router.Login))
	mux.HandleFunc("/authenticator", handler.Log(router.Authenticator))
	mux.HandleFunc("/item", handler.Log(router.Item))
	mux.HandleFunc("/index", handler.Log(router.Index))

	server := &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
		TLSConfig: &tls.Config{
			ClientAuth:   tls.RequestClientCert,
			ClientCAs:    clientCerts,
			Certificates: []tls.Certificate{serverCert},
		},
	}

	fmt.Println(serverAddr)
	err = server.ListenAndServeTLS("", "")
	checkerr(err)
}
