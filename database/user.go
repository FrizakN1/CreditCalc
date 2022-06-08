package database

import (
	"creditCalc/utils"
	"database/sql"
	"time"
)

type User struct {
	ID       int
	Name     string
	Surname  string
	Email    string
	Password string
}

type Session struct {
	Hash   string
	User   User
	Date   string
	Exists bool
}

var query map[string]*sql.Stmt
var sessionMap map[string]Session

func prepareUser() []string {
	sessionMap = make(map[string]Session)
	query = make(map[string]*sql.Stmt)
	errors := make([]string, 0)
	var e error

	query["SessionSelect"], e = Link.Prepare(`SELECT "Hash", "ID", "Name", "Surname", "Email", "Date" FROM "Sessions" AS s INNER JOIN "User" AS u ON u."ID" = s."User"`)
	if e != nil {
		errors = append(errors, e.Error())
	}

	query["SessionInsert"], e = Link.Prepare(`INSERT INTO "Sessions" ("Hash", "User", "Date") VALUES ($1, $2, CURRENT_TIMESTAMP)`)
	if e != nil {
		errors = append(errors, e.Error())
	}

	query["SessionDelete"], e = Link.Prepare(`DELETE FROM "Sessions" WHERE "Hash" = $1`)
	if e != nil {
		errors = append(errors, e.Error())
	}

	query["AddUser"], e = Link.Prepare(`INSERT INTO "User"("Name", "Surname", "Email", "Password") VALUES ($1,$2,$3,$4)`)
	if e != nil {
		errors = append(errors, e.Error())
	}

	query["CheckEmail"], e = Link.Prepare(`SELECT "ID" FROM "User" WHERE "Email" = $1`)
	if e != nil {
		errors = append(errors, e.Error())
	}

	query["Login"], e = Link.Prepare(`SELECT "ID", "Name", "Surname" FROM "User" WHERE "Email" = $1 AND "Password" = $2`)
	if e != nil {
		errors = append(errors, e.Error())
	}

	return errors
}

func (user *User) LoginCheck() bool {
	stmt, ok := query["Login"]
	if !ok {
		return false
	}
	row := stmt.QueryRow(user.Email, user.Password)
	e := row.Scan(&user.ID, &user.Name, &user.Surname)
	if e != nil {
		utils.Logger.Println(e)
		return false
	}
	return true
}

func (user *User) AddUser() bool {
	stmt, ok := query["AddUser"]
	if !ok {
		return false
	}

	var e error

	user.Password, e = utils.Encrypt(user.Password)
	if e != nil {
		utils.Logger.Println(e)
		return false
	}

	_, e = stmt.Exec(user.Name, user.Surname, user.Email, user.Password)
	if e != nil {
		utils.Logger.Println(e)
		return false
	}

	return true

}

func CheckEmail(email string) bool {
	stmt, ok := query["CheckEmail"]
	if !ok {
		return false
	}

	var check int
	row := stmt.QueryRow(email)
	e := row.Scan(&check)
	if e != nil {
		return true
	}

	return false
}

func GetSession(hash string) *Session {
	session, ok := sessionMap[hash]
	if ok {
		return &session
	}

	return nil
}

func (s *Session) DeleteSession() {
	stmt, ok := query["SessionDelete"]
	if !ok {
		return
	}

	_, e := stmt.Exec(s.Hash)
	if e != nil {
		utils.Logger.Println(e)
	}

	return
}

func CreateSession(user *User) (string, bool) {
	stmt, ok := query["SessionInsert"]
	if !ok {
		return "", false
	}

	hash, e := utils.GenerateHash(user.Email)
	if e != nil {
		utils.Logger.Println(e)
		return "", false
	}

	_, e = stmt.Exec(hash, user.ID)
	if e != nil {
		utils.Logger.Println(e)
		return "", false
	}

	if sessionMap != nil {
		sessionMap[hash] = Session{
			Hash: hash,
			User: User{
				ID:       user.ID,
				Name:     user.Name,
				Surname:  user.Surname,
				Email:    user.Email,
				Password: "",
			},
			Date: time.Now().String()[:19],
		}
	}

	return hash, true
}

func LoadSession(m map[string]Session) {
	stmt, ok := query["SessionSelect"]
	if !ok {
		return
	}

	rows, e := stmt.Query()
	if e != nil {
		utils.Logger.Println(e)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var session Session
		e = rows.Scan(&session.Hash, &session.User.ID, &session.User.Name, &session.User.Surname, &session.User.Email, &session.Date)
		if e != nil {
			utils.Logger.Println(e)
			return
		}

		m[session.Hash] = session //БАГ
	}
}
