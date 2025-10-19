package main

import (
	"html/template"
	"log"
	mrand "math/rand/v2"
	"net/http"
	"fmt"
)

var tmplIndex *template.Template

// Data-driven settings UI
type Option struct {
	Key   string
	Title string
	On    bool
}

type Group struct {
	Label   string
	Options []Option
}

type PageData struct {
	Chords   []string
	Settings Settings
	Groups   []Group
}

func main() {
	parseTemplates()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Full page
	http.HandleFunc("/", handleIndex)

	// HTMX fragments
	http.HandleFunc("/main", handleMain)
	http.HandleFunc("/grid", handleGrid)
	http.HandleFunc("/settings", handleSettings)

	log.Println("Listening on :8875")
	if err := http.ListenAndServe(":8875", nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	s := parseSettings(r)
	data := PageData{
		Chords: getRandomChords(s),
		Settings: s,
		Groups: buildGroups(s),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplIndex.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func isHX(r *http.Request) bool { return r.Header.Get("HX-Request") == "true" }

func fragmentFallback(w http.ResponseWriter, r *http.Request) bool {
	if !isHX(r) {
		http.Redirect(w, r, "/?"+r.URL.RawQuery, http.StatusFound)
		return true
	}
	return false
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	if fragmentFallback(w, r) {
		return
	}
	s := parseSettings(r)
	data := PageData{Chords: getRandomChords(s), Settings: s, Groups: buildGroups(s)}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplIndex.ExecuteTemplate(w, "main", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func handleGrid(w http.ResponseWriter, r *http.Request) {
	if fragmentFallback(w, r) {
		return
	}
	s := parseSettings(r)
	data := PageData{Chords: getRandomChords(s), Settings: s}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplIndex.ExecuteTemplate(w, "grid", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	if fragmentFallback(w, r) {
		return
	}
	s := parseSettings(r)
	data := PageData{Settings: s, Groups: buildGroups(s)}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplIndex.ExecuteTemplate(w, "settings_host", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// ---------------- UI model ----------------
func buildGroups(s Settings) []Group {
	return []Group{
		{Label: "Roots", Options: []Option{
			{Key: "flats", Title: "Flats", On: s.IncludeFlats},
			{Key: "sharps", Title: "Sharps", On: s.IncludeSharps},
		}},
		{Label: "Triads", Options: []Option{
			{Key: "maj", Title: "maj", On: s.IncludeMajTriad},
			{Key: "min", Title: "min", On: s.IncludeMinTriad},
			{Key: "aug", Title: "aug", On: s.IncludeAug},
			{Key: "dim", Title: "dim", On: s.IncludeDim},
		}},
		{Label: "Sevenths", Options: []Option{
			{Key: "maj7", Title: "maj7", On: s.IncludeMaj7},
			{Key: "dom7", Title: "7", On: s.IncludeDom7},
			{Key: "min7", Title: "min7", On: s.IncludeMin7},
			{Key: "m7b5", Title: "m7♭5", On: s.IncludeMin7b5},
		}},
		{Label: "Ninths", Options: []Option{
			{Key: "maj9", Title: "maj9", On: s.IncludeMaj9},
			{Key: "dom9", Title: "9", On: s.IncludeDom9},
			{Key: "min9", Title: "min9", On: s.IncludeMin9},
		}},
		{Label: "Extensions", Options: []Option{
			{Key: "maj7sharp11", Title: "maj7♯11", On: s.IncludeMaj7Sharp11},
		}},
		{Label: "Altered", Options: []Option{
			{Key: "alt", Title: "alt", On: s.IncludeAlt},
			{Key: "sevenb9", Title: "7♭9", On: s.Include7b9},
			{Key: "sevensharp11", Title: "7♯11", On: s.Include7Sharp11},
			{Key: "sevensharp5", Title: "7♯5", On: s.Include7Sharp5},
		}},
	}
}

// builds
// var roots []string like Ab, C#, G, etc
// var suffixes []string like min, Maj7, Alt, etc
// and uses both to create a var chords []string to render the chords.
func getRandomChords(s Settings) []string {
	allRoots := []string{
		"A♭",
		"A",
		"A♯",
		"B♭",
		"B",
		"C",
		"C♯",
		"D♭",
		"D",
		"D♯",
		"E♭",
		"E",
		"F",
		"F♯",
		"G♭",
		"G",
		"G♯",
	}

	// Filter roots by accidental settings
	var roots []string
	for _, root := range allRoots {
		hasFlat, hasSharp := false, false
		for _, ch := range root {
			if ch == '♭' {
				hasFlat = true
				break
			} else if ch == '♯' {
				hasSharp = true
				break
			}
		}
		if hasFlat && !s.IncludeFlats {
			continue
		}
		if hasSharp && !s.IncludeSharps {
			continue
		}
		roots = append(roots, root)
	}
	// Build suffix set based on toggles
	var suffixes []string
	if s.IncludeMajTriad {
		suffixes = append(suffixes, "")
	}
	if s.IncludeMinTriad {
		suffixes = append(suffixes, "min")
	}
	if s.IncludeAug {
		suffixes = append(suffixes, "aug")
	}
	if s.IncludeDim {
		suffixes = append(suffixes, "dim")
	}
	if s.IncludeMaj7 {
		suffixes = append(suffixes, "maj7")
	}
	if s.IncludeDom7 {
		suffixes = append(suffixes, "7")
	}
	if s.IncludeMin7 {
		suffixes = append(suffixes, "min7")
	}
	if s.IncludeMin7b5 {
		suffixes = append(suffixes, "min7♭5")
	}
	if s.IncludeMaj9 {
		suffixes = append(suffixes, "maj9")
	}
	if s.IncludeDom9 {
		suffixes = append(suffixes, "9")
	}
	if s.IncludeMin9 {
		suffixes = append(suffixes, "min9")
	}
	if s.IncludeMaj7Sharp11 {
		suffixes = append(suffixes, "maj7♯11")
	}
	if s.IncludeAlt {
		suffixes = append(suffixes, "alt")
	}
	if s.Include7b9 {
		suffixes = append(suffixes, "7♭9")
	}
	if s.Include7Sharp11 {
		suffixes = append(suffixes, "7♯11")
	}
	if s.Include7Sharp5 {
		suffixes = append(suffixes, "7♯5")
	}

	if len(suffixes) == 0 {
		suffixes = append(suffixes, "")
	}

	chords := make([]string, s.Count)
	fmt.Println("roots:", roots)
	fmt.Println("suffixes:", suffixes)
	pFreq := make(map[int]int)
	for i := 0; i < s.Count; i++ {
		p := mrand.IntN(len(roots))
		pFreq[p]++
		root := roots[p]
		k := mrand.IntN(len(suffixes))
		suf := suffixes[k]
		chords[i] = root + suf
	}
	fmt.Println(pFreq)
	for i, root := range roots {
		fmt.Printf("%2d\t%s\t%2d\n", i, root, pFreq[i])
	}
	return chords
}
