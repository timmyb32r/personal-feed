package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
	"personal-feed/pkg/util"
	"strings"
	"time"
)

type Repo struct {
	config *RepoConfigPG
	logger *logrus.Logger
	conn   *pgx.Conn
}

func (r *Repo) GenerateLiquibaseProperties() (string, error) {
	result := ""
	result += fmt.Sprintf("changeLogFile:dbchangelog.yml\n")
	result += fmt.Sprintf("url: jdbc:postgresql://%s:%d/%s\n", r.config.Host, r.config.Port, r.config.Name)
	result += fmt.Sprintf("username: %s\n", r.config.User)
	result += fmt.Sprintf("password: %s\n", r.config.Password)
	return result, nil
}

func (r *Repo) NewTx(ctx context.Context) (repo.Tx, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, xerrors.Errorf("unable to start tx, err: %w", err)
	}
	return tx, err
}

func (r *Repo) GetUserInfo(tx repo.Tx, ctx context.Context, userEmail string) (*model.User, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		ctx,
		`SELECT id, email, tg_chat_id, nickname, pass_hash FROM users WHERE email=$1;`,
		userEmail,
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select user info: %w", err)
	}
	var user *model.User = nil
	for rows.Next() {
		user = new(model.User)
		err = rows.Scan(&user.ID, &user.Email, &user.TgChatID, &user.Nickname, &user.PassHash)
		if err != nil {
			return nil, xerrors.Errorf("unable to scan, err: ", err)
		}
	}
	if rows.Err() != nil {
		return nil, xerrors.Errorf("got some error during reading, err: %w", err)
	}
	return user, nil
}

func (r *Repo) UpdateUserInfo(tx repo.Tx, ctx context.Context, userEmail string, user *model.User) error {
	query := `UPDATE users SET tg_chat_id=$1, nickname=$2, pass_hash=$3 WHERE email=$4;`
	unwrappedTx := tx.(pgx.Tx)
	_, err := unwrappedTx.Exec(ctx, query, user.TgChatID, user.Nickname, user.PassHash, userEmail)
	return err
}

func (r *Repo) ListSources(ctx context.Context) ([]model.Source, error) {
	rows, err := r.conn.Query(
		ctx,
		`SELECT id, description, crawler_id, crawler_meta, schedule, num_should_be_matched, history_state FROM source;`,
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select nodes: %w", err)
	}
	result := make([]model.Source, 0)
	for rows.Next() {
		var source model.Source
		err = rows.Scan(&source.ID, &source.Description, &source.CrawlerID, &source.CrawlerMeta, &source.Schedule, &source.NumShouldBeMatched, &source.HistoryState)
		if err != nil {
			return nil, xerrors.Errorf("unable to scan, err: ", err)
		}
		result = append(result, source)
	}
	if rows.Err() != nil {
		return nil, xerrors.Errorf("got some error during reading, err: %w", err)
	}
	return result, nil
}

