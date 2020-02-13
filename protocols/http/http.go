// http this simple http protocol.
// Only supports some methods, including: get, post, put, head, options and delete, etc..
package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"goes/lib"
	"goes/protocols"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

// methods the method of http request.
var methods = [...]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}

// httpCode the http code of response.
var httpCode = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	306: "(Unused)",
	307: "Temporary Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Request Entity Too Large",
	414: "Request-URL Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	422: "Unprocessable Entity",
	423: "Locked",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "Http Version Not Supported",
}

// Http struct of http.
type Http struct {
	// header header key-value.
	header               map[string]interface{}
	gzip                 bool
	sessionPath          string
	sessionName          string
	sessionGcProbability int
	sessionGcMaxLifeTime int
	sessionStarted       bool
	// sessionFile session file.
	sessionFile string
	// post post data.
	//post map[string]interface{}
	post *url.Values
}

// GetPost return the data of post method.
func (h *Http) GetPost() url.Values {
	return *h.post
}

// Input check the integrity of the package.
func (h *Http) Input(recvBuffer []byte) interface{} {
	position := bytes.IndexAny(recvBuffer, "\r\n\r\n")
	if position == -1 {
		if len(recvBuffer) >= protocols.MaxPackageSize {
			return false
		}
		return 0
	}

	header := bytes.SplitN(recvBuffer, []byte("\r\n\r\n"), 2)[0]
	method := string(bytes.Split(header, []byte(" "))[0])

	for _, v := range methods {
		if v == method {
			return h.getRequestSize(header, method)
		}
	}
	return 0
}

// Encode http encode, the type of return value is string.
func (h *Http) Encode(data []byte) interface{} {
	// default http-code.
	var header string
	if _, ok := h.header["Http-Code"]; !ok {
		header = "HTTP/1.1 200 OK\r\n"
	} else {
		header = h.header["Http-Code"].(string) + "\r\n"
		delete(h.header, "Http-Code")
	}

	// content-type.
	if _, ok := h.header["Content-Type"]; !ok {
		header += "Content-Type: text/html;charset=utf-8\r\n"
	}

	// other headers key-value.
	for k, v := range h.header {
		if k == "Set-Cookie" && reflect.TypeOf(v).Kind() == reflect.Map {

		} else {
			header += v.(string) + "\r\n"
		}
	}
	if h.gzip {
		header += "Content-Encoding: gzip\r\n"
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err := zw.Write(data)
		if err != nil {
			lib.Error("gzip content error: %s", err)
		}
		if err := zw.Close(); err != nil {
			lib.Warn(err.Error())
		}
		data = buf.Bytes()
	}

	// header.
	header += "Server: Goes/0.1" + "\r\nContent-Length: " + strconv.Itoa(len(data)) + "\r\n\r\n"

	// the whole http package.
	return header + string(data)
}

