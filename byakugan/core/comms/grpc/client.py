"""gRPC client implementation"""
import grpc
import logging
from typing import Optional

from byakugan.config.comms import CommsConfig
from byakugan.core.comms.grpc.proto import scanner_pb2, scanner_pb2_grpc
from byakugan.core.comms.exceptions import ConnectionError

class ScannerClient:
    def __init__(self, config: CommsConfig):
        self.config = config
        self.channel = None
        self.stub = None
        self.logger = logging.getLogger("scanner_client")

    def connect(self) -> bool:
        """Establish gRPC connection"""
        try:
            address = f"{self.config.grpc_host}:{self.config.grpc_port}"
            self.logger.info(f"Connecting to gRPC server at {address}")
            
            opts = [
                ('grpc.max_receive_message_length', self.config.grpc_max_message_size),
                ('grpc.max_send_message_length', self.config.grpc_max_message_size),
            ]
            
            self.channel = grpc.insecure_channel(address, options=opts)
            grpc.channel_ready_future(self.channel).result(timeout=self.config.grpc_timeout)
            self.stub = scanner_pb2_grpc.ScannerServiceStub(self.channel)
            return True
            
        except Exception as e:
            self.logger.error(f"Connection failed: {str(e)}")
            return False

    def close(self):
        """Close client connection"""
        if self.channel:
            self.channel.close()
            self.channel = None

    def execute_scan(self, task: scanner_pb2.ScanTask) -> scanner_pb2.ScanResult:
        """Execute scan task"""
        if not self.stub:
            raise RuntimeError("Client not connected")
        
        try:
            return self.stub.ExecuteScan(task)
        except Exception as e:
            self.logger.error(f"Scan execution failed: {str(e)}")
            return scanner_pb2.ScanResult(
                task_id=task.id,
                success=False,
                findings=[],
                error_details=str(e)
            )