package trafikverket_test

import (
	"tagtider/pkg/trafikverket"
	"testing"
)

// TestGenerateRequestPayload is a testing method for GenerateRequestPayload
func TestGenerateRequestPayload(t *testing.T) {
	apiKey := "testApiKey"
	currentDate := "2021-11-01"
	fromTime := "00:00:00"
	toTime := "23:59:59"
	locationSignature := "Kn"

	expectedOutput := `<REQUEST><LOGIN authenticationkey='testApiKey'/><QUERY  lastmodified='true' orderby='AdvertisedTimeAtLocation,EstimatedTimeAtLocation' objecttype='TrainAnnouncement' schemaversion='1.6' includedeletedobjects='false' sseurl='true'><FILTER><AND><EQ name='LocationSignature' value='Kn'/><EQ name='Advertised' value='true'/><EQ name='ActivityType' value='Avgang'/><OR><AND><GT name='AdvertisedTimeAtLocation' value='00:00:00'/><LT name='AdvertisedTimeAtLocation' value='23:59:59'/></AND><GT name='EstimatedTimeAtLocation' value='2021-11-01'/></OR></AND></FILTER></QUERY></REQUEST>`

	result := trafikverket.GenerateRequestPayload(apiKey, currentDate, fromTime, toTime, locationSignature)

	if result != expectedOutput {
		t.Errorf("Test failed: \nGot: \n%s \nExpected: \n%s", result, expectedOutput)
	}
}
