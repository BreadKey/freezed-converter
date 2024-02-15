# Example
```cd cli;go run . ../test_freezed.dart```
## Input
```dart
@freezed
class TestData with _$TestData {
  const factory TestData({
    required String id,
    required int a,
    required String b,
    required double c,
    int? d,
    @JsonKey(includeToJson: false) int? e,
    String? refId,
    required bool idle,
  }) = _TestData;
}
```
## Output
```go
type TestData struct {
	ID    string  `firestore:"id"`
	A     int     `firestore:"a"`
	B     string  `firestore:"b"`
	C     float64 `firestore:"c"`
	D     *int    `firestore:"d,omitempty"`
	RefID *string `firestore:"refId,omitempty"`
	Idle  bool    `firestore:"idle"`
}
```
# Supported Languages
## Go
