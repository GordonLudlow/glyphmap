package hello

import (
    "fmt"
    "io"
    //"io/ioutil"    
    "net/http"
    "os"
    
    //"golang.org/x/net/context"
    "google.golang.org/appengine"
    //"google.golang.org/appengine/file"
    "google.golang.org/appengine/log"
    //"google.golang.org/cloud/storage"  
    
    //"google.golang.org/appengine/cloudsql"    
    "database/sql"
    _ "github.com/ziutek/mymysql/godrv"
)

var bucket = "runmap-140616.appspot.com"

func init() {
    http.HandleFunc("/", handler)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    /*
    ctx := appengine.NewContext(r)
    if bucket == "" {
        var err error
        if bucket, err = file.DefaultBucketName(ctx); err != nil {
            log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
            return
        }
    }
    */

    f, err := os.OpenFile("portals.json", os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    _, err = io.Copy(f, r.Body)
    if err != nil {
        panic(err)
    }      
}

func handler(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)
    db, err := sql.Open("mymysql", "cloudsql:runmap-140616:us-central1:portals*portals/web/ALL_LOWER_CASE_NO_UNDERSCORES")

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