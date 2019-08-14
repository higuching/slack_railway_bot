package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const TABLE_NAME = "TB_RAILWAYS"

type TbRailWays struct {
	ID   int
	NAME string
}

type RailwaysDb struct {
	DB *sql.DB
}

func NewRailways() *RailwaysDb {
	// データベースのコネクションを開く
	db, err := sql.Open("sqlite3", "/tmp/railways.db")
	if err != nil {
		panic(err)
	}
	o := RailwaysDb{}
	o.DB = db
	return &o
}

// テーブル作成
func (r *RailwaysDb) Create() error {
	_, err := r.DB.Exec(
		`CREATE TABLE IF NOT EXISTS "` + TABLE_NAME + `" ("ID" INTEGER PRIMARY KEY, "NAME" VARCHAR(255))`,
	)
	if err != nil {
		return err
	}
	return nil
}

// レコード1件取得(存在していたらTRUEを返す)
func (r *RailwaysDb) Get(id int) bool {
	row := r.DB.QueryRow(
		`SELECT ID FROM `+TABLE_NAME+` WHERE ID=?`,
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

// レコード全件取得
func (r *RailwaysDb) GetAll() []TbRailWays {
	rows, err := r.DB.Query(
		`SELECT ID, NAME FROM ` + TABLE_NAME,
	)
	if err != nil {
		panic(err)
	}

	var rs []TbRailWays
	defer rows.Close()
	for rows.Next() {
		var r TbRailWays
		if err := rows.Scan(&r.ID, &r.NAME); err != nil {
			panic(err)
		}
		rs = append(rs, r)
	}
	return rs
}

// レコード作成
func (r *RailwaysDb) Insert(id int, name string) error {
	_, err := r.DB.Exec(
		`INSERT INTO "`+TABLE_NAME+`" ("ID", "NAME") VALUES(?, ?)`,
		id,
		name,
	)
	if err != nil {
		return err
	}
	return nil
}

// レコード削除
func (r *RailwaysDb) Delete(id int) error {
	_, err := r.DB.Exec(
		`DELETE FROM "`+TABLE_NAME+`" WHERE ID=?`,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

// レコード全削除
func (r *RailwaysDb) DeleteAll() error {
	_, err := r.DB.Exec(
		`DELETE FROM "` + TABLE_NAME + `"`,
	)
	if err != nil {
		return err
	}
	return nil
}
