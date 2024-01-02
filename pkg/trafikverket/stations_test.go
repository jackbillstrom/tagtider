package trafikverket_test

import (
	"bytes"
	"io"
	"net/http"
	"tagtider/pkg/trafikverket"
	"testing"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockDo MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

// TestFetchTrainStations tests the FetchTrainStations function
func TestFetchTrainStations(t *testing.T) {
	jsonResponse := `{"RESPONSE":{"RESULT":[{"TrainStation":[{"AdvertisedLocationName":"Test Station","LocationSignature":"TS"}]}]}}`
	r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))

	mockClient := &MockClient{
		MockDo: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	stations, err := trafikverket.FetchTrainStations(mockClient, "testKey", "testUrl")
	if err != nil {
		t.Errorf("FetchTrainStations() returned an error: %v", err)
	}

	if len(stations) != 1 {
		t.Errorf("FetchTrainStations() = %v, want 1 station", len(stations))
	}

	if stations[0].AdvertisedLocationName != "Test Station" || stations[0].LocationSignature != "TS" {
		t.Errorf("FetchTrainStations() = %v, want station with name 'Test Station' and signature 'TS'", stations[0])
	}
}
