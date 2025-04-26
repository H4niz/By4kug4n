"""Task execution orchestrator"""
from dataclasses import dataclass
from typing import Dict, Any, Optional, List
import time
import logging
import asyncio
from datetime import datetime

from byakugan.config.comms import CommsConfig
from byakugan.core.comms.grpc.client import ScannerClient
# Import directly from proto package
from byakugan.core.comms.grpc.proto.scanner_pb2 import ScanTask, ScanResult, TaskStatus
from byakugan.core.comms.exceptions import ConnectionError

@dataclass
class ExecutionResult:
    """Result of task execution"""
    success: bool
    message: str
    data: Optional[Dict[str, Any]] = None
    error: Optional[Exception] = None

class TaskExecutor:
    """Executes scan tasks using scanner nodes"""

    def __init__(self, config: CommsConfig):
        """Initialize task executor"""
        self.config = config
        self.logger = logging.getLogger(__name__)  # Add logger initialization
        self._client = None
        self._initialize_client()
        self._active_tasks: Dict[str, ScanTask] = {}
        self._task_statuses: Dict[str, TaskStatus] = {}
        self._lock = asyncio.Lock()

    def _initialize_client(self) -> None:
        """Initialize scanner client with retries"""
        retries = 0
        max_retries = self.config.connection_retries
        while retries < max_retries:
            try:
                client = ScannerClient(self.config)
                if client.connect():
                    self._client = client
                    return
            except Exception as e:
                self.logger.error(f"Connection attempt {retries+1}/{max_retries} failed: {str(e)}")
            retries += 1
            time.sleep(1)  # Wait before retry

        raise ConnectionError(f"Failed to initialize scanner client after {max_retries} attempts")

    async def execute_batch(self, tasks: List[ScanTask]) -> List[ScanResult]:
        """Execute a batch of tasks concurrently"""
        async with self._lock:
            results = []
            tasks_to_execute = []
            
            for task in tasks:
                if self._validate_task(task):
                    self._active_tasks[task.id] = task
                    tasks_to_execute.append(self._execute_single(task))
                else:
                    results.append(ScanResult(
                        task_id=task.id,
                        success=False,
                        message="Invalid task parameters"
                    ))
            
            if tasks_to_execute:
                batch_results = await asyncio.gather(
                    *tasks_to_execute, 
                    return_exceptions=True
                )
                results.extend(batch_results)
            
            return results

    async def _execute_single(self, task: ScanTask) -> ScanResult:
        """Execute a single task"""
        try:
            # Update task status
            self._update_task_status(
                task.id, 
                "running",
                progress=0.0
            )

            # Execute scan
            result = await self._client.execute_scan_async(task)
            
            # Update final status
            self._update_task_status(
                task.id,
                "completed" if result.success else "failed",
                progress=100.0,
                message=result.message
            )

            return result

        except Exception as e:
            self.logger.error(f"Task {task.id} failed: {str(e)}")
            self._update_task_status(
                task.id,
                "failed",
                message=str(e)
            )
            return ScanResult(
                task_id=task.id,
                success=False,
                message=str(e)
            )

    def get_task_status(self, task_id: str) -> Optional[TaskStatus]:
        """Get status of a specific task"""
        return self._task_statuses.get(task_id)

    def _update_task_status(
        self, 
        task_id: str,
        status: str,
        progress: float = 0.0,
        message: str = ""
    ) -> None:
        """Update task status"""
        self._task_statuses[task_id] = TaskStatus(
            task_id=task_id,
            status=status,
            progress=progress,
            timestamp=datetime.now().timestamp(),
            message=message
        )

    def execute_task(self, task: ScanTask) -> ScanResult:
        """Execute a scan task"""
        try:
            if not self._validate_task(task):
                return ScanResult(
                    task_id=task.id,
                    success=False,
                    message="Invalid task parameters"
                )

            result = self._client.execute_scan(task)
            # Ensure success flag is properly set
            if result:
                result.success = True
            return result

        except Exception as e:
            self.logger.error("Task execution failed: %s", str(e))
            return ScanResult(
                task_id=task.id if task.id else "",
                success=False,
                message=str(e)
            )

    def _validate_task(self, task: ScanTask) -> bool:
        """Validate task parameters"""
        return (
            task.id and
            task.target and
            task.target.url and
            task.target.method and
            task.payload and
            task.payload.insertion_points
        )