package core

import (
	"errors"
	"log"
	"os"
	"time"
)

var GlobalServerConfig ServerConfig

type ServerConfig struct {
	// Required
	Token       string
	ProxyServer string

	// can be left empty
	Host string

	// One of two is required
	Port       string
	UnixSocket string

	BaseURL        string
	UserAgent      string
	AcceptLanguage string

	StartingTime string
	Version      string
}

func (s *ServerConfig) InitializeConfig() error {
	token, hasToken := os.LookupEnv("PIXIVFE_TOKEN")
	if !hasToken {
		log.Fatalln("PIXIVFE_TOKEN is required, but was not set.")
		return errors.New("PIXIVFE_TOKEN is required, but was not set.\n")
	}
	s.SetToken(token)

	proxyServer, hasProxyServer := os.LookupEnv("PIXIVFE_IMAGEPROXY")
	if !hasProxyServer {
		log.Fatalln("PIXIVFE_IMAGEPROXY is required, but was not set.")
		return errors.New("PIXIVFE_IMAGEPROXY is required, but was not set.\n")
	}
	s.SetProxyServer(proxyServer)

	hostname, hasHostname := os.LookupEnv("PIXIVFE_HOST")
	if hasHostname {
		log.Printf("Set TCP hostname to: %s\n", hostname)
		s.Host = hostname
	}

	port, hasPort := os.LookupEnv("PIXIVFE_PORT")
	if hasPort {
		s.SetPort(port)
	}

	socket, hasSocket := os.LookupEnv("PIXIVFE_UNIXSOCKET")
	if hasSocket {
		s.SetUnixSocket(socket)
	}

	if !hasPort && !hasSocket {
		log.Fatalln("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
		return errors.New("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
	}

	// Not required
	s.SetBaseURL(os.Getenv("PIXIVFE_BASEURL"))

	userAgent, hasUserAgent := os.LookupEnv("PIXIVFE_USERAGENT")
	if !hasUserAgent {
		userAgent = "Mozilla/5.0"
	}
	s.SetUserAgent(userAgent)

	acceptLanguage, hasAcceptLanguage := os.LookupEnv("PIXIVFE_ACCEPTLANGUAGE")
	if !hasAcceptLanguage {
		acceptLanguage = "en-US,en;q=0.5"
	}
	s.SetAcceptLanguage(acceptLanguage)

	s.setStartingTime()
	s.setVersion()

	return nil
}

func (s *ServerConfig) SetToken(v string) {
	s.Token = v
	log.Printf("Set token to: %s\n", v)
}

func (s *ServerConfig) SetBaseURL(v string) {
	s.BaseURL = v
	log.Printf("Set base URL to: %s\n", v)
}

func (s *ServerConfig) SetProxyServer(v string) {
	s.ProxyServer = v
	log.Printf("Set image proxy server to: %s\n", v)
}

func (s *ServerConfig) SetPort(v string) {
	s.Port = v
	log.Printf("Set TCP port to: %s\n", v)
}

func (s *ServerConfig) SetUnixSocket(v string) {
	s.UnixSocket = v
	log.Printf("Set UNIX socket path to: %s\n", v)
}

func (s *ServerConfig) SetUserAgent(v string) {
	s.UserAgent = v
	log.Printf("Set user agent to: %s\n", v)
}

func (s *ServerConfig) SetAcceptLanguage(v string) {
	s.AcceptLanguage = v
	log.Printf("Set Accept-Language header to: %s\n", v)
}

func (s *ServerConfig) setStartingTime() {
	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")
	log.Printf("Set starting time to: %s\n", s.StartingTime)
}

func (s *ServerConfig) setVersion() {
	s.Version = "v2.0.1"
	log.Printf("Set server version to: %s\n", s.Version)
}
