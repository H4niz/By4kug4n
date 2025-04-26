package scanner

import (
	"context"
	"log"
	"time"

	"github.com/avast/retry-go"
	"github.com/haniz/byakugan/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type ScannerService struct {
	pb.UnimplementedScannerServiceServer
	engine *Engine
	config *Config
}

func NewScannerService(config *Config) *ScannerService {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(config.MaxMessageSize),
		grpc.MaxSendMsgSize(config.MaxMessageSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 30 * time.Second,
			Time:              30 * time.Second,
			Timeout:           10 * time.Second,
		}),
	}

	if config.TLSEnabled {
		creds, err := credentials.NewServerTLSFromFile(
			config.TLSCertPath,
			config.TLSKeyPath,
		)
		if err != nil {
			log.Fatalf("Failed to load TLS credentials: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	return &ScannerService{
		server: grpc.NewServer(opts...),
		config: config,
		engine: NewEngine(config),
	}
}

func (s *ScannerService) ExecuteScan(ctx context.Context, task *pb.ScanTask) (*pb.ScanResult, error) {
	// Validate task
	if err := s.validateTask(task); err != nil {
		return nil, err
	}

	// Create scan context
	scanCtx := &ScanContext{
		ID:          task.Id,
		Target:      task.Target,
		RuleContext: task.RuleContext,
		Config:      task.Config,
	}

	// Execute scan
	findings, err := s.engine.ExecuteScan(ctx, scanCtx)
	if err != nil {
		return nil, err
	}

	// Process and return results
	return &pb.ScanResult{
		TaskId:   task.Id,
		Success:  true,
		Findings: findings,
	}, nil
}

func (s *ScannerService) StreamResults(stream pb.ScannerService_StreamResultsServer) error {
	for {
		result, err := stream.Recv()
		if err != nil {
			return err
		}

		// Process result
		processed := s.engine.ProcessResult(result)

		if err := stream.Send(processed); err != nil {
			return err
		}
	}
}
