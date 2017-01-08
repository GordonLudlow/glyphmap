package hello

import (
    "fmt"
    "net/http"
    "net/url"
    "encoding/json"

    "google.golang.org/appengine"
    "google.golang.org/appengine/log"
    
    "database/sql"
    _ "github.com/ziutek/mymysql/godrv"
)

var bucket = "runmap-140616.appspot.com"

func init() {
    http.HandleFunc("/", handler)
}

type coordinateList [][]float64

func handlePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    ctx := appengine.NewContext(r)
    decoder := json.NewDecoder(r.Body)
    var city string
    var coords coordinateList   
    err := decoder.Decode(&city)
    if err != nil {
        panic(err)
    }
    err = decoder.Decode(&coords)
    if err != nil {
        panic(err)
    }
    insert, err := db.Prepare("INSERT INTO portals (lat,lng,city) VALUES (?, ?,?)") 
    if err != nil {
        panic(err)
    }
    var insertCount = 0
    for i := range coords {
        _, err = insert.Exec(coords[i][0], coords[i][1], city)
        if err == nil {
            insertCount++
        } else {
            log.Errorf(ctx, "error adding %f,%f: %v", coords[i][0], coords[i][1], err)
            fmt.Fprintf(w, "error adding %f,%f: %v", coords[i][0], coords[i][1], err)
        }
    }   
    fmt.Fprintf(w, "%i portals added", insertCount)
}

func handler(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("mymysql", "cloudsql:runmap-140616:us-central1:portals*portals/web/ALL_LOWER_CASE_NO_UNDERSCORES")
    if err != nil {
        panic(err)
    }
    defer db.Close()    
    if r.Method == "POST" {
        handlePost(w,r,db)
        return
    }

    var city string
    m, _ := url.ParseQuery(r.URL.RawQuery)
    city = m["city"][0]
    
    ctx := appengine.NewContext(r)
    select := fmt.Sprintf("SELECT lat, lng FROM portals WHERE city = '%s'", city)
    rows, err := db.Query(select)
    if err != nil {
        log.Errorf(ctx, "db.Query: %v", err)
        panic(err)
    }
    defer rows.Close()

    var leadingComma = ""
    for rows.Next() {
        var lat float64
        var lng float64
        if err := rows.Scan(&lat, &lng); err != nil {
            log.Errorf(ctx, "rows.Scan: %v", err)
            continue
        }
        fmt.Fprintf(w, "%s[%f,%f]", leadingComma, lat, lng)
        log.Infof(ctx, "%s[%f,%f]", leadingComma, lat, lng)
        leadingComma = ","
    }
    if err := rows.Err(); err != nil {
        log.Errorf(ctx, "Row error: %v", err)
    }
}