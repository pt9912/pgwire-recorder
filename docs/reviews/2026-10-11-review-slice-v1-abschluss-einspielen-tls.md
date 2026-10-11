# Review-Report: slice-v1-abschluss-einspielen-tls — 2026-10-11

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-20.a` (*Start*, *TLS*, *Abbruch im Aufbau*, Tabelle *Fehlerregeln*) und `LH-FA-17.a` (*Wirkung einer URL*, *Fehler*), gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und [ADR-0019](../plan/adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) und gegen `AGENTS.md` §3.

**Gegenstand:** `3d0bd0c..c84df26`: Code, Tests und Testhilfen `147afea` (+2288 −130), Belege des Implementers `24bc173` (Plan §7), Entscheidungen des Architect `c84df26` (Plan §6, `LH-FA-20.a` *TLS*, `spec/architecture.md` §6). Geändert sind `internal/adapters/driven/postgres/einspielen_tls.go` (neu) und `einspielen.go`, `internal/adapters/driving/cli/` (`cli.go`, `upstream.go`, `verbindung.go`, neu `zertifizierungsstelle.go`), `internal/bootstrap/bootstrap.go`, `internal/hexagon/model/fehler.go`, neu `internal/bootstrap/tlsproxy` und `internal/testpki`, die Tests dazu und `test/integration/play_tls_e2e_test.go`.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-tls.md` ganz gelesen am Stand `c84df26`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-20.a` (*Randformen je Schritt*, *Fehlerregeln*), `SPEC-038` (*Warten in Tests*); [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten); `spec/architecture.md` §4.4 und §6.
- ADRs wie oben, Abschnitt *Entscheidung*.
- `AGENTS.md` §3.3 bis §3.13, `.claude/commands/implement-slice.md` (*Randform-Rückgabe*).
- Adresse des Plans: `slice-v1-abschluss-einspielen-tls-doku` (gelesen auf die Sendung; der Diff ändert ihn nicht).
- Vorheriger Report am selben Modul: `slice-v1-abschluss-einspielen-laufsteuerung` (bis F-568); die Nummern dieses Laufs beginnen bei F-601.

**Ausgeführte Läufe** (alles Docker-only, Scratch-Präfix `rev-tls-`, nichts im Repo geändert):

- **Arbeitsweise der Mutanten.** Je Mutant eine frische Kopie eines `git archive`-Stands von `c84df26` per `cp -r` ohne `-p`, genau eine Ersetzung, `gofmt -l` auf der Datei (sauber; Ausnahme: bei M22 und M23 bleibt die Ausrichtung der Felder in `bootstrap.go` ungeordnet, die Läufe sind davon unberührt), `go test -count=1` des Pakets im gepinnten Image mit Modul-Cache, ohne Netz. 32 Läufe, drei davon ungültig (M17, M24 und M24b übersetzen nicht; M17 ist durch M17b ersetzt, M24 deckt §7 ab), 29 gültig.
  - **Rot (26):** M1 `InsecureSkipVerify` (15 Fälle, `TestEinspielTLSZertifikatFehler`); M2 `N` setzt ohne TLS fort (Downgrade; `TestEinspielTLSAntwort`, 3 Fälle); M3 Bytes nach `S` ungeprüft; M4 `N` als `PGR-E4002`; M5 Servername fest `127.0.0.1` (`TestEinspielTLSName`); M6 `--upstream-ca` nicht in den Speicher; M7 nur CA ohne Systemspeicher (`TestEinspielTLSZertifikatsspeicher`); M9 Senden des `SSLRequest` als `PGR-E4005`; M10 fremdes Byte als `PGR-E4005`; M11 Fehler der Aushandlung als `PGR-E4002`; M13 Ende vor der Antwort als `PGR-E4005`; M15 Typ der Blöcke ungeprüft; M16 kein Block erlaubt; M17b reguläre Datei ungeprüft (FIFO endet nicht binnen 10 s); M18 nur der erste Block; M19 Lesefehler des X.509 ignoriert; M20 `sslmode` entscheidet immer (55 Fälle); M21 `--upstream-ca` ohne TLS erlaubt; M22 und M23 Bootstrap reicht TLS beziehungsweise CA nicht durch (`TestRunPlayTLS`); M25 `Lstat` statt `Stat`; M26 Pfad in der Ursache der Meldung; M28 `PGR-E2004` statt `PGR-E2007`; M29 `Format` von `postgres.Zertifikate` gibt Subjekt aus; M31 Startup auf der Klartext-Verbindung nach der Aushandlung; M34 `n > 2` statt `n > 1`.
  - **Grün (3):** M8 `HandshakeContext(context.Background())` (laut §7 absichtlich: `AfterFunc` in `Verbinde` schließt ebenfalls; kein Mangel). M12 `MinVersion: tls.VersionTLS10` (F-604). M14 der Standardpfad ersetzt `x509.SystemCertPool` durch einen leeren Speicher (F-603).
- **Wiederholung gegen Flackern.** `go test -count=15` über Upstream-Adapter, CLI-Adapter und Bootstrap: grün; `go test -race -count=1` (golang:1.27 mit gcc) über dieselben drei Pakete: grün, keine Meldung.
- **Integration.** `tools/test/run-integration-tests.sh` in einer Kopie (Image-Name und Lauf-Präfix auf `rev-tls-` umbenannt): `run-integration-tests: gruen`; `TestE2EPlayTLS`, `…Anmeldung`, `…Abgelehnt`, `…ZertifikatFehler`, `TestE2EPlayKlartextAbgelehnt`, `TestE2EPlayUpstreamCAFehler`, `TestE2EPlayTLSAbbruch` bestanden.
- **Binary.** `go list -deps ./cmd/...` am Stand: weder `tlsproxy` noch `testpki` noch `testing`. Mit einem Blank-Import von `tlsproxy` in `bootstrap.go` in einer Kopie: `make a-check` meldet `gesamt: 0 Befund(e)`, `go list -deps` nennt `testing`, `testpki` und `tlsproxy` (F-602).
- **Namensabgleich.** Alle in §7 genannten `Test…`-Funktionen und Teilfälle per `grep` im Repo gesucht (72 Namen, Teilfälle mit Komma einzeln nachgeprüft): alle gefunden bis auf `TestE2EPlayZwischenstand`, den §7 als entfernt nennt.
- **Weitere.** `make kopf-check` in der Kopie: grün; `make abdeckung-check` am Stand: grün. Zertifikatsspeicher: ein kleines Programm mit `SSL_CERT_FILE=/dev/null` liest in einem frischen Prozess 0 Zertifikate, ohne die Variable 150 (F-603).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-601 | HIGH | Die Sicht nennt einen Aufruf des Go-Werkzeugs samt Paketmuster (`go list -deps ./cmd/...`) und spricht vom Linken des Binaries; das ist ein sprachspezifisches Konstrukt in `spec/architecture.md`. | `AGENTS.md` §3.4; Reviewer-Skill, HIGH *Spec-Stratum nennt Artefakt unterhalb oder außerhalb* | `spec/architecture.md` · `(`go list -deps ./cmd/...`) linkt sie nicht.` | ja — `grep -n 'go list' spec/architecture.md` | Sicht trägt Sprach-Werkzeugaufruf |
| F-602 | MEDIUM | „Nur Testcode importiert das Paket“ steht in beiden `doc.go`, in `spec/architecture.md` §6 und im Plan, und nichts prüft es: Ein Import von `tlsproxy` in `bootstrap.go` läuft durch `make a-check` (`internal/bootstrap/**` ist ausgenommen) mit 0 Befunden und linkt `testing`, `testpki` und den Proxy ins Binary. Der Auftrag des Architect, die Zusage zu verengen, ist in den `doc.go` nicht umgesetzt (`c84df26` berührt sie nicht). | `AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `internal/bootstrap/tlsproxy/doc.go` · `nur Testcode importiert\n// das Paket`; `internal/testpki/doc.go` · `Nur Testcode importiert das Paket.`; `spec/architecture.md` · `importiert nur Testcode` | ja — Kopie mit Blank-Import in `bootstrap.go`: `make a-check` grün, `go list -deps ./cmd/...` nennt `testing` | Zusage im Kommentar weiter als Prüfung |
| F-603 | MEDIUM | Der Standardpfad des Zertifikatsspeichers (`wurzeln = x509.SystemCertPool`) hat keinen Test: Der Mutant mit leerem Speicher als Standard bleibt in allen Paketen grün, und damit lehnte `play` jedes Serverzertifikat einer öffentlichen Zertifizierungsstelle ab, ohne dass ein Gate rot wird. Das akzeptierte Negativ begründet die Lücke mit „`SSL_CERT_FILE` wirkt erst bei der ersten Ladung im Prozess“; `starteBis` startet das Binary je Aufruf als frischen Prozess mit der Umgebung des Tests, dort liest die Variable ein neues Laden (gemessen: 0 gegen 150 Zertifikate). §7 behauptet „kein grüner Mutant“ und mutiert nur *CA ersetzt den Speicher* und *CA nicht aufgenommen*, nie den Standard. | `AGENTS.md` §3.10; `LH-FA-20.a` *TLS*; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | `internal/adapters/driven/postgres/einspielen_tls.go` · `wurzeln = x509.SystemCertPool`; Plan §6 · `Zertifikatsspeicher des Systems im Integrationstest` | ja — Mutant M14 (`wurzeln = func() (*x509.CertPool, error) { return x509.NewCertPool(), nil }`): `go test` der drei Pakete und die Integration grün | Standardpfad einer Zusage ohne Mutation |
| F-604 | LOW | Der Fall `S, dann TLS 1.0` scheitert, weil die Gegenstelle `MaxVersion = TLS10` setzt und ihr eigener Standard (ab TLS 1.2) die Aushandlung ablehnt; die Untergrenze des Clients prüft er nicht. Der Mutant `MinVersion: tls.VersionTLS10` im Client bleibt grün, obwohl §7 die Zeile unter *Fehler der Aushandlung … Version* führt. | `AGENTS.md` §3.11; Plan §6 *TLS-Version und Verfahren* (keine Zusage) | `internal/adapters/driven/postgres/einspielen_tls_test.go` · `{"S, dann TLS 1.0", tlsStelle{antwort: []byte("S"), config: alt}, model.CodeLogin}` | ja — Mutant M12 | Testname sagt mehr als der Test prüft |
| F-605 | MEDIUM | Der Code-Commit `147afea` entscheidet vier Randformen, die §6 nicht nannte, und gibt sie erst danach zurück: das Senden des `SSLRequest` als `PGR-E4002` (Zeile und Test `TestEinspielTLSSendenUndLesen` im Code), den Ort der Testhilfen (`internal/bootstrap/tlsproxy`, `internal/testpki`), den Typ der Zertifikate, die Hook-Signatur. Der Architect hat sie am selben Stand bestätigt; die Entscheidung stand damit im Code, bevor sie in §6 stand. | `AGENTS.md` §3.12; `.claude/commands/implement-slice.md` *Randform-Rückgabe*; `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` | Plan §6 · `Rückgabe des Implementers (geprüft vom Architect am 2026-10-11, nach dem Code)`; `einspielen_tls.go` · `SSLRequest nicht zu senden` | nein — kein Gate fängt die Reihenfolge von Entscheidung und Code | Randform im Code entschieden, dann zurückgegeben |
| F-606 | INFO | `close(p.Hello)` im hängenden Proxy ist nicht gegen eine zweite Verbindung geschützt; ein zweites Byte nach `S` auf einer zweiten Verbindung desselben Proxys ließe den Testprozess mit „close of closed channel“ enden. Die heutigen Tests öffnen je hängendem Proxy eine Verbindung. | Maintainability | `internal/bootstrap/tlsproxy/proxy.go` · `close(p.Hello)` | nein | Schließen eines Kanals ohne Schutz in Testhilfe |
| F-607 | INFO | `LH-FA-20.a` *TLS* hängt „(Grenze: ohne Test, …)“ an den Satz über Bytes nach `S`, die mit dem `S` eintreffen; dieser Fall hat Tests (`TestEinspielTLSAntwort/S_mit_einem_Byte_dahinter`, `/S_mit_TLS-Record_dahinter`), ungeprüft ist nur der Fall in zwei Segmenten. Die Aussage ist enger zu lesen als geschrieben. | `AGENTS.md` §3.11 | `spec/spezifikation.md` · `(Grenze: ohne Test, weil die Segmentierung des Netzes nicht steuerbar ist)` | nein | Grenze am falschen Satzteil |
| F-608 | INFO | Der Diff umfasst 2443 Zeilen (rund 14 % über den ~2000 Zeilen, die der Plan als Maß für eine Sitzung nennt), davon etwa 350 Produktcode; der Implementer hat nicht angehalten und das in §7 benannt. In einer Sitzung prüfbar war er hier; Zuständig für das Maß ist der Planner. | Plan §8; `v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice | Plan §7 · `Das liegt rund 14 % über den ~2000 Zeilen` | nein | Größe über der Sitzungsgrenze |

## Negativbefunde

Zu den Stufen der Aufgabe, je betrachteter Bereich:

| Bereich | Ergebnis |
|---|---|
| (1) `internal/adapters/driven/postgres/einspielen_tls.go`: `SSLRequest` über `pgproto3` (acht Bytes im Test), `S`/`N`/anderes Byte/Ende, Bytes nach `S`, Fehlerabbildung `PGR-E4005`/`PGR-E4002`, kein `Terminate`, Abbruch über `HandshakeContext(ctx)` und `AfterFunc` | geprüft, ohne Befund (Mutanten M2 bis M4, M9 bis M11, M13, M34 rot; M8 absichtlich grün) |
| (1) Konfiguration: kein `InsecureSkipVerify`, `ServerName` aus dem Host der Adresse (IPv4, IPv6 mit und ohne Zone), `RootCAs` aus Systemspeicher und `--upstream-ca`, nicht ladbarer Speicher gilt als leer; keine Nonce und kein Geheimnis in einer Meldung; `Format` auf `Zertifikate` | geprüft, ohne Befund außer F-603 (Standardpfad) und F-604 (Untergrenze, keine Zusage) |
| (2) `internal/adapters/driving/cli/zertifizierungsstelle.go`: reguläre Datei (Link gefolgt, FIFO ohne Warten), mindestens ein Block, jeder `CERTIFICATE`, X.509 lesbar, Text außerhalb unbeachtet, `PGR-E2007` ohne Pfad und Inhalt | geprüft, ohne Befund (M15 bis M19, M25, M26, M28 rot) |
| (2) `upstream.go`, `cli.go`, `verbindung.go`: `--upstream-tls` gesetzt (Option, Umgebung, Schlüssel, auch `false`) vor `sslmode`, `--upstream-ca` ohne TLS `PGR-E2001` nach `--upstream` und vor den Variablen, die Datei als letzte Prüfung vor dem Laden, leerer Wert, `config show` liest die Datei nicht; die Teilung von `lies` ändert die Reihenfolge der Quellen nicht | geprüft, ohne Befund (M20, M21 rot) |
| (3) `internal/bootstrap/bootstrap.go`, `internal/hexagon/model/fehler.go` | geprüft, ohne Befund: Verdrahtung ohne Logik (zwei Felder, eine Typumwandlung, eine Konstante); M22, M23 rot |
| (4) `internal/testpki`, `internal/bootstrap/tlsproxy`: Zertifikate zur Testzeit mit `crypto/x509`, nichts eingecheckt; Fristen aller neuen Warte-Stellen als Literal ≤ 60 s (`SPEC-038`); kein Flackern in 15 Läufen, `-race` ohne Meldung; Proxy leitet nur nach ausgehandeltem TLS weiter und lehnt Verbindungen ohne `SSLRequest` mit 28000 ab; Binary ohne beide Pakete am Stand | geprüft, ohne Befund außer F-602 (Zusage ohne Prüfung) und F-606 |
| (6) §3.12 gegen §6: jede Operation des Diffs hat eine Randform in §6 oder die Rückgabe in `c84df26`; Downgrade-Schutz (`N` bei `sslmode=require` ist `PGR-E4005`, kein Startup im Klartext; M2 rot, `TestE2EPlayTLSAbgelehnt` grün) | geprüft; Befund F-605 zur Reihenfolge |
| (6) §3.11: Hilfetext von `play` (`--upstream-tls`, `--upstream-ca`) gegen das Verhalten; Kommentare in `einspielen_tls.go`, `zertifizierungsstelle.go`, `cli.go` | geprüft, ohne Befund außer F-602, F-604, F-607 |
| (7) §3.9: Plan §1, §3, §6 folgen dem Diff (Testhilfen, Typ der Zertifikate, Hook-Signatur, Ort der CA-Datei); `make kopf-check` grün | geprüft, ohne Befund |
| (7) §3.13: Sendung an `slice-v1-abschluss-einspielen-tls-doku` besteht seit vor dem Diff, der Diff ändert den Nehmer nicht; Zählung im Plan §8: Liefer-Punkte 2, Schichten 2 (Upstream-Adapter, CLI-Adapter), Bootstrap und Konstante ohne Logik, Testhilfen ohne Schicht | geprüft, ohne Befund |
| `spec/spezifikation.md` `LH-FA-20.a` (Änderung in `c84df26`) | geprüft, ohne Befund außer F-607; keine ADR-, Slice- oder Wellen-Kennung im Stratum |
| [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) und [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) (Core-Reinheit, `pgproto3` nur in den beiden PGWire-Adaptern): `internal/hexagon/**` importiert weder `crypto/x509` noch `crypto/tls` | geprüft, ohne Befund |
| `//nolint`, Lockerung von `.golangci.yml`, `.a-check.yml` | geprüft, ohne Befund (keine Änderung am Profil oder an den Regeln) |
| `docs/user/abdeckung-*.md` (erzeugt) | geprüft, ohne Befund; `make abdeckung-check` grün |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Sicht trägt Sprach-Werkzeugaufruf · Zusage im Kommentar weiter als Prüfung · Standardpfad einer Zusage ohne Mutation · Testname sagt mehr als der Test prüft · Randform im Code entschieden, dann zurückgegeben · Schließen eines Kanals ohne Schutz in Testhilfe · Grenze am falschen Satzteil · Größe über der Sitzungsgrenze

## Verdikt

**Merge-blockierend:** ja — F-601 (HIGH), F-602 und F-603 (MEDIUM). F-605 ist nach der Bestätigung des Architect ohne Änderung am Code erledigt und zählt für den Steering-Loop; er blockiert allein nicht. Die Prüfung des Zertifikats selbst (Kette, Name, Ablauf, unbekannte Stelle, kein Überspringen), der Downgrade-Schutz und die CA-Datei sind durch Mutanten gefangen; die Blockierung betrifft die Sicht, die Zusage über die Testhilfen und den Standardpfad des Systemspeichers.

**Übergabe:**

- **F-601** an den **Architect** (Stratum `spec/architecture.md`, Artefakt dieser Übergabe: dieser Report, Zeile F-601). Kein Widerspruch des Implementers liegt vor; der Konflikt-Pfad ist nicht eröffnet.
- **F-602** an den **Implementer**, mit dem Auftrag des Architect (verengen) und dem Befund, dass `make a-check` den Import nicht fängt; ob die Zusage enger gefasst oder mit einem Test versehen wird (`AGENTS.md` §3.11), entscheidet der Implementer, ein Gate-Eingriff (`AGENTS.md` §3.6) gehört dem Architect.
- **F-603** an den **Implementer**: ein Test des Standardpfads und die Mutation dazu in §7; das akzeptierte Negativ in §6 ist neu zu begründen oder zu streichen (der Architect hat es am 2026-10-11 bestätigt, daher Rückfrage an ihn).
- **F-604** bis **F-608** an den **Implementer** beziehungsweise den Planner (F-608); ohne Blockade.
- Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier.
