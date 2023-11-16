## decrypt.sh
infrastructureリポジトリのシークレットを復号化するスクリプト

開発環境で利用するシークレットを復号化
```shell
bash scripts/decrypt.sh
```

プロダクション環境で利用するシークレットを復号化
```shell
bash scripts/decrypt.sh production
```
### encrypt.sh
planner リポジトリのシークレットを暗号化するスクリプト

開発環境で利用するシークレットを暗号化
```shell
bash scripts/encrypt.sh
```

プロダクション環境で利用するシークレットを暗号化
```shell
bash scripts/encrypt.sh production
```