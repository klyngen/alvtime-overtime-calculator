package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type AvailableHoursRepsonse struct {
	AvailableHoursBeforeCompensation float64 `json:"availableHoursBeforeCompensation"`
	AvailableHoursAfterCompensation  float64 `json:"availableHoursAfterCompensation"`
}

func (a *AvailableHoursRepsonse) PrintAmount(rate float64) {
	fmt.Printf("Antall timer tilgjengelig: %v \nEstimert i kroner: %v Kr,- \n", a.AvailableHoursBeforeCompensation, a.CalculateCurrency(rate))
}

func (a *AvailableHoursRepsonse) CalculateCurrency(rate float64) float64 {
	return a.AvailableHoursAfterCompensation * rate
}

func FetchAvailableHours() (*AvailableHoursRepsonse, error) {
	client := http.Client{}
	request, _ := http.NewRequest("GET", "https://api.alvtime.no/api/user/AvailableHours", nil)

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("ALVTIME_TOKEN")))

	response, err := client.Do(request)

	if err != nil {
		return nil, errors.New("Request to alvtime failed")
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.New("Unable to read the resonse")
	}

	hours := AvailableHoursRepsonse{}

	err = json.Unmarshal(data, &hours)

	if err != nil {
		return nil, errors.New("Unable to unmarshal")
	}

	return &hours, err

}

func getRateFromEnvironment() float64 {
	rate := os.Getenv("ALVTIME_RATE")
	amount, _ := strconv.ParseFloat(rate, 64)
	return amount
}

func main() {
	hours, err := FetchAvailableHours()

	if err != nil {
		log.Fatal()
	}

	hours.PrintAmount(getRateFromEnvironment())
}
