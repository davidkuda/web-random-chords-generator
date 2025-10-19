package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"html/template"
	"log"
	mrand "math/rand"
	"net/http"
	"strconv"
	"time"
)

var tmplIndex *template.Template

type Settings struct {
	// Root accidentals
	IncludeFlats  bool // include ♭ roots
	IncludeSharps bool // include ♯ roots

	// Triads
	IncludeMajTriad bool // ""
	IncludeMinTriad bool // "min"
	IncludeAug      bool // "aug"
	IncludeDim      bool // "dim"

	// Sevenths
	IncludeMaj7   bool // "maj7"
	IncludeDom7   bool // "7"
	IncludeMin7   bool // "min7"
	IncludeMin7b5 bool // "min7♭5"

	// Ninths
	IncludeMaj9 bool // "maj9"
	IncludeDom9 bool // "9"
	IncludeMin9 bool // "min9"

	// Extensions
	IncludeMaj7Sharp11 bool // "maj7♯11"

	// Altered (dominant)
	IncludeAlt      bool // "alt"
	Include7b9      bool // "7♭9"
	Include7Sharp11 bool // "7♯11"
	Include7Sharp5  bool // "7♯5"

	// UI
	ShowSettings bool
	Count        int
}

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
	seedRand()
	parseTemplates()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Full page
	http.HandleFunc("/", handleIndex)

	// HTMX fragments
	http.HandleFunc("/main", handleMain)         // renders only <main>
	http.HandleFunc("/grid", handleGrid)         // renders only tiles grid
	http.HandleFunc("/settings", handleSettings) // renders settings host (+ OOB toggle)

	log.Println("Listening on :8875")
	if err := http.ListenAndServe(":8875", nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	s := parseSettings(r)
	data := PageData{Chords: renderRandomChords(s), Settings: s, Groups: buildGroups(s)}
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
	data := PageData{Chords: renderRandomChords(s), Settings: s, Groups: buildGroups(s)}
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
	data := PageData{Chords: renderRandomChords(s), Settings: s}
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

// ---------------- settings parsing ----------------
func parseSettings(r *http.Request) Settings {
	q := r.URL.Query()

	val := func(key string, def bool) bool {
		if v := q.Get(key); v != "" {
			return v == "on"
		}
		return def
	}

	count := 16
	if v := q.Get("count"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 400 {
			count = n
		}
	}

	return Settings{
		// Triads
		IncludeMajTriad: val("maj", true),
		IncludeMinTriad: val("min", true),
		IncludeAug:      val("aug", false),
		IncludeDim:      val("dim", false),

		// 7ths
		IncludeMaj7:   val("maj7", false),
		IncludeDom7:   val("dom7", false),
		IncludeMin7:   val("min7", false),
		IncludeMin7b5: val("m7b5", false),

		// 9ths
		IncludeMaj9: val("maj9", false),
		IncludeDom9: val("dom9", false),
		IncludeMin9: val("min9", false),

		// Extensions
		IncludeMaj7Sharp11: val("maj7sharp11", false),

		// Altered
		IncludeAlt:      val("alt", false),
		Include7b9:      val("sevenb9", false),
		Include7Sharp11: val("sevensharp11", false),
		Include7Sharp5:  val("sevensharp5", false),

		// Accidentals
		IncludeFlats:  val("flats", true),
		IncludeSharps: val("sharps", true),

		ShowSettings: q.Get("settings") == "on",
		Count:        count,
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

// ---------------- chord gen ----------------
func renderRandomChords(s Settings) []string {
	// Start with full chromatic set (prefer flats for E♭/A♭/B♭)
	allRoots := []string{"C", "C♯", "D", "E♭", "E", "F", "F♯", "G", "A♭", "A", "B♭", "B"}

	// Filter roots by accidental settings
	var roots []string
	for _, r := range allRoots {
		hasFlat, hasSharp := false, false
		for _, ch := range r {
			if ch == '♭' {
				hasFlat = true
			} else if ch == '♯' {
				hasSharp = true
			}
		}
		if hasFlat && !s.IncludeFlats {
			continue
		}
		if hasSharp && !s.IncludeSharps {
			continue
		}
		roots = append(roots, r)
	}
	if len(roots) == 0 { // user turned both flats & sharps off -> naturals only
		roots = []string{"C", "D", "E", "F", "G", "A", "B"}
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
		suffixes = []string{""}
	}

	chords := make([]string, s.Count)
	for i := 0; i < s.Count; i++ {
		root := roots[mrand.Intn(len(roots))]
		suf := suffixes[mrand.Intn(len(suffixes))]
		chords[i] = root + suf
	}
	return chords
}

func seedRand() {
	var b [8]byte
	if _, err := crand.Read(b[:]); err == nil {
		mrand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
		return
	}
	mrand.Seed(time.Now().UnixNano())
}
