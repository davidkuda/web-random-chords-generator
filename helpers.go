package main

import (
	"net/url"
	"strconv"
)

func (s Settings) toValues() url.Values {
	v := url.Values{}
	// Triads
	v.Set("maj", onOff(s.IncludeMajTriad))
	v.Set("min", onOff(s.IncludeMinTriad))
	if s.IncludeAug {
		v.Set("aug", "on")
	}
	if s.IncludeDim {
		v.Set("dim", "on")
	}
	// 7ths
	if s.IncludeMaj7 {
		v.Set("maj7", "on")
	}
	if s.IncludeDom7 {
		v.Set("dom7", "on")
	}
	if s.IncludeMin7 {
		v.Set("min7", "on")
	}
	if s.IncludeMin7b5 {
		v.Set("m7b5", "on")
	}
	// 9ths
	if s.IncludeMaj9 {
		v.Set("maj9", "on")
	}
	if s.IncludeDom9 {
		v.Set("dom9", "on")
	}
	if s.IncludeMin9 {
		v.Set("min9", "on")
	}
	// extensions
	if s.IncludeMaj7Sharp11 {
		v.Set("maj7sharp11", "on")
	}
	// altered
	if s.IncludeAlt {
		v.Set("alt", "on")
	}
	if s.Include7b9 {
		v.Set("sevenb9", "on")
	}
	if s.Include7Sharp11 {
		v.Set("sevensharp11", "on")
	}
	if s.Include7Sharp5 {
		v.Set("sevensharp5", "on")
	}
	// accidentals
	v.Set("flats", onOff(s.IncludeFlats))
	v.Set("sharps", onOff(s.IncludeSharps))

	if s.ShowSettings {
		v.Set("settings", "on")
	}
	v.Set("count", strconv.Itoa(s.Count))
	return v
}

func onOff(b bool) string {
	if b {
		return "on"
	} else {
		return "off"
	}
}

// Pretty hrefs for progressive enhancement
func hrefCurrent(s Settings) string { return "/?" + s.toValues().Encode() }
func hrefFlip(s Settings, key string) string {
	v := s.toValues()
	cur := v.Get(key) == "on"
	if cur {
		v.Set(key, "off")
	} else {
		v.Set(key, "on")
	}
	return "/?" + v.Encode()
}
func hrefCount(s Settings, n int) string {
	v := s.toValues()
	v.Set("count", strconv.Itoa(n))
	return "/?" + v.Encode()
}

// Query strings for HTMX endpoints
func queryCurrent(s Settings) string { return s.toValues().Encode() }
func queryFlip(s Settings, key string) string {
	v := s.toValues()
	cur := v.Get(key) == "on"
	if cur {
		v.Set(key, "off")
	} else {
		v.Set(key, "on")
	}
	return v.Encode()
}
func queryCount(s Settings, n int) string {
	v := s.toValues()
	v.Set("count", strconv.Itoa(n))
	return v.Encode()
}
