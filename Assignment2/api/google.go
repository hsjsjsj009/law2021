package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GoogleAPI struct {
	Token string
	Client *http.Client
}

func NewGoogleAPI(token string) (*GoogleAPI,error) {
	if token == "" {
		return nil,fmt.Errorf("token empty")
	}
	return &GoogleAPI{Token: token,Client: http.DefaultClient},nil
}

type GoogleData struct {
	FamilyName string `json:"family_name"`
	Name string `json:"name"`
	Picture	string	`json:"picture"`
	Locale string `json:"locale"`
	Email string `json:"email"`
	GivenName string `json:"given_name"`
	ID string `json:"id"`
	VerifiedEmail bool `json:"verified_email"`
}

func (g *GoogleAPI) GetUserData() (*GoogleData,error) {
	req,err := http.NewRequest(http.MethodGet,"https://www.googleapis.com/userinfo/v2/me",nil)
	if err != nil {
		return nil,err
	}
	req.Header.Add("Authorization",fmt.Sprintf("Bearer %s",g.Token))
	res,err := g.Client.Do(req)
	if err != nil {
		return nil,err
	}
	if res.StatusCode != 200 {
		return nil,fmt.Errorf("token invalid")
	}
	data := &GoogleData{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(data)
	if err != nil {
		return nil,err
	}
	return data,nil
}
