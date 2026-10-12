# Slice slice-v1-abschluss-zeitangaben-doku: Handbuch und README zu: Zeitangaben beim Aufzeichnen und Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben`)

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-zeitangaben-pruefen` (das Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`) liefern, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11). Der Beginn der Messung von `offset_ms` ist in der Lesart des Abschnitts *Messpunkt* der Spezifikation zu beschreiben (Beginn der Session nach TLS und Startup des Clients, nicht der Socket; Sendung von `slice-v1-abschluss-zeitangaben-aufzeichnen`, 2026-10-12, V-189); nachgezählt: Liefer-Punkte und Schichten unverändert, der Satz gehört zu DoD-Punkt 1. Dazu die Lesarten aus `slice-v1-abschluss-zeitangaben` (Sendung vom 2026-10-12, Closure jenes Slice): die Ordnung der Startfehler (eine beschädigte Interaktion ist `PGR-E3003`, Exit-Code 3, und geht `PGR-E2003` vor; Abschnitt *Prüfung der Werte* der Spezifikation), die Regeln der Kombinationen (`--timing-mode` oder `--timing-reference` ohne `--keep-timing` und `--timing-reference` zusammen mit `relative` sind `PGR-E2001`; „gesetzt“ ist, was irgendeine Quelle nennt, auch den Wert des Standards) und die Wertemenge (nur Kleinbuchstaben; jeder andere Wert, auch der leere, ist als Option und Variable `PGR-E2001`, als Schlüssel der Datei `PGR-E2004`; die Meldung nennt den Wert von Variable und Datei nicht); nachgezählt: Liefer-Punkte und Schichten unverändert, die Sätze gehören zu DoD-Punkt 1.

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-zeitangaben`, das am 2026-10-11 nach Entscheidung des Nutzers (vor dem Anspruch zu schneiden, vier Schichten) in `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-uhr`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` geschnitten wurde; der Doku-Teil blieb ungeteilt hier, weil das Handbuch die Zeitangaben als Ganzes beschreibt (Zuschnitt: Planner). Abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-zeitangaben` (damals noch ungeteilt) über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt die drei Teile mit sichtbarem Verhalten voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt hinter dem letzten, `slice-v1-abschluss-zeitangaben`. Der Uhr-Port (`slice-v1-abschluss-zeitangaben-uhr`) ändert kein Verhalten des Binaries; das Handbuch hat dazu nichts zu beschreiben. Zwischen den Code-Slices und diesem hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen. Das Handbuch beschreibt die Zeitangaben nach dem letzten Teil: Es sagt nichts über eine Aufzeichnung, die sich schreiben, aber nicht einspielen lässt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` liefern.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` liefern, im Ist-Zustand des gebauten Binaries: Die Abschnitte, die Kommandos, Optionen, Meldungscodes, Beispiele und Grenzen dieses Slice betreffen, stehen im Handbuch, und die Grenz-Sätze, die er überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7). Die Lesarten aus §1 (Ordnung der Startfehler, Kombinationen, Wertemenge) gehören dazu: Jede steht im Handbuch, geprobt am gebauten Binary.
- [ ] `README.md` nennt, was `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` liefern, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
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
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` liefern: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt; ein neuer Abschnitt *Zeitangaben aufzeichnen und einspielen* mit den Lesarten aus §1, die Zeilen der Tabelle in §5 und der *Fehlercodes*, die SQLite-Tabelle, Einleitung, FAQ und Glossar |
| `README.md` | update | Sätze, die jene drei überholen, im Ist-Zustand |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile *In Arbeit* nennt den Slice (eigener Commit) |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` liegen in `done/` (Voraussetzung, dazu `slice-v1-abschluss-zeitangaben-uhr`, das jene voraussetzen; Schritt 22 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

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

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob einer der drei Slices `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` oder `slice-v1-abschluss-zeitangaben` eine weitere Randform liefert, die dort nicht steht.

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 3×, nicht zugewiesen) — **Ausgang:** offen bis Closure.
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

### Belege des Implementers

**Stand.** Geändert sind `docs/user/benutzerhandbuch.md` (Version 0.4: neuer Abschnitt *Zeitangaben aufzeichnen und einspielen*, vier Zeilen der Tabelle in §5, drei Zeilen und eine neue Zeile in *Fehlercodes*, Zeile `interaction` der SQLite-Tabelle, Einleitung, Hinweis bei *Eine Aufzeichnung wiedergeben*, FAQ, Glossar), `README.md`, `docs/plan/planning/in-progress/roadmap.md` (Zeile *In Arbeit*, eigener Commit) und dieser Plan (§3, §7, §8). Produkt-Code, Tests, Gates und `docs/user/abdeckung-*.md` sind unverändert.

**Probe-Aufbau.** `make build` am Stand `15583ea` (Image `pgwire-recorder:dev`, `sha256:0a130c4a8cda…`, Produktcode unverändert). Proben in einem Scratch-Verzeichnis außerhalb des Repos. Netz `zd-net`, ein PostgreSQL 17.11 aus dem gepinnten Image von `harness/mk/integration.mk` (`zd-pg`, Servermodus `trust`, mit `log_statement=all`, `log_connections=on`, `log_min_duration_statement=0` und Zeitstempel mit Millisekunden; im Netz zusätzlich unter dem Namen `postgres`), `record`, `replay` und `play` als `docker run … pgwire-recorder:dev` im selben Netz, Client `psql` aus dem PostgreSQL-Image (Pausen des Clients mit `\! sleep`), für Pausen vor der Startup-Nachricht ein Roh-Client in Python (`SSLRequest`, Startup, Query). Gemessen wird **am Server**: `play` sendet an `zd-pg`, die Zeitstempel von „connection received“ und „statement“ im Server-Log sind Beginn der Verbindung und Beginn der Anfrage (die Zeit steht vor dem Warten auf eine Sperre); „ENDE“ ist die Zeile `duration` der Anfrage. Aufzeichnungen: `r1` (drei Einfügungen, Pausen 1 s und 1 s: `offset_ms` 4, 1012, 2016), `r2` (Sitzung 1 nach 2 s Pause: 2005, 3010; Sitzung 2: 2, 1008), `r3` (Pausen 8 s: 3, 8014, 16021), `fast` (4, 5, 6), `s1` (`r1` als SQLite). Die Blöcke des neuen Abschnitts liefen **wörtlich** (Dateien `blk0` bis `blk4`, aus dem Handbuch gezogen; `pgwire-recorder` als Funktion, die das Image mit dem Netz startet; `record` im Hintergrund mit `psql` als Anwendung, `play` im Vordergrund). Container (mit `docker rm -fv`) und Netz sind danach per Name entfernt; die Zahl der Volumes ist vor und nach den Proben 51, kein neues blieb liegen. Die Dateien liegen im Scratch-Verzeichnis der Sitzung außerhalb des Repos. Nur Linux im Container, lokales Dateisystem (`BEO-REPO/verhalten-nur-unter-linux-geprueft`).

**Beispiele des Abschnitts (wörtlich).**

| Block | Lauf | Ergebnis |
|---|---|---|
| `record … --record-timing` (YAML) | Block, dann `psql` mit zwei Einfügungen, 1 s Pause, `docker stop` | Exit 0; `sequence: 2` trägt `offset_ms: 1011` unter derselben Einrückung wie der Ausschnitt |
| `play … --keep-timing` | Block gegen die Aufzeichnung | Exit 0; Server-Log: Anfrage 1 bei 2 ms, Anfrage 2 bei 1009 ms nach der Verbindung |
| Konfigurationsdatei mit vier Schlüsseln | `record --config`, danach `play --config` | Aufzeichnung mit 2 Zeilen `offset_ms`; `play` Exit 0, Anfragen bei 2 ms und 1011 ms (`first-request`, `absolute`) |
| `sqlite3 … select session_id, sequence, offset_ms from interaction` | Aufzeichnung mit `--format sqlite --record-timing` | `1|1|3` und `1|2|1010` |

**Proben je Aussage.** `Z` = die Aufzeichnungen oben; „Server-Log“ = Zeit seit „connection received“ der Verbindung von `play`.

| Aussage im Handbuch | Probe | Ergebnis (Gegenprobe) |
|---|---|---|
| `--record-timing` hält an jeder Anfrage die Zeit seit Beginn der Verbindung fest; ganze Zahl ab 0 | `r1` (Client: Einfügung, 1 s, Einfügung, 1 s, Einfügung); SQLite `typeof` | 4, 1012, 2016; `integer` (Pausen von 1 s ergeben Abstände von 1008 und 1004) |
| Zeit bis zum Senden der Startup-Nachricht zählt nicht, auch nicht nach einer Anfrage nach Verschlüsselung | Roh-Client: 2 s vor Startup; `SSLRequest`, `N`, dann 2 s vor Startup; 2 s vor dem `SSLRequest`; 2 s nach Verbindungsaufbau vor dem `SSLRequest` | je `offset_ms` 4 |
| Zeit zwischen Anmeldung und erster Anfrage zählt mit | Roh-Client: Startup sofort, nach `ReadyForQuery` 2 s Pause | 2005 (Gegenprobe zu den vier Fällen oben) |
| Wert der ersten Anfrage nicht notwendig 0 | `psql` und Roh-Client ohne Pause | 4 und 3 |
| Jede Verbindung zählt ab eigenem Beginn | `r2`: Sitzung 1 endet bei 3010, Sitzung 2 folgt | Sitzung 2: 2 und 1008 (nicht nach 3010) |
| Anfragen beider Protokolle und Lebendprüfung tragen das Feld | `psql` mit `\bind` (erweitert) und Roh-Client mit `-- ping` | `type: extended` mit 1012; `-- ping` mit 1504; einfache Anfrage 4 |
| Ohne `--record-timing` kein Feld (YAML), `NULL` (SQLite), Standard | `record` ohne Option, YAML und SQLite | `grep -c offset_ms` 0; SQLite `typeof` `null`, Wert leer; mit Option `integer` |
| Zwei Aufzeichnungen desselben Ablaufs können sich in `offset_ms` unterscheiden | zwei Läufe mit `r1`-Ablauf, `diff` | drei Zeilen verschieden (4/3, 1013/1012, 2018/2017), alle `offset_ms`; keine andere Zeile |
| `RECORD_TIMING=true` genügt; Schlüssel `record_timing: true` genügt | `record` mit der Variablen; `record --config` | je 1 Zeile `offset_ms` (Gegenprobe: ohne beides 0) |
| Priorität Option vor Umgebung vor Datei; leere Umgebungsvariable gilt als nicht gesetzt | Umgebung `true` mit `--record-timing=false`; Datei `true` mit Option `false`; Umgebung `false` mit Datei `true`; Umgebung leer mit Datei `true` | 0, 0, 0 Zeilen; 1 Zeile (leer: Datei gilt) |
| `--record-timing` nimmt keinen Wert, `=true`, `=false`; ungültig `=1`, `=True`, `=`: `PGR-E2001`; Umgebungsvariable `1`, `Geheim123`: `PGR-E2001`, auch mit Option; Datei `yes`: `PGR-E2004` | je ein `record` | Option ohne Wert, `=true`, `=false` starten; `=1`, `=True`, `=` Exit 2 `PGR-E2001`; Variable `1`, `Geheim123` Exit 2 `PGR-E2001`, auch mit `--record-timing`; Datei `yes` Exit 2 `PGR-E2004` |
| `--keep-timing` wartet vor jeder Anfrage; ohne die Option wartet `play` nie, auch mit Zeitangaben | `play` von `r1` ohne und mit Option, Server-Log | ohne: 1, 9, 15 ms; mit: 3, 1012, 2016 ms (Abstände 1009 und 1004 bei aufgezeichneten 1008 und 1004) |
| Für `--keep-timing` braucht jede Anfrage eine Zeitangabe; mit `--record-timing` aufzeichnen | `play --keep-timing` von `s0` (ohne Option aufgezeichnet) und von `r1`/`s1` | `s0`: `PGR-E2003`, Exit 2; `r1`, `s1`: Exit 0 |
| `KEEP_TIMING=true` und `keep_timing: true` genügen; `timing_mode`, `timing_reference` als Umgebung und Datei | `play` mit Umgebung (`absolute`, `first-request`, `r2`), `play --config` mit `keep_timing: true` (`r1`) | Anfragen bei 2, 1008 / 1015, 2021 ms (`first-request`) bzw. 2, 1011, 2016 ms (Datei) |
| Zeitangaben gelten für SQLite | `play --keep-timing` von `s1` | 2, 1012, 2017 ms |
| `relative`: Standard; erste Anfrage ohne Wartezeit; sonst mindestens der aufgezeichnete Unterschied seit dem Beginn der vorigen Anfrage | `play --keep-timing` von `r1`, `r2` | `r1` siehe oben; `r2`: 3 ms, 1008 ms (aufgezeichneter Unterschied 1005), zweite Verbindung bei 1013 ms, ihre erste Anfrage 2 ms nach ihrem Beginn |
| Dauert die vorige Anfrage länger als aufgezeichnet: nächste, sobald beantwortet; Verspätung später nicht aufgeholt | Tabelle `zlog` 3 s gesperrt (zweite Sitzung des Clients), `play --keep-timing` von `r1`, Server-Log | Anfrage 1 bei 2 ms, Ende 1693 ms; Anfrage 2 bei 1694 ms (1 ms nach dem Ende); Anfrage 3 bei 2698 ms = Beginn von 2 plus 1004 |
| `absolute`: Bezugspunkt; Verspätung wird aufgeholt, Abstand kann kürzer sein | dieselbe Sperre, `--timing-mode absolute` | Anfrage 2 bei 1650 ms (1 ms nach dem Ende von 1), Anfrage 3 bei 2016 ms: Abstand 366 ms bei aufgezeichneten 1004 ms (Gegenprobe: `relative` 1004 ms) |
| `connect` (Standard): Beginn des Verbindungsaufbaus, wartet auch vor der ersten Anfrage; `first-request` und `relative` ohne Wartezeit | `r2` mit `absolute` ohne `--timing-reference`, mit `first-request`, mit `relative` | `absolute`: 2003, 3009 ms (aufgezeichnet 2005, 3010), die zweite Verbindung bei 3014 ms, ihre Anfragen bei 3016 und 4021 ms; `first-request`: 2, 1007 ms, zweite Verbindung 1012 ms, Anfragen 1013 und 2020 ms; `relative`: 3, 1008 ms |
| Jede Verbindung beginnt mit eigener Uhr; Verbindungen nacheinander | `r2` mit `absolute`: zweite Verbindung | öffnet bei 3014 ms, nachdem die letzte Anfrage der ersten bei 3009 ms begann; ihre Anfragen 2 ms und 1007 ms nach ihrem Beginn (nicht nach 3010) |
| Werte nur Kleinbuchstaben; `Absolute`, `first_request`, leer: `PGR-E2001` (Option, Umgebungsvariable), `PGR-E2004` (Datei); Meldung nennt den Wert von Variable und Datei nicht; leere Umgebungsvariable gilt als nicht gesetzt | Option `Absolute`, `=`, `""`; `first_request`, `=`; Umgebung `Geheim123` bei Modus und `KEEP_TIMING`, `first_request` bei Referenz; Datei `Absolute`, `first_request`, `""`, `keep_timing: yes`; Umgebung leer | Optionen Exit 2 `PGR-E2001`; Umgebung Exit 2 `PGR-E2001`, `Geheim123` und `first_request` kommen in keiner Meldung der Umgebung vor, `Absolute` und `yes` in keiner der Datei; Datei Exit 2 `PGR-E2004` mit Schlüssel (`play.timing_mode` usw.), ohne Wert; leere Variable: Exit 0 |
| `--keep-timing`: Wertemenge, Priorität, Umgebungsvariable wie bei `--record-timing` | `=1`, `=True`, `=`; Umgebung `1`, `True`, `Geheim123`, mit Option; Umgebung `false` mit `--keep-timing`; Umgebung `true` mit `--keep-timing=false`; Datei `keep_timing: true` und `timing_mode: absolute` mit Umgebung `false` und mit leerer Umgebung | `=1`, `=True`, `=` und die Umgebungswerte Exit 2 `PGR-E2001`, auch mit Option; Option `--keep-timing` mit Umgebung `false`: Exit 0; Umgebung `true` mit `=false` und `absolute`: `PGR-E2001`; Umgebung `false` vor Datei: `PGR-E2001`; leere Umgebung: Exit 0 |
| Ungültig: `--timing-mode`/`--timing-reference` ohne `--keep-timing`; `--timing-reference` mit `relative`; als gesetzt gilt auch der Wert des Standards; Reihenfolge ändert nichts | `play` mit den Kombinationen (Option, Umgebung, Datei) | Exit 2 `PGR-E2001` bei: Modus allein, Referenz allein (`connect`, `first-request`), `--keep-timing --timing-reference connect`, `--timing-mode relative --timing-reference first-request`, `--keep-timing=false --timing-mode absolute`, Umgebung `TIMING_MODE=absolute` ohne `--keep-timing`, Datei `timing_mode: absolute` ohne `keep_timing`, Datei mit `keep_timing: false`, Datei `first-request` mit Option `--timing-mode relative`. Exit 0 bei: `--keep-timing`, `--keep-timing --timing-mode relative`, `--keep-timing=false`, `--timing-mode absolute --keep-timing`, `--timing-reference connect --timing-mode absolute --keep-timing`, beide Referenzen mit `absolute`, Datei mit Option `--keep-timing` |
| `record` kennt die drei Optionen nicht, `play` nicht `--record-timing`, `replay` keine der vier | je ein Aufruf | Exit 2 `PGR-E2001` (`flag provided but not defined`) in allen acht Fällen |
| `PGR-E2003`: beim Start, Exit 2, vor einer Verbindung; nennt Sitzung und erste Anfrage ohne Angabe; ohne `--keep-timing` wird eingespielt | `v-miss2` (Anfrage 2 ohne); `v-r2-miss-s2` (Sitzung 2, Anfrage 2); `v-r2-miss-s1i1` (Sitzung 1 Anfrage 1 und Sitzung 2 Anfrage 2); Server-Log; ohne `--keep-timing` | `Session 1, Interaktion 2`; `Session 2, Interaktion 2`; `Session 1, Interaktion 1` (die erste in der Reihenfolge); Server-Log: 0 Verbindungen; ohne `--keep-timing` Exit 0 |
| Abbruchsignal beendet jede Wartezeit sofort; ohne Option keine weitere Anfrage; mit Option übrige ohne Wartezeit; Exit-Code 0 | `r3` (Pausen 8 s), Signal nach 3 s: `SIGTERM` ohne und mit `--finish-session-on-interrupt`; `SIGINT` bei erster Wartezeit von `absolute connect` (20 s), bei Wartezeit nach der ersten mit `first-request`, mit Option | ohne: Ende nach 0,3 s, nur Anfrage 1 am Server, Exit 0; mit: Anfragen 2 und 3 bei 3063 und 3066 ms, Exit 0; `connect`: Ende ohne Anfrage, Exit 0; `first-request`: nur Anfrage 1, Exit 0; `connect` mit Option: alle drei bei 2557 bis 2564 ms |
| `replay` wartet nie auf Zeitangaben | `replay` von `r1` (Spanne 2 s), Client ohne Pause | drei Einfügungen in 0,57 s einschließlich des Starts des Clients, Exit 0 |
| Prüfung beim Laden für `replay` und `play`, mit und ohne `--keep-timing`, YAML und SQLite | `v-desc`, `q-desc`, `q-neg` mit beiden Befehlen, ohne und mit Option | je Exit 3 `PGR-E3003` (SQLite mit `--keep-timing` für `play`) |
| YAML: gültig `0` und Ziffern ohne führende Null bis 9223372036854775807; ungültig `1.5`, `5.0`, `1e3`, `0x10`, `007`, `+5`, `-1`, `"5"`, `true`, `null`, 9223372036854775808 | `r1` mit einem Wert an Stelle 2 bzw. 3; `play` ohne Option, `replay` für vier | Gültig (Exit 0): `5`, `0` an Stelle 1, 9223372036854775807 an Stelle 3; ungültig: je Exit 3 `PGR-E3003`, die letzte „außerhalb des Bereichs“, `replay` für `1.5`, `-1`, `null` und den Wert über dem Bereich ebenso |
| SQLite: `NULL` oder ganze Zahl ab 0; negative Zahl, Text, Dezimalzahl beschädigt | `s1` mit `-1`, `'abc'`, `1.5`, `NULL` | `-1`, `abc`, `1.5`: Exit 3 `PGR-E3003` (`play`, `replay`); `NULL`: Exit 0, mit `--keep-timing` `PGR-E2003` |
| Absteigend (kleiner als der vorige derselben Sitzung) ist beschädigt; gleiche Werte gültig; kleinerer Wert in der nächsten Sitzung gültig | `v-desc` (2016 → 5), `v-gleich` (1012, 1012), `q-gleich`, `r2` (Sitzung 2 beginnt bei 2 nach 3010) | `PGR-E3003` mit `Session 1, Interaktion 3`; `v-gleich`, `q-gleich`: Exit 0 (`play`, auch mit `--keep-timing`), `replay` läuft; `r2` Exit 0 |
| Anfrage ohne `offset_ms` zwischen zwei mit Wert ist gültig; `play --keep-timing` lehnt sie ab | `v-miss2` mit `play`, `replay`, `play --keep-timing` | Exit 0; `replay` läuft; `PGR-E2003` |
| Meldung nennt Nummer von Sitzung und Anfrage, nie den Wert | `v-desc-wert` (31337 → 271), `v-bad-wert` (4242.5), `q-wert` | `Session 1, Interaktion 3` bzw. `2`; weder `31337` noch `271` noch `4242.5` in der Meldung |
| Beschädigung geht `PGR-E2003` vor, auch in einer anderen Sitzung; ohne Beschädigung `PGR-E2003` | `v-miss1-desc3` (Anfrage 1 ohne, Anfrage 3 absteigend), `v-r2-miss-s1-desc-s2` (Sitzung 1 Anfrage 1 ohne, Sitzung 2 absteigend), je mit `--keep-timing` | Exit 3 `PGR-E3003` (`Session 1, Interaktion 3` bzw. `Session 2, Interaktion 2`); Gegenprobe `v-miss2`: Exit 2 `PGR-E2003` |
| Tabelle §5: Betriebsarten, Variablen, Standardwerte der vier Optionen | `--help` von `record` und `play` | `--record-timing` in `record`; `--keep-timing`, `--timing-mode`, `--timing-reference` in `play`; Standards `false`, `false`, `relative`, `connect` (Verhalten ohne die Option: siehe Zeilen oben) |
| Zeile `PGR-E2001`, `PGR-E2003`, `PGR-E2004`, `PGR-E3003` in *Fehlercodes* | dieselben Proben | Codes, Exit-Codes 2, 2, 2, 3 wie in den Zeilen oben; die Schlüssel `record_timing`, `keep_timing`, `timing_mode`, `timing_reference` je `PGR-E2004` bei ungültigem Wert |
| FAQ: bei jedem Aufzeichnen mit `--record-timing` können sich die Werte ändern | Zeile *Zwei Aufzeichnungen …* | siehe dort |
| SQLite-Tabelle: `offset_ms` ohne `--record-timing` ohne Wert, also `NULL` | `s0` | `null`, Wert leer (Zeile oben) |
| Hinweis bei *wiedergeben*: `replay` wartet nie, prüft beim Laden | Zeilen *replay wartet nie* und *Prüfung beim Laden* | siehe dort |
| keine rechnerspezifischen Angaben (FAQ, unverändert; liest sich mit den neuen Werten) | `grep -a` nach Namen des Containers, des Netzes und `/rec` in `da.yaml` und `s1.sqlite` | 0 Treffer in beiden |

**Nicht in das Handbuch übernommen** (Probe zeigt es nicht oder die Grenze steht nur in der Spezifikation): dass der Wert der ersten Anfrage die Zeit für den Aufbau zur Datenbank einschließt (die Probe misst nur 3 und 4 ms; ein langsamer Aufbau ist nicht erzeugt), die Reihenfolge der Werte einer Verbindung beim Aufzeichnen (alle Proben stiegen auf; zugesagt ist es für Anfragen, die nacheinander eintreffen), die Wartezeit bei großen Werten, `--compare-responses`, Windows und macOS, Netzwerk-Dateisysteme, der Wortlaut der Meldungen, eine Sitzung ohne Anfrage (das Binary lehnt sie in YAML als `PGR-E3003` ab, ohne mit Zeitangaben zu tun zu haben).

**Befehlsfolge aus `slice-doku-ist-stand` §7** (mit der benannten Ausnahme `PGWIRE_RECORDER_PASSWORD`, Variablen-Zeile `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`), am Stand dieses Commits:

| Prüfung | Ergebnis |
|---|---|
| Optionen des Handbuchs gegen `--help` | nur `--h`, `--help`, `--version` (Nicht-Optionen im Handbuch), `--rm`, `--name` (`docker run`) |
| Variablen ohne Option | leer |
| Codes außerhalb des Katalogs | leer |
| Mutant: `--format-bogus` angehängt | `--format-bogus` in der ersten Ausgabe |
| Mutant: `PGWIRE_RECORDER_BOGUS` angehängt | in der zweiten Ausgabe |
| Mutant: `PGR-E9999` angehängt | in der dritten Ausgabe |
| Mutant: `--keep-timing` aus der Hilfe entfernt | `--keep-timing` in der ersten Ausgabe (die Option steht im Handbuch) |

**`grep` über Handbuch und README.** Im Handbuch: kein Treffer für `ADR-`, `SPEC-`, `LH-`, `[Ss]lice`, `Welle`, `Spezifikation`, `Lastenheft`, `[Rr]eview`; in den neuen Zeilen von Handbuch und README kein „noch nicht“, „kommt“, „geplant“, „bisher“, „später“, `compare-responses`. Im README bleiben die Verweise auf `spec/` und `docs/plan/` ohne Aussage über einen Stand.

**Namensabgleich.** `docs/user/benutzerhandbuch.md`, `README.md`, `internal/hexagon/model/fehler.go`, `harness/mk/integration.mk`, `done/slice-doku-ist-stand.md` und die vier Slices `slice-v1-abschluss-zeitangaben-pruefen`, `-uhr`, `-aufzeichnen` und `slice-v1-abschluss-zeitangaben` in `done/`: per `ls` und `grep` gefunden; die Anker `#zeitangaben-aufzeichnen-und-einspielen` (Handbuch, mehrfach), `#mit-einem-datenbanktreiber-arbeiten` löst `make docs-check` auf.

