package main

import (
	"database/sql"
	"github.com/carloseduribeiro/working-with-grpc-go/internal/database"
	"github.com/carloseduribeiro/working-with-grpc-go/internal/pb"
	"github.com/carloseduribeiro/working-with-grpc-go/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if _, err := db.Exec(`create table if not exists category(id string, name string, description string)`); err != nil {
		panic(err)
	}
	if _, err := db.Exec(`create table if not exists course(id string, name string, description string, category_id string)`); err != nil {
		panic(err)
	}

	categoryDB := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDB)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
