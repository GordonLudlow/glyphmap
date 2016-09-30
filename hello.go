package hello

import (
    "fmt"
    "io"   
    "net/http"
    "os"
    "encoding/json"
    
    "google.golang.org/appengine"
    "google.golang.org/appengine/log"
    
    //"google.golang.org/appengine/cloudsql"    
    "database/sql"
    _ "github.com/ziutek/mymysql/godrv"
)

var bucket = "runmap-140616.appspot.com"
var ctx context.Context
var db *DB

func init() {
    http.HandleFunc("/", handler)
}

type coordinateList [][]float64

func handlePost(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var coords coordinateList   
    err := decoder.Decode(&coords)
    if err != nil {
        panic(err)
    }
    insert, err := db.Prepare("INSERT INTO portals VALUES (?, ?)") 
    if err != nil {
        panic(err)
    }    
    for i := range coords {
        _, err = insert.Run(coords[i][0], coords[i][1])
        log.Infof(ctx, "Added [%f,%f]", coords[i][0], coords[i][1])
    }   
}

func handler(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)
    db, err := sql.Open("mymysql", "cloudsql:runmap-140616:us-central1:portals*portals/web/ALL_LOWER_CASE_NO_UNDERSCORES")
    if err != nil {
        panic(err)
    }   
    if r.Method == "POST" {
        handlePost(w,r)
        //return
    }
    
    rows, err := db.Query("SELECT lat, lng FROM portals")
    if err != nil {
        log.Errorf(ctx, "db.Query: %v", err)
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