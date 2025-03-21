package response

type HealthCheckResponseDTO ResponseDTO[string]

type ReserveStockResponseDTO ResponseDTO[ReserveStockInfo]

type ReserveStockInfo struct {
	ReservationID string `json:"reservation_id"`
}
