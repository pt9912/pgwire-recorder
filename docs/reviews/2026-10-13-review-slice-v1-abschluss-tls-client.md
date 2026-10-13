# Review slice-v1-abschluss-tls-client

**Rolle:** Reviewer. **Datum:** 2026-10-13.
**Gegenstand:** Diff `4223e5e..464ca8a` (Move-Commit bis HEAD, sieben Commits) gegen den Plan `docs/plan/planning/done/slice-v1-abschluss-tls-client.md`, [ADR-0018](../plan/adr/0018-tls-zum-client.md), [LH-FA-23](../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client), `LH-FA-23.a`, `LH-FA-17.a` und die Randformen der Spezifikation, `AGENTS.md` §3.2, §3.6, §3.7, §3.9 bis §3.14.
**Summary:** 0 HIGH · 3 MEDIUM · 2 LOW · 6 INFO.

## Findings

### F-1 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10 (je Zusage eine Mutation), §3.12 (Randform im Code entschieden)
- `pfad`: `internal/adapters/driving/pgwire/server.go` · "if laenge == 8 && (code == codeSSLRequest || code == codeGSSEncRequest) {"
- `befund`: Der Diff liest `SSLRequest` und `GSSENCRequest` selbst und entscheidet dabei, dass eine Anfrage mit diesem Code und einer Länge ungleich 8 nicht hier, sondern über `ReceiveStartupMessage` läuft. §6 des Plans nennt diese Randform nicht, die Spezifikation auch nicht, und kein Test hat sie. Das Entfernen von `laenge == 8 &&` lässt das ganze Paket `pgwire` grün (nachgefahren). Szenario: Eine Anfrage mit Code `SSLRequest` und Länge 12 wird dann nach `Discard(8)` als Anfrage behandelt, vier Restbytes bleiben im Leser; das Verhalten wechselt von der Meldung der Bibliothek (`PGR-E6001`) zu `N` beziehungsweise `PGR-W3002`, ohne dass ein Test es merkt.
- `verifizierbar`: ja (Mutation, `go test` im Paket `pgwire`)
- `klasse`: Randform im Code entschieden, nicht in §6, ohne Mutation

### F-2 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10 (Aufzählung zählt je Element), §3.9
- `pfad`: `internal/adapters/driving/cli/zertifizierungsstelle.go` · "if !info.Mode().IsRegular() {" (neu in `liesDatei`)
- `befund`: `LH-FA-23.a` *Dateien* zählt "Verzeichnis, FIFO, Gerät oder Socket" auf; §6 des Plans nennt nur "FIFO, Verzeichnis", die Tests (`TestTLSDateiKeineRegulaere`) haben Verzeichnis, Link auf Verzeichnis und FIFO. Ersetzt man die Bedingung durch `info.IsDir() || info.Mode()&os.ModeNamedPipe != 0`, bleibt `TestTLS…` im CLI-Adapter grün (nachgefahren). Szenario: Eine Gerätedatei (zum Beispiel `/dev/zero`) als `--tls-key` wird gelesen, der Start liest ohne Ende, statt `PGR-E2007` zu melden. Der Plan ist enger als die Spezifikation, die der Slice liefert.
- `verifizierbar`: ja (Mutation, `go test` im CLI-Adapter)
- `klasse`: Element einer Aufzählung ohne Fall und Mutation

