package httptool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github/lhh-gh/go-mall/comon/errcode"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/comon/util"

	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var _Client *http.Client

func getHttpClient() *http.Client {
	if _Client != nil {
		return _Client
	}

	tr := &http.Transport{
		//Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接, 0 表示不限制
		IdleConnTimeout:       90 * time.Second, // 空闲连接在连接池中的最大存活时间, 0表示不限制
		MaxIdleConnsPerHost:   50,               //每个Host (host + port) 的最大空闲连接, 不设置默认用 DefaultMaxIdleConnsPerHost
		MaxConnsPerHost:       50,               // 每个Host 的最大连接, 0 表示不限制 与 MaxIdleConnsPerHost 相等(尽量使用空闲连接)
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second, // waiting time for TLS handshake, zero means no timeout
		ExpectContinueTimeout: 1 * time.Second,
	}

	_Client = &http.Client{Transport: tr}
	return _Client
}

func Request(method string, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	start := time.Now()
	reqOpts := defaultRequestOptions() // 默认的请求选项
	for _, opt := range options {      // 在reqOpts上应用通过options设置的选项
		err = opt.apply(reqOpts)
		if err != nil {
			return
		}
	}
	log := logger.New(reqOpts.ctx)
	defer func() {
		if err != nil {
			log.Error("HTTP_REQUEST_ERROR_LOG", "method", method, "url", url, "body", reqOpts.data, "reply", respBody, "err", err)
		}
	}()
	// 创建请求对象
	req, err := http.NewRequest(method, url, bytes.NewReader(reqOpts.data))
	if err != nil {
		return
	}
	reqOpts.ctx, _ = context.WithTimeout(reqOpts.ctx, reqOpts.timeout) // 给 Request 设置Timeout
	req = req.WithContext(reqOpts.ctx)
	defer req.Body.Close()

	// 在Header中添加追踪信息 把内部服务串起来
	traceId, spanId, _ := util.GetTraceInfoFromCtx(reqOpts.ctx)
	reqOpts.headers["traceid"] = traceId
	reqOpts.headers["spanid"] = spanId
	if len(reqOpts.headers) != 0 { // 设置请求头
		for key, value := range reqOpts.headers {
			req.Header.Add(key, value)
		}
	}
	// 发起请求
	client := getHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 记录请求日志
	dur := time.Since(start).Milliseconds()
	defer func() {
		if dur >= 3000 { // 超过 3s 返回, 记一条 Warn 日志
			log.Warn("HTTP_REQUEST_SLOW_LOG", "method", method, "url", url, "body", reqOpts.data, "reply", respBody, "err", err, "dur/ms", dur)
		} else {
			log.Debug("HTTP_REQUEST_DEBUG_LOG", "method", method, "url", url, "body", string(reqOpts.data), "reply", string(respBody), "err", err, "dur/ms", dur)
		}
	}()

	httpStatusCode = resp.StatusCode
	if httpStatusCode != http.StatusOK {
		// 返回非 200 时Go的 http 库不回返回error, 这里处理成error 调用方好判断
		err = errcode.Wrap("request api error", errors.New(fmt.Sprintf("non 200 response, response code: %d", httpStatusCode)))
		return
	}

	respBody, _ = ioutil.ReadAll(resp.Body)
	return
}

// Get 发起GET请求
func Get(ctx context.Context, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	options = append(options, WithContext(ctx))
	return Request("GET", url, options...)
}

// Post 发起POST请求
func Post(ctx context.Context, url string, data []byte, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	// 默认自带Header Content-Type: application/json 可通过 传递 WithHeaders 增加或者覆盖Header信息
	defaultHeader := map[string]string{"Content-Type": "application/json"}
	var newOptions []Option
	newOptions = append(newOptions, WithHeaders(defaultHeader), WithData(data), WithContext(ctx))
	newOptions = append(newOptions, options...)

	httpStatusCode, respBody, err = Request("POST", url, newOptions...)
	return
}

// 针对可选的HTTP请求配置项，模仿gRPC使用的Options设计模式实现
type requestOption struct {
	ctx     context.Context
	timeout time.Duration
	data    []byte
	headers map[string]string
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(option *requestOption) error

func (f optionFunc) apply(opts *requestOption) error {
	return f(opts)
}

func defaultRequestOptions() *requestOption {
	return &requestOption{ // 默认请求选项
		ctx:     context.Background(),
		timeout: 5 * time.Second,
		data:    nil,
		headers: map[string]string{},
	}
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.ctx = ctx
		return
	})
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.timeout, err = timeout, nil
		return
	})
}

func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		for k, v := range headers {
			opts.headers[k] = v
		}
		return
	})
}

func WithData(data []byte) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.data, err = data, nil
		return
	})
}
