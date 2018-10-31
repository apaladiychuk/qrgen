package serverapi

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mozillazg/request"
	"net/http"
	"time"
)

var certificate = `-----BEGIN CERTIFICATE-----
MIICATCCAWoCCQDBKmdEUmMP8DANBgkqhkiG9w0BAQsFADBFMQswCQYDVQQGEwJV
QTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0
cyBQdHkgTHRkMB4XDTE3MDMxMDA4MjQxN1oXDTE4MDMxMDA4MjQxN1owRTELMAkG
A1UEBhMCVUExEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0
IFdpZGdpdHMgUHR5IEx0ZDCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA6nvL
psguckrNIiPJV7nJyU3yMqp5klDLE9gommJH3UMfNYyxxfbOFVWDD8+9xVEH5SdW
2Vweznt2h+aN4UJWB6AijH++mTW8+eFqGOPzadv7W6O7BXjcHlCuRLV2bmaS/Kfw
PNX32rkNhOEC025beb4JzvXO1efUQUyT125ZZG8CAwEAATANBgkqhkiG9w0BAQsF
AAOBgQA849zA97iozln+MpxsPbLVlWO+oMgBIoTUkSqLXNUOQiM9gxm5f3o4vm8A
0TJG3psR5wUwSkyas5119gBAPxEbNChpm5Yq/F1t1+bG5fz1Hm1k+185I+jZVjxt
KxH6CejI0IkKSEDBuQrOLbKKdQySpIT2QQv/n1tZVC7wKVp3hA==
-----END CERTIFICATE-----
`
var (
	tlsConf   *tls.Config
	client    *http.Client
	BaseUrl   string
	sessionId string
	UserId    string
	CloudHost string
	CloudPort string
)

// go build -ldflags "-s -w"
func init() {
	CloudHost = "https://34.208.211.74"
	//
	// CloudHost = "https://localhost"
	CloudPort = "10443"
	BaseUrl = fmt.Sprintf("%s:%s", CloudHost, CloudPort)

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(certificate))
	if !ok {
		panic("failed to parse root certificate")
	}
	tlsConf = &tls.Config{RootCAs: roots,
		InsecureSkipVerify: true}

	transport := &http.Transport{TLSClientConfig: tlsConf}

	client = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	Connect()
}

func SetHeader(router *request.Request) {
	router.Headers = map[string]string{
		"sessionId":  sessionId,
		"customerId": UserId,
	}

}

//  template for execute query
func ExecQuery(r func(req *request.Request) (*request.Response, error)) (*request.Response, []byte) {
	if client == nil {
		fmt.Println(" httpclient is null ")
		return nil, nil
	}

	req := request.NewRequest(client)
	SetHeader(req)

	resp, err := r(req)
	if err != nil {
		fmt.Println("[REQUEST]", err.Error())

		return resp, []byte("serverError")
	}

	var body []byte
	body = make([]byte, resp.ContentLength)

	resp.Body.Read(body)
	resp.Body.Close()
	return resp, body

}

// Connect to server
func Connect() {
	//req := request.NewRequest(client)
	//req.Data = map[string]string{
	//	"login":    "manuf",
	//	"password": "manuf",
	//}
	//resp, err := req.Post(BaseUrl + "/v1/security/")
	//if err != nil {
	//	fmt.Println("POST request ", err.Error())
	//} else {
	//	j, err := resp.Json()
	//	resp.Body.Close()
	//	if err != nil {
	//		fmt.Println("JSON conv ", err.Error())
	//	} else {
	//		sessionId, _ = j.Get("SessionId").String()
	//		UserId, _ = j.Get("UserId").String()
	//		fmt.Println("JSON ", j)
	//		fmt.Println(" session id = ", sessionId)
	//		fmt.Println(" user id = ", UserId)
	//	}
	//}
}

func UploadInventory(modelId string, modelName string) {
	resp, body := ExecQuery(func(req *request.Request) (*request.Response, error) {
		req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		params := make(map[string]string)
		params["modelId"] = modelId
		params["modelName"] = modelName

		fmt.Println(">> Send ", BaseUrl+"/v2/manuf/inventory/")
		res, err := req.PostForm(BaseUrl+"/v2/manuf/inventory/", params)
		if err != nil {
			fmt.Println("<< RECV ", err.Error())
		}
		return res, err
	})
	fmt.Printf("<<< Resp   %d ", resp.StatusCode)
	if resp.StatusCode == 200 {
		fmt.Println(string(body))

	} else {
		fmt.Printf("<<< ERROR  %d ", resp.StatusCode)
		fmt.Println("")
		fmt.Println("<<< ", string(body))
		// todo Save in temporary file for resend

	}

}
