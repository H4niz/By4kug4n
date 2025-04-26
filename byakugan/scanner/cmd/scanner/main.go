package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/haniz/byakugan/scanner/config"
	"github.com/haniz/byakugan/scanner/engine"
	"github.com/haniz/byakugan/scanner/proto"
	"google.golang.org/grpc"
)

func main() {
	// Parse flags
	configPath := flag.String("config", "config/scanner.yaml", "Path to config file")
	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create engine
	scanEngine := engine.NewEngine(&engine.Config{
		WorkerCount: cfg.Engine.Workers,
		MaxRetries:  cfg.Engine.MaxRetries,
		Timeout:     cfg.Engine.Timeout,
	})

	// Create gRPC server
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxMessageSize),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxMessageSize),
	)

	// Register scanner service
	scannerService := NewScannerService(scanEngine, cfg)
	proto.RegisterScannerServiceServer(server, scannerService)

	// Start server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Handle shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal %v, initiating shutdown", sig)
		server.GracefulStop()
	}()

	log.Printf("Starting scanner service on port %d", cfg.GRPC.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type ScannerService struct {
	proto.UnimplementedScannerServiceServer
	engine *engine.Engine
	config *config.Config
}

func NewScannerService(engine *engine.Engine, cfg *config.Config) *ScannerService {
	return &ScannerService{
		engine: engine,
		config: cfg,
	}
}

// ExecuteScan handles scan requests
func (s *ScannerService) ExecuteScan(ctx context.Context, task *proto.ScanTask) (*proto.ScanResult, error) {
	// Validate task
	if err := s.validateTask(task); err != nil {
		return nil, err
	}

	// Create scan context
	scanCtx := &engine.ScanContext{
		ID:          task.Id,
		Target:      task.Target,
		RuleContext: task.RuleContext,
		Config:      task.Config,
	}

	// Execute scan
	findings, err := s.engine.ExecuteScan(ctx, scanCtx)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %v", err)
	}

	return &proto.ScanResult{
		TaskId:   task.Id,
		Success:  true,
		Findings: findings,
	}, nil
}

func (s *ScannerService) validateTask(task *proto.ScanTask) error {
	if task == nil {
		return fmt.Errorf("task is nil")
	}
	if task.Target == nil {
		return fmt.Errorf("target is nil")
	}
	if task.RuleContext == nil {
		return fmt.Errorf("rule context is nil")
	}
	if task.Target.Url == "" {
		return fmt.Errorf("target URL is empty")
	}
	return nil
}