### F-3 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10, Plan §6 *Geheimnis* · `LH-RB-01`
- `pfad`: `docs/plan/planning/done/slice-v1-abschluss-tls-client.md` · "der Schlüssel erscheint in keiner Log-Zeile (auch `debug`)"
- `befund`: Die Zusage nennt jede Log-Zeile, auch `debug`. Kein Test fährt `record` oder `replay` mit `--log-level=debug` und TLS und sucht den Schlüssel in der Ausgabe; die Logger der PGWire-Tests laufen auf Stufe `INFO` (`slog.NewTextHandler(log, nil)`), und die Warnung prüft nur auf `BEGIN`. Die Belege in §7 führen für diese Zusage keine Zeile (nur "keine Textform"). Szenario: Eine spätere Debug-Zeile im Start oder im Handshake, die Optionen oder Konfiguration ausgibt, bricht die Zusage, und nichts wird rot. Heute leckt der Code nicht (Text-Handler, `Format` des Typs `Schluessel`).
- `verifizierbar`: ja (Test mit Stufe `debug` über `bootstrap.Run`)
- `klasse`: Sicherheitszusage ohne Test auf der strengsten Stufe

### F-4 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.11 (Hilfetext und Spezifikationsaussage nur mit Test), §3.10
- `pfad`: `internal/adapters/driving/cli/cli.go` · "PEM-Datei mit dem Zertifikat des Servers und seiner Kette"
- `befund`: Hilfetext und `LH-FA-23.a` ("Die Zertifikatskette steht in der Zertifikatsdatei") setzen voraus, dass die Kette an `NeuesTLS` geht. Kein Test prüft, dass der Server die Kette ausliefert (`tlsEin` baut mit einem einzelnen Blatt; kein `PeerCertificates` im Repo); das Abschneiden auf das erste Zertifikat (`Certificate: paar.Certificate[:1]`) lässt `pgwire` grün (nachgefahren). Ebenso hat die Kette in `TestTLSDateienGueltig` nur den Fall "abgelaufen"; ein "noch nicht gültiges" zweites Zertifikat fehlt (Mutation `NotBefore` des letzten Zertifikats statt des ersten bleibt grün; `NotAfter` wird gefangen).
- `verifizierbar`: ja
- `klasse`: Hilfetext-Zusage ohne Test

### F-5 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.11 (Kommentar sagt nur zu, was ein Test prüft)
- `pfad`: `internal/adapters/driving/pgwire/tls.go` · "Bytes, die der Client nach dem SSLRequest vor der Antwort sendet, beenden die Verbindung mit PGR-W3002, bevor ein Handshake beginnt."
- `befund`: Der Kommentar von `sslAnfrage` sagt zu, dass alle Bytes vor der Antwort so enden. Der Code erkennt nur, was beim Senden des `S` schon im Puffer liegt (`k.br.Buffered()`); die Spezifikation hat diesen Satz nach der Rückgabe des Implementers enger gefasst ("beim Senden des `S` schon empfangen", später eintreffende Bytes scheitern im Handshake). Der Kommentar zog nicht nach.
- `verifizierbar`: nein
- `klasse`: Kommentar weiter als Spezifikation und Test

### F-6 INFO

- `quelle`: `AGENTS.md` §3.9, §3.14 · `pfad`: `docs/plan/planning/done/slice-v1-abschluss-tls-client.md` · "vor den Commits 2094 Zeilen"
- `befund`: §7 nennt die Größe 2094. Die Metrik des Plans (`git diff --numstat 4223e5e..HEAD -- internal test docs/user`) misst am letzten Commit 2186 (2094 Einfügungen, 92 Löschungen), mit den Pfaden des Plans ohne `spec` und Plan. §7 nennt keinen Stand nach der Runde der Rückgabe. Keine Aktion für die Marke (2300), aber die Zahl im Plan folgt dem Stand nicht.
- `klasse`: Zahl im Plan nach Messung nicht nachgezogen

### F-7 INFO

