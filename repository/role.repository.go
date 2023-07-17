package repository

import (
	"context"
	"errors"
	"github.com/gocroot/gocroot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RoleRepository struct {
	db *mongo.Database
	c  *mongo.Collection
}

func NewRoleRepository(database *mongo.Database, collectionName string) *RoleRepository {
	return &RoleRepository{
		db: database,
		c:  database.Collection(collectionName),
	}
}

func (rr *RoleRepository) Create(role models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := rr.c.InsertOne(ctx, role)
	return err
}

func (rr *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": name}
	var role models.Role
	err := rr.c.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &role, nil
}

func (rr *RoleRepository) GetRoleById(roleId string) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(roleId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	var role models.Role
	err = rr.c.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &role, nil
}
