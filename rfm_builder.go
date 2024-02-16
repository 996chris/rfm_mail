package main

import "time"

type RFMBuilder interface {
	BuildDepositedRFM(start time.Time, end time.Time) (RFM[Customer], error)
	BuildNoneDepositedRFM(start time.Time, end time.Time) (RFM[Customer], error)
}
type rfmBuilder struct {
	repo MysqlRepository
}

func NewRfmBuilder(repo MysqlRepository) RFMBuilder {
	r := &rfmBuilder{
		repo: repo,
	}

	return r
}

func (r *rfmBuilder) BuildDepositedRFM(start time.Time, end time.Time) (RFM[Customer], error) {
	c, err := r.repo.GetDepositCustomer(start, end)
	if err != nil {
		return nil, err
	}
	rfm := NewRFM(c)

	return rfm, nil

}

func (r *rfmBuilder) BuildNoneDepositedRFM(start time.Time, end time.Time) (RFM[Customer], error) {
	c, err := r.repo.GetNoneDepositCustomer(start, end)
	if err != nil {
		return nil, err
	}
	rfm := NewRFM(c)

	return rfm, nil

}
