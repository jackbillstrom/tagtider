package trafikverket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FetchTrainStations makes a callAPI call and returns the train stations
func FetchTrainStations(client HttpClient, apiKey, url string) ([]TrainStation, error) {
	requestBody := `<REQUEST>
<LOGIN authenticationkey='` + apiKey + `'/>
<QUERY objecttype='TrainStation' schemaversion='1'>
<FILTER/>
</QUERY>
</REQUEST>`

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response
	var data struct {
		RESPONSE struct {
			RESULT []struct {
				TrainStation []TrainStation `json:"TrainStation"`
			} `json:"RESULT"`
		} `json:"RESPONSE"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if len(data.RESPONSE.RESULT) == 0 || len(data.RESPONSE.RESULT[0].TrainStation) == 0 {
		return nil, fmt.Errorf("no train stations found")
	}

	return data.RESPONSE.RESULT[0].TrainStation, nil
}
