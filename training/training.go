package training

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chazari-x/training-api-bot/model"
)

type Training struct {
	cfg model.URLs
}

func NewTraining(cfg model.URLs) *Training {
	return &Training{cfg: cfg}
}

func (c *Training) GetUser(user string) (User model.TrainingUser, err error) {
	r, err := http.Get(fmt.Sprintf("%s/user/%s", c.cfg.Training, user))
	if err != nil {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	switch r.StatusCode {
	case 200:
		var data []byte
		data, err = io.ReadAll(r.Body)
		if err != nil {
			return
		}

		err = json.Unmarshal(data, &User)
	case 404:
		return
	default:
		err = fmt.Errorf("%d", r.StatusCode)
	}

	return
}

func (c *Training) GetAdmins() (Admin []model.TrainingAdmin, err error) {
	r, err := http.Get(fmt.Sprintf("%s/admin", c.cfg.Training))
	if err != nil {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	switch r.StatusCode {
	case 200:
		var data []byte
		data, err = io.ReadAll(r.Body)
		if err != nil {
			return
		}

		err = json.Unmarshal(data, &Admin)
	case 404:
		return
	default:
		err = fmt.Errorf("%d", r.StatusCode)
	}

	return
}
