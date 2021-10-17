package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type ResponseTest struct {
	Status		string
	Code string
}

type UserFull struct {
	Id 				  int
	Login   		  string
	Password 		  string
	Developer   	  int
	Ban				  int
	Group       	[]int
	Undesirable 	  int
	UserKey			  int
}

func main() {
	port := os.Getenv("PORT")
	log.Println("server start with port: " +port)
	http.HandleFunc("/api", api)
	http.ListenAndServe(":" +port, nil)
}


func api(w http.ResponseWriter, r *http.Request) {
	 ErrUrl := [...]string{
		"Error 101:\nBad title ( Empty )",
		"Error 102:\nBad token",
	}

	log.Print("new request")
	query := r.URL.Query()

	TitleByte := query.Get("title")
	TokenByte := query.Get("token")
	EmailByte := query.Get("email")
	//UserKeyByte := query.Get("key")
	PasswordByte := query.Get("password")
	LoginByte := query.Get("login")


	title := string(TitleByte[:])		//err id 1
	token := string(TokenByte[:])		//err id 2
	email := string(EmailByte[:])		//err id 3
	//userkey := string(UserKeyByte[:])	//err id 4
	password := string(PasswordByte[:])	//err id 5
	login := string(LoginByte[:])		//err id 6

	if title == "" {
		log.Println(ErrUrl[0])
		fmt.Fprintf(w, ErrUrl[0])
	}

	if token == "" {
		log.Println(ErrUrl[1])
		fmt.Fprintf(w, ErrUrl[1])
	}


	if title == "reg" {
		if token != "cardinal" {
			responsetest := ResponseTest{"ERR",  "291 Bad token cardinal"}

			js, err := json.Marshal(responsetest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			log.Println(js)
		}else{

			if login == "" {
				responsetest := ResponseTest{"ERR",  "106 Login empty"}

				js, err := json.Marshal(responsetest)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
				log.Println(js)
			}else{
				if email == "" {
					responsetest := ResponseTest{"ERR",  "103 Email empty"}

					js, err := json.Marshal(responsetest)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
					log.Println(js)
				}else{
					if password == "" {
						responsetest := ResponseTest{"ERR",  "105  Password empty"}

						js, err := json.Marshal(responsetest)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						w.Header().Set("Content-Type", "application/json")
						w.Write(js)
						log.Println(js)
					}else{
						//db connect then reg user and send response with code
						db, err := sql.Open("mysql", "site:xLb43XEDkr8R4O@tcp(185.219.40.250)/site")


						if err != nil {
							panic(err.Error())
						}


						defer db.Close()

						query := fmt.Sprintf("SELECT * FROM `users` WHERE `login` = ?")
						rows, err := db.Query(query, login)
						if err != nil {
							if err == sql.ErrNoRows {
								query_email := fmt.Sprintf("SELECT * FROM `users` WHERE `email` = ?")
								rows_email, err := db.Query(query_email, email)

								if err != nil {
									if err == sql.ErrNoRows {
										db, err := sql.Open("mysql", "site:xLb43XEDkr8R4O@tcp(185.219.40.250)/site")


										if err != nil {
											panic(err.Error())
										}

										md5_userkey := rand.Intn(9999999999)

										md5_password := GetMD5Hash(password)
										var null = "0"

										res, err := db.Exec("INSERT INTO `users` (`login`, `password`, `userkey`, `ban`, `group`, `developer`, `undesirable`) VALUES(?, ?, ?, ?, ?, ?, ?)", login, md5_password, md5_userkey, null, null, null, null)
										if err != nil {

										}
										log.Println("Create user login: " + login + " | password: "+ password)
										spew.Dump(res)

										responsetest := ResponseTest{"OK",  "10 Create user"}

										js, err := json.Marshal(responsetest)
										if err != nil {
											http.Error(w, err.Error(), http.StatusInternalServerError)
											return
										}

										w.Header().Set("Content-Type", "application/json")
										w.Write(js)
										log.Println(js)
									}else{
										responsetest := ResponseTest{"ERR",  "999 Unknown error ( бтв такой пользователь уже есть )"}

										js, err := json.Marshal(responsetest)
										if err != nil {
											http.Error(w, err.Error(), http.StatusInternalServerError)
											return
										}

										w.Header().Set("Content-Type", "application/json")
										w.Write(js)
										log.Println(js)
									}
								}
								defer rows_email.Close()

							}else{
								responsetest := ResponseTest{"ERR",  "999 Unknown error ( бтв такой пользователь уже есть )"}

								js, err := json.Marshal(responsetest)
								if err != nil {
									http.Error(w, err.Error(), http.StatusInternalServerError)
									return
								}

								w.Header().Set("Content-Type", "application/json")
								w.Write(js)
								log.Println(js)
							}
						}

						defer rows.Close()

					}
				}
			}

		}
	}

	if title == "auth" {
			if token != "cardinal" {

				responsetest := ResponseTest{"ERR",  "201 Incorrect Token Cardinal"}

				js, err := json.Marshal(responsetest)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
				log.Println(js)
		}else{
			if login == "" {
				responsetest := ResponseTest{"ERR",  "??? Empty url login"} //change err code

				js, err := json.Marshal(responsetest)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
				log.Println(js)
			}else{
				if password == ""{
					responsetest := ResponseTest{"ERR",  "??? Empty url password"} //change err code

					js, err := json.Marshal(responsetest)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					w.Header().Set("Content-Type", "application/json")
					w.Write(js)
					log.Println(js)
				}else{
					// db connect, find user, if *status* - response
					db, err := sql.Open("mysql", "site:xLb43XEDkr8R4O@tcp(185.219.40.250)/site")


					if err != nil {
						panic(err.Error())
					}


					defer db.Close()

					query := fmt.Sprintf("SELECT * FROM `users` WHERE `login` = ?")
					rows, err := db.Query(query, login)
					if err != nil {

						responsetest := ResponseTest{"ERR",  "304 User not found"}

						js, err := json.Marshal(responsetest)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						w.Header().Set("Content-Type", "application/json")
						w.Write(js)
						log.Println(js)
					}
					defer rows.Close()

					for rows.Next() {
							var user UserFull
							err = rows.Scan(user.Id, user.Login, user.Password, user.Developer, user.Group, user.Ban, user.Undesirable)
							if err != nil {
								log.Println(err)
							}



							if password == user.Password {
								responsetest := ResponseTest{"OK",  "1 TRUE"}

								js, err := json.Marshal(responsetest)
								if err != nil {
									http.Error(w, err.Error(), http.StatusInternalServerError)
									return
								}

								w.Header().Set("Content-Type", "application/json")
								w.Write(js)
								log.Println(js)
							}else{
								responsetest := ResponseTest{"ERR",  "302 Bad password"}

								js, err := json.Marshal(responsetest)
								if err != nil {
									http.Error(w, err.Error(), http.StatusInternalServerError)
									return
								}

								w.Header().Set("Content-Type", "application/json")
								w.Write(js)
								log.Println(js)
							}
					}
				}
			}

		}
	}

}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}