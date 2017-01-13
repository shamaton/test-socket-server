package socket

var (
	users      = map[int]*Client{}
	user2name  = map[int]string{}
	user2group = map[int]int{}
	group2user = map[int]map[int]bool{}
)

func IsExistUser(userId int) bool {
	// check exist
	_, isExist := users[userId]
	return isExist
}
