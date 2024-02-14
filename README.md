# Example
## Input
```dart
@freezed
class TestData with _$TestData {
  const factory TestData({
    required int a,
    required String b,
    required double c,
    int? d,
    @JsonKey(includeToJson: false)
    int? e,
  }) = _TestData;
}
```
## Output
```go
type TestData struct {
	A int     `firestore:"a"`
	B string  `firestore:"b"`
	C float64 `firestore:"c"`
	D *int    `firestore:"d,omitempty"`
}
```