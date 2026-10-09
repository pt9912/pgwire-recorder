# Slice slice-v1-abschluss-einspielen: Einspielen einer Aufzeichnung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten), [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md), [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md)

**Berührte Spec-Stellen:** `LH-FA-20.a` · `LH-FA-17.a` · `LH-FA-03.b` · `LH-FA-02.b` · `LH-FA-01.a` · `LH-FA-12.a` · `LH-FA-14.a` · `LH-FA-18.a` · `SPEC-017` · `SPEC-022` · `SPEC-026` · `SPEC-028` · `SPEC-033` · `SPEC-034` · `SPEC-041` · `ARC-002` · `ARC-003` · `ARC-005` · `ARC-007` · `ARC-009`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `pgwire-recorder play` führt die Client-Anfragen einer Aufzeichnung (einfach und Extended) gegen einen PostgreSQL-Server aus, authentifiziert sich als Client, verbindet sich auf Wunsch mit TLS und verhält sich bei Serverfehlern und Abbruchsignalen wie spezifiziert.

**Übernommen aus `slice-replay-semantik-mismatch`:** `play` kennt die Option `--fail-on-unconsumed` nicht (`PGR-E2001`, Exit-Code 2) und lässt ihre Umgebungsvariable `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` unbeachtet, auch mit ungültigem Wert (`LH-FA-03.b` §Andere Kommandos, `LH-FA-17.a`); der Test dafür gehört zu den Optionen von `play`.

