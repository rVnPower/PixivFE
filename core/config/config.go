package core

import (
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/doc"
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
	s.setVersion()

	doc.CollectAllEnv()

	token, hasToken := doc.LookupEnv("PIXIVFE_TOKEN")
	if !hasToken {
		log.Fatalln("PIXIVFE_TOKEN is required, but was not set.")
		return errors.New("PIXIVFE_TOKEN is required, but was not set.\n")
	}
	// TODO Maybe add some testing?
	s.Token = strings.Split(token, ",")

	port, hasPort := doc.LookupEnv("PIXIVFE_PORT")
	socket, hasSocket := doc.LookupEnv("PIXIVFE_UNIXSOCKET")
	if !hasPort && !hasSocket {
		log.Fatalln("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
		return errors.New("Either PIXIVFE_PORT or PIXIVFE_UNIXSOCKET has to be set.")
	}
	s.Port = port
	s.UnixSocket = socket

	_, hasDev := doc.LookupEnv("PIXIVFE_DEV")
	s.InDevelopment = hasDev

	s.Host = doc.GetEnv("PIXIVFE_HOST")

	s.UserAgent = doc.GetEnv("PIXIVFE_USERAGENT")

	s.AcceptLanguage = doc.GetEnv("PIXIVFE_ACCEPTLANGUAGE")

	s.SetRequestLimit(doc.GetEnv("PIXIVFE_REQUESTLIMIT"))

	s.SetProxyServer(doc.GetEnv("PIXIVFE_IMAGEPROXY"))

	doc.AnnounceAllEnv()

	s.setStartingTime()

	return nil
}

func (s *ServerConfig) SetProxyServer(v string) {
	proxyUrl, err := url.Parse(v)
	if err != nil {
		panic(err)
	}
	s.ProxyServer = *proxyUrl
	if (proxyUrl.Scheme == "") != (proxyUrl.Host == "") {
		log.Panicf("proxy server url is weird: %s\nPlease specify e.g. https://example.com", proxyUrl.String())
	}
	if strings.HasSuffix(proxyUrl.Path, "/") {
		log.Panicf("proxy server path (%s) has cannot end in /: %s\nPixivFE does not support this now, sorry", proxyUrl.Path, proxyUrl.String())
	}
	log.Printf("Set image proxy server to: %s\n", proxyUrl.String())
}

func (s *ServerConfig) SetRequestLimit(v string) {
	t, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	s.RequestLimit = t
}

func (s *ServerConfig) setStartingTime() {
	s.StartingTime = time.Now().UTC().Format("2006-01-02 15:04")
	log.Printf("Set starting time to: %s\n", s.StartingTime)
}

func (s *ServerConfig) setVersion() {
	s.Version = "v2.3"
	log.Printf("Set server version to: %s\n", s.Version)
}
