package main

import (
	"html/template"
	"log"
)

// ---------------- template parsing ----------------
func parseTemplates() {
	funcs := template.FuncMap{
		"mul": func(i, ms int) int { return i * ms },
		"splitChord": func(chord string) []string {
			if chord == "" {
				return []string{"", ""}
			}
			r := []rune(chord)
			root := string(r[0])
			if len(r) > 1 && (r[1] == '♭' || r[1] == '♯' || r[1] == '#') {
				if r[1] == '#' {
					root += "♯"
					return []string{root, string(r[2:])}
				}
				root += string(r[1])
				return []string{root, string(r[2:])}
			}
			return []string{root, string(r[1:])}
		},
		// URL + query helpers
		"hrefCurrent":  hrefCurrent,
		"hrefFlip":     hrefFlip,
		"hrefCount":    hrefCount,
		"queryCurrent": queryCurrent,
		"queryFlip":    queryFlip,
		"queryCount":   queryCount,
		// small util to build slices inline in templates
		"slice": func(ns ...int) []int { return ns },
	}
	var err error
	tmplIndex, err = template.New("index.html").Funcs(funcs).ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("parse index template: %v", err)
	}
}
