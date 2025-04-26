"""Protocol buffer modules"""

import os
import sys

# Add proto directory to path
proto_dir = os.path.dirname(os.path.abspath(__file__))
if proto_dir not in sys.path:
    sys.path.append(proto_dir)

# Import generated code
from . import scanner_pb2
from . import scanner_pb2_grpc

# Export key types
__all__ = [
    'scanner_pb2',
    'scanner_pb2_grpc',
    'ScannerServiceStub',
    'ScanTask',
    'ScanResult'
]