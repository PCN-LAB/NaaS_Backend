package models

import "time"

// Defining different structs that'll be used in the api
type InputVars struct {
	Timeframe string `json:"timeframe"`
	Location  string `json:"location"`
	Keywords  string `json:"keywords"`
}

type Location struct {
	Name          string `json:"name"`
	Location_type string `json:"location_type"`
}

type InitialData struct {
	StartTime string     `json:"startTime"`
	EndTime   string     `json:"endTime"`
	Location  []Location `json:"locations"`
	Topics    []string   `json:"topics"`
}

type Request struct {
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Location  string   `json:"location"`
	Topics    []string `json:"topics"`
}

type DataRequest struct {
	Source string `json:"source"`
}

type NewsData struct {
	FocusTime     string `json:"focusTime"`
	FocusLocation string `json:"focusLocation"`
	Header        string `json:"header"`
	Link          string `json:"link"`
	Category      string `json:"category"`
	Coordinates   string `json:"coordinates"`
	Topics        string `json:"topics"`
	LocationType  string `json:"locationType"`
	Sentiment     string `json:"sentiment"`
	CreationDate  string `json:"creationDate"`
	RacialBias    string `json:"racial_bias"`
	HateSpeech    string `json:"hate_speech"`
	PoliticalBias string `json:"political_bias"`
}

type NewsDataTribune struct {
	FocusTime     string `json:"focusTime"`
	FocusLocation string `json:"focusLocation"`
	Header        string `json:"header"`
	Link          string `json:"link"`
	Category      string `json:"category"`
	Coordinates   string `json:"coordinates"`
	Topics        string `json:"topics"`
	LocationType  string `json:"locationType"`
	Picture       string `json:"picture"`
	Sentiment     string `json:"sentiment"`
	CreationDate  string `json:"creationDate"`
	RacialBias    string `json:"racial_bias"`
	HateSpeech    string `json:"hate_speech"`
	PoliticalBias string `json:"political_bias"`
}

type News struct {
	Header        string    `json:"header"`
	Category      string    `json:"category"`
	Link          string    `json:"link"`
	Topics        []string  `json:"topics"`
	FocusLocation string    `json:"focus_location"`
	FocusTime     string    `json:"focus_time"`
	LocationType  string    `json:"location_type"`
	Sentiment     string    `json:"sentiment"`
	CreationDate  time.Time `json:"creation_date"`
	Province      string    `json:"province"`
	District      string    `json:"district"`
	RacialBias    string    `json:"racial_bias"`
	HateSpeech    string    `json:"hate_speech"`
	PoliticalBias string    `json:"political_bias"`
}
type SearchParams struct {
	CreationDateStart time.Time
	CreationDateEnd   time.Time
}


type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type SignupInfo struct {
    Email      string `json:"email"`
    Password   string `json:"password"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Occupation string `json:"occupation"`
}

