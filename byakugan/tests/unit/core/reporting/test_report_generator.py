import pytest
from datetime import datetime
from pathlib import Path
from byakugan.core.reporting.report_generator import ReportGenerator, ReportConfig

@pytest.fixture
def template_dir(tmp_path):
    """Create temporary template directory with sample templates"""
    template_dir = tmp_path / "templates"
    template_dir.mkdir()
    
    # Create HTML template
    html_template = """
    <html>
    <body>
        <h1>Scan Report</h1>
        <h2>Summary</h2>
        <p>Target: {{ scan_info.target }}</p>
        <h2>Findings</h2>
        {% for finding in findings %}
        <div class="finding" style="color: {{ finding.cvss_color }}">
            <h3>{{ finding.title }}</h3>
            <p>{{ finding.description_html }}</p>
            {% if finding.formatted_evidence %}
            <pre>{{ finding.formatted_evidence.request }}</pre>
            {% endif %}
        </div>
        {% endfor %}
    </body>
    </html>
    """
    (template_dir / "report_template.html").write_text(html_template)
    
    return template_dir

@pytest.fixture
def sample_scan_results():
    """Sample scan results for testing"""
    return {
        "metadata": {
            "target": "https://api.example.com",
            "start_time": "2024-04-10T10:00:00",
            "end_time": "2024-04-10T10:05:00"
        },
        "summary": {
            "execution_time": 300,
            "total_findings": 2,
            "severity_counts": {"high": 1, "medium": 1}
        },
        "findings": [
            {
                "title": "JWT None Algorithm",
                "severity": "high",
                "cvss_score": 9.1,
                "description": "**Critical** authentication bypass found",
                "evidence": {
                    "request": {
                        "method": "POST",
                        "url": "/auth/verify",
                        "headers": {"Content-Type": "application/json"},
                        "body": {"token": "eyJ0..."}
                    },
                    "response": {
                        "status_code": 200,
                        "body": {"authenticated": true}
                    }
                }
            },
            {
                "title": "Information Disclosure",
                "severity": "medium",
                "cvss_score": 5.5,
                "description": "Sensitive information exposed"
            }
        ],
        "errors": []
    }

@pytest.fixture
def report_generator(template_dir):
    config = ReportConfig(
        template_dir=str(template_dir),
        output_format="html",
        include_evidence=True,
        severity_threshold="low"
    )
    return ReportGenerator(config)

def test_generate_report(report_generator, sample_scan_results):
    """Test generating a complete report"""
    report = report_generator.generate(sample_scan_results)
    
    # Verify basic content
    assert "Scan Report" in report
    assert "https://api.example.com" in report
    assert "JWT None Algorithm" in report
    assert "Information Disclosure" in report
    
    # Verify formatting
    assert "#cc0000" in report  # CVSS color for critical finding
    assert "<strong>Critical</strong>" in report  # Markdown conversion
    assert "POST /auth/verify HTTP/1.1" in report  # Evidence formatting

def test_filtering_by_severity(report_generator, sample_scan_results):
    """Test filtering findings by severity threshold"""
    report_generator.config.severity_threshold = "high"
    report = report_generator.generate(sample_scan_results)
    
    # Should only include high severity finding
    assert "JWT None Algorithm" in report
    assert "Information Disclosure" not in report

def test_evidence_formatting(report_generator):
    """Test HTTP message formatting"""
    evidence = {
        "request": {
            "method": "POST",
            "url": "/api/test",
            "headers": {"X-Test": "Value"},
            "body": {"key": "value"}
        }
    }
    
    formatted = report_generator._format_evidence(evidence)
    request_text = formatted["request"]
    
    assert "POST /api/test HTTP/1.1" in request_text
    assert "X-Test: Value" in request_text
    assert '"key": "value"' in request_text

def test_cvss_color_coding(report_generator):
    """Test CVSS score to color mapping"""
    assert report_generator._get_cvss_color(9.5) == "#cc0000"  # Critical
    assert report_generator._get_cvss_color(7.5) == "#ff6600"  # High
    assert report_generator._get_cvss_color(5.0) == "#ffcc00"  # Medium
    assert report_generator._get_cvss_color(3.0) == "#00cc00"  # Low