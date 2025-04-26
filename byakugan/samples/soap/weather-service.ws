dl
<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"
             xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/"
             xmlns:tns="http://weather.example.com/"
             targetNamespace="http://weather.example.com/">

    <types>
        <schema xmlns="http://www.w3.org/2001/XMLSchema">
            <element name="GetWeatherRequest">
                <complexType>
                    <sequence>
                        <element name="city" type="string"/>
                        <element name="country" type="string"/>
                    </sequence>
                </complexType>
            </element>
            
            <element name="GetWeatherResponse">
                <complexType>
                    <sequence>
                        <element name="temperature" type="decimal"/>
                        <element name="humidity" type="decimal"/>
                        <element name="conditions" type="string"/>
                        <element name="windSpeed" type="decimal"/>
                    </sequence>
                </complexType>
            </element>
        </schema>
    </types>

    <message name="GetWeatherInput">
        <part name="parameters" element="tns:GetWeatherRequest"/>
    </message>

    <message name="GetWeatherOutput">
        <part name="parameters" element="tns:GetWeatherResponse"/>
    </message>

    <portType name="WeatherPortType">
        <operation name="GetWeather">
            <input message="tns:GetWeatherInput"/>
            <output message="tns:GetWeatherOutput"/>
        </operation>
    </portType>

    <binding name="WeatherSoapBinding" type="tns:WeatherPortType">
        <soap:binding style="document" 
                     transport="http://schemas.xmlsoap.org/soap/http"/>
        <operation name="GetWeather">
            <soap:operation soapAction="http://weather.example.com/GetWeather"/>
            <input>
                <soap:body use="literal"/>
            </input>
            <output>
                <soap:body use="literal"/>
            </output>
        </operation>
    </binding>

    <service name="WeatherService">
        <port name="WeatherPort" binding="tns:WeatherSoapBinding">
            <soap:address location="http://weather.example.com/ws"/>
        </port>
    </service>

</definitions>