package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/axengine/utils/sign"
	"github.com/shopspring/decimal"
)

// APISign API Param Sign
//  1. rawStr eg: GET http://example.com/hello?n=1&a=2  Key["n","a"]-ASC Sort["a","n"] GetParam(a) a=2&n=1 param key string attaches the methods
//  2. other request http method，for Content-Type: application/json {"n":"m","a":2} Key ASC Sort,param key string attaches the methods => {"a":2,"n":"m"} => a=2&n=m
//  3. rawStr+timestamp => a=2&n=m1626167650 (1626167650 is unix timestamp), verify sign time valid（default 10s）
//  4. Sign Method： Method(rawStr+timestamp, secretKey) signed text encode [Base64, Hex(default)]
//     Method=[HMAC-SHA256,HMAC-SHA1] Encode=[Base64,Hex] Default = HMAC-SHA256-HEX
//  5. default: signStr=Hex(HMAC-SHA256(rawStr+timestamp,secretKey))
//  6. Sign http request Header X-Signature=accessKey:signStr:timestamp (: split elem)
type APISign struct {
	AccessKey    string
	AccessSecret string
	TTL          int64
	Method
}

type Method string

const (
	HmacSha256    Method = "HMAC-SHA256-BASE64"
	HmacSha1      Method = "HMAC-SHA1-BASE64"
	HmacSha1Hex   Method = "HMAC-SHA1-HEX"
	HmacSha256Hex Method = "HMAC-SHA256-HEX"
)

func NewAPISign(accessKey, accessSecret string, ttl int64, method Method) *APISign {
	return &APISign{
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
		TTL:          ttl,
		Method:       method,
	}
}

// Verify param sign result verify
// req:the http.Request
// authHeaderName:header name,like X-Signature,value:accessKey:signature:deadline
func (p *APISign) Verify(req *http.Request, authHeaderName string) error {
	authStr := req.Header.Get(authHeaderName)
	if authStr == "" {
		return fmt.Errorf("signature header required:%s", authHeaderName)
	}
	ss := strings.Split(authStr, ":")
	if len(ss) != 3 {
		return fmt.Errorf("signature header invalid:%s", authHeaderName)
	}
	accessKey, signature, deadlines := ss[0], ss[1], ss[2]

	deadline, err := strconv.ParseInt(deadlines, 10, 64)
	if err != nil || deadline < time.Now().Unix() {
		return errors.New("signature expired")
	}
	if accessKey != p.AccessKey {
		return errors.New("invalid AccessKey")
	}

	raw, err := p.ToSignRaw(req)
	if err != nil {
		return err
	}
	raw = fmt.Sprintf("%s%d", raw, deadline)
	checkSignature := signHash(p.Method, []byte(raw), []byte(p.AccessSecret))
	if checkSignature != signature {
		return fmt.Errorf("sign method invalid raw:%s", raw)
	}
	return nil
}

// Sign sign and return signature,deadline is validity period of signature
func (p *APISign) Sign(req *http.Request, deadline int64) (string, error) {
	raw, err := p.ToSignRaw(req)
	if err != nil {
		return "", err
	}
	if deadline == 0 {
		deadline = time.Now().Unix() + p.TTL
	}
	raw = fmt.Sprintf("%s%d", raw, deadline)
	return signHash(p.Method, []byte(raw), []byte(p.AccessSecret)), nil
}

// ToSignRaw 从req解析并生成已排序待签名参数字符串
// 如果contentType是JSON，只取JSON参数
// 其他：取URL PARAMS和BODY PARAMS
func (p *APISign) ToSignRaw(req *http.Request) (string, error) {
	var raw string
	contentType := req.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "application/json"):
		bz, err := io.ReadAll(req.Body)
		if err != nil {
			_ = req.Body.Close()
			return "", err
		}
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(bz))

		var reqBody = make(requestBodyMap)
		if err := json.Unmarshal(bz, &reqBody); err != nil {
			return "", err
		}
		raw, _ = reqBody.SortToString("&")
	default:
		if err := req.ParseForm(); err != nil {
			return "", err
		}

		if req.Form != nil && len(req.Form) > 0 {
			var paramNames []string
			for k := range req.Form {
				paramNames = append(paramNames, k)
			}
			sort.Strings(paramNames)

			var query []string
			for _, k := range paramNames {
				query = append(query, url.QueryEscape(k)+"="+url.QueryEscape(req.Form.Get(k)))
			}
			raw = strings.Join(query, "&")
		}
	}
	return raw, nil
}

// signHash hash and encode
func signHash(method Method, rawStr, secretKey []byte) (hash string) {
	switch method {
	case HmacSha1:
		hash = sign.HMACSha1B64(rawStr, secretKey)
	case HmacSha256:
		hash = sign.HMACSha256B64(rawStr, secretKey)
	case HmacSha1Hex:
		hash = sign.HMACSha1Hex(rawStr, secretKey)
	case HmacSha256Hex:
		hash = sign.HMACSha256Hex(rawStr, secretKey)
	default:
		hash = sign.HMACSha256Hex(rawStr, secretKey)
	}
	return
}

type requestBodyMap map[string]interface{}

// SortToString request body param sort format
func (r requestBodyMap) SortToString(separator string) (string, error) {
	if len(r) == 0 {
		return "", nil
	}
	kvs := make(KvSlice, 0)
	for k, v := range r {
		kvs = append(kvs, Kv{Key: k, Value: v})
	}

	sort.Sort(kvs)
	var s = make([]string, 0, len(kvs))
	for _, v := range kvs {
		switch v.Value.(type) {
		case float64:
			s = append(s, fmt.Sprintf("%s=%s", v.Key, decimal.NewFromFloat(v.Value.(float64)).String()))
		case float32:
			s = append(s, fmt.Sprintf("%s=%s", v.Key, decimal.NewFromFloat(float64(v.Value.(float32))).String()))
		case *float64:
			s = append(s, fmt.Sprintf("%s=%s", v.Key, decimal.NewFromFloat(*v.Value.(*float64)).String()))
		case *float32:
			s = append(s, fmt.Sprintf("%s=%s", v.Key, decimal.NewFromFloat(float64(*v.Value.(*float32))).String()))
		case string:
			s = append(s, fmt.Sprintf("%s=%s", v.Key, v.Value))
		case *string:
			s = append(s, fmt.Sprintf("%s=%s", v.Key, *v.Value.(*string)))
		default:
			buf := make([]byte, 0)
			buffer := bytes.NewBuffer(buf)
			if err := json.NewEncoder(buffer).Encode(v.Value); err != nil {
				return "", err
			}
			s = append(s, fmt.Sprintf("%s=%s", v.Key, string(r.trimNewline(buffer.Bytes()))))
		}
	}
	return strings.Join(s, separator), nil
}

func (r requestBodyMap) trimNewline(buf []byte) []byte {
	if i := len(buf) - 1; i >= 0 {
		if buf[i] == '\n' {
			buf = buf[:i]
		}
	}
	return buf
}

type Kv struct {
	Key   string
	Value interface{}
}
type KvSlice []Kv

func (s KvSlice) Len() int           { return len(s) }
func (s KvSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s KvSlice) Less(i, j int) bool { return s[i].Key < s[j].Key }
