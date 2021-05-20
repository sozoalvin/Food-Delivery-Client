package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func rs3(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "assets/rs3.png")
}

func userprofile(w http.ResponseWriter, req *http.Request) {

	var apiKey string

	var me = make(map[string]string) //make map for error

	session, _ := store.Get(req, "KVapi")

	if storedApiKey, ok := session.Values["api"]; !ok {

		apiKey = ""

	} else {

		apiKey = fmt.Sprintf("%v", storedApiKey)

	}

	parseData := struct {
		APIKey string
		Data   map[string]string
	}{
		apiKey, nil,
	}

	if req.Method == http.MethodPost {

		unapiKey := req.FormValue("apiKey")
		apiKey = html.EscapeString(unapiKey)

		if len(apiKey) != 32 {
			me["apiKey"] = "API Key should only contain 32 characters"
			parseData := struct {
				APIKey string
				Data   map[string]string
			}{
				"", me,
			}
			tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
			return
		}

		url := protocol + domain + "/api/v1/apivalidation?api=" + apiKey
		response, err := http.Get(url)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s \n", err)
		} else {
			fmt.Println(url)
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println("[GET] API validation from client dashboard [My API] made:", response.StatusCode)

			_ = data // fmt.Println(string(data))
			response.Body.Close()

			if response.StatusCode == 200 {

				MyHandler(w, req, apiKey)
				me["apiKey"] = "Your API Key has been validated. Thank you"
				parseData := struct {
					APIKey string
					Data   map[string]string
				}{
					apiKey, me,
				}
				tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
				return

			} else {

				me["apiKey"] = "Your API Key is wrong. Please follow the instructions below"
				parseData := struct {
					APIKey string
					Data   map[string]string
				}{
					"", me,
				}
				tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
				return
			}
		}
	}
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData)
}

func MyHandler(w http.ResponseWriter, req *http.Request, s string) {

	session, _ := store.Get(req, "KVapi")
	session.Values["api"] = s

	err := session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func additems(w http.ResponseWriter, req *http.Request) {

	var apiKey string

	var merchant merchantInformation

	session, _ := store.Get(req, "KVapi")

	if storedApiKey, ok := session.Values["api"]; !ok {
		// apiKey := ""
		http.Redirect(w, req, "/userprofile", http.StatusSeeOther)
		return

	} else {
		apiKey = fmt.Sprintf("%v", storedApiKey)
	}

	url := protocol + domain + "/api/v1/nameaddress?api=" + apiKey

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s \n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("[GET] API validation from client dashboard [Add Items] made:", response.StatusCode)
		response.Body.Close()

		if response.StatusCode == 200 {

			err = json.Unmarshal(data, &merchant)
			if err != nil {
				fmt.Println("error processing json files")
			}
		} else {
			session.Options.MaxAge = -1
			err = session.Save(req, w)
			if err != nil {
				fmt.Printf("Failed to delete session", err)
			}
		}
	}

	var me = make(map[string]string) //make map for error

	if req.Method == http.MethodPost {
		// get form values
		dnUnsanitized := req.FormValue("dishName")
		upnsanitized := req.FormValue("unitPrice")

		dnNoCaseCheck := strings.TrimSpace(html.EscapeString(dnUnsanitized))
		upoCaseCheck := html.EscapeString(upnsanitized)

		dn := strings.ToLower(dnNoCaseCheck)
		up := strings.ToLower(upoCaseCheck)

		if len(dn) > 50 || len(dn) == 0 {

			me["dishName"] = "Dishname contains invalid or too many characters"

			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "additems.gohtml", parseData)
			return
		}

		if !validateStringtoFloat(up) {

			me["unitPrice"] = "Price contains invalid or too many characters. Only 4 characters e.g. 5.80 and 5 characters e.g. 10.90 are allowed."

			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "additems.gohtml", parseData)
			return
		}

		verifyData := struct {
			Foodname string
			Price    string
		}{
			dn, up,
		}

		jsonValue, _ := json.Marshal(verifyData)
		url := protocol + domain + "/api/v1/additems?api=" + apiKey
		response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println("[POST] Request from client dashboard [Add Items] made:", response.StatusCode)
			response.Body.Close()

			if response.StatusCode == 201 {
				me["dishName"] = string(data)
				fmt.Println(string(data))
			} else if response.StatusCode == 204 {
				me["dishName"] = "This Food Item already exists. Please enter a new dish item that doesn't exist"
			}

		}

		parseData := struct {
			Merchant merchantInformation
			Me       map[string]string
		}{
			merchant, me,
		}
		tpl.ExecuteTemplate(w, "additems.gohtml", parseData)
		return
	}
	showSessions() // for demonstration purposes

	parseData := struct {
		Merchant merchantInformation
		Me       map[string]string
	}{
		merchant, nil,
	}

	tpl.ExecuteTemplate(w, "additems.gohtml", parseData)
}

