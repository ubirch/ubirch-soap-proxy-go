# ubirch-soap-proxy-go

Simple SOAP proxy for UBIRCH client requests.

```
(SOAP-Client) -> (Proxy/XML) -> (UBIRCH CLient/JSON)
```

### Configuration
A configuration file `config.json` is required and should be located in the working directory.
```json
{
  "verificationBaseURL": "https://SOMEDOMAIN.COM/verify",
  "ubirchClientURL": "http://localhost:8080"
}
```

### Request
- Protocol: `http`
- Address: `:8090`
- Headers:
    - `X-UUID`: `<UUID>`
    - `X-Auth-Token`: `<auth token>`
- Body: 
    ```xml
    <?xml version="1.0" encoding="utf-8"?>
    <soap:Envelope xmlns:soap='http://schemas.xmlsoap.org/soap/envelope/'>
      <soap:Body>
        <ubirch:Document xmlns:ubirch='http://ubirch.com/wsdl/1.0'>
          <ActionReferenceNumber>a</ActionReferenceNumber>
          <ActionID>1234567890</ActionID>
          <SpecialUseDesc>C32-cb12347-test</SpecialUseDesc>
          <PeriodBeginDate>2020-11-10</PeriodBeginDate>
          <PeriodBeginTime>11:30</PeriodBeginTime>
          <PeriodEndDate>2020-12-10</PeriodEndDate>
          <PeriodEndTime>12:35</PeriodEndTime>
          <PostCode>10997</PostCode>
          <City>Berlin</City>
          <District>Kreuzberg</District>
          <Street>Eisenbahnstr.</Street>
          <FromHouseNumber>42</FromHouseNumber>
          <ToHouseNumber>43</ToHouseNumber>
          <FromCrossroad>Muskauer Str.</FromCrossroad>
          <ToCrossroad>Wrangelstr.</ToCrossroad>
          <!-- optional -->
          <LicensePlate>B-PL 1234</LicensePlate>
          <!-- optional -->
          <GeoAreaCoordinates>52.5021851,13.4296059</GeoAreaCoordinates>
          <!-- optional -->
          <GeoOverviewCoordinates>52.5021851,13.4296059</GeoOverviewCoordinates>
        </ubirch:Document>
      </soap:Body>
    </soap:Envelope>
    ```
  
### Response
- On success:
  
  Status code: `200`
  
  Body:
  ```xml
  <CertificationResponse>
     <Hash>696vMmuQpXgNPugqCmRjKyItKtab8WAZYaHr1neZFPg=</Hash>
     <Upp>liPEEHkErEkhS0uZl3pDzlbAItvEQIrhN1/890jw19ZIxfwMkMg9tTFgpqIt5zHTOMc/qW4edTUsKwM0LbgOWwWO4RK/69VtHGa1BCEnWsOISowzUSEAxCDr3q8ya5CleA0+6CoKZGMrIi0q1pvxYBlhoevWd5kU+MRA8o/CimZXNcnq4OOWXpWLLiRCXUemnkY1ZEZEAxqHN5IyXYX5qWE7Y5WAz/VkeE5iZL6wlq2ot4Ga7LmNbBinng==</Upp>
     <Response>liPEEJ08eP8i80RBpdGFxjbUhv/EQPKPwopmVzXJ6uDjll6Viy4kQl1Hpp5GNWRGRAMahzeSMl2F+alhO2OVgM/1ZHhOYmS+sJatqLeBmuy5jWwYp54AxBCbcpUxcZ5HRoWGuX5tNarKxEgwRgIhANuUX2gPlYBj9r8x6882dhPvn4a5b0W3mDxqUxGrr3YjAiEAmiJ0UiMJ0QmGeDnQ4KmGcCRx7L2FzUIxGErxRgS0QSg=</Response>
     <VerificationURL>https://ubirch.com/gelsenkirchen#FromHouseNumber:42;ToHouseNumber:43;GeoAreaCoordinates:52.5021851,13.4296059;PeriodEndTime:12:35;Street:Eisenbahnstr.;District:Kreuzberg;FromCrossroad:Muskauer Str.;SpecialUseDesc:C32-cb12347-test;PeriodBeginTime:11:30;PeriodEndDate:2020-12-10;City:Berlin;LicensePlate:B-PL 1234;ActionReferenceNumber:a;ActionID:1234567890;ToCrossroad:Wrangelstr.;GeoOverviewCoordinates:52.5021851,13.4296059;PeriodBeginDate:2020-11-10;PostCode:10997</VerificationURL>
  </CertificationResponse>
  ```
  
- On fail:

  Status code: `>= 300`

  Body:
  ```xml
  <fault>
      <faultcode>soap:Server</faultcode>
      <faultstring>error message</faultstring>
  </fault>
  ```

# Example

Run both the [ubirch-client-go](https://github.com/ubirch/ubirch-client-go)
and this proxy. Make sure the `ubirchClientURL`above is correct.

## Sending Requests

The SOAP client must send authentication headers with the POST request:

```
X-UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
X-Auth-Token: xxxx
```

An example [request](example_request.xml) and [response](example_response.xml) are provided.
The response contains a URL that is then to be rendered into a QR code. Below is the
URL from the example response:

![Example QR Code](example_qrcode.png)

