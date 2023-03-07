# httpx

## APISign API参数签名验证
-  1. rawStr eg: GET http://example.com/hello?n=1&a=2  Key["n","a"]-ASC Sort["a","n"] GetParam(a) a=2&n=1 param key string attaches the methods
-  2. other request http method，for Content-Type: application/json {"n":"m","a":2} Key ASC Sort,param key string attaches the methods => {"a":2,"n":"m"} => a=2&n=m
-  3. rawStr+timestamp => a=2&n=m1626167650 (1626167650 is unix timestamp), verify sign time valid（default 10s）
-  4. Sign Method： Method(rawStr+timestamp, secretKey) signed text encode [Base64, Hex(default)]
//     Method=[HMAC-SHA256,HMAC-SHA1] Encode=[Base64,Hex] Default = HMAC-SHA256-HEX
-  5. default: signStr=Hex(HMAC-SHA256(rawStr+timestamp,secretKey))
-  6. Sign http request Header X-Signature=accessKey:signStr:timestamp (: split elem)

### USEAGE
see test example