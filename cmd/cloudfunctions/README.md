# Cloud Functions

### テスト方法
SEE: [関数をローカルでビルドしてテストする](https://cloud.google.com/functions/docs/create-deploy-http-go?hl=ja#build_and_test_your_function_locally)  

**`HelloHTTP`という関数をテストする場合**
以下の方法でサーバを起動し`https://localhost:8080`にアクセスすると関数が実行される。

```shell
export FUNCTION_TARGET=HelloHTTP
go run cmd/cloudfunctions/main.go
```
