# 同じ階層に planner と infrastructure　があることを想定しています。
# - directory
#   - planner
#   - infrastructure

gcloud kms encrypt \
  --location "asia-northeast1" \
  --keyring "planner_api_key_ring" \
  --key "planner_api_crypt_key" \
  --plaintext-file ./.env.local \
  --ciphertext-file ../infrastructure/roles/app/planner/production/.env.local.enc