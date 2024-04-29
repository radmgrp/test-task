package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	TimeRange int
	Hosts     []string
}

func main() {
	// Get initial config
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Start the HTTP server
	http.HandleFunc("/ping", pingHandler)
	go func() {
		if err := http.ListenAndServe(config.Port, nil); err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			//clear envs
			os.Clearenv()
			// reload envs (to be aple apply changes in configmap without pod restarting)
			config, err = loadConfig()
			if err != nil {
				log.Fatalf("error loading config: %v", err)
			}
			// // Invoke other hosts' GET /ping route every Nth amount of time
			for _, host := range config.Hosts {
				hostName := strings.Split(host, ":")
				// get all ips for the headless service
				ips, err := net.LookupIP(hostName[0])
				fmt.Printf("%v ips connected to the service %v found.\n", len(ips), hostName[0])
				// make ping
				for _, ip := range ips {
					fmt.Printf("contact ip: %v\n", ip.String()+":"+hostName[1])
					resp, err := http.Get(fmt.Sprintf("http://%s/ping", ip.String()+":"+hostName[1]))
					if err != nil {
						fmt.Printf("Error connecting to %s: %v\n", host, err)
						continue
					}
					defer resp.Body.Close()
					// output to stdout
					fmt.Printf("Response from %s: %s\n", ip.String()+":"+hostName[1], readResponseBody(resp))
				}

				if err != nil {
					fmt.Printf("unexpected error while looking up ips: %v", err)
				}
			}
			// output to stdout Timeout message
			fmt.Println("Wait " + strconv.Itoa(config.TimeRange) + " seconds")
			time.Sleep(time.Duration(config.TimeRange) * time.Second)
		}
	}()
	// initial params output
	fmt.Println("Server started at :" + config.Port)
	fmt.Println("Ping period: " + strconv.Itoa(config.TimeRange))
	select {}
}

// Helper function to load config

func loadConfig() (Config, error) {

	err := godotenv.Load("/config/.env")
	//err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := ":" + os.Getenv("port")
	hostsStr := os.Getenv("hosts")
	timeRange, err := strconv.Atoi(os.Getenv("timeRange"))

	hostsSlice := strings.Split(hostsStr, ",")
	var hosts []string
	for _, h := range hostsSlice {
		hosts = append(hosts, strings.TrimSpace(h))
	}

	config := Config{
		Port:      port,
		TimeRange: timeRange,
		Hosts:     hosts,
	}

	return config, nil

}

// Handler for GET /ping route
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

// Helper function to read response body

func readResponseBody(resp *http.Response) string {
	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	return string(buf[:n])
}