**Nacharbeit zum Review (F-660 bis F-663), Stand `445f074`.** Geändert sind Handbuch und README (Sätze, keine neue Aussage außer dem Log-Satz in *Ergebnis*) und dieser Plan. Aufbau wie oben (Image `pgwire-recorder:dev` unverändert, neues Netz, neuer PostgreSQL mit Server-Log, Container mit `docker rm -fv`, Netz per Name entfernt).

| Aussage im Handbuch (neu oder geändert) | Probe | Ergebnis (Gegenprobe) |
|---|---|---|
| F-660: Einleitung, Abschnittsanfang, *Ergebnis*, README: `play` wartet „vor den Anfragen“; beim Standard `relative` ist die erste Anfrage jeder Verbindung ohne Wartezeit | `play` von `r2` (erste Anfrage aufgezeichnet bei 2005 ms) ohne Option; `--keep-timing`; `--timing-mode absolute`; `absolute` mit `first-request`; Server-Log | ohne Option 3, 11, 20, 24 ms; `relative`: erste Anfrage 3 ms (nicht 2005), zweite 1008 ms, zweite Verbindung 1013 ms, ihre erste Anfrage 1015 ms (2 ms nach Beginn); `absolute` mit `connect`: erste Anfrage bei 2004 ms (Gegenprobe: die erste Anfrage wartet, wo der Bezugspunkt `connect` ist); `first-request`: 2 ms |
| F-661: Abbruchsignal in der Wartezeit: Exit-Code 0, nach einem früheren Fehler 4 | Aufzeichnung `err20` (Anfrage 1 `select * from nirgends`, Antwort der Datenbank mit Fehler, Anfragen 2 und 3 bei 20000 ms): `SIGTERM` nach 3 s mit `--continue-on-error`; `SIGINT` mit `--continue-on-error --finish-session-on-interrupt` | Exit 4 mit `PGR-E4004` in beiden Fällen (Gegenprobe ohne früheren Fehler, `r3` mit `SIGTERM`, ohne und mit der Option: Exit 0; ohne `--continue-on-error` endet `err20` schon an der Fehlerantwort mit `PGR-E4004`, Exit 4) |
| F-662: `offset_ms` zählt seit der Startup-Nachricht der Verbindung (Einleitung, Abschnittsanfang, Hinweis, Glossar, README) | Roh-Client: 2 s Pause vor der Startup-Nachricht; `SSLRequest`, `N`, dann 2 s Pause; 2 s vor dem `SSLRequest`; 2 s Pause nach `ReadyForQuery` | 2 ms, 4 ms, 4 ms; 2004 ms (Gegenprobe: die Zeit nach der Anmeldung zählt) |
| F-663: *Ergebnis* zum Einspielen mit beobachtbarer Ausgabe: erste und letzte Anfrage eine Sekunde auseinander, `relative`: zwischen `play gestartet` und `play beendet` mindestens eine Sekunde, ohne `--keep-timing` weniger | Aufzeichnung `e1` (`offset_ms` 4 und 1011), `play` ohne und mit Option | ohne: 16 ms zwischen den beiden Zeilen; mit: 1,014 s (Gegenprobe: ohne Option weniger als eine Sekunde) |
| F-663: Gliederung, keine neue Aussage | Kombinationen als eigener Abschnitt *Ungültige Kombinationen* mit fünf Punkten, *Fehlende Zeitangaben und Abbruch* als eigener Abschnitt | Wortlaut der Sätze unverändert; Proben der Kombinationen, `PGR-E2003` und des Abbruchsignals stehen in der Tabelle oben |

