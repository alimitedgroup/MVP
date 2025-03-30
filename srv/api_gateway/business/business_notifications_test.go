package business

//go:generate go run go.uber.org/mock/mockgen@latest -destination mock_notifications.go -package business github.com/alimitedgroup/MVP/srv/api_gateway/portout NotificationPortOut
