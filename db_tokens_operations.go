package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func createTokens(client *mongo.Client, collection *mongo.Collection, tokens []interface{}) {
	var err error
	var insertManyResult *mongo.InsertManyResult
	var insertSession mongo.Session

	if insertSession, err = client.StartSession(); err != nil {
		log.Println(err)
	}
	if err = insertSession.StartTransaction(); err != nil {
		log.Println(err)
	}
	if err = mongo.WithSession(context.TODO(), insertSession, func(sc mongo.SessionContext) error {
		if insertManyResult, err = collection.InsertMany(context.TODO(), tokens); err != nil {
			log.Println(err)
		}
		if len(insertManyResult.InsertedIDs) != len(tokens) {
			fmt.Printf("Insert failed, expected %v but got %v.\n", len(tokens), len(insertManyResult.InsertedIDs))
		}

		if err = insertSession.CommitTransaction(sc); err != nil {
			log.Println(err)
		}
		return nil
	}); err != nil {
		log.Println(err)
	}
	insertSession.EndSession(context.TODO())

	fmt.Println("Inserted documents: \n", insertManyResult.InsertedIDs)
}

func updateTokens(client *mongo.Client, collection *mongo.Collection, filter bson.D, update bson.D) {
	var err error
	var updateResult *mongo.UpdateResult
	var updateSession mongo.Session

	if updateSession, err = client.StartSession(); err != nil {
		log.Println(err)
	}
	if err = updateSession.StartTransaction(); err != nil {
		log.Println(err)
	}
	if err = mongo.WithSession(context.TODO(), updateSession, func(sc mongo.SessionContext) error {
		if updateResult, err = collection.UpdateMany(context.TODO(), filter, update); err != nil {
			log.Println(err)
		}
		if updateResult.ModifiedCount != updateResult.MatchedCount {
			fmt.Printf("Update failed, expected %v but got %v.\n", int(updateResult.MatchedCount), int(updateResult.ModifiedCount))
		}

		if err = updateSession.CommitTransaction(sc); err != nil {
			log.Println(err)
		}
		return nil
	}); err != nil {
		log.Println(err)
	}
	updateSession.EndSession(context.TODO())

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func deleteTokens(client *mongo.Client, collection *mongo.Collection, filter bson.D) {
	var err error
	var deleteManyResult *mongo.DeleteResult
	var deleteSession mongo.Session

	if deleteSession, err = client.StartSession(); err != nil {
		log.Println(err)
	}
	if err = deleteSession.StartTransaction(); err != nil {
		log.Println(err)
	}
	if err = mongo.WithSession(context.TODO(), deleteSession, func(sc mongo.SessionContext) error {
		if deleteManyResult, err = collection.DeleteMany(context.TODO(), filter); err != nil {
			log.Println(err)
		}

		if err = deleteSession.CommitTransaction(sc); err != nil {
			log.Println(err)
		}
		return nil
	}); err != nil {
		log.Println(err)
	}
	deleteSession.EndSession(context.TODO())

	fmt.Printf("Deleted %v documents.\n", int(deleteManyResult.DeletedCount))
}

func findTokensByGuid(collection *mongo.Collection, guid string) []*Token {
	filter := bson.D{{"guid", guid}}

	// Array with the decoded documents
	var results []*Token

	// Passing nil as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Token
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple tokens by guid %s: %+v\n", guid, results)
	return results
}
