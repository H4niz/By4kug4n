byakugan/core/
├── __init__.py
├── orchestrator/
│   ├── __init__.py
│   ├── scan_manager.py
│   ├── task_coordinator.py
│   └── exception_handler.py
├── parser/
│   ├── __init__.py 
│   ├── base.py
│   ├── openapi.py
│   ├── postman.py
│   ├── graphql.py
│   └── normalizer.py
├── rule_engine/
│   ├── __init__.py
│   ├── loader.py
│   ├── validator.py
│   ├── matcher.py
│   └── executor.py
├── scheduler/
│   ├── __init__.py
│   ├── task.py
│   ├── queue.py
│   └── balancer.py
├── auth/
│   ├── __init__.py
│   ├── manager.py
│   ├── oauth.py
│   ├── jwt.py
│   └── session.py
└── config/
    ├── __init__.py
    ├── settings.py
    └── constants.py