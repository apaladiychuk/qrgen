package serverapi

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/mozillazg/request"
	"net/http"
	"strconv"
	"webapp/controllers"
	"webapp/models"
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
	CloudPort int
)

func Init() {
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
	reqCnt := 1
rerun:

	resp, err := r(req)
	if err != nil {
		fmt.Println("[REQUEST]", err.Error())
		if reqCnt < 3 {
			reqCnt++
			Connect()
			goto rerun
		}
		resp = &request.Response{}
		resp.Status = "200"
		resp.StatusCode = 200
		return resp, []byte("serverError")
	}

	if resp.StatusCode == 403 && reqCnt < 3 {
		reqCnt++
		Connect()
		goto rerun
	}
	var body []byte
	body = make([]byte, resp.ContentLength)

	resp.Body.Read(body)
	resp.Body.Close()
	return resp, body

}

// Connect to server
func Connect() {
	req := request.NewRequest(client)
	req.Data = map[string]string{
		"login":    "manuf",
		"password": "manuf",
	}
	resp, err := req.Post(BaseUrl + "/v1/security/")
	if err != nil {
		fmt.Println("POST request ", err.Error())
	} else {
		j, err := resp.Json()
		resp.Body.Close()
		if err != nil {
			fmt.Println("JSON conv ", err.Error())
		} else {
			sessionId, _ = j.Get("SessionId").String()
			UserId, _ = j.Get("UserId").String()
			fmt.Println("JSON ", j)
			fmt.Println(" session id = ", sessionId)
			fmt.Println(" user id = ", UserId)
		}
	}
}

func UploadInventory(modelId string, modelName string) {
	resp, body := ExecQuery(func(req *request.Request) (*request.Response, error) {
		req.Headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8"
		params := make(map[string]string)
		params["modelId"] = modelId
		params["modelName"] = modelName

		return req.PostForm(BaseUrl+"/v1/manuf/invento", params)
	})
	if resp.StatusCode == 200 {
		var respObj controllers.ImportMutableResponse
		if err := json.Unmarshal(body, &respObj); err != nil {
			return err
		} else {
			if respObj.Status == models.RESPONSE_OK {
				return nil
			} else {
				return errors.New("[error]" + respObj.Description)
			}
		}

	} else {
		return errors.New("web error " + resp.Status)
	}

}
