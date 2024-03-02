package array

// Find はスライスと条件関数を引数に取り、スライスの各要素を引数に条件関数を実行
// 条件を満たす値を返し、満たさない値はnilで返す
func Find[T any](slice []T, condition func(T) bool) (T, bool) {
	var zero T
	for _, v := range slice {
		if condition(v) {
			return v, true
		}
	}
	return zero, false
}
