package waifu

type WebApi interface {
	GetImages(keyword string, amount int) ([]string, error)
}
