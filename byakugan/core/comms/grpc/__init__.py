"""gRPC communication package"""

# First load proto files
from byakugan.core.comms.grpc.proto import scanner_pb2, scanner_pb2_grpc

# Then load implementations
from byakugan.core.comms.grpc.client import ScannerClient
from byakugan.core.comms.grpc.server import ScannerServer

__all__ = [
    'ScannerClient',
    'ScannerServer',
    'scanner_pb2',
    'scanner_pb2_grpc'
]