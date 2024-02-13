package main

import (
	"database/sql"
	"reflect"
)

type FactorProperty struct {
	Name string
	//true:該因子越小得分越多,false:該因子越大得分越少
	ReversWeight bool
}

// Customer的行為，各個群體資料，需實作該介面才可以使用RFM模型
type Customer interface {
	GetUserID() string
	GetUserName() string
	GetUserEmail() string
	GetFactorName() []FactorProperty
}
type EmailFormat struct {
	UserEmail string `json:"email"`
	UserName  string `json:"account"`
}
type DepositCustomer struct {
	UserID             string  `json:"user_id"`
	UserEmail          string  `json:"user_email"`
	UserName           string  `json:"user_name"`
	TotalDepositTime   float64 `json:"total_deposit_time"`
	TotalDepositAmount float64 `json:"total_deposit_amount"`
	LastTimeDepositDay float64 `json:"last_time_deposit_day"`
}
type NoneDepositCustomer struct {
	UserID        string  `json:"user_id"`
	UserEmail     string  `json:"user_email"`
	UserName      string  `json:"user_name"`
	TotalLoginDay float64 `json:"total_login_day"`
	TotalBetDay   float64 `json:"total_bet_day"`
	TotalBetTime  float64 `json:"total_bet_time"`
}
type MysqlRepository interface {
	GetAllDepositCustomer() ([]Customer, error)
	GetAllNoneDepositCustomer() ([]Customer, error)
}
type mysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(db *sql.DB) MysqlRepository {
	repo := &mysqlRepository{
		db: db,
	}

	return repo
}

func (d DepositCustomer) GetUserID() string {
	return d.UserID
}
func (d DepositCustomer) GetUserName() string {
	return d.UserName
}
func (d DepositCustomer) GetUserEmail() string {
	return d.UserEmail
}
func (d DepositCustomer) GetFactorName() []FactorProperty {
	f := make([]FactorProperty, 0)
	value := reflect.TypeOf(d)
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Field(i).Name
		if fieldName != "UserID" && fieldName != "UserName" && fieldName != "UserEmail" {

			if fieldName == "LastTimeDepositDay" {
				f = append(f, FactorProperty{Name: fieldName, ReversWeight: true})
			} else {
				f = append(f, FactorProperty{Name: fieldName, ReversWeight: false})
			}
		}
	}
	return f
}
func ToEmailFormat(cs []Customer) []EmailFormat {
	result := make([]EmailFormat, 0)

	// Populate the new slice with the desired fields
	for _, customer := range cs {
		result = append(result, EmailFormat{
			UserEmail: customer.GetUserEmail(),
			UserName:  customer.GetUserName(),
		})
	}
	return result
}

func (d NoneDepositCustomer) GetUserID() string {
	return d.UserID
}
func (d NoneDepositCustomer) GetUserName() string {
	return d.UserName
}
func (d NoneDepositCustomer) GetUserEmail() string {
	return d.UserEmail
}
func (d NoneDepositCustomer) GetFactorName() []FactorProperty {
	f := make([]FactorProperty, 0)
	value := reflect.TypeOf(d)
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Field(i).Name
		if fieldName != "UserID" && fieldName != "UserName" && fieldName != "UserEmail" {
			f = append(f, FactorProperty{Name: fieldName, ReversWeight: false})
		}
	}
	return f
}
func (m *mysqlRepository) GetAllDepositCustomer() ([]Customer, error) {

	row, err := m.db.Query(`select 
	m.id as user_id,
	m.email as user_email,
	m.username as user_name,
	COALESCE(count(dl.id),0)as total_deposit_time,
	COALESCE(sum(dl.amount),0)as total_deposit_amount,
	(UNIX_TIMESTAMP() - COALESCE(max(dl.updated_at),0)) DIV (60 * 60 * 24) as last_time_deposit_day
	from katsu.member as m
	left join katsu.deposit_log as dl
	on dl.member_id = m.id
	where m.status = 1 and m.email is not null and m.email != ''
	group by user_id
	having total_deposit_time != 0
	order by user_id;`)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	body := make([]Customer, 0)
	for row.Next() {
		var c DepositCustomer
		if err = row.Scan(&c.UserID, &c.UserEmail, &c.UserName, &c.TotalDepositTime, &c.TotalDepositAmount, &c.LastTimeDepositDay); err != nil {
			return nil, err
		}
		body = append(body, c)
	}
	return body, nil

}

func (m *mysqlRepository) GetAllNoneDepositCustomer() ([]Customer, error) {

	row, err := m.db.Query(`SELECT
			m.id as user_id,
			m.email as user_email,
			m.username as user_name,
			COALESCE(mb.login_days_count,0) as total_login_day,
			coalesce(gbr.total_bet_day,0) as total_bet_day,
			coalesce(gbr.total_bet_time,0) as total_bet_time
		FROM
			katsu.member as m
		left join
			katsu.deposit_log as dl on m.id = dl.member_id
		left join
			(SELECT
				member_id,
				COUNT(DISTINCT DATE(FROM_UNIXTIME(created_at))) AS login_days_count
			FROM
				katsu.member_login
			GROUP BY
				member_id
			) as mb on m.id = mb.member_id
		left join
			(SELECT
				member_id,
				COUNT(DISTINCT DATE(FROM_UNIXTIME(created_at))) AS total_bet_day,
				count(created_at) as total_bet_time
			FROM
				katsu.game_bet_report
			GROUP BY
				member_id
		) as gbr on m.id = gbr.member_id
		WHERE m.status = 1
		and m.email is not null
		and m.email != ''
		and FROM_UNIXTIME(m.created_at) BETWEEN '2023-12-01' AND '2023-12-31'
		GROUP BY user_id
		having COALESCE(count(dl.id),0) = 0;`)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	body := make([]Customer, 0)
	for row.Next() {
		var c NoneDepositCustomer
		if err = row.Scan(&c.UserID, &c.UserEmail, &c.UserName, &c.TotalLoginDay, &c.TotalBetDay, &c.TotalBetTime); err != nil {
			return nil, err
		}
		body = append(body, c)
	}
	return body, nil

}
