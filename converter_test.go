package freezedconverter

import "testing"

func TestParse(t *testing.T) {
	result := Parse("test_freezed.dart")

	if len(result) != 1 {
		t.Errorf("Wrong result! %d", len(result))
	}

	freezed := result[0]

	if freezed.Name != "TestData" {
		t.Errorf("Wrong name! %s", freezed.Name)
	}

	if len(freezed.Parameters) != 3 {
		t.Errorf("Wrong parameter length! %d\n", len(freezed.Parameters))

		for _, p := range freezed.Parameters {
			t.Errorf("%s, %s\n", p.Type, p.Name)
		}
	}

	expectedParameters := [3]ParameterToken{
		{
			Type: "int",
			Name: "a",
		},
		{
			Type: "String",
			Name: "b",
		},
		{
			Type: "double",
			Name: "c",
		},
	}

	for i, p := range freezed.Parameters {
		expected := expectedParameters[i]

		if expected.Type != p.Type || expected.Name != p.Name {
			t.Errorf("Wrong parameter! expected: %s, %s, actual: %s, %s", expected.Type, expected.Name, p.Type, p.Name)
		}
	}
}
