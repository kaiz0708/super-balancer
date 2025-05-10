package response

import (
	"Go/config"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type MetricsPageData struct {
	Backends  map[string]*config.BackendMetrics
	Algorithm string
}

//go:embed templates/*
var templates embed.FS

func HandleStatusHTML(w http.ResponseWriter, r *http.Request) {
	data := MetricsPageData{
		Backends:  config.MetricsMap,
		Algorithm: config.LoadBalancerDefault,
	}

	templatePath := filepath.Join("templates", "metrics.html")
	fmt.Println("templatePath : ", templatePath)
	tmpl, err := template.ParseFS(templates, "templates/metrics.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering HTML:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
