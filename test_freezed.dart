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