**Hinweise an Architect und Planner** (keine Randform dieses Slice, kein Befund gegen die Spezifikation):

1. *Sitzung ohne Interaktion in YAML:* `interactions: []` und ein fehlender Schlüssel enden mit `PGR-E3003` („Session 2 ohne Interaktion“). Die Spezifikation von `LH-FA-21.a` sagt, eine Sitzung ohne Interaktion verlange bei `PGR-E2003` nichts; im Binary ist der Fall schon beim Laden beschädigt. Das Handbuch sagt dazu nichts.
2. *Umgebungsvariable einer anderen Betriebsart:* `PGWIRE_RECORDER_RECORD_TIMING=true` bei `play` bleibt ohne Wirkung und ohne Meldung. Das Handbuch sagt dazu nichts.

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
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (3×, nicht zugewiesen; der Lese-Schritt der Closure von welle-v1-abschluss weist zu) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code. Nachgezählt beim Eintragen der Lesarten aus `slice-v1-abschluss-zeitangaben` (Planner, 2026-10-12, Closure jenes Slice): Liefer-Punkte 2 unverändert (die Lesarten gehören zu DoD-Punkt 1), Schichten 1. Nachgezählt beim Umsetzen (Implementer, 2026-10-12): Der Diff ändert `docs/user/benutzerhandbuch.md`, `README.md`, die Roadmap-Zeile *In Arbeit* und diesen Plan, keinen Code; Liefer-Punkte 2, Schichten 1 (Nutzerdokumentation; Planung zählt nicht).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
