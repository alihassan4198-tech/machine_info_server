package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

const timeLayout = "Jan 2, 2006 at 3:04pm (MST)"

var DBName string

/**
* Global logger
 */
var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
}

// uploader
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	dst, err := os.Create("aaa")
	if err != nil {
		logger.Printf("Error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// copy each part to destination
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Printf("Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// if part.FileName() is empty, skip this iteration
		if part.FileName() == "" {
			continue
		}

		var user_dir string
		if runtime.GOOS == "linux" {

			user_dir = "/home/machineinfoserver/"

		} else {
			user_dir = "/Users/Shared/machineinfoserver/"
		}

		// Check if dir exists, if not create it
		if _, err := os.Stat(user_dir); os.IsNotExist(err) {
			err := os.Mkdir(user_dir, 0750)
			if err != nil {
				logger.Printf("Error: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Wrap dst file creation and copy in function with immediate execution, so when
		// it returns the deferred dst.Close() is called
		err = func() error {
			dst, err := os.Create(user_dir + part.FileName())
			if err != nil {
				return err
			}
			defer dst.Close()

			logger.Println(user_dir+part.FileName(), "has been created")

			if _, err := io.Copy(dst, part); err != nil {
				return err
			}

			return nil
		}()

		if err != nil {
			logger.Printf("Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add the files to the user's file db.
		file_name := part.FileName()

		// Change Permissions
		os.Chmod(file_name, 0600)
	}
}

func main() {
	// Define a custom port from .env file.
	var port string = os.Getenv("PORT")
	if port == "" {
		logger.Println("no port name provided, using 3010")
		port = "3010"
	}
	logger.Println("Using port: ", port)
	logger.Println("Starting Uploader...")

	// Start Web Server
	http.HandleFunc("/uploadfile", uploadFileHandler)
	logger.Fatal(http.ListenAndServe(":"+port, nil))
}
