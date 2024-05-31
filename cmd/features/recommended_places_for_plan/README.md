# おすすめの場所
### 概要
場所を指定してプランを作成するときの候補として、ユーザーに提示する場所を管理する

### 利用方法
1. 登録したい場所を名前で検索する
```sh
 go run cmd/features/recommended_places_for_plan/main.go \
  -name "<場所の名前>"
```
1. 登録したい場所の名前を確認する
```sh
 go run cmd/features/recommended_places_for_plan/main.go \
  -place <場所のID> 
```

2. 場所を登録する  
- ユーザに提示するときに必ず画像と一緒に提示するために、画像が紐づいていない場所は登録できないようになっています。
```sh
 go run cmd/features/recommended_places_for_plan/main.go \
  -place <場所のID> \
  -register 
```

3. 登録した場所を削除する
```sh
 go run cmd/features/recommended_places_for_plan/main.go \
  -place <場所のID> \
  -delete 
```