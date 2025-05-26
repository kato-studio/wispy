package template

import (
	"crypto/tls"
	"log"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"strings"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/filters"
	"github.com/kato-studio/wispy/template/tags"
	"github.com/kato-studio/wispy/wispy_common/structure"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// var SiteMap = map[string]*structure.SiteStructure{}
var Logger *slog.JSONHandler

/*
=================================================================
Core External Functions
=================================================================
*/
func init() {
	// -------------
	// Setup Logger
	// -------------
	// logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer logFile.Close()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// EngineCtx is the engine context which holds site mappings and configuration.
// var Wispy = structure.WispyConfig{
// 	SITE_DIR:         "./sites",
// 	PAGE_FILE_NAME:   "page",
// 	FILE_EXT:         ".hstm",
// 	SITE_CONFIG_NAME: "config.toml",
// }

var DefaultTemplateFilters = []structure.TemplateFilter{
	filters.UpcaseFilter,
	filters.DowncaseFilter,
	filters.CapitalizeFilter,
	filters.StripFilter,
	filters.TruncateFilter,
	filters.SliceFilter,
}

var DefaultTemplateTags = []structure.TemplateTag{
	tags.IfTag,
	tags.EachTag,
	tags.CommentTag,
	tags.DefineTag,
	tags.BlockTag,
}

var DefaultEngineTags = []structure.TemplateTag{
	tags.IfTag,
	tags.EachTag,
	tags.PartialTag,
	tags.CommentTag,
	tags.DefineTag,
	tags.BlockTag,
	tags.ExtendsTag,
	tags.LayoutTag,
	tags.PassedTag,
	//
	tags.HeadTag,
	tags.FooterAssetsTag,
	tags.TitleTag,
	tags.MetaTag,
	tags.CSSTag,
	tags.JSTag,
	tags.ImportTag,
	tags.AssignTag,
}

func StartDefaultEngine() *structure.TemplateEngine {
	var engine = structure.TemplateEngine{}
	return engine.Init(DefaultEngineTags, DefaultTemplateFilters)
}

func StartHttpServer(r http.Handler, engine structure.TemplateEngine) {
	// Autocert manager
	var domainWhiteList = []string{}
	for domain := range maps.Keys(engine.SiteMap) {
		if strings.Contains(domain, "localhost:") {
			continue
		}
		domainWhiteList = append(domainWhiteList, domain)
		dotCount := strings.Count(domain, ".")
		if dotCount == 1 {
			domainWhiteList = append(domainWhiteList, "www."+domain)
		}
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domainWhiteList...),
		Email:      os.Getenv("ACME_CONTACT_EMAIL"),
		Cache:      autocert.DirCache("/var/www/.cache"),
		Client: &acme.Client{
			DirectoryURL: os.Getenv("ACME_DIRECTORY_URL"),
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
		},
	}

	ENV := os.Getenv("ENV")
	if strings.ToLower(ENV) == "PROD" {
		// ---------------
		//  PRODUCTION
		// ---------------
		// HTTP server for ACME challenges
		go func() {
			log.Println("Starting HTTP server on :80")
			err := http.ListenAndServe(":80", m.HTTPHandler(http.NotFoundHandler()))
			if err != nil {
				log.Fatalf("HTTP server error: %v", err)
			}
		}()

		// HTTPS server
		srv := &http.Server{
			Addr:    ":443",
			Handler: r,
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
				MinVersion:     tls.VersionTLS12,
			},
		}
		log.Println("Running in PROD mode")
		log.Println("Starting HTTPS server on :443")
		log.Fatal(srv.ListenAndServeTLS("", ""))
	} else {
		// ---------------
		//  Develop / Local
		// ---------------
		srv := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}
		log.Println("Running in DEV mode")
		log.Println("Starting HTTPS server on :8080")
		log.Fatal(srv.ListenAndServe())
	}
}

var Render = core.Render