// Decode parse POST, GET, COOKIE.
func (h *Http) Decode(recvBuffer []byte) []byte {
	if h.header == nil {
		h.header = make(map[string]interface{})
	}
	if h.post == nil {
		h.post = new(url.Values)
	}
	h.header["Connection"] = "Connection: keep-alive"
	server := map[string]interface{}{
		"QUERY_STRING":         "",
		"REQUEST_METHOD":       "",
		"REQUEST_URL":          "",
		"SERVER_PROTOCOL":      "",
		"SERVER_SOFTWARE":      "Goes",
		"SERVER_NAME":          "",
		"HTTP_HOST":            "",
		"HTTP_USER_AGENT":      "",
		"HTTP_ACCEPT":          "",
		"HTTP_ACCEPT_LANGUAGE": "",
		"HTTP_ACCEPT_ENCODING": "",
		"HTTP_COOKIE":          "",
		"HTTP_CONNECTION":      "",
		"CONTENT_TYPE":         "",
		"REMOTE_ADDR":          "",
		"REMOTE_PORT":          "",
		"REQUEST_TIME":         time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
	}

	// parse header.
	packages := bytes.SplitN(recvBuffer, []byte("\r\n\r\n"), 2)
	header, body := packages[0], packages[1]
	headerData := bytes.Split(header, []byte("\r\n"))
	headerData2 := bytes.Split(headerData[0], []byte(" "))
	server["REQUEST_METHOD"], server["REQUEST_URL"], server["REQUEST_PROTOCOL"] = string(headerData2[0]), string(headerData2[1]), string(headerData2[2])

	// parse general header.
	httpPostBoundary := ""
	for i, v := range headerData {
		if i == 0 || len(v) == 0 {
			continue
		}
		keyValue := bytes.SplitN(v, []byte(":"), 2)
		key := string(bytes.ToUpper(bytes.Replace(keyValue[0], []byte("-"), []byte("_"), 1)))
		value := bytes.Trim(keyValue[1], " ")
		server["HTTP_"+key] = string(value)
		switch key {
		// HTTP_HOST.
		case "HOST":
			tmp := bytes.Split(value, []byte(":"))
			server["SERVER_NAME"] = string(tmp[0])
			if len(tmp) > 1 {
				server["SERVER_PORT"] = string(tmp[1])
			}
		// cookie.
		case "COOKIE":
			server["HTTP_COOKIE"] = string(bytes.Replace(value, []byte("; "), []byte("&"), -1))
		// content-type.
		case "CONTENT_TYPE":
			re := regexp.MustCompile("boundary=\\S+")
			matches := re.FindSubmatch(v)
			if matches == nil {
				if pos := bytes.IndexAny(v, ";"); pos != -1 {
					server["CONTENT_TYPE"] = value[:pos]
				} else {
					server["CONTENT_TYPE"] = string(value)
				}
			} else {
				server["CONTENT_TYPE"] = "multipart/form-data"
				httpPostBoundary = "--" + string(matches[0])
			}
		// content-length.
		case "CONTENT_LENGTH":
			server["CONTENT_LENGTH"] = string(value)
		// upgrade.
		case "upgrade":
		// query_string.
		case "REFERER":
			values, _ := url.Parse(string(v))
			server["QUERY_STRING"] = values.Query().Encode()
		default:
		}
	}
	if v, ok := server["HTTP_ACCEPT_ENCODING"]; ok && v == "gzip" {
		h.gzip = true
	}

	// parse post.
	if server["REQUEST_METHOD"] == "POST" {
		if server["CONTENT_TYPE"] != "" {
			switch server["CONTENT_TYPE"] {
			case "multipart/form-data":
				h.parseUploadFile(body, httpPostBoundary)
			case "application/json":
				err := json.Unmarshal(body, &h.post)
				if err != nil {
					lib.Warn("decode data of post error: %v", err)
				}
			case "application/x-www-form-urlencoded":
				*h.post, _ = url.ParseQuery(string(body))
			}
		}
	}

	result, err := json.Marshal(map[string]interface{}{
		"server": server,
		"post":   h.post,
	})
	if err != nil {
		lib.Warn("Unable marshal data, err: %v", err)
		return []byte{}
	}
	return result
}

// parseUploadFile parse file.
func (h *Http) parseUploadFile(body []byte, httpPostBoundary string) {

}

// getRequestSize get whole size of the request.
// includes the request headers and request body.
func (h *Http) getRequestSize(header []byte, method string) int {
	if method == "GET" || method == "OPTIONS" || method == "HEAD" || method == "DELETE" {
		return len(header) + 4
	}

	re := regexp.MustCompile("Content-Length: \\d+")
	matches := re.FindSubmatch(header)
	if matches != nil {
		contentLength := bytes.Trim(bytes.Split(matches[0], []byte(":"))[1], " ")
		v, _ := strconv.ParseInt(string(contentLength), 10, 32)
		return len(header) + int(v) + 4
	}

	return 0
}

// NewHttpProtocol
func NewHttpProtocol() *Http {
	return &Http{}
}
