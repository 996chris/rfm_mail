package main

type RFMBuilder interface {
	BuildDepositedRFM() (RFM[Customer], error)
	BuildNoneDepositedRFM() (RFM[Customer], error)
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

func (r *rfmBuilder) BuildDepositedRFM() (RFM[Customer], error) {
	c, err := r.repo.GetAllDepositCustomer()
	if err != nil {
		return nil, err
	}
	rfm := NewRFM(c)

	return rfm, nil

}

func (r *rfmBuilder) BuildNoneDepositedRFM() (RFM[Customer], error) {
	c, err := r.repo.GetAllNoneDepositCustomer()
	if err != nil {
		return nil, err
	}
	rfm := NewRFM(c)

	return rfm, nil

}
