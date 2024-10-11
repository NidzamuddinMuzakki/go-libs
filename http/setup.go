package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"go-cimb-lib/common"
	"go-cimb-lib/env"
	"go-cimb-lib/log"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type (
	HeaderHostItem struct {
		Path  string
		Value string
	}

	ResponseInfoSendHttp struct {
		Error      interface{} `json:"error"`
		StatusCode int         `json:"status_code"`
		Status     string      `json:"status"`
		Proto      string      `json:"proto"`
		Time       string      `json:"time"`
		ReceivedAt time.Time   `json:"received_at"`
		Body       interface{} `json:"body"`
	}

	RequestInfoSendHttp struct {
		DNSLookup      string      `json:"dns_lookup"`
		ConnTime       string      `json:"conn_time"`
		TCPConnTime    string      `json:"tcp_conn_time"`
		TLSHandshake   string      `json:"tls_handshake"`
		ServerTime     string      `json:"server_time"`
		ResponseTime   string      `json:"response_time"`
		TotalTime      string      `json:"total_time"`
		IsConnReused   bool        `json:"is_conn_reused"`
		IsConnWasIdle  bool        `json:"is_conn_was_idle"`
		ConnIdleTime   string      `json:"conn_idle_time"`
		RequestAttempt int         `json:"request_attempt"`
		RemoteAddr     net.Addr    `json:"remote_addr"`
		Body           interface{} `json:"body"`
	}

	InfoSendHttp struct {
		ResponseInfo ResponseInfoSendHttp `json:"response_info"`
		RequestInfo  RequestInfoSendHttp  `json:"request_info"`
	}
)

func Send(method string, url string, header http.Header, req interface{}) (res []byte, info InfoSendHttp, err error) {

	var response *resty.Response
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	request := client.R().EnableTrace().SetHeaderMultiValues(header)

	switch method {
	case http.MethodGet:
		if !common.IsNilInterface(req) {
			v, _ := query.Values(req)
			url = url + "?" + v.Encode()
			// request.SetQueryParams(req.(map[string]string))
		}
		response, err = request.Get(url)
	case http.MethodPost:
		request.SetBody(req)
		response, err = request.Post(url)
	case http.MethodDelete:
		request.SetBody(req)
		response, err = request.Delete(url)
	case http.MethodPut:
		request.SetBody(req)
		response, err = request.Put(url)
	default:
		err = errors.New("method not supported")
		return
	}
	if err != nil {
		return
	}
	ti := response.Request.TraceInfo()
	resLog := ResponseInfoSendHttp{
		Error:      response.Error(),
		StatusCode: response.StatusCode(),
		Status:     response.Status(),
		Proto:      response.Proto(),
		Time:       response.Time().String(),
		ReceivedAt: response.ReceivedAt(),
		Body:       response.String(),
	}

	reqLog := RequestInfoSendHttp{
		DNSLookup:      ti.DNSLookup.String(),
		ConnTime:       ti.ConnTime.String(),
		TCPConnTime:    ti.TCPConnTime.String(),
		TLSHandshake:   ti.TLSHandshake.String(),
		ServerTime:     ti.ServerTime.String(),
		ResponseTime:   ti.ResponseTime.String(),
		TotalTime:      ti.TotalTime.String(),
		IsConnReused:   ti.IsConnReused,
		IsConnWasIdle:  ti.IsConnWasIdle,
		ConnIdleTime:   ti.ConnIdleTime.String(),
		RequestAttempt: ti.RequestAttempt,
		RemoteAddr:     ti.RemoteAddr,
		Body:           req,
	}
	info.RequestInfo = reqLog
	info.ResponseInfo = resLog

	res = []byte(response.String())

	return
}

func SendFormData(method string, url string, header http.Header, req map[string]string) (res []byte, info InfoSendHttp, err error) {
	var resp *resty.Response
	client := resty.New()

	request := client.R().
		EnableTrace().
		SetHeaderMultiValues(header)
	request.SetFormData(req)
	switch method {

	case http.MethodGet:

		resp, err = request.Get(url)

	case http.MethodPost:
		resp, err = request.Post(url)
	case http.MethodDelete:
		resp, err = request.Delete(url)
	case http.MethodPut:
		resp, err = request.Put(url)
	default:
		err = errors.New("method not supported")
		return
	}
	if err != nil {
		return
	}
	ti := resp.Request.TraceInfo()
	resLog := ResponseInfoSendHttp{
		Error:      err,
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Proto:      resp.Proto(),
		Time:       resp.Time().String(),
		ReceivedAt: resp.ReceivedAt().Local(),
		Body:       resp.String(),
	}
	reqLog := RequestInfoSendHttp{
		DNSLookup:      ti.DNSLookup.String(),
		ConnTime:       ti.ConnTime.String(),
		TCPConnTime:    ti.TCPConnTime.String(),
		TLSHandshake:   ti.TLSHandshake.String(),
		ServerTime:     ti.ServerTime.String(),
		ResponseTime:   ti.ResponseTime.String(),
		TotalTime:      ti.TotalTime.String(),
		IsConnReused:   ti.IsConnReused,
		IsConnWasIdle:  ti.IsConnWasIdle,
		ConnIdleTime:   ti.ConnIdleTime.String(),
		RequestAttempt: ti.RequestAttempt,
		RemoteAddr:     ti.RemoteAddr,
		Body:           req,
	}
	info.RequestInfo = reqLog
	info.ResponseInfo = resLog
	res = []byte(resp.String())

	return
}

func GenerateHeader(traceID string, skipCaller int, envName string) http.Header {
	logging := log.NewLogging(skipCaller + 1)
	var header []HeaderHostItem

	ByteHeader, _ := json.Marshal(env.Interface(envName, nil))
	if len(ByteHeader) > 0 {
		err := json.Unmarshal(ByteHeader, &header)
		if err != nil {
			logging.Error("unmarshal", err.Error(), nil)
		}
	}

	res := make(http.Header)
	for _, item := range header {
		res.Add(item.Path, item.Value)
	}
	res.Add("Trace-ID", traceID)

	return res
}
