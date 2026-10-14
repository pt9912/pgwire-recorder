# Slice slice-v1-abschluss-tls-client-doku: Handbuch und README zu: TLS zum Client bei record und replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-23`](../../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client), [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-tls-client`)

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-tls-client` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Lieferumfang im Handbuch** (`AGENTS.md` §3.9, nachgezogen nach Review F-1): der Abschnitt *Verschlüsselte Verbindungen der Anwendung annehmen*, die Grenz-Sätze in §1, §4 und §7, die Einstellungen in §5, die Fehlercode-Tabelle, die Sitzungs-Sätze unter *Mehrere Verbindungen aufzeichnen und wiedergeben* (Verbindungen mit `PGR-W3002` und `PGR-E6003` zählen bei `record` und bei `replay` mit `first-request` nicht als Sitzung) und der Satz, dass `config show` die Dateien `tls_cert` und `tls_key` nicht liest.

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-tls-client`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-tls-client` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Übernommen aus `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, offene Frage zu V-3; Planner, 2026-10-12, `AGENTS.md` §3.13): Das Handbuch nennt `--tls-cert`, `--tls-key` und `--allow-plaintext` erst mit `slice-v1-abschluss-tls-client` und nennt dann zu jeder Option ihre Umgebungsvariable mit den Kommandos, bei denen sie gilt; die Variable beim anderen Kommando ist `PGR-E2001` (Satz zur Meldung `PGR-E2001` im Handbuch, §5 *Einstellungen*). Das Handbuch beschreibt die Variable, nachdem das gebaute Binary sie liest.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-tls-client` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-tls-client` liefert, im Ist-Zustand des gebauten Binaries (zu `--tls-cert`, `--tls-key` und `--allow-plaintext` mit Umgebungsvariable und dem Verhalten beim anderen Kommando, aus `slice-v1-abschluss-umgebung-anderes-kommando`; Probe am gebauten Binary): Aus `slice-doku-ist-stand` (dort §6, *Teilweise geliefert*): Der Slice ersetzt im Handbuch die Grenz-Sätze zu TLS zum Client bei `record` und `replay` — §1 *Voraussetzungen* („Die Verbindung zum Werkzeug läuft ohne Verschlüsselung“), §4 *Mit einem Datenbanktreiber arbeiten* (Voraussetzung „ohne Verschlüsselung“, Schritt 3 mit `sslmode=disable`/`prefer`) und §7 *Die Anwendung kann sich nicht verbinden* — und im README den Satz „Alle Verbindungen laufen unverschlüsselt“, soweit er die Verbindung der Anwendung betrifft. Dazu gehören die Sitzungs-Sätze unter *Mehrere Verbindungen aufzeichnen und wiedergeben* (`PGR-W3002` und `PGR-E6003` zählen bei `record` und bei `replay` mit `first-request` nicht als Sitzung) und der Satz zu `config show` (liest `tls_cert` und `tls_key` nicht). Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] `README.md` nennt, was `slice-v1-abschluss-tls-client` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; der Satz „`record` und `replay` nehmen Verbindungen nur ohne TLS an …“ ist durch den Satz ersetzt, der TLS ab 1.2 mit `--tls-cert` und `--tls-key`, die Abweisung einer Anwendung ohne TLS mit `PGR-E6003` und `--allow-plaintext` nennt („Alle Verbindungen laufen unverschlüsselt“ stand schon vor diesem Slice nicht mehr im README; Review F-2), Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
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
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-tls-client` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt; Sitzungs-Sätze unter *Mehrere Verbindungen aufzeichnen und wiedergeben* (`PGR-W3002`, `PGR-E6003` zählen nicht als Sitzung); Satz zu `config show` (liest die Dateien nicht) |
| `README.md` | update | Sätze, die `slice-v1-abschluss-tls-client` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-tls-client` liegt in `done/` (Voraussetzung; Schritt 28 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Handbuch und README getrennt.
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

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect hat `slice-v1-abschluss-tls-client` am 2026-10-13 geprüft; die Randformen dieses Typs Slice sind entschieden. Was das Binary zeigt, gehört in dieses Handbuch im Präsens und nur, soweit der Code-Slice es mit Test liefert (Zusagen in `LH-FA-23.a`, Test im Code-Slice): die Form beider Dateien (zwei Dateien, Zertifikat zuerst mit der Kette, Schlüssel unverschlüsselt, keine Datei für beide), die Quellen (Option, Umgebungsvariable, Schlüssel `tls_cert`, `tls_key`, `allow_plaintext` in `record:` und `replay:`) und das Verhalten beim anderen Kommando, `PGR-E2001` bei einer Kombination, `PGR-E2007` bei einer Datei (ohne Pfad in der Meldung), die Prüfung des Ablaufs nur beim Start, die Abweisung mit `PGR-E6003` und `--allow-plaintext`, `PGR-W3002` ohne Wirkung auf den Exit-Code, `config show` liest die Dateien nicht, und als Grenze: mindestens TLS 1.2, kein Client-Zertifikat, kein direktes TLS (`sslnegotiation=direct`, der Client wird nicht bedient; Grenze, Entscheidung des Nutzers vom 2026-10-13), eine Datei für Zertifikat und Schlüssel zugleich ist ungültig (ebenso). Die Sätze dazu gelten erst, wenn der Code-Slice den Fall mit einem Test belegt (`AGENTS.md` §3.11). Aus der Closure von `slice-v1-abschluss-tls-client` (Planner, 2026-10-13, `AGENTS.md` §3.13), soweit das Binary daran beobachtbar reagiert und der Code-Slice den Fall mit Test belegt: die Dateiarten (Verzeichnis, FIFO, Gerät, Socket sind keine lesbare Datei, der Start wartet nie auf Eingabe), Zertifikatsdatei nur `CERTIFICATE`-Blöcke mit der Kette nach dem Zertifikat des Servers (die Kette wird ausgeliefert, geprüft wird der Ablauf nur am ersten Zertifikat, Beginn und Ende eingeschlossen), Schlüsseldatei genau ein Block (`PRIVATE KEY`, `RSA PRIVATE KEY`, `EC PRIVATE KEY`; RSA, ECDSA, Ed25519; verschlüsselt ungültig), Text außerhalb der PEM-Blöcke unbeachtet, Namen, Verwendung und Aussteller des Zertifikats ungeprüft, der Recorder sendet beim Schließen kein `close_notify`, ein Klartext-Client wird auch nach `GSSENCRequest` mit `PGR-E6003` abgewiesen, ein direkter TLS-Handshake ist `PGR-W3003` auch mit `--allow-plaintext`, und der Schlüssel erscheint in keiner Log-Zeile (auch `--log-level=debug`). Der Weg mit einem selbst signierten Zertifikat (Risiko unten) steht als Beispiel, das als Datei gegen das Binary läuft.

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** offen bis Closure.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** offen bis Closure.
- Clients mit eigener Zertifikatsprüfung scheitern an einem selbst signierten Zertifikat; das Handbuch muss den Weg beschreiben — **Ausgang:** offen bis Closure.

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

### Belege des Implementers

Kein Produkt-Code, kein Test, kein neuer Wächter: Der Slice ändert `docs/user/benutzerhandbuch.md` und `README.md` (Commit `955951b`) und diesen Plan; eine Mutation des Produkts entfällt (`AGENTS.md` §3.10 gilt für neue Verträge). Die Befehlsfolge aus Plan §3 ist kein neuer Wächter; ihre Mutanten stehen in `done/slice-v1-abschluss-einspielen-anmeldung-doku` §7, hier einmal wiederholt (unten), weil Optionen, Variablen und Codes neu hinzukamen.

**Aufbau der Proben** (alles unter Linux im Container, `BEO-REPO/verhalten-nur-unter-linux-geprueft`; nach den Proben entfernt: alle Container, das Netz `impl-tlscd-net` und das Scratch-Verzeichnis, nur Namen mit Präfix `impl-tlscd-`). Binary: `pgwire-recorder:dev` (Stufe `runtime`, `make build` auf `59e14db`; `955951b` ändert nur Handbuch und README). Server: `postgres:17-alpine` mit dem Digest aus `harness/mk/integration.mk` (`trust`), im eigenen Docker-Netz. Clients: `psql` (libpq 17) aus demselben Image, und `python:3.13-alpine` (Modul `ssl`, OpenSSL 3.5.8) für Rohnachrichten (`GSSENCRequest`, zweiter `SSLRequest`, TLS-Versionen, direktes TLS, Ende ohne `close_notify`, Zeit zwischen Handshake und Startnachricht). Zertifikate und Schlüssel mit OpenSSL 3.3.7 (`alpine:3.20`): Wurzel, Zwischenstelle, Blätter (EC, RSA, Ed25519; gültig, abgelaufen, noch nicht gültig, von der Zwischenstelle), Schlüssel als `PRIVATE KEY` (EC, RSA, Ed25519, X25519), `RSA PRIVATE KEY`, `EC PRIVATE KEY`, `ENCRYPTED PRIVATE KEY`, `EC PRIVATE KEY` mit `Proc-Type: 4,ENCRYPTED`. **Ein Fehler der ersten Läufe, behoben:** Die Recorder-Container, gegen die Clients liefen, trugen denselben Netz-Alias `recorder`; liefen zwei zugleich, löste der Alias auf beide auf, und einzelne Ergebnisse stammten vom falschen Container (die ersten Läufe von P4 und P5 mit zwei Containern sind verworfen). P3 lief mit genau einem Recorder im Netz; P4, P5 (wiederholt), P6, P7, P8 und P10 liefen nach einem `clean` vor jedem Start (genau ein Recorder), P1, P2, P9 und P11 starten ihre Container ohne Alias. Die Tabelle nennt nur diese Ergebnisse.

| Nr. | Probe | Aufruf (Kurzform) |
|---|---|---|
| P1 | Dateien beim Start | `replay --tls-cert A --tls-key B` für die Fälle der Tabelle, Ausgang nach 3 s `SIGTERM` (Exit 124 des Zeitgebers heißt: lief) |
| P2 | Kombinationen, Quellen | Optionen, Umgebung, Datei, `--config` mit `-w` (aktuelles Verzeichnis), `play` mit Option, Variable, Schlüssel |
| P3 | Verbindungen zu `record` (Beispiel des Handbuchs) | `psql` mit `require`, `verify-full`, `prefer`, `disable`; danach Stopp und Exit-Code |
| P4 | Kette | `replay` mit Zertifikatsdatei mit und ohne Zwischenzertifikat, `psql verify-full` gegen die Wurzel |
| P5 | Rohnachrichten | `python`: `gssenc`, `second-ssl`, `second-gss`, `maxtls11`, `maxtls12`, `direct`, `closenotify`, `timing`, `tls-close` |
| P6 | direktes TLS | `psql` mit `sslnegotiation=direct`, mit und ohne `--allow-plaintext` |
| P7 | Ablauf nach dem Start | Zertifikat mit Gültigkeit von 14 s; Verbindung vor und nach dem Ablauf |
| P8 | Schlüssel im Log | `record` und `replay` mit `--log-level debug`, drei Verbindungen (gültig, ohne TLS, nicht vertraut), `grep -F` der zweiten Zeile des Schlüssels |
| P9 | Beispiele als Datei | die Blöcke des Abschnitts per `awk` aus dem Handbuch gezogen und ausgeführt (`openssl`-Aufruf, `record`, `psql`, `replay --config`) |
| P10 | Sitzungen | `record --record-empty-sessions` und `replay` mit TLS: Abbruch nach dem `S`, Handshake und Ende, `sslmode=disable`, dann eine gültige Verbindung |
| P11 | Container | Schlüssel mit Rechten `0600` des Benutzers 1000 (Image läuft als 65532), relativer Pfad ohne und mit `-w` |

**Aussage im Handbuch · Test des Code-Slice · Probe · Ergebnis** (jeder Satz mit Verhaltensaussage hat eine Zeile; Tests in `internal/adapters/driving/cli` (CLI), `internal/adapters/driving/pgwire/server_tls_test.go` (PGW), `test/integration/tls_client_e2e_test.go` (E2E)):

| Aussage (Abschnitt *Verschlüsselte Verbindungen der Anwendung annehmen*, sofern nicht anders genannt) | Test | Probe | Ergebnis |
|---|---|---|---|
| Ohne die Dateien nehmen `record` und `replay` nur Verbindungen ohne Verschlüsselung an (Einleitung) | Bestand (`SSLRequest` → `N`) | — | unverändert aus dem Bestand, eng gefasst auf „ohne `--tls-cert`“ |
| Mit beiden Dateien: Verbindung verschlüsselt, Ablauf wie unverschlüsselt; einfache Anfragen ergeben dieselbe Aufzeichnung | `TestE2ERecordTLS` (Aufzeichnung byte-gleich), `TestE2EReplayTLS` | P3, P9 | `require` und `verify-full`: Ergebnis, Aufzeichnung enthält die Anfrage |
| `replay` liefert über beide Verbindungsarten dieselben Antworten | `TestE2EReplayTLS`, `TestTLSReplay` | P9 | `replay --config` mit `psql verify-full`: Zeile `1` |
| `--tls-cert` allein, `--tls-key` allein: `PGR-E2001` | CLI `TestTLSKombination` | P2 | Exit 2, „--tls-cert verlangt --tls-key“ und umgekehrt |
| `--allow-plaintext=true` ohne die Dateien `PGR-E2001`, `=false` gültig | `TestTLSKombination` | P1, P2 | Exit 2 (Option und Umgebung `true`); `=false` und Umgebung `false`: lief |
| Quellen mischen (Zertifikat aus der Umgebung, Schlüssel aus der Datei); Umgebungsvariablen, Schlüssel `tls_cert`, `tls_key`, `allow_plaintext` | `TestTLSQuellen/record`, `…/replay` | P2 | gemischt, nur Umgebung, nur Datei: lief |
| Relativer Pfad ab dem aktuellen Verzeichnis, auch in der Datei | `TestTLSRelativ` | P11 | `-w /pki`: lief; ohne `-w`: Exit 2, „Datei nicht vorhanden“ (beide Richtungen) |
| `play`: Optionen unbekannt (`PGR-E2001`); Variablen `PGR-E2001`; Schlüssel in `play:` `PGR-E2004` | `TestLeserUmgebungAndererKommandos` (Tabelle führt die drei Optionen), `TestTLSQuellen/play` | P2 | Option: „flag provided but not defined“; alle drei Variablen: Exit 2, `PGR-E2001`; `play.tls_cert`: `PGR-E2004` „unbekannter Schlüssel“ |
| Zertifikatsdatei: mindestens ein `CERTIFICATE`, kein Block anderen Typs | `TestTLSDateienFehler/*/--tls-cert/Schlüsselblock` (und 6 weitere), `TestTLSDateienGueltig` | P1 | Schlüssel als Zertifikat und `beides.pem`: Exit 2, „PEM-Block hat nicht den Typ CERTIFICATE“; leere Datei: „kein PEM-Block“ |
| Erster Block ist das Zertifikat des Servers, die weiteren die Kette, die ausgeliefert wird | `TestTLSKetteWirdAusgeliefert` | P4 | Kette in der Datei, Client mit Wurzel `verify-full`: Zeile; nur das Blatt der Zwischenstelle: „certificate verify failed“, Log `PGR-W3002`; Blatt direkt von der Wurzel: Zeile (beide Richtungen) |
| Vertraut die Anwendung nur der Wurzel, gehören die Zwischenzertifikate hinter das Zertifikat (libpq, `verify-full`) | derselbe | P4 | wie oben, beide Richtungen |
| Schlüsseldatei genau ein Block `PRIVATE KEY`, `RSA PRIVATE KEY`, `EC PRIVATE KEY`; RSA, ECDSA, Ed25519 | `TestTLSDateienGueltig/*/…`, `TestTLSDateienFehler/*/--tls-key/zwei_Schlüssel` | P1 | EC (PKCS #8, SEC1), RSA (PKCS #8, PKCS #1), Ed25519: lief; zwei Blöcke: „mehr als ein PEM-Block“; X25519: „weder RSA noch ECDSA noch Ed25519“ |
| Schlüssel nicht verschlüsselt | `TestTLSDateienFehler/*/--tls-key/ENCRYPTED_PRIVATE_KEY`, `…/Proc-Type_4,ENCRYPTED` | P1, P11 | beide Formen: Exit 2, „der Schlüssel ist verschlüsselt“, ohne Rückfrage (auch mit `-i` und leerer Standardeingabe, mit `--log-level debug` genau eine Zeile) |
| Eine Datei für Zertifikat und Schlüssel ist ungültig | `TestTLSDateienFehler/*/--tls-cert/Schlüsselblock` | P1 | `beides.pem` für beide Optionen: Exit 2 |
| Text vor, zwischen, nach den Blöcken unbeachtet, in beiden Dateien | `TestTLSTextAusserhalbDerBloecke/*/vor`, `…/zwischen`, `…/nach` | P1 | Text vor, zwischen, nach den Blöcken (Zertifikatsdatei) und vor, nach (Schlüsseldatei): lief |
| Start: zuerst Zertifikat, dann Schlüssel, erster Fehler gemeldet | `TestTLSKombination/*/beide_Dateien` | P1 | abgelaufenes Zertifikat und verschlüsselter Schlüssel: Meldung nennt `--tls-cert` (Gegenrichtung: gültiges Zertifikat, verschlüsselter Schlüssel: `--tls-key`) |
| `PGR-E2007` bei: nicht lesbar, ohne gültigen Inhalt, abgelaufen, noch nicht gültig, verschlüsselter Schlüssel, Schlüssel passt nicht | `TestTLSDateienFehler`, `TestTLSGueltigkeit/*/Beginn`, `…/Ende`, `TestTLSDateienFehler/*/--tls-key/anderer_Schlüssel_derselben_Art`, `…/Schlüssel_einer_anderen_Art` | P1, P11 | je Exit 2, Option und Ursache in der Meldung, ohne Pfad; Datei mit Rechten `0600` für den Benutzer des Images: „Datei nicht lesbar: permission denied“; Schlüssel einer anderen Art zum EC-Zertifikat: „passt nicht zum Zertifikat“ |
| Jede Datei, Links gefolgt, regulär; Verzeichnis, FIFO, Gerät, Socket `PGR-E2007`; Start wartet nie | `TestTLSDateiKeineRegulaere` | P1 | Verzeichnis, FIFO (beide Optionen), Socket, `/dev/null`: Exit 2 sofort, „keine reguläre Datei“ |
| Die Meldung nennt die Option und die Ursache, nicht den Pfad | `TestE2ETLSStartfehler`, `TestTLSDateienFehler` | P1 | keine Meldung der Läufe nennt einen Pfad |
| Ursache bei mehreren Mängeln nicht festgelegt | — (Spezifikation, §6 des Code-Slice) | — | Satz sagt nur „nicht festgelegt“; keine Zusage |
| Prüfung der Dateien vor der der Ausgabedatei; vorhandene Ausgabedatei unberührt; nicht `PGR-E2002` | `TestE2ETLSStartfehler` | — | Test (vorhandene Datei bleibt, Ausgabe ohne `PGR-E2002`) |
| `config show` liest die Dateien nicht | `TestConfigShowLiestTLSNicht` | P2 | `config show` mit Pfaden, die nicht existieren: Exit 0, Ausgabe nennt die Pfade (Handbuch §4 *Die gewählte Konfigurationsdatei anzeigen*) |
| Gültigkeit nur beim Start, nur am Zertifikat des Servers, Beginn und Ende eingeschlossen | `TestTLSGueltigkeit/*/Beginn`, `…/Ende`, `TestTLSDateienGueltig/*/Kette_mit_abgelaufenem_zweiten_Zertifikat`, `…/Kette_mit_noch_nicht_gültigem_zweiten_Zertifikat` | P1 | Blatt abgelaufen: Exit 2; Blatt noch nicht gültig: Exit 2; Kette mit abgelaufenem zweiten Zertifikat: lief (beide Richtungen) |
| Läuft das Zertifikat nach dem Start ab, werden neue Verbindungen weiter angenommen | `TestTLSZertifikatUngeprueft` | P7 | Gültigkeit 14 s: Verbindung vor und nach dem Ablauf mit `require` (Sitzung 1 und 2 der Aufzeichnung): je eine Zeile; `verify-full` danach: Client meldet „certificate verify failed“; Neustart mit demselben Zertifikat: Exit 2 „abgelaufen“ |
| Namen, Verwendung, Aussteller ungeprüft; Sache der Anwendung | `TestTLSDateienGueltig/*/ohne_Namen`, `…/Verwendung_eines_Clients`, `TestTLSKombination/*/beide_Dateien` (Aussteller) | P1 | Blatt, dessen Wurzel nicht in der Datei steht: lief |
| Antwort `S`, TLS mindestens 1.2; Client mit höchstens 1.1 scheitert | `TestTLSAushandlung`, `TestTLSMindestversion` | P5 | `maxtls12`: TLSv1.2; `maxtls11`: „tlsv1 alert protocol version“, Log `PGR-W3002`, „offered only unsupported versions“ |
| Kein Zertifikat der Anwendung angefordert | `TestTLSAushandlung` | P3 | `psql` ohne Client-Zertifikat verbindet |
| Handshake scheitert (nicht vertraut, abbricht, Klartext nach `S`): `PGR-W3002` mit `code`, `remote`, Ursache; nie Zertifikat oder Schlüssel; Exit-Code unverändert | `TestTLSHandshakeScheitert/*` (Abbruch nach dem `S`, Klartext nach dem `S`, Bytes vor dem `S`), `TestTLSMindestversion` | P3, P4, P5 | Client ohne Vertrauen: Log `level=WARN code=PGR-W3002 remote=… error=…`; `ssl-then-plain`: ebenso; Stopp danach Exit 0 |
| Direktes TLS: `PGR-W3003`, auch mit `--allow-plaintext` | `TestTLSDirekt` | P5, P6 | `direct` (Python) und `psql sslnegotiation=direct` mit `--allow-plaintext`: Log `PGR-W3003`, Client „EOF“ |
| Weiterer `SSLRequest` / `GSSENCRequest` auf der verschlüsselten Verbindung: `PGR-E6001` | `TestTLSZweiterAufbau` | P5 | `second-ssl`, `second-gss`: ErrorResponse und Log `PGR-E6001` |
| Ohne `--allow-plaintext`: Klartext-Anwendung `PGR-E6003`, Verbindung endet, nicht aufgezeichnet, Exit-Code 6 beim Beenden | `TestE2EKlartextAbgewiesen`, `TestTLSUnverschluesselt` | P3, P10 | `sslmode=disable`: `PGR-E6003`, Stopp Exit 6; mit `--record-empty-sessions` keine Sitzung |
| … auch nach `GSSENCRequest` (Antwort `N`) | `TestTLSUnverschluesselt` (`nachGSS`) | P5 | `gssenc`: `N`, dann `PGR-E6003`, Log, Exit 6 |
| … nur für Protokollversion 3.0, sonst `PGR-E6002` | `TestTLSRangfolge/andere_Version` | — | Test |
| Mit `--allow-plaintext` beide Arten | `TestE2ERecordTLS` (gemischt), `TestTLSUnverschluesselt` | P6 | `psql require` und `disable` gegen `--allow-plaintext`: beide angenommen |
| Kein `close_notify`; Anwendung sieht das Ende der TCP-Verbindung | `TestTLSEndeOhneAlert` | P5 | `closenotify`: nach `Terminate` „EOF occurred in violation of protocol“ |
| Ende ohne `close_notify` ist ein Verbindungsende wie ohne TLS | `TestTLSEndeOhneCloseNotify` | — | Test |
| Ruhende verschlüsselte Verbindung hält das Beenden nicht auf; abgeschlossene Anfragen in der Aufzeichnung | `TestE2ERecordTLSBeendenMitOffenerVerbindung` | — | Test |
| Schlüssel in keiner Log-Zeile, auch `debug` | `TestRunTLSLogOhneSchluessel/record`, `…/replay` | P8 | 0 Treffer (record 5, replay 6 Zeilen), Gegenprobe des Suchmusters auf die Schlüsseldatei: 1 Treffer |
| Zeit zwischen Handshake und Startnachricht zählt nicht zu `offset_ms` (Bestandssatz, Handbuch *Zeitangaben aufzeichnen*, unverändert) | kein Test mit TLS; der Satz bleibt durch die Probe belegt | P5 `timing` | 1,5 s nach dem Handshake bis zur Startnachricht, `offset_ms: 3` der ersten Anfrage |
| Einstellungen: Tabelle in §5, Variablen bei `play` | `TestLeserUmgebungAndererKommandos`, `TestTLSHilfe` | P2, Befehlsfolge | Hilfe → Handbuch leer |
| Fehlercodes `PGR-E6003`, `PGR-W3002`, `PGR-W3003`, `PGR-E2007`, `PGR-E2001`, `PGR-E2004`, `PGR-E6001` | siehe die Zeilen oben | siehe oben | siehe oben |
| Sitzungen: Verbindung mit gescheitertem TLS-Aufbau oder mit `PGR-E6003` ist keine Sitzung / zählt nicht (`record`, `replay` `first-request`) | `TestTLSHandshakeScheitert` und `TestTLSUnverschluesselt` (Sessions 0) | P10 | `record --record-empty-sessions`: Abbruch nach dem `S`, Handshake und Ende, `disable`: nur die gültige Verbindung ist Sitzung 1; `replay first-request`: danach erhält die gültige Verbindung Sitzung 1 |
| Fehlerbehebung: `PGR-E6003` heißt Anwendung beginnt ohne Verschlüsselung; `PGR-W3002`; `PGR-W3003` | wie oben | P3, P5, P6 | wie oben |
| Beispiele als Datei (`openssl`, `record`, `psql`, `replay --config`) | — (Gate: Schritt P9) | P9 | die vier Blöcke unverändert ausgeführt (`openssl`: `server.pem` beginnt mit `BEGIN CERTIFICATE`, `server.key` mit `BEGIN PRIVATE KEY`; `record` Exit 0 mit Aufzeichnung `sql: select 1 as eins`; `psql verify-full` Zeile `1`; `replay --config` mit `psql` Zeile `1`; `config show` mit der Datei ohne Fehler) |

**Quantoren-Durchsicht des Bestands** (Merkmal: die Allgemeinregeln, die der neue Abschnitt einschränkt: Verschlüsselung der Anwendungsseite, Zählung von Verbindungen, Sitzungen, Exit-Code, „nie nach Eingaben fragen“; Suchbegriffe `grep -n -i -E "verschlüssel|\btls\b|unverschlüsselt"` über das Handbuch und die Quantoren „jede/alle/immer/nie/keine/nur“ über den Diff):

| Fundstelle | Entscheidung |
|---|---|
| §1 *Voraussetzungen*, zwei Sätze („Weitere Änderungen … nicht nötig“, „läuft ohne Verschlüsselung“) | eng gefasst („Außer an der Verschlüsselung“, „solange … ohne `--tls-cert`“), verbunden mit dem neuen Abschnitt |
| §4 *Mit einem Datenbanktreiber arbeiten* (Einleitung „nur Host und Port“, Voraussetzung, Schritt 3) | eng gefasst, verbunden |
| §4 *Mehrere Verbindungen*, Aufzeichnen („Keine Sitzung ist …“), die Zeile zur Anfrage nach Verschlüsselung, Wiedergeben („zählt ab ihrer Startnachricht“) | eng gefasst auf „ohne `--tls-cert`“; ergänzt um `PGR-W3002` und `PGR-E6003` mit P10 |
| §4 *Zeitangaben aufzeichnen*, „auch nicht nach einer Anfrage nach Verschlüsselung“ | unverändert, mit Probe (P5 `timing`) |
| §4 *Einspielen*, „Ohne Wunsch nach TLS … ohne TLS mit der Datenbank“, `record` verbindet ohne TLS zur Datenbank (§5) | unverändert: Aussagen über die Seite der Datenbank, die TLS zur Anwendung nicht berührt (`record` mit TLS und ohne TLS zur Datenbank in P3, P9) |
| §3 *Das Werkzeug fragt nie nach Eingaben* (Testautomatisierung) | unverändert, mit Probe: verschlüsselter Schlüssel endet mit Exit 2 ohne Rückfrage (P1) |
| §5 Einleitung („Die Umgebungsvariable einer Option, die die Betriebsart nicht kennt …“) | ergänzt um die drei Variablen bei `play` (P2) |
| §5 *Konfigurationsdatei* („Ein relativer Pfad gilt relativ zum aktuellen Verzeichnis“, „Wahrheitswerte“) | unverändert: gilt für `tls_cert`, `tls_key` (P2, P11) und `allow_plaintext` (`yes` → `PGR-E2004`, P2) |
| §5 *Log-Ausgaben*, „Scheitert der Start, … genau den Fehlertext“ | unverändert, mit Probe: `PGR-E2007` mit `--log-level debug` eine Zeile (P1) |
| §5 *Exit-Codes*, „6 nicht unterstützte Funktion“ | unverändert: `PGR-E6003` ist Klasse „nicht unterstützt“, Exit 6 (P3, P5) |
| §6 *Rollen und Rechte*, „nimmt jede Anmeldung an“ | unverändert: TLS fordert kein Zertifikat der Anwendung an (P3) |
| §7 *Fehlercodes*, `PGR-E4005` / `PGR-E4002` (TLS zur Datenbank) | unverändert: Seite der Datenbank |

**Befehlsfolge aus Plan §3** (Optionen, Variablen, Codes des Handbuchs gegen `--help` und Katalog, Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`, ergänzt um die Gegenrichtungen Hilfe → Handbuch und Katalog → Handbuch):

| Stand | Ergebnis |
|---|---|
| `59e14db` (vorher) | rot: Hilfe → Handbuch: `--allow-plaintext`, `--tls-cert`, `--tls-key`, Variablen `PGWIRE_RECORDER_ALLOW_PLAINTEXT`, `_TLS_CERT`, `_TLS_KEY`; Katalog → Handbuch: `PGR-E6003`, `PGR-W3002` |
| `955951b` (nachher) | Optionen des Handbuchs, die `--help` nicht kennt: nur `--h`, `--help`, `--version`, `--rm`, `--name`; Variablen mit der benannten Ausnahme, Codes und beide Gegenrichtungen: leer |
| Mutant (Kopie des Handbuchs mit angehängter Zeile `--tls-ca` `PGWIRE_RECORDER_TLS_CA` `PGR-E2099`) | je rot: `--tls-ca` unter den Optionen, `PGWIRE_RECORDER_TLS_CA` unter den Variablen, `PGR-E2099` unter den Codes |

`grep -n -i` im Handbuch nach Slice, Welle, ADR-, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review, „noch nicht“, „kommt“, „geplant“, „später“, „bisher“, „künftig“: Treffer nur im Bestand (Zeilen 14 „später“, 167 und 1214/1595 „noch nicht“ als Zustand des Zertifikats beziehungsweise der Zieldatei, 533/613 „später“ im Sinn der Reihenfolge, 712/1109/1375/1219 „kommt“ im Sinn von Herkunft, Eintreffen und Reihenfolge); keiner im Sinn von Künftigem. README: kein Treffer auf Slice, Welle, „noch nicht“, „geplant“, „bisher“. `make docs-check`: 732 Dateien, 0 Befunde.

**Abweichungen und Funde für den Planner:**

- **DoD-Zeile 2 (README):** Der Satz „Alle Verbindungen laufen unverschlüsselt“ steht im README nicht mehr (er war schon vor diesem Slice ersetzt); ersetzt habe ich den Satz „`record` und `replay` nehmen Verbindungen nur ohne TLS an …“. Der Wortlaut der DoD ist veraltet.
- **Nicht gefahren, nicht genannt:** Windows und macOS; Treiber außer libpq 17 und Python `ssl` (pgx, JDBC); TLS-Aufbau von `record`/`replay` mit `--session-assignment connection` und gescheitertem Handshake; Schlüssel aus PKCS #8 mit RSA-4096; die Behauptung zu `sslnegotiation=direct` gilt als Beleg für libpq 17 (der Handbuch-Satz nennt die Version nicht).
- **Verhalten neben der Spezifikation:** keines gefunden.
- **Randformen:** Keine nicht genannte oder offene Randform entschieden; der Handbuch-Satz „Ursache bei mehreren Mängeln nicht festgelegt“ gibt die Spezifikation wieder.
- **Abdeckungstabellen:** unberührt (`make abdeckung-check` Teil von `make gates`).

**Läufe und Größe:**

- `make gates` (voll) auf dem sauberen Baum `dbe6232` (Handbuch und README seit `955951b` unverändert, danach nur dieser Plan): grün, Exit 0, 7 min 32 s; darin `make docs-check` (732 Dateien, 0 Befunde), `make test-integration` (`run-integration-tests: gruen`), `make a-check-negativ` (23 Fälle), `make baseline-verify`, `make lint-gegenprobe`, `make stop-hook-gegenprobe` (139 Fälle). Der Stop-Hook gibt auf diesem Stand frei (`{"decision":"approve"}`). Dieser Absatz ist der letzte Commit und ändert nur §7 (Positivliste); die Übergabe-Messung ist `make gates` auf `dbe6232`.
- Engster Sensor je Änderung: `make docs-check` nach jeder Änderung am Handbuch (zuletzt 732 Dateien, 0 Befunde), die Befehlsfolge oben, die Beispiele als Datei (P9).
- Größe (`AGENTS.md` §3.14, Summe aus Einfügungen und Löschungen von `git diff --numstat 59e14db..dbe6232`): Handbuch 226 (196 + 30), README 5 (4 + 1), dieser Plan 106 (105 + 1); ohne den Plan 231, mit dem Plan 337; der Messbereich endet bei `dbe6232`, am letzten Commit `77efd06` des Implementers stehen die Pläne bei 113 (112 + 1) und die Summe bei 344, den Plan-Commits nach dem Review (diese Nachträge) zählt die Messung nicht. Eine Marke nennt §4 nicht in Zeilen („nicht in einer Review-Sitzung prüfbar“).
- Liefer-Punkte 2, Schichten eine (Dokumentation), unverändert gegenüber §8.

Wird bei Closure um die Closure-Notiz ergänzt (vor dem `git mv` nach `done/`).

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
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-tls-client`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code. Nachgezählt beim Eintragen der Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12): Sie gehört zu DoD-Punkt 1 (Handbuch), kein neuer Liefer-Punkt; Liefer-Punkte weiter 2, Schichten weiter eine.

Nachgezählt beim Eintragen der Sendung aus der Closure von `slice-v1-abschluss-tls-client` (Planner, 2026-10-13, `AGENTS.md` §3.13): Die Randformen in §6 gehören zu DoD-Punkt 1 (Handbuch) und fügen keinen Liefer-Punkt hinzu; Liefer-Punkte weiter 2, Schichten weiter eine (Dokumentation).

Nachgezählt beim Eintragen der Handbuch-Stellen aus Review F-1 (Planner, 2026-10-14, `AGENTS.md` §3.13): Sitzungs-Sätze und `config show` gehören zu DoD-Punkt 1 (Handbuch) und fügen keinen Liefer-Punkt hinzu; Liefer-Punkte weiter 2 (≤ 3), Schichten weiter eine (Dokumentation, ≤ 2).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
