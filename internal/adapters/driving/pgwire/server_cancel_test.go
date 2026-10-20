package pgwire_test

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/testpki"
)

// cancelNachricht ist ein CancelRequest der Länge laenge im Kopf, dessen Rest
// geliefert Bytes mit dem Wert fuell hat.
func cancelNachricht(laenge, geliefert int, fuell byte) []byte {
	return rohKopf(laenge, pgwire.CodeCancelRequest, bytes.Repeat([]byte{fuell}, geliefert)...)
}

// cancelVoll ist ein CancelRequest der Länge laenge mit vollständigem Rest.
func cancelVoll(laenge int) []byte {
	return cancelNachricht(laenge, laenge-8, 0x2a)
}

// cancelErwartung ist, was der Client von einem CancelRequest sieht: Warnung
// PGR-W3001 ohne Antwort, oder eine ErrorResponse mit PGR-E6001 und gemerktem Fehler.
type cancelErwartung struct {
	name  string
	daten []byte
	// halb schließt nach dem Senden die Schreibseite (der Rest wird nicht geliefert).
	halb bool
	// fehler ist PGR-E6001 mit ErrorResponse; sonst gilt die Warnung PGR-W3001 ohne Antwort.
	fehler bool
	// restUngelesen: Der Client hat mehr Bytes gesendet, als der Recorder bei der Ablehnung
	// gelesen hat; die ErrorResponse ist dann nicht zugesagt (LH-FA-05.e), geprüft werden das
	// Ende der Verbindung, PGR-E6001 im Log und der gemerkte Fehler.
	restUngelesen bool
}

func cancelFaelle() []cancelErwartung {
	var fs []cancelErwartung
	for _, l := range []int{12, 16, 17, 20, 268} {
		fs = append(fs, cancelErwartung{name: fmt.Sprintf("Länge %d", l), daten: cancelVoll(l)})
	}
	fs = append(fs,
		cancelErwartung{name: "Länge 12, Prozess-ID und Schlüssel 0", daten: cancelNachricht(12, 4, 0)},
		cancelErwartung{name: "Länge 12, Prozess-ID 0xff", daten: cancelNachricht(12, 4, 0xff)},
		cancelErwartung{name: "Länge 16, Schlüssel 0xff", daten: cancelNachricht(16, 8, 0xff)},
	)
	for _, l := range []int{8, 9, 10, 11, 13, 14, 15} {
		fs = append(fs, cancelErwartung{name: fmt.Sprintf("Länge %d", l), daten: cancelVoll(l), fehler: true})
	}
	for _, l := range []int{269, 270, 10000} {
		fs = append(fs, cancelErwartung{name: fmt.Sprintf("Länge %d", l), daten: cancelVoll(l), fehler: true, restUngelesen: true})
	}
	for _, l := range []int{12, 16, 268} {
		for _, g := range []int{0, 1, 3} {
			fs = append(fs, cancelErwartung{name: fmt.Sprintf("Länge %d, %d von %d Bytes Rest", l, g, l-8), daten: cancelNachricht(l, g, 0), halb: true, fehler: true})
		}
	}
	for _, l := range []int{9, 11, 13, 15, 269} {
		fs = append(fs, cancelErwartung{name: fmt.Sprintf("Länge %d, Rest nicht geliefert", l), daten: cancelNachricht(l, 0, 0), halb: true, fehler: true})
	}
	return fs
}

// cancelStelle ist die Stelle in der Folge der Nachrichten: was vor dem CancelRequest
// steht, wie viele „N“ der Client dafür liest, und ob TLS ausgehandelt ist.
type cancelStelle struct {
	name     string
	tls      bool
	klartext bool
	vor      []byte
	n        int
}

func cancelStellen() []cancelStelle {
	gss := rohKopf(8, pgwire.CodeGSSEncRequest)
	return []cancelStelle{
		{name: "erste Nachricht"},
		{name: "nach SSLRequest mit N", vor: sslAnfrage(), n: 1},
		{name: "nach GSSENCRequest mit N", vor: gss, n: 1},
		{name: "über TLS nach S", tls: true},
		{name: "unverschlüsselt bei TLS-Pflicht", tls: true, klartext: true},
		{name: "nach GSSENCRequest bei TLS-Pflicht", tls: true, klartext: true, vor: gss, n: 1},
	}
}

