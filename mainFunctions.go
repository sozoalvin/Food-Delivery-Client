package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func generateApiKey(s string) string {

	var timeNow string = time.Now().Format("010206")

	result := strings.Split(s, "")

	result1 := result[0:26]

	result1WithDate := append(result1, timeNow)

	userApiKey := strings.Join(result1WithDate, "")

	return userApiKey
}

func retrieveUserApiKey(un string) string {

	var apiKey string

	rows, err := db.Query(`SELECT apiKey FROM usersDB where un =?`, un)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&apiKey)
		check(err)
	}
	return apiKey
}

func showSessions() {
	fmt.Println("********")
	fmt.Println("")
}

func checkPostalCodeInputValidity(s []string) bool {

	for _, v := range s {

		if len(v) != 9 {
			return false
		}
	}

	for _, v := range s {

		OperatingHours := strings.Split(v, "")
		startingTimeSlice := OperatingHours[:4] //gives me the first 4 digits representing the operating hours
		endingTimeASlice := OperatingHours[6:]  //gives me the last 4 digits representing the closing hours

		startingTimeString := strings.Join(startingTimeSlice, "")
		endingTimeString := strings.Join(endingTimeASlice, "")

		startingTimeInt, err := strconv.Atoi(startingTimeString)
		if err != nil {
			fmt.Println("error occured when converting starting time's string to int")
			return false
		}
		endingTimeInt, err2 := strconv.Atoi(endingTimeString)

		if err2 != nil {
			fmt.Println("error occured when converting starting time's string to int")
			return false
		}

		if startingTimeInt < endingTimeInt {
			return false
		}

	}

	return true
}

func toTitle(s string) string {

	slice := strings.Split(s, "")
	var pos []int
	for i := 0; i < len(slice); i++ {

		if slice[i] == " " {
			pos = append(pos, i+1)
			i++
		}

	}
	pos = append(pos, 0)
	for i := 0; i < len(pos); i++ {

		slice[pos[i]] = strings.ToUpper(slice[pos[i]])

	}
	result := strings.Join(slice, "")

	return result
}
