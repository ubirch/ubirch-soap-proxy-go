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

### Mapping

The XML tags are mapped to short keys for the URL so it is not unnecessary long.

| XML Tag | URL field | Description |
|---------|-----------|-------------|
|ActionReferenceNumber| arn | reference number of the certification |
|ActionID| id | internal identity number |
|SpecialUseDesc| ud | description of the case |
|PeriodBeginDate| bd | start date |
|PeriodBeginTime| bt | start time |
|PeriodEndDate| ed | end date |
|PeriodEndTime| et | end time |
|PostCode| pc | post code|
|City| c | address city |
|District| d | address distict name |
|Street| s | address street name |
|FromHouseNumber| fn | starting house number |
|ToHouseNumber| tn| ending house number |
|FromCrossroad| fc | from crossroad |
|ToCrossroad| tc | to crossroad |
|*LicensePlate*| lp | license plate (**UNUSED**) |
|*GeoAreaCoordinates*| gac | area (**UNUSED**) |
|*GeoOverviewCoordinates*| goc | overview area (**UNUSED**) |

### Response

- On success:

  Status code: `200`

  Body:
  ```xml
  <?xml version="1.0" encoding="UTF-8"?>
  <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
      <soap:Body>
          <ubirch:CertificationResponse xmlns:ubirch="http://ubirch.com/wsdl/1.0">
              <Hash>OfS7M7hrQpvQDuXqHRt0f2FltnI94iyaTAyHZCy1ju8=</Hash>
              <UPP>liPEEHkErEkhS0uZl3pDzlbAItvEQHq8w1BkzlHevdxbfPRqA8fOorJL89QqiDJZeLXM4wp7tFhST+Ma1awnPbCEgIMABigFHDYHGQONXGemQotOQeQAxCA59LszuGtCm9AO5eodG3R/YWW2cj3iLJpMDIdkLLWO78RAfD/T4FgZfLbjpzZSw1YuNmfsCM1rCdWZFlB0Y07Jg04SX+DXJH3JrSGqDUv514005wrFPLSFH91/7zll1P4F2Q==</UPP>
              <Response>liPEEJ08eP8i80RBpdGFxjbUhv/EQHw/0+BYGXy246c2UsNWLjZn7AjNawnVmRZQdGNOyYNOEl/g1yR9ya0hqg1L+deNNOcKxTy0hR/df+85ZdT+BdkAxBAp232Tr9hALp+jdS83oJVjxEcwRQIhAP4a3txN+jwDBYETo2q5hDrzXS5OPxMejyCSROSdYqCkAiB4mNFBa68ASzsbJVvvEmNUS/H+nNApd1oNmdU/yvSjhQ==</Response>
              <URL>https://ubirch.com/gelsenkirchen#arn=a;pc=10997;s=Eisenbahnstr.;tn=43;fc=Muskauer%20Str.;lp=B-PL%201234;id=1234567890;tc=Wrangelstr.;bt=11:30;ed=2020-12-11;et=12:35;c=Berlin;d=Kreuzberg;fn=42;ud=C32-cb12347-test;bd=2020-11-09</URL>
          </ubirch:CertificationResponse>
      </soap:Body>
  </soap:Envelope>
  ```

- On fail:

  Status code: `>= 300`

  Body:
  ```xml
  <soap:Fault>
      <faultcode>soap:Server</faultcode>
      <faultstring>error message</faultstring>
  </soap:Fault>
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

An example [request](example_request.xml) and [response](example_response.xml) are provided. The response contains a URL
that is then to be rendered into a QR code. Below is the URL from the example response:

![Example QR Code](example_qrcode.png)

## Copyright

```
Copyright (c) 2020 ubirch GmbH

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
