package response

import (
	"Go/config"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type MetricsPageData struct {
	Backends  map[string]*config.BackendMetrics
	Algorithm string
	Errors    []config.ErrorBackend
	State     string
	SmartMode bool
}

//go:embed templates/*
var templates embed.FS

func HandleStatusHTML(w http.ResponseWriter, r *http.Request) {
	errors := config.GlobalDB.ReadMetrics()
	data := MetricsPageData{
		Backends:  config.MetricsMap,
		Algorithm: config.ConfigSystem.Algorithm,
		Errors:    errors,
		State:     config.StateSystem,
		SmartMode: config.ConfigSystem.SmartMode,
	}

	if !config.ConfigSystem.ActiveLogin {
		tmpl, err := template.ParseFS(templates, "templates/login.html")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, "")
		if err != nil {
			log.Println("Error rendering HTML:", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	} else {
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
}
