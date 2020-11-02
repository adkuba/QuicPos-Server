package geoloc

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

//GetLocalization form ip
func GetLocalization(ip string) (string, float64, float64, error) {
	//plik jest wykonywany z glownej sciezki dlatego mam internal geoloc
	db, err := geoip2.Open("./internal/geoloc/GeoLite2-City.mmdb")
	if err != nil {
		return "", -1, -1, err
	}
	//defer db.Close()
	ipParsed := net.ParseIP(ip)
	city, err := db.City(ipParsed)
	if err != nil {
		return "", -1, -1, err
	}
	subdivisions := ""
	for _, subdivision := range city.Subdivisions {
		subdivisions += subdivision.Names["en"] + " "
	}
	locString := city.Continent.Names["en"] + " " + city.Country.Names["en"] + " " + subdivisions + city.City.Names["en"]
	return locString, city.Location.Latitude, city.Location.Longitude, nil
}
