package socket

var (
	users      = map[int]*Client{}
	user2group = map[int]int{}
	group2user = map[int][]int{}
)

func IsExistUser(userId int) bool {
	// check exist
	_, isExist := users[userId]
	return isExist
}