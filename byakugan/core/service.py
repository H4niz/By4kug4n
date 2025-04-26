import asyncio
import logging
from typing import Dict, List
from grpc import aio
from byakugan.proto import scanner_pb2 as pb
from byakugan.core.rule_engine import RuleEngine
from byakugan.core.scanner import ScannerClient

class CoreService:
    """Core service for managing scans and rules"""
    
    def __init__(self, config: Dict):
        self.config = config
        self.logger = logging.getLogger(__name__)
        self.rule_engine = RuleEngine(config['rules_dir'])
        self.scanner_client = ScannerClient(config['scanner'])

    async def execute_scan(self, task: pb.ScanTask) -> pb.ScanResult:
        """Execute scan task"""
        try:
            # Validate task 
            if not self._validate_task(task):
                return pb.ScanResult(
                    task_id=task.id,
                    success=False,
                    message="Invalid task parameters"
                )

            # Load rules
            rule = self.rule_engine.get_rule_by_id(task.rule_context.id)
            if not rule:
                return pb.ScanResult(
                    task_id=task.id,
                    success=False, 
                    message=f"Rule {task.rule_context.id} not found"
                )

            # Create scanner task
            scanner_task = self._create_scanner_task(task, rule)

            # Execute scan
            result = await self.scanner_client.execute_scan(scanner_task)
            
            # Process findings
            findings = self._process_findings(result, rule)

            return pb.ScanResult(
                task_id=task.id,
                success=True,
                findings=findings
            )

        except Exception as e:
            self.logger.error(f"Scan failed: {str(e)}")
            return pb.ScanResult(
                task_id=task.id,
                success=False,
                message=str(e)
            )

    def _validate_task(self, task: pb.ScanTask) -> bool:
        """Validate task parameters"""
        return (
            task.id and
            task.target and 
            task.target.url and
            task.rule_context and
            task.rule_context.id
        )

    def _create_scanner_task(self, task: pb.ScanTask, rule: Dict) -> pb.ScanTask:
        """Create scanner task from core task"""
        return pb.ScanTask(
            id=task.id,
            target=task.target,
            auth_context=task.auth_context,
            rule_context=pb.RuleContext(
                id=rule['id'],
                category=rule['category'],
                severity=rule['severity']
            ),
            payload=self._create_payload(rule),
            validation=self._create_validation(rule)
        )

    def _create_payload(self, rule: Dict) -> pb.Payload:
        """Create payload from rule definition"""
        return pb.Payload(
            insertion_points=[
                pb.InsertionPoint(
                    location=point['location'],
                    type=point['type'],
                    payloads=point['payloads']
                )
                for point in rule['insertion_points']
            ]
        )

    def _create_validation(self, rule: Dict) -> pb.Validation:
        """Create validation from rule definition"""
        return pb.Validation(
            success_conditions=pb.SuccessConditions(
                status_codes=rule['validation']['status_codes'],
                response_patterns=rule['validation']['patterns']
            )
        )