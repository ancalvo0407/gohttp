package main

//Type for URLs return data
type JsonItems struct {
    Data []struct {
        URL            string  `json:"url"`
        Views          int     `json:"views"`
        RelevanceScore float64 `json:"relevanceScore"`
    } `json:"data"`
    Count int `json:count`
}
