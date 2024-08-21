package db

import (
	"context"
	"errors"
	"fmt"
	"go-rest-api/internal/apperror"
	"go-rest-api/internal/user"
	"go-rest-api/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) FindAll(ctx context.Context) (users []user.User, err error) {
	d.logger.Debug("finding all users")
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return users, fmt.Errorf("can't find users: %w", err)
	}
	if err = cursor.All(ctx, &users); err != nil {
		return users, fmt.Errorf("can't decode users: %w", err)
	}
	return users, nil
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("creating user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("can't create user: %w", err)
	}
	d.logger.Debug("converting InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("can't convert object id to string")
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	d.logger.Debugf("finding user by id: %s", id)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("can't convert string to object id: %w", err)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("can't find user: %w", result.Err())
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("can't decode user: %w", err)
	}
	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	d.logger.Debug("updating user by id: %s", user.ID)
	if user.ID == "" {
		return fmt.Errorf("user ID is not set or invalid")
	}
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("can't convert user ID to ObjectID: %w", err)
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{"$set": bson.M{
		"username": user.Username,
		"email":    user.Email,
		"password": user.PasswordHash,
	}}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}
	d.logger.Tracef("Matched: %d, Modified: %d", result.MatchedCount, result.ModifiedCount)
	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	d.logger.Debugf("deleting user by id: %s", id)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("can't convert userID to ObjectID: %w, userID=%s", err, id)
	}
	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("can't delete user: %w", err)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}
	d.logger.Tracef("Deleted: %d", result.DeletedCount)
	return nil
}