func validateStringtoFloat(s string) bool {

	result := strings.Split(s, "")
	count := len(result)
	if count == 5 {
		_, err1 := strconv.Atoi(result[0])
		_, err2 := strconv.Atoi(result[1])
		_, err3 := strconv.Atoi(result[3])
		_, err4 := strconv.Atoi(result[4])

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return false
		}

	} else if count == 4 {
		_, err1 := strconv.Atoi(result[0])
		_, err2 := strconv.Atoi(result[2])
		_, err3 := strconv.Atoi(result[3])

		if err1 != nil || err2 != nil || err3 != nil {
			return false
		}
	} else {
		return false
	}
	return true
}

func validateStringtoFloat2(s string) bool {

	result := strings.Split(s, "")
	count := len(result)
	if count == 5 {
		_, err1 := strconv.Atoi(result[0])
		_, err2 := strconv.Atoi(result[1])
		_, err3 := strconv.Atoi(result[3])
		_, err4 := strconv.Atoi(result[4])

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return false
		}

	} else if count == 4 {
		_, err1 := strconv.Atoi(result[0])
		_, err2 := strconv.Atoi(result[2])
		_, err3 := strconv.Atoi(result[3])

		if err1 != nil || err2 != nil || err3 != nil {
			return false
		}

	} else if count == 0 {
		return true
	} else {
		return false
	}
	return true
}
func checkRecords(w http.ResponseWriter, req *http.Request) {

	var apiKey string

	type retrieveall struct {
		MerchantName string
		Address      string
		FoodInfo     []MerchantFoodInfo
	}

	var result retrieveall

	session, _ := store.Get(req, "KVapi")

	if storedApiKey, ok := session.Values["api"]; !ok {
		http.Redirect(w, req, "/userprofile", http.StatusSeeOther)
		return

	} else {
		apiKey = fmt.Sprintf("%v", storedApiKey)
	}

	url := protocol + domain + "/api/v1/retrieveall?api=" + apiKey

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s \n", err)
	} else {
		fmt.Println("[GET] API validation from client dashboard [Check Records] made:", response.StatusCode)
		data, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()
		fmt.Println(response.StatusCode)

		if response.StatusCode == 200 {

			err = json.Unmarshal(data, &result)
			if err != nil {
				fmt.Println("error processing json files")
			}

		} else {
			session.Options.MaxAge = -1
			err = session.Save(req, w)
			if err != nil {
				fmt.Printf("Failed to delete session", err)
			}
		}
	}

	showSessions() // for demonstration purposes

	parseData := struct {
		MerchantName string
		Address      string
		FoodInfo     []MerchantFoodInfo
	}{
		result.MerchantName, result.Address, result.FoodInfo,
	}

	tpl.ExecuteTemplate(w, "checkrecords.gohtml", parseData)
}

