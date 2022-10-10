package main

import (
    "fmt"
    "net/http"
    "strconv"
    "encoding/json"
    "io/ioutil"
    "sync"
    "sort"
    "os"
    "strings"
    "time"
)

var retryTimes = []time.Duration{
    1 * time.Second,
    3 * time.Second,
    10 * time.Second,
}

type Server struct {
}

func get_url_dataWRetries (url string, aData *[]JsonItems, dataLoc int, wg *sync.WaitGroup) {
    var resBody []byte
    var err error
    var res *http.Response

    for _, tryTime := range retryTimes {
        res, resBody, err = get_url_data(url)

        if err == nil {
	    if res.StatusCode == 200 {
                break
            } else {
                fmt.Fprintf(os.Stderr, "Status code error for %s\n", url)
	    }
	} else {
           fmt.Fprintf(os.Stderr, "Request error for %s: %+v\n", url, err)
	}

        fmt.Fprintf(os.Stderr, "Retrying for %s in %v\n", url, tryTime)
        time.Sleep(tryTime)
    }

    // All retries failed
    if err != nil || res.StatusCode != 200 {
        fmt.Fprintf(os.Stderr, "Max retries reached, was not able to retrieve data from %s\n", url)
        wg.Done()
	return
    }

    var data JsonItems
    if err := json.Unmarshal(resBody, &data); err != nil {
        fmt.Fprintf(os.Stderr, "Error reading retrieved data for %s\n", url)
	return
    }
    //Add the data retrived to the aData variable
   (*aData)[dataLoc] = data

    wg.Done()
}



//Function get the data from the URL and returns it on the aData variable
func get_url_data (url string)(*http.Response,[]byte, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Printf("Client: Error creating request: %s\n", err)
        return nil, nil, err
    }

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Printf("Client: Error making http request: %s\n", err)
        return nil, nil, err
    }

    fmt.Printf("Client: status code for %s: %d\n", url, res.StatusCode)

    defer res.Body.Close()

    resBody, err := ioutil.ReadAll(res.Body)

    if err != nil {
        fmt.Printf("Client: Error reading response body: %s\n", err)
        return nil, nil, err
    }

    return res, resBody, nil
}

//Server function that will process the request
func (s Server) gohttp(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("server: %s\n", r.Method)
    limit := r.URL.Query().Get("limit")
    sortKey := r.URL.Query().Get("sortKey")

    envLimit, isLimit := os.LookupEnv("limit")
    envCheckURLs, isURLs := os.LookupEnv("checkURLs")

    //Initializes and sets default value
    limitInt := 200
    //Checks if environment variable exists and is a number and if so, set the value
    if envTempLimit, err := strconv.Atoi(envLimit); err == nil && isLimit {
        limitInt = envTempLimit
    }

//Initializes and sets default values
    splitURLs := []string{"https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json", "https://raw.githubusercontent.com/assignment132/assignment/main/google.json", "https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json"}
//Checks if environment variable exists and if so, set the values
    if isURLs {
        splitURLs = strings.Split(envCheckURLs,",")
    }

    aAllData := make([]JsonItems, len(splitURLs))
    var wg sync.WaitGroup
    wg.Add(len(splitURLs))

    //Gets the values for each of the URLs
    for i := 0; i < len(splitURLs); i++ {
        //go get_url_data(splitURLs[i], &aAllData, i, &wg)
        go get_url_dataWRetries(splitURLs[i], &aAllData, i, &wg)
    }

    //Waits for all calls to finish before moving forward
    wg.Wait()

    //Merges all data retrieved
    var allData JsonItems
    for i := 0; i < len(splitURLs); i++ {
        allData.Data = append(allData.Data, aAllData[i].Data...)
    }

    //Sorts data by relevanceScore or view
    if sortKey == "relevanceScore" {
        sort.SliceStable(allData.Data, func(i, j int) bool {
            return allData.Data[i].RelevanceScore < allData.Data[j].RelevanceScore
	})
    } else if sortKey == "views" {
        sort.SliceStable(allData.Data, func(i, j int) bool {
            return allData.Data[i].Views < allData.Data[j].Views
	})
    }

    //Checks limits in the parameters is a number and lower than high limit or higher than 1
    if limitTemp, err := strconv.Atoi(limit); err == nil {
        if limitInt > limitTemp && limitTemp > 1 {
            limitInt = limitTemp
	}
    }

    //Limits returned results if more than items in the slice
    if len(allData.Data) > limitInt {
        allData.Data = allData.Data[:limitInt]
    }
    allData.Count = len(allData.Data)

    returnData, err := json.Marshal(allData)
    if err != nil {
        fmt.Printf("Error happened in JSON marshal. Err: %s", err)
     }

    w.Write(returnData)
    return
}
