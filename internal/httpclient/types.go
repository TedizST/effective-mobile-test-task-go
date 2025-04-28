package httpclient

type (
	APIType   string
	HTTPError struct {
		StatusCode int
		Body       string
	}
	AgifyResponse struct {
		Count uint   `json:"count"`
		Name  string `json:"name"`
		Age   uint64 `json:"age"`
	}
	GenderizeResponse struct {
		Count       uint    `json:"count"`
		Name        string  `json:"name"`
		Gender      string  `json:"gender"`
		Probability float64 `json:"probability"`
	}
	CountryProbability struct {
		CountryId   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	}
	NationalizeResponse struct {
		Count     uint                 `json:"count"`
		Name      string               `json:"name"`
		Countries []CountryProbability `json:"country"`
	}
	PredictorResponse interface {
		AgifyResponse | GenderizeResponse | NationalizeResponse
	}
)

const (
	Agify       APIType = "agify"
	Genderize   APIType = "genderize"
	Nationalize APIType = "nationalize"
)
