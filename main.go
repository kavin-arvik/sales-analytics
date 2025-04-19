package main

import (
	db "SalesAnalytics/DB"
	"SalesAnalytics/handlers"
	"SalesAnalytics/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

func autoRestart() {
	for {
		now := time.Now()
		//resart the program everyday at 4am
		//at 3am, the program goes for 1 minute sleep and after that it will restart
		if now.Hour() == 3 {
			//sleep for an minute
			//in the loop does not continue again in next iteration
			time.Sleep(60 * 5 * time.Second)
			fmt.Println(now.Hour(), now.Minute(), now.Second())
			log.Println(now.Hour(), now.Minute(), now.Second())
			// Restart the program
			fmt.Println("Restarting the program...")
			log.Println("Restarting the program...")
			execPath, err := os.Executable()
			if err != nil {
				fmt.Println("Error getting executable path:", err)
				log.Println("Error getting executable path:", err)
				return
			}
			cmd := exec.Command(execPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Start()
			if err != nil {
				fmt.Println("Error restarting program:", err)
				log.Println("Error restarting program:", err)
				return
			}
			os.Exit(0)
		}
		time.Sleep(60 * 30 * time.Second)
	}
}

func main() {
	log.Println("Server Started on http://localhost:8080")

	f, err := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	//Automatically restarting the program for everyday @ 3am
	go autoRestart()

	//Automatic refresh logic for everyday
	go handlers.AutoRefresh()

	lErr := db.BuildConnecrtion()
	if lErr != nil {
		log.Println("Error while executing DB connection.")
	} else {
		//Close the DB connection at the end of main function
		defer db.GDBConnection.Close()
	}
	r := mux.NewRouter()
	router.RegisterRoutes(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
