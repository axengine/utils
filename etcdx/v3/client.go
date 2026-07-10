package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(cacertFile, certFile, keyFile string, username, password string, endpoints []string) (*clientv3.Client, error) {
	var tlsConfig *tls.Config
	if cacertFile != "" && certFile != "" && keyFile != "" {
		caCert, err := os.ReadFile(cacertFile)
		if err != nil {
			return nil, err
		}

		clientCert, err := os.ReadFile(certFile)
		if err != nil {
			return nil, err
		}

		clientKey, err := os.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}

		// Create a certificate pool and add CA certificates
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, err
		}
		// Create a TLS configuration
		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
			Certificates: []tls.Certificate{
				cert,
			},
		}
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}
