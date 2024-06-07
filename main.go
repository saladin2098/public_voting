package main

import (
	"path/filepath"
	"runtime"
	cf "user-service/config"
	"user-service/config/logger"
	"user-service/db/postgresql"
	service "user-service/services"

	"net"

	pb "user-service/services/genproto"

	"google.golang.org/grpc"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	config := cf.Load()
	logger := logger.NewLogger(basepath, config.LOG_PATH)
	em := cf.NewErrorManager(logger)

	db, err := postgresql.ConnectDB(config)
	em.CheckErr(err)
	defer db.Close()

	listener, err := net.Listen("tcp", config.TCP_PORT)
	em.CheckErr(err)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, service.NewUserService(db))
	logger.INFO.Printf("Starting server on port %v", listener.Addr())
	em.CheckErr(err)

	if err := s.Serve(listener); err != nil {
		logger.ERROR.Panic(err.Error())
	}
}
