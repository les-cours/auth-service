package service

import (
	"github.com/gorilla/mux"
	"github.com/les-cours/auth-service/api/auth"
	"github.com/les-cours/auth-service/api/server"
	"github.com/les-cours/auth-service/api/users"
	"github.com/les-cours/auth-service/env"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
)

var (
	registry       = prometheus.NewRegistry()
	requestCounter = prometheus.NewGauge(prometheus.GaugeOpts{Name: "request_counter", Help: "request counter"})
	memoryUsage    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "memory_usage", Help: "memory usage"})
	goRoutineNum   = prometheus.NewGauge(prometheus.GaugeOpts{Name: "go_routines_num", Help: "the number of go routine "})
	cpuPercentage  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "cpu_percentage", Help: "cpu percentage"})
)

func monitoringMiddleware(originalHandler http.Handler) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memoryUsage.Set(float64(m.Alloc))
		goRoutineNum.Set(float64(runtime.NumGoroutine()))
		percent, _ := cpu.Percent(0, false)
		cpuPercentage.Set(percent[0])
		originalHandler.ServeHTTP(w, r)
	})
}

func loggerInit() *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	return logger
}

func Start() {

	logger := loggerInit()
	defer logger.Sync()

	registry.MustRegister(requestCounter, memoryUsage, cpuPercentage, goRoutineNum)
	router := mux.NewRouter()

	db, err := StartDB()
	if err != nil {
		log.Printf("Failed to start database, error: %v", err)
		return
	}

	defer db.Close()

	userConnectionService, err := grpc.Dial(env.Settings.UserService.Host+":"+env.Settings.UserService.Port, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Connection to agent domain service faild ", err)
	}

	userServiceClient := users.NewUserServiceClient(userConnectionService)
	s := server.GetInstance(userServiceClient, db, logger)

	router.HandleFunc("/login", cors(s.LoginHandler))
	router.HandleFunc("/teacher-login", cors(s.LoginTeacherHandler))
	router.HandleFunc("/login-admin", cors(s.LoginAdminHandler))
	router.HandleFunc("/token-health", cors(s.TokenHealthHandler))
	router.HandleFunc("/logout", cors(s.LogoutHandler))
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	router.HandleFunc("/metrics", cors(monitoringMiddleware(promHandler)))
	log.Println("Auth API listen at :" + env.Settings.HTTPPort)
	go func() {
		log.Fatal(http.ListenAndServe(":"+env.Settings.HTTPPort, router))
	}()

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, s)
	lis, err := net.Listen("tcp", ":"+env.Settings.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen on port %v: %v", env.Settings.GrpcPort, err)
	}
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server on port %v: %v", env.Settings.GrpcPort, err)
	}
}
