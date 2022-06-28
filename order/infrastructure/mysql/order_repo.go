package mysql

import (
	"database/sql"

	"github.com/imylam/delivery-test/order"

	"github.com/jmoiron/sqlx"
)

type orderRepoMysql struct {
	MysqlConn *sqlx.DB
}

// NewOrderRepositoryMysql will create an object that represent the order.OrderRepository interface
func NewOrderRepositoryMysql(mysqlConn *sqlx.DB) order.OrderRepository {
	return &orderRepoMysql{mysqlConn}
}

func (repo *orderRepoMysql) Create(order *order.Order) error {
	q1 := "INSERT INTO orders (distance, status, created_at, updated_at) VALUES (?,?,now(),now())"
	q2 := "SELECT * FROM orders WHERE id=?"

	insertStmt, err := repo.MysqlConn.Prepare(q1)
	if err != nil {
		return err
	}

	result, err := insertStmt.Exec(order.Distance, order.Status)
	if err != nil {
		return err
	}

	order.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	err = repo.MysqlConn.QueryRowx(q2, order.ID).StructScan(order)
	if err != nil {
		return err
	}

	return nil
}

func (repo *orderRepoMysql) UpdateStatusByID(id int64) error {
	q := "UPDATE orders SET status=? WHERE id=? AND status=?"

	updateStmt, err := repo.MysqlConn.Prepare(q)
	if err != nil {
		return err
	}

	result, err := updateStmt.Exec(order.StatusTaken, id, order.StatusUnassigned)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (repo *orderRepoMysql) FindByID(id int64) (*order.Order, error) {
	q := "SELECT * FROM orders WHERE id=?"

	var order order.Order
	err := repo.MysqlConn.QueryRowx(q, id).StructScan(&order)
	if err != nil {
		return nil, err
	}

	return &order, err
}

func (repo *orderRepoMysql) FindRange(limit, offset int) (*[]order.Order, error) {
	q := "SELECT * FROM orders LIMIT ? OFFSET ?"

	rows, err := repo.MysqlConn.Queryx(q, limit, offset)
	if err != nil {
		return nil, err
	}

	var orders []order.Order
	for rows.Next() {
		var order order.Order
		err = rows.StructScan(&order)
		if err != nil {

			return nil, err
		}
		orders = append(orders, order)
	}

	return &orders, nil
}
