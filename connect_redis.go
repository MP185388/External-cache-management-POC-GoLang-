package main

import (
		"fmt"
		"github.com/go-redis/redis/v7"
		"encoding/json"
		"io/ioutil"
        "log"
        "net/http"
        "github.com/gorilla/mux"
		
       )


type transaction struct{
	TxnID          string `json:"ID"`
		Itemcode       string  `json:"Itemcode"`
		CustomerID     string `json:"CustomerID"`
		OrgID	       string `json:"OrgID`
}

func createEvent(w http.ResponseWriter, r *http.Request) {
        var newtransaction transaction
                reqBody, err := ioutil.ReadAll(r.Body)
                if err != nil {
                        fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
                }

        json.Unmarshal(reqBody, &newtransaction)
		w.WriteHeader(http.StatusCreated)

        json.NewEncoder(w).Encode(newtransaction)

		Setval(newtransaction)
}

func Setval(trans transaction){

				client := rClient()
				json, err := json.Marshal(&trans)
				if err != nil {
					panic(err)
				}
				
				client.Set(trans.TxnID, json, 0)
}


func Getval(w http.ResponseWriter, r *http.Request)   {
		        var newtransaction transaction
				TransactionID := mux.Vars(r)["id"]
				client := rClient()
				
				val, errget := client.Get(TransactionID).Result()
                if errget != nil {
                	fmt.Println(errget)
                }
				//finval := json.Unmarshal([]byte(val), &newtransaction)
	err := json.Unmarshal([]byte(val), &newtransaction)
	if err != nil {
		panic(err)
	}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(val)
				
}               				
func rClient() *redis.Client {
				client := redis.NewClient(&redis.Options{
										   Addr: "localhost:6379",
										   DB:       6,
										 })

	return client
}

func main(){
			
			router := mux.NewRouter().StrictSlash(true)
			router.HandleFunc("/events/{id}", Getval).Methods("GET")
            router.HandleFunc("/event", createEvent).Methods("POST")
			log.Fatal(http.ListenAndServe(":8060", router))

	rClient()
}
