"""Protocol buffer modules"""
from pathlib import Path
import sys

# Add proto directory to path
proto_dir = Path(__file__).parent
if str(proto_dir) not in sys.path:
    sys.path.insert(0, str(proto_dir))

from .scanner_pb2 import *  # noqa
from .scanner_pb2_grpc import *  # noqa

__all__ = [
    'scanner_pb2',
    'scanner_pb2_grpc'
]