package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)
// データを保持するためのマップ
var data = make(map[string][]byte)
// キーの文字列パターンをチェックする正規表現
var validKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func main() {
	// /objects/パスに対するリクエストハンドラを登録
	http.HandleFunc("/objects/", handleObjectsRequest)
	// サーバが起動したことをコンソールに出力
	fmt.Println("Server listening on port 8000")
	// ポート8000でサーバを起動
	err := http.ListenAndServe(":8000", nil)
	// エラー発生時の処理
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func handleObjectsRequest(w http.ResponseWriter, r *http.Request) {
	// リクエストURLからキー部分を抽出
	key := strings.TrimPrefix(r.URL.Path, "/objects/")

	// キーの長さと正規表現のチェック
	if len(key) > 10 || !validKeyRegex.MatchString(key) {
		// 制約に合わない場合は404を返す
		http.NotFound(w, r)
		return
	}

// リクエストメソッドに応じた処理
	switch r.Method {
	case http.MethodPut:
		// リクエストボディを読み取り
		body, err := ioutil.ReadAll(r.Body)
		
		if err != nil {
			// エラー発生時の処理
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// データをマップに保存
		data[key] = body
		// 成功のステータスコードを返す
		w.WriteHeader(http.StatusOK)


	case http.MethodGet:
		// マップからデータを取得
		value, ok := data[key]
		// データがない場合は404を返す
		if !ok {
			http.NotFound(w, r)
			return
		}
		// データをレスポンスボディに書き込む
		w.Write(value)

	default:
		// その他のメソッドの場合は405を返す
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}