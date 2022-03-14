package main

type Entry struct {
	IP string `json:"ip"`

	Serial string `json:"serial"`
}

type NS struct {
	Name string `json:"name"`

	List []Entry `json:"list"`
}

type Data struct {
	Master string `json:"master"`

	Slaves []NS `json:"slaves"`
}