func (r *Repo) InsertNewTreeNodes(ctx context.Context, sourceID int, nodes []model.DBTreeNode) error {
	rollbacks := util.Rollbacks{}
	defer rollbacks.Do()

	tx, err := r.NewTx(ctx)
	if err != nil {
		return xerrors.Errorf("unable to create transaction, err: %w", err)
	}

	rollbacks.Add(func() { _ = tx.Rollback(ctx) })

	err = r.InsertNewTreeNodesTx(tx, ctx, sourceID, nodes)
	if err != nil {
		return xerrors.Errorf("unable to insert nodes, err: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return xerrors.Errorf("unable to commit, err: %w", err)
	}

	rollbacks.Cancel()
	return nil
}

func (r *Repo) InsertNewTreeNodesTx(tx repo.Tx, ctx context.Context, sourceID int, nodes []model.DBTreeNode) error {
	if len(nodes) == 0 {
		return nil
	}
	eventsElems := make([]string, 0, len(nodes))
	eventsArgs := make([]interface{}, 0, len(nodes)*5)
	eventsIndex := 1

	docsElems := make([]string, 0, len(nodes))
	docsArgs := make([]interface{}, 0, len(nodes)*5)
	docsIndex := 1

	for _, node := range nodes {
		currentNodeJSON := node.CurrentNodeJSON
		if node.IsDoc {
			currentNodeJSON = "{}"
		}

		eventsElems = append(eventsElems, fmt.Sprintf("($%d, $%d, $%d, $%d, now(), $%d)", eventsIndex, eventsIndex+1, eventsIndex+2, eventsIndex+3, eventsIndex+4))
		eventsArgs = append(eventsArgs, sourceID, node.Depth, node.CurrentFullKey, currentNodeJSON, node.BusinessTime)
		eventsIndex += 5

		docsElems = append(docsElems, fmt.Sprintf("($%d, $%d)", docsIndex, docsIndex+1))
		docsArgs = append(docsArgs, node.CurrentFullKey, node.CurrentNodeJSON)
		docsIndex += 2
	}

	var query string
	var err error
	unwrappedTx := tx.(pgx.Tx)

	query = `INSERT INTO events (source_id, depth, current_full_key, current_node_json, insert_date, business_time) VALUES ` + strings.Join(eventsElems, ",") + ";"
	_, err = unwrappedTx.Exec(ctx, query, eventsArgs...)
	if err != nil {
		return xerrors.Errorf("unable to exec query, query:%s, err:%w", query, err)
	}

	query = `INSERT INTO events_doc (current_full_key, body) VALUES ` + strings.Join(docsElems, ",") + ";"
	_, err = unwrappedTx.Exec(ctx, query, docsArgs...)
	if err != nil {
		return xerrors.Errorf("unable to exec query, query:%s, err:%w", query, err)
	}

	return err
}

func (r *Repo) ExtractTreeNodes(ctx context.Context, sourceID int) ([]model.DBTreeNode, error) {
	rollbacks := util.Rollbacks{}
	defer rollbacks.Do()

	tx, err := r.NewTx(ctx)
	if err != nil {
		return nil, xerrors.Errorf("unable to create transaction, err: %w", err)
	}

	rollbacks.Add(func() { _ = tx.Rollback(ctx) })

	result, err := r.ExtractTreeNodesTx(tx, ctx, sourceID)
	if err != nil {
		return nil, xerrors.Errorf("unable to extract nodes, err: %w", err)
	}

	return result, nil
}

func (r *Repo) ExtractTreeNodesTx(tx repo.Tx, ctx context.Context, sourceID int) ([]model.DBTreeNode, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		ctx,
		fmt.Sprintf(`SELECT depth, current_full_key, current_node_json FROM events WHERE source_id=%d;`, sourceID),
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select nodes: %w", err)
	}
	result := make([]model.DBTreeNode, 0)
	for rows.Next() {
		var node model.DBTreeNode
		err = rows.Scan(&node.Depth, &node.CurrentFullKey, &node.CurrentNodeJSON)
		if err != nil {
			return nil, xerrors.Errorf("unable to scan, err: ", err)
		}
		result = append(result, node)
	}
	if rows.Err() != nil {
		return nil, xerrors.Errorf("got some error during reading, err: %w", err)
	}
	return result, nil
}

func (r *Repo) GetNextCronPeriod(ctx context.Context) (lastRunTime *time.Time, currentTime time.Time, err error) {
	query := fmt.Sprintf(
		"SELECT (SELECT last_run_time FROM %s.cron), now() AT TIME ZONE 'utc'",
		r.config.Schema)
	row := r.conn.QueryRow(ctx, query)
	if err := row.Scan(&lastRunTime, &currentTime); err != nil {
		return nil, time.Time{}, xerrors.Errorf("unable to get last_run_time and now(), err: %w", err)
	}
	return lastRunTime, currentTime, nil
}

func (r *Repo) SetCronLastRunTime(ctx context.Context, cronLastRunTime time.Time) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.cron(last_run_time) VALUES ($1) ON CONFLICT(id) DO UPDATE SET last_run_time = $1",
		r.config.Schema)
	if _, err := r.conn.Exec(ctx, query, cronLastRunTime); err != nil {
		return xerrors.Errorf("unable to set cron last_run_time, err: %w", err)
	}
	return nil
}

func (r *Repo) SetState(ctx context.Context, sourceID int, state string) error {
	query := `UPDATE source SET history_state=$1 WHERE id=$2;`
	res, err := r.conn.Exec(ctx, query, state, sourceID)
	if err != nil {
		return xerrors.Errorf("unable to update state, err: %w", err)
	}
	if res.RowsAffected() != int64(1) {
		return xerrors.Errorf("SetState updated wrong amount of rows: %d", res.RowsAffected())
	}
	return nil
}

func (r *Repo) InsertSourceIterationTx(tx repo.Tx, ctx context.Context, sourceID int, link, body string) error {
	query := fmt.Sprintf(
		"INSERT INTO %s.events_iteration(source_id, insert_timestamp, link, body) VALUES ($1, now(), $2, $3)",
		r.config.Schema)
	unwrappedTx := tx.(pgx.Tx)
	if _, err := unwrappedTx.Exec(ctx, query, sourceID, link, body); err != nil {
		return xerrors.Errorf("unable to insert event into events_iteration, err: %w", err)
	}
	return nil
}

func (r *Repo) InsertSourceIteration(ctx context.Context, sourceID int, link, body string) error {
	rollbacks := util.Rollbacks{}
	defer rollbacks.Do()

	tx, err := r.NewTx(ctx)
	if err != nil {
		return xerrors.Errorf("unable to create transaction, err: %w", err)
	}

	rollbacks.Add(func() { _ = tx.Rollback(ctx) })

	err = r.InsertSourceIterationTx(tx, ctx, sourceID, link, body)
	if err != nil {
		return xerrors.Errorf("unable to insert source iteration in tx, err: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return xerrors.Errorf("unable to commit tx, err: %w", err)
	}

	rollbacks.Cancel()
	return nil
}

func (r *Repo) TestExtractAllTreeNodes(tx repo.Tx, ctx context.Context) ([]model.DBTreeNode, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		ctx,
		fmt.Sprintf(`SELECT e.depth, e.current_full_key, d.body FROM events e INNER JOIN events_doc d ON e.current_full_key=d.current_full_key WHERE e.depth=2 ORDER BY e.business_time DESC LIMIT 10;`),
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select nodes: %w", err)
	}
	result := make([]model.DBTreeNode, 0)
	for rows.Next() {
		var node model.DBTreeNode
		err = rows.Scan(&node.Depth, &node.CurrentFullKey, &node.CurrentNodeJSON)
		if err != nil {
			return nil, xerrors.Errorf("unable to scan, err: ", err)
		}
		result = append(result, node)
	}
	if rows.Err() != nil {
		return nil, xerrors.Errorf("got some error during reading, err: %w", err)
	}
	return result, nil
}

func (r *Repo) Close(ctx context.Context) error { // absent in abstract.go
	return r.conn.Close(ctx)
}

func NewRepo(ctx context.Context, cfg interface{}, logger *logrus.Logger) (repo.Repo, error) {
	cfgUnwrapped := cfg.(*RepoConfigPG) // unpack
	cfgUnwrappedCopy := *cfgUnwrapped
	conn, err := pgx.Connect(ctx, cfgUnwrapped.ToConnString())
	if err != nil {
		return nil, xerrors.Errorf("unable to connect to the database, err: %w", err)
	}
	return &Repo{
		config: &cfgUnwrappedCopy,
		logger: logger,
		conn:   conn,
	}, nil
}

func init() {
	repo.Register(NewRepo, &RepoConfigPG{})
}
