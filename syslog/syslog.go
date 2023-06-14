package syslog

import (
	"database/sql"
	"log"
	"time"
)

type logdata struct {
	logid      int
	logTime    string
	logCommand string
}
type movieLog struct {
	ImdbID      string  `json:"imdbID"`
	Title       string  `json:"title"`
	Year        int     `json:"year"`
	Rating      float32 `json:"rating"`
	IsSuperHero bool    `json:"isSuperHero"`
	Status      string  `json :"status"`
}

func Logsave(logMovie movieLog, command string) {
	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	defer db.Close()
	currentTime := time.Now()
	log_time := currentTime.Format("2006-01-02 15:04:05")
	logScript := `INSERT INTO change_log(log_time,command) VALUES (?,?);`
	logExec, err := db.Prepare(logScript)
	if err != nil {
		log.Fatal("Error :Log Pre ERROR")
	}
	_, err = logExec.Exec(log_time, command)
	if err != nil {
		log.Fatal("Error : Log ERROR")
	}
}

/* FOR CHECK SYSLOG
func Getlog(number int) {
	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	defer db.Close()
	var Getlog logdata
	GetlogScript := db.QueryRow(`SELECT logid, log_time,command
	FROM change_log WHERE logid=?`, number)
	err = GetlogScript.Scan(&Getlog.logid, &Getlog.logTime, &Getlog.logCommand)
	fmt.Println(Getlog)

}*/
