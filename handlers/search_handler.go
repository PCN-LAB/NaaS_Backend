package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/database"
	"main/models"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func searchDataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content=Type", "application/json")

	type RequestData struct {
		CreationDateStart string   `json:"creation_date_start"`
		CreationDateEnd   string   `json:"creation_date_end"`
		Keywords          []string `json:"keywords"`
	}

	var requestData RequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Construct the SQL query
	query := `
        SELECT n.header, n.category, n.link, n.location_type, n.sentiment, n.creation_date, n.province, n.district, n.focus_location, n.focus_time
        FROM News_Dawn as n 
        INNER JOIN keywords as k ON n.id = k.dawn_id
        WHERE n.creation_date >= $1 AND n.creation_date <= $2 AND k.word IN (`

	// Add placeholders for the IN clause
	placeholders := make([]string, len(requestData.Keywords))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+3)
	}

	query += strings.Join(placeholders, ",") + ")"

	// Execute the query
	args := make([]interface{}, len(requestData.Keywords)+2)
	args[0] = requestData.CreationDateStart
	args[1] = requestData.CreationDateEnd
	for i, keyword := range requestData.Keywords {
		args[i+2] = keyword
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Process the result set
	var newsList []models.News
	for rows.Next() {
		var news models.News
		// Scan each row into a News struct
		err := rows.Scan(&news.Header, &news.Category, &news.Link, &news.LocationType, &news.Sentiment, &news.CreationDate, &news.Province, &news.District, &news.FocusLocation, &news.FocusTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Append the News struct to the results slice
		newsList = append(newsList, news)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert results to JSON and write response
	jsonResponse, err := json.Marshal(newsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Function to get the timeframe when user inputs
func getTimeFrame(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	params := mux.Vars(r)
	param := params["timeframe"]

	json.Unmarshal([]byte(param), &req)
	fmt.Println(req.Location, req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content=Type", "application/json")

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
	)
	var rows *sql.Rows
	if req.Location == "" {
		if req.StartDate != "" && req.EndDate != "" {
			rows, err = db.Query("select distinct(n.header), n.focus_time::date, n.category, n.link, n.focus_location, concat (p.coordinates, d.coordinates) as coordinates from news as n left join district as d on d.name = n.focus_location left join province as p on p.name = n.focus_location where n.focus_time::date between $1 and $2 ;", req.StartDate, req.EndDate)
		}
	} else {
		err = db.QueryRow("select location_type from locations where name= $1", req.Location).Scan(&location_type)
		query := fmt.Sprintf("select distinct(n.header), n.focus_time::date,  n.category, n.link, n.focus_location, concat (p.coordinates, d.coordinates) as coordinates from news as n left join tehsil as t on t.name = n.focus_location left join district as d on d.name = n.focus_location left join province as p on p.name = n.focus_location where n.%s = $1;", location_type)
		query_time := fmt.Sprintf("select distinct(n.header), n.focus_time::date,  n.category, n.link, n.focus_location, concat (p.coordinates, d.coordinates) as coordinates from news as n left join tehsil as t on t.name = n.focus_location left join district as d on d.name = n.focus_location left join province as p on p.name = n.focus_location where n.%s = $1 and n.focus_time::date between $2 and $3 ;", location_type)
		if req.Location != "" && (req.EndDate == "" || req.StartDate == "") {
			rows, err = db.Query(query, req.Location)
		} else if req.Location != "" && req.EndDate != "" && req.StartDate != "" {
			// fmt.Println("Query was run till here")
			rows, err = db.Query(query_time, req.Location, req.StartDate, req.EndDate)
		}
	}

	if err != nil {
		log.Fatal(err)

	}
	defer rows.Close()
	var news []models.NewsData
	for rows.Next() {
		err := rows.Scan(&header, &focus_time, &category, &link, &focus_location, &coordinates)
		if err != nil {
			log.Fatal(err)
		}
		var temp models.NewsData
		temp.Header = header
		temp.FocusLocation = focus_location
		temp.FocusTime = focus_time[:10]
		temp.Coordinates = coordinates
		temp.Category = category
		temp.Link = link

		news = append(news, temp)
	}
	json.NewEncoder(w).Encode(news)
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	params := mux.Vars(r)
	param := params["location"]

	json.Unmarshal([]byte(param), &req)
	fmt.Println(req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content=Type", "application/json")

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
		picture        string
		sentiment      string
		creationDate   string
	)

	defer db.Close()
	var rows *sql.Rows
	if req.Location == "" {
		if len(req.Topics) > 0 {
			rows, err = db.Query(`   SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.picture, n.sentiment, n.creation_date,
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
			news_tribune n 
			LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
			LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
			LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
			where n.focus_time::date between $1 and $2 and $3 && n.topics;`, req.StartDate, req.EndDate, pq.Array(req.Topics))
		} else {
			rows, err = db.Query(`   SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.picture, n.sentiment, n.creation_date,
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
			news_tribune n 
			LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
			LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
			LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
			where n.focus_time::date between $1 and $2;`, req.StartDate, req.EndDate)
		}
	} else {
		err = db.QueryRow("select location_type from locations where name= $1", req.Location).Scan(&location_type)
		if len(req.Topics) > 0 {
			query := fmt.Sprintf(`   SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.picture, n.sentiment, n.creation_date,
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
			news_tribune n 
			LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
			LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
			LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
			where n.%s = $1 and n.focus_time::date between $2 and $3 and $4 && n.topics;`, location_type)

			rows, err = db.Query(query, req.Location, req.StartDate, req.EndDate, pq.Array(req.Topics))
		} else {
			query := fmt.Sprintf(`   SELECT n.header, n.focus_time::date, n.category, n.link, n.topics, n.location_type, n.picture, n.sentiment, n.creation_date,
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
			news_tribune n 
			LEFT JOIN province p ON n.focus_location = p.name AND n.location_type = 'Province' 
			LEFT JOIN district d ON n.focus_location = d.name AND n.location_type = 'District'
			LEFT JOIN tehsil t ON n.focus_location = t.name AND n.location_type = 'Tehsil'
			where n.%s = $1 and n.focus_time::date between $2 and $3;`, location_type)
			fmt.Println(query)

			rows, err = db.Query(query, req.Location, req.StartDate, req.EndDate)
		}
	}
	if err != nil {
		log.Fatal(err)

	}
	defer rows.Close()
	var news []models.NewsDataTribune
	for rows.Next() {
		err := rows.Scan(&header, &focus_time, &category, &link, &topics, &location_type, &picture, &sentiment, &creationDate, &focus_location, &coordinates)
		if err != nil {
			log.Fatal(err)
		}
		var temp models.NewsDataTribune
		temp.Header = header
		temp.FocusLocation = focus_location
		temp.FocusTime = focus_time[:10]
		temp.Coordinates = coordinates
		temp.LocationType = location_type
		temp.Topics = topics
		temp.Category = category
		temp.Link = link
		temp.Picture = picture
		temp.Sentiment = sentiment
		temp.CreationDate = creationDate[:10]
		// fmt.Println(coordinates)
		// json.Unmarshal([]byte(coordinates), &temp.Coordinates)
		news = append(news, temp)
	}
	json.NewEncoder(w).Encode(news)
}

// Function to get all the data from the user in the initialData struct and generate a query based on that data
func getInitialData(w http.ResponseWriter, r *http.Request) {
	var req models.DataRequest

	params := mux.Vars(r)
	param := params["initialData"]

	json.Unmarshal([]byte(param), &req)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content=Type", "application/json")

	// Connect to the database
	db, err := database.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// var data initialData
	var (
		name          string
		startTime     string
		endTime       string
		topic         string
		frequency     int
		location_type string
	)
	query := fmt.Sprintf("SELECT distinct(focus_location), location_type from news_%s where focus_location is not null;", req.Source)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)

	}

	defer rows.Close()
	var poke models.InitialData
	for rows.Next() {
		err := rows.Scan(&name, &location_type)
		if err != nil {
			log.Fatal(err)
		}
		var location models.Location
		location.Name = name
		location.Location_type = location_type
		poke.Location = append(poke.Location, location)
	}
	query = fmt.Sprintf("Select min(focus_time)::date, max(focus_time)::date from NEWS_%s;", req.Source)
	rows, err = db.Query(query)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&startTime, &endTime)
		if err != nil {
			log.Fatal(err)
		}
		startTime = startTime[:10]
		endTime = endTime[:10]
	}
	poke.StartTime = startTime
	poke.EndTime = endTime
	query = fmt.Sprintf(`SELECT 
	topic, COUNT(*) AS frequency 
  FROM 
	(SELECT 
	  (UNNEST(topics)) AS topic 
	FROM 
	  news_%s n) 
	AS extracted_topics 
  GROUP BY topic 
  ORDER BY frequency DESC;`, req.Source)
	rows, err = db.Query(query)
	if err != nil {
		log.Fatal(err)

	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&topic, &frequency)
		if err != nil {
			log.Fatal(err)
		}
		poke.Topics = append(poke.Topics, topic)
	}
	// fmt.Println(poke)
	json.NewEncoder(w).Encode(poke)
}
