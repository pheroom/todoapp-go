package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/k0kubun/pp/v3"
	"net/http"
	"os"
	"time"
)

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	urlExample := "postgres://admin:admin@localhost:5432/testapidb"
	pp.Println(os.Getenv("DATABASE_URL"))
	return pgx.Connect(ctx, urlExample)
}

type TaskModel struct {
	ID          int
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func CreateTable(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
	    id SERIAL PRIMARY KEY,
	    title VARCHAR(255) NOT NULL,
	    completed BOOLEAN NOT NULL,
	    created_at TIMESTAMP NOT NULL,
	    completed_at TIMESTAMP
	)`

	_, err := conn.Exec(ctx, sqlQuery)
	return err
}

func InsertRaw(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	INSERT INTO tasks (title, completed, created_at)
	VALUES ($1, $2, $3)`

	_, err := conn.Exec(ctx, sqlQuery, "новая заметка2", false, time.Now())
	return err
}

func UpdateRaw(ctx context.Context, conn *pgx.Conn, id int, model TaskModel) error {
	sqlQuery := `
	UPDATE tasks 
	SET title = $2, completed = $4, created_at = $3, completed_at = $5
	WHERE id = $1`

	_, err := conn.Exec(ctx, sqlQuery, id, model.Title, model.Completed, model.CreatedAt, model.CompletedAt)
	return err
}

func GetRaws(ctx context.Context, conn *pgx.Conn) ([]TaskModel, error) {
	sqlQuery := `
	select  id, title, completed, created_at, completed_at
	from tasks
	order by completed_at desc 
	limit 100`

	rows, err := conn.Query(ctx, sqlQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	res := []TaskModel{}
	for rows.Next() {
		var task TaskModel
		rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.CompletedAt)
		res = append(res, task)
	}
	return res, nil
}

func main() {
	//ctx := context.Background()
	//conn, err := CreateConnection(ctx)
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close(ctx)
	//
	//if err := CreateTable(ctx, conn); err != nil {
	//	panic(err)
	//}
	//
	//if err := InsertRaw(ctx, conn); err != nil {
	//	panic(err)
	//}
	//
	//if err := UpdateRaw(ctx, conn); err != nil {
	//	panic(err)
	//}
	//
	//res, err := GetRaws(ctx, conn)
	//if err != nil {
	//	panic(err)
	//}
	//pp.Println(res)

	fmt.Println(os.Getpid())

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	fmt.Println("Server starting")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		fmt.Println("listen err:", err)
	}
}
