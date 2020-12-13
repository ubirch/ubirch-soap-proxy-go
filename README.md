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
    <Upp>liPEEHkErEkhS0uZl3pDzlbAItvEQPKPwopmVzXJ6uDjll6Viy4kQl1Hpp5GNWRGRAMahzeSMl2F+alhO2OVgM/1ZHhOYmS+sJatqLeBmuy5jWwYp54AxCDr3q8ya5CleA0+6CoKZGMrIi0q1pvxYBlhoevWd5kU+MRAjlqV/SXR/DcK4D5MOfkr4RjRP1gd5v3nMx5yZ01EKJOQdXvSKxVy7fib+eSce/MZiI4/zVjDNsPyR5p13R86DQ==</Upp>
    <Response>liPEEJ08eP8i80RBpdGFxjbUhv/EQI5alf0l0fw3CuA+TDn5K+EY0T9YHeb95zMecmdNRCiTkHV70isVcu34m/nknHvzGYiOP81YwzbD8keadd0fOg0AxBAhKEg3IDZKBJ+hkKmieI0nxEYwRAIgOOxMHw7kASyTFLWEFihDX8HKyJo6duVuYVRqhlHCIO0CIEamhxCKNeSuNkvXy8bJqIOvkD3iGc4A7JQaOqmavzSt</Response>
    <VerificationURL>https://ubirch.com/gelsenkirchen#ActionReferenceNumber=a;PeriodEndDate=2020-12-10;ToHouseNumber=43;LicensePlate=B-PL 1234;SpecialUseDesc=C32-cb12347-test;PeriodBeginTime=11:30;PeriodEndTime=12:35;PostCode=10997;City=Berlin;Street=Eisenbahnstr.;GeoAreaCoordinates=52.5021851,13.4296059;GeoOverviewCoordinates=52.5021851,13.4296059;ActionID=1234567890;PeriodBeginDate=2020-11-10;District=Kreuzberg;FromHouseNumber=42;FromCrossroad=Muskauer Str.;ToCrossroad=Wrangelstr.</VerificationURL>
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

## Copyright

```
Copyright (c) 2019-2020 ubirch GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
