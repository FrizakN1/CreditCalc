package database

import (
	"creditCalc/utils"
	"database/sql"
	"fmt"
	"math"
)

type Credit struct {
	CreditSum      int64
	CreditDuration float64
	CreditRate     float64
	CreditPayment  int64
	DateApply      string
}

func prepareCredit() []string {
	errors := make([]string, 0)
	if query == nil {
		query = make(map[string]*sql.Stmt)
	}
	var e error

	query["AddCredit"], e = Link.Prepare(`INSERT INTO "Credit"("User", "Date_apply", "Credit_amount", "Credit_duration", "Credit_rate", "Credit_payment") VALUES ($1, CURRENT_TIMESTAMP, $2, $3, $4, $5)`)
	if e != nil {
		errors = append(errors, e.Error())
	}

	query["GetCredits"], e = Link.Prepare(`SELECT "Date_apply", "Credit_amount", "Credit_duration", "Credit_rate", "Credit_payment" FROM "Credit" WHERE "User" = $1`)

	return errors
}

func (s *Session) GetCredits() []Credit {
	var credits []Credit

	stmt, ok := query["GetCredits"]
	if !ok {
		return nil
	}

	rows, e := stmt.Query(s.User.ID)
	if e != nil {
		utils.Logger.Println(e)
		return nil
	}

	for rows.Next() {
		var credit Credit
		e = rows.Scan(&credit.DateApply, &credit.CreditSum, &credit.CreditDuration, &credit.CreditRate, &credit.CreditPayment)
		if e != nil {
			utils.Logger.Println(e)
			return nil
		}
		credit.DateApply = credit.DateApply[0:10]
		fmt.Println(credit.CreditSum)
		credits = append(credits, credit)
	}

	return credits
}

func (s *Session) ApplyCredit(creditData Credit) bool {
	stmt, ok := query["AddCredit"]
	if !ok {
		return false
	}

	payment := math.Round(float64(creditData.CreditSum) * (0.159 / 12) / (1 - 1/math.Pow(1+0.159/12, creditData.CreditDuration)))

	_, e := stmt.Exec(s.User.ID, creditData.CreditSum, creditData.CreditDuration, 15.9, payment)
	if e != nil {
		utils.Logger.Println(e)
		return false
	}

	return true
}
