# Example
```cd cli;go run . ../test_freezed.dart```
## @freezed
### Input
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
### Output
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
## l10n
### Input
```json
{
    "helloPet": "Hello {species, select, dog{Dog} cat{Cat} others{Others} other{?}} World!"
}
```
### Output
```go
type Species int
const (
	Dog Species = iota
	Cat
	Others
)
func SpeciesToL10n(species Species) string {
	switch {
	case Dog:
		return "dog"
	case Cat:
		return "cat"
	case Others:
		return "others"
	default:
		return "other"
	}
}

func HelloPet(species Species) string {
	switch species {
	case Dog:
		return "Hello Dog World!"
	case Cat:
		return "Hello Cat World!"
	case Others:
		return "Hello Others World!"
	default:
		return "?"
	}
}
```
# Supported Languages
## Go

# Supported Formats
## firestore
## json