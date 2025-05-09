# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import scanner_pb2 as scanner__pb2


class ScannerServiceStub(object):
    """Scanner service definition
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.ExecuteScan = channel.unary_unary(
                '/scanner.ScannerService/ExecuteScan',
                request_serializer=scanner__pb2.ScanTask.SerializeToString,
                response_deserializer=scanner__pb2.ScanResult.FromString,
                )
        self.StreamResults = channel.stream_stream(
                '/scanner.ScannerService/StreamResults',
                request_serializer=scanner__pb2.ScanResult.SerializeToString,
                response_deserializer=scanner__pb2.StreamAck.FromString,
                )
        self.GetTaskStatus = channel.unary_unary(
                '/scanner.ScannerService/GetTaskStatus',
                request_serializer=scanner__pb2.TaskStatusRequest.SerializeToString,
                response_deserializer=scanner__pb2.TaskStatus.FromString,
                )
        self.GetScanStatus = channel.unary_unary(
                '/scanner.ScannerService/GetScanStatus',
                request_serializer=scanner__pb2.StatusRequest.SerializeToString,
                response_deserializer=scanner__pb2.ScanStatus.FromString,
                )
        self.Heartbeat = channel.unary_unary(
                '/scanner.ScannerService/Heartbeat',
                request_serializer=scanner__pb2.HeartbeatRequest.SerializeToString,
                response_deserializer=scanner__pb2.HeartbeatResponse.FromString,
                )


class ScannerServiceServicer(object):
    """Scanner service definition
    """

    def ExecuteScan(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def StreamResults(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetTaskStatus(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetScanStatus(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Heartbeat(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ScannerServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'ExecuteScan': grpc.unary_unary_rpc_method_handler(
                    servicer.ExecuteScan,
                    request_deserializer=scanner__pb2.ScanTask.FromString,
                    response_serializer=scanner__pb2.ScanResult.SerializeToString,
            ),
            'StreamResults': grpc.stream_stream_rpc_method_handler(
                    servicer.StreamResults,
                    request_deserializer=scanner__pb2.ScanResult.FromString,
                    response_serializer=scanner__pb2.StreamAck.SerializeToString,
            ),
            'GetTaskStatus': grpc.unary_unary_rpc_method_handler(
                    servicer.GetTaskStatus,
                    request_deserializer=scanner__pb2.TaskStatusRequest.FromString,
                    response_serializer=scanner__pb2.TaskStatus.SerializeToString,
            ),
            'GetScanStatus': grpc.unary_unary_rpc_method_handler(
                    servicer.GetScanStatus,
                    request_deserializer=scanner__pb2.StatusRequest.FromString,
                    response_serializer=scanner__pb2.ScanStatus.SerializeToString,
            ),
            'Heartbeat': grpc.unary_unary_rpc_method_handler(
                    servicer.Heartbeat,
                    request_deserializer=scanner__pb2.HeartbeatRequest.FromString,
                    response_serializer=scanner__pb2.HeartbeatResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'scanner.ScannerService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ScannerService(object):
    """Scanner service definition
    """

    @staticmethod
    def ExecuteScan(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/scanner.ScannerService/ExecuteScan',
            scanner__pb2.ScanTask.SerializeToString,
            scanner__pb2.ScanResult.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def StreamResults(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(request_iterator, target, '/scanner.ScannerService/StreamResults',
            scanner__pb2.ScanResult.SerializeToString,
            scanner__pb2.StreamAck.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetTaskStatus(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/scanner.ScannerService/GetTaskStatus',
            scanner__pb2.TaskStatusRequest.SerializeToString,
            scanner__pb2.TaskStatus.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetScanStatus(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/scanner.ScannerService/GetScanStatus',
            scanner__pb2.StatusRequest.SerializeToString,
            scanner__pb2.ScanStatus.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def Heartbeat(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/scanner.ScannerService/Heartbeat',
            scanner__pb2.HeartbeatRequest.SerializeToString,
            scanner__pb2.HeartbeatResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
