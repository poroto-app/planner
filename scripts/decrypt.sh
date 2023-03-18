# ローカルで実行するときはこの行のコメントアウトを外す
#git clone --depth 1 git@github.com:poroto-app/infrastructure.git

gcloud kms decrypt \
  --location "asia-northeast1" \
  --keyring "planner_api_key_ring" \
  --key "planner_api_crypt_key" \
  --plaintext-file ./.env.local \
  --ciphertext-file ../infrastructure/roles/app/planner/production/.env.local.enc

rm -rf infrastructure