- `quelle`: `AGENTS.md` §3.14, Plan §4 · `pfad`: Plan §4 · "über 2300"
- `befund`: Größe 2186 gegen Marke 2300: nicht eingetreten, Abstand 114. Die Schätzung des Architect (1400 bis 1800) ist um 386 (21 %) überschritten; Tests und Gegenbelege tragen rund zwei Drittel (`server_tls_test.go` 657, `tls_test.go` 482, `tls_client_e2e_test.go` 156, `bootstrap/tls_test.go` 108). Die Meldeschwelle 1900 steht in keinem Artefakt dieses Repos (Plan, `AGENTS.md`, `.claude/commands`, `harness/`; `grep` ohne Treffer), ich kann ihre Verletzung nicht gegen eine Regel prüfen; wenn sie Weisung des Planners oder Architect war, ist das Weiterlaufen über sie ohne Anhalten eine Abweichung im Vorgehen, nicht im Diff. Der Slice ist in einer Sitzung noch prüfbar, aber an der Grenze; weiteres Wachstum durch diese Review-Runde hätte die Rückführung ausgelöst.
- `klasse`: Schätzung in der Metrik unterschritten, Meldeschwelle nicht verankert

### F-8 INFO

- `quelle`: Maintainability · `pfad`: `internal/adapters/driving/cli/tls.go` · "func (z TLSZugang) Gesetzt() bool"; `internal/adapters/driving/cli/fremd_test.go` · "return map[string]bool{}"
- `befund`: `Gesetzt()` hat im Produktionscode keinen Aufrufer (nur Tests); `NeuesTLS` entscheidet "keine Verschlüsselung" selbst. `nichtGelesen()` liefert nun eine leere Menge, `TestLeserUmgebungNichtGelesenerOptionen` läuft damit ohne Fall, wie es der Plan für diese Bindung vorsieht. Keine Aktion.
- `klasse`: nur von Tests benutzte API / leere Hilfsmenge

### F-9 INFO

- `quelle`: `AGENTS.md` §3.12 · `pfad`: Commit-Folge `a65a095` (Code, `tlsVerbindung`) vor `246dce2` (Entscheidung in §6 und Spezifikation)
- `befund`: Der Code enthielt "kein `close_notify`" und die Behandlung von Text außerhalb der Blöcke, bevor der Architect sie in §6 und der Spezifikation entschied; §3 des Plans nennt `tlsVerbindung` schon in `179a76f`. Aus dem Diff ist nicht ablesbar, ob der Implementer vor dem Code angehalten hat; die Rückgabe steht in §7, die Entscheidung deckt den Code, und `246dce2` trägt weder Code noch Tests. Kein Verstoß gegen die Regel "Randform im Code-Commit", aber die Reihenfolge ist Code, dann Entscheidung.
- `klasse`: Entscheidung des Architect nach dem Code

### F-10 INFO

- `quelle`: `AGENTS.md` §3.10 · `pfad`: `internal/adapters/driving/pgwire/tls.go` · "if len(zertifikat) == 0 && len(schluessel) == 0 {"; `internal/adapters/driving/cli/tls.go` · `block.Headers["Proc-Type"] == "4,ENCRYPTED"`
- `befund`: `NeuesTLS` mit nur einem der beiden Werte ist weder in der Spezifikation noch im Plan entschieden (Fehler oder `nil`); die Mutation `&&` zu `||` bleibt grün. Der CLI-Adapter schließt den Fall aus. Ebenso hat die Prüfung der Kopfzeile nur den Wert `4,ENCRYPTED` im Test; ein anderer Wert von `Proc-Type` ist nicht zugesagt. Keine Aktion für diesen Diff.
- `klasse`: nicht zugesagte Eingabe ohne Entscheidung

### F-11 INFO

- `quelle`: `LH-RB-01` · `pfad`: `internal/adapters/driving/cli/tls.go` · "Jede Formatierung gibt festen Text aus"
- `befund`: `Schluessel` schützt über `fmt.Formatter` (getestet: `%v %+v %#v %s %q %x %d`, auch Zeiger und Teilstrukturen). Ein `encoding/json`- oder `yaml`-Aufruf auf einem `Command` gäbe den Inhalt aus, weil der String-Typ keinen `MarshalText` trägt; das Repo ruft beides nicht auf Optionen, und der Log-Handler ist ein Text-Handler. Zugesagt ist nur die fmt-Form ("keine Textform des `Command`, Fall `%v` und `%+v`").
- `klasse`: Schutz nur für fmt

