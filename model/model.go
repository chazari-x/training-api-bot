package model

type Log struct {
	Level string `yaml:"level"`
}

type Discord struct {
	Token  string
	Prefix string `yaml:"prefix"`

	Channel struct {
		Log   string `yaml:"log"`
		DOLog string `yaml:"do_log"`
	} `yaml:"channel"`
}

type URLs struct {
	Training string `yaml:"training"`
}

type TrainingUser struct {
	Data struct {
		ID         int    `json:"id"`
		Login      string `json:"login"`
		Access     int    `json:"access"`
		Moder      int    `json:"moder"`
		Verify     int    `json:"verify"`
		VerifyText string `json:"verifyText"`
		Mute       int    `json:"mute"`
		Online     int    `json:"online"`
		Playerid   int    `json:"playerid"`
		Regdate    string `json:"regdate"`
		Lastlogin  string `json:"lastlogin"`
		Warn       []struct {
			Reason  string `json:"reason"`
			Admin   string `json:"admin"`
			Bantime string `json:"bantime"`
		} `json:"warn"`
	} `json:"data"`
}

type TrainingAdmin struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	LastLogin int    `json:"lastLogin"`
	Warn      int    `json:"warn"`
}
