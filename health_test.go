package health

import "testing"

func TestNewHealth(t *testing.T) {
	h := NewHealth()

	if !h.IsUnknown() {
		t.Errorf("NewHealth().IsDown() == %t, want %t", h.IsUnknown(), true)
	}
}

func Test_Health_Unknown(t *testing.T) {
	h := NewHealth()
	h.Unknown()

	if !h.IsUnknown() {
		t.Errorf("NewHealth().IsUnknown() == %t, want %t", h.IsUnknown(), true)
	}
}

func Test_Health_Up(t *testing.T) {
	h := NewHealth()
	h.Up()

	if !h.IsUp() {
		t.Errorf("NewHealth().IsUp() == %t, want %t", h.IsUp(), true)
	}
}

func Test_Health_Down(t *testing.T) {
	h := NewHealth()
	h.Up()
	h.Down()

	if !h.IsDown() {
		t.Errorf("NewHealth().IsDown() == %t, want %t", h.IsDown(), true)
	}
}

func Test_Health_OutOfService(t *testing.T) {
	h := NewHealth()
	h.OutOfService()

	if !h.IsOutOfService() {
		t.Errorf("NewHealth().IsOutOfService() == %t, want %t", h.IsOutOfService(), true)
	}
}

func Test_Health_IsUp(t *testing.T) {
	h := NewHealth()
	h.Up()

	if !h.IsUp() {
		t.Errorf("NewHealth().status == %v, want %s", h.IsUp(), up)
	}
}

func Test_Health_IsDown(t *testing.T) {
	h := NewHealth()
	h.Down()

	if !h.IsDown() {
		t.Errorf("NewHealth().status == %v, want %s", h.IsDown(), down)
	}
}

func Test_Health_IsOutOfService(t *testing.T) {
	h := NewHealth()
	h.OutOfService()

	if !h.IsOutOfService() {
		t.Errorf("NewHealth().status == %v, want %s", h.IsOutOfService(), outOfService)
	}
}

func Test_Health_MarshalJSON(t *testing.T) {
	h := NewHealth()
	h.Up()

	h.AddInfo("status", "Should not render")

	json, err := h.MarshalJSON()

	if err != nil {
		t.Errorf("err != nil, wants nil")
	}

	expected := `{"status":"UP"}`

	if string(json) != expected {
		t.Errorf("h.MarshalJSON() == %s, wants %s", string(json), expected)
	}
}

func Test_Health_AddInfo(t *testing.T) {
	h := NewHealth()

	h.AddInfo("key", "value")

	val, ok := h.GetInfo("key").(string)

	if !ok {
		t.Error("h.AddInfo() should add a key value to the map")
		return
	}
	if val != "value" {
		t.Errorf("h.AddInfo() expect value got %s", val)
	}
}

func Test_Health_AddInfo_null_map(t *testing.T) {
	h := health{}

	h.AddInfo("key", "value")

	_, ok := h.info["key"]

	if !ok {
		t.Error("h.AddInfo() should add a key value to the map")
	}
}

func Test_Health_GetInfo(t *testing.T) {
	h := NewHealth()

	h.AddInfo("key", "value")

	value := h.GetInfo("key")

	if value != "value" {
		t.Errorf(`h.GetInfo("key") == %s, wants %s`, value, "value")
	}
}

func Test_Health_GetInfo_null_map(t *testing.T) {
	h := health{}

	value := h.GetInfo("key")

	if value != nil {
		t.Errorf(`h.GetInfo("key") == %v, wants %v`, value, nil)
	}
}
