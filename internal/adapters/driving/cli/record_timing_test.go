package cli_test

import (
	"strings"
	"testing"
)

// Abdeckung: LH-FA-21/Happy, LH-FA-17/Boundary — --record-timing, seine
// Umgebungsvariable und der Schlüssel record_timing im Abschnitt record: setzen
// RecordTiming; ohne alle drei ist es false (LH-FA-21.a *Aufzeichnen*, LH-FA-17.a).
func TestRecordTimingQuellen(t *testing.T) {
	leere(t, "record")
	args := basis("record", "")
	cmd, err := lese(args...)
	if err != nil || cmd.Record.RecordTiming {
		t.Fatalf("ohne Option: RecordTiming %v, %v", cmd.Record.RecordTiming, err)
	}
	datei := schreibe(t, "record:\n  record_timing: true\n")
	for _, extra := range []string{"--record-timing", "--record-timing=true", "--config=" + datei} {
		cmd, err := lese(append(append([]string{}, args...), extra)...)
		if err != nil || !cmd.Record.RecordTiming {
			t.Errorf("%s: RecordTiming %v, %v", extra, cmd.Record.RecordTiming, err)
		}
	}
	t.Setenv("PGWIRE_RECORDER_RECORD_TIMING", "true")
	if cmd, err := lese(args...); err != nil || !cmd.Record.RecordTiming {
		t.Errorf("Umgebungsvariable: RecordTiming %v, %v", cmd.Record.RecordTiming, err)
	}
	if cmd, err := lese(append(append([]string{}, args...), "--record-timing=false")...); err != nil || cmd.Record.RecordTiming {
		t.Errorf("Option false vor Umgebungsvariable true: RecordTiming %v, %v", cmd.Record.RecordTiming, err)
	}
}

// Abdeckung: LH-FA-01/Happy, LH-FA-21/Happy — die Hilfe von record nennt
// --record-timing und PGWIRE_RECORDER_RECORD_TIMING, die von replay und play
// nicht (LH-FA-21.a *Aufzeichnen*, LH-FA-01.a).
func TestRecordTimingHilfe(t *testing.T) {
	for kommando, erwartet := range map[string]bool{"record": true, "replay": false, "play": false} {
		text := hilfe(t, kommando, "--help")
		for _, teil := range []string{"--record-timing", "PGWIRE_RECORDER_RECORD_TIMING"} {
			if strings.Contains(text, teil) != erwartet {
				t.Errorf("%s: Hilfe nennt %q: %v, erwartet %v", kommando, teil, !erwartet, erwartet)
			}
		}
	}
}
