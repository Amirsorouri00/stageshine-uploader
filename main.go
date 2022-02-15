
// main.go
package main

import (
    "log"
    "net/http"
    "fmt"
    "os"
    "os/exec"
    "io"
    "path/filepath"


    "github.com/gorilla/mux"
)

// This is where we must upload the 
// file and save it into a filder with
// name = /videos/{video.uuid}/*.{mp4, mkv, ...}
func uploadNewVideo(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    uuid := vars["uuid"]

    r.ParseMultipartForm(32 << 20)

    file, handler, err := r.FormFile("file") // Retrieve the file from form data
    if err != nil {
        fmt.Printf("%+v\n", err)
    }
    defer file.Close()                       // Close the file when we finish
    fmt.Printf("%+v\n%s", handler,uuid)
    extension := filepath.Ext(handler.Filename)


    // This is path which we want to store the file
    f, err := os.OpenFile("/Users/amirhossiensorouri/Desktop/projects/stageshine/create-rest-api-in-go-tutorial/videos/"+uuid+extension, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Printf("%+v\n", err)
    }

    // Copy the file to the destination path
    io.Copy(f, file)
    
    // this is where we must run our script to 
    // transcode our uploaded video into .m3u8 file.

    cmd := exec.Command("/Users/amirhossiensorouri/Desktop/projects/stageshine/transcoder/transcoder", uuid+extension)
    cmd.Stdout = os.Stdout
    err = cmd.Start()
    if err != nil {
    log.Fatal(err)
    }
    log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)

    // done.
}

func handleRequests() {
    addr := ":10000"

    myRouter := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/video/{uuid}", uploadNewVideo).Methods("POST")
    
    log.Println("listen on", addr)
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    handleRequests()
}