## Eigene Prüfung (nicht auf Zusage)

- Mutationen nachgefahren (Kopie des Arbeitsbaums im Image `pgr-deps`, `go test -count=1`, Arbeitsbaum unberührt, `git status` sauber). Rot, wie behauptet: `Buffered() > 0` zu `false` (`TestTLSHandshakeScheitert/Bytes_vor_dem_S`); `Close` mit `close_notify` (`TestTLSZwangsende`, `TestTLSEndeOhneAlert`); `GSSENCRequest` aus der Sonderbehandlung (`TestTLSUnverschluesselt`); Prüfung des `GSSENCRequest` auf der verschlüsselten Verbindung entfernt (`TestTLSZweiterAufbau`); `MinVersion` TLS 1.1 (`TestTLSMindestversion`); Ablauf gegen das letzte statt das erste Zertifikat (`TestTLSDateienGueltig/*/Kette_mit_abgelaufenem_zweiten_Zertifikat`). Grün (Befunde): `laenge == 8` entfernt (F-1), Gerät/Socket (F-2), Kette abgeschnitten und Beginn des letzten Zertifikats (F-4), `&&` zu `||` in `NeuesTLS` und `Proc-Type` (F-10).
- Architekturregel: `.a-check.yml` und `.golangci.yml` im Diff unverändert (`git diff --name-only`); der CLI-Adapter importiert `crypto`, `crypto/ecdsa`, `crypto/ed25519`, `crypto/rsa`, `crypto/x509` und `encoding/pem`, nicht `crypto/tls`; `crypto/tls` steht nur in `pgwire/tls.go` (Test: `make a-check`, `make a-check-negativ` in der schnellen Kette).
- Schlüssel und Pfade: Meldungen von `liesDatei`, `zertifikate`, `privaterSchluessel`, `tlsSchluessel` nennen Option und feste Ursache, nie Pfad oder Inhalt (`GEHEIM`, `BEGIN`, Pfad im Test geprüft, auch FIFO und Verzeichnis); `NeuesTLS` wirft den Fehler von `tls.X509KeyPair` nicht durch; die Warnung `PGR-W3002` trägt `code`, `remote`, `error` des Handshakes (Go-Fehlertexte nennen keinen Schlüssel); das Debug-Log beim Herunterfahren trägt nur den Fehlertext. Lücke: F-3.
- Reihenfolge am Start: Kombination (`tlsKombination`) vor der `zuletzt`-Schleife, Zertifikat (`tlsZertifikat`) vor Schlüssel (`tlsSchluessel`) als letzte Optionen der Tabelle, `pgwire.NeuesTLS` im Bootstrap vor `recording.New` und vor dem Listen-Port; `TestTLSReihenfolge` und `TestE2ETLSStartfehler` decken Kombination vor `--upstream`, Zertifikat vor Schlüssel und `PGR-E2007` vor der Ausgabedatei. Die Kopplung "tlsSchluessel liest `zugang(c).Zertifikat`" ist im Kommentar genannt und durch die Vertauschung gefangen.
- `nichtGelesen()` / `zuordnung()`: `nichtGelesen()` ist leer, die drei Optionen stehen in der Tabelle für `record` und `replay`; `TestLeserUmgebungAndererKommandos` hält `play` bei `PGR-E2001`, `TestTLSQuellen/play` bei `PGR-E2004`; `TestLeserOptionen` hält die Reihenfolge der Tabelle.
- Handshake: `MinVersion` TLS 1.2 (rot bei 1.1), `ClientAuth` nicht gesetzt (Test: `GetClientCertificate` nie aufgerufen), `NextProtos` nicht gesetzt (Test: Client bietet `h2`, ausgehandelt `""`); Bytes vor dem `S` beenden ohne `S` und ohne Handshake (`Buffered`); Klartext nach dem `S` und Abbruch nach dem `S` sind `PGR-W3002`; Herunterfahren und Zwangsende melden nicht. Rangfolge: `PGR-E6003` nur für eine `StartupMessage` der Version 3.0 (Test mit vier Nachbarfällen); `PGR-E6003` kommt aus `s.fail` mit Exit-Code 6.
- `AGENTS.md` §3.2: keine `//nolint`-Direktive im Diff (`grep` auf `internal`, `cmd`, `test` leer). §3.6: keine Schwelle, Regel oder Ausnahme in `.a-check.yml`, `.golangci.yml`, `Makefile` gesenkt. §3.7: neue Kommentare Indikativ, ohne Alternativen; Ausnahme F-5. §3.11: Handbuch und README unberührt (Doku-Folge-Slice `slice-v1-abschluss-tls-client-doku`); Hilfetext ohne Chronik, Aussagen in `TestTLSHilfe` und den Tests der Optionen belegt, Ausnahme F-4.
- §3.9: Kopf, §1, §3, §6, §8 folgen dem Diff (Nachzählungen in §8: drei Liefer-Punkte, zwei Schichten plus Verdrahtung; der Bootstrap hat einen `if err != nil` je Konstruktor, reine Weitergabe wie die Nachbarn); Ausnahme F-6. §3.13: die Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` ist getragen (Leser angebunden, `nichtGelesen()` gekürzt). §3.3: der Move `4223e5e` ist rein; die sechs Folgecommits enthalten kein `R`/`D`. §3.5: keine ADR geändert.
- Randform "Rückgabe des Implementers": fünf Aussagen der Spezifikation (kein `close_notify`; Text außerhalb der Blöcke; Ursache bei mehreren Mängeln nicht zugesagt; Bytes nach dem `SSLRequest`; PKCS #8 anderer Art) haben je einen Test (`TestTLSEndeOhneAlert`, `TestTLSTextAusserhalbDerBloecke`, X25519-Fall, `TestTLSHandshakeScheitert`); der Test der Ursache prüft sie nur bei Dateien mit einem Mangel, wie zugesagt.
- Unbenutzter Code: außer `Gesetzt()` (F-8) keiner; `tlsVerbindung` hat einen Aufrufer, `kanal.verschluesselt` wird gelesen und gesetzt, `ParseMitUhr` ist Testhilfe.
- Nicht ausgeführt von mir: `make test-integration` (braucht den gepinnten PostgreSQL); die Integrationstests habe ich gelesen, nicht gefahren.

## Geprüft, ohne Befund

- `internal/adapters/driving/cli/` (`tls.go`, `cli.go`, `zertifizierungsstelle.go`, Tests): F-2, F-4, F-8, F-10, F-11; sonst ohne Befund.
- `internal/adapters/driving/pgwire/` (`tls.go`, `server.go`, Tests): F-1, F-5; sonst ohne Befund.
- `internal/bootstrap/`: ohne Befund (Verdrahtung, Reihenfolge vor Ausgabedatei und Port).
- `internal/hexagon/model/fehler.go`: ohne Befund (zwei Konstanten).
- `internal/testpki/`: ohne Befund (Zertifikate im Testlauf erzeugt, keine eingecheckte Datei mit Schlüssel).
- `test/integration/`: ohne Befund (gelesen, nicht gefahren).
- `spec/spezifikation.md`: ohne Befund (kein Artefakt unterhalb oder außerhalb genannt).
- `docs/user/abdeckung-*.md`: ohne Befund (erzeugt, `make abdeckung-check` in der Kette).
- `.a-check.yml`, `.golangci.yml`, `Makefile`: unverändert, ohne Befund.
- `docs/plan/planning/done/slice-v1-abschluss-tls-client.md`: F-3 (Quelle), F-6, F-7.
