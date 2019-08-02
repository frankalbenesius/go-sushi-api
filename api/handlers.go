package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Roll is model for sushi
type Roll struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Ingredients string             `json:"ingredients,omitempty" bson:"ingredients,omitempty"`
}

const dbName = "test"

func (s *SushiAPI) handleGetRolls() http.HandlerFunc {
	collection := s.dbClient.Database(dbName).Collection("rolls")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var results []*Roll
		cursor, err := collection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}
		for cursor.Next(context.TODO()) {
			var elem Roll
			err := cursor.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, &elem)
		}
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
		cursor.Close(context.TODO())
		json.NewEncoder(w).Encode(results)
	}
}

func (s *SushiAPI) handleGetRoll() http.HandlerFunc {
	collection := s.dbClient.Database(dbName).Collection("rolls")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var roll Roll
		filter := bson.D{{"_id", id}}
		err := collection.FindOne(context.TODO(), filter).Decode(&roll)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}

		json.NewEncoder(w).Encode(roll)
	}
}

func (s *SushiAPI) handleCreateRoll() http.HandlerFunc {
	collection := s.dbClient.Database(dbName).Collection("rolls")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var roll Roll
		json.NewDecoder(r.Body).Decode(&roll)
		result, _ := collection.InsertOne(context.TODO(), roll)
		json.NewEncoder(w).Encode(result) // returns {  "InsertedID": "000000000000000000000000" }
	}
}

func (s *SushiAPI) handleUpdateRoll() http.HandlerFunc {
	collection := s.dbClient.Database(dbName).Collection("rolls")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])

		var roll Roll
		json.NewDecoder(r.Body).Decode(&roll)
		filter := bson.D{{"_id", id}}
		update := bson.D{
			{"$set", bson.D{
				{"name", roll.Name},
				{"ingredients", roll.Ingredients},
			}},
		}
		result, _ := collection.UpdateOne(context.TODO(), filter, update)
		json.NewEncoder(w).Encode(result)
		// returns an object that looks like:
		// {
		// 	"MatchedCount": 1,
		// 	"ModifiedCount": 1,
		// 	"UpsertedCount": 0,
		// 	"UpsertedID": null
		// }
		// it may make more sense to return the updated document?
	}
}

func (s *SushiAPI) handleDeleteRoll() http.HandlerFunc {
	collection := s.dbClient.Database(dbName).Collection("rolls")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, _ := primitive.ObjectIDFromHex(params["id"])
		filter := bson.D{{"_id", id}}
		result, _ := collection.DeleteOne(context.TODO(), filter)
		json.NewEncoder(w).Encode(result)
		/*
			{
			  "DeletedCount": 1
			}
		*/
	}
}