// cancelSenden schickt daten und liest, ohne die Schreibseite zu schließen (außer halb),
// bis der Server die Verbindung schließt; ein Warten bis zur Frist ist ein Fehler. Ein
// Rücksetzen der Verbindung (beim Schreiben ECONNRESET oder EPIPE, beim Lesen
// connection reset) gilt als Ende: Schließt der Server mit ungelesenen Bytes, kann es als
// RST ankommen (LH-FA-05.e, Ablehnung aus dem Kopf bei ungelesenem Rest).
func cancelSenden(r *randServer, daten []byte, halb, klartext bool) []byte {
	r.t.Helper()
	tcp, c := r.verbinde(klartext)
	if _, err := c.Write(daten); err != nil && !errors.Is(err, syscall.ECONNRESET) && !errors.Is(err, syscall.EPIPE) {
		r.t.Fatal(err)
	}
	if halb {
		if tc, ok := c.(*tls.Conn); ok {
			_ = tc.CloseWrite()
		} else {
			_ = tcp.CloseWrite()
		}
	}
	roh, err := io.ReadAll(c)
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		r.t.Fatalf("Verbindung nicht geschlossen, gelesen %q", roh)
	}
	return roh
}

// Abdeckung: LH-FA-05/Negative, LH-FA-05/Boundary — die Länge des CancelRequest
// (LH-FA-05.e *Länge des CancelRequest*): 12 (Prozess-ID ohne Schlüssel) und 16 bis 268
// sind PGR-W3001 ohne Antwort und ohne Fehler, die Verbindung wird geschlossen, es
// entsteht keine Session; 8 bis 11 und 13 bis 15 sind PGR-E6001 mit ErrorResponse und
// ohne Warnung PGR-W3001; 269, 270 und 10000 mit vollem Rest sind PGR-E6001 mit Ende der
// Verbindung und Eintrag im Log, die ErrorResponse wird nicht verlangt (ungelesener Rest,
// Zustellung nicht zugesagt; geprüft nur unter Linux); ein Rest, den die Verbindung nicht ganz
// liefert (Länge 12, 16, 268 mit 0, 1 und 3 Bytes), ist PGR-E6001; Prozess-ID und
// Schlüssel werden nicht ausgewertet (0, 0xff, 0x2a); in record und replay, an jeder
// Stelle der Folge (erste Nachricht, nach SSLRequest mit N, nach GSSENCRequest, über TLS
// nach S, unverschlüsselt bei TLS-Pflicht ohne PGR-E6003).
func TestCancelRequestLaenge(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	for _, replay := range []bool{false, true} {
		for _, st := range cancelStellen() {
			for _, f := range cancelFaelle() {
				t.Run(fmt.Sprintf("replay=%v/%s/%s", replay, st.name, f.name), func(t *testing.T) {
					var zt *pgwire.TLS
					if st.tls {
						zt = z
					}
					r := neuerRandServer(t, replay, zt)
					roh := cancelSenden(r, append(append([]byte{}, st.vor...), f.daten...), f.halb, st.klartext)
					b := r.befund(roh)
					warn := strings.Count(r.log.String(), "level=WARN")
					if f.restUngelesen {
						if b.sessions != 0 || b.erster != model.CodeUnsupported || !strings.Contains(r.log.String(), "code="+model.CodeUnsupported) || strings.Contains(r.log.String(), model.CodeCancelRequest) {
							t.Fatalf("Befund %+v, erwartet PGR-E6001 im Log und als gemerkter Fehler, keine Session, keine Warnung %s:\n%s", b, model.CodeCancelRequest, r.log.String())
						}
						return
					}
					if b.n != st.n || b.sessions != 0 {
						t.Fatalf("Befund %+v, erwartet n=%d, keine Session", b, st.n)
					}
					if f.fehler {
						if b.typen != "E" || !strings.Contains(b.fehler, model.CodeUnsupported) || b.erster != model.CodeUnsupported {
							t.Fatalf("Befund %+v, erwartet ErrorResponse mit %s", b, model.CodeUnsupported)
						}
						if strings.Contains(r.log.String(), model.CodeCancelRequest) {
							t.Fatalf("Warnung %s trotz Fehler:\n%s", model.CodeCancelRequest, r.log.String())
						}
						return
					}
					if b.typen != "" || b.erster != "" {
						t.Fatalf("Befund %+v, erwartet keine Antwort und keinen Fehler", b)
					}
					if warn != 1 || !strings.Contains(r.log.String(), "code="+model.CodeCancelRequest) {
						t.Fatalf("Log ohne genau eine Warnung %s:\n%s", model.CodeCancelRequest, r.log.String())
					}
				})
			}
		}
	}
}

