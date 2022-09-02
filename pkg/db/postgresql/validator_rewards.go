package postgresql

import (
	"context"

	"github.com/cortze/eth2-state-analyzer/pkg/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func (p *PostgresDBService) createRewardsTable(ctx context.Context, pool *pgxpool.Pool) error {
	// create the tables
	_, err := pool.Exec(ctx, model.CreateValidatorRewardsTable)
	if err != nil {
		return errors.Wrap(err, "error creating rewards table")
	}
	return nil
}

func (p *PostgresDBService) InsertNewValidatorRow(valRewardsObj model.ValidatorRewards) error {

	_, err := p.psqlPool.Exec(p.ctx, model.InsertNewValidatorLineTable,
		valRewardsObj.ValidatorIndex,
		valRewardsObj.Slot,
		valRewardsObj.Epoch,
		valRewardsObj.ValidatorBalance,
		valRewardsObj.Reward,
		valRewardsObj.MaxReward,
		valRewardsObj.AttSlot,
		valRewardsObj.InclusionDelay,
		valRewardsObj.BaseReward,
		valRewardsObj.MissingSource,
		valRewardsObj.MissingTarget,
		valRewardsObj.MissingHead)
	if err != nil {
		return errors.Wrap(err, "error inserting row in validator rewards table")
	}
	return nil
}

func (p *PostgresDBService) UpdateValidatorRowReward(valRewardsObj model.ValidatorRewards) error {

	_, err := p.psqlPool.Exec(p.ctx, model.UpdateValidatorLineTable, valRewardsObj.ValidatorIndex, valRewardsObj.Slot, valRewardsObj.Reward)
	if err != nil {
		return errors.Wrap(err, "error inserting row in validator rewards table")
	}
	return nil
}

func (p PostgresDBService) AddtoQueueVal(queryID int, valRewardsObj model.ValidatorRewards, batch *pgx.Batch) {

	if queryID == 0 {
		batch.Queue(model.VALIDATOR_QUERIES[queryID],
			valRewardsObj.ValidatorIndex,
			valRewardsObj.Slot,
			valRewardsObj.Epoch,
			valRewardsObj.ValidatorBalance,
			valRewardsObj.Reward,
			valRewardsObj.MaxReward,
			valRewardsObj.AttSlot,
			valRewardsObj.InclusionDelay,
			valRewardsObj.BaseReward,
			valRewardsObj.MissingSource,
			valRewardsObj.MissingTarget,
			valRewardsObj.MissingHead)
	}

	if queryID == 1 {
		batch.Queue(model.VALIDATOR_QUERIES[queryID],
			valRewardsObj.ValidatorIndex, valRewardsObj.Slot, valRewardsObj.Reward)
	}

}
