from typing import Dict, List
import json
from datetime import datetime
from jinja2 import Environment, FileSystemLoader
import markdown2
from dataclasses import dataclass

@dataclass
class ReportConfig:
    template_dir: str
    output_format: str = "html"
    include_evidence: bool = True
    severity_threshold: str = "low"

class ReportGenerator:
    def __init__(self, config: ReportConfig):
        self.config = config
        self.jinja_env = Environment(
            loader=FileSystemLoader(config.template_dir)
        )
        
    def generate(self, scan_results: Dict) -> str:
        """Generate formatted report from scan results"""
        template = self._get_template()
        
        # Filter findings based on severity threshold
        filtered_findings = self._filter_findings(
            scan_results["findings"]
        )
        
        # Process and format findings
        formatted_findings = self._format_findings(filtered_findings)
        
        # Generate report content
        context = {
            "scan_info": {
                "target": scan_results["metadata"]["target"],
                "start_time": scan_results["metadata"]["start_time"],
                "end_time": scan_results["metadata"]["end_time"],
                "duration": scan_results["summary"]["execution_time"]
            },
            "summary": scan_results["summary"],
            "findings": formatted_findings,
            "errors": scan_results["errors"],
            "generated_at": datetime.utcnow().isoformat()
        }
        
        return template.render(context)
    
    def _get_template(self) -> Environment:
        """Get appropriate template based on output format"""
        template_map = {
            "html": "report_template.html",
            "pdf": "report_template.html",  # Will be converted to PDF
            "json": "report_template.json",
            "markdown": "report_template.md"
        }
        return self.jinja_env.get_template(
            template_map[self.config.output_format]
        )
    
    def _filter_findings(self, findings: List[Dict]) -> List[Dict]:
        """Filter findings based on severity threshold"""
        severity_levels = {
            "critical": 4,
            "high": 3,
            "medium": 2,
            "low": 1,
            "info": 0
        }
        threshold = severity_levels[self.config.severity_threshold]
        
        return [
            finding for finding in findings
            if severity_levels[finding["severity"]] >= threshold
        ]
    
    def _format_findings(self, findings: List[Dict]) -> List[Dict]:
        """Format and enrich findings data"""
        for finding in findings:
            # Convert description to HTML if markdown
            if "description" in finding:
                finding["description_html"] = markdown2.markdown(
                    finding["description"]
                )
            
            # Add CVSS score visualization
            if "cvss_score" in finding:
                finding["cvss_color"] = self._get_cvss_color(
                    finding["cvss_score"]
                )
            
            # Format evidence if included
            if self.config.include_evidence and "evidence" in finding:
                finding["formatted_evidence"] = self._format_evidence(
                    finding["evidence"]
                )
        
        return findings
    
    def _get_cvss_color(self, score: float) -> str:
        """Get color code for CVSS score"""
        if score >= 9.0:
            return "#cc0000"  # Critical - Red
        elif score >= 7.0:
            return "#ff6600"  # High - Orange
        elif score >= 4.0:
            return "#ffcc00"  # Medium - Yellow
        else:
            return "#00cc00"  # Low - Green
    
    def _format_evidence(self, evidence: Dict) -> Dict:
        """Format evidence data for display"""
        return {
            "request": self._format_http_message(evidence.get("request", {})),
            "response": self._format_http_message(evidence.get("response", {})),
            "proof": evidence.get("proof"),
            "screenshots": evidence.get("screenshots", [])
        }
    
    def _format_http_message(self, message: Dict) -> str:
        """Format HTTP request/response for display"""
        if not message:
            return ""
            
        formatted = []
        
        # Add request/response line
        if "method" in message:  # Request
            formatted.append(f"{message['method']} {message['url']} HTTP/1.1")
        elif "status_code" in message:  # Response
            formatted.append(f"HTTP/1.1 {message['status_code']}")
            
        # Add headers
        for header, value in message.get("headers", {}).items():
            formatted.append(f"{header}: {value}")
            
        # Add body
        if "body" in message:
            formatted.append("")  # Empty line between headers and body
            if isinstance(message["body"], (dict, list)):
                formatted.append(
                    json.dumps(message["body"], indent=2)
                )
            else:
                formatted.append(str(message["body"]))
                
        return "\n".join(formatted)