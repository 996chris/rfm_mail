package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RFMHandler struct {
	builder RFMBuilder
}

type PostClassRequest struct {
	IsDeposit bool   `json:"is_deposit" validate:"boolean"`
	Class     string `json:"class" validate:"required,min=1,max=1"`
	StartDate string `json:"start_date" validate:"min=10,max=10"`
	EndDate   string `json:"end_date" validate:"min=10,max=10"`
}

func NewRFMHttpDeliver(e *gin.Engine, b RFMBuilder) {
	handler := &RFMHandler{
		builder: b,
	}
	e.POST("/rfm", handler.GetClass)
}
func (r *RFMHandler) GetClass(c *gin.Context) {
	var body PostClassRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	start, err := time.Parse(time.DateOnly, body.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	end, err := time.Parse(time.DateOnly, body.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if start.After(end) {
		c.JSON(http.StatusBadRequest, "start date can not after end date")
		return
	}
	var rfm RFM[Customer]
	if body.IsDeposit {
		rfm2, err := r.builder.BuildDepositedRFM(start, end)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		rfm = rfm2
	} else {
		rfm2, err := r.builder.BuildNoneDepositedRFM(start, end)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		rfm = rfm2
	}
	class := StringToClass(body.Class)
	data := rfm.GetClass(class)
	b := ToEmailFormat(data)
	log.Println("class : ", class, "len : ", len(b))
	c.JSON(http.StatusOK, b)
}
