package freezedconverter

import "testing"

func TestToken(t *testing.T) {
	tokenizer := NewTokenizer("test_freezed.dart")

	expectedTokens := []string{
		"@freezed",
		"class",
		"TestData",
		"with",
		"_$TestData",
		"{",
		"const",
		"factory",
		"TestData",
		"(",
		"{",
		"required",
		"String",
		"id",
		"required",
		"int",
		"a",
		"required",
		"String",
		"b",
		"required",
		"double",
		"c",
		"int?",
		"d",
		"@JsonKey",
		"(",
		"includeToJson",
		"false",
		")",
		"int?",
		"e",
		"String?",
		"refId",
		"required",
		"bool",
		"idle",
		"List<int>",
		"list",
		"}",
		")",
		"=",
		"_TestData",
		"}",
	}

	for _, token := range expectedTokens {
		actual := tokenizer.Next()

		if actual != token {
			t.Errorf("Wrong result! %s expected but %s!", token, actual)
		}
	}
}
