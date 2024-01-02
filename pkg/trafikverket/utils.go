package trafikverket

import "fmt"

// GenerateRequestPayload Generates the XML payload for the API request
func GenerateRequestPayload(apiKey, currentDate, fromTime, toTime, locationSignature string) string {
	payload := fmt.Sprintf(`<REQUEST><LOGIN authenticationkey='%s'/><QUERY  lastmodified='true' orderby='AdvertisedTimeAtLocation,EstimatedTimeAtLocation' objecttype='TrainAnnouncement' schemaversion='1.6' includedeletedobjects='false' sseurl='true'><FILTER><AND><EQ name='LocationSignature' value='%s'/><EQ name='Advertised' value='true'/><EQ name='ActivityType' value='Avgang'/><OR><AND><GT name='AdvertisedTimeAtLocation' value='%s'/><LT name='AdvertisedTimeAtLocation' value='%s'/></AND><GT name='EstimatedTimeAtLocation' value='%s'/></OR></AND></FILTER></QUERY></REQUEST>`, apiKey, locationSignature, fromTime, toTime, currentDate)
	return payload
}
