package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/xerrors"
	"personal-feed/pkg/model"
	"strings"
	"time"
)

type PgClient struct {
	conn *pgx.Conn
}

func (c *PgClient) NewTx() (model.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tx, err := c.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite})
	if err != nil {
		return nil, xerrors.Errorf("unable to start tx, err: %w", err)
	}
	return tx, err
}

func (c *PgClient) GetUserInfo(tx model.Tx, userEmail string) (*model.User, error) {
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

func (c *PgClient) UpdateUserInfo(tx model.Tx, userEmail string, user *model.User) error {
	query := `UPDATE users SET tg_chat_id=$1, nickname=$2, pass_hash=$3 WHERE email=$4;`
	unwrappedTx := tx.(pgx.Tx)
	_, err := unwrappedTx.Exec(context.Background(), query, user.TgChatID, user.Nickname, user.PassHash, userEmail)
	return err
}

func (c *PgClient) ListSources(tx model.Tx) ([]model.Source, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		context.Background(),
		`SELECT id, description, crawler_id, crawler_meta, schedule FROM source;`,
	)
	if err != nil {
		return nil, xerrors.Errorf("unable to select nodes: %w", err)
	}
	result := make([]model.Source, 0)
	for rows.Next() {
		var source model.Source
		err = rows.Scan(&source.ID, &source.Description, &source.CrawlerID, &source.CrawlerMeta, &source.Schedule)
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

func (c *PgClient) InsertNewTreeNodes(tx model.Tx, sourceID int, nodes []model.DBTreeNode) error {
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

func (c *PgClient) ExtractTreeNodes(tx model.Tx, sourceID int) ([]model.DBTreeNode, error) {
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

func (c *PgClient) TestExtractAllTreeNodes(tx model.Tx) ([]model.DBTreeNode, error) {
	unwrappedTx := tx.(pgx.Tx)
	rows, err := unwrappedTx.Query(
		context.Background(),
		fmt.Sprintf(`SELECT depth, parent_full_key, current_node_json FROM events LIMIT 10;`),
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

func (c *PgClient) Close() error {
	return c.conn.Close(context.Background())
}

func NewPgClient(cfg *config) (*PgClient, error) {
	conn, err := pgx.Connect(context.Background(), cfg.ToConnString())
	if err != nil {
		return nil, xerrors.Errorf("unable to connect to the database, err: %w", err)
	}
	return &PgClient{
		conn: conn,
	}, nil
}
