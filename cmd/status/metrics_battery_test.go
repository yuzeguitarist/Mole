package main

import (
	"math"
	"testing"
)

func TestParseAppleSmartBatteryThermalKeepsBatteryTemperatureOutOfCPUTemp(t *testing.T) {
	out := `
  | |   "Temperature" = 3055
  | |   "SystemPowerIn" = 19967
  | |   "BatteryPower" = 13654
  | |   "AdapterDetails" = {"Watts" = 96}
`

	thermal := parseAppleSmartBatteryThermal(out)

	if thermal.CPUTemp != 0 {
		t.Fatalf("expected cpu temp to stay unset, got %v", thermal.CPUTemp)
	}
	if math.Abs(thermal.BatteryTemp-30.55) > 0.001 {
		t.Fatalf("expected battery temp 30.55, got %v", thermal.BatteryTemp)
	}
	if math.Abs(thermal.SystemPower-19.967) > 0.001 {
		t.Fatalf("expected system power 19.967W, got %v", thermal.SystemPower)
	}
	if thermal.AdapterPower != 96 {
		t.Fatalf("expected adapter power 96W, got %v", thermal.AdapterPower)
	}
	if math.Abs(thermal.BatteryPower-13.654) > 0.001 {
		t.Fatalf("expected battery power 13.654W, got %v", thermal.BatteryPower)
	}
	if math.Abs(thermal.CurrentPower-19.967) > 0.001 {
		t.Fatalf("expected current power 19.967W, got %v", thermal.CurrentPower)
	}
	if thermal.PowerSource != "system" {
		t.Fatalf("expected system power source, got %q", thermal.PowerSource)
	}
}

func TestParseAppleSmartBatteryThermalParsesTwosComplementBatteryPower(t *testing.T) {
	out := `
  | |   "BatteryPower"=18446744073709539271
`

	thermal := parseAppleSmartBatteryThermal(out)

	if math.Abs(thermal.BatteryPower-(-12.345)) > 0.001 {
		t.Fatalf("expected battery power -12.345W, got %v", thermal.BatteryPower)
	}
	if math.Abs(thermal.CurrentPower-12.345) > 0.001 {
		t.Fatalf("expected current charging power 12.345W, got %v", thermal.CurrentPower)
	}
	if thermal.PowerSource != "charging" {
		t.Fatalf("expected charging power source, got %q", thermal.PowerSource)
	}
}

func TestParseAppleSmartBatteryThermalDerivesBatteryWattsFromVoltageAndAmperage(t *testing.T) {
	out := `
  | |   "Voltage" = 12000
  | |   "InstantAmperage" = -1500
`

	thermal := parseAppleSmartBatteryThermal(out)

	if math.Abs(thermal.BatteryPower-18.0) > 0.001 {
		t.Fatalf("expected derived battery power 18W, got %v", thermal.BatteryPower)
	}
	if math.Abs(thermal.CurrentPower-18.0) > 0.001 {
		t.Fatalf("expected current power 18W, got %v", thermal.CurrentPower)
	}
	if thermal.PowerSource != "battery" {
		t.Fatalf("expected battery power source, got %q", thermal.PowerSource)
	}
}

func TestParsePowermetricsPowerPrefersCombinedPower(t *testing.T) {
	out := `
CPU Power: 1200 mW
GPU Power: 300 mW
ANE Power: 100 mW
Combined Power (CPU + GPU + ANE): 1900 mW
`

	thermal := parsePowermetricsPower(out)

	if math.Abs(thermal.CurrentPower-1.9) > 0.001 {
		t.Fatalf("expected combined current power 1.9W, got %v", thermal.CurrentPower)
	}
	if math.Abs(thermal.SystemPower-1.9) > 0.001 {
		t.Fatalf("expected system power 1.9W, got %v", thermal.SystemPower)
	}
	if thermal.PowerSource != "powermetrics" {
		t.Fatalf("expected powermetrics source, got %q", thermal.PowerSource)
	}
}

func TestParsePowermetricsPowerSumsComponentsWhenCombinedMissing(t *testing.T) {
	out := `
CPU Power: 1.5 W
GPU Power: 500 mW
ANE Power: 250 mW
`

	thermal := parsePowermetricsPower(out)

	if math.Abs(thermal.CurrentPower-2.25) > 0.001 {
		t.Fatalf("expected summed current power 2.25W, got %v", thermal.CurrentPower)
	}
}
