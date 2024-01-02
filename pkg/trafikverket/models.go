package trafikverket

import (
	"encoding/json"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func UnmarshalResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Response is the response from the Trafikverket API
type Response struct {
	Response RESPONSEClass `json:"RESPONSE"`
}

// RESPONSEClass is the response from the Trafikverket API
type RESPONSEClass struct {
	Result []Result `json:"RESULT"`
}

// Result is the result from the Trafikverket API
type Result struct {
	TrainAnnouncement []TrainAnnouncement `json:"TrainAnnouncement,omitempty"`
	TrainStation      []TrainStation      `json:"TrainStation,omitempty"`
	Info              Info                `json:"INFO,omitempty"`
}

// Info is the info from the Trafikverket API
type Info struct {
	Lastmodified Lastmodified `json:"LASTMODIFIED"`
	Sseurl       string       `json:"SSEURL"`
}

// Lastmodified is the lastmodified from the Trafikverket API
type Lastmodified struct {
	AttrDatetime string `json:"_attr_datetime"`
}

// TrainAnnouncement is the train announcement from the Trafikverket API
type TrainAnnouncement struct {
	ActivityID                            string            `json:"ActivityId"`
	ActivityType                          ActivityType      `json:"ActivityType"`
	Advertised                            bool              `json:"Advertised"`
	AdvertisedTimeAtLocation              string            `json:"AdvertisedTimeAtLocation"`
	AdvertisedTrainIdent                  string            `json:"AdvertisedTrainIdent"`
	Canceled                              bool              `json:"Canceled"`
	Deleted                               bool              `json:"Deleted"`
	EstimatedTimeIsPreliminary            bool              `json:"EstimatedTimeIsPreliminary"`
	FromLocation                          []Location        `json:"FromLocation"`
	InformationOwner                      InformationOwner  `json:"InformationOwner"`
	LocationSignature                     LocationSignature `json:"LocationSignature"`
	ModifiedTime                          string            `json:"ModifiedTime"`
	NewEquipment                          int64             `json:"NewEquipment"`
	Operator                              Operator          `json:"Operator"`
	EstimatedTimeAtLocation               string            `json:"EstimatedTimeAtLocation,omitempty"`
	PlannedEstimatedTimeAtLocationIsValid bool              `json:"PlannedEstimatedTimeAtLocationIsValid"`
	ProductInformation                    []Booking         `json:"ProductInformation"`
	ScheduledDepartureDateTime            string            `json:"ScheduledDepartureDateTime"`
	TechnicalDateTime                     string            `json:"TechnicalDateTime"`
	TechnicalTrainIdent                   string            `json:"TechnicalTrainIdent"`
	ToLocation                            []Location        `json:"ToLocation"`
	TrackAtLocation                       string            `json:"TrackAtLocation"`
	TrainOwner                            TrainOwner        `json:"TrainOwner"`
	TypeOfTraffic                         []Booking         `json:"TypeOfTraffic"`
	WebLink                               string            `json:"WebLink"`
	WebLinkName                           WebLinkName       `json:"WebLinkName"`
	MobileWebLink                         *string           `json:"MobileWebLink,omitempty"`
	OtherInformation                      []Booking         `json:"OtherInformation,omitempty"`
	ViaToLocation                         []Location        `json:"ViaToLocation,omitempty"`
	Booking                               []Booking         `json:"Booking,omitempty"`
	TrainComposition                      []Booking         `json:"TrainComposition,omitempty"`
}

// TrainStation is the train station from the Trafikverket API
type TrainStation struct {
	Advertised                  bool     `json:"Advertised"`
	AdvertisedLocationName      string   `json:"AdvertisedLocationName"`
	AdvertisedShortLocationName string   `json:"AdvertisedShortLocationName"`
	CountyNo                    []int64  `json:"CountyNo,omitempty"`
	Deleted                     bool     `json:"Deleted"`
	Geometry                    Geometry `json:"Geometry"`
	LocationSignature           string   `json:"LocationSignature"`
	PlatformLine                []string `json:"PlatformLine,omitempty"`
	Prognosticated              bool     `json:"Prognosticated"`
	ModifiedTime                string   `json:"ModifiedTime"`
	LocationInformationText     *string  `json:"LocationInformationText,omitempty"`
}

// Geometry is the geometry from the Trafikverket API
type Geometry struct {
	Sweref99Tm *string `json:"SWEREF99TM,omitempty"`
	Wgs84      string  `json:"WGS84"`
}

type Booking struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type Location struct {
	LocationName string `json:"LocationName"`
	Priority     int64  `json:"Priority"`
	Order        int64  `json:"Order"`
}

type ActivityType string
type InformationOwner string
type LocationSignature string
type Operator string
type TrainOwner string
type WebLinkName string
