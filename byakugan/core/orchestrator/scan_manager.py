import logging
from uuid import UUID, uuid4
from datetime import datetime
from typing import Dict, List, Optional
from .models import ScanConfig, ScanTask
from ..parser import ApiParser
from ..rule_engine import RuleEngine
from .task_coordinator import TaskCoordinator
from .result_aggregator import ResultAggregator
from .types import ScanStatus

class ScanManager:
    """Manages overall scan workflow"""
    
    def __init__(self, config: ScanConfig):
        self.config = config
        self.logger = logging.getLogger(__name__)
        self.parser = ApiParser()
        self.rule_engine = RuleEngine(config)
        self.task_coordinator = TaskCoordinator(config)
        self.result_aggregator = ResultAggregator()
        
        # Track active scans
        self.active_scans: Dict[UUID, ScanStatus] = {}

    def get_scan_status(self, scan_id: UUID) -> Optional[ScanStatus]:
        """Get current scan status"""
        return self.active_scans.get(scan_id)

    async def start_scan(self) -> UUID:
        """Execute complete scan workflow"""
        try:
            scan_id = uuid4()
            self.active_scans[scan_id] = ScanStatus(
                id=scan_id,
                status="running",
                progress=0.0,
                start_time=datetime.now()
            )

            # Parse API spec and load rules
            api_def = self.parser.parse(self.config.api_definition)
            rules = self.rule_engine.load_rules(self.config.rules_dir)

            # Generate and execute tasks
            tasks = self._generate_tasks(api_def, rules, scan_id)
            results = await self.task_coordinator.execute_tasks(tasks)

            # Process results
            self.result_aggregator.process_results(results)
            
            self.active_scans[scan_id].status = "completed"
            return scan_id

        except Exception as e:
            if scan_id in self.active_scans:
                self.active_scans[scan_id].status = "failed"
                self.active_scans[scan_id].error = str(e)
            raise

    def _generate_tasks(self, api_def: Dict, rules: List[Dict], scan_id: UUID) -> List[ScanTask]:
        """Generate scan tasks for endpoints and rules"""
        tasks = []
        for endpoint in api_def.endpoints:
            matching_rules = self.rule_engine.get_matching_rules(endpoint)
            for rule in matching_rules:
                tasks.append(ScanTask(
                    id=str(uuid4()),
                    scan_id=scan_id,
                    endpoint=endpoint,
                    rule=rule
                ))
        return tasks