package geoloc

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

//GetLocalization form ip
func GetLocalization(ip string) (string, error) {
	db, err := geoip2.Open("/home/kuba/Documents/gitfolders/QuicPos-Server/internal/geoloc/GeoLite2-City.mmdb")
	if err != nil {
		return "", err
	}
	//defer db.Close()
	ipParsed := net.ParseIP(ip)
	city, err := db.City(ipParsed)
	if err != nil {
		return "", err
	}
	subdivisions := ""
	for _, subdivision := range city.Subdivisions {
		subdivisions += subdivision.Names["en"] + " "
	}
	return city.Continent.Names["en"] + " " + city.Country.Names["en"] + " " + subdivisions + city.City.Names["en"], nil
}
