package main

import (
	"net/http"
	"strconv"
)

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

