package controller

import "github.com/gin-gonic/gin"

type ReservationController struct {
}

func NewReservationController() *ReservationController {
	return &ReservationController{}
}

func (c *ReservationController) CreateReservationHandler(ctx *gin.Context) {
}
