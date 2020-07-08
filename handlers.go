package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

// Refresh operation
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	client := connect()
	collection := client.Database(databaseName).Collection("tokens")

	params := mux.Vars(r)
	guid := params["guid"]
	fmt.Printf("GUID: %v\n", guid)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Println("Refresh token wasn't read!")
	} else {
		oldRefreshToken, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Old refresh token: %v\n", string(oldRefreshToken))
		userTokens := findTokensByGuid(collection, guid) // get all tokens for user
		for _, token := range userTokens {
			// update refresh and access token if refresh token from cookie was find
			if verifyRefreshToken([]byte(token.RefreshToken), oldRefreshToken) {
				refreshToken, refreshExpiresAt := tokenGenerator(), time.Now().Add(time.Hour*72)
				accessToken, accessExpiresAt := newJwtToken(guid, signingKey, time.Minute*30)

				filter := bson.D{{"refreshtoken", token.RefreshToken}}
				hashedRefreshToken := hashToken([]byte(refreshToken)) // to hash new refresh token
				update := bson.D{
					{"$set", bson.D{
						{"refreshtoken", hashedRefreshToken},
					}},
				}
				updateTokens(client, collection, filter, update)

				addCookie("access_token", accessToken, accessExpiresAt, false, w)
				// add refresh token to cookies with httpOnly
				addCookie("refresh_token", base64.StdEncoding.EncodeToString([]byte(refreshToken)),
					refreshExpiresAt, true, w)
				break
			}
		}
	}
	// redirect to application server page. Disabled for watching cookies after request
	//http.Redirect(w, r, "/", 301)
	disconnect(client)
}

// Create tokens by GUID
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	client := connect()
	collection := client.Database(databaseName).Collection("tokens")

	params := mux.Vars(r)
	guid := params["guid"]
	fmt.Printf("GUID: %v\n", guid)

	refreshToken, refreshExpiresAt := tokenGenerator(), time.Now().Add(time.Hour*72)
	accessToken, accessExpiresAt := newJwtToken(guid, signingKey, time.Minute*30)
	hashedRefreshToken := hashToken([]byte(refreshToken))

	document := Token{guid, hashedRefreshToken}
	tokens := []interface{}{document}

	createTokens(client, collection, tokens)

	addCookie("access_token", accessToken, accessExpiresAt, false, w)
	addCookie("refresh_token", base64.StdEncoding.EncodeToString([]byte(refreshToken)),
		refreshExpiresAt, true, w)

	//http.Redirect(w, r, "/", 301)
	disconnect(client)
}

// Delete token by refresh
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	client := connect()
	collection := client.Database(databaseName).Collection("tokens")

	params := mux.Vars(r)
	guid := params["guid"]
	fmt.Printf("GUID: %v\n", guid)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Println("Refresh token wasn't read!")
	} else {
		oldRefreshToken, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Refresh token for deleting: %v\n", string(oldRefreshToken))

		userTokens := findTokensByGuid(collection, guid)
		for _, token := range userTokens {
			if verifyRefreshToken([]byte(token.RefreshToken), oldRefreshToken) {
				filter := bson.D{{"refreshtoken", token.RefreshToken}}
				deleteTokens(client, collection, filter)
				break
			}
		}
	}
	http.Redirect(w, r, "/", 301)
	disconnect(client)
}

// Delete all tokens by GUID
func DeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	client := connect()
	collection := client.Database(databaseName).Collection("tokens")

	params := mux.Vars(r)
	guid := params["guid"]
	fmt.Printf("GUID: %v\n", guid)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Println("Refresh token wasn't read!")
	} else {
		oldRefreshToken, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Refresh token: %v\n", string(oldRefreshToken))

		userTokens := findTokensByGuid(collection, guid)
		for _, token := range userTokens {
			if verifyRefreshToken([]byte(token.RefreshToken), oldRefreshToken) {
				filter := bson.D{{"guid", guid}}
				deleteTokens(client, collection, filter)
				break
			}
		}
	}
	disconnect(client)
	http.Redirect(w, r, "/", 301)
}