**Übernommen aus `slice-replay-semantik-meldungscodes`:** `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `play`, mit derselben Wertemenge und Strenge wie bei `record` und `replay` (Optionstabelle, `LH-FA-17.a`).

**Bereinigung aus `slice-replay-semantik-fehlerreplay`:** Die Teile von `LH-FA-20.a` Schritt 3, die aufgezeichnete Interaktionen ohne `ReadyForQuery` behandeln, sind seit der Entscheidung in `LH-FA-02.b` §Fehlerantwort vor dem Abbruch nicht mehr erreichbar: Der Recorder schreibt keine solche Interaktion, und der Leser lehnt sie als beschädigt ab (`PGR-E3003`). Bereinigt hat sie der Architect vor dem Code (2026-10-09, Schritt 3 von `LH-FA-20.a`); die Zeile *Aufgezeichnetes Ende* der Tabelle gilt nur mit Vergleich und bleibt bei `slice-v1-abschluss-antwortvergleich` (dort §6). Die DoD sagt für sie nichts zu.

**Übernommen aus `slice-v1-abschluss-konfiguration`** (dort §1, Abgrenzung): Die Optionen von `play` werden über den allgemeinen Leser jenes Slice gelesen, in derselben Priorität wie bei `record` und `replay`; das gehört zu den Optionen von `play` oben.

**Übernommen aus `slice-v1-abschluss-konfigurationsdatei`** (dort §1, Abgrenzung): der Abschnitt `play:` der Konfigurationsdatei, gelesen über das Laden jenes Slice; er gehört zu den Optionen von `play`.

**Übernommen aus `slice-v1-abschluss-upstream-verbinden`** (dort §1, Abgrenzung; hervorgegangen aus den Schnitten von `slice-v1-abschluss-konfigurationsdatei` und `slice-v1-abschluss-verbindungen-platzhalter` vom 2026-10-09; bis zum zweiten Schnitt gab `slice-v1-abschluss-verbindungen-platzhalter` den Punkt): die Wirkung einer benannten Verbindung bei `play`, die `slice-v1-abschluss-verbindungen-platzhalter` beim Laden prüft und jener Slice auflöst und einsetzt — Host und Port, Benutzer und Datenbank der URL mit dem Vorrang von `--user` und `--database`, das Passwort aus dem eingesetzten Platzhalter und TLS zum Upstream nach `sslmode`, wenn `--upstream-tls` nicht gesetzt ist (`LH-FA-17.a`, *Wirkung einer URL*); das gehört zum Ziel oben (*authentifiziert sich als Client, verbindet sich auf Wunsch mit TLS*). Dazu, mit der Kennung `slice-v1-abschluss-upstream-verbinden` (dort §6, U8): Bei `play` werden die Platzhalter aller Teile der benutzten Verbindung eingesetzt, auch in Benutzer, Passwort und Datenbank, über das Einsetzen jenes Slice; eine nicht gesetzte oder leere Variable dort ist `PGR-E2005` mit der Verbindung und dem Namen der ersten Variable in der Reihenfolge der URL (Test bei den Optionen von `play`). Dazu, mit der Kennung `slice-v1-abschluss-upstream-verbinden` (dort §1, Abgrenzung; Befund F-533 aus dessen Review): im Benutzerhandbuch §5 *Konfigurationsdatei* die Wirkung einer Verbindung bei `play`, wie geliefert — das Passwort aus dem Platzhalter, `PGWIRE_RECORDER_PASSWORD`, wenn die Verbindung kein Passwort schreibt, und `sslmode=require` mit TLS und Prüfung des Zertifikats, ein gesetztes `--upstream-tls` vor `sslmode`; jener Slice beschreibt dort nur `record`. Dazu, mit der Kennung `slice-v1-abschluss-upstream-verbinden` (Befund V-125 aus dessen Verifikation): in Beispiel und Abschnittsliste von §5 *Konfigurationsdatei* wieder der Abschnitt `play:` (mit `upstream` und `input`), den jener Slice durch `record:` ersetzt hat, weil das Laden `play:` vor diesem Slice als unbekannten Schlüssel ablehnt; das Beispiel als Datei startet mit `config show` ohne Meldung.

**Übernommen aus `slice-v1-abschluss-anmeldung`** (dort §1, Abgrenzung): die Anmeldung beim Einspielen (`play`); sie gehört zum Ziel oben (*authentifiziert sich als Client*).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Zeitangaben und zeitgetreues Einspielen — `slice-v1-abschluss-zeitangaben`.
- Vergleich der Serverantworten mit der Aufzeichnung — `slice-v1-abschluss-antwortvergleich`.
- Paralleles Einspielen — Out-of-Scope von LH-FA-20.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): Eine Aufzeichnung mit DDL- und DML-Anweisungen, einfach und Extended, wird gegen eine leere Instanz eingespielt, und die Datenbank enthält danach deren Wirkung (Abnahmeszenario 12); mehrere Sessions laufen über eigene Verbindungen nacheinander; Authentifizierung mit Klartext, MD5 und SCRAM-SHA-256 sowie `--upstream-tls` (auch mit `--upstream-ca` für eine eigene Zertifizierungsstelle) funktionieren, eine fehlgeschlagene Anmeldung oder TLS-Pflicht ohne Option meldet `PGR-E4005` (auch bei abgelaufenem oder ungültigem Serverzertifikat), eine unlesbare CA-Datei beim Start `PGR-E2007`, `--upstream-ca` ohne TLS `PGR-E2001`, `--fail-on-unconsumed` ist bei `play` unbekannt (`PGR-E2001`) und ihre Umgebungsvariable bleibt dort unbeachtet, auch mit ungültigem Wert; `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` wirken bei `play` mit derselben Wertemenge, Strenge und Schwelle wie bei `record` und `replay`, ein ungültiger Wert ist `PGR-E2001`, auch in der Umgebungsvariable neben gültiger Option (Integrationstest).
- [ ] Eine Fehlerantwort des Servers bricht ab (`PGR-E4004`, Exit-Code 4); `--continue-on-error` läuft weiter und endet mit Exit-Code 4; `--allow-recorded-errors` lässt aufgezeichnete Fehler zu; ein Verbindungsfehler bricht immer ab (Test).
- [ ] `SIGINT` und `SIGTERM` beenden nach der laufenden Interaktion, mit `--finish-session-on-interrupt` nach der laufenden Session; der Exit-Code ist 0 ohne vorherigen Fehler, sonst 4; ohne Vergleich ist jedes Verbindungsende nach dem ersten `ReadyForQuery` `PGR-E4003`; die Rangfolge mit Vergleich (Exit-Code 5) prüft `slice-v1-abschluss-antwortvergleich` (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` (Play-Service), `internal/hexagon/ports/driving` | neu | Einspiel-Use-Case; nutzt nur Driven Ports |
| `internal/adapters/driving/cli` | update | Kommando `play`, Optionen, Passwort aus der Umgebung, Konfigurationsdatei; `--fail-on-unconsumed` bei `play` unbekannt, ihre Umgebungsvariable unbeachtet (Test, `LH-FA-03.b`); `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `play` wie bei `record` und `replay` (übernommen aus `slice-replay-semantik-meldungscodes`, Test) |
| `internal/adapters/driven/postgres` | update | Authentifizierung und TLS als Client, Nachrichten der Gruppen senden |
| `internal/bootstrap` | update | `play` verdrahten: Upstream-Adapter mit Adresse, Passwort und Zertifikaten, Play-Service, Signale (erstes, zweites), Exit-Code |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-20 |
| `docs/user/benutzerhandbuch.md` | update | §5 *Konfigurationsdatei*: Wirkung einer Verbindung bei `play` (Passwort aus dem Platzhalter, `PGWIRE_RECORDER_PASSWORD`, `sslmode=require`) und der Abschnitt `play:` in Beispiel und Abschnittsliste (V-125), übernommen aus `slice-v1-abschluss-upstream-verbinden` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`, und `slice-v1-abschluss-herunterfahren`, `slice-v1-abschluss-konfiguration`, `slice-v1-abschluss-konfigurationsdatei`, `slice-v1-abschluss-verbindungen-platzhalter` und `slice-v1-abschluss-upstream-verbinden` sind `done` (Signalbehandlung, allgemeiner Leser, Konfigurationsdatei, Laden und Benutzen benannter Verbindungen; alle aus `slice-v1-abschluss-betrieb` hervorgegangen), dazu nach der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) `slice-v1-abschluss-schreiben`, und [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die Authentifizierungsverfahren des Servers sprengen den Slice — zurück zur Zerlegung.
- `in-progress` → `open`: Das Recording enthält Interaktionen, die sich nicht einspielen lassen — Carveout.


## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — je Operation des Laufs, wo sie entschieden ist. Geprüft
vom Architect am 2026-10-09 vor dem ersten Code-Commit, gegen `LH-FA-20.a`, `LH-FA-17.a`,
`LH-FA-14.a`, `SPEC-034`, [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md),
[ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md),
[ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) und den Bestand
(`cli.go`: `optionen`, `leserKommandos`, `lies`; `verbindung.go`: `einsetzen`, `adresseRecord`;
`postgres/upstream.go`: `Open`, `lies`, `verbindungsende`; `bootstrap.go`: `Run`, `fail`). Neu
entschieden heißt: im Commit dieser Prüfung in `LH-FA-20.a` *Randformen je Schritt*, sonst an
der genannten Stelle. Die Marke in eckigen Klammern sagt, zu welchem Teil des vorgeschlagenen
Schnitts (Risiken unten) die Randform gehört: **[K]** Kern, **[A]** Anmeldung, **[T]** TLS.
Offen ist keine.

*Lesen der Optionen (CLI-Adapter)*

- **Optionen von `play`** [K] — `--upstream`, `--input`, `--user`, `--database`,
  `--continue-on-error`, `--allow-recorded-errors`, `--finish-session-on-interrupt`,
  `--log-level`, `--config` am allgemeinen Leser, in der Reihenfolge der Optionstabelle,
  Abschnitt `play:`; [A] ohne eigene Option, [T] `--upstream-tls`, `--upstream-ca`; bestätigt,
  `LH-FA-17.a`. Werte von `--user` und `--database`: jeder nicht leere Text (`artText`).
- **Optionen der Folge-Slices** (`--keep-timing`, `--timing-mode`, `--timing-reference`,
  `--compare-responses`) — bis zu `slice-v1-abschluss-zeitangaben` und
  `slice-v1-abschluss-antwortvergleich` unbekannt (`PGR-E2001`), ihre Schlüssel im Abschnitt
  `play:` unbekannt (`PGR-E2004`), ihre Umgebungsvariablen unbeachtet; derselbe Stand wie bei
  `record` (`RecordOptions`: „soweit dieser Stand sie kennt“). Beide Slices führen die Optionen
  in ihrem §3 (`internal/adapters/driving/cli`); keine Zuweisung.
- **`--fail-on-unconsumed` und `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED`** [K] — unbekannt
  beziehungsweise unbeachtet, auch mit ungültigem Wert; bestätigt, `LH-FA-03.b` §Andere
  Kommandos.
- **Gewöhnliches Argument** [K] — `PGR-E2001`; neu entschieden in `LH-FA-01.a`.
- **Abschnitt `play:` bei anderen Kommandos** [K] — mit `play` am allgemeinen Leser prüft das
  Laden ihn auch bei `record` und `replay`; bestätigt, `LH-FA-17.a` („unabhängig vom
  Kommando“).
- **`--upstream-tls` ausdrücklich `false`** [T] — geht `sslmode=require` vor, aus jeder Quelle
  (Option, Umgebungsvariable, Schlüssel); neu entschieden in `LH-FA-17.a` *Wirkung einer URL*.
  *Hinweis an den Implementer:* Der Leser muss „gesetzt“ vom Standardwert unterscheiden
  (`gelesen.cliOk`, `envOk`, `datei.wert`).
- **`--upstream-ca` ohne TLS** [T] — `PGR-E2001`, in der Reihenfolge nach `--upstream` und vor
  den Variablen der Platzhalter, weil TLS vom `sslmode` der benutzten Verbindung abhängt; neu
  entschieden in `LH-FA-17.a` *Fehler*. Die Meldung nennt die Option, nicht den Pfad.
- **Variablen aller Teile** (U8 aus `slice-v1-abschluss-upstream-verbinden`) [K] —
  `PGR-E2005` mit Verbindung und erster Variable in der Reihenfolge der URL, auch in einem
  Teil, den `--user` oder `--database` überschreibt, und im Passwort auch, wenn der Server
  keines verlangt; neu entschieden (Überschreiben) in `LH-FA-17.a` *Wirkung einer URL*, der
  Rest bestätigt (`einsetzen` über alle Teile).
- **Leere `PGWIRE_RECORDER_PASSWORD`** [A] — nicht gesetzt; neu entschieden in `LH-FA-20.a`
  *Anmeldung*.

*Start (nach den Optionen)*

- **Datei aus `--upstream-ca`** [T] — gelesen nach den Prüfungen von `LH-FA-17.a`, vor der
  Aufzeichnung; nur eine reguläre Datei (Links gefolgt), FIFO, Verzeichnis, Gerät, Socket nicht
  lesbar, ohne Warten; mindestens ein PEM-Block, jeder `CERTIFICATE` mit lesbarem X.509, Text
  außerhalb unbeachtet; sonst `PGR-E2007`, Meldung ohne Pfad und Inhalt; neu entschieden in
  `LH-FA-20.a` *Start*. Ablauf eines Zertifikats erst beim Aufbau (`PGR-E4005`), bestätigt.
- **Aufzeichnung nicht ladbar** [K] — wie bei `replay`, Exit-Code 3; bestätigt, Schritt 1.
- **Aufzeichnung ohne Session mit Interaktion** [K] — keine Verbindung, Exit-Code 0; neu
  entschieden in `LH-FA-20.a` *Start*. Eine Session nur aus Lebendprüfungen wird eingespielt
  (bestätigt, `LH-FA-12.a`).
- **Ausgabe eines Startfehlers** [K] — Zeile beim Prozessende (`fail`), keine Verbindung; neu
  entschieden in `LH-FA-20.a` *Start*.
- **Signal während des Starts** [K] — bricht ihn nicht ab, danach keine Session; neu
  entschieden in `LH-FA-20.a` *Start*.

*Verbindungsaufbau je Session (Upstream-Adapter)*

- **Startup-Daten** [K] — aufgezeichnete Parameter unverändert, `user` und `database` nach dem
  Vorrang Option, URL, Aufzeichnung; ohne Wert fehlt der Parameter; Serverparameter nicht
  gesendet; neu entschieden in `LH-FA-20.a` *Startup-Daten*.
- **Server nicht erreichbar, Verbindungsende vor dem ersten `ReadyForQuery` ohne
  Fehlerantwort, unerwartete oder nicht lesbare Nachricht im Aufbau** [K] — `PGR-E4002`; neu
  entschieden in `LH-FA-20.a` *Aufbau*. *Hinweis:* `Open` meldet ein Verbindungsende im Aufbau
  heute als `PGR-E4003` und eine Anmeldeanforderung als `PGR-E6001` (für `record`); `play`
  braucht die Einstufung nach `LH-FA-20.a`.
- **Fehlerantwort im Aufbau** [K] — Klasse 28 `PGR-E4005`, sonst `PGR-E4002`, gleich welcher
  Schweregrad; neu entschieden in `LH-FA-20.a` *Aufbau* (die Tabelle nannte `FATAL`).
- **`BackendKeyData`, `CancelRequest`** [K] — verworfen, nie gesendet; neu entschieden in
  `LH-FA-20.a` *Aufbau*.
- **Antwort auf `SSLRequest`** [T] — `S` Aushandlung, `N` `PGR-E4005`, anderes Byte oder Ende
  davor `PGR-E4002`, Bytes nach `S` vor der Aushandlung `PGR-E4002`; jeder Fehler der
  Aushandlung `PGR-E4005`; neu entschieden in `LH-FA-20.a` *TLS*.
- **Name im Zertifikat** [T] — gegen den eingesetzten Host, IPv6 ohne Zone gegen die
  IP-Adressen; neu entschieden in `LH-FA-20.a` *TLS*.
- **Zertifikatsspeicher des Systems nicht ladbar** [T] — gilt als leer (Grenze); neu
  entschieden in `LH-FA-20.a` *TLS*.
- **Passwort fehlt, Server verlangt eines** [A] — `PGR-E4005`, nichts gesendet; **Server
  verlangt keines** — keines gesendet; neu entschieden in `LH-FA-20.a` *Anmeldung*.
- **Nicht unterstütztes Verfahren** (Kerberos, GSSAPI, SSPI, SASL ohne `SCRAM-SHA-256`) [A] —
  `PGR-E4005`, nichts gesendet; SCRAM ohne Channel Binding; neu entschieden in `LH-FA-20.a`
  *Anmeldung*.
- **Fehler im SCRAM-Austausch** (unpassende oder nicht lesbare Nachricht, falsche
  Serversignatur) [A] — `PGR-E4005`; neu entschieden in `LH-FA-20.a` *Anmeldung*.
- **SASLprep** [A] — nicht angewandt (Grenze); **Klartext ohne TLS** — gesendet, wenn
  verlangt (Grenze); neu entschieden in `LH-FA-20.a` *Anmeldung*; das Handbuch nennt beide.

*Interaktion (Play-Service und Upstream-Adapter)*

- **Warten innerhalb einer Extended-Interaktion** [K] — nach `Sync` auf `ReadyForQuery`, nach
  `Flush` auf die Antwort jeder Client-Nachricht der Gruppe (Tabelle der Antworten), nicht
  nach den aufgezeichneten Server-Nachrichten; nach einer `ErrorResponse` nur noch auf das
  `ReadyForQuery`, übrige Gruppen ohne Warten; neu entschieden in `LH-FA-20.a` *Gruppen*.
  *Hinweis:* Die Antwort hängt nicht an der Datenlage (Zeilenzahl), darum nicht an der
  Aufzeichnung.
- **Andere Server-Nachrichten** (`NoticeResponse`, `ParameterStatus`, `NotificationResponse`,
  Zeilen) [K] — gelesen und verworfen; neu entschieden in `LH-FA-20.a` *Interaktion*.
- **`Copy…Response`, nicht lesbare Nachricht** [K] — `PGR-E6001`, Exit-Code 6, sofort; neu
  entschieden in `LH-FA-20.a` (Tabellenzeile, *Interaktion*) und `SPEC-034`.
- **Fehlerantwort `FATAL` oder `PANIC`** [K] — sofort `PGR-E4003`, kein `PGR-E4004`, kein
  Warten auf das Ende; Schweregrad aus `V`, ohne es `S`; neu entschieden in `LH-FA-20.a`
  *Interaktion*. Damit hat das zweite Risiko unten seinen Ort.
- **Erwarteter Fehler bei `--allow-recorded-errors`** [K] — entscheidet allein, ob die
  aufgezeichnete Interaktion irgendwo eine `error_response` trägt, in jeder Gruppe; keine
  Meldung; neu entschieden in `LH-FA-20.a` *Interaktion*.
- **Senden scheitert** [K] — `PGR-E4003`, kein Weiterlesen; **Lesen endet** nach dem ersten
  `ReadyForQuery` — `PGR-E4003`; neu entschieden beziehungsweise bestätigt (Tabelle) in
  `LH-FA-20.a`.
- **Abbruch nach `PGR-E4004`** [K] — keine weitere Antwort gelesen, keine weitere Nachricht oder
  Gruppe, `Terminate`; neu entschieden (Lesen) in `LH-FA-20.a` *Interaktion*, sonst bestätigt,
  Schritt 6.
- **Fortsetzung** (`--continue-on-error`, erwarteter Fehler) [K] — Rest der Interaktion wie
  aufgezeichnet, Exit-Code 4 am Ende nach einem `PGR-E4004`; bestätigt, Schritt 6 und Tabelle.
- **Frist** [K] — keine eigene, auf Antworten wie auf den Aufbau; neu entschieden in
  `LH-FA-20.a` *Interaktion*.

*Meldungen und Log (Bootstrap und Play-Service)*

- **Fehler nach dem Start** [K] — je Fehler, auch je `PGR-E4004` bei `--continue-on-error`,
  eine Log-Zeile `error` mit `code` und `error`; nennt `id`, `sequence`, SQLSTATE und `M` der
  Fehlerantwort, keine weiteren Felder, nie Passwort oder Parameterwerte; neu entschieden in
  `LH-FA-20.a` *Meldungen* und `LH-FA-14.a`.
- **Zeilen der Stufe `info`** [K] — Start (`upstream=host:port`, Aufzeichnung; nie Benutzer,
  Passwort, Datenbank), Ende, erstes Abbruchsignal; neu entschieden in `LH-FA-14.a` und
  `LH-FA-20.a` *Meldungen*.

*Ende einer Session, Signal, Exit-Code*

- **`Terminate` und Schließen scheitern** [K] — keine Meldung, Exit-Code unverändert; offene
  Transaktion setzt der Server zurück; neu entschieden in `LH-FA-20.a` *Ende einer Session*.
- **Erstes Signal im Aufbau** [K] — der Aufbau läuft zu Ende, ein Fehler darin zählt; danach
  ohne `--finish-session-on-interrupt` keine Interaktion, mit ihr die ganze Session; neu
  entschieden in `LH-FA-20.a` *Abbruchsignal*. **Zwischen Sessions** — keine neue (Session
  läuft ab ihrem Aufbau). **In einer Extended-Interaktion** — bis zu ihrem `ReadyForQuery`.
- **Zweites Signal** [K] — `Terminate`, soweit sofort möglich, Schließen, Ende; die
  unterbrochene Interaktion ist kein Fehler, Exit-Code nach der Zeile *Abbruchsignal* (0 ohne
  vorherigen Fehler); weitere Signale ohne Wirkung; neu entschieden in `LH-FA-20.a`
  *Abbruchsignal* (die Tabelle deckte das zweite Signal schon, der Satz macht es ausdrücklich).
- **Exit-Code beim Abbruch nach einem früheren `PGR-E4004`** [K] — der des abbrechenden Fehlers
  (4 oder 6); neu entschieden in `LH-FA-20.a` *Exit-Code*.

*Akzeptierte Negative der Prüfung vom 2026-10-09* (keine Folgepflicht, einmalig und harmlos):

- **Große Ergebnismenge** — keine Zusage, dass `play` eine Antwort streamt: `SPEC-032` ist ein
  Ziel ohne Zusage, und dieselbe Antwort hielt `record` beim Aufzeichnen schon im Speicher.
- **NUL in `--user` oder `--database`** — nur über ein Escape in der Konfigurationsdatei
  erreichbar; der Server lehnt den Aufbau ab (`PGR-E4002` oder `PGR-E4005`), keine eigene
  Prüfung.
- **Mehrere Fehlerantworten in einer Interaktion** — PostgreSQL sendet höchstens eine je
  `Sync` beziehungsweise `Query`; „je Fehlerantwort eine Meldung“ ist über mehrere
  Interaktionen prüfbar, innerhalb einer nicht.
- **Kein Abbruch der Anfrage beim zweiten Signal** — `play` sendet kein `CancelRequest`, der
  Server führt die Anfrage zu Ende, bis er das Verbindungsende bemerkt; kein Fehler des
  Einspielens.
- **TLS-Version und Verfahren** — die Voreinstellung der Standardbibliothek, keine Zusage.
- **Weitere aufgezeichnete Startup-Parameter** (etwa `replication`) — unverändert gesendet,
  wie aufgezeichnet; was der Server daraus macht, ist seine Antwort.

**Risiken:**

- Das Einspielen verändert eine Datenbank; ein Fehlgebrauch gegen eine falsche Instanz ist durch das Handbuch nur gewarnt — **Ausgang:** offen bis Closure.
- Eine Serverantwort mit `FATAL` beendet die Verbindung; die Behandlung gemäß Spezifikation ist erst im Test belegbar — **Ausgang:** offen bis Closure.
- **Größe** (Prüfung des Architect vom 2026-10-09): Der Slice berührt CLI-Adapter, Play-Service
  mit Port, Upstream-Adapter und Bootstrap — mehr als zwei Schichten — und trägt mit Anmeldung
  (drei Verfahren, SCRAM selbst), TLS mit eigener Zertifizierungsstelle, Fehlerregeln,
  Signalen, Konfiguration und Handbuch rund 80 Zusagen; geschätzt **2500 bis 3500 Zeilen**.
  `slice-v1-abschluss-upstream-verbinden` lag mit rund 30 Zusagen bei 712 Zeilen,
  `slice-v1-abschluss-konfigurationsdatei` mit rund 70 bei 1066 und war die Grenze einer
  Review-Sitzung. Vorschlag (Entscheidung des Nutzers, Schnitt durch den Planner): drei
  Slices nach den Marken oben — **Kern** (dieser Slice: Kommando, Optionen, Abschnitt `play:`,
  Verbindung ohne Passwort und ohne TLS, Fehlerregeln, Signale, Abnahmeszenario 12, Handbuch
  V-125), **Anmeldung** (Passwortquellen, Klartext, MD5, SCRAM-SHA-256, `PGR-E4005` der
  Anmeldung, Handbuch zum Passwort aus F-533) und **TLS** (`--upstream-tls`, `sslmode=require`,
  `--upstream-ca`, `PGR-E2007`, `PGR-E4005` von TLS und Zertifikat, Handbuch zu `sslmode` aus
  F-533) — **Ausgang:** offen bis zur Entscheidung des Nutzers vor dem Code.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
