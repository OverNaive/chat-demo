package main

type Message struct {
	To int64 `json:"to"`
	From int64 `json:"from"`
	Content string `json:"content"`
}
