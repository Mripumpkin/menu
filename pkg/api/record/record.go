package record

import (
	"context"
	"fmt"
	"menu/pkg/models/menu"
	"time"

	"github.com/go-redis/redis/v8"
)

func CheckChoseCount(rdb *redis.Client) (*menu.Record, bool) {
	// ctx := context.Background()
	recordResult := menu.Record{}
	// result, err := rdb.Get(ctx, "menu_record").Result()
	// if err != nil {
	// 	return recordResult, false
	// }
	// if result == "" {
	// 	recordResult.ChooseCount = 0
	// 	recordResult.TodayMenu.MainCourse = ""
	// 	recordResult.TodayMenu.SideDish = ""
	// 	return recordResult, true
	// }
	// err = json.Unmarshal([]byte(result), &recordResult)
	// if err != nil {
	// 	return recordResult, false
	// }

	// if recordResult.ChooseCount <= 3 {
	// 	return recordResult, true
	// }
	return &recordResult, false
}

func GetTodayMenu(rdb *redis.Client) interface{} {
	ctx := context.Background()
	result, err := rdb.Get(ctx, "menu_record").Result()
	if err != nil {
		fmt.Print(result)
	}
	return nil
}

func SetTodayMenu(rdb *redis.Client, menu interface{}) {
	ctx := context.Background()
	rdb.Set(ctx, "menu_record", menu, 100*time.Hour)
}
