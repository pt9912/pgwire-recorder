# Review slice-v1-abschluss-anmeldung

**Rolle:** Reviewer. **Datum:** 2026-10-14.
**Gegenstand:** Diff `fa5a64c..bf62b3d` (Move-Commit bis HEAD, 16 Commits) gegen den Plan `docs/plan/planning/in-progress/slice-v1-abschluss-anmeldung.md`, [ADR-0041](../plan/adr/0041-anmeldung-im-record-als-rueckruf-des-upstream-ports.md), [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md), [LH-FA-05](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-05.b`](../../spec/spezifikation.md#lh-fa-05b--authentifizierung) mit dem Block *Randformen der Vermittlung (Record)*, `spec/architecture.md` §4.4 und `AGENTS.md` §3.2, §3.6, §3.7, §3.9 bis §3.14.
**Summary:** 0 HIGH · 2 MEDIUM · 1 LOW · 4 INFO.

Vorgehen: Diff und Plan §7 gelesen; die Tests des Pakets in einer Kopie des HEAD-Stands (`git archive`, Image `deps`) gefahren (grün); 17 Mutanten am Quelltext nachgefahren, die Kopie danach verworfen, der Arbeitsbaum blieb unberührt; `make a-check` im Arbeitsbaum: `gesamt: 0 Befund(e)`. Mutanten dieses Laufs (Kennung M, nicht die des Implementers): 13 rot, 4 grün plus 2 ohne Aussage (M6 Build-Fehler, M18 ohne Wirkung gewählt); die grünen sind F-1 bis F-3.

## Findings

### F-1 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10 (je Zusage eine Mutation), [`LH-FA-05.b`](../../spec/spezifikation.md#lh-fa-05b--authentifizierung) *Ende des Servers*
- `pfad`: `internal/adapters/driven/postgres/vermittlung.go` · "if istVerbindungsende(err) {" (in `lesefehler`)
- `befund`: Die Spezifikation sagt zu: Endet die Verbindung des Servers in der Anmeldung, ist das `PGR-E4003`. Der Slice liefert dafür nur den Fall "Senden der Antwort scheitert" (`TestVermittlungSchreibfehlerAmServer`, Mutation Y4); die Lesehälfte hat keinen Fall. Ersetzt man die Bedingung durch `if false {`, bleibt das ganze Paket `postgres` grün (M5, nachgefahren); das Ende des Servers wird dann als `PGR-E6001` ("nicht lesbar") gemeldet. Szenario: Ein Server schließt die Verbindung nach der Anforderung oder mitten im SCRAM-Austausch; der Client bekommt eine falsche Fehlerklasse (Exit-Code 6 statt 4), und nichts wird rot. Die Tabellen in §7 führen für diese Hälfte keine Zeile.
- `verifizierbar`: ja (Mutation, `go test ./internal/adapters/driven/postgres`)
- `klasse`: Hälfte einer Zusage ohne Fall und Mutation

### F-2 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10 (Aufzählung zählt je Element), [`LH-FA-05.b`](../../spec/spezifikation.md#lh-fa-05b--authentifizierung) *Reihenfolge des Servers*, Plan §7 ("jede Zelle der Tabelle")
- `pfad`: `internal/adapters/driven/postgres/upstream_anmeldung_test.go` · "{rSASL, rSASLWeiter, rSASLAbschluss, rSASLAbschluss}, {rSASL, rSASLWeiter, rSASLAbschluss, rSASLWeiter}" (Liste in `TestVermittlungReihenfolge`)
- `befund`: Die Spezifikation zählt auf: eine weitere Anforderung (Code 3, 5, 10), nachdem eine beantwortet ist oder während ein SASL-Austausch läuft. Die Liste hat nicht jedes Element in jedem Stand: Code 10 und Code 5 im laufenden Austausch nach Code 11 (`[10 11 10]`, `[10 11 5]`) und nach Code 12 (`[10 11 12 10]`, `[10 11 12 5]`) fehlen. Vier Mutanten bleiben grün (nachgefahren): Code 10 zusätzlich in `standSASLLaeuft` (M7) und in `standSASLAbschluss` (M9), Code 5 zusätzlich in `standSASLLaeuft` (M8) und in `standSASLAbschluss` (M10). Szenario: Ein Server (oder ein Fehler im Code) schickt nach dem ersten SASL-Schritt eine neue SASL-Anforderung; der Recorder reicht sie an den Client weiter, statt `PGR-E4002` zu melden. Das Fangen der Zeilen `[10 11 3]` und `[10 11 12 3]` deckt nur die Mutation des Ausdrucks für Code 3. Die Aussage des Plans, jede Zelle der Tabelle sei durch R1 bis R22 gebrochen, gilt nicht für diese Zellen, weil die Ausdrücke je Code getrennt stehen.
- `verifizierbar`: ja (Mutation, `go test ./internal/adapters/driven/postgres`)
- `klasse`: Element einer Aufzählung ohne Fall und Mutation

### F-3 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.11 (Kommentar sagt nur zu, was ein Test prüft), §3.10
- `pfad`: `internal/adapters/driving/pgwire/server.go` · "die Lesefrist des Wächters gilt für sie nicht" (Kommentar von `startupBisSignal`) und "_ = roh.SetReadDeadline(time.Time{})"
- `befund`: Der Kommentar sagt zu, dass die Lesefrist des Wächters für die Anmeldung nicht gilt. `TestAnmeldungWartetOhneFrist` fängt nur eine Frist, die der Code selbst nach der Startnachricht setzt (X5); der Zweig, der eine vom Wächter gesetzte Frist zurücknimmt (`if ok && gewecktWorden`), hat keinen Test: Entfernt man die Zeile, bleibt das Paket `pgwire` grün (M14, nachgefahren). Szenario: Das erste Signal trifft ein, nachdem die Startnachricht gelesen ist und bevor `close(fertig)` läuft; die Lesefrist bleibt gesetzt, und die Anmeldung bricht beim Warten auf den Client ab, obwohl sie zu Ende laufen soll. Das Fenster ist schmal und nicht stabil zu treffen; ein akzeptiertes Negativ mit diesem Grund steht in §7 nicht.
- `verifizierbar`: nein (Zeitfenster), die grüne Mutation ja
- `klasse`: Zusage im Kommentar über den getesteten Teil hinaus

### F-4 INFO

- `quelle`: `AGENTS.md` §3.9 · `pfad`: Plan §7 · "`make gates` (voll, ein Lauf, keine parallelen Läufe) | `4c38fb7`"
- `befund`: Der volle Gate-Lauf in §7 steht auf `4c38fb7`. Danach änderte `e7e6953` Code (`vermittlung.go`) und Tests, `e8790e7` die Abdeckungstabellen. Ein Lauf auf dem letzten Stand ist in §7 nicht belegt; die Liste der Commits im Absatz *Stand* nennt `e7e6953` und `e8790e7` nicht. Die Frage, ob DoD-Punkt 4 erfüllt ist, gehört dem Verifier.
- `klasse`: Beleg älter als der letzte Code-Commit

### F-5 INFO

- `quelle`: [ADR-0041](../plan/adr/0041-anmeldung-im-record-als-rueckruf-des-upstream-ports.md) · `pfad`: [ADR-0041](../plan/adr/0041-anmeldung-im-record-als-rueckruf-des-upstream-ports.md) · "Der Record-Service stellt den Rückruf bereit"
- `befund`: Im Diff stellt der PGWire-Adapter den Rückruf bereit (`func (s *Server) anmeldung`), und der Record-Service reicht ihn nur an den Upstream durch (`s.upstream.Open(ctx, startup, anmeldung)`). Die Entscheidung (Rückruf in `Open`, undurchsichtige Werte, Kern ohne PGWire-Typen, Driven Adapter ohne Driving Adapter) ist eingehalten; `make a-check` und die Mutanten des Implementers belegen sie. Die Wortwahl der ADR weicht ab; die ADR ist `Accepted` und wird nicht überschrieben (`AGENTS.md` §3.5). Keine Aktion.
- `klasse`: Wortlaut der ADR enger als die Umsetzung

### F-6 INFO

- `quelle`: Maintainability, [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) · `pfad`: `internal/adapters/driven/postgres/vermittlung.go` · "unknown authentication type: %d"
- `befund`: `unlesbarerCode` liest den Code der Nachricht `R`, die `pgproto3` nicht liest, aus dem Text des Fehlers der Bibliothek ("unknown authentication type", Endung "is unimplemented"). Ändert ein Update der Bibliothek den Text, fallen Code 6, 7, 9 und die unbekannten Codes auf die allgemeine `PGR-E6001`-Meldung ohne Code; die Tests (`TestVermittlungNichtVermittelt`, M12 rot) würden es bei einem Update melden. Keine Aktion; die Abhängigkeit vom Text steht in keinem Kommentar.
- `klasse`: Abhängigkeit vom Fehlertext einer Bibliothek

### F-7 INFO

- `quelle`: `AGENTS.md` §3.14, Plan §4 · `pfad`: Plan §4 · "über 2500"
- `befund`: Nachgemessen: `git diff --numstat fa5a64c..HEAD -- internal test docs/user` ergibt 1576 Einfügungen und 89 Löschungen, Summe 1665, wie §7 sagt; unter der Meldeschwelle 2000 und der Marke 2500, in der Schätzung des Architect (1800 bis 2400) nicht erreicht. Die Zahl 1623 in §6 ("Stand 1623") nennt den Stand vor den Entscheidungen und ist durch §7 überholt. Keine Aktion.
- `klasse`: Größe innerhalb der Marke

## Negativbefunde (geprüft, ohne Befund)

- `internal/hexagon/` (Modell `anmeldung.go`, Ports, Service): Kern ohne PGWire-Typen und ohne `pgproto3`; `.a-check.yml` unverändert; `make a-check` `gesamt: 0 Befund(e)` (die 33 Hinweise "in keiner Schicht" sind Bestand). Der Rückruf kommt vom Aufrufer bis zum Upstream (S1 des Implementers); Service hält weder Anforderung noch Antwort (S2); `ErrClientEnde` erzeugt keine Session auch mit `LeereSessions` (S3).
- `internal/adapters/driven/postgres/`: Fehlerklassen E6001 (Codes 2, 6, 7, 9, unbekannt, nicht lesbar), E4002 (Reihenfolge, Code ohne Verfahren nach `AuthenticationOk`, zweites Code 11) und E4003 (Senden an den Server) stimmen mit der Spezifikation überein, soweit F-1 und F-2 nicht greifen; nachgefahren rot: M1 (`AuthenticationOk` im SASL-Austausch), M2 (Code 8), M3 (Offset des Rests), M4 (Stelle vor Code), M11 (Code 11 nach Klartext), M12 (Code 6 und 9 aus dem Frontend), M13 (kein Schließen nach dem Abbruch). Kein `Terminate` nach Abbruch; Fehlerantwort des Servers geht mit allen Feldern weiter, ohne Log-Zeile.
- `internal/adapters/driving/pgwire/`: Client-Ereignisse (`Terminate`, Ende des Datenstroms, Reset, Schreibfehler an den Client) liefern `ErrClientEnde` ohne Meldung und ohne Session (M15 und M16 rot); `SetAuthType` wird nach der Anmeldung zurückgesetzt (M17 rot); Wächter endet mit der Startnachricht, Signal und Zwangsende nach dem Plan (außer F-3); kein Aufruf des Loggers im neuen Code, keine Meldung mit Inhalt (Meldungen nennen `%T` oder den Code).
- Zugangsdaten: in `vermittlung.go`, `anmeldung.go` und im Service kein Log-Aufruf, keine Meldung mit Salz, Nonce, Hash, Passwort oder Bytes; die Aufzeichnung enthält keine Antworten (Test `TestRecordAnmeldungNichtInDerAufzeichnung`, Spuren im Integrationstest); die Mutanten des Implementers (S2, E4, W18, W19) sind in den Tests nachvollziehbar, nachgefahren habe ich sie nicht.
- Lesepuffer des Aufbaus: `Open` liefert die Session mit dem Frontend des Aufbaus (`newSession` vor `aufbauen`); Test `TestOpenSessionTraegtLesepufferDesAufbaus` vorhanden.
- `spec/`: Die Entscheidungen des Architect stehen in `468997b` vor dem ersten Code-Commit dazu (`e7e6953`); keine Randform im Code-Commit neu in §6. Kein Artefakt unterhalb in der Spezifikation genannt.
- Plan §3, §6, §7 folgen dem Diff (Dateien, Entscheidungen, Rückgaben); keine Adresse, kein Nehmer neu genannt, der Doku-Folge-Slice und die Matrix-Slices stehen in §1 als Abgrenzung mit Kennung.
- `docs/user/abdeckung-*.md`: nur die Tabellen, die `make abdeckung` erzeugt; `docs/user/benutzerhandbuch.md` und `README.md` unberührt (Doku-Folge-Slice, `AGENTS.md` §3.11).
- `tools/test/run-integration-tests.sh`: Zeile `reject` für `play_reject`; Kommentar passt.
- Suppression-Verbot (§3.2): keine `nolint`-Direktive im Diff; `.golangci.yml` unverändert, keine Gate-Lockerung (§3.6).
- Tote Pfade: keiner gefunden; `inReihenfolge` und `anmeldecode` werden benutzt, `codeSASLAbschluss` im PGWire-Adapter ebenfalls.
