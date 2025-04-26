import pytest
import xmltodict
from byakugan.core.parser import SOAPParser

@pytest.fixture
def soap_wsdl():
    return """<?xml version="1.0" encoding="UTF-8"?>
    <definitions xmlns="http://schemas.xmlsoap.org/wsdl/"
                 xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/">
        <portType name="TestPort">
            <operation name="SearchUsers">
                <input message="SearchRequest"/>
            </operation>
        </portType>
        <message name="SearchRequest">
            <part name="searchQuery" type="string"/>
        </message>
    </definitions>
    """

def test_parse_operations(soap_wsdl):
    """Test operation parsing"""
    parser = SOAPParser()
    api_def = parser.parse(soap_wsdl)
    
    assert len(api_def.endpoints) == 1
    endpoint = api_def.endpoints[0]
    
    assert endpoint.path == "/soap/SearchUsers"
    assert endpoint.method == "POST"
    assert len(endpoint.parameters) == 1
    
    param = endpoint.parameters[0]
    assert param.name == "searchQuery"
    assert param.type == "string"
    
    assert len(param.insertion_points) == 1
    point = param.insertion_points[0]
    assert point.param_type == "sql_injection"
    assert point.location == "body"

def test_soap_schema_generation(soap_wsdl):
    """Test SOAP envelope schema generation"""
    parser = SOAPParser()
    api_def = parser.parse(soap_wsdl)
    
    endpoint = api_def.endpoints[0]
    assert endpoint.request_body is not None
    assert endpoint.request_body.content_type == "application/soap+xml"
    
    schema = endpoint.request_body.schema
    assert "Envelope" in schema["properties"]
    assert "Body" in schema["properties"]["Envelope"]["properties"]