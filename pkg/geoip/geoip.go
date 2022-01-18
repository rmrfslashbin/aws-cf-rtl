package geoip

import (
	"log"
	"net"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

const (
	geodb = "" // path to 'GeoLite2-City.mmdb'
)

var (
	db *maxminddb.Reader
)

// Record defines the fields to fetch from the GeoIP database.
type Record struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`

	Continent struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"continent"`

	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`

	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`

	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`

	Subdivisions []struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"subdivisions"`
}

// GeoIPData represents the data returned.
type GeoIPData struct {
	IP          net.IP  `json:"ip"`
	City        string  `json:"city_name"`
	Continent   string  `json:"continent_code"`
	Country     string  `json:"country_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   uint    `json:"metro_code"`
	TimeZone    string  `json:"time_zone"`
	PostalCode  string  `json:"postal_code"`
	Subdivision string  `json:"subdivision_code"`
}

// init sets up the database connection
func init() {
	db, err := maxminddb.Open(geodb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// Lookip GeoIP data for the given IP address.
func Lookup(ip net.IP) (*GeoIPData, error) {
	var data Record = Record{}

	err := db.Lookup(ip, &data)
	if err != nil {
		return nil, err
	}

	subdivs := []string{}
	for _, s := range data.Subdivisions {
		subdivs = append(subdivs, s.IsoCode)
	}

	return &GeoIPData{
		City:        data.City.Names["en"],
		Continent:   data.Continent.Code,
		Country:     data.Country.IsoCode,
		Latitude:    data.Location.Latitude,
		Longitude:   data.Location.Longitude,
		MetroCode:   data.Location.MetroCode,
		TimeZone:    data.Location.TimeZone,
		PostalCode:  data.Postal.Code,
		Subdivision: strings.Join(subdivs, ";"),
	}, nil
}
