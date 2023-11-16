function decrypt() {
  environment=$1
  file=$2

  echo "[${environment}] Decrypting ${file} ..."
  gcloud kms decrypt \
    --location "asia-northeast1" \
    --keyring "planner_api_key_ring" \
    --key "planner_api_crypt_key" \
    --plaintext-file "${file}" \
    --ciphertext-file "./infrastructure/roles/app/planner/${environment}/${file}.enc" || exit 1
  echo ">> [${environment}] Decrypted ${file}!"
}

# infrastructureリポジトリをclone
if [ ! -d infrastructure ]; then
  git clone --depth 1 git@github.com:poroto-app/infrastructure.git
fi

# 環境設定（デフォルトはdevelopment）
env=$1
if [ -z "${env}" ]; then
  env="development"
fi

# 共通のシークレットを復号化
secrets=(".env.local")
for file in "${secrets[@]}"; do
  decrypt "production" "${file}"
done

# 環境ごとのシークレットを復号化
if [ ! -d secrets ]; then
  mkdir secrets
fi

case "$env" in
"development")
  secrets_variable=("secrets/google-credential.json" ".env.development.local")
  ;;
"staging")
  secrets_variable=("secrets/google-credential.json" ".env.staging.local")
  ;;
"production")
  secrets_variable=("secrets/google-credential.json" ".env.production.local")
  ;;
esac

for file in "${secrets_variable[@]}"; do
  decrypt "${env}" "${file}"
done

rm -rf infrastructure
