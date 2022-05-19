package mongo

import (
	"context"
	"time"

	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListPrs() ([]*model.Pr, error) {
	ctx := context.TODO()

	var prs []*model.Pr

	cur, err := prCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when listing prs")
		return nil, err
	}

	for cur.Next(ctx) {
		var pr model.Pr
		err := cur.Decode(&pr)
		if err != nil {
			log.Error().Err(err).Msgf("there was an error when decoding pr")
			return nil, err
		}

		prs = append(prs, &pr)
	}

	log.Info().Msgf("found %s pr envs", len(prs))

	return prs, nil
}

func InsertPr(pr *model.Pr) error {
	ctx := context.TODO()

	pr.ID = primitive.NewObjectID()
	pr.UpdatedAt = time.Now()
	pr.CreatedAt = time.Now()

	result, err := prCollection.InsertOne(ctx, pr)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when upserting pr with id %s", pr.ID)
		return err
	}

	log.Info().Msgf("%v id, document was inserted", result.InsertedID)
	return nil
}

func UpdatePr(pr *model.Pr) error {
	ctx := context.TODO()

	pr.UpdatedAt = time.Now()
	result, err := prCollection.ReplaceOne(ctx, bson.D{{"_id", pr.ID}}, pr)

	if err != nil {
		log.Error().Err(err).Msgf("there was an error when updating pr with id %s", pr.ID)
		return err
	}

	log.Info().Msgf("successfully replaced %d documents", result.MatchedCount)
	return nil
}

func DeletePr(pr *model.Pr) error {
	ctx := context.TODO()

	result, err := prCollection.DeleteOne(ctx, bson.D{{"_id", pr.ID}})
	if err != nil {
		log.Error().Err(err).Msgf("error when deleting pr with id %s", pr.ID)
		return err
	}

	log.Info().Msgf("successfully deleted %d documents with id %s", result.DeletedCount, pr.ID)
	return nil
}
