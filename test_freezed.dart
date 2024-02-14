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
