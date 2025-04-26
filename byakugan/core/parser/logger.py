import logging
import os
from datetime import datetime
from typing import Dict

def setup_parser_logger(config: Dict) -> logging.Logger:
    """Setup parser logging"""
    
    # Create logs directory if not exists
    log_dir = "logs"
    if not os.path.exists(log_dir):
        os.makedirs(log_dir)
        
    # Create logger
    logger = logging.getLogger("parser")
    logger.setLevel(config.get("level", "INFO"))
    
    # Create file handler
    log_file = os.path.join(
        log_dir,
        f"parser_{datetime.now().strftime('%Y%m%d')}.log"
    )
    handler = logging.FileHandler(log_file)
    
    # Create formatter
    formatter = logging.Formatter(
        config.get("format", '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    )
    handler.setFormatter(formatter)
    
    # Add handler to logger
    logger.addHandler(handler)
    
    return logger