package trafikverket

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// attemptFetchAPIInfo attempts to fetch API details such as endpoint and apikey
func attemptFetchAPIInfo(client *http.Client, station string) (string, string, error) {
	uri := "https://www.trafikverket.se/trafikinformation/tag/?Station=" + url.QueryEscape(station) + "&ArrDep=departure"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", "", err
	}

	// Headers for the scraping request
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,sv;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "TrvCookieConsent=functional=true&analytical=true; ASP.NET_SessionId=hsqdbwx4lihatzvs4hbtpp23; extweb=ffffffff09140a4c45525d5f4f58455e445a4a42378b; _pk_id.14.3162=2da6cfd35441c18a.1704146770.; _pk_ses.14.3162=1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Read HTML body
	body, _ := io.ReadAll(resp.Body)

	// Create a goquery document from the HTTP response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", "", err
	}

	// Find the API key
	apiKey := doc.Find("traincomponent").AttrOr("apikey", "")

	// Find the API URL
	apiURL := doc.Find("traincomponent").AttrOr("apiurl", "")

	return apiKey, apiURL, nil
}

// FetchAPIInfo gets API details such as endpoint and apikey
func FetchAPIInfo(client *http.Client, station string) (string, string, error) {
	// Antal försök
	maxRetries := 3

	var err error
	for i := 0; i < maxRetries; i++ {
		var apiKey, apiURL string
		apiKey, apiURL, err = attemptFetchAPIInfo(client, station)
		if err == nil {
			return apiKey, apiURL, nil
		}
		time.Sleep(time.Second * 2) // Vänta i 2 sekunder mellan försök
	}
	return "", "", err
}

// MakeAPICall posts and queries the Trafikverket API
func MakeAPICall(client *http.Client, apiURL, requestBody string) string {
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return ""
	}

	req.Header.Set("authority", "api.trafikinfo.trafikverket.se")
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,sv;q=0.8")
	req.Header.Set("content-type", "text/plain")
	req.Header.Set("origin", "https://www.trafikverket.se")
	req.Header.Set("referer", "https://www.trafikverket.se")
	req.Header.Set("sec-ch-ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(responseBody)
}
