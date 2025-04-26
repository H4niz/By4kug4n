import grpc
import yaml
import json
import uuid
from byakugan.proto import scanner_pb2, scanner_pb2_grpc

class CoreClient:
    def __init__(self, addr: str):
        self.config = self._load_config()
        options = self._get_channel_options()
        
        # Create gRPC channel
        creds = self._get_credentials() 
        self.channel = grpc.secure_channel(addr, creds, options) if creds else \
                      grpc.insecure_channel(addr, options)
                      
        # Initialize stub with fixed imports
        self.stub = scanner_pb2_grpc.ScannerServiceStub(self.channel)

    def _load_config(self):
        with open('config/byakugan.yaml') as f:
            config = yaml.safe_load(f)
            return config['comms']['grpc']

    def _get_channel_options(self):
        service_config = {
            'methodConfig': [{
                'name': [{'service': 'ScannerService'}],
                'retryPolicy': {
                    'maxAttempts': self.config.get('retries', 3),
                    'initialBackoff': f"{self.config.get('retry_delay', 5)}s",
                    'maxBackoff': '30s',
                    'backoffMultiplier': 2,
                    'retryableStatusCodes': ['UNAVAILABLE']
                },
                'timeout': f"{self.config.get('timeout', 30)}s"
            }]
        }

        return [
            ('grpc.enable_retries', 1),
            ('grpc.max_receive_message_length', self.config['max_message_size']),
            ('grpc.max_send_message_length', self.config['max_message_size']),
            ('grpc.keepalive_time_ms', 30000),
            ('grpc.keepalive_timeout_ms', 10000),
            # Convert dict to JSON string
            ('grpc.service_config_json', json.dumps(service_config))
        ]

    def _get_credentials(self):
        if self.config.get('security', {}).get('tls_enabled'):
            cert_path = self.config['security']['tls_cert_path']
            return grpc.ssl_channel_credentials(
                root_certificates=open(cert_path, 'rb').read()
            )
        return None

    def execute_scan(self, target: str, rules: list) -> scanner_pb2.ScanResult:
        """Execute security scan on target"""
        # Create scan task
        task = scanner_pb2.ScanTask(
            id=str(uuid.uuid4()),
            target=scanner_pb2.Target(
                url=target,
                method="GET",
                protocol="http"
            ),
            rule_context=scanner_pb2.RuleContext(
                id=rules[0],  # Primary rule
                severity="HIGH",
                category="security",
                required_evidence=rules[1:]  # Additional rules as evidence
            )
        )

        try:
            result = self.stub.ExecuteScan(task)
            return result
        except grpc.RpcError as e:
            status_code = e.code()
            if status_code == grpc.StatusCode.UNAVAILABLE:
                raise ConnectionError(f"Scanner unavailable: {e.details()}")
            raise RuntimeError(f"Scan failed: {e.details()}")

    def stream_results(self):
        """Stream scan results"""
        try:
            return self.stub.StreamResults()
        except grpc.RpcError as e:
            raise RuntimeError(f"Failed to stream results: {e.details()}")

    def get_task_status(self, task_id: str):
        """Get status of specific task"""
        request = scanner_pb2.TaskStatusRequest(task_id=task_id)
        try:
            return self.stub.GetTaskStatus(request)
        except grpc.RpcError as e:
            raise RuntimeError(f"Failed to get task status: {e.details()}")

    def close(self):
        """Close gRPC channel"""
        self.channel.close()