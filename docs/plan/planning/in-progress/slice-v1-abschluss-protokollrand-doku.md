# Slice slice-v1-abschluss-protokollrand-doku: Handbuch und README zu: Protokollversion und Protokollrand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-23`](../../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client), [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-protokollrand`)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-protokollrand` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Inhalt der Sendung:** Der Inhalt aus der Sendung von `slice-v1-abschluss-protokollrand` (Geber, `AGENTS.md` §3.13) ist:
(1) ein Startparameter mit dem Präfix `_pq_.` ist `PGR-E6001`, in `record` und `replay`, auch an einem Recorder mit TLS-Pflicht (dort statt `PGR-E6003`);
(2) die Zählung abgelehnter Verbindungen: Version außer 3.0, nicht lesbarer Rumpf, `_pq_.` und Klartext trotz TLS-Pflicht zählen nicht und verbrauchen im `replay` keine Session;
(3) `NotificationResponse`, auch durch `LISTEN`/`NOTIFY`, ist `PGR-E6001` und die Session wird nicht aufgezeichnet;
(4) die feste Auswahl im Replay: `AuthenticationOk` ohne Anforderung, eine Nachricht `p` nach dem Aufbau ist `PGR-E6001`, `require_auth` als Grenze;
(5) die COPY-Antworten des Servers. Die Beispiele laufen als Datei gegen das gebaute Binary, das Handbuch beschreibt nur, was das Binary zeigt (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-protokollrand`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-protokollrand` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-protokollrand` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-protokollrand` liefert (die fünf Inhalte unter *Inhalt der Sendung* in §1: `_pq_.`-Ablehnung `PGR-E6001` in `record` und `replay`, Zählung abgelehnter Verbindungen, `NotificationResponse`/`LISTEN`-`NOTIFY` mit `PGR-E6001` und nicht aufgezeichneter Session, feste Auswahl im Replay, COPY-Antworten), im Ist-Zustand des gebauten Binaries: Die Abschnitte, die Kommandos, Optionen, Meldungscodes, Beispiele und Grenzen dieses Slice betreffen, stehen im Handbuch, und die Grenz-Sätze, die er überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] `README.md` nennt, was `slice-v1-abschluss-protokollrand` liefert (soweit es die fünf Inhalte aus §1 betrifft), im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-protokollrand` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-protokollrand` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Proben P1 bis P7 und die Randformen R1 bis R8 aus §6 laufen und stehen mit Aufruf und Ergebnis in §7; die Tests, die §6 als Beleg nennt, werden nicht ausgeführt (kein Gate-Lauf in diesem Schritt), der Implementer prüft aber, dass jeder genannte Test die Aussage trägt, für die er genannt ist.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-protokollrand` liegt in `done/` (Voraussetzung; Schritt 33 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von `git diff --numstat cf08408..<letzter Commit des Slice> -- docs/user/benutzerhandbuch.md README.md` liegt über 200 (Metrik nach `AGENTS.md` §3.14; Schätzung des Architect 70 bis 140: Handbuch rund 6 Stellen mit je 5 bis 15 Zeilen, Tabelle der Meldungscodes 4 Zeilen, README 5 bis 10 Zeilen; Prosa ohne Tabellen- oder Gegenprobe-Code, Faktor 1), oder der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Handbuch und README getrennt.
- `in-progress` → `open` (blockiert — Carveout?): Eine Aussage lässt sich nur mit einer Änderung am Binary wahr machen, und der Nutzer will sie nicht streichen; dann zuerst die Entscheidung des Nutzers, welcher Slice das Verhalten liefert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, die Proben und `grep`-Läufe in §7, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect hat vor dem ersten Commit geprüft (2026-10-15): `slice-v1-abschluss-protokollrand` liefert die folgenden weiteren Randformen; sie sind entschieden, der Implementer entscheidet keine davon und hält bei einer nicht genannten an (Rückgabe an den Architect).

**Randformen des Handbuch-Teils** (Entscheidung des Architect nach `AGENTS.md` §3.11 und §3.12; Quelle des Ist-Zustands `spec/spezifikation.md` *Protokollrand*, *Randformen des Protokollrands*, *Zählung der Verbindungen* im Replay und die Tests, die die Zeilen unten nennen; das Handbuch verweist auf keine davon):

| Nr. | Aussage im Handbuch | Stelle im Handbuch | Beleg (Datei oder Probe) |
|---|---|---|---|
| R1 | Ein Parameter mit dem Präfix `_pq_.` in der Startnachricht ist `PGR-E6001` (Fehlerantwort an die Anwendung, Exit-Code 6), in `record` und `replay`; `record` baut dafür keine Verbindung zur Datenbank auf. | Hinweise zum Aufzeichnen (Mehrere Verbindungen), Hinweise zum Wiedergeben, Tabelle der Meldungscodes (`PGR-E6001`) | `TestE2EProtokollrandBeideModi` (beide Modi), `TestE2ERecordErweiterungOhneUpstream` (kein `PGR-E4002`, Exit-Code 6) in `test/integration/protokollrand_e2e_test.go`: genau diese Aussage, Beleg genügt |
| R2 | Mit TLS-Pflicht (`--tls-cert`, ohne `--allow-plaintext`) ist eine unverschlüsselte Startnachricht mit `_pq_.` `PGR-E6001`, nicht `PGR-E6003`. | Verschlüsselte Verbindungen der Anwendung annehmen, Hinweise; Tabelle `PGR-E6003` | Kein vorhandener Test: Probe P1 (Pflicht). Gelingt sie nicht oder zeigt sie anderes, entfällt der Satz und der Fund geht als Zeile in §7 an den Planner. |
| R3 | Nicht gezählt, keine Sitzung und im `replay` keine verbrauchte Sitzung: Startnachricht mit Version außer 3.0, nicht lesbarer Rumpf, `_pq_.`; dazu die Verbindung ohne Startnachricht (steht im Handbuch, bleibt). Die Zählung im `replay` gilt bei `first-request` und `connection`. | Hinweise zum Aufzeichnen (Aufzählung der Nicht-Sitzungen), Hinweise zum Wiedergeben (Aufzählung der Zählung) | Replay: `TestE2EProtokollrandBeideModi` (nach `Version 3.2`, `_pq_.`, Rumpf Länge 9 erhält die reguläre Verbindung die erste Session, mit `connection`). Das Handbuch nennt die Version außer 3.0 mit dem Beispiel 3.2 oder ohne Beispiel, nicht mit einer Liste von Nummern. `first-request` und die Zählung der TLS-Pflicht-Ablehnung: Handbuch-Bestand (`PGR-E6003`), Beleg die vorhandenen E2E-Tests, kein neuer Satz. Record mit `--record-empty-sessions`: Probe P2 (Pflicht); ohne sie steht „auch mit der Option“ nicht dabei. |
| R4 | `LISTEN` und `NOTIFY` allein sind gewöhnliche Anfragen und werden aufgezeichnet. Sendet die Datenbank dazu eine `NotificationResponse` an eine Session, die zuhört (LISTEN und NOTIFY auf denselben Kanal in einer Anfrage, in zwei Anfragen derselben Session oder zwischen zwei Sessions), ist das `PGR-E6001` (Exit-Code 6); die Session, die sie erhalten sollte, wird nicht aufgezeichnet, die sendende zweier Sessions schon. Zeitpunkt: spätestens bei der nächsten Anfrage der zuhörenden Session; kein Satz sagt, dass das Werkzeug sie schon beim Eintreffen liest. Trifft sie nach der letzten Anfrage der zuhörenden Session ein, ist es kein Fehler (Exit-Code 0). Im `replay` sind `LISTEN` und `NOTIFY` aufgezeichnete Anfragen wie jede andere und liefern keine `NotificationResponse`. | Hinweise zum Aufzeichnen; Verbindungsfehler; Tabelle `PGR-E6001`; Grenzen der Aufzeichnung (Hinweise des Abschnitts Mit einem Datenbanktreiber arbeiten) | `TestE2ERecordListenNotify` (PostgreSQL 17): genau diese Fälle. Der Satz zum `replay` trägt `TestE2EProtokollrandBeideModi` nicht; er entfällt, oder Probe P3 (`LISTEN x;` aus einer Aufzeichnung wiedergeben, Antwort ohne `NotificationResponse`, Exit-Code 0). Ein Satz über Postgres-Versionen steht nicht dabei. |
| R5 | `replay` sendet `AuthenticationOk` als erste Nachricht ohne Anforderung eines Passworts; Benutzer, Datenbank und Optionen der Startnachricht prüft es nicht (ein Benutzer, den die Aufzeichnung nicht kennt, wird bedient). Eine Nachricht `p` nach dem Aufbau ist `PGR-E6001`. Ein Client, der eine Anmeldung verlangt (libpq `require_auth` mit einem Verfahren), bricht selbst ab; `replay` meldet das nicht, Exit-Code 0, keine Meldung `PGR-`. `require_auth` ist damit eine Grenze für `replay`, ein weiteres Verfahren (Klartext-Passwort, MD5, SCRAM) bietet es nicht an. | Eine Aufzeichnung wiedergeben, Hinweise; Verschlüsselte Verbindungen (der Satz „Nichts wird geprüft“ gilt mit und ohne TLS, nur wenn P4 es zeigt) | `TestE2EReplayAuswahl` (ohne TLS; Benutzer, Datenbank, Option, `require_auth=scram-sha-256` über pgconn); die Nachricht `p`: `TestE2EProtokollrandBeideModi` (`p` ohne Nullbyte, nach dem Aufbau, beide Modi). Mit TLS: Probe P4 (Pflicht, wenn der Satz „mit und ohne TLS“ im Handbuch steht; sonst steht „ohne TLS geprüft“ dabei nicht, der Satz nennt TLS nicht). |
| R6 | Server-COPY-Antworten: Antwortet die Datenbank im `record` auf eine Anfrage mit einem COPY-Datenstrom (`COPY … TO STDOUT`, `COPY … FROM STDIN`), erhält die Anwendung die Fehlerantwort `PGR-E6001` statt der Antwort der Anfrage, die Sitzung wird nicht aufgezeichnet, eine andere Verbindung schon. Eine COPY-Nachricht der Anwendung (`CopyData`) ist `PGR-E6001`; `CopyDone` und `CopyFail` stehen nicht im Handbuch. `CopyBothResponse` (Replikation) und ob ein `Terminate` zur Datenbank vorausgeht, stehen nicht im Handbuch. | Hinweise zum Aufzeichnen (zwei Stellen), Verbindungsfehler, Tabelle `PGR-E6001`; `play` (Bestand, bleibt) | `TO STDOUT`: `test/integration/record_e2e_test.go` (Zeile 192 ff.), `einfach_e2e_test.go` (Zeile 279 ff.), `sessions_e2e_test.go` (Zeile 491, andere Verbindung aufgezeichnet). `FROM STDIN` im `record`: Probe P5 (Pflicht, sonst steht nur `TO STDOUT` im Handbuch). `CopyData` der Anwendung: `TestE2EProtokollrandBeideModi`. |
| R7 | `FunctionCall` und ein `NotificationResponse` im Aufbau sind nicht Teil der Sendung (§1). | — | Steht nicht im Handbuch; kein Satz dazu. |
| R8 | Abbruchwunsch (`CancelRequest`): Das Werkzeug leitet ihn nicht an die Datenbank weiter, beendet die Verbindung (`PGR-W3001`, kein Fehler, Exit-Code unverändert) und **wertet weder Prozess-ID noch Schlüssel aus**: Eine laufende Anfrage einer anderen Verbindung wird dadurch nicht abgebrochen. | Tabelle `PGR-W3001`; Hinweise zum Aufzeichnen, Zeile zum Abbruchwunsch | `TestCancelRequest` (`server_test.go`) belegt Weiterleitung, Ende der Verbindung und `PGR-W3001`, nicht „wertet nicht aus“. Der Halbsatz „wertet … nicht aus“ steht nur nach Probe P6 (Pflicht); ohne sie steht nur, was `TestCancelRequest` belegt. Eine Länge außer 16 steht nicht im Handbuch (keine Zusage). |

**Proben (Pflicht, einmalig vor der Übergabe, im Scratch-Verzeichnis, Ergebnis je Probe in §7):**

- **P1:** `record --tls-cert … --tls-key …` (ohne `--allow-plaintext`), rohe unverschlüsselte Startnachricht 3.0 mit `_pq_.x=1`: Fehlerantwort mit `PGR-E6001` (nicht `PGR-E6003`), Exit-Code 6 beim Beenden; gleich mit `replay`.
- **P2:** `record --record-empty-sessions` gegen PostgreSQL: Startnachricht Version 3.2, Startnachricht mit `_pq_.x`, Rumpf Länge 9 und danach eine reguläre Verbindung mit einer Anfrage: die Aufzeichnung hat eine Sitzung (`grep -c '^  - id:'` oder die Zeile des Formats), Exit-Code 6.
- **P3:** `replay` mit einer Aufzeichnung, die `LISTEN x;` enthält; Client `psql` oder pgconn: keine `NotificationResponse`, Exit-Code 0. Entfällt mit dem Satz.
- **P4:** wie `TestE2EReplayAuswahl`, aber über TLS (`replay --tls-cert … --tls-key …`): `AuthenticationOk` ohne Anforderung, unbekannter Benutzer bedient.
- **P5:** `record`, Anwendung sendet `COPY t FROM STDIN` (Tabelle vorher angelegt): Fehlerantwort `PGR-E6001`, Sitzung nicht in der Datei, Exit-Code 6.
- **P6:** `record`, Sitzung A läuft mit `SELECT pg_sleep(2)`, Sitzung B sendet einen `CancelRequest` mit Prozess-ID und Schlüssel von A (aus `BackendKeyData`, das A erhielt): A endet ohne Fehler (nicht `57014`), B endet mit `PGR-W3001`. Zeigt A `BackendKeyData` nicht (die Probe lässt sich nicht bauen), entfällt der Halbsatz „wertet … nicht aus“ und geht als Zeile in §7 an den Planner.
- **P7:** Treiber mit der Option für die Protokollversion (libpq 18 `max_protocol_version=3.0` oder pgconn) gegen `record` und `replay`: Verbindung gelingt. Nur dann steht die Option im Hinweis zu `PGR-E6002`.

Jede Probe ist ein Aufruf des gebauten Binarys gegen die genannte Datei oder Verbindung; der Implementer trägt Aufruf und Ergebnis in §7 ein. Die Beispiele des Handbuchs (Aufrufe wie `pgwire-recorder replay --input … --listen …`) laufen als Datei: Die Handbuch-Beispiele dieses Slice sind die vorhandenen Aufrufe der Abschnitte; neu hinzu kommt keine Konfigurationsdatei, der Implementer prüft jede Zeile in einem Codeblock, die er neu schreibt oder ändert, mit `--help` und einem Lauf.

**Meldungscodes-Tabelle** (Abschnitt Meldungscodes, Zeilen `PGR-E6001`, `PGR-E6002`, `PGR-E6003`, `PGR-W3001`): `PGR-E6001` führt zusätzlich: Parameter `_pq_.` in der Startnachricht; nicht lesbare Startnachricht (Länge 8 oder 9, Rumpf ohne abschließendes Nullbyte); `NotificationResponse` der Datenbank beim Aufzeichnen; COPY-Antwort der Datenbank beim Aufzeichnen; Nachricht `p` an `replay` nach dem Aufbau; Hinweis je Fall, was die Anwendung ändert (`_pq_.`: den Treiber ohne Protokollerweiterung konfigurieren; `LISTEN`/`NOTIFY`: nicht in einer aufgezeichneten Sitzung nutzen). `PGR-E6002` führt „Version außer 3.0, auch 3.2“ und den Hinweis, dass das Werkzeug keine Version aushandelt (Treiber auf 3.0 festlegen, bei libpq `max_protocol_version=3.0`; die Option des Treibers steht nur, wenn Probe P7, ein Lauf von `psql` 18 oder pgconn mit der Option gegen das Binary, sie belegt, sonst steht nur „Verwenden Sie einen Treiber, der Version 3.0 nutzt“). `PGR-E6003`: unverändert; der Satz zur Reihenfolge mit `_pq_.` steht dort mit einem Verweis auf `PGR-E6001`. `PGR-W3001`: unverändert, plus Halbsatz nach R8.

**CancelRequest-Wortlaut** (nur nach P6): „Das Werkzeug leitet den Abbruchwunsch nicht an die Datenbank weiter und wertet weder Prozess-ID noch Schlüssel aus; es schließt die Verbindung und warnt mit `PGR-W3001`. Eine laufende Anfrage einer anderen Verbindung bleibt unberührt.“ Ohne P6: „Das Werkzeug leitet den Abbruchwunsch nicht weiter und schließt die Verbindung (`PGR-W3001`).“

**Plattform** (`BEO-REPO/verhalten-nur-unter-linux-geprueft`): Plattformabhängig ist allein, wie eine Verbindung ohne Startnachricht endet (Schließen oder Zurücksetzen der TCP-Verbindung nach 0 bis 7 Bytes, `TestE2EProtokollrandTCPProbe`) und das Warten auf die Fehlerantwort; die Bytefolgen der Protokollnachrichten sind es nicht. Das Handbuch nennt die geprüfte Umgebung in einem Satz im Stil der Zeile zu SQLite: „Das Verhalten am Protokollrand ist unter Linux im Container geprüft; für andere Systeme ist es nicht geprüft.“ Der Satz steht einmal, in den Hinweisen zum Aufzeichnen neben der Aufzählung der Nicht-Sitzungen, ohne Version, Tag, Datum oder „bisher“. Das Risiko behält seinen Ausgang offen bis Closure; ein Test unter einem anderen System ist kein Gegenstand dieses Slice.

**Nicht entschieden, an den Planner:** keine Randform bleibt offen. Der Planner entscheidet nur, falls P1, P4 oder P6 ein Verhalten zeigen, das `spec/spezifikation.md` anders regelt (Verhalten neben der Spezifikation, Zeile in §7).

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** offen bis Closure.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — genau der Gegenstand dieses Slice: Das Handbuch sagt nur zu, was das Binary zeigt.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-protokollrand`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code.

Nachgezählt beim Eintragen der Sendung von `slice-v1-abschluss-protokollrand` (Architect, 2026-10-15, `AGENTS.md` §3.13, Rückgabe V-5): Die fünf Inhalte liegen in den zwei vorhandenen Liefer-Punkten (Handbuch; README); Liefer-Punkte 2 (höchstens 3), Schichten eine (Handbuch und README zählen als eine Schicht; höchstens 2). Der Nehmer liegt nicht darüber, kein Schnitt vorgemerkt.

Zählung des Architect (2026-10-15, nach der Prüfung der Sendung, `AGENTS.md` §3.13): Liefer-Punkte 2 (Handbuch mit Tabelle der Meldungscodes; README), Schichten eine (Dokumentation; kein Code, kein Test, kein Gate). Die Sendung R1 bis R8 in §6 liegt in diesen zwei Punkten; kein Schnitt vorgemerkt. Größe der Rückführungs-Marke in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
