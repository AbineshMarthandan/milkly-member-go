package repository

import (
	"context"
	"time"

	"milkly-member/config"
	"milkly-member/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MemberRepository interface {
	Create(ctx context.Context, member *entity.Member) (*entity.Member, error)
	GetByID(ctx context.Context, id string) (*entity.Member, error)
	GetByEmail(ctx context.Context, email string) (*entity.Member, error)
	Update(ctx context.Context, id string, member *entity.Member) (*entity.Member, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Member, error)
}

type memberRepository struct {
	collection *mongo.Collection
}

func NewMemberRepository(cfg *config.Config) MemberRepository {
	collection := cfg.MongoDB.Collection("members")
	return &memberRepository{
		collection: collection,
	}
}

func (r *memberRepository) Create(ctx context.Context, member *entity.Member) (*entity.Member, error) {
	member.ID = primitive.NewObjectID()
	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()
	member.Status = entity.MemberStatusActive

	result, err := r.collection.InsertOne(ctx, member)
	if err != nil {
		return nil, err
	}

	member.ID = result.InsertedID.(primitive.ObjectID)
	return member, nil
}

func (r *memberRepository) GetByID(ctx context.Context, id string) (*entity.Member, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var member entity.Member
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&member)
	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (r *memberRepository) GetByEmail(ctx context.Context, email string) (*entity.Member, error) {
	var member entity.Member
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&member)
	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (r *memberRepository) Update(ctx context.Context, id string, member *entity.Member) (*entity.Member, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	member.UpdatedAt = time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": member}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *memberRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *memberRepository) List(ctx context.Context, limit, offset int) ([]*entity.Member, error) {
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []*entity.Member
	for cursor.Next(ctx) {
		var member entity.Member
		if err := cursor.Decode(&member); err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	return members, cursor.Err()
}