// Abdeckung: LH-FA-05/Boundary — eine Länge außerhalb von 8 bis 10000 ist PGR-W3003, auch
// mit dem Code des CancelRequest: Die Länge wird vor der Länge des CancelRequest geprüft,
// und die Antwort ist leer (LH-FA-05.e *Länge des CancelRequest*, Reihenfolge der Gründe).
func TestCancelRequestLaengeAusserhalb(t *testing.T) {
	for _, replay := range []bool{false, true} {
		for _, l := range []int{0, 4, 7, 10001} {
			t.Run(fmt.Sprintf("replay=%v/Länge %d", replay, l), func(t *testing.T) {
				r := neuerRandServer(t, replay, nil)
				roh := cancelSenden(r, rohKopf(l, pgwire.CodeCancelRequest), false, false)
				b := r.befund(roh)
				if len(roh) != 0 || b.erster != "" || b.sessions != 0 {
					t.Fatalf("Antwort %q, Befund %+v", roh, b)
				}
				log := r.log.String()
				if !strings.Contains(log, "code="+model.CodeForeignProtocol) || strings.Contains(log, model.CodeCancelRequest) {
					t.Fatalf("Log ohne %s oder mit %s:\n%s", model.CodeForeignProtocol, model.CodeCancelRequest, log)
				}
			})
		}
	}
}

// Abdeckung: LH-FA-05/Boundary — ein CancelRequest berührt eine zweite Verbindung nicht
// (LH-FA-05.e *Länge des CancelRequest*, *Modus*): Eine reguläre Verbindung, die schon
// offen ist, sendet nach dem CancelRequest ihre Anfrage und endet regulär, und die nächste reguläre Verbindung
// öffnet die zweite Session; der CancelRequest selbst öffnet keine, in record und replay,
// für die Längen 12 und 16.
func TestCancelRequestZweiteVerbindungUnberuehrt(t *testing.T) {
	for _, replay := range []bool{false, true} {
		for _, l := range []int{12, 16} {
			t.Run(fmt.Sprintf("replay=%v/Länge %d", replay, l), func(t *testing.T) {
				r := neuerRandServer(t, replay, nil)
				_, reg := r.verbinde(false)
				_, _ = reg.Write(rohStart("user", "app"))
				if _, err := io.ReadFull(reg, make([]byte, 1)); err != nil {
					t.Fatal(err)
				}
				roh := cancelSenden(r, cancelVoll(l), false, false)
				if len(roh) != 0 || r.zaehler() != 1 {
					t.Fatalf("Antwort %q, Sessions %d", roh, r.zaehler())
				}
				_, _ = reg.Write(append(rohNachricht('Q', "SELECT 1\x00"), rohNachricht('X', "")...))
				rest, _ := io.ReadAll(reg)
				if !strings.Contains(string(rest), "SELECT 1") || strings.Contains(string(rest), "PGR-") {
					t.Fatalf("Antwort der regulären Verbindung %q", rest)
				}
				_, zweite := r.verbinde(false)
				_, _ = zweite.Write(append(rohStart("user", "app"), rohNachricht('X', "")...))
				_, _ = io.ReadAll(zweite)
				if n := r.zaehler(); n != 2 {
					t.Fatalf("Sessions %d, erwartet 2", n)
				}
			})
		}
	}
}

// cancelFrist ist die Zeit, die ein Test auf eine Antwort wartet, die kommen soll
// oder nicht kommen darf.
const cancelFrist = 2 * time.Second

