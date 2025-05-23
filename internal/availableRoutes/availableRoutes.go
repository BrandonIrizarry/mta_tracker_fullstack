package availableRoutes

type AvailableRoutes struct {
	Code        int   `json:"code"`
	CurrentTime int64 `json:"currentTime"`
	Data        struct {
		LimitExceeded bool `json:"limitExceeded"`
		List          []struct {
			AgencyID    string `json:"agencyId"`
			Color       string `json:"color"`
			Description string `json:"description"`
			ID          string `json:"id"`
			LongName    string `json:"longName"`
			ShortName   string `json:"shortName"`
			TextColor   string `json:"textColor"`
			Type        int    `json:"type"`
			URL         string `json:"url"`
		} `json:"list"`
		References struct {
			Agencies []struct {
				Disclaimer     string `json:"disclaimer"`
				ID             string `json:"id"`
				Lang           string `json:"lang"`
				Name           string `json:"name"`
				Phone          string `json:"phone"`
				PrivateService bool   `json:"privateService"`
				Timezone       string `json:"timezone"`
				URL            string `json:"url"`
			} `json:"agencies"`
			Routes     []any `json:"routes"`
			Situations []any `json:"situations"`
			Stops      []any `json:"stops"`
			Trips      []any `json:"trips"`
		} `json:"references"`
	} `json:"data"`
	Text    string `json:"text"`
	Version int    `json:"version"`
}
