package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Bank struct {
	VnName    string `json:"vn_name"`
	ShortName string `json:"shortName"`
}

func main() {
	records := struct {
		Banks []Bank `json:"banksnapas"`
	}{}

	file, err := os.Open("etc/banks.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bz, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bz, &records); err != nil {
		panic(err)
	}

	var res string

	for _, rc := range records.Banks {
		res = res + rc.VnName + " - " + rc.ShortName + "\n"
	}

	f, err := os.Create("etc/converted_banks.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(res)
	if err != nil {
		panic(err)
	}

	fmt.Println("convert successfully")
}
