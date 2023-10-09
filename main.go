package main

import (
	"encoding/json"
	"math/rand"
	"menu/config"
	"menu/pkg/api/record"
	"menu/pkg/models"
	"menu/pkg/models/menu"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func readMenu(menufile string) (*menu.Menus, error) {
	f, err := os.OpenFile(menufile, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var menu menu.Menus
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&menu); err != nil {
		return nil, err
	}
	return &menu, nil
}

func randMenu(strings []string) string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(strings))
	selectedString := strings[randomIndex]
	return selectedString
}

func menuHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	// ctx := r.Context()

	todayMenu := menu.TodayMenu{}

	// 检查菜单是否已选择
	menuRecord, isCheck := record.CheckChoseCount(rdb)
	if isCheck {
		todayMenu = menuRecord.TodayMenu
	} else {
		menufile := "menu.json"
		menus, err := readMenu(menufile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logrus.Error("无法读取菜单:", err)
			return
		}

		mainCourse := randMenu(menus.MainCourses)
		sideDish := randMenu(menus.SideDishs)

		todayMenu.MainCourse = mainCourse
		todayMenu.SideDish = sideDish

		menuRecord.ChooseCount += 1
		menuRecord.TodayMenu.MainCourse = mainCourse
		menuRecord.TodayMenu.SideDish = sideDish
		// record.SetTodayMenu(rdb, menuRecord)
	}

	// 设置HTTP响应头
	w.Header().Set("Content-Type", "application/json")

	// 编码并发送JSON响应
	if err := json.NewEncoder(w).Encode(todayMenu); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("无法编码JSON响应:", err)
		return
	}
}

func main() {
	// 加载配置
	cfg := config.LoadConfigProvider()

	// 创建Redis客户端
	rdb := models.GetRedis(cfg)

	// 提供前端文件
	http.Handle("/", http.FileServer(http.Dir("web")))

	// 设置菜单处理程序
	http.HandleFunc("/getMenu", func(w http.ResponseWriter, r *http.Request) {
		menuHandler(w, r, rdb)
	})

	// 启动服务器并监听端口
	addr := "0.0.0.0:8080"
	logrus.Infof("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatalf("Server error: %s", err)
	}
}
