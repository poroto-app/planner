# おすすめの場所
### 概要
場所を指定してプランを作成するときの候補として、ユーザーに提示する場所を管理する

### 利用方法
1. 登録したい場所の名前を確認する
```sh
 go run cmd/features/recommended_places_for_plan/main.go \
  -place <場所のID> 
```

2. 場所を登録する
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