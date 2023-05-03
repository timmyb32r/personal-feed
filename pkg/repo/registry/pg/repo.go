package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
	"personal-feed/pkg/repo"
	"strings"
	"time"
)

type Repo struct {
	config *RepoConfigPG
	logger *logrus.Logger
	conn   *pgx.Conn
}

func (r *Repo) NewTx() (repo.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, xerrors.Errorf("unable to start tx, err: %w", err)
	}
	return tx, err
}

func (r *Repo) GetUserInfo(tx repo.Tx, userEmail string) (*model.User, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		context.Background(),
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

func (r *Repo) UpdateUserInfo(tx repo.Tx, userEmail string, user *model.User) error {
	query := `UPDATE users SET tg_chat_id=$1, nickname=$2, pass_hash=$3 WHERE email=$4;`
	unwrappedTx := tx.(pgx.Tx)
	_, err := unwrappedTx.Exec(context.Background(), query, user.TgChatID, user.Nickname, user.PassHash, userEmail)
	return err
}

func (r *Repo) ListSources() ([]model.Source, error) {
	rows, err := r.conn.Query(
		context.Background(),
		`SELECT id, description, crawler_id, crawler_meta, schedule, num_should_be_matched FROM source;`,
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select nodes: %w", err)
	}
	result := make([]model.Source, 0)
	for rows.Next() {
		var source model.Source
		err = rows.Scan(&source.ID, &source.Description, &source.CrawlerID, &source.CrawlerMeta, &source.Schedule, &source.NumShouldBeMatched)
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

func (r *Repo) InsertNewTreeNodes(tx repo.Tx, sourceID int, nodes []model.DBTreeNode) error {
	if len(nodes) == 0 {
		return nil
	}
	elems := make([]string, 0, len(nodes))
	args := make([]interface{}, 0, len(nodes)*4)
	index := 1
	for _, node := range nodes {
		elems = append(elems, fmt.Sprintf("($%d, $%d, $%d, $%d, now())", index, index+1, index+2, index+3))
		args = append(args, sourceID, node.Depth, node.ParentFullKey, node.CurrentNodeJSON)
		index += 4
	}
	query := `INSERT INTO events (source_id, depth, parent_full_key, current_node_json, insert_date) VALUES ` + strings.Join(elems, ",") + ";"
	unwrappedTx := tx.(pgx.Tx)
	_, err := unwrappedTx.Exec(context.Background(), query, args...)
	return err
}

func (r *Repo) ExtractTreeNodes(tx repo.Tx, sourceID int) ([]model.DBTreeNode, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		context.Background(),
		fmt.Sprintf(`SELECT depth, parent_full_key, current_node_json FROM events WHERE source_id=%d;`, sourceID),
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select nodes: %w", err)
	}
	result := make([]model.DBTreeNode, 0)
	for rows.Next() {
		var node model.DBTreeNode
		err = rows.Scan(&node.Depth, &node.ParentFullKey, &node.CurrentNodeJSON)
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

func (r *Repo) TestExtractAllTreeNodes(tx repo.Tx) ([]model.DBTreeNode, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		context.Background(),
		fmt.Sprintf(`SELECT depth, parent_full_key, current_node_json, insert_date FROM events ORDER BY insert_date DESC LIMIT 10;`),
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select nodes: %w", err)
	}
	result := make([]model.DBTreeNode, 0)
	for rows.Next() {
		var node model.DBTreeNode
		err = rows.Scan(&node.Depth, &node.ParentFullKey, &node.CurrentNodeJSON)
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

func (r *Repo) Close() error {
	return r.conn.Close(context.Background())
}

func NewRepo(cfg interface{}, logger *logrus.Logger) (repo.Repo, error) {
	cfgUnwrapped := cfg.(*RepoConfigPG) // unpack
	cfgUnwrappedCopy := *cfgUnwrapped
	conn, err := pgx.Connect(context.Background(), cfgUnwrapped.ToConnString())
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