// Abdeckung: LH-FA-05/Negative, LH-FA-05/Boundary — die Länge des CancelRequest ist aus dem
// Kopf entschieden (LH-FA-05.e *Länge des CancelRequest*, Reihenfolge der Gründe): Bei offener
// Verbindung, die nach dem Kopf von 8 Bytes nichts mehr sendet, lehnt der Recorder die Längen
// 8 bis 11, 13 bis 15, 269, 270 und 10000 sofort mit ErrorResponse (PGR-E6001) ab und schließt;
// bei den Längen 12, 16, 17 und 268 mit fehlendem Rest (0 und 3 Bytes) antwortet er bis zum
// Ablauf von zwei Sekunden nicht und fehlerfrei erst mit dem Ende der Verbindung (PGR-E6001);
// in record und replay, an jeder der sechs Stellen der Folge.
func TestCancelRequestKopf(t *testing.T) {
	z, _ := tlsEin(t, false, testpki.Blatt{})
	for _, replay := range []bool{false, true} {
		for _, st := range cancelStellen() {
			for _, l := range []int{8, 9, 10, 11, 13, 14, 15, 269, 270, 10000} {
				t.Run(fmt.Sprintf("replay=%v/%s/Länge %d nur Kopf, Verbindung offen", replay, st.name, l), func(t *testing.T) {
					t.Parallel()
					r := neuerRandServer(t, replay, tlsWenn(st, z))
					tcp, c := r.verbinde(st.klartext)
					_ = c.SetReadDeadline(time.Now().Add(cancelFrist))
					if _, err := c.Write(append(append([]byte{}, st.vor...), rohKopf(l, pgwire.CodeCancelRequest)...)); err != nil {
						t.Fatal(err)
					}
					roh, err := io.ReadAll(c)
					var ne net.Error
					if errors.As(err, &ne) && ne.Timeout() {
						t.Fatalf("keine Ablehnung aus dem Kopf, gelesen %q", roh)
					}
					_ = tcp.Close()
					b := r.befund(roh)
					if b.n != st.n || b.typen != "E" || !strings.Contains(b.fehler, model.CodeUnsupported) || b.erster != model.CodeUnsupported || b.sessions != 0 {
						t.Fatalf("Befund %+v", b)
					}
				})
			}
			for _, l := range []int{12, 16, 17, 268} {
				for _, g := range []int{0, 3} {
					t.Run(fmt.Sprintf("replay=%v/%s/Länge %d, %d Bytes Rest, Verbindung offen", replay, st.name, l, g), func(t *testing.T) {
						t.Parallel()
						r := neuerRandServer(t, replay, tlsWenn(st, z))
						tcp, c := r.verbinde(st.klartext)
						if _, err := c.Write(append(append([]byte{}, st.vor...), cancelNachricht(l, g, 0)...)); err != nil {
							t.Fatal(err)
						}
						_ = c.SetReadDeadline(time.Now().Add(cancelFrist))
						gelesen, err := io.ReadAll(c)
						var ne net.Error
						if !errors.As(err, &ne) || !ne.Timeout() || len(gelesen) != st.n {
							t.Fatalf("Antwort vor dem Ende der Verbindung: gelesen %q, Fehler %v", gelesen, err)
						}
						_ = tcp.Close()
						b := r.befund(nil)
						if b.erster != model.CodeUnsupported || b.sessions != 0 {
							t.Fatalf("Befund nach dem Ende der Verbindung %+v", b)
						}
					})
				}
			}
		}
	}
}

// tlsWenn liefert z, wenn die Stelle TLS verlangt, sonst nil.
func tlsWenn(st cancelStelle, z *pgwire.TLS) *pgwire.TLS {
	if st.tls {
		return z
	}
	return nil
}

// Abdeckung: LH-FA-05/Boundary — ein CancelRequest (Länge 12 und 16) auf einer Verbindung
// beendet eine laufende Interaktion einer zweiten Verbindung nicht (LH-FA-05.e *Länge des
// CancelRequest*, *Modus*), in record: Die Anfrage der zweiten Verbindung läuft im Recorder,
// der CancelRequest wird behandelt und schließt seine Verbindung, danach endet die Anfrage
// regulär mit ihrer Antwort und CommandComplete. Die Hälfte für replay steht in
// TestCancelRequestLaufendeInteraktionUnberuehrtReplay.
func TestCancelRequestLaufendeInteraktionUnberuehrt(t *testing.T) {
	for _, l := range []int{12, 16} {
		t.Run(fmt.Sprintf("Länge %d", l), func(t *testing.T) {
			r := neuerRandServer(t, false, nil)
			r.rec.queryHalt = make(chan struct{})
			_, reg := r.verbinde(false)
			_, _ = reg.Write(append(rohStart("user", "app"), rohNachricht('Q', "SELECT 1\x00")...))
			for ende := time.Now().Add(cancelFrist); ; time.Sleep(10 * time.Millisecond) {
				r.rec.mu.Lock()
				laeuft := slices.Contains(r.rec.ereignisse, "q:SELECT 1")
				r.rec.mu.Unlock()
				if laeuft {
					break
				}
				if time.Now().After(ende) {
					t.Fatal("die Anfrage läuft nicht")
				}
			}
			if roh := cancelSenden(r, cancelVoll(l), false, false); len(roh) != 0 {
				t.Fatalf("Antwort auf den CancelRequest %q", roh)
			}
			close(r.rec.queryHalt)
			_, _ = reg.Write(rohNachricht('X', ""))
			rest, _ := io.ReadAll(reg)
			if !strings.Contains(string(rest), "SELECT 1") || strings.Contains(string(rest), "PGR-") {
				t.Fatalf("Antwort der laufenden Anfrage %q", rest)
			}
		})
	}
}

// rohTypen liefert die Typbytes der Nachrichten eines Stroms hinter dem Verbindungsaufbau.
func rohTypen(roh []byte) string {
	var typen []byte
	for len(roh) >= 5 {
		l := int(binary.BigEndian.Uint32(roh[1:5]))
		if l < 4 || 1+l > len(roh) {
			typen = append(typen, '?')
			break
		}
		typen = append(typen, roh[0])
		roh = roh[1+l:]
	}
	return string(typen)
}

// Abdeckung: LH-FA-05/Boundary — ein CancelRequest (Länge 12 und 16) auf einer Verbindung
// beendet eine laufende Interaktion einer zweiten Verbindung nicht (LH-FA-05.e *Länge des
// CancelRequest*, *Modus*), in replay: Die zweite Verbindung hat eine Extended-Folge bis vor
// Sync begonnen (Parse, Bind, Execute), der CancelRequest wird behandelt und schließt seine
// Verbindung, ohne dass der Use Case eine Verbindung beendet, abbricht oder eine weitere
// Session öffnet; danach bekommt die zweite Verbindung mit Sync die vollständigen
// aufgezeichneten Antworten (ParseComplete, BindComplete, DataRow, CommandComplete,
// ReadyForQuery), und der Use Case sieht genau die vier Nachrichten in ihrer Reihenfolge.
func TestCancelRequestLaufendeInteraktionUnberuehrtReplay(t *testing.T) {
	for _, l := range []int{12, 16} {
		t.Run(fmt.Sprintf("Länge %d", l), func(t *testing.T) {
			r := neuerRandServer(t, true, nil)
			_, reg := r.verbinde(false)
			_, _ = reg.Write(append(rohStart("user", "app"), rohNachricht('P', "\x00SELECT 1\x00\x00\x00")...))
			_, _ = reg.Write(rohNachricht('B', "\x00\x00\x00\x00\x00\x00\x00\x00"))
			_, _ = reg.Write(rohNachricht('E', "\x00\x00\x00\x00\x00"))
			for ende := time.Now().Add(cancelFrist); len(r.rep.nachrichten()) < 3; time.Sleep(10 * time.Millisecond) {
				if time.Now().After(ende) {
					t.Fatal("die Extended-Folge erreicht den Use Case nicht")
				}
			}
			if roh := cancelSenden(r, cancelVoll(l), false, false); len(roh) != 0 {
				t.Fatalf("Antwort auf den CancelRequest %q", roh)
			}
			r.rep.mu.Lock()
			geoeffnet, geschlossen, erzwungen, shutdowns := r.rep.geoeffnet, r.rep.closed, r.rep.forced, r.rep.shutdowns
			r.rep.mu.Unlock()
			if geoeffnet != 1 || geschlossen != 0 || erzwungen != 0 || shutdowns != 0 {
				t.Fatalf("nach dem CancelRequest: geöffnet %d, geschlossen %d, erzwungen %d, Shutdown %d; erwartet 1, 0, 0, 0", geoeffnet, geschlossen, erzwungen, shutdowns)
			}
			_, _ = reg.Write(append(rohNachricht('S', ""), rohNachricht('X', "")...))
			rest, _ := io.ReadAll(reg)
			if typen := rohTypen(rest); !strings.HasSuffix(typen, "Z12DCZ") || strings.Contains(typen, "E") || strings.Contains(string(rest), "PGR-") {
				t.Fatalf("Antworten der zweiten Verbindung %q (Typen %q), erwartet bis ReadyForQuery, 1, 2, D, C, Z", rest, typen)
			}
			var folge []model.ClientMessageType
			for _, m := range r.rep.nachrichten() {
				folge = append(folge, m.Type)
			}
			if !slices.Equal(folge, []model.ClientMessageType{model.ClientParse, model.ClientBind, model.ClientExecute, model.ClientSync}) {
				t.Fatalf("Nachrichten beim Use Case %v", folge)
			}
		})
	}
}
