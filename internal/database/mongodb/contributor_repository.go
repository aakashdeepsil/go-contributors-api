package mongodb

import (
	"context"
	"time"

	"github.com/aakashdeepsil/go-contributors-api/internal/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContributorRepository struct {
	collection *mongo.Collection
}

func NewContributorRepository(db *mongo.Database) *ContributorRepository {
	collection := db.Collection("contributors")

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		panic(err)
	}

	return &ContributorRepository{collection: collection}
}

func (r *ContributorRepository) Create(ctx context.Context, input *models.ContributorInput) (*models.Contributor, error) {
	contributor := &models.Contributor{
		Username:  input.Username,
		Email:     input.Email,
		Name:      input.Name,
		AvatarURL: input.AvatarURL,
		Projects:  input.Projects,
		JoinedAt:  time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, contributor)
	if err != nil {
		return nil, err
	}

	contributor.ID = result.InsertedID.(primitive.ObjectID)
	return contributor, nil
}

func (r *ContributorRepository) GetByID(ctx context.Context, id string) (*models.Contributor, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var contributor models.Contributor
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&contributor)
	if err != nil {
		return nil, err
	}

	return &contributor, nil
}

func (r *ContributorRepository) GetByUsername(ctx context.Context, username string) (*models.Contributor, error) {
	var contributor models.Contributor
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&contributor)
	if err != nil {
		return nil, err
	}

	return &contributor, nil
}

func (r *ContributorRepository) Update(ctx context.Context, id string, input *models.ContributorInput) (*models.Contributor, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"username":  input.Username,
			"email":     input.Email,
			"name":      input.Name,
			"avatarUrl": input.AvatarURL,
			"projects":  input.Projects,
			"updatedAt": time.Now(),
		},
	}

	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var contributor models.Contributor
	if err := result.Decode(&contributor); err != nil {
		return nil, err
	}

	return &contributor, nil
}

func (r *ContributorRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *ContributorRepository) List(ctx context.Context, limit, offset int) ([]*models.Contributor, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "joinedAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contributors []*models.Contributor
	if err = cursor.All(ctx, &contributors); err != nil {
		return nil, err
	}

	return contributors, nil
}
