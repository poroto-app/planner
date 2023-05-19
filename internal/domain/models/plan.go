package models

type Plan struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Places        []Place `json:"places"`
	TimeInMinutes uint    `json:"time_in_minutes"` // MEMO: 複数プレイスを扱うようになったら，区間ごとの移動時間も保持したい
}
