package models

type Plan struct {
	Name          string  `json:"name"`
	Places        []Place `json:"places"`
	TimeInMinutes float64 `json:"time_in_minutes` // MEMO: 複数プレイスを扱うようになったら，区間ごとの移動時間も保持したい
}
