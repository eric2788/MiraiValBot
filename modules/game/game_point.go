package game

import "github.com/eric2788/MiraiValBot/internal/file"

func DepositPoint(uid, p int64) {
	InitPointAccount(uid)
	file.UpdateStorage(func() {
		file.DataStorage.Points[uid] += p
	})
}

func InitPointAccount(uid int64) bool {
	_, ok := file.DataStorage.Points[uid]
	if ok {
		return false
	}
	file.UpdateStorage(func() {
		file.DataStorage.Points[uid] = 1000
	})
	return true
}

func GetPoint(uid int64) int64 {
	InitPointAccount(uid)
	return file.DataStorage.Points[uid]
}

func WithdrawPoint(uid, p int64) bool {
	InitPointAccount(uid)
	if GetPoint(uid) < p {
		return false
	}
	file.UpdateStorage(func() {
		file.DataStorage.Points[uid] -= p
	})
	return true
}

func SetPoint(uid, p int64) {
	file.UpdateStorage(func() {
		file.DataStorage.Points[uid] = p
	})
}