func editRecords(w http.ResponseWriter, req *http.Request) {

	var apiKey string

	type MerchantFoodInfo2 struct {
		FoodName string
		Price    string
		PID      int
	}
	type retrieveall struct {
		MerchantName string
		Address      string
		FoodInfo     []MerchantFoodInfo
	}

	var result retrieveall

	session, _ := store.Get(req, "KVapi")

	if storedApiKey, ok := session.Values["api"]; !ok {
		http.Redirect(w, req, "/userprofile", http.StatusSeeOther)
		return

	} else {
		apiKey = fmt.Sprintf("%v", storedApiKey)
	}

	url := protocol + domain + "/api/v1/retrieveall?api=" + apiKey

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s \n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("[GET] API validation from client dashboard [Edit Items] made:", response.StatusCode)
		response.Body.Close()

		if response.StatusCode == 200 {

			err = json.Unmarshal(data, &result)
			if err != nil {
				fmt.Println("error processing json files")
			} else {
				// fmt.Println("parsed results", result)
			}
		} else {
			session.Options.MaxAge = -1
			err = session.Save(req, w)
			if err != nil {
				fmt.Printf("Failed to delete session", err)
			}
		}
	}

	showSessions() // for demonstration purposes

	var result2 []MerchantFoodInfo2
	var count int = 1

	for _, v := range result.FoodInfo {
		result2 = append(result2, MerchantFoodInfo2{toTitle(v.FoodName), v.Price, count})
		// result2 = append(result2, MerchantFoodInfo2{v.FoodName, v.Price, count})
		count++
	}

	if req.Method == http.MethodPost {

		pid := req.FormValue("pid")
		pidint, err := strconv.Atoi(pid)
		if err != nil {
			fmt.Println("error when converting PID value from string to int")
		}

		session, _ := store.Get(req, "KVpid")
		session.Values["pid"] = result2[pidint-1].FoodName

		// Save it before we write to the response/return from the handler.
		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/edititem", http.StatusSeeOther)

	} // end req.Method

	parseData := struct {
		MerchantName string
		Address      string
		Result2      []MerchantFoodInfo2
	}{
		result.MerchantName, result.Address, result2,
	}

	tpl.ExecuteTemplate(w, "editrecords.gohtml", parseData)

}

func editItem(w http.ResponseWriter, req *http.Request) {

	var apiKey string

	var merchant merchantInformation

	session, _ := store.Get(req, "KVapi")

	if storedApiKey, ok := session.Values["api"]; !ok {
		// apiKey := ""
		http.Redirect(w, req, "/userprofile", http.StatusSeeOther)
		return

	} else {
		apiKey = fmt.Sprintf("%v", storedApiKey)
	}

	url := protocol + domain + "/api/v1/nameaddress?api=" + apiKey
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s \n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("[GET] API validation from client dashboard [Edit Item] made:", response.StatusCode)
		response.Body.Close()

		if response.StatusCode == 200 {

			err = json.Unmarshal(data, &merchant)
			if err != nil {
				fmt.Println("error processing json files")
			}
		} else {
			session.Options.MaxAge = -1
			err = session.Save(req, w)
			if err != nil {
				fmt.Printf("Failed to delete session", err)
			}
		}
	}

	session, _ = store.Get(req, "KVpid")
	dnc := fmt.Sprintf("%v", session.Values["pid"])
	// apiKey = fmt.Sprintf("%v", storedApiKey)

	var me = make(map[string]string) //make map for error

	if req.Method == http.MethodPost {
		// get form values
		dnUnsanitized := req.FormValue("dishName")
		upnsanitized := req.FormValue("unitPrice")

		dnNoCaseCheck := strings.TrimSpace(html.EscapeString(dnUnsanitized))
		upoCaseCheck := html.EscapeString(upnsanitized)

		dn := strings.ToLower(dnNoCaseCheck)
		up := strings.ToLower(upoCaseCheck)

		if len(dn) > 50 {

			me["dishName"] = "Dishname contains invalid or too many characters"

			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "additems.gohtml", parseData)
			return
		}

		if !validateStringtoFloat2(up) {

			me["unitPrice"] = "Price contains invalid or too many characters. Only 4 characters e.g. 5.80 and 5 characters e.g. 10.90 are allowed."

			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "additems.gohtml", parseData)
			return
		}

		if len(dn) == 0 {
			dn = dnc
		}

		verifyData := struct {
			OldFoodname string
			NewFoodname string
			Price       string
		}{
			dnc, strings.ToLower(dn), up,
		}

		fmt.Println(dnc, dn, up)

		fmt.Printf("\n%+v\n", verifyData)
		jsonValue, _ := json.Marshal(verifyData)
		url := protocol + domain + "/api/v1/edititems?api=" + apiKey

		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println("[Put] Request from client dashboard [Edit Item] made:", response.StatusCode)
			fmt.Println(string(data))
			response.Body.Close()

			if response.StatusCode == 200 {
				me["dishName"] = string(data)
			} else if response.StatusCode == 204 {
				me["dishName"] = "An error occured. Please try again shortly."
			}
		}

		parseData := struct {
			MerchantName    string
			MerchantAddress string
			CurrentEdit     interface{}
			Me              map[string]string
		}{
			merchant.MerchantName, merchant.Address, session.Values["pid"], me,
		}
		tpl.ExecuteTemplate(w, "edititem.gohtml", parseData)
		return
	}
	showSessions() // for demonstration purposes

	parseData := struct {
		MerchantName    string
		MerchantAddress string
		CurrentEdit     interface{}
		Me              map[string]string
	}{
		merchant.MerchantName, merchant.Address, session.Values["pid"], nil,
	}

	tpl.ExecuteTemplate(w, "edititem.gohtml", parseData)
}

func deleteItems(w http.ResponseWriter, req *http.Request) {

	var apiKey string

	type MerchantFoodInfo2 struct {
		FoodName string
		Price    string
		PID      int
	}
	type retrieveall struct {
		MerchantName string
		Address      string
		FoodInfo     []MerchantFoodInfo
	}

	var result retrieveall
	var me = make(map[string]string) //make map for error
	// var dataResponse string

	session, _ := store.Get(req, "KVapi")

	if storedApiKey, ok := session.Values["api"]; !ok {
		http.Redirect(w, req, "/userprofile", http.StatusSeeOther)
		return

	} else {
		apiKey = fmt.Sprintf("%v", storedApiKey)
	}

	url := protocol + domain + "/api/v1/retrieveall?api=" + apiKey

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s \n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("[GET] API validation from client dashboard [Delete Items] made:", response.StatusCode)
		response.Body.Close()

		if response.StatusCode == 200 {

			err = json.Unmarshal(data, &result)
			if err != nil {
				fmt.Println("error processing json files")
			}
		} else {
			session.Options.MaxAge = -1
			err = session.Save(req, w)
			if err != nil {
				fmt.Printf("Failed to delete session", err)
			}
		}
	}

	showSessions() // for demonstration purposes

	var result2 []MerchantFoodInfo2
	var count int = 1

	for _, v := range result.FoodInfo {
		result2 = append(result2, MerchantFoodInfo2{toTitle(v.FoodName), v.Price, count})
		// result2 = append(result2, MerchantFoodInfo2{(v.FoodName), v.Price, count})
		count++
	}

	if req.Method == http.MethodPost {

		pid := req.FormValue("pid")
		pidint, err := strconv.Atoi(pid)
		if err != nil {
			fmt.Println("error when converting PID value from string to int")
		}

		jsonValue, err := json.Marshal(result2[pidint-1].FoodName)

		if err != nil {
			fmt.Println("There is an error while marshaling the data into json", err)
			http.Redirect(w, req, "/deleteitems", http.StatusSeeOther)

		} else {

			url := protocol + domain + "/api/v1/deleteitems?api=" + apiKey

			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonValue))
			client := &http.Client{}
			response, err := client.Do(request)

			if err != nil {
				fmt.Printf("The hTTP request failedwith error %s", err)
			} else {
				fmt.Println("[DELETE] Request from client dashboard [Delete Items] made:", response.StatusCode)
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(response.StatusCode)
				fmt.Println(string(data))
				// dataResponse = string(data)
				response.Body.Close()

				if response.StatusCode == 200 {
					me["responseMessage"] = string(data)
				} else if response.StatusCode == 204 {
					me["responseMessage"] = string(data)
				}

			}

			parseData := struct {
				MerchantName string
				Address      string
				Result2      []MerchantFoodInfo2
				Data         map[string]string
			}{
				result.MerchantName, result.Address, result2, me,
			}

			tpl.ExecuteTemplate(w, "deleteitems.gohtml", parseData)

		}
	} // end req.Method

	parseData := struct {
		MerchantName string
		Address      string
		Result2      []MerchantFoodInfo2
		Data         map[string]string
	}{
		result.MerchantName, result.Address, result2, me,
	}

	tpl.ExecuteTemplate(w, "deleteitems.gohtml", parseData)

}

func merchantsetup(w http.ResponseWriter, req *http.Request) {

	var apiKey string

	var merchant merchantInformation

	session, _ := store.Get(req, "KVapi")

	if storedApiKey, ok := session.Values["api"]; !ok {
		// apiKey := ""
		http.Redirect(w, req, "/userprofile", http.StatusSeeOther)
		return

	} else {
		apiKey = fmt.Sprintf("%v", storedApiKey)
	}

	url := protocol + domain + "/api/v1/nameaddress?api=" + apiKey

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s \n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("[GET] API validation from client dashboard [Merchant Setup] made:", response.StatusCode)
		response.Body.Close()

		if response.StatusCode == 200 {

			err = json.Unmarshal(data, &merchant)
			if err != nil {
				fmt.Println("error processing json files")
			}
		} else {
			session.Options.MaxAge = -1
			err = session.Save(req, w)
			if err != nil {
				fmt.Printf("Failed to delete session", err)
			}
		}
	}

	var me = make(map[string]string) //make map for error

	if req.Method == http.MethodPost {
		// get form values
		r1NoCaseCheck := strings.TrimSpace(html.EscapeString(req.FormValue("merchantName")))
		r2 := html.EscapeString(req.FormValue("detailedLocation"))
		r3 := html.EscapeString(req.FormValue("postalCode"))
		r4 := html.EscapeString(req.FormValue("monWH"))
		r5 := html.EscapeString(req.FormValue("tuesWH"))
		r6 := html.EscapeString(req.FormValue("wedWH"))
		r7 := html.EscapeString(req.FormValue("thursWH"))
		r8 := html.EscapeString(req.FormValue("friWH"))
		r9 := html.EscapeString(req.FormValue("satWH"))
		r10 := html.EscapeString(req.FormValue("sunWH"))
		r11 := html.EscapeString(req.FormValue("phWH"))
		r12 := html.EscapeString(req.FormValue("cot"))

		r1 := strings.ToLower(r1NoCaseCheck)

		if len(r1) > 45 || len(r1) == 0 {
			me["merchantName"] = "Your Brand's Name is invalid. It cannot be blank, more than 20 characters or have any symbols."
			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "merchantsetup.gohtml", parseData)
			return
		}

		if len(r2) > 45 || len(r2) == 0 {
			me["detailedLocation"] = "Your location address is invalid. It cannot be blank or contain more than 45 characters"
			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "merchantsetup.gohtml", parseData)
			return
		}

		if len(r3) > 7 || len(r3) == 0 {
			me["postalCode"] = "Your postal code is invalid. Maximum Postal code should only be 6-7 characters long"
			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "merchantsetup.gohtml", parseData)
			return
		}

		var operatingHours []string
		operatingHours = append(operatingHours, r4, r5, r6, r7, r8, r9, r10, r11)

		if !checkPostalCodeInputValidity(operatingHours) {
			me["operatingHours"] = "One of your Operating Hours are invalid. Hours should be in 24 hours format. e.g. 0900-2200 (9 characters only) representing 9am to 10pm. "
			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "merchantsetup.gohtml", parseData)
			return
		}

		_, err := strconv.Atoi(r12)

		if err != nil {
			me["cot"] = "Your cut off timing is invalid. e.g. Enter 30 if your last order is 30 minutes before your closing time"
			parseData := struct {
				Merchant merchantInformation
				Me       map[string]string
			}{
				merchant, me,
			}
			tpl.ExecuteTemplate(w, "merchantsetup.gohtml", parseData)
			return
		}

		verifyData := struct {
			MerchantName     string
			DetailedLocation string
			PostalCode       string
			MonWH            string
			TuesWH           string
			WedWH            string
			ThursWH          string
			FriWH            string
			SatWH            string
			SunWH            string
			PhWH             string
			Cot              string
		}{
			r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12,
		}

		fmt.Println("this is the json value before marshalling", verifyData)
		jsonValue, _ := json.Marshal(verifyData)
		url := protocol + domain + "/api/v1/merchantsetup?api=" + apiKey
		response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println("[POST] Request from client dashboard [Merchant Setup] made:", response.StatusCode)
			response.Body.Close()

			if response.StatusCode == 201 {
				me["responseMessage"] = string(data)
				fmt.Println(string(data))
			} else if response.StatusCode == 204 {
				me["responseMessage"] = "couldn't be added"
			}

		}

		parseData := struct {
			Merchant merchantInformation
			Me       map[string]string
		}{
			merchant, me,
		}
		tpl.ExecuteTemplate(w, "merchantsetup.gohtml", parseData)
		return
	}
	showSessions() // for demonstration purposes

	parseData := struct {
		Merchant merchantInformation
		Me       map[string]string
	}{
		merchant, nil,
	}

	tpl.ExecuteTemplate(w, "merchantsetup.gohtml", parseData)
}
