import grpc
from typing import Dict, Any
from . import scanner_pb2, scanner_pb2_grpc

class ScannerServicer(scanner_pb2_grpc.ScannerServiceServicer):
    """gRPC service implementation for scanner"""
    
    def __init__(self):
        self.engine = None

    def ExecuteScan(self, request: scanner_pb2.ScanTask, context) -> scanner_pb2.ScanResult:
        """Execute scan task"""
        try:
            # Convert protobuf task to internal format
            task = self._convert_task(request)
            
            # Execute scan using engine
            result = self.engine.execute_task(task)
            
            # Convert result back to protobuf
            return self._convert_result(result)
            
        except Exception as e:
            if context:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details(str(e))
            return scanner_pb2.ScanResult()

    def _convert_task(self, proto_task: scanner_pb2.ScanTask) -> Dict[str, Any]:
        """Convert protobuf task to internal format"""
        return {
            "id": proto_task.id,
            "target": {
                "url": proto_task.target.url,
                "method": proto_task.target.method,
                "protocol": proto_task.target.protocol
            },
            "auth_context": {
                "type": proto_task.auth_context.type,
                "headers": dict(proto_task.auth_context.headers)
            },
            "payload": {
                "headers": dict(proto_task.payload.headers),
                "query_params": dict(proto_task.payload.query_params),
                "insertion_points": [
                    {
                        "location": point.location,
                        "type": point.type,
                        "payloads": list(point.payloads)
                    }
                    for point in proto_task.payload.insertion_points
                ]
            },
            "validation": {
                "success_conditions": {
                    "status_codes": list(proto_task.validation.success_conditions.status_codes),
                    "timing_threshold": proto_task.validation.success_conditions.timing_threshold,
                    "error_patterns": list(proto_task.validation.success_conditions.error_patterns)
                },
                "evidence_collection": {
                    "save_request": proto_task.validation.evidence_collection.save_request,
                    "save_response": proto_task.validation.evidence_collection.save_response,
                    "timing_analysis": proto_task.validation.evidence_collection.timing_analysis
                }
            }
        }

    def _convert_result(self, result: Dict[str, Any]) -> scanner_pb2.ScanResult:
        """Convert internal result to protobuf format"""
        findings = []
        for finding in result.get("findings", []):
            evidence = scanner_pb2.Evidence(
                data=finding["evidence"]["data"],
                description=finding["evidence"]["description"]
            )
            findings.append(scanner_pb2.Finding(
                id=finding["id"],
                rule_id=finding["rule_id"],
                severity=finding["severity"],
                evidence=evidence
            ))

        return scanner_pb2.ScanResult(
            task_id=result["task_id"],
            success=result["success"],
            findings=findings
        )