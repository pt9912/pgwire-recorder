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
- [x] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-protokollrand` liefert (die fünf Inhalte unter *Inhalt der Sendung* in §1: `_pq_.`-Ablehnung `PGR-E6001` in `record` und `replay`, Zählung abgelehnter Verbindungen, `NotificationResponse`/`LISTEN`-`NOTIFY` mit `PGR-E6001` und nicht aufgezeichneter Session, feste Auswahl im Replay, COPY-Antworten), im Ist-Zustand des gebauten Binaries: Die Abschnitte, die Kommandos, Optionen, Meldungscodes, Beispiele und Grenzen dieses Slice betreffen, stehen im Handbuch, und die Grenz-Sätze, die er überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [x] `README.md` nennt, was `slice-v1-abschluss-protokollrand` liefert (soweit es die fünf Inhalte aus §1 betrifft), im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [x] `make gates` grün. Beleg: voller Lauf nach `393c17d`, Exit 0, Stand `393c17d`, Dauer rund 10 bis 12 Minuten (Entscheidung des Nutzers vom 2026-10-15 zur Closure).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

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
- Die Proben P1 bis P7 und die Randformen R1 bis R8 aus §6 laufen und stehen mit Aufruf und Ergebnis in §7; für Sätze, die kein Test aus §6 trägt, laufen weitere Proben (P8 bis P16, §7); die Tests, die §6 als Beleg nennt, werden nicht ausgeführt (kein Gate-Lauf in diesem Schritt), der Implementer prüft aber, dass jeder genannte Test die Aussage trägt, für die er genannt ist.
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
| R3 | Nicht gezählt, keine Sitzung und im `replay` keine verbrauchte Sitzung: Startnachricht mit Version außer 3.0, nicht lesbarer Rumpf, `_pq_.`; dazu die Verbindung ohne Startnachricht (steht im Handbuch, bleibt). Die Zählung im `replay` gilt bei `first-request` und `connection`. | Hinweise zum Aufzeichnen (Aufzählung der Nicht-Sitzungen), Hinweise zum Wiedergeben (Aufzählung der Zählung) | Replay: `TestE2EProtokollrandBeideModi` (nach `Version 3.2`, `_pq_.`, Rumpf Länge 9 erhält die reguläre Verbindung die erste Session, mit `connection`). Das Handbuch nennt die Version außer 3.0 mit dem Beispiel 3.2 oder ohne Beispiel, nicht mit einer Liste von Nummern. `first-request` und die Zählung der TLS-Pflicht-Ablehnung: Das Handbuch führte die Ablehnung (`PGR-E6003`) nur für `first-request`; für `connection` steht sie seit der Nacharbeit zu V-1 im Satz (Probe P12, beide Zuordnungen), kein Test trägt sie. Record mit `--record-empty-sessions`: Probe P2 (Pflicht); ohne sie steht „auch mit der Option“ nicht dabei. |
| R4 | `LISTEN` und `NOTIFY` allein sind gewöhnliche Anfragen und werden aufgezeichnet. Sendet die Datenbank dazu eine `NotificationResponse` an eine Session, die zuhört (LISTEN und NOTIFY auf denselben Kanal in einer Anfrage, in zwei Anfragen derselben Session oder zwischen zwei Sessions), ist das `PGR-E6001` (Exit-Code 6); die Session, die sie erhalten sollte, wird nicht aufgezeichnet, die sendende zweier Sessions schon. Zeitpunkt: spätestens bei der nächsten Anfrage der zuhörenden Session; kein Satz sagt, dass das Werkzeug sie schon beim Eintreffen liest. Jeder Satz, der die Benachrichtigung als Grund nennt, die Session nicht aufzuzeichnen, trägt die Grenze „bis zum Ende ihrer letzten Anfrage“ (Nacharbeit V-3, Probe P14); die Session ist es auch mit `--record-empty-sessions` nicht. Trifft sie nach der letzten Anfrage der zuhörenden Session ein, ist es kein Fehler (Exit-Code 0). Im `replay` sind `LISTEN` und `NOTIFY` aufgezeichnete Anfragen wie jede andere und liefern keine `NotificationResponse`. | Hinweise zum Aufzeichnen; Verbindungsfehler; Tabelle `PGR-E6001`; Grenzen der Aufzeichnung (Hinweise des Abschnitts Mit einem Datenbanktreiber arbeiten) | `TestE2ERecordListenNotify` (PostgreSQL 17): genau diese Fälle. Der Satz zum `replay` trägt `TestE2EProtokollrandBeideModi` nicht; er entfällt, oder Probe P3 (`LISTEN x;` aus einer Aufzeichnung wiedergeben, Antwort ohne `NotificationResponse`, Exit-Code 0). Ein Satz über Postgres-Versionen steht nicht dabei. |
| R5 | `replay` sendet `AuthenticationOk` als erste Nachricht ohne Anforderung eines Passworts; Benutzer, Datenbank und Optionen der Startnachricht prüft es nicht (ein Benutzer, den die Aufzeichnung nicht kennt, wird bedient). Eine Nachricht `p` nach dem Aufbau ist `PGR-E6001`. Ein Client, der eine Anmeldung verlangt (libpq `require_auth` mit einem Verfahren), bricht selbst ab; `replay` meldet das nicht, Exit-Code 0; die Sitzung, die der abbrechende Client nicht verbrauchte, erscheint beim Beenden als `PGR-W2001`, wenn keine andere Verbindung sie verbrauchte (Probe P15, Nacharbeit V-5). `require_auth` ist damit eine Grenze für `replay`, ein weiteres Verfahren (Klartext-Passwort, MD5, SCRAM) bietet es nicht an. | Eine Aufzeichnung wiedergeben, Hinweise; Verschlüsselte Verbindungen (der Satz „Nichts wird geprüft“ gilt mit und ohne TLS, nur wenn P4 es zeigt) | `TestE2EReplayAuswahl` (ohne TLS; Benutzer, Datenbank, Option, `require_auth=scram-sha-256` über pgconn); die Nachricht `p`: `TestE2EProtokollrandBeideModi` (`p` ohne Nullbyte, nach dem Aufbau, beide Modi). Mit TLS: Probe P4 (Pflicht, wenn der Satz „mit und ohne TLS“ im Handbuch steht; sonst steht „ohne TLS geprüft“ dabei nicht, der Satz nennt TLS nicht). |
| R6 | Server-COPY-Antworten: Antwortet die Datenbank im `record` auf eine Anfrage mit einem COPY-Datenstrom (`COPY … TO STDOUT`, `COPY … FROM STDIN`), erhält die Anwendung die Fehlerantwort `PGR-E6001` statt der Antwort der Anfrage, die Sitzung wird nicht aufgezeichnet, eine andere Verbindung schon. Eine COPY-Nachricht der Anwendung (`CopyData`) ist `PGR-E6001`, ebenso eine Nachricht `p` nach dem Aufbau, und die Verbindung ist keine Sitzung, auch nicht mit `--record-empty-sessions` und auch nicht mit einer abgeschlossenen Anfrage davor (Nacharbeit V-2, Probe P13); `CopyDone` und `CopyFail` stehen nicht im Handbuch. `CopyBothResponse` (Replikation) und ob ein `Terminate` zur Datenbank vorausgeht, stehen nicht im Handbuch. | Hinweise zum Aufzeichnen (zwei Stellen), Verbindungsfehler, Tabelle `PGR-E6001`; `play` (Bestand, bleibt) | `TO STDOUT`: `test/integration/record_e2e_test.go` (Zeile 192 ff.), `einfach_e2e_test.go` (Zeile 279 ff.), `sessions_e2e_test.go` (Zeile 491, andere Verbindung aufgezeichnet). `FROM STDIN` im `record`: Probe P5 (Pflicht, sonst steht nur `TO STDOUT` im Handbuch). `CopyData` der Anwendung: `TestE2EProtokollrandBeideModi`. |
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

**Plattform** (`BEO-REPO/verhalten-nur-unter-linux-geprueft`): Plattformabhängig ist allein, wie eine Verbindung ohne Startnachricht endet (Schließen oder Zurücksetzen der TCP-Verbindung nach 0 bis 7 Bytes, `TestE2EProtokollrandTCPProbe`) und das Warten auf die Fehlerantwort; die Bytefolgen der Protokollnachrichten sind es nicht. Das Handbuch nennt die geprüfte Umgebung in einem Satz im Stil der Zeile zu SQLite: „Das Verhalten am Protokollrand ist unter Linux geprüft; für andere Systeme ist es nicht geprüft.“ (Nacharbeit F-2: ohne „im Container“, weil §7 für den Container keinen Lauf der Proben trug; die Proben P1 bis P11 liefen als Prozess unter Linux, §7) Der Satz steht einmal, in den Hinweisen zum Aufzeichnen neben der Aufzählung der Nicht-Sitzungen, ohne Version, Tag, Datum oder „bisher“. Das Risiko behält seinen Ausgang offen bis Closure; ein Test unter einem anderen System ist kein Gegenstand dieses Slice.

**Nicht entschieden, an den Planner:** keine Randform bleibt offen. Der Planner entscheidet nur, falls P1, P4 oder P6 ein Verhalten zeigen, das `spec/spezifikation.md` anders regelt (Verhalten neben der Spezifikation, Zeile in §7).

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, Registerstand 3× erreicht, Ausgang dort nicht zugewiesen; dieser Slice nennt die Grenze selbst, kein neues Auftreten) — **Ausgang:** *weiter offen* → Register (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, Vermerk in `state.md`): Alle Proben P1 bis P16 liefen als Prozess unter Linux (Implementer und Verifier); das Handbuch sagt „unter Linux geprüft; für andere Systeme nicht geprüft“, also kein Auftreten und kein Entfallen. Der Zähler bleibt bei 3; der Ausgang ist dem Lese-Schritt der Closure von welle-v1-abschluss zugewiesen, nicht diesem Slice.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** *entfallen*: Er trat nicht ein. Die Verifikation bestätigt alle fünf Inhalte der Sendung gegen das gebaute Binary; V-1 (Klartext bei TLS-Pflicht stand nur für `first-request` im Satz der Zählung) fand sie vor der Closure, und die Nacharbeit `ef8f5f2` hat ihn behoben, der Slice schloss nicht mit der Lücke. Die Suche nach den ersetzten Grenz-Sätzen findet nichts Verbliebenes (Befehlsfolge und `grep` in §7).

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

### Belege des Implementers

**Umgebung der Proben** (2026-10-15): Binary aus `make build` (Image `pgwire-recorder:dev`, `pgwire-recorder version` meldet `dev`), mit `docker cp` aus dem Image in ein Scratch-Verzeichnis außerhalb des Repos kopiert und dort als Prozess gegen PostgreSQL gestartet. Upstream: `postgres:17-alpine@sha256:b0f9560a…` (17.11, `POSTGRES_HOST_AUTH_METHOD=trust`) in einem eigenen Docker-Netz, Port auf `127.0.0.1` veröffentlicht. Clients: ein Python-Skript mit rohen Nachrichten (Scratch, nicht im Repo), `psql` 17.11 (`postgres:17-alpine`) und `psql` 18.6 (`postgres:18-alpine`, Image-ID `c293117fcecd`) mit `--network host`. Ein Lauf zur Zeit; Container, Netz und Prozesse danach entfernt. Zeilen, die nicht auf eine Probe zeigen, tragen den Test aus §6.

**Größe** (Metrik §4): `git diff --numstat cf08408..7a0be02 -- docs/user/benutzerhandbuch.md README.md` = 78+23 und 9+1, Summe 111 (Marke über 200 nicht erreicht).

**Proben** (Aufruf, Ergebnis):

| Probe | Aufruf | Ergebnis |
|---|---|---|
| P1 | `record` und `replay` je mit `--tls-cert server.pem --tls-key server.key` (ohne `--allow-plaintext`), rohe unverschlüsselte Startnachricht 3.0 mit `_pq_.x=1`; Gegenprobe ohne `_pq_.`-Parameter | beide Modi: Fehlerantwort `0A000` mit `PGR-E6001` („Protokollerweiterung … `_pq_.`“), Verbindung endet, Exit-Code 6; ohne `_pq_.`: `PGR-E6003`, Exit-Code 6. Im `replay` meldet das Log danach `PGR-W2001` „2 aufgezeichnete Session(s) nie zugeordnet“: die Verbindung hat keine Session verbraucht |
| P2 | `record --record-empty-sessions`: Startnachricht 3.2, Startnachricht mit `_pq_.x`, Rumpf Länge 9 (ein Nullbyte), danach eine reguläre Verbindung mit `SELECT 1;` | 3.2: `PGR-E6002`; `_pq_.x`: `PGR-E6001`; Rumpf: `PGR-E6001`; die reguläre Verbindung läuft durch; Exit-Code 6; Aufzeichnung `empty_sessions: true` mit genau 1 Sitzung (`grep -c '- id:'`) |
| P3 | Aufzeichnung mit zwei Sitzungen (`LISTEN x;`, `NOTIFY x;`, je eigene Verbindung, Exit-Code 0); `replay` ohne Zusatzoptionen: rohe Anfragen je Sitzung; dazu `psql` 17.11 mit `-c "LISTEN x;"` | Antworten `C`, `Z` je Sitzung, keine Nachricht `A`; `psql`: Exit-Code 0, Ausgabe `LISTEN`; `replay` Exit-Code 0, keine Meldung `PGR-` |
| P4 | `replay --tls-cert … --tls-key …`: roh über TLS (SSLRequest, Antwort `S`, Benutzer `gibtesnicht`, Datenbank `nix`, Option `-c x=1`); dazu `psql` 17.11 `sslmode=require` mit `require_auth=password`, `md5`, `scram-sha-256`, `none`; dieselben vier ohne TLS (`sslmode=disable`) | TLS roh: erste Nachricht `R` mit Code 0 (`AuthenticationOk`), `Z`, die Anfrage `SELECT 1;` wird beantwortet. `require_auth=password`/`md5`/`scram-sha-256`: `psql` bricht ab („authentication method requirement … failed: server did not complete authentication“), mit und ohne TLS; `require_auth=none`: Ergebnis `1`, Exit-Code 0; `replay` Exit-Code 0 in beiden Läufen; Meldung `PGR-`: gemessen in der Nacharbeit zu V-5 (P15): keine, wenn eine andere Verbindung die Sitzung zuvor verbrauchte; sonst `PGR-W2001` für die nie zugeordnete Sitzung, nicht für den Abbruch. Läufe ohne den rohen Vorlauf, der die einzige Sitzung verbraucht hätte |
| P5 | `record`, rohe Verbindung A mit `COPY t FROM STDIN` (Tabelle `t` vorher angelegt), danach Verbindung B mit `SELECT 7;`; Wiederholung mit `COPY t TO STDOUT` | A: Fehlerantwort `PGR-E6001` („Serverantwort *pgproto3.CopyInResponse“ bzw. `CopyOutResponse`) statt der Antwort; B läuft durch; Exit-Code 6; Aufzeichnung: 1 Sitzung, `COPY` steht nicht darin, `SELECT 7` steht darin |
| P6 | `record`: A (`SELECT pg_sleep(2);`), B mit `CancelRequest` und Prozess-ID und Schlüssel von A aus `BackendKeyData` | **nicht baubar**: A erhält kein `BackendKeyData` (Nachrichten bis `ReadyForQuery`: `R`, 14 mal `S`, `Z`; Kontrolle direkt gegen PostgreSQL: Abbruch mit `57014` gelingt). Der Halbsatz „wertet … nicht aus“ und der Satz zur anderen Verbindung entfallen; das Handbuch behält `PGR-W3001` unverändert (Wortlaut ohne P6) |
| P7 | `psql` 18.6 mit `max_protocol_version=3.0`, `latest`, `3.2` gegen `record` und `replay` | `3.0`: Verbindung gelingt, Ergebnis `1`, in beiden Modi; `latest` und `3.2`: `PGR-E6002` („Protokollversion 3.2“), `record` und `replay` Exit-Code 6 |
| P8 | `replay` (Aufzeichnung mit `--record-empty-sessions`, zwei Sitzungen `SELECT 1;`, `SELECT 2;`) mit `first-request` und mit `connection`: Startnachricht 3.2, `_pq_.x`, Rumpf Länge 9, danach zwei reguläre Verbindungen mit `SELECT 1;`, `SELECT 2;` | beide Zuordnungen: die abgelehnten Verbindungen erhalten `PGR-E6002`/`PGR-E6001`/`PGR-E6001`, die regulären `T D C Z` ohne `PGR-E5001`/`PGR-E5003`: keine Sitzung verbraucht; Exit-Code 6 |
| P9 | rohe Nachricht `p` mit Rumpf `geheim\0` und rohe Nachricht `d` (`CopyData`) nach dem Aufbau, `record` und `replay` | `PGR-E6001` („Client-Nachricht *pgproto3.PasswordMessage“ bzw. `CopyData`), Verbindung endet (EOF), Exit-Code 6, in beiden Modi |
| P10 | `record`, Startnachrichten Version 3.0: Länge 8 ohne Rumpf; Länge 9 mit Rumpf `x`; Länge 9 mit Rumpf Nullbyte; `user\0postgres\0` ohne abschließendes Nullbyte; gültige Startnachricht (Gegenprobe) | die vier ersten: `PGR-E6001` („Startnachricht nicht lesbar“), EOF; die gültige: `AuthenticationOk` und weiter; Exit-Code 6 |
| P11 | Startnachricht Version 2.0, 3.1, 3.2, 4.0 gegen `record` und `replay` | jede: `PGR-E6002` („Protokollversion … wird nicht unterstützt“), Exit-Code 6 |
| P12 (V-1) | `replay --input` (zwei Sitzungen `SELECT 1;`, `SELECT 2;`, mit `--record-empty-sessions` aufgezeichnet) `--tls-cert server.pem --tls-key server.key`, je `--session-assignment connection` und `first-request`: rohe unverschlüsselte Startnachricht 3.0 (Benutzer `u`, Datenbank `d`), danach `psql` 17.11 mit `sslmode=require`, `SELECT 1;` und in einer zweiten Verbindung `SELECT 2;` | beide Zuordnungen: unverschlüsselt `PGR-E6003`; die beiden TLS-Verbindungen erhalten `1` und `2`, also Sitzung 1 und 2: die Klartext-Verbindung hat bei `connection` wie bei `first-request` keine Sitzung verbraucht; Gegenprobe (`connection`): geht eine reguläre TLS-Verbindung mit `SELECT 1;` voraus, erhält die zweite mit `SELECT 1;` `PGR-E5001`, Exit-Code 5: Das Muster der Probe wird rot, wenn eine vorausgehende Verbindung eine Sitzung verbraucht |
| P13 (V-2) | `record --record-empty-sessions` und, getrennt, ohne die Option; rohe Verbindungen mit Benutzer `postgres`: A `SELECT 101;` (abgeschlossen) dann `p`; B nur `p`; C nur `d` (`CopyData`); D `SELECT 104;` dann `d`; E `SELECT 105;` dann `X`; F (nur mit Option) nur `X` | Mit Option: A bis D `PGR-E6001` (`PasswordMessage` bzw. `CopyData`), Verbindung endet; Datei mit genau 2 Sitzungen, E und F (leer); `SELECT 101` und `SELECT 104` stehen nicht darin. Ohne Option (A `SELECT 201;`+`p`, D `SELECT 204;`+`d`, E `SELECT 205;`): 1 Sitzung, nur `SELECT 205;`. Eine abgeschlossene Anfrage davor macht die Verbindung nicht zur Sitzung |
| P14 (V-3) | `record`: A `LISTEN k;`, B `NOTIFY k;`, danach (a) beide `X` ohne weitere Anfrage; (b) A `SELECT 301;` dann `X`; mit `--record-empty-sessions`: (c) A `LISTEN k; NOTIFY k;` in einer Anfrage, (d) B `LISTEN k;`, C `NOTIFY k;`, B `SELECT 401;` | (a) Exit-Code 0, beide Sitzungen in der Datei (`LISTEN k;`, `NOTIFY k;`); (b) A erhält auf `SELECT 301;` `PGR-E6001` (`NotificationResponse`), Exit-Code 6, Datei: 1 Sitzung (`NOTIFY k;`); (c) und (d): A und B erhalten `PGR-E6001`, Exit-Code 6, Datei mit genau 1 Sitzung (C, `NOTIFY k;`): die zuhörende Verbindung ist auch mit der Option keine Sitzung |
| P15 (V-5) | `replay` (Aufzeichnung mit einer Sitzung `SELECT 1;`) ohne TLS: (a) nur `psql` 17.11 mit `require_auth=password`; (b) zuerst `psql` mit `SELECT 1;` (Ergebnis `1`), dann `require_auth=password` | (a) `psql` bricht ab, `replay` Exit-Code 0, eine Zeile `PGR-W2001` (Sitzung nie zugeordnet); (b) `psql` bricht ab, `replay` Exit-Code 0, keine Zeile `PGR-`. Die Aussage „keine Meldung `PGR-`“ aus P4 gilt nur, wenn eine andere Verbindung die Sitzung verbraucht hat |
| P16 (V-3, Quantor) | `replay --tls-cert`: `CancelRequest` (Länge 16); Startnachricht 3.0 mit `_pq_.x`, 3.2 mit `_pq_.x`, 3.2 ohne, alle unverschlüsselt | `CancelRequest`: Verbindung endet ohne Antwort, kein `AuthenticationOk` (`PGR-W3001`); 3.0 mit `_pq_.x`: `PGR-E6001`; 3.2 mit und ohne `_pq_.x`: `PGR-E6002` |

**Sätze und Belege** (je Verhaltensaussage eine Zeile; Teststellen in `test/integration/protokollrand_e2e_test.go` und den in §6 genannten Dateien, von mir gelesen, nicht gefahren):

| Satz im Handbuch/README | Beleg | Gegenprobe (rot gesehen, wenn Code oder Aufbau anders wäre) |
|---|---|---|
| `_pq_.` in der Startnachricht ist `PGR-E6001` in `record` und `replay`, Exit-Code 6, Fehlerantwort an die Anwendung | `TestE2EProtokollrandBeideModi` (Fall „Parameter _pq_.“, `fehlerAuf` prüft `PGR-E6001`, `p.stop(t, 6)`); P1, P2, P8 | Ohne `_pq_.` zeigt P1 `PGR-E6003`; Fehlerantwort liest der Test aus Feld `M` |
| `record` baut dafür keine Verbindung zur Datenbank auf | `TestE2ERecordErweiterungOhneUpstream` (Upstream `127.0.0.1:1`, Exit-Code 6, kein `PGR-E4002` im Log) | Würde `record` verbinden, wäre es `PGR-E4002`, Exit-Code 4 (Text der Prüfung) |
| mit TLS-Pflicht ist `_pq_.` `PGR-E6001` statt `PGR-E6003`, Exit-Code 6 | P1 (beide Modi) | P1 ohne `_pq_.` ergibt `PGR-E6003` |
| Startnachricht mit Version außer 3.0 (Beispiel 3.2) ist `PGR-E6002`, Exit-Code 6, Verbindung endet | `TestE2EProtokollrandBeideModi` (Fall „Version 3.2“); P2, P7, P11 (2.0, 3.1, 3.2, 4.0) | — |
| Rumpf nicht lesbar (zum Beispiel ohne abschließendes Nullbyte) ist `PGR-E6001` | Test (Fall „Rumpf nicht lesbar“, Länge 9); P2, P10 (vier Formen; die gültige Startnachricht als Gegenprobe läuft weiter) | P10, letzte Zeile |
| diese Verbindungen sind keine Sitzung, auch mit `--record-empty-sessions` | P2 (1 Sitzung in der Datei); `TestE2EProtokollrandBeideModi` (Aufnahme mit Option, reguläre Verbindung danach) | P2 zählt Sitzungen mit `grep`, nicht den Exit-Code |
| im `replay` verbrauchen sie keine Sitzung, bei `first-request` und `connection`; ebenso eine unverschlüsselte Verbindung, die das Werkzeug mit TLS abweist (`PGR-E6003`), bei beiden Zuordnungen (Nacharbeit V-1) | P12 (beide Zuordnungen) und `connection`: `TestE2EProtokollrandBeideModi` (reguläre Verbindung erhält Session 1); beide: P8 (Antworten ohne `PGR-E5001`/`PGR-E5003`), P1 (`PGR-W2001`: Sitzungen nie zugeordnet) | Hätte die abgelehnte Verbindung eine Sitzung verbraucht, wäre die zweite reguläre Anfrage `PGR-E5003` (P8); P12-Gegenprobe: eine reguläre Verbindung davor ergibt `PGR-E5001` |
| „die nächste reguläre Verbindung erhält die Sitzung, die sie ohne sie erhielte“ | P8 | wie vorige Zeile |
| Eine Verbindung, die nach dem Aufbau `p` oder `CopyData` sendet, ist keine Sitzung, auch nicht mit `--record-empty-sessions` und nicht mit einer abgeschlossenen Anfrage davor (Nacharbeit V-2; Handbuch, README, §4) | P13 (A, B, C, D; mit und ohne Option) | P13: dieselbe Verbindung mit `SELECT` und danach `p` fehlt in der Datei, E (`SELECT 105;`) steht darin |
| Nachricht `p` nach dem Aufbau ist `PGR-E6001` (`record` und `replay`), Verbindung endet | `TestE2EProtokollrandBeideModi` („Nachricht p ohne Nullbyte“); P9 (mit Nullbyte, EOF) | P9: ohne `p` kein Fehler (Aufbau in P1 bis P8) |
| `CopyData` der Anwendung ist `PGR-E6001`, Verbindung endet | Test (Fall „CopyData“); P9 | — |
| `LISTEN`, `NOTIFY` allein werden aufgezeichnet | `TestE2ERecordListenNotify` („LISTEN allein“, „NOTIFY allein“) | — |
| Benachrichtigung ist `PGR-E6001` (Exit-Code 6): `LISTEN` und `NOTIFY` in einer Anfrage, in zwei Anfragen derselben Verbindung, auf zwei Verbindungen; die empfangende wird nicht aufgezeichnet, die sendende schon; spätestens auf die nächste Anfrage der zuhörenden | `TestE2ERecordListenNotify` („beides in einer Anfrage“, „zwei Anfragen derselben Session“, „zwei Sessions“: `nicht`/`sieht` und `exit` 6) | Die Fälle „LISTEN allein“ und „Benachrichtigung nach der letzten Interaktion“ (Exit-Code 0) |
| die Benachrichtigung bis zum Ende der letzten Anfrage der zuhörenden Verbindung: die Verbindung ist auch mit `--record-empty-sessions` keine Sitzung (Nacharbeit V-3; Handbuch an sechs Stellen und README) | P14 (b), (c), (d) | P14 (a): trifft sie nach der letzten Anfrage ein, stehen beide Sitzungen in der Datei, Exit-Code 0 |
| nach der letzten Anfrage der zuhörenden Verbindung: kein Fehler, Exit-Code 0 | `TestE2ERecordListenNotify` („Benachrichtigung nach der letzten Interaktion“, beide Sitzungen in der Datei) | — |
| im `replay` sind `LISTEN`/`NOTIFY` aufgezeichnete Anfragen ohne Benachrichtigung | P3 (roh und `psql` 17.11) — der Satz steht im Handbuch **nicht** (Platz im Abschnitt gespart; kein Test trägt ihn über P3 hinaus) | entfällt |
| `replay` beginnt jede angenommene Verbindung mit `AuthenticationOk` ohne Passwort-Anforderung, mit und ohne `--tls-cert`; eine abgelehnte Verbindung erhält sofort die Fehlerantwort (Nacharbeit F-1) | `TestE2EReplayAuswahl` (ohne TLS: `R` mit Code 0); P4 (TLS roh); abgelehnte Verbindung: P1, P8, P10, Nacharbeit F-1 (Lauf als Container) | Nacharbeit F-1: der Satz mit „jede Verbindung“ ist durch die erste Nachricht `E` statt `R` widerlegt |
| Benutzer, Datenbank und Optionen werden nicht geprüft, ein unbekannter Benutzer wird bedient | `TestE2EReplayAuswahl` (`gibtesnicht`, `nix`, `options`); P4 (TLS) | — |
| `require_auth` mit `password`, `md5`, `scram-sha-256`: Client bricht ab, `replay` Exit-Code 0, keine Meldung zum Abbruch, mit psql 17.11 | `TestE2EReplayAuswahl` (`scram-sha-256`, pgconn, kein `PGR-` im Log); P4 (alle drei, mit und ohne TLS) | P4: `require_auth=none` gelingt, `replay` Exit-Code 0; ohne P4 stünde nur `scram-sha-256` |
| `replay` verlangt kein Passwort und bedient eine angenommene Verbindung mit jedem Benutzer und jeder Datenbank (Satz in der Passwort-Anmeldung und §6) | wie die Zeile zu `AuthenticationOk` | der alte Satz „nimmt jedes Passwort an“ ist enger gefasst (Quantor, unten) |
| COPY-Antwort der Datenbank (`TO STDOUT`, `FROM STDIN`) im `record`: `PGR-E6001` statt der Antwort, Sitzung nicht aufgezeichnet, andere Verbindung schon, Exit-Code 6 | `TestE2ERecordNichtUnterstuetzt`, `TestE2ERecordWeitergabeNichtUnterstuetzt` (nur die Fehlerantwort, keine Ergebnisse), `sessions_e2e_test.go` Z. 491 (andere Verbindung); `FROM STDIN`: P5 | P5 (`FROM STDIN` und `TO STDOUT`) |
| `max_protocol_version=3.0` bei libpq 18 (geprüft 18.6) verbindet mit `record` und `replay`; `latest` erhält `PGR-E6002` | P7 | P7 `latest`/`3.2` |
| das Werkzeug handelt keine Version aus | P7, P11 (jede andere Version endet mit `PGR-E6002`, keine Antwort `NegotiateProtocolVersion`) | — |
| Treiber „fordert keine Protokollerweiterung an“ (Voraussetzung), Ratschläge („stellen Sie … ein“) | Handlungsanweisung, keine Verhaltensaussage; Wirkung: P1, P2 | — |
| Verhalten am Protokollrand ist unter Linux geprüft; für andere Systeme nicht | Proben P1 bis P11 als Prozess unter Linux (`uname -sr`: Linux 6.8.0-139-generic), `TestE2EProtokollrandTCPProbe` für die Plattformabhängigkeit; „im Container“ steht nicht mehr im Satz (Nacharbeit F-2); der Lauf von `replay` als Container steht in der Nacharbeit unten | ein Satz „im Container“ hätte einen Lauf aller Proben als Container verlangt; ohne ihn ist der Satz enger gefasst |
| Meldungscodes-Tabelle E6001 bis E6003 | Sätze oben; Codes gegen den Katalog (`comm`, unten) | — |
| README: dieselben Aussagen in Kürze | die Zeilen oben | — |

Der Satz zu `LISTEN`/`NOTIFY` im `replay` (R4, P3) steht nicht im Handbuch: Er hat mit P3 einen Beleg, aber keinen Test, und wurde zugunsten kurzer Hinweise weggelassen. (Planner: Wenn er stehen soll, trägt ihn P3.)

**Quantoren-Durchsicht** (Merkmal: die Allgemeinregeln, die der neue Abschnitt einschränkt — Anmeldung/Passwort beim Wiedergeben, was `record` weiterleitet und aufzeichnet, was eine Sitzung ist, was am Treiber zu ändern ist; Suche nach `jede/alle/immer/stets/nie/keine/jeweils` zusammen mit Verbindung, Sitzung, Anmeldung, Passwort, Startnachricht, Version, COPY, Protokoll, Antwort; je Zeile und Absatz):

| Stelle (vorher) | Entscheidung |
|---|---|
| §1 „Das Werkzeug leitet alles an die echte Datenbank weiter“ | eng gefasst: „die Kommunikation“ (COPY-Datenstrom und Benachrichtigung leitet es nicht weiter) |
| §1 Voraussetzungen „weitere Änderungen an der Anwendung nicht nötig“ | eng gefasst: nur nötig, wenn die Anwendung eine Grenze am Rand des Protokolls berührt, mit Verweis |
| §4 „wird als eigene Sitzung aufgezeichnet … nutzt sie `COPY`, wird sie nicht aufgezeichnet“ (zwei Stellen, README) | verbunden mit dem neuen Satz: COPY-Datenstrom oder Benachrichtigung; „nutzt COPY“ ersetzt, weil nicht der Aufruf, sondern die Antwort der Datenbank die Ursache ist |
| §4 Passwort-Anmeldung „`replay` … nimmt dabei jedes Passwort an“ | eng gefasst: „verlangt dabei kein Passwort“ (P4/Test: kein Verfahren angeboten) |
| §6 „beim Wiedergeben nimmt es jede Anmeldung an, gleich mit welchem Benutzer …“ | eng gefasst: „verlangt es kein Passwort und bedient jeden Benutzer und jede Datenbank“ |
| Hinweise zum Aufzeichnen „Mit der Option wird sie eine Sitzung ohne Interaktionen“ und „mit der Option legt es [die Zieldatei] an“ | eng gefasst: „außer in den Fällen unten“ bzw. „wenn die Verbindung eine Sitzung wird“ (die Ausnahmen stehen in den Aufzählungen darunter, darunter die neue) |
| Hinweise zum Wiedergeben „Verbindung zählt ab ihrer Startnachricht“ | verbunden: die neuen abgelehnten Startnachrichten stehen im selben Satz (P8) |
| §1 FAQ „Die Wiedergabe beantwortet alle aufgezeichneten Anfragen aus der Datei“; Glossar „nicht jede solche Verbindung ist eine Sitzung“; §7 `PGR-E3004` | unverändert: betreffen weder Protokollrand noch Anmeldung |
| Verbindungsfehler „Nutzt beim Aufzeichnen eine Verbindung `COPY` …“ | verbunden mit Benachrichtigung (siehe oben) |
| Sätze mit `play` („jede Verbindung“ …) | unverändert: Betriebsart `play`, vom Protokollrand dieses Slice nicht berührt |

**Greps und Befehlsfolge** (Stand `7a0be02`, `W` im Scratch):

| Prüfung | Ergebnis |
|---|---|
| Befehlsfolge aus `done/slice-doku-ist-stand` §7: Optionen | nur `--h`, `--help`, `--version` (Nicht-Optionen im Handbuch), `--rm`, `--name` (`docker run`/`docker create`), wie dort eingefroren; keine neue Option, kein neuer Codeblock |
| Variablen (`comm -23 $W/v-hb $W/v-help \| grep -vx PGWIRE_RECORDER_PASSWORD`) | leer |
| Codes gegen `internal/hexagon/model/fehler.go` | leer |
| `grep -ciE 'spezifikation\|lastenheft\|\bADR-\|slice\|welle\|review'` Handbuch | 0 (vorher 0) |
| `grep -niE 'noch nicht\|noch keine\|kommt mit\|geplant\|veröffentlichte Version\|ghcr\|homebrew\|macOS\|windows'` Handbuch | 4 Treffer wie vorher (Zieldatei „existiert noch nicht“, „noch nicht gültig“ ×2, „Es gibt keine veröffentlichte Version.“) |
| README `grep -ciE "noch\|erste Version\|Beginn der Umsetzung\|auf Wunsch\|TLS"` | 6 (vorher 6; keine neue Zeile) |
| neue oder geänderte Zeilen in Codeblöcken | keine (nur Prosa und Tabellenzeilen) |
| `make docs-check` | 787 Dateien, 0 Befunde |

**Funde für den Planner:** (1) P6 nicht baubar: Das Binary reicht `BackendKeyData` nicht an die Anwendung weiter, daher entfallen der Halbsatz „wertet weder Prozess-ID noch Schlüssel aus“ und der Satz zur anderen Verbindung; das widerspricht der Spezifikation nicht (Abbruchkennung ist nach `SPEC-004` nicht Teil der Antworten). (2) P1, P4 und P6 zeigen kein Verhalten, das die Spezifikation anders regelt. (3) Der Satz zu `LISTEN`/`NOTIFY` im `replay` steht nicht im Handbuch (siehe oben). (4) Ein Parameter `_pq_.` mit TLS-Pflicht: Der Satz steht im TLS-Abschnitt und in der Tabelle. (5) Die Größe: 111 von 200.

**Mutationen am Code** entfallen: Der Slice liefert keinen Vertrag und ändert keinen Code (§1); die Gegenproben sind die Gegenfälle der Tabellen.

### Nacharbeit zum Review (F-1, F-2)

**F-1 (LOW), Klasse „Quantor in neu geschriebenem Satz nicht gegen die Ausnahmen derselben Seite geprüft“.** Merkmal: ein Quantor über „Verbindung“ im Wiedergabe-Satz, den dieselbe Seite mit abgelehnten Startnachrichten einschränkt. Behoben als Klasse: `replay beginnt jede Verbindung mit AuthenticationOk` ist in README und Handbuch zu „jede angenommene Verbindung“ gefasst (Handbuch zusätzlich: „eine abgelehnte Verbindung erhält sofort die Fehlerantwort“); die zweite Stelle desselben Merkmals, §6 „bedient jeden Benutzer und jede Datenbank“, ist zu „bedient eine angenommene Verbindung mit jedem Benutzer und jeder Datenbank“ gefasst.

Quantoren-Durchsicht der neu geschriebenen Zeilen dieses Diffs (`git diff cf08408..HEAD` auf Handbuch und README, Treffer auf `jede/jeder/jeden/alle/keine/keinen/nie/immer/nur`, je Treffer einzeln):

| Neu geschriebene Stelle | Entscheidung |
|---|---|
| README und Handbuch „`replay` beginnt jede Verbindung mit `AuthenticationOk`“ | eng gefasst: „jede angenommene Verbindung“ (F-1) |
| Handbuch §6 „bedient jeden Benutzer und jede Datenbank“ | eng gefasst: „eine angenommene Verbindung mit jedem Benutzer und jeder Datenbank“ |
| Handbuch „Ein Benutzer, den die Aufzeichnung nicht kennt, wird bedient“ | unverändert: steht in der Aufzählung hinter „jede angenommene Verbindung“ und nennt keinen Quantor; P4 |
| „Keine Sitzung ist, auch mit der Option, eine Verbindung, deren Startnachricht das Werkzeug ablehnt“ | unverändert: P2 (Datei mit genau 1 Sitzung), P8 |
| „verbraucht im `replay` keine“ / „Sie verbraucht keine Sitzung“ (README, Hinweise) | unverändert: P8 (`first-request` und `connection`), P1 (`PGR-W2001`) |
| „Daten zählen bei keiner der beiden Zuordnungen als Verbindung; ebenso eine Verbindung, deren Startnachricht das Werkzeug ablehnt“ | unverändert: P8 beide Zuordnungen |
| „handelt keine andere [Version] aus“ | unverändert: P7, P11 (jede andere Version endet mit `PGR-E6002`) |
| „`record` baut dafür keine Verbindung zur Datenbank auf“ | unverändert: `TestE2ERecordErweiterungOhneUpstream` |
| „`replay` meldet diesen Abbruch nicht und zählt ihn nicht als Fehler“ | unverändert: P4 (Exit-Code 0, keine Meldung `PGR-`) |
| „`replay` verlangt kein Passwort“, „bricht er selbst ab“, „keinen Ablauf aufzeichnen“, „keine Protokollerweiterung“ | unverändert: Verneinung der Aussage, kein Quantor über eine Menge, die eine Ausnahme einschränkt (Ratschläge: Handlungsanweisung); P4 |
| „andere Verbindungen laufen weiter“, „die sendende von zwei Verbindungen schon“, „außer in den Fällen unten“ | unverändert: P2, P5, `TestE2ERecordListenNotify`; kein Quantor |
| „`CopyData` / `p` ... das Werkzeug endet mit Exit-Code 6“ | unverändert: P9 (beide Modi) |

Lauf gegen das gebaute Binary (2026-10-15, `make build`, Image `pgwire-recorder:dev`, Image-ID `9a27f080a356`; `replay` **als Container** mit `--network host`, Aufzeichnung mit einer Sitzung `x`, Python-Skript mit rohen Nachrichten im Scratch): Startnachricht 3.2 ergibt als erste Nachricht `E` mit `PGR-E6002`; `_pq_.x=1` ergibt `E` mit `PGR-E6001`; Rumpf Länge 9 ergibt `E` mit `PGR-E6001`; die reguläre Verbindung danach beginnt mit `R` (`AuthenticationOk`), `Z`. Der Klartext-Fall mit TLS-Pflicht ist durch P1 belegt (`PGR-E6003`, Fehlerantwort), nicht neu gefahren. Das alte „jede Verbindung“ wäre mit der ersten Nachricht `E` statt `R` widerlegt. `make docs-check` und `make gates-schnell` siehe Bericht.

**F-2 (INFO).** Entscheidung: „im Container“ aus dem Satz gestrichen (Handbuch: „ist unter Linux geprüft; für andere Systeme ist es nicht geprüft“), statt einen Lauf aller Proben als Container nachzuholen; das Kleinste, das `AGENTS.md` §3.11 wahrt, denn die Proben P1 bis P11 liefen als Prozess unter Linux (Zeile oben). Der Lauf von `replay` als Container steht oben, trägt den Satz aber nicht allein.

**Größe** (Metrik §4, `git diff --numstat cf08408..HEAD -- docs/user/benutzerhandbuch.md README.md`): 78+23 und 9+1, Summe 111 (Marke über 200 nicht erreicht); die Nacharbeit tauscht Zeilen im Bestand des Diffs.

### Nacharbeit zur Verifikation (V-1 bis V-5)

Umgebung wie oben (Binary aus `make build`, Image-ID `9a27f080a356`, als Prozess unter Linux, PostgreSQL 17.11 Image `postgres:17-alpine@sha256:b0f9560a…`, `psql` 17.11; rohe Nachrichten mit einem Python-Skript im Scratch). Ein Lauf zur Zeit; Container und Netz danach entfernt.

| Finding | Ausgang | Probe |
|---|---|---|
| V-1 (LOW) Klartext trotz TLS-Pflicht verbraucht bei `connection` keine Sitzung, das Handbuch sagte es nur für `first-request` | behoben: die Aufzählung der Zählung nennt „eine unverschlüsselte, die das Werkzeug mit TLS abweist (`PGR-E6003`)“ für beide Zuordnungen; der Zusatz bei `first-request` führt nur noch den gescheiterten TLS-Aufbau (`PGR-W3002`, Bestand) | P12 samt Gegenprobe |
| V-2 (LOW) Verbindung mit `p` oder `CopyData` nach dem Aufbau ist keine Sitzung | behoben: der Satz zu `p`/`CopyData` im Handbuch trägt „keine Sitzung, auch nicht mit der Option und auch nicht nach einer abgeschlossenen Anfrage“; §4, *Mehrere Verbindungen* und README nennen `p` und `CopyData` neben COPY-Datenstrom und Benachrichtigung als Grund, eine Verbindung nicht aufzuzeichnen | P13: die Probe trägt den Satz in beiden Fällen (mit und ohne Option, mit und ohne abgeschlossene Anfrage davor) |
| V-3 (LOW) Quantor „Benachrichtigung, Verbindung nicht aufgezeichnet“ ohne die Ausnahme derselben Seite | behoben: sechs Stellen im Handbuch (§4, *Mehrere Verbindungen* Einleitung, *Verbindungsfehler*, Treiber-Hinweise, Tabelle `PGR-E6001`, der Bullet selbst) und das README sagen „bis zum Ende ihrer letzten Anfrage“; der Bullet nennt zusätzlich „auch nicht mit der Option“ | P14 (a) bis (d) |
| V-4 (INFO) „1×“ für `BEO-REPO/verhalten-nur-unter-linux-geprueft` | berichtigt in §6 (Risiko) und §8 auf den Registerstand: 3× erreicht, Ausgang nicht zugewiesen, kein neues Auftreten; `folge-slice-liefert-handbuch-teil-nicht` bleibt 1× (eine Datei unter `evidence/`) | Zählung der Dateien unter `evidence/` (3 bzw. 1) |
| V-5 (INFO) P4 „keine Meldung `PGR-`“ | präzisiert in P4 und §6 R5: gilt nur, wenn eine andere Verbindung die Sitzung zuvor verbraucht hat. Gemessen: auch mit einer Aufzeichnung aus nur einer Sitzung meldet `replay` beim Beenden `PGR-W2001`, wenn der abbrechende Client die einzige Sitzung nicht verbraucht; die Sitzungszahl entscheidet nicht, sondern ob sie verbraucht wurde. Der Handbuchsatz „`replay` meldet diesen Abbruch nicht“ bleibt wahr und sagt nichts zu `PGR-W2001` | P15 (a), (b) |
| V-6 (INFO) | nichts zu tun | — |

**Klasse „Quantor gegen Ausnahme derselben Seite“, Durchsicht aller neu geschriebenen Sätze des Diffs** (`git diff -U0 cf08408 -- docs/user/benutzerhandbuch.md README.md`, `+`-Zeilen, Treffer auf `jede/jeder/jedes/jeden/jedem/alle/allen/immer/stets/nie/nur/keine/keinen/kein/jeweils/beide/ebenso/auch/ohne/gleich/sonst/außer`; je Treffer einzeln, auch Sätze, die der erste Durchgang als Verneinung abtat; Merkmal: ein Satz über „Verbindung“, „Sitzung“, „aufgezeichnet“, „zählt“ oder „Anmeldung“, dessen Seite eine Ausnahme führt):

| Neuer oder geänderter Satz | Entscheidung |
|---|---|
| §4 / *Mehrere Verbindungen* / README: „zeichnet eine Verbindung, die eine Anfrage abgeschlossen hat, als eigene Sitzung auf“ mit der Aufzählung der Ausnahmen | eng gefasst: Aufzählung um `p`/`CopyData` und die Grenze der Benachrichtigung ergänzt (V-2, V-3); der erste Halbsatz zu abgeschlossenen Anfragen ist Bestand, die übrigen Ausnahmen (gescheiterter Aufbau, zwangsweises Beenden) stehen im Abschnitt dahinter |
| *Mehrere Verbindungen*: „Mit der Option wird sie, außer in den Fällen unten, eine Sitzung ohne Interaktionen“ | verbunden: die Fälle unten nennen jetzt `p`/`CopyData` und Benachrichtigung (V-2, V-3) |
| *Mehrere Verbindungen*: „Keine Sitzung ist, auch mit der Option, eine Verbindung, deren Startnachricht das Werkzeug ablehnt …“ | unverändert: P2, P10, P11 |
| Neuer Bullet `p`/`CopyData` | eng gefasst und um den Satz zur Sitzung ergänzt (V-2, P13) |
| Neuer Bullet `LISTEN`/`NOTIFY` („diese Verbindung wird nicht aufgezeichnet“) | eng gefasst: „bis zum Ende ihrer letzten Anfrage“, „auch nicht mit der Option“ (V-3, P14) |
| *Verbindungsfehler*: „wird die andere aufgezeichnet“ | eng gefasst mit der Grenze der Benachrichtigung (P14) |
| Treiber-Hinweise und Tabelle `PGR-E6001`: Benachrichtigung | eng gefasst (V-3) |
| „Zeichnen Sie keinen Ablauf auf, in dem eine Verbindung eine Benachrichtigung erhält“ (Bullet und Tabelle) | unverändert: Handlungsanweisung an den Schreiber, enger als die Zusage; ein Ablauf mit Benachrichtigung nach der letzten Anfrage bleibt dadurch verboten, nicht zugesagt |
| *Wiedergeben*, Zählung: „ebenso eine Verbindung, deren Startnachricht das Werkzeug ablehnt …“ | verbunden mit der TLS-Pflicht-Ablehnung bei beiden Zuordnungen (V-1, P12) |
| `replay` „beginnt jede angenommene Verbindung mit `AuthenticationOk`“ (README, Handbuch) und „bedient eine angenommene Verbindung mit jedem Benutzer und jeder Datenbank“ | unverändert: P16 (ein `CancelRequest` erhält kein `AuthenticationOk`; das Handbuch zählt ihn nicht als Verbindung, `PGR-W3001`); abgelehnte Verbindungen stehen im selben Absatz |
| Tabelle `PGR-E6003`: „Steht in der Startnachricht ein Parameter `_pq_.`, ist es `PGR-E6001`“ und Satz zum TLS-Abschnitt | unverändert: P1; P16 (3.2 mit `_pq_.` ist `PGR-E6002`; der Satz steht im Handbuch hinter „nur für eine Startnachricht der Protokollversion 3.0“) |
| Tabelle `PGR-E6002`: „handelt keine Version aus“ | unverändert: P7, P11, P16 |
| Voraussetzungen: „weitere Änderungen an der Anwendung nur nötig, wenn sie eine der Grenzen am Rand des Protokolls berührt“ | unverändert: der Satz verweist auf *Mit einem Datenbanktreiber arbeiten*, wo die Grenzen stehen |
| §1: „leitet die Kommunikation an die echte Datenbank weiter und speichert sie“ | unverändert: gleiche Aussage wie vor der Änderung („speichert die Kommunikation“); das Passwort fehlt in der Datei, steht im Handbuch bei der Anmeldung |
| „Das Verhalten am Protokollrand ist unter Linux geprüft; für andere Systeme nicht geprüft“ | unverändert: P1 bis P16 als Prozess unter Linux (`uname -sr`: Linux 6.8.0-139-generic) |
| „`record` baut dafür keine Verbindung zur Datenbank auf“, „`replay` verlangt kein Passwort“, „bricht er selbst ab“, „handelt keine andere aus“ | unverändert: Verneinung ohne Ausnahme auf der Seite; `TestE2ERecordErweiterungOhneUpstream`, P4, P7 |
| Übrige Treffer (`auch nicht mit --log-level debug` u. a.) | Bestand oder ohne Quantor über die Menge der Verbindungen |

**Größe** (Metrik §4, `git diff --numstat cf08408 -- docs/user/benutzerhandbuch.md README.md` im Arbeitsbaum, nach den Änderungen der Nacharbeit): 85+26 und 11+1, Summe 123 (Marke über 200 nicht erreicht).

### Closure

**Review und Verifikation.** Review (`docs/reviews/2026-10-15-review-slice-v1-abschluss-protokollrand-doku.md`: 0 HIGH, 0 MEDIUM, 1 LOW, 1 INFO): F-1 (LOW, „jede Verbindung“ neben den abgelehnten Startnachrichten) behoben in `12b5a72`, Belege in `ffcfdbf`; F-2 (INFO, „im Container“ ohne Lauf) behoben durch Streichen des Zusatzes. Verifikation (`docs/reviews/2026-10-15-verify-slice-v1-abschluss-protokollrand-doku.md`: PASS, 0 HIGH, 0 MEDIUM, 3 LOW, 3 INFO): DoD-Punkte 1, 2 und 4 bestätigt, Punkt 3 durch den Hook getragen und `make gates-schnell` Exit 0 (voller Lauf vom Verifier nicht gefahren). V-1 bis V-3 (LOW) und V-4, V-5 (INFO) behoben in `ef8f5f2` und `28895f1` (Proben P12 bis P16 in §7); V-6 (INFO) ohne Handlung, die Grenze ist richtig eng gefasst. Die Nacharbeit ist vom Verifier nicht erneut gegengelesen; ihre Belege sind die Proben P12 bis P16 des Implementers.

**Größe** (`AGENTS.md` §3.14, `git diff --numstat cf08408..<letzter Commit am Handbuch> -- docs/user/benutzerhandbuch.md README.md`): 85+26 und 11+1, Summe 123 nach der Nacharbeit; Marke 200 nicht erreicht, Schätzung des Architect 70 bis 140 getroffen. Liefer-Punkte 2, Schichten eine (Dokumentation), wie §8.

**Entscheidungen zu den Übergaben.**

- **V-6, P4 „keine Meldung `PGR-`“ (V-5):** V-6 ohne Handlung; V-5 ist in P4 und R5 auf „wenn eine andere Verbindung die Sitzung verbraucht hat“ präzisiert, der Handbuchsatz blieb wahr.
- **`PGR-W2001` beim Abbruch des Clients im `replay`** (Hinweis des Implementers, P15 (a)): Das Handbuch nennt den Fall nicht ausdrücklich — bricht ein Client allein ab und verbraucht die einzige Sitzung nicht, meldet `replay` beim Beenden `PGR-W2001` bei Exit-Code 0. Kein Satz ist falsch, und die allgemeinen Sätze zu nie zugeordneten Sitzungen sowie die Tabellenzeile `PGR-W2001` decken den Fall. Eine Adresse dafür gibt es nicht: `slice-v1-abschluss-verbindungsende-doku` (`grep -n "W2001"` leer) hat das Ende durch den Server zum Gegenstand und nimmt den Abbruch des Clients nicht an; ein Slice dafür wäre ohne Anlass geplant. Ausgang: *weiter offen*, eigener Eintrag `BEO-REPO/replay-meldet-w2001-beim-abbruch-des-clients` (1×, offen, ohne Adresse; Beleg P15, Entscheidung des Nutzers vom 2026-10-15), Verweis in `state.md` von `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (kein Auftreten).
- **Die Mitteilung des Verifiers zu den Risiko-Ausgängen** war ein Vorschlag; die Ausgänge in §6 sind die Entscheidung dieser Closure (Linux: weiter offen; Handbuch-Teil: entfallen).

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge): siehe §6; Beispiel nur unter Linux: *weiter offen* (Register); Teil von Handbuch und README nicht geliefert: *entfallen*.

**Register** (Zähler = Dateien unter `evidence/`, gelesen am Stand dieser Closure):

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — neue Datei `evidence/slice-v1-abschluss-protokollrand-doku.md` (F-1, F-2, V-1, V-2, V-3, V-5): 57 → 58. Verkörpert (`AGENTS.md` §3.11, `.claude/commands/implement-slice.md` Schritt 17 und 19); die Regel zur Quantoren-Durchsicht galt, sie prüfte neu geschriebene Sätze nicht gegen die Ausnahmen derselben Seite. Der neue Satz dazu steht seit Commit `076ded3` in Schritt 17 (siehe Lerneintrag).
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` — kein neues Auftreten; vermerkt in `state.md`. Der Zähler bleibt bei 3; der Ausgang ist dem Lese-Schritt der Closure von welle-v1-abschluss zugewiesen (dort steht er, hier nicht zu lösen).
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` — kein Auftreten (V-1 als „kein Auftreten“ bestätigt, Entscheidung des Nutzers vom 2026-10-15); vermerkt „Teil geliefert“ und den Verweis auf den neuen Eintrag in `state.md`, bleibt 1×, offen.
- `BEO-REPO/replay-meldet-w2001-beim-abbruch-des-clients` — neu angelegt (Kürzel `REPO`), `evidence/slice-v1-abschluss-protokollrand-doku.md` (P15): 0 → 1×, offen, ohne Adresse (Entscheidung des Nutzers vom 2026-10-15).
- `BEO-REPO/schichtteilung-je-plan-verschieden` — kein Auftreten (Zählung in §8 nach `AGENTS.md` §3.13, eine Schicht); nicht geändert.
- Kein nicht verkörperter Eintrag erreicht mit diesem Slice 3× neu (`BEO-REPO/verhalten-nur-unter-linux-geprueft` steht seit `slice-v1-abschluss-sqlite-format-lesen` auf 3×).

**Lerneintrag — geschärfte Regel (verkörpert).** Der Slice schrieb neue Sätze mit Quantor („jede Verbindung“, „wird aufgezeichnet“, „sie wird eine Sitzung“) auf Seiten, die dieselbe Verbindung an anderer Stelle mit Ausnahmen führen (abgelehnte Startnachricht, `p` oder `CopyData` nach dem Aufbau, Benachrichtigung nach der letzten Anfrage); die Durchsicht der Quantoren suchte nach dem Bestand und nach den Wörtern der zuletzt gefundenen Ausnahme, nicht nach den Ausnahmen, die dieselbe Seite schon aufzählt, und fand F-1, V-2 und V-3 erst Review und Verifikation in drei Gestalten derselben Klasse. Verkörpert ist ein Satz in `.claude/commands/implement-slice.md` Schritt 17 (Commit `076ded3`, Entscheidung des Nutzers vom 2026-10-15): Für jeden neu geschriebenen Satz mit Quantor listet §7 zusätzlich die Ausnahmen, die dieselbe Seite für denselben Gegenstand schon aufzählt, und gibt je Ausnahme eine Probe oder fasst den Satz eng; derselbe Satz steht an allen Stellen gleich, an denen er wiederkehrt. **Kein neuer Sensor:** Welche Ausnahmen ein Satz meint, sieht kein Skript. Herkunfts-Anker im Register: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `seit slice-v1-abschluss-protokollrand-doku`. `AGENTS.md` ist nicht geändert.

