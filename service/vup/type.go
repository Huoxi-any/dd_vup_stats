package vup

import (
	"time"
)

type (
	ListResp struct {
		Page    int         `json:"page"`
		Size    int         `json:"size"`
		MaxPage int64       `json:"max_page"`
		Total   int64       `json:"total"`
		List    []*UserResp `json:"list"`
	}

	UserResp struct {
		UserInfo
		Listening      bool      `json:"listening"`
		LastListenedAt time.Time `json:"last_listened_at"`
	}

	UserInfo struct {
		SimpleUserInfo
		FirstListenAt   time.Time  `json:"first_listen_at"`
		LastBehaviourAt *time.Time `json:"last_behaviour_at"`
		DDCount         int64      `json:"dd_count"`
	}

	SimpleUserInfo struct {
		Uid    int64  `json:"uid"`
		Name   string `json:"name"`
		Face   string `json:"face"`
		RoomId int64  `json:"room_id"`
		Sign   string `json:"sign"`
	}

	Analysis struct {
		TopDDVups    []*SimpleUserInfo `json:"top_dd_vups"`
		TopGuestVups []*SimpleUserInfo `json:"top_guest_vups"`
	}
)
