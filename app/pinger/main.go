package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	netUrl "net/url"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
)

var usage = `Usage: pinger [options...]

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

	stopChan := make(chan os.Signal, 10)
	signal.Notify(stopChan, os.Interrupt)

	go ping(viper.GetStringMapString("ponger"))
	<-stopChan
	log.Println("Shutting down pinger...")
}

func ping(pongerConfig map[string]string) {
	url := fmt.Sprintf("%v", pongerConfig["url"])
	log.Printf("Starting to ping  %v", url)

	u, err := netUrl.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	tr := new(http.Transport)
	if certFile, ok := pongerConfig["acceptcert"]; ok && u.Scheme == "https" {
		certPool, err := createCertPool(certFile)
		if err != nil {
			log.Fatalf("Failed to generate certificate pool %v", err.Error())
		}
		config := &tls.Config{
			RootCAs:    certPool,
			MinVersion: tls.VersionTLS12,
		}
		tr = &http.Transport{TLSClientConfig: config}
	}

	client := &http.Client{Transport: tr}

	for {
		r := rand.Intn(1000)
		sleepPeriod := time.Duration(r) * time.Millisecond
		time.Sleep(sleepPeriod)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/ping", url), nil)
		if err != nil {
			log.Printf("error creating http request %v", err.Error())
		}
		resp, err := client.Do(req)
		log.Printf("Sent ping")
		if err != nil {
			log.Printf("Request failed %v", err.Error())

			continue
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			log.Printf("error with buffer %v", err.Error())

			continue
		}
		newStr := buf.String()
		log.Printf("Got %v", newStr)
	}
}

func createCertPool(certFile string) (*x509.CertPool, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("failed to append %q to rootcas: %v", certFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	return rootCAs, nil
}