**Paarungen.** (a) Anker: Die Regel ist verkörpert; Zielort `.claude/commands/implement-slice.md` Schritt 17 trägt den Satz mit `seit slice-v1-abschluss-protokollrand-doku` (`grep` nach der Kennung dort trifft). (b) Folge-Slice: Diese Closure weist keinem Slice etwas zu; `slice-v1-abschluss-verbindungsende-doku` wurde für den Fall zu `PGR-W2001` geprüft und nimmt ihn nicht an, deshalb steht er als eigener Register-Eintrag ohne Adresse und nicht als Zuweisung. (c) Register: die genannten Kennungen existieren als Verzeichnis mit nicht leerem `evidence/` (`zusage-im-kommentar-weiter-als-pruefung`, `verhalten-nur-unter-linux-geprueft`, `folge-slice-liefert-handbuch-teil-nicht`, `replay-meldet-w2001-beim-abbruch-des-clients`, `schichtteilung-je-plan-verschieden`).

**Nachzählen.** Keine Sendung in einen Nehmer, daher kein Nachzählen von Liefer-Punkten und Schichten. `slice-v1-abschluss-verbindungsende-format` (Schritt 35) nennt diesen Slice nur als Start-Voraussetzung und ist unverändert.

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
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (Registerstand 3× erreicht, Ausgang nicht zugewiesen) — Proben laufen nur unter Linux (§6, §7); dieser Slice nennt die Grenze im Handbuch, kein neues Auftreten.
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; `verhalten-nur-unter-linux-geprueft` steht schon auf 3×.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code.

Nachgezählt beim Eintragen der Sendung von `slice-v1-abschluss-protokollrand` (Architect, 2026-10-15, `AGENTS.md` §3.13, Rückgabe V-5): Die fünf Inhalte liegen in den zwei vorhandenen Liefer-Punkten (Handbuch; README); Liefer-Punkte 2 (höchstens 3), Schichten eine (Handbuch und README zählen als eine Schicht; höchstens 2). Der Nehmer liegt nicht darüber, kein Schnitt vorgemerkt.

Zählung des Architect (2026-10-15, nach der Prüfung der Sendung, `AGENTS.md` §3.13): Liefer-Punkte 2 (Handbuch mit Tabelle der Meldungscodes; README), Schichten eine (Dokumentation; kein Code, kein Test, kein Gate). Die Sendung R1 bis R8 in §6 liegt in diesen zwei Punkten; kein Schnitt vorgemerkt. Größe der Rückführungs-Marke in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
