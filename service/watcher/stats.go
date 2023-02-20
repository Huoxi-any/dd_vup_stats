package watcher

import (
	"fmt"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"

	"gorm.io/gorm"
)

func GetTotalCommandStats(uid int64) ([]stats.TotalStats, error) {

	var totalStatus []stats.TotalStats

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"command",
			"COUNT(*) as count",
			"SUM(price) as price",
		}).
		Where("uid = ?", uid).
		Group("command").
		Find(&totalStatus).
		Error

	if err != nil {
		logger.Errorf("獲取 %v 的行為時出現錯誤: %v", uid, err)
		return nil, err
	}

	return totalStatus, nil
}

func GetStatsCommand(uid int64, limit int, command string, price bool) ([]AnalysisVupInfo, error) {

	var mostDDVups []AnalysisVupInfo

	r := db.Database.Model(&db.WatcherBehaviour{})

	orderBy := "count"

	if price {
		orderBy = "price"
		r = r.Where("watcher_behaviours.price > 0")
	}

	r = r.Session(&gorm.Session{})

	// D 最多
	err := r.
		Where("watcher_behaviours.uid = ? AND watcher_behaviours.command = ?", uid, command).
		Joins("inner join vups on vups.uid = watcher_behaviours.target_uid").
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.face",
			"COUNT(*) as count",
			"SUM(watcher_behaviours.price) as price",
		}).
		Group("watcher_behaviours.target_uid, vups.uid").
		Order(fmt.Sprintf("%s desc", orderBy)).
		Limit(limit).
		Find(&mostDDVups).
		Error

	if err != nil {
		return nil, err
	}

	return mostDDVups, nil
}

func GetStats(uid int64, limit int) (*Analysis, error) {

	var mostDDVups []AnalysisVupInfo

	// D 最多
	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Joins("inner join vups on vups.uid = watcher_behaviours.target_uid").
		Where("watcher_behaviours.uid = ?", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.face",
			"COUNT(*) as count",
			"SUM(watcher_behaviours.price) as price",
		}).
		Group("watcher_behaviours.target_uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDVups).
		Error

	if err != nil {
		return nil, err
	}

	var mostSpentVups []PricedVupInfo

	// 花費最多
	err = db.Database.
		Model(&db.WatcherBehaviour{}).
		Joins("inner join vups on vups.uid = watcher_behaviours.target_uid").
		Where("watcher_behaviours.uid = ? and watcher_behaviours.price > 0", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.face",
			"SUM(watcher_behaviours.price) as spent",
		}).
		Group("watcher_behaviours.target_uid, vups.uid").
		Order("spent desc").
		Limit(limit).
		Find(&mostSpentVups).
		Error

	if err != nil {
		return nil, err
	}

	return &Analysis{
		TopDDVups:    mostDDVups,
		TopSpentVups: mostSpentVups,
	}, nil
}


// GetMostBehavioursByVup 返回该 vup 中最高互动的 watcher
func GetMostBehavioursByVup(uid int64, limit int) ([]AnalysisWatcherInfo, error) {
	
	var mostDDWatchers []AnalysisWatcherInfo

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			"(array_agg(u_name order by created_at desc))[1] as u_name",
			"COUNT(*) as count",
		}).
		Where("target_uid = ?", uid).
		Group("uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDWatchers).
		Error

	if err != nil {
		return nil, err
	}

	return mostDDWatchers, nil
}

// GetMostSpentByVup 返回该 vup 中花费最多的 watcher
func GetMostSpentByVup(uid int64, limit int) ([]AnalysisWatcherInfo, error) {
	
	var mostSpentWatchers []AnalysisWatcherInfo

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			"(array_agg(u_name order by created_at desc))[1] as u_name",
			"SUM(price) as price",
		}).
		Where("target_uid = ? and price > 0", uid).
		Group("uid").
		Order("price desc").
		Limit(limit).
		Find(&mostSpentWatchers).
		Error

	if err != nil {
		return nil, err
	}

	return mostSpentWatchers, nil
}