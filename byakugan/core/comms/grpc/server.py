from concurrent import futures
import grpc
from typing import Optional
import time

from ..config import CommsConfig
from ..exceptions import ConnectionError, SecurityError
from . import scanner_pb2_grpc, scanner_pb2
from .servicer import ScannerServicer

class ScannerServer:
    """gRPC server for scanner nodes"""
    
    def __init__(self, config: CommsConfig):
        self.config = config
        self._server: Optional[grpc.Server] = None
        self._servicer: Optional[ScannerServicer] = None
        
    def start(self) -> bool:
        """Start gRPC server"""
        try:
            # Create server
            self._server = grpc.server(
                futures.ThreadPoolExecutor(
                    max_workers=self.config.grpc_max_workers
                ),
                options=[
                    ('grpc.max_message_length', 
                     self.config.grpc_max_message_size)
                ]
            )

            # Create and add servicer
            self._servicer = ScannerServicer()
            scanner_pb2_grpc.add_ScannerServiceServicer_to_server(
                self._servicer, 
                self._server
            )

            # Configure security
            if self.config.tls_enabled:
                self._configure_security()
            else:
                self._server.add_insecure_port(
                    f"{self.config.grpc_host}:{self.config.grpc_port}"
                )

            # Start server
            self._server.start()
            return True
            
        except Exception as e:
            if self._server:
                self.stop()
            raise ConnectionError(f"Failed to start server: {str(e)}")
            
    def stop(self):
        """Stop gRPC server"""
        if self._server:
            self._server.stop(grace=5)
            self._server = None

    def _configure_security(self):
        """Configure TLS security"""
        if not (self.config.tls_cert_path and self.config.tls_key_path):
            raise SecurityError("TLS enabled but cert/key paths not provided")
            
        try:
            with open(self.config.tls_cert_path, 'rb') as f:
                cert = f.read()
            with open(self.config.tls_key_path, 'rb') as f:
                key = f.read()
                
            creds = grpc.ssl_server_credentials([(key, cert)])
            self._server.add_secure_port(
                f"{self.config.grpc_host}:{self.config.grpc_port}", 
                creds
            )
        except Exception as e:
            raise SecurityError(f"Failed to configure TLS: {str(e)}")

    def execute_scan(self, task: scanner_pb2.ScanTask) -> scanner_pb2.ScanResult:
        """Execute scan task"""
        if not self._servicer:
            return scanner_pb2.ScanResult()
        
        return self._servicer.ExecuteScan(task, None)