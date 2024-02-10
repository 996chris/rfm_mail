package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RFMHandler struct {
	builder RFMBuilder
}

func NewRFMHttpDeliver(e *gin.Engine, b RFMBuilder) {
	handler := &RFMHandler{
		builder: b,
	}
	e.GET("/rfm", handler.GetClass)
}
func (r *RFMHandler) GetClass(c *gin.Context) {
	isDeposit := c.Query("is_deposit")
	var rfm RFM[Customer]
	if isDeposit == "1" {
		rfm2, err := r.builder.BuildDepositedRFM()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		rfm = rfm2
	} else {
		rfm2, err := r.builder.BuildNoneDepositedRFM()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		rfm = rfm2
	}
	sClass := c.Query("class")
	class := StringToClass(sClass)
	data := rfm.GetClass(class)
	b := ToEmailFormat(data)
	log.Println("class : ", class, "len : ", len(b))
	c.JSON(http.StatusOK, b)
}
