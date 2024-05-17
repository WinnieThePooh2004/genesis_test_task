package rates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test_project/settings"
)

type NbuRatesService struct {
	settings settings.AppSettings
}

func NewService(settings settings.AppSettings) *NbuRatesService {
	return &NbuRatesService{settings: settings}
}

func (s *NbuRatesService) GetRate() (float64, error) {
	response, err := http.Get(s.settings.RatesUrl)
	if err != nil {
		return 0, err
	}

	if response.StatusCode != 200 {
		return 0, fmt.Errorf("nbu API is not currently available, status code: %d", response.StatusCode)
	}

	rate := [1]NbuRate{}
	err = json.NewDecoder(response.Body).Decode(&rate)
	if err != nil {
		return 0, err
	}

	return rate[0].Rate, nil
}
