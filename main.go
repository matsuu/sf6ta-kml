package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/twpayne/go-kml/v3"
)

const (
	JsonUrl = "https://www.taito.co.jp/api/NxLShopList/?category=sf6_ac&pref="
	KmlUrl  = "https://matsuu.github.io/sf6ta-kml/arcades.kml"
)

type Arcade struct {
	Addr string  `json:"ADDR"`
	Cnt  string  `json:"CNT"`
	Lat  float64 `json:"LAT"`
	Lng  float64 `json:"LNG"`
	Pref string  `json:"PREF"`
	Name string  `json:"TNAME"`
}

func WriteNetworkKML(w io.Writer) error {
	k := kml.KML(
		kml.NetworkLink(
			kml.Name("STREET FIGHTER 6 TYPE ARCADE 稼働店舗"),
			kml.Link(
				kml.Href(KmlUrl),
			),
		),
	)
	if err := k.WriteIndent(w, "", "  "); err != nil {
		return fmt.Errorf("Failed to output: %w", err)
	}

	return nil
}

func WriteArcadeKML(w io.Writer) error {
	res, err := http.Get(JsonUrl)
	if err != nil {
		return fmt.Errorf("Failed to get %s: %w", JsonUrl, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Unknown status code for %s: %d", JsonUrl, res.StatusCode)
	}

	var as []Arcade
	if err := json.NewDecoder(res.Body).Decode(&as); err != nil {
		return fmt.Errorf("Failed to decode from %s: %w", JsonUrl, err)
	}

	if len(as) == 0 {
		return fmt.Errorf("No arcades found from %s", JsonUrl)
	}

	d := kml.Document()
	for _, g := range as {
		addr := fmt.Sprintf("%s%s\n設置台数 %s", g.Pref, g.Addr, g.Cnt)
		p := kml.Placemark(
			kml.Name(g.Name),
			// kml.Address(addr),
			kml.Description(addr),
			kml.Point(
				kml.Coordinates(kml.Coordinate{Lon: g.Lng, Lat: g.Lat}),
			),
		)
		d.Append(p)
	}

	k := kml.KML(d)
	if err := k.WriteIndent(w, "", "  "); err != nil {
		return fmt.Errorf("Failed to output: %w", err)
	}

	return nil
}

func main() {
	n, err := os.Create("public/sf6ta.kml")
	if err != nil {
		log.Fatal(err)
	}
	defer n.Close()

	if err := WriteNetworkKML(n); err != nil {
		log.Fatal(err)
	}

	a, err := os.Create("public/arcades.kml")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	if err := WriteArcadeKML(a); err != nil {
		log.Fatal(err)
	}
}
