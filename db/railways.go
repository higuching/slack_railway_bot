package db

import (
	"database/sql"
	"fmt"

	// go-sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// TableRailwaysName テーブル名
const TableRailwaysName = "TB_RAILWAYS"

// TRailWays スキーマ
type TRailWays struct {
	ID   int
	NAME string
}

// RailwaysDb インスタンス的なもの
type RailwaysDb struct {
	DB *sql.DB
}

// NewRailways データベースのコネクションを開く
func NewRailways() *RailwaysDb {
	db, err := sql.Open("sqlite3", "/tmp/railways.db")
	if err != nil {
		panic(err)
	}
	o := RailwaysDb{}
	o.DB = db
	return &o
}

// Create テーブル作成
func (r *RailwaysDb) Create() error {
	_, err := r.DB.Exec(
		`CREATE TABLE IF NOT EXISTS "` + TableRailwaysName + `" ("ID" INTEGER PRIMARY KEY, "NAME" VARCHAR(255))`,
	)
	if err != nil {
		return err
	}
	return nil
}

// Get レコード1件取得(存在していたらTRUEを返す)
func (r *RailwaysDb) Get(id int) bool {
	row := r.DB.QueryRow(
		`SELECT ID FROM `+TableRailwaysName+` WHERE ID=?`,
		id,
	)

	err := row.Scan(&id)

	// エラー内容による分岐
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("Not found: %d\n", id)
	case err != nil:
		panic(err)
	default:
		fmt.Printf("Finded!: %d\n", id)
		return true
	}
	return false
}

// GetAll レコード全件取得
func (r *RailwaysDb) GetAll() []TRailWays {
	rows, err := r.DB.Query(
		`SELECT ID, NAME FROM ` + TableRailwaysName,
	)
	if err != nil {
		panic(err)
	}

	var rs []TRailWays
	defer rows.Close()
	for rows.Next() {
		var r TRailWays
		if err := rows.Scan(&r.ID, &r.NAME); err != nil {
			panic(err)
		}
		rs = append(rs, r)
	}
	return rs
}

// Insert レコード作成
func (r *RailwaysDb) Insert(id int, name string) error {
	_, err := r.DB.Exec(
		`INSERT INTO "`+TableRailwaysName+`" ("ID", "NAME") VALUES(?, ?)`,
		id,
		name,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete レコード削除
func (r *RailwaysDb) Delete(id int) error {
	_, err := r.DB.Exec(
		`DELETE FROM "`+TableRailwaysName+`" WHERE ID=?`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll レコード全削除
func (r *RailwaysDb) DeleteAll() error {
	_, err := r.DB.Exec(
		`DELETE FROM "` + TableRailwaysName + `"`,
	)
	if err != nil {
		return err
	}
	return nil
}
