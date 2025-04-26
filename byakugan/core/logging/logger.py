import logging
import json
from datetime import datetime
from typing import Dict, Any
from pathlib import Path

class ScanLogger:
    def __init__(self, config: Dict[str, Any]):
        self.log_dir = Path(config.get("log_dir", "logs"))
        self.log_dir.mkdir(exist_ok=True)
        
        # Configure file handler
        self.file_handler = logging.FileHandler(
            self.log_dir / f"scan_{datetime.now().strftime('%Y%m%d_%H%M%S')}.log"
        )
        self.file_handler.setFormatter(
            logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
        )
        
        # Configure logger
        self.logger = logging.getLogger("byakugan.scan")
        self.logger.setLevel(logging.DEBUG)
        self.logger.addHandler(self.file_handler)
        
    def log_scan_start(self, scan_config: Dict):
        """Log scan initialization"""
        self.logger.info(
            "Starting new scan",
            extra={
                "event": "scan_start",
                "config": json.dumps(scan_config)
            }
        )
    
    def log_finding(self, finding: Dict):
        """Log discovered vulnerability"""
        self.logger.warning(
            f"Found {finding['severity']} severity issue: {finding['title']}",
            extra={
                "event": "finding",
                "finding_data": json.dumps(finding)
            }
        )
    
    def log_error(self, error: Exception, context: Dict = None):
        """Log error with context"""
        self.logger.error(
            str(error),
            extra={
                "event": "error",
                "error_type": error.__class__.__name__,
                "context": json.dumps(context) if context else None
            },
            exc_info=True
        )