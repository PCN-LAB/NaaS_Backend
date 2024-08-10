package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/database"
	"main/models"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Keyword struct {
	Word string `json:"word"`
}

func GetKeywords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var dates map[string]string
	if err := json.NewDecoder(r.Body).Decode(&dates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate := dates["startDate"]
	endDate := dates["endDate"]
	source := dates["source"]

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var keywords []Keyword

	if source == "Dawn" {

		rows, err := db.Query(`
           SELECT (d.topics)
            FROM news_dawn as d
            WHERE d.focus_time >= $1 AND d.focus_time <= $2
            
        `, startDate, endDate)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var word string
			if err := rows.Scan(&word); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			keywords = append(keywords, Keyword{Word: word})
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(keywords); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if source == "Tribune" {
		rows, err := db.Query(`
			SELECT (d.topics)
            FROM news_tribune as d
            WHERE d.focus_time >= $1 AND d.focus_time <= $2
    `, startDate, endDate)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var word string
			if err := rows.Scan(&word); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			keywords = append(keywords, Keyword{Word: word})
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(keywords); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getKeywords(w http.ResponseWriter, r *http.Request) {
	var req models.Request

	params := mux.Vars(r)
	param := params["keywords"]

	json.Unmarshal([]byte(param), &req)
	fmt.Println(req)

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var (
		header         string
		focus_time     string
		focus_location string
		link           string
		category       string
		coordinates    string
		location_type  string
		topics         string
		sentiment      string
		creationDate   string
		RacialBias     string
		HateSpeech     string
		PoliticalBias  string
	)

	var rows *sql.Rows
	if req.Location == "" {
		if len(req.Topics) > 0 {
			// fmt.Println("Test 1")
			rows, err = db.Query(`
                SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.sentiment, n.creation_date,n.racial_bias,n.hate_speech,n.political_bias,
                CASE 
                    WHEN n.location_type = 'Province' THEN p.name 
                    WHEN n.location_type = 'District' THEN d.name
                    WHEN n.location_type = 'Tehsil' THEN t.name 
                    ELSE NULL 
                END AS location, 
                CASE 
                    WHEN n.location_type = 'Province' THEN p.coordinates 
                    WHEN n.location_type = 'District' THEN d.coordinates
                    WHEN n.location_type = 'Tehsil' THEN t.coordinates 
                    ELSE NULL 
                END AS coordinates 
                FROM 
                news_dawn n 
                LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
                LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
                LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
                WHERE n.focus_time::date BETWEEN $1 AND $2 AND $3 && n.topics;
            `, req.StartDate, req.EndDate, pq.Array(req.Topics))
		} else {
			// fmt.Println("Test 2")
			rows, err = db.Query(`
                SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.sentiment, n.creation_date,n.racial_bias,n.hate_speech,n.political_bias,
                CASE 
                    WHEN n.location_type = 'Province' THEN p.name 
                    WHEN n.location_type = 'District' THEN d.name
                    WHEN n.location_type = 'Tehsil' THEN t.name 
                    ELSE NULL 
                END AS location, 
                CASE 
                    WHEN n.location_type = 'Province' THEN p.coordinates 
                    WHEN n.location_type = 'District' THEN d.coordinates
                    WHEN n.location_type = 'Tehsil' THEN t.coordinates 
                    ELSE NULL 
                END AS coordinates
                FROM 
                news_dawn n 
                LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
                LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
                LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
                WHERE n.focus_time::date BETWEEN $1 AND $2;
            `, req.StartDate, req.EndDate)

			// printing the row
			fmt.Println(rows)
		}
	} else {
		// fmt.Println("Test 3")
		err = db.QueryRow("SELECT location_type FROM locations WHERE name = $1", req.Location).Scan(&location_type)
		if len(req.Topics) > 0 {
			query := fmt.Sprintf(`
                SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.sentiment, n.creation_date,n.racial_bias,n.hate_speech,n.political_bias,
                CASE 
                    WHEN n.location_type = 'Province' THEN p.name 
                    WHEN n.location_type = 'District' THEN d.name
                    WHEN n.location_type = 'Tehsil' THEN t.name 
                    ELSE NULL 
                END AS location,
                CASE 
                    WHEN n.location_type = 'Province' THEN p.coordinates 
                    WHEN n.location_type = 'District' THEN d.coordinates
                    WHEN n.location_type = 'Tehsil' THEN t.coordinates 
                    ELSE NULL 
                END AS coordinates 
                FROM 
                news_dawn n 
                LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
                LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
                LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
                WHERE n.%s = $1 AND n.focus_time::date BETWEEN $2 AND $3 AND $4 && n.topics;
            `, location_type)

			rows, err = db.Query(query, req.Location, req.StartDate, req.EndDate, pq.Array(req.Topics))
		} else {
			// fmt.Println("Test 4")
			query := fmt.Sprintf(`
                SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.sentiment, n.creation_date,n.racial_bias,n.hate_speech,n.political_bias,
                CASE 
                    WHEN n.location_type = 'Province' THEN p.name 
                    WHEN n.location_type = 'District' THEN d.name
                    WHEN n.location_type = 'Tehsil' THEN t.name 
                    ELSE NULL 
                END AS location,
                CASE 
                    WHEN n.location_type = 'Province' THEN p.coordinates 
                    WHEN n.location_type = 'District' THEN d.coordinates
                    WHEN n.location_type = 'Tehsil' THEN t.coordinates 
                    ELSE NULL 
                END AS coordinates 
                FROM 
                news_dawn n 
                LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
                LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
                LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
                WHERE n.%s = $1 AND n.focus_time::date BETWEEN $2 AND $3;
            `, location_type)

			rows, err = db.Query(query, req.Location, req.StartDate, req.EndDate)
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var news []models.NewsData
	for rows.Next() {

		rows.Scan(&header, &focus_time, &category, &link, &topics, &location_type, &sentiment, &creationDate, &RacialBias, &HateSpeech, &PoliticalBias, &focus_location, &coordinates)

		var temp models.NewsData
		temp.Header = header
		temp.FocusLocation = focus_location
		temp.FocusTime = focus_time[:10]
		temp.Coordinates = coordinates
		temp.LocationType = location_type
		temp.Topics = topics
		temp.Category = category
		temp.Link = link
		temp.Sentiment = sentiment
		temp.CreationDate = creationDate[:10]
		temp.RacialBias = RacialBias
		temp.HateSpeech = HateSpeech
		temp.PoliticalBias = PoliticalBias

		news = append(news, temp)
	}

	json.NewEncoder(w).Encode(news)

}