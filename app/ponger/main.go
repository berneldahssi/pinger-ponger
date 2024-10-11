package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var usage = `Usage: ponger [options...]

Options:
	-c config file (default: config.yaml)

`
var configFile = flag.String("c", "config.yaml", "Configuration file")

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %s", err))
	}

	go prom(viper.GetStringMapString("metrics"))

	stopChan := make(chan os.Signal, 10)
	signal.Notify(stopChan, os.Interrupt)
	go serve(viper.GetStringMapString("service"))

	<-stopChan
	log.Println("Shutting down server...")
}

func serve(serviceConfig map[string]string) error {
	port := fmt.Sprintf(":%v", serviceConfig["port"])

	log.Printf("Server started on %v", port)
	http.HandleFunc("/ping", ping)

	if cert, key, err := checkTLS(serviceConfig); err == nil {
		log.Println("Mode: HTTPS")
		if err := http.ListenAndServeTLS(port, cert, key, nil); err != nil {
			panic(err)
		}
	} else {
		log.Println("Mode: HTTP")
		if err := http.ListenAndServe(port, nil); err != nil {
			panic(err)
		}
	}
        return fmt.Errorf("ponger has stopped")
}

func checkTLS(tlsConfig map[string]string) (cert string, key string, err error) {
	cert, certPresent := tlsConfig["tlscertificate"]
	key, keyPresent := tlsConfig["tlsprivatekey"]
	if certPresent && keyPresent {
		if !fileExists(cert) || !fileExists(key) {
			return "", "", fmt.Errorf("certificate or private key files not present")
		}
		log.Println("certificate and private key present, enabling tls...")

		return cert, key, nil
	}

	return "", "", fmt.Errorf("configuration for tls not present")
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func ping(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %v %v", r.Method, r.URL.Path)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Printf("error with write %v", err.Error())
	}
}

func prom(metricsConfig map[string]string) {
	endpoint := metricsConfig["endpoint"]
	port := fmt.Sprintf(":%v", metricsConfig["port"])
	http.Handle(endpoint, promhttp.Handler())
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("error with listen and serve %v", err.Error())
	}
}
