from abc import ABC, abstractmethod
from typing import Dict, Any, Optional

class CommunicationBase(ABC):
    """Base class for communication implementations"""
    
    @abstractmethod
    def connect(self) -> bool:
        """Establish connection"""
        pass
        
    @abstractmethod
    def close(self):
        """Close connection"""
        pass
        
    @abstractmethod 
    def is_connected(self) -> bool:
        """Check if connected"""
        pass