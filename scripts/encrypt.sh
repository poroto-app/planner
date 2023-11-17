# 同じ階層に planner と infrastructure　があることを想定しています。
# - directory
#   - planner
#   - infrastructure

function encrypt() {
  environment=$1
  file=$2

  echo "[${environment}] Encrypting ${file} ..."
  gcloud kms encrypt \
    --location "asia-northeast1" \
    --keyring "planner_api_key_ring" \
    --key "planner_api_crypt_key" \
    --plaintext-file "${file}" \
    --ciphertext-file "../infrastructure/roles/app/planner/${environment}/${file}.enc" || exit 1
  echo ">> [${environment}] Encrypted ${file}!"
}

# 環境設定（デフォルトはdevelopment）
env=$1
if [ -z "${env}" ]; then
  env="development"
fi

# 共通のシークレットを復号化
secrets=(".env.local")
for file in "${secrets[@]}"; do
  encrypt "production" "${file}"
done

case "$env" in
"development")
  secrets_variable=(".env.development.local")
  ;;
"staging")
  secrets_variable=(".env.staging.local")
  ;;
"production")
  secrets_variable=(".env.production.local")
  ;;
esac

for file in "${secrets_variable[@]}"; do
  encrypt "${env}" "${file}"
done
