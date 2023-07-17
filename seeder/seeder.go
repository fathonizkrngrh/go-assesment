package seeder

import (
	"context"
	"fmt"
	"github.com/gocroot/gocroot/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func GeneratePassword(pw string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password")
	}
	password := string(hashedPassword)

	return password
}

func SeedData(db *mongo.Database) error {
	roleCol := db.Collection("roles")
	userCol := db.Collection("users")

	// Insert roles
	roleIDs := make(map[string]primitive.ObjectID)
	roles := []models.Role{
		{Name: "USER"},
		{Name: "ADMIN"},
		{Name: "SUPERADMIN"},
	}

	for _, role := range roles {
		roleID := GenerateObjectID()
		roleIDs[role.Name] = roleID
		_, err := roleCol.InsertOne(context.Background(), role)
		if err != nil {
			return err
		}
		fmt.Printf("Inserted Role: %s with ID: %s\n", role.Name, roleID.Hex())
	}

	// Insert users with RoleIDs
	users := []models.User{
		{
			Username: "user",
			Email:    "user@email.com",
			RoleID:   roleIDs["USER"],
			Password: GeneratePassword("1234"),
		},
		{
			Username: "admin",
			Email:    "admin@email.com",
			RoleID:   roleIDs["ADMIN"],
			Password: GeneratePassword("1234"),
		},
		{
			Username: "superadmin",
			Email:    "superadmin@email.com",
			RoleID:   roleIDs["SUPERADMIN"],
			Password: GeneratePassword("1234"),
		},
	}

	for _, user := range users {
		_, err := userCol.InsertOne(context.Background(), user)
		if err != nil {
			return err
		}
		fmt.Printf("Inserted User: %s with ID: %s and RoleID: %s\n", user.Username, user.ID.Hex(), user.RoleID.Hex())
	}

	return nil
}
