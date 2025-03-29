package adapterout

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/portout"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

func NewRuleRepository(brk *broker.NatsMessageBroker) (*RuleRepositoryImpl, error) {
	kv, err := brk.Js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:  "notifications",
		History: 1,
	})
	if err != nil {
		return nil, err
	}
	return &RuleRepositoryImpl{kv: kv}, nil
}

type RuleRepositoryImpl struct {
	kv jetstream.KeyValue
}

var _ portout.RuleRepository = (*RuleRepositoryImpl)(nil)

func (r RuleRepositoryImpl) AddRule(data types.QueryRule) (uuid.UUID, error) {
	id := uuid.New()
	bytes, err := json.Marshal(data)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = r.kv.Create(context.Background(), id.String(), bytes)
	if err != nil {
		return uuid.Nil, types.ErrRuleExists
	}

	return id, nil
}

func (r RuleRepositoryImpl) GetRule(id uuid.UUID) (types.QueryRule, error) {
	get, err := r.kv.Get(context.Background(), id.String())
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		return types.QueryRule{}, types.ErrRuleNotExists
	}
	if err != nil {
		return types.QueryRule{}, err
	}

	var res types.QueryRule
	err = json.Unmarshal(get.Value(), &res)
	return res, err
}

func (r RuleRepositoryImpl) ListRules() ([]types.QueryRuleWithId, error) {
	keys, err := r.kv.Keys(context.Background())
	if err != nil {
		return nil, err
	}

	// TODO: implementare meglio, magari leggendo direttamente i messaggi dallo stream
	var res []types.QueryRuleWithId
	for _, key := range keys {
		bytes, err := r.kv.Get(context.Background(), key)
		if err != nil {
			return nil, err
		}

		var val types.QueryRuleWithId
		err = json.Unmarshal(bytes.Value(), &val)
		if err != nil {
			return nil, err
		}
		val.RuleId, err = uuid.Parse(key)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}

	return res, nil
}

func (r RuleRepositoryImpl) EditRule(id uuid.UUID, data types.EditRule) error {
	bytes, err := r.kv.Get(context.Background(), id.String())
	if err != nil {
		return err
	}

	var val types.QueryRule
	err = json.Unmarshal(bytes.Value(), &val)
	if err != nil {
		return err
	}

	if data.Operator != nil {
		val.Operator = *data.Operator
	}
	if data.GoodId != nil {
		val.GoodId = *data.GoodId
	}
	if data.Threshold != nil {
		val.Threshold = *data.Threshold
	}

	serialized, err := json.Marshal(val)
	if err != nil {
		return err
	}

	_, err = r.kv.Put(context.Background(), id.String(), serialized)
	return err
}

func (r RuleRepositoryImpl) RemoveRule(id uuid.UUID) error {
	return r.kv.Delete(context.Background(), id.String())
}
