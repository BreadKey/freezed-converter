import 'test.dart';

final a = "test{Hello}";

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
    List<int>? list,
  }) = _TestData;
}
