package repository

import (
	"context"
	"fmt"
	"time"
	"write_base/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct{
	UserCollection *mongo.Collection
	RefreshTokenCollection *mongo.Collection
	EmailTokenCollection *mongo.Collection
}
func NewUserRepository(db *mongo.Database) domain.IUserRepository {
    refreshTokenCol := db.Collection("RefreshToken")
    emailTokenCol := db.Collection("EmailToken")

    ctx := context.Background()
    // TTL index for refresh tokens
    _, err := refreshTokenCol.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{{Key: "expires_at", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(0),
    })
    if err != nil {
        fmt.Println("Failed to create TTL index for RefreshToken:", err)
    }

    // TTL index for email tokens
    _, err = emailTokenCol.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{{Key: "expires_at", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(0),
    })
    if err != nil {
        fmt.Println("Failed to create TTL index for EmailToken:", err)
    }

    return &UserRepository{
        UserCollection: db.Collection("users"),
        RefreshTokenCollection: refreshTokenCol,
        EmailTokenCollection: emailTokenCol,
    }
}
func (r *UserRepository) DeleteUnverifiedExpiredUsers(ctx context.Context, expiration time.Duration) error {
    threshold := primitive.NewDateTimeFromTime(time.Now().Add(-expiration))
    _, err := r.UserCollection.DeleteMany(ctx, bson.M{
        "verified": false,
        "created_at": bson.M{"$lt": threshold},
    })
    return err
}
func (r *UserRepository) DeleteOldRevokedTokens(ctx context.Context, olderThan time.Duration) error {
    threshold := primitive.NewDateTimeFromTime(time.Now().Add(-olderThan))
    _, err := r.RefreshTokenCollection.DeleteMany(ctx, bson.M{
        "revoked": true,
        "revoked_at": bson.M{"$lt": threshold},
    })
    return err
}
// CreateUser inserts a new user into the users collection.
func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
    dto := userDomainToDTO(user)
    _, err := r.UserCollection.InsertOne(ctx, dto)
    return err
}

// GetByEmail finds a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    var dto userDTO
    err := r.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&dto)
	if err != nil {
        return nil, err
    }
    return userDTOToDomain(dto), nil
}

// GetByID finds a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    var dto userDTO
    err := r.UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&dto)
    if err != nil {
        return nil, err
    }
    return userDTOToDomain(dto), nil
}

// GetByUsername finds a user by username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
    var dto userDTO
    err := r.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&dto)
    if err != nil {
        return nil, err
    }
    return userDTOToDomain(dto), nil
}

// UpdateUser updates an existing user.
func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
    dto := userDomainToDTO(user)
    _, err := r.UserCollection.UpdateOne(
        ctx,
        bson.M{"_id": user.ID},
        bson.M{"$set": dto},
    )
    return err
}

// DeleteUser deletes a user by ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
    _, err := r.UserCollection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

// PromoteToAdmin sets the user's role to admin.
func (r *UserRepository) PromoteToAdmin(ctx context.Context, userID string) error {
    _, err := r.UserCollection.UpdateOne(
        ctx,
        bson.M{"_id": userID},
        bson.M{"$set": bson.M{"role": "admin"}},
    )
    return err
}

// DemoteToUser sets the user's role to user.
func (r *UserRepository) DemoteToUser(ctx context.Context, userID string) error {
    _, err := r.UserCollection.UpdateOne(
        ctx,
        bson.M{"_id": userID},
        bson.M{"$set": bson.M{"role": "user"}},
    )
    return err
}
// StoreToken inserts a new refresh token.
func (r *UserRepository) StoreToken(ctx context.Context, token *domain.RefreshToken) error {
    dto := refreshTokenDomainToDTO(token)
    _, err := r.RefreshTokenCollection.InsertOne(ctx, dto)
    return err
}

func (r *UserRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
    var dto refreshTokenDTO
    err := r.RefreshTokenCollection.FindOne(ctx, bson.M{"token": token}).Decode(&dto)
    if err != nil {
        return nil, err
    }
    return refreshTokenDTOToDomain(dto), nil
}

// GetValidByUser returns all valid (not revoked, not expired) refresh tokens for a user.
func (r *UserRepository) GetValidByUser(ctx context.Context, userID string) ([]*domain.RefreshToken, error) {
    now := primitive.NewDateTimeFromTime(time.Now())
    cursor, err := r.RefreshTokenCollection.Find(ctx, bson.M{
        "user_id": userID,
        "revoked": false,
        "expires_at": bson.M{"$gt": now},
    })
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    var tokens []refreshTokenDTO
    if err := cursor.All(ctx, &tokens); err != nil {
        return nil, err
    }
    var result []*domain.RefreshToken
    for _, dto := range tokens {
        result = append(result, refreshTokenDTOToDomain(dto))
    }
    return result, nil
}

// RevokeToken sets the revoked flag to true for a given token.
func (r *UserRepository) RevokeToken(ctx context.Context, token string) error {
    now := primitive.NewDateTimeFromTime(time.Now())
    _, err := r.RefreshTokenCollection.UpdateOne(
        ctx,
        bson.M{"token": token},
        bson.M{"$set": bson.M{"revoked": true, "revoked_at": now}},
    )
    return err
}

// DeleteExpiredTokens removes all tokens that have expired.
func (r *UserRepository) DeleteExpiredTokens() error {
    now := primitive.NewDateTimeFromTime(time.Now())
    _, err := r.RefreshTokenCollection.DeleteMany(
        context.Background(),
        bson.M{"expires_at": bson.M{"$lte": now}},
    )
    return err
}

// RevokeAllByUser revokes all tokens for a user.
func (r *UserRepository) RevokeAllByUser(ctx context.Context, userID string) error {
    now := primitive.NewDateTimeFromTime(time.Now())
    _, err := r.RefreshTokenCollection.UpdateMany(
        ctx,
        bson.M{"user_id": userID, "revoked": false},
        bson.M{"$set": bson.M{"revoked": true, "revoked_at": now}},
    )
    return err
}

// SaveVerificationToken inserts a new email verification token.
func (r *UserRepository) SaveVerificationToken(ctx context.Context, emailToken *domain.EmailVerificationToken) error {
    dto := emailVerificationTokenDomainToDTO(emailToken)
    _, err := r.EmailTokenCollection.InsertOne(ctx, dto)
    return err
}

// GetVerificationToken finds an email verification token by its token string.
func (r *UserRepository) GetVerificationToken(ctx context.Context, emailToken string) (*domain.EmailVerificationToken, error) {
    var dto emailVerificationTokenDTO
    err := r.EmailTokenCollection.FindOne(ctx, bson.M{"token": emailToken}).Decode(&dto)
    if err != nil {
        return nil, err
    }
    return emailVerificationTokenDTOToDomain(dto), nil
}
func (r *UserRepository) DeleteVerificationToken(ctx context.Context, token string) error {
    _, err := r.EmailTokenCollection.DeleteOne(ctx, bson.M{"token": token})
    return err
}