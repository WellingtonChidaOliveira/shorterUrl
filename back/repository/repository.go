package repository

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/wellingtonchida/shortner/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type DataBase interface {
	GetAll(ctx context.Context) ([]models.Shorter, error)
	GetById(ctx context.Context, id string) (models.Shorter, error)
	Create(ctx context.Context, shorter models.Shorter) (string, error)
	Update(ctx context.Context, id string, shorter models.Shorter) error
	Delete(ctx context.Context, id string) error
	Inactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error 
}

type dataBase struct {
	coll *mongo.Collection
}

func NewDataBase(ctx context.Context) DataBase {
	client := Connect(ctx)
	dbname := os.Getenv("MONGO_DB_NAME")

	return &dataBase{
		coll: client.Database(dbname).Collection(dbname),
	}
}

func Connect(ctx context.Context) *mongo.Client {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	pass := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)

	opt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opt)

	if err != nil {
		panic("Error to connect to database!")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return client
}

func (db *dataBase) GetAll(ctx context.Context) ([]models.Shorter, error) {

	cursor, err := db.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []models.Shorter
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *dataBase) GetById(ctx context.Context, id string) (models.Shorter, error) {

	var filter bson.M
	oid, err := transFormId(id)
	if err == nil {
		filter = bson.M{"_id": oid}
	} else {
		filter = bson.M{"name": id}
	}

	var shorter models.Shorter
	err = db.coll.FindOne(ctx, filter).Decode(&shorter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return shorter, mongo.ErrNoDocuments
		}
		return shorter, fmt.Errorf("unspected error %v",err.Error())
	}

	return shorter, nil
}

func (db *dataBase) Create(ctx context.Context, shorter models.Shorter) (string, error) {
	result, err := db.coll.InsertOne(ctx, shorter)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", nil
	}
	return oid.Hex(), nil
}

func (db *dataBase) Update(ctx context.Context, id string, updateData models.Shorter) error {
	oid, err := transFormId(id)
	if err != nil {
		return err
	}

	filter := bson.M{}

	if updateData.Shorter != "" {
		filter["name"] = updateData.Shorter
	}
	if updateData.OriginUrl != "" {
		filter["url"] = updateData.OriginUrl
	}
	if updateData.MostUse != 0 {
		filter["most_use"] = updateData.MostUse
	}
	
	if len(filter) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	_, err = db.coll.UpdateOne(ctx,
		bson.M{"_id": oid}, 
		bson.M{"$set": filter},
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *dataBase) Delete(ctx context.Context, id string) error {
	oid, err := transFormId(id)
	if err != nil {
		return err
	}

	_, err = db.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}

func (db *dataBase) Inactivate(ctx context.Context, id string) error {
	oid, err := transFormId(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"status": 0,
		},
	}

	_, err = db.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (db *dataBase) Activate(ctx context.Context, id string) error {
	oid, err := transFormId(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"status": 1,
		},
	}

	_, err = db.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func transFormId(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Printf("\n error to convert: %v\n", id)
		fmt.Printf("\n error : %v\n", err.Error())
		return primitive.ObjectID{}, err
	}

	fmt.Printf("\n convert: %v\n", oid)
	return oid, nil
}


