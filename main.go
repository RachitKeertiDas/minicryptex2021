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
	"github.com/dgrijalva/jwt-go"
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
	answers["0"] = "ladygodiva"
	answers["1"] = "triskaidekaphobia"
	answers["2"] = "beatles"
	answers["3"] = "fcuk"
	answers["4"] = "stanlee"
	answers["5"] = "pi"
	answers["6"] = "pisces"
	answers["7"] = "pabloescobar"
	answers["8"] = "welovecryptex67435"
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
	router.HandleFunc("/backend/whichlevel/{clientid}", LevelQueryHandler)
	router.HandleFunc("/backend/css", CSSHandler)
	router.HandleFunc("/backend/doesUsernameExist/{username}", DoesUsernameExistHandler)
	router.Handle("/backend/adduser/{ID}/{username}/{name1}/{name2}/{name3}/{name4}/{name5}", negroni.New(
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
	router.Handle("/backend/acceptedrules", negroni.New(
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
		router.Handle("/backend/answer/{level}/{answer}", negroni.New(
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
	router.Handle("/backend/level", negroni.New(
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
						resp = LevelResponse{0, "Level 1:<br/><img width=\"auto\" height=\"450px\" margin-top=\"50px\" src=\"https://res.cloudinary.com/ddc4oysii/image/upload/v1579860040/F8258061161AFD91C588E71FF6DC5AC103FAE288E9D4FDB1110A35618E88AF7C/7496CE0698D7F149CFEF1C096EAF0944F6851C01696E8143BA6ABB69C466290E_a7h02o.png\"> "}
					} else if current.Level == 1 {
						resp = LevelResponse{1, "Level 2:<br/><img width=\"auto\" height=\"450px\" margin-top=\"50px\" src=\"https://res.cloudinary.com/drgddftct/image/upload/v1547292346/QPADBgJd8EkeBut6.png\">"}
					} else if current.Level == 2 {
						resp = LevelResponse{2, "Level 3: <br/><a href='https://drive.google.com/file/d/1UQGx-2_ZZxnEo65-oYXZgzJ5u2GdLabM/view?usp=sharing'>Click here. </a> "}
					} else if current.Level == 3 {
						resp = LevelResponse{3, "Level 4:<br/> <img width=\"auto\" height=\"400px\" margin-top=\"50px\" src=\"https://res.cloudinary.com/dmridruee/image/upload/v1547192728/fpF6juWJPP7D2S9BJWcc/LQtD12ldlFRZ4OT90cDj.png\"> "}
					} else if current.Level == 4 {
						resp = LevelResponse{4,"Level 5:You took this long to get to the fifth question?<br>"+
 			"HAAHHAHA HAAAAHHH HAAHHHAA HAAHHAHA HAAHAAHH HAAAHHAA HAAHAHHA HAAHAAAA<br/>"+
 			"HAAAHHAH HHAHAAHH HHAHHHHH HAAHAHHH HAAHHAHA HHAHHHHH HAAHHAHA HAAAAHHH<br/>"+
 			"HAAHHHAA HAAHAAHH HAAHHHHA HAAHAHHA HAAHAAHA HAAHHAHA HAAHHAHH HHHHAHAH<br/>"+
 			"HAHHAHHH HAAHAHHA HAAAHHAA HHAHHHHH HAAHAAHH HAAHAHHA HAAHHAAH HAAHHAHA<br/>"+
 			"HHAHHHHH HAAHHHHA HAAHAAAH HAAHHAHH HHAHHHHH HAAAHHAA HAAAHAHH HAAHAAAA<br/>"+
			"HAAAHHAH HAAHAHHA HAAHHAHA HAAAHHAA HHAHAAHH HHAHHHHH HAAHHHAA HAAHAAAA<br/>"+
 			"HAAHAAHA HAAHAHHA HAAHAAAH HAAHHAAA HHAHHHHH HAAAHAHH HAAHAAAA HHAHHHHH<br/>"+
 			"HAAHHHHA HHAHHHHH HAAHAHHH HAAHHHHA HAAHAAHH HAAAHAHH HHHHAHAH<br/>"+
 			"You two-dimensional, depth lacking loser!"}
					} else if current.Level == 5 {
						resp = LevelResponse{5, "Level 6:<br/> <img width=\"700px\" height=\"auto\" src=\"https://res.cloudinary.com/ddc4oysii/image/upload/v1579860364/F8258061161AFD91C588E71FF6DC5AC103FAE288E9D4FDB1110A35618E88AF7C/7B2E50D7BDBDA6B2D92268D2498E80471C9D8EE0BD5E708FC28E8E3F00E91322_xez7qs.jpg\"> "}
					} else if current.Level == 6 {
						resp = LevelResponse{6, "Level 7:<img src=\"https://res.cloudinary.com/drgddftct/image/upload/v1605967394/question_dg.png\">"}
					} else if current.Level == 7 {
						resp = LevelResponse{7, "Level 8:<br/> <img width=\"700px\" height=\"auto\" src=\"https://res.cloudinary.com/dmridruee/image/upload/v1547211291/0PNQNGAOck2NQwyb6hQV.png\"> "}
					} else if current.Level == 8 {
						resp = LevelResponse{8,"Level 9: <br/> <img width=\"700px\" height=\"auto\" src=\"https://res.cloudinary.com/drgddftct/image/upload/v1606051168/level9.png\"> "}
					}else {
						resp = LevelResponse{9, "Congrats, You have Won."}
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					jData, _ := json.Marshal(resp)
					w.Write(jData)
				}
			}
		}))))
	router.HandleFunc("/backend/leaderboard", LeaderboardHandler)
	router.HandleFunc("/backend/leaderboardtable", LeaderboardTableHandler)
	router.HandleFunc("/backend/rules", RulesHandler)

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
		//fmt.Println(elem)
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
