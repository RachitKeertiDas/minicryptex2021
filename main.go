package main

import (
	"encoding/json"
	"strconv"
	"time"

	// "github.com/graphql-go/handler"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "github.com/codegangsta/negroni"
	// "github.com/auth0/go-jwt-middleware"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/joho/godotenv"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type user struct {
	Username string `json:"username"`
	Level    int    `json:"level"`
}

type DatabaseUserObject struct {
	ClientID string `json:"clientID"`
	Username string `json:"username"`
	Level    int    `json:"level"`
	Name1    string `json:"name1"`
	Name2    string `json:"name2"`
	Name3    string `json:"name3"`
	Name4    string `json:"name4"`
	Name5    string `json:"name5"`
}

type LevelResponse struct {
	Level int
	URL   string
}

var answers map[string]string

// var context.TODO(), _ = context.WithTimeout(context.Background(), 10*time.Second)
// var client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
// var collection = client.Database("cryptex").Collection("users")

var collection *mongo.Collection

func main() {
	fmt.Println("Server started... ")
	fmt.Println("To do : Protect all endpoints with JWT Auth")
	fmt.Println("Change level type to int. It's string rn. ")
	answers = make(map[string]string)
	answers["0"] = "zozimus"
	answers["1"] = "goldenfleece"
	answers["2"] = "221"
	answers["3"] = "icveg"
	answers["4"] = "hailhydra"
	answers["5"] = "riptide"
	answers["6"] = "pixar"
	answers["7"] = "lemonysnicket"
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection = client.Database("Cryptex").Collection("users")
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := "https://dev-l0ini8h1.us.auth0.com/api/v2/"
			fmt.Println("here")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := "https://dev-l0ini8h1.us.auth0.com/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	// jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options {
	//     ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	//         // Verify 'aud' claim
	//         aud := "https://cryptex2020.auth0.com/api/v2/"
	//         checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	//         if !checkAud {
	//             return token, errors.New("Invalid audience.")
	//         }
	//         // Verify 'iss' claim
	//         iss := "https://cryptex2020.auth0.com/"
	//         checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	//         if !checkIss {
	//             return token, errors.New("Invalid issuer.")
	//         }

	//         cert, err := getPemCert(token)
	//         if err != nil {
	//             panic(err.Error())
	//         }

	//         result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	//         return result, nil
	//     },
	//     SigningMethod: jwt.SigningMethodRS256,
	// })
	router := mux.NewRouter()

	// router.Handle("/callback", http.ServeFile())
	// Without JWT middleware check
	// router.Handle("/things", ThingsHandler).Methods("GET")

	// ALL API CALLS (GraphQL) are defined here
	router.HandleFunc("/whichlevel/{clientid}", LevelQueryHandler)
	router.HandleFunc("/css", CSSHandler)
	router.HandleFunc("/doesUsernameExist/{username}", DoesUsernameExistHandler)
	router.Handle("/adduser/{ID}/{username}/{name1}/{name2}/{name3}/{name4}/{name5}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Here")
			vars := mux.Vars(r)
			find, _ := collection.Find(context.TODO(), bson.M{"email": vars["ID"]})
			JSOND, _ := json.Marshal(find.Next(context.TODO()))
			UserStatus := string(JSOND)
			if strings.Compare(UserStatus, "false") == 0 {
				_, _ = collection.InsertOne(context.TODO(), bson.M{"clientID": vars["ID"], "username": vars["username"], "level": -1, "name1": vars["name1"], "name2": vars["name2"], "name3": vars["name3"], "name4": vars["name4"], "name5": vars["name5"], "lastModified": time.Now().UTC()})
			}
		}))))
	router.Handle("/acceptedrules", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.URL.Query()["id_token"][0]
			claims := jwt.MapClaims{}
			jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				x, err := getPemCert(token)
				return []byte(x), err
			})
			email := claims["email"].(string)
			filter := bson.D{{"clientID", email}}
			update := bson.D{
				{"$set", bson.D{
					{"level", 0},
				}},
			}
			_, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Fatal(err)
			}
		}))))
	router.Handle("/answer/{level}/{answer}", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.URL.Query()["id_token"][0]
			claims := jwt.MapClaims{}
			jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				x, err := getPemCert(token)
				return []byte(x), err
			})
			email := claims["email"].(string)
			find, _ := collection.Find(context.TODO(), bson.M{"clientID": email})
			JSOND, _ := json.Marshal(find.Next(context.TODO()))
			vars := mux.Vars(r)
			if strings.Compare(string(JSOND), "true") == 0 {
				if val, ok := answers[vars["level"]]; ok {
					var current DatabaseUserObject
					err := find.Decode(&current)
					if err != nil {
						fmt.Println("Error decoding database object ", err)
					}
					fmt.Println(current.Username, " ", current.Level, " ", vars["answer"])
					if strings.Compare(strconv.Itoa(current.Level), vars["level"]) == 0 {
						if strings.Compare(val, vars["answer"]) == 0 {
							filter := bson.D{{"clientID", email}}
							update := bson.D{
								{"$inc", bson.D{
									{"level", 1},
								}},
							}
							_, err := collection.UpdateOne(context.TODO(), filter, update)
							update = bson.D{
								{"$set", bson.D{
									{"lastModified", time.Now().UTC()},
								}},
							}
							_, err = collection.UpdateOne(context.TODO(), filter, update)
							if err != nil {
								fmt.Println("Error updating ", err)
								responseJSON("DatabaseError", w, http.StatusInternalServerError)
							} else {
								responseJSON("Correct", w, http.StatusOK)
							}
						} else {
							responseJSON("Wrong", w, http.StatusOK)
						}
					} else {
						responseJSON("LevelNoMatch", w, http.StatusOK)
					}
				} else {
					responseJSON("InvalidLevel", w, http.StatusOK)
				}
			} else {
				responseJSON("InvalidToken", w, http.StatusOK)
			}
		}))))
	router.Handle("/level", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.URL.Query()["id_token"][0]
			claims := jwt.MapClaims{}
			jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				x, err := getPemCert(token)
				return []byte(x), err
			})
			email := claims["email"].(string)
			find, _ := collection.Find(context.TODO(), bson.M{"clientID": email})
			JSOND, _ := json.Marshal(find.Next(context.TODO()))
			if strings.Compare(string(JSOND), "true") == 0 {
				var current DatabaseUserObject
				err := find.Decode(&current)
				if err != nil {
					fmt.Println("Not able to read database object")
					responseJSON("DatabaseError", w, http.StatusInternalServerError)
				} else {
					var resp LevelResponse
					if current.Level == 0 {
						resp = LevelResponse{0, "Level 0: The Flyer contains a clue. "}
					} else if current.Level == 1 {
						resp = LevelResponse{1, "Level 1: The Zozimus Logo"}
					} else if current.Level == 2 {
						resp = LevelResponse{2, "Level 2: The room number. "}
					} else if current.Level == 3 {
						resp = LevelResponse{3, "Level 3: A code. "}
					} else if current.Level == 4 {
						resp = LevelResponse{4, "Level 4: <a href='https://docs.google.com/document/d/e/2PACX-1vRbmcKXV1JTztY31VXhdcNj6jKLWJHaMzaHTul3uIdZYYmdhLzuJS55mb4I2YTn-wLGfvs0-uDjHItu/pub'>Click here. </a>"}
					} else if current.Level == 5 {
						resp = LevelResponse{5, "Level 5: The Modern name of the hairpin. "}
					} else if current.Level == 6 {
						resp = LevelResponse{6, "Level 6: Connect. "}
					} else if current.Level == 7 {
						resp = LevelResponse{7, "Level 7: The End. "}
					} else {
						resp = LevelResponse{8, "Won"}
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					jData, _ := json.Marshal(resp)
					w.Write(jData)
				}
			}
		}))))
	router.HandleFunc("/leaderboard", LeaderboardHandler)
	router.HandleFunc("/leaderboardtable", LeaderboardTableHandler)
	router.HandleFunc("/rules", RulesHandler)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("build")))
	http.ListenAndServe(":8080", router)
}
func CSSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./src/App.css")
}
func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://dev-l0ini8h1.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}
func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
func responseJSON(message string, w http.ResponseWriter, statusCode int) {
	enableCors(&w)
	response := Response{message}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

func LeaderboardHandler(w http.ResponseWriter, request *http.Request) {
	options := options.Find()
	options.SetSort(bson.D{{"level", -1}, {"lastModified", 1}})
	find, _ := collection.Find(context.TODO(), bson.M{}, options)
	var results []user
	for find.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem user
		err := find.Decode(&elem)
		fmt.Println(elem)
		if err != nil {
			fmt.Println("Error decoding leaderboard item")
		}
		results = append(results, elem)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jData, _ := json.Marshal(results)
	w.Write(jData)
}

func LevelQueryHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	find, _ := collection.Find(context.TODO(), bson.M{"clientID": vars["clientid"]})
	JSOND, _ := json.Marshal(find.Next(context.TODO()))
	// Returning the level of the queried user
	if strings.Compare(string(JSOND), "true") == 0 {
		var current DatabaseUserObject
		_ = find.Decode(&current)
		responseJSON(strconv.Itoa(current.Level), w, http.StatusOK)
	} else {
		responseJSON("-2", w, http.StatusOK)
	}
}

func DoesUsernameExistHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	find, _ := collection.Find(context.TODO(), bson.M{"username": vars["username"]})
	JSOND, _ := json.Marshal(find.Next(context.TODO()))
	// Returning the level of the queried user
	if strings.Compare(string(JSOND), "true") == 0 {
		responseJSON("true", w, http.StatusOK)
	} else {
		responseJSON("false", w, http.StatusOK)
	}
}

func LeaderboardTableHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "leaderboard.html")
}

func RulesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "rules.html")
}

// func submitAnswer(w http.ResponseWriter, request *http.Request) {
//     vars := mux.Vars(request)
//     client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
//     context.TODO(), _ := context.WithTimeout(context.Background(), 10*time.Second)
//     _ = client.Connect(context.TODO())
//     collection := client.Database("Cryptex").Collection("users")
//     context.TODO(), _ = context.WithTimeout(context.Background(), 5*time.Second)

// }
// Function is obsoelete, implemented using GraphQL in main()
// func RetrieveLevel(w http.ResponseWriter, request *http.Request) {
//     vars := mux.Vars(request)
//     client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
//     context.TODO(), _ := context.WithTimeout(context.Background(), 10*time.Second)
//     err = client.Connect(context.TODO())
//     fmt.Println(err)
//     collection := client.Database("Cryptex").Collection("users")
//     context.TODO(), _ = context.WithTimeout(context.Background(), 5*time.Second)
//     filter := bson.M{"clientID" : vars["ID"]}
//     var result map[string]interface{}
//     err = collection.FindOne(context.TODO(), filter).Decode(&result)

// }

// Does not provide any fine tuning. Adjust CORS funciton later.
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
