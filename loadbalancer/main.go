package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func heartbeatHandler(instances *map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend alive"))

		//type HeartbeatRequest struct {
		//	InstanceID string `json:"instance_id"`
		//}
		//var req HeartbeatRequest
		//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//	log.Printf("Error decoding heartbeat request: %v", err)
		//	return
		//}

		//remove port from remote address
		address := strings.Split(r.RemoteAddr, ":")[0]
		(*instances)[address] = true
		log.Printf("Received heartbeat from %s", r.RemoteAddr)

	}
}

// instanceID identifies which container answered a request. Set via the
// INSTANCE_ID env var in docker-compose.yml so you can watch the gateway's
// load balancer spread requests across backend-1/2/3.
var instanceID = "unknown"

// withInstanceHeader tags every response with X-Instance-Id so a client
// (or curl -i) can see which backend instance actually served it.
func withInstanceHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("LB-Instance-Id", instanceID)
		next(w, r)
	}
}

func HeartbeatRequests(instances *map[string]bool) {
	go func() {
		for {
			if len(*instances) == 0 {
				log.Println("No backend instances registered yet")
			} else {
				for id := range *instances {
					resp, err := http.Get("http://" + id + ":8080/heartbeat")
					if err != nil {
						log.Printf("Error checking heartbeat for %s: %v", id, err)
						(*instances)[id] = false
						continue
					}
					if resp.StatusCode != http.StatusOK {
						(*instances)[id] = false
						log.Printf("Backend %s is not alive", id)
					}
					if resp.StatusCode == http.StatusOK {
						log.Printf("Backend %s is alive", id)
					}
				}

			}
			for id, active := range *instances {
				log.Printf("Instance %s, is alive: %t", id, active)
			}
			time.Sleep(20 * time.Second)
		}
	}()
}

func loadBalancerHandler(instances *map[string]bool, instanceNum int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(*instances) == 0 {
			http.Error(w, "No backend instances available", http.StatusServiceUnavailable)
			return
		}

		// Simple round-robin load balancing
		var instanceIDs []string
		for id := range *instances {
			instanceIDs = append(instanceIDs, id)
		}
		sort.Strings(instanceIDs)
		instanceNum = (instanceNum + 1) % len(instanceIDs)
		proxyURL := "http://" + instanceIDs[instanceNum] + ":8080" + r.URL.Path
		req, err := http.NewRequest(r.Method, proxyURL, r.Body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}
		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error forwarding request to: "+proxyURL, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Copy the response headers and body back to the original client
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		io.Copy(w, resp.Body)

		w.WriteHeader(resp.StatusCode)
	}
}

func instancesHandler(instances *map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(*instances)
	}
}

func main() {
	if v := os.Getenv("INSTANCE_ID"); v != "" {
		instanceID = v
	}

	instances := make(map[string]bool)
	var instanceNum int = 0

	HeartbeatRequests(&instances)

	mux := http.NewServeMux()
	mux.HandleFunc("/", withInstanceHeader(loadBalancerHandler(&instances, instanceNum)))
	mux.HandleFunc("/instances", withInstanceHeader(instancesHandler(&instances)))
	mux.HandleFunc("/heartbeat", withInstanceHeader(heartbeatHandler(&instances)))

	log.Printf("backend %s listening on :8080", instanceID)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
