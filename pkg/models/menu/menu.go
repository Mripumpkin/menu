package menu

type Menus struct {
	MainCourses []string `json:"main_course"`
	SideDishs   []string `json:"side_dish"`
}

type TodayMenu struct {
	MainCourse string `json:"main_course"`
	SideDish   string `json:"side_dish"`
}

type Record struct {
	ChooseCount int       `json:"choose_count"`
	TodayMenu   TodayMenu `json:"menu"`
}
