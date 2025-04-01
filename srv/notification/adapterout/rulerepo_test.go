package adapterout

import (
	"testing"

	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/notification/types"
	"github.com/stretchr/testify/require"
)

func TestRuleRepoAddRule(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	defer ns.Close()

	broker := broker.NewTest(t, ns)

	repo, err := NewRuleRepository(broker)
	require.NoError(t, err)
	require.NotNil(t, repo)

	uuid, err := repo.AddRule(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, uuid)
}

func TestRuleRepoGetRule(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	defer ns.Close()

	broker := broker.NewTest(t, ns)

	repo, err := NewRuleRepository(broker)
	require.NoError(t, err)
	require.NotNil(t, repo)

	uuid, err := repo.AddRule(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, uuid)

	rule, err := repo.GetRule(uuid)
	require.NoError(t, err)
	require.Equal(t, rule.GoodId, "1")
}

func TestRuleRepoListRules(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	defer ns.Close()

	broker := broker.NewTest(t, ns)

	repo, err := NewRuleRepository(broker)
	require.NoError(t, err)
	require.NotNil(t, repo)

	uuid, err := repo.AddRule(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, uuid)

	rules, err := repo.ListRules()
	require.NoError(t, err)
	require.Len(t, rules, 1)
	require.Equal(t, rules[0].GoodId, "1")
}

func TestRuleRepoRemoveRule(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	defer ns.Close()

	broker := broker.NewTest(t, ns)

	repo, err := NewRuleRepository(broker)
	require.NoError(t, err)
	require.NotNil(t, repo)

	uuid, err := repo.AddRule(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, uuid)

	err = repo.RemoveRule(uuid)
	require.NoError(t, err)

	rule, err := repo.GetRule(uuid)
	require.Error(t, types.ErrRuleNotExists, err)
	require.Equal(t, types.QueryRule{}, rule)
}

func TestRuleRepoEditRule(t *testing.T) {
	ns, _ := broker.NewInProcessNATSServer(t)
	defer ns.Close()

	broker := broker.NewTest(t, ns)

	repo, err := NewRuleRepository(broker)
	require.NoError(t, err)
	require.NotNil(t, repo)

	uuid, err := repo.AddRule(types.QueryRule{
		GoodId:    "1",
		Operator:  "<",
		Threshold: 10,
	})
	require.NoError(t, err)
	require.NotEmpty(t, uuid)

	goodId := "2"
	operator := "<="
	quantity := 10

	err = repo.EditRule(uuid, types.EditRule{
		GoodId:    &goodId,
		Operator:  &operator,
		Threshold: &quantity,
	})
	require.NoError(t, err)
}
