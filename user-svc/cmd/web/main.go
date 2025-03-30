package main

import (
	"be-yourmoments/user-svc/internal/adapter"
	"be-yourmoments/user-svc/internal/config"
	http "be-yourmoments/user-svc/internal/delivery/http/controller"
	"be-yourmoments/user-svc/internal/delivery/http/middleware"
	"be-yourmoments/user-svc/internal/delivery/http/route"
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/helper/discovery"
	"log"
	"os"
	"os/signal"
	"syscall"

	"be-yourmoments/user-svc/internal/helper/discovery/consul"
	"be-yourmoments/user-svc/internal/helper/logger"
	"be-yourmoments/user-svc/internal/repository"
	"be-yourmoments/user-svc/internal/usecase"

	"context"
	"fmt"
	"strconv"
	"time"
)

var logs = logger.New("main")

func webServer() error {
	app := config.NewApp()

	serverConfig := config.NewServerConfig()
	dbConfig := config.NewDB()
	// minioConfig := config.NewMinio()
	redisConfig := config.NewRedisClient()

	registry, err := consul.NewRegistry(serverConfig.ConsulAddr, serverConfig.Name)
	if err != nil {
		logs.Error("Failed to create consul registry for category service" + err.Error())
		return err
	}

	// GRPCserviceID := discovery.GenerateServiceID(serverConfig.Name + "-grpc")
	HTTPserviceID := discovery.GenerateServiceID(serverConfig.Name + "-http")

	// grpcPortInt, _ := strconv.Atoi(serverConfig.GRPCPort)
	httpPortInt, _ := strconv.Atoi(serverConfig.HTTPPort)

	ctx := context.Background()

	// err = registry.RegisterService(ctx, serverConfig.Name+"-grpc", GRPCserviceID, serverConfig.GRPCAddr, grpcPortInt, []string{"grpc"})
	// if err != nil {
	// 	logs.Error("Failed to register gRPC book service to consul")
	// 	return err
	// }

	err = registry.RegisterService(ctx, serverConfig.Name+"-http", HTTPserviceID, serverConfig.HTTPAddr, httpPortInt, []string{"http"})
	if err != nil {
		logs.Error("Failed to register category service to consuls")
		return err
	}

	// go func() {
	// 	failureCount := 0
	// 	const maxFailures = 5
	// 	for {
	// 		err := registry.HealthCheck(GRPCserviceID, serverConfig.Name+"-grpc")
	// 		if err != nil {
	// 			logs.Error(fmt.Sprintf("Failed to perform health check for gRPC service: %v", err))
	// 			failureCount++
	// 			if failureCount >= maxFailures {
	// 				logs.Error("Max health check failures reached for gRPC service. Exiting health check loop.")
	// 				break
	// 			}
	// 		} else {
	// 			failureCount = 0
	// 		}
	// 		time.Sleep(time.Second * 2)
	// 	}
	// }()
	// defer registry.DeregisterService(ctx, GRPCserviceID)

	go func() {
		failureCount := 0
		const maxFailures = 5
		for {
			err := registry.HealthCheck(HTTPserviceID, serverConfig.Name)
			if err != nil {
				logs.Error(fmt.Sprintf("Failed to perform health check: %v", err))
				failureCount++
				if failureCount >= maxFailures {
					logs.Error("Max health check failures reached for HTTP service. Exiting health check loop.")
					break
				}
			} else {
				failureCount = 0
			}
			time.Sleep(time.Second * 2)
		}
	}()
	defer registry.DeregisterService(ctx, HTTPserviceID)

	cacheAdapter := adapter.NewCacheAdapter(redisConfig)
	emailAdapter := adapter.NewEmailAdapter()
	googleTokenAdapter := adapter.NewGoogleTokenAdapter()
	jwtAdapter := adapter.NewJWTAdapter()
	securityAdapter := adapter.NewSecurityAdapter()
	customValidator := helper.NewCustomValidator()

	logs.Log(fmt.Sprintf("Succsess connected http service at port: %v", serverConfig.HTTP))

	userRepository, err := repository.NewUserRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	userProfileRepository, err := repository.NewUserProfileRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	emailVerificationRepository, err := repository.NewEmailVerificationRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	resetPasswordRepository, err := repository.NewResetPasswordRepository(dbConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	authUseCase := usecase.NewAuthUseCase(dbConfig, userRepository, userProfileRepository, emailVerificationRepository, resetPasswordRepository,
		googleTokenAdapter, emailAdapter, jwtAdapter, securityAdapter, cacheAdapter)

	authController := http.NewAuthController(authUseCase, customValidator)

	authMiddleware := middleware.NewUserAuth(authUseCase, customValidator)

	routeConfig := route.RouteConfig{
		App:            app,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.Setup()

	// go func() {
	// 	grpcServer := grpc.NewServer()
	// 	reflection.Register(grpcServer)

	// 	l, err := net.Listen("tcp", serverConfig.GRPC)
	// 	if err != nil {
	// 		logs.Error(fmt.Sprintf("Failed to listen: %v", err))
	// 	}
	// 	logs.Log(fmt.Sprintf("gRPC server started on %s", serverConfig.GRPC))
	// 	defer l.Close()

	// 	// grpcHandler.NewPhotoGRPCHandler(grpcServer, photoUsecase, faceCamUseCase, userSimilarPhotoUsecase)

	// 	if err := grpcServer.Serve(l); err != nil {
	// 		logs.Error(fmt.Sprintf("Failed to start gRPC category server: %v", err))
	// 	}
	// }()

	// photoController.Route(app)
	logs.Log(fmt.Sprintf("Succsess connected http service at port: %v", serverConfig.HTTP))

	err = app.Listen(serverConfig.HTTP)

	if err != nil {
		logs.Error(fmt.Sprintf("Failed to start HTTP category server: %v", err))
		return err
	}
	return nil
}

func main() {
	if err := webServer(); err != nil {
		logs.Error(err)
	}

	logs.Log("Api gateway server started")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigchan
	logs.Log(fmt.Sprintf("Received signal: %s. Shutting down gracefully...", sig))
}
