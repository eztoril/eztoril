package database

import (
	"database/sql"
	"fmt"
	"time"
)

func writeDB(body cumInvRTDataReqResp) int {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "solar:solar@tcp(127.0.0.1:3306)/solardb")
	defer db.Close()

	if err != nil {
		fmt.Println("sql.Open failure: ", err)
		return -1
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("db.Ping failure: ", err)
		return -1
	}
	currentTime := time.Now()
	fmt.Printf("Total power production at %s was: %d\n",
		currentTime.Format("2006-01-02 3:4"), body.Body.Data.DayEnergy.Values.Num1)
	query := fmt.Sprintf("INSERT INTO Production(Day,Power) values(curdate(), %d);",
		body.Body.Data.DayEnergy.Values.Num1)

	_, err = db.Query(query)

	if err != nil {
		fmt.Println("db.Query failure: ", err)
		return -1
	}
	return 0
}
