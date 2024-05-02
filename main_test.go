package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleObjectsRequest(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name           string
		method         string
		url            string
		body           []byte
		expectedStatus int
		expectedBody   []byte
	}{
		// PUT リクエストのテスト
		{
			name:           "PUT valid key",
			method:         http.MethodPut,
			url:            "/objects/validkey",
			body:           []byte("test data"),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "PUT invalid key (too long)",
			method:         http.MethodPut,
			url:            "/objects/invalidkeylong",
			body:           []byte("test data"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "PUT invalid key (regex)",
			method:         http.MethodPut,
			url:            "/objects/invalid!key",
			body:           []byte("test data"),
			expectedStatus: http.StatusNotFound,
		},

		// GET リクエストのテスト
		{
			name:           "GET existing key",
			method:         http.MethodGet,
			url:            "/objects/validkey",
			expectedStatus: http.StatusOK,
			expectedBody:   []byte("test data"),
		},
		{
			name:           "GET non-existing key",
			method:         http.MethodGet,
			url:            "/objects/nonexistingkey",
			expectedStatus: http.StatusNotFound,
		},

		// その他のメソッドのテスト
		{
			name:           "Other method",
			method:         http.MethodDelete,
			url:            "/objects/validkey",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	// テストケースの実行
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// リクエストの作成
			req, err := http.NewRequest(test.method, test.url, bytes.NewBuffer(test.body))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			// レスポンスレコーダの作成
			rr := httptest.NewRecorder()

			// ハンドラの呼び出し
			handler := http.HandlerFunc(handleObjectsRequest)
			handler.ServeHTTP(rr, req)

			// ステータスコードのチェック
			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("unexpected status code: got %v want %v", status, test.expectedStatus)
			}

			// レスポンスボディのチェック
			if test.expectedBody != nil {
				body, err := ioutil.ReadAll(rr.Body)
				if err != nil {
					t.Fatalf("Error reading response body: %v", err)
				}
				if !bytes.Equal(body, test.expectedBody) {
					t.Errorf("unexpected response body: got %v want %v", body, test.expectedBody)
				}
			}
		})
	}
}