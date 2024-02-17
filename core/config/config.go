package core

import (
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var GlobalServerConfig ServerConfig

type ServerConfig struct {
	// Required
	Token []string

	ProxyServer url.URL // proxy server, may contain prefix as well

	// can be left empty
	Host string

	// One of two is required
	Port       string
	UnixSocket string

	UserAgent      string
	AcceptLanguage string
	RequestLimit   int

	StartingTime  string
	Version       string
	InDevelopment bool
}

func (s *ServerConfig) InitializeConfig() error {
	_, hasDev := os.LookupEnv("PIXIVFE_DEV")
	s.InDevelopment = hasDev
	if s.InDevelopment {
		log.Printf("Set server to development mode\n")
	}

	token, hasToken := os.LookupEnv("PIXIVFE_TOKEN")
	if !hasToken {
		log.Fatalln("PIXIVFE_TOKEN is required, but was not set.")
		return errors.New("PIXIVFE_TOKEN is required, but was not set.\n")
	}
	s.SetToken(token)

	proxyServer, hasProxyServer := os.LookupEnv("PIXIVFE_IMAGEPROXY")
	if hasProxyServer {
		s.SetProxyServer(proxyServer)
	} else {
		s.ProxyServer = url.URL{Path: "/proxy/i.pximg.net"}
	}

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

	requestLimit, hasRequestLimit := os.LookupEnv("PIXIVFE_REQUESTLIMIT")
	if hasRequestLimit {
		s.SetRequestLimit(requestLimit)
	} else {
		s.RequestLimit = 15
	}

	s.setStartingTime()
	s.setVersion()

	return nil
}

func (s *ServerConfig) SetToken(v string) {
	// TODO Maybe add some testing?
	s.Token = strings.Split(v, ",")
	log.Printf("Set token to: %s\n", v)
}

func (s *ServerConfig) SetProxyServer(v string) {
	proxyUrl, err := url.Parse(v)
	if err != nil {
		panic(err)
	}
	s.ProxyServer = *proxyUrl
	if proxyUrl.Scheme == "" {
		log.Panicf("proxy server url has no scheme: %s\nPlease specify e.g. https://example.com", proxyUrl.String())
	}
	if proxyUrl.Host == "" {
		log.Panicf("proxy server url has no host: %s\nPlease specify e.g. https://example.com", proxyUrl.String())
	}
	if strings.HasSuffix(proxyUrl.Path, "/") {
		log.Panicf("proxy server path (%s) has cannot end in /: %s\nPixivFE does not support this now, sorry", proxyUrl.Path, proxyUrl.String())
	}
	log.Printf("Set image proxy server to: %s\n", proxyUrl.String())
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

func (s *ServerConfig) SetRequestLimit(v string) {
	t, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	s.RequestLimit = t
	log.Printf("Set request limit to %s requests per 30 seconds\n", v)
}

func (s *ServerConfig) setStartingTime() {
	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")
	log.Printf("Set starting time to: %s\n", s.StartingTime)
}

func (s *ServerConfig) setVersion() {
	s.Version = "v2.3"
	log.Printf("Set server version to: %s\n", s.Version)
}
