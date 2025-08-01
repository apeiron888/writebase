package repository

import (
	"context"
	"errors"
	"starter/internal/domain"
	dtodbrep "starter/internal/repository/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoReportRepository struct {
	collection *mongo.Collection
}

func NewMongoReportRepository(collection *mongo.Collection) *MongoReportRepository {
	return &MongoReportRepository{collection: collection}
}

func (r *MongoReportRepository) CreateReport(ctx context.Context, report *domain.Report) error {
	dto := dtodbrep.ReportResponse{
		ID:         report.ID,
		ReporterID: report.ReporterID,
		TargetID:   report.TargetID,
		TargetType: report.TargetType,
		Reason:     report.Reason,
		Status:     string(report.Status),
		CreatedAt:  report.CreatedAt,
		ResolvedAt: report.ResolvedAt,
	}
	_, err := r.collection.InsertOne(ctx, dto)
	return err
}

func (r *MongoReportRepository) GetReports(ctx context.Context, filter map[string]interface{}) ([]*domain.Report, error) {
	var results []*domain.Report
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var dto dtodbrep.ReportResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		results = append(results, &domain.Report{
			ID:         dto.ID,
			ReporterID: dto.ReporterID,
			TargetID:   dto.TargetID,
			TargetType: dto.TargetType,
			Reason:     dto.Reason,
			Status:     domain.ReportStatus(dto.Status),
			CreatedAt:  dto.CreatedAt,
			ResolvedAt: dto.ResolvedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoReportRepository) UpdateReportStatus(ctx context.Context, reportID string, status domain.ReportStatus) error {
	filter := bson.M{"_id": reportID}
	update := bson.M{"$set": bson.M{"status": string(status)}}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrReportNotFound
	}
	return nil
}
