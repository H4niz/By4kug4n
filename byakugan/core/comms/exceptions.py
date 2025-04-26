class CommsError(Exception):
    """Base exception for communication errors"""
    pass

class ConnectionError(CommsError):
    """Error establishing connection"""
    pass

class TimeoutError(CommsError):
    """Operation timed out"""
    pass

class SecurityError(CommsError):
    """Security-related error"""
    pass