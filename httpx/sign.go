package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axengine/utils/sign"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

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
	TTL          time.Duration
	Method
}

type Method string

const (
	HmacSha256    Method = "HMAC-SHA256-BASE64"
	HmacSha1      Method = "HMAC-SHA1-BASE64"
	HmacSha1Hex   Method = "HMAC-SHA1-HEX"
	HmacSha256Hex Method = "HMAC-SHA256-HEX"
)

func NewAPISign(accessKey, accessSecret string, ttl time.Duration, method Method) *APISign {
	return &APISign{
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
		TTL:          ttl,
		Method:       method,
	}
}

// Verify param sign result verify
func (apiSign *APISign) Verify(req *http.Request, authHeaderKey string) error {
	val := req.Header.Get(authHeaderKey)
	if val == "" {
		return fmt.Errorf("signature header required:%s", authHeaderKey)
	}
	strs := strings.Split(val, ":")
	if len(strs) != 3 {
		return fmt.Errorf("signature header invalid:%s", authHeaderKey)
	}
	accessKey, signature, timestamp := strs[0], strs[1], strs[2]

	ttl, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || ttl < time.Now().Unix() {
		return errors.New("signature expired")
	}
	if accessKey != apiSign.AccessKey {
		return errors.New("invalid AccessKey")
	}

	var rawStr string
	switch strings.ToUpper(req.Method) {
	default:
		// not support http method
		return nil
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		ct := filterFlags(req.Header.Get("Content-Type"))
		switch strings.ToLower(ct) {
		case "application/json":
			byt, err := io.ReadAll(req.Body)
			if err != nil {
				_ = req.Body.Close()
				return err
			}
			_ = req.Body.Close()

			reqBody := make(RequestBodyMap)
			if err := json.Unmarshal(byt, &reqBody); err != nil {
				return err
			}
			bodyStr, err := reqBody.SortToString("&")
			if err != nil {
				return fmt.Errorf("SortToString %v", err)
			}
			rawStr = bodyStr
			req.Body = io.NopCloser(bytes.NewBuffer(byt))
		case "multipart/form-data":
			rawStr, err = SortParamForm(req)
			if err != nil {
				return err
			}
		}

	case http.MethodGet:
		rawStr, err = SortParamForm(req)
		if err != nil {
			return err
		}
	}
	rawStr = rawStr + timestamp
	signStrDist := SignHash(apiSign.Method, []byte(rawStr), []byte(apiSign.AccessSecret))
	if signStrDist != signature {
		return fmt.Errorf("sign method invalid rawStr:%s", rawStr)
	}
	return nil
}

// SortParamForm sort and format  URL | form-data param
func SortParamForm(req *http.Request) (string, error) {
	resource := ""
	switch filterFlags(req.Header.Get("Content-Type")) {
	case "multipart/form-data":
		err := req.ParseMultipartForm(10 << 20)
		if err != nil {
			return "", err
		}
	default:
		err := req.ParseForm()
		if err != nil {
			return "", err
		}
	}

	var paramNames []string
	if req.Form != nil && len(req.Form) > 0 {
		for k := range req.Form {
			paramNames = append(paramNames, k)
		}
		sort.Strings(paramNames)

		var query []string
		for _, k := range paramNames {
			query = append(query, url.QueryEscape(k)+"="+url.QueryEscape(req.Form.Get(k)))
		}
		resource = strings.Join(query, "&")
	}

	return resource, nil
}

// SignHash hash and encode
func SignHash(method Method, rawStr, secretKey []byte) (hash string) {
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

type RequestBodyMap map[string]interface{}

func (r RequestBodyMap) GetStringValue(key string) (string, error) {
	val, ok := r[key]
	if !ok {
		return "", fmt.Errorf("request body miss %s key", key)
	}
	v, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("request body %s key not string type", key)
	}
	return v, nil
}

// SortToString request body param sort format
func (r RequestBodyMap) SortToString(separator string) (string, error) {
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
			s = append(s, fmt.Sprintf("%s=%s", v.Key, string(r.TrimNewline(buffer.Bytes()))))
		}
	}
	return strings.Join(s, separator), nil
}

func (r RequestBodyMap) TrimNewline(buf []byte) []byte {
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

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}
