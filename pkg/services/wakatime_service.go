package services

import (
	"ekel-backend/pkg/entities"
	"ekel-backend/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WakatimeService struct {
}

func NewWakatimeService() *WakatimeService {
	return &WakatimeService{}
}

func (s *WakatimeService) GetWakatimeStats() (*entities.WakaTimeResponse, error) {
	var wakatimeResponse *entities.WakaTimeResponse
	fmt.Println(utils.Env().WAKATIME_API_URL)
	req, err := http.NewRequest("GET", utils.Env().WAKATIME_API_URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+utils.Env().WAKATIME_API_KEY)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &wakatimeResponse)
	if err != nil {
		return nil, err
	}

	return wakatimeResponse, nil
}
