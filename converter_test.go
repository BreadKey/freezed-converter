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

	expectedParameters := [8]ParameterToken{
		{
			Type: "String",
			Name: "id",
		},
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
		{
			Type:     "int",
			Name:     "d",
			Nullable: true,
		},
		{
			Type:     "String",
			Name:     "refId",
			Nullable: true,
		},
		{
			Type: "bool",
			Name: "idle",
		},
		{
			Type:     "List<int>",
			Name:     "list",
			Nullable: true,
		},
	}

	if len(freezed.Parameters) != len(expectedParameters) {
		t.Errorf("Wrong parameter length! %d\n", len(freezed.Parameters))

		for _, p := range freezed.Parameters {
			t.Errorf("%s, %s\n", p.Type, p.Name)
		}
	}

	for i, p := range freezed.Parameters {
		expected := expectedParameters[i]

		if expected.Type != p.Type || expected.Name != p.Name {
			t.Errorf("Wrong parameter! expected: %s, %s, actual: %s, %s", expected.Type, expected.Name, p.Type, p.Name)
		}
	}
}

func TestTranslateToGo(t *testing.T) {
	result := Parse("test_freezed.dart")

	freezed := result[0]

	expected := `type TestData struct {
	ID    string  ` + "`firestore:\"id\"`" + `
	A     int     ` + "`firestore:\"a\"`" + `
	B     string  ` + "`firestore:\"b\"`" + `
	C     float64 ` + "`firestore:\"c\"`" + `
	D     *int    ` + "`firestore:\"d,omitempty\"`" + `
	RefID *string ` + "`firestore:\"refId,omitempty\"`" + `
	Idle  bool    ` + "`firestore:\"idle\"`" + `
	List  *?      ` + "`firestore:\"list,omitempty\"`" + `
}`

	translated := TranslateToGo(&freezed)

	if translated != expected {
		t.Errorf("Wront result! %v", translated)
	}
}
