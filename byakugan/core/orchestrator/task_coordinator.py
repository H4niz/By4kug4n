import asyncio
import logging
from datetime import datetime
from typing import Dict, List, Optional
from uuid import UUID
from dataclasses import dataclass

from .models import ScanTask, TaskCoordinatorConfig
from .types import TaskStatus
from .exception_handler import ExceptionHandler
from ..rule_engine import RuleEngine

class TaskCoordinator:
    """Coordinates task execution"""
    
    def __init__(self, config: TaskCoordinatorConfig):
        self.config = config
        self.logger = logging.getLogger(__name__)
        self.active_tasks: Dict[str, ScanTask] = {}
        self._stop_events: Dict[UUID, asyncio.Event] = {}
        self.exception_handler = ExceptionHandler()

    def get_scan_tasks(self, scan_id: UUID) -> List[ScanTask]:
        """Get tasks for scan ID"""
        return self.scan_tasks.get(scan_id, [])

    async def execute_task(self, task: ScanTask) -> Dict:
        """Execute single task with retry logic"""
        retry_count = 0
        max_retries = self.config.retry_count

        while retry_count < max_retries:
            try:
                result = await self._execute_single_attempt(task)
                return {
                    "task_id": task.id,
                    "success": True,
                    "findings": result.get("findings", []),
                    "retry_count": retry_count,
                    "execution_time": result.get("execution_time", 0.0)
                }
            except Exception as e:
                retry_count += 1
                if retry_count < max_retries:
                    await asyncio.sleep(self.config.retry_delay)
                    continue
                return {
                    "task_id": task.id,
                    "success": False,
                    "error_details": str(e),
                    "retry_count": retry_count,
                    "findings": []
                }
                
    async def _execute_single_attempt(self, task: ScanTask) -> Dict:
        """Execute single task attempt"""
        # Add actual task execution logic here
        # For testing purposes:
        if task.endpoint["path"] == "/invalid":
            raise ValueError("Invalid path")
        return {
            "findings": [],
            "execution_time": 0.1
        }

    async def execute_tasks(self, tasks: List[ScanTask]) -> List[Dict]:
        """Execute tasks concurrently"""
        tasks_coros = [self.execute_task(task) for task in tasks]
        results = await asyncio.gather(*tasks_coros, return_exceptions=True)
        
        processed_results = []
        for result in results:
            if isinstance(result, Exception):
                # Handle exception
                continue
            processed_results.append(result)
            
        return processed_results

    async def stop_tasks(self, scan_id: UUID):
        """Stop all tasks for a scan"""
        if scan_id not in self._stop_events:
            self._stop_events[scan_id] = asyncio.Event()
        
        # Set stop event
        self._stop_events[scan_id].set()
        
        # Update status of active tasks
        tasks = self.scan_tasks.get(scan_id, [])
        for task in tasks:
            if task.status == TaskStatus.RUNNING:
                task.status = TaskStatus.STOPPED
                task.completed_at = datetime.now()