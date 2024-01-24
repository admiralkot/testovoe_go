package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	//"log"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age,omitempty"`
	//мы живём в современном обществе, возможно у нас не 2 гендера :D
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type notAllowedHandler struct{}

// SliceToJSON encodes a slice with JSON records
func SliceToJSON(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(slice)
}

func (h notAllowedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	MethodNotAllowedHandler(rw, r)
}

func DefaultHandler(rw http.ResponseWriter, r *http.Request) {
	log.Debug("DefaultHandler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	rw.WriteHeader(http.StatusNotFound)
	Body := r.URL.Path + " is not supported. Thanks for visiting!\n"
	fmt.Fprintf(rw, "%s", Body)
}

// MethodNotAllowedHandler is executed when the HTTP method is incorrect
func MethodNotAllowedHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	rw.WriteHeader(http.StatusNotFound)
	Body := "Method not allowed!\n"
	fmt.Fprintf(rw, "%s", Body)
}

// TimeHandler is for handling /time – it works with plain text
func TimeHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("TimeHandler Serving:", r.URL.Path, "from", r.Host)
	rw.WriteHeader(http.StatusOK)
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is: " + t + "\n"
	fmt.Fprintf(rw, "%s", Body)
}

// AddHandler is for adding a new user
func AddHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("AddHandler Serving:", r.URL.Path, "from", r.Host)
	d, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln(err)
		return
	}

	if len(d) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln("No input!")
		return
	}

	var users = User{}
	err = json.Unmarshal(d, &users)
	if err != nil {
		log.Errorln(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	var (
		urlage         = "https://api.agify.io/?name=" + users.Name
		urlgander      = "https://api.genderize.io/?name=" + users.Name
		urlnationality = "https://api.nationalize.io/?name=" + users.Name
	)
	log.Debugln("urlage=", urlage, "\n\nurlgander=", urlgander, "\n\nurlnationality=", urlnationality)

	respAge, err := http.Get(urlage)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln(err)
		return
	}
	//respGander, err := http.Get(urlgander)
	//if err != nil {
	//	rw.WriteHeader(http.StatusBadRequest)
	//	log.Errorln(err)
	//	return
	//}
	//respNationality, err := http.Get(urlnationality)
	//if err != nil {
	//	rw.WriteHeader(http.StatusBadRequest)
	//	log.Errorln(err)
	//	return
	//}

	log.Debugln("resp=", respAge)

	age, err := io.ReadAll(respAge.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln(err)
		return
	}
	log.Debugln(age)
	err = json.Unmarshal(age, &users)
	if err != nil {
		log.Errorln(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Debug(users)

	//result := InsertUser(users)
	//if !result {
	//	rw.WriteHeader(http.StatusBadRequest)
	//}
}

// DeleteHandler is for deleting an existing user + DELETE
func DeleteHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("DeleteHandler Serving:", r.URL.Path, "from", r.Host)

	// Get the ID of the user to be deleted
	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Errorln("ID value not set!")
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	var user = User{}
	err := user.FromJSON(r.Body)
	if err != nil {
		log.Errorln(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Errorln("id", err)
		return
	}

	t := FindUserID(intID)
	if t.Name != "" {
		log.Infoln("About to delete:", t)
		deleted := DeleteUser(intID)
		if deleted {
			log.Infoln("User deleted:", id)
			rw.WriteHeader(http.StatusOK)
			return
		} else {
			log.Errorln("User ID not found:", id)
			rw.WriteHeader(http.StatusNotFound)
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

// GetAllHandler is for getting all data from the user database
func GetAllHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("GetAllHandler Serving:", r.URL.Path, "from", r.Host)
	d, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln(err)
		return
	}

	if len(d) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln("No input!")
		return
	}

	var user = User{}
	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Errorln(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = SliceToJSON(ListAllUsers(), rw)
	if err != nil {
		log.Errorln(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetIDHandler returns the ID of an existing user
func GetIDHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("GetIDHandler Serving:", r.URL.Path, "from", r.Host)
	//
	//username, ok := mux.Vars(r)["username"]
	//if !ok {
	//	log.Println("ID value not set!")
	//	rw.WriteHeader(http.StatusNotFound)
	//	return
	//}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln(err)
		return
	}

	if len(d) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln("No input!")
		return
	}

	var user = User{}
	err = json.Unmarshal(d, &user)
	if err != nil {
		log.Errorln(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	t := FindUserName(user.Name)
	if t.ID != 0 {
		err := t.ToJSON(rw)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Errorln(err)
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
		log.Errorln("User " + user.Name + "not found")
	}
}

// GetUserDataHandler + GET returns the full record of a user
func GetUserDataHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("GetUserDataHandler Serving:", r.URL.Path, "from", r.Host)
	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Errorln("ID value not set!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Errorln("id", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	t := FindUserID(intID)
	if t.ID != 0 {
		err := t.ToJSON(rw)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Errorln(err)
		}
		return
	}

	log.Errorln("User not found:", id)
	rw.WriteHeader(http.StatusBadRequest)
}

// UpdateHandler is for updating the data of an existing user + PUT
func UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	log.Infoln("UpdateHandler Serving:", r.URL.Path, "from", r.Host)
	d, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln(err)
		return
	}

	if len(d) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		log.Errorln("No input!")
		return
	}

	var users = User{}
	err = json.Unmarshal(d, &users)
	if err != nil {
		log.Errorln(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debugln(users)
	t := FindUserName(users.Name)
	users.ID = t.ID
	if !UpdateUser(users) {
		log.Errorln("Update failed:", users)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infoln("Update successful:", users)
	rw.WriteHeader(http.StatusOK)
}
