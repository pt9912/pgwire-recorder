# Slice slice-v1-abschluss-anmeldung-doku: Handbuch und README zu: Anmeldung des Clients im Record-Modus vermitteln

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler)

**Berührte Spec-Stellen:** `LH-FA-05.b` (Punkt *Aufzeichnung*: der Satz zum Passwort als Startup-Parameter ist ersetzt durch „Startup-Parameter stehen unverändert in der Aufzeichnung, der Recorder reicht sie unverändert durch“; geändert vom Architect nach Review F-2 dieses Slice, `6a44529` und `0aadcac`); sonst beschreibt der Slice das gelieferte Verhalten von `slice-v1-abschluss-anmeldung`

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-anmeldung` liefert (im Handbuch der Abschnitt *Eine Anwendung mit Passwort-Anmeldung aufzeichnen*, dazu die ersetzten Grenz-Sätze), im Ist-Zustand des gebauten Binaries; Sätze über die Datenbank oder den Treiber nennen die Version, gegen die die Probe lief (PostgreSQL 17, libpq 17): ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-anmeldung`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-anmeldung` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan; die eine Änderung an `LH-FA-05.b` (siehe Kopf) hat der Architect gemacht, nicht der Implementer. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-anmeldung` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Die Anmeldung betreffende Randformen, die das Handbuch bewusst nicht nennt (Review F-7): eine nicht lesbare Nachricht des Servers im Aufbau und eine Nachricht des Clients in der Anmeldung außer `p` und `Terminate` (beide `PGR-E6001`, `LH-FA-05.b`) — Lücke ohne Probe am Binary; sie bekommt eine, wenn ein Folge-Lauf sie belegt, bis dahin sagt das Handbuch nichts dazu.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [x] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-anmeldung` liefert, im Ist-Zustand des gebauten Binaries: Aus `slice-doku-ist-stand` (dort §6, *Teilweise geliefert*): Der Slice ersetzt im Handbuch die Grenz-Sätze zum Passwort bei `record` an den neun Stellen aus §3 (§1 *Voraussetzungen*, §4 zwei Hinweise und die Zeile zu den Zeitangaben, §6 *Rollen und Rechte*, §7 die Zeilen `PGR-E4002`, `PGR-E4003`, `PGR-E6001` und *Die Anwendung kann sich nicht verbinden*) — und im README den Satz, dass die Datenbank den Benutzer ohne Passwort anmelden muss, soweit er `record` betrifft. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [x] `README.md` nennt, was `slice-v1-abschluss-anmeldung` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [x] `make gates` grün.
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
| `docs/user/benutzerhandbuch.md` | update | neuer Abschnitt *Eine Anwendung mit Passwort-Anmeldung aufzeichnen* (Beispiele als Datei), ersetzte Grenz-Sätze in §1, §4 (zwei Hinweise, die Zeile zu den Zeitangaben), §6, §7 (`PGR-E6001`, `PGR-E4002`, `PGR-E4003`, *Die Anwendung kann sich nicht verbinden*), Versionszeile 0.8 |
| `README.md` | update | Sätze, die `slice-v1-abschluss-anmeldung` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-anmeldung` liegt in `done/` (Voraussetzung; Schritt 30 der Reihenfolge in §5 von [welle-v1-abschluss](welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Handbuch und README getrennt. **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/user README.md`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 300** geht der Slice zurück (Empfehlung des Architects, 2026-10-14). **Meldeschwelle** neben der Marke: Erreicht die Summe **220**, meldet der Implementer den Stand dem Architect, bevor er weiterarbeitet. **Schätzung** (Architect, 2026-10-14, in dieser Metrik, Prosa ohne Faktor für Tests): Handbuch rund 100 bis 160 Zeilen (fünf Stellen ersetzt, ein Abschnitt zur Anmeldung neu), README rund 10; mit der Review-Runde **120 bis 200**.
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

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect hat geprüft (2026-10-14): `slice-v1-abschluss-anmeldung` liefert weitere Randformen, entschieden in `LH-FA-05.b` *Randformen der Vermittlung (Record)* (Spezifikation); das Handbuch nennt davon nur, was das Binary zeigt: die drei Verfahren, dass `record` keine Option und keine Variable für die Anmeldung kennt und `PGWIRE_RECORDER_PASSWORD` sowie Benutzer, Passwort und Datenbank der Verbindung dort unbeachtet bleiben, dass die Fehlerantwort des Servers bei fehlgeschlagener Anmeldung unverändert beim Client ankommt, dass ein Verfahren außer Klartext, MD5 und SCRAM-SHA-256 (Kerberos, GSSAPI, SSPI) mit `PGR-E6001` endet und eine verletzte Reihenfolge des Servers mit `PGR-E4002`, dass ein Client, der in der Anmeldung abbricht, keine Meldung und keine Aufzeichnung hinterlässt, dass Passwort und Anmeldenachrichten nicht in der Aufzeichnung stehen (die Startup-Parameter stehen unverändert darin, ohne Aussage über Zugangsdaten darunter), und dass Channel Binding nicht gelingt. Jede dieser Aussagen gilt erst mit Probe gegen das Binary (`AGENTS.md` §3.11).

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** *weiter offen* → Register (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, Vermerk in `state.md`): Beide Beispielblöcke und alle Proben liefen unter Linux im Container (Implementer und Verifier); das Handbuch sagt zu Windows und macOS nichts, also kein Auftreten und kein Entfallen.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** *entfallen*: Er trat nicht ein. Die Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-anmeldung` (neun Stellen, Befehlsfolge, Beispiele als Datei); die Lücke zur nicht lesbaren Nachricht (F-7, V-1) ist eine benannte Lücke in §1 und kein nicht gelieferter Teil, und der zu weite Satz ist behoben (`67cc189`).
- Handbuchsätze über das Verhalten des Servers (Verfahren, Fehlerantwort, Meldung) nennen keine Version, belegt ist allein die gepinnte PostgreSQL-Version 17, die Matrix weiterer Versionen liefert `slice-v1-abschluss-postgres-versionen` (Empfehlung des Planners aus der Closure von `slice-v1-abschluss-anmeldung`, nach `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `slice-v1-abschluss-tls-client-doku` V-1; keine Entscheidung des Nutzers) — **Ausgang:** *weiter offen* → Register (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, Vermerk in `state.md`; Adresse der Beobachtung bleibt `slice-v1-abschluss-postgres-versionen`). Das Handbuch nennt „PostgreSQL 17“ an allen Stellen zur Anmeldung (Review F-3 behoben, `1291bcd`); die Matrix weiterer Versionen ist nicht gefahren. Ein Eintrag in den Nehmer unterbleibt, weil er keine Adresse ist, die annehmen kann: `slice-v1-abschluss-postgres-versionen` zählt vier Schichten (§8 dort, über der Grenze von zwei), und eine Sendung zählte nach `AGENTS.md` §3.13 über der Grenze; ob er geschnitten wird, entscheidet der Planner vor dem ersten Code-Commit dort. Die Sätze „geprüft mit PostgreSQL 17“ im Handbuch bleiben richtig, solange nur 17 gefahren ist; liefert der Nehmer die Matrix, gehören sie in dessen Handbuch-Teil (Frage an den Nutzer in der Closure-Mitteilung, kein Eintrag).

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

Kein Produkt-Code, kein Test, kein neuer Wächter: Der Slice ändert `docs/user/benutzerhandbuch.md` und `README.md` (Commit `e12ab97`) und diesen Plan; eine Mutation des Produkts entfällt (`AGENTS.md` §3.10 gilt für neue Verträge). Statt ihrer steht je Zeile unten das Gegenbeispiel, das der Satz ausschließt. Der Bericht des Implementers verweist hierher.

**Aufbau der Proben** (alles unter Linux im Container, `BEO-REPO/verhalten-nur-unter-linux-geprueft`; ein Lauf zur Zeit, ein Recorder zur Zeit, kein Netz-Alias, kein `pkill`; danach entfernt: alle Container `impl-anmd-*` und das Netz `impl-anmd-net`, `docker ps -a` und `docker network ls` ohne Treffer). Binary: `pgwire-recorder:dev` (`make build` auf `f149fda`, Stufe `runtime`). Server: `postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24` (PostgreSQL 17.11, der Digest aus `harness/mk/integration.mk`) im eigenen Docker-Netz. Client: `psql` 17.11 (libpq 17.11) aus demselben Image. Rohnachrichten und Stand-in-Server: `python:3.13-alpine` (Quellen unten). Zertifikat: OpenSSL 3.3.7 aus `alpine:3.20` (`apk add openssl`, mit Netz), EC P-256, selbst signiert, CN und SAN `impl-anmd-rec`.

**Genannte Versionen und der Lauf, der sie belegt** (`AGENTS.md`-Lehre aus `slice-v1-abschluss-tls-client-doku` V-1/V-2; ein Verifier fährt sie mit dem Aufbau und den Befehlen unten nach):

| Handbuch / README nennt | belegt durch | Nachfahren |
|---|---|---|
| PostgreSQL 17 (Handbuch *Voraussetzung* des neuen Abschnitts, §1 *Voraussetzungen*, §4 Hinweis, §6; README) | Server `postgres:17-alpine` (17.11), Proben P1 bis P5, P7, P8, P10, P11 | Aufbau, dann P1 |
| libpq 17 / `psql` aus libpq 17 (Handbuch *Voraussetzung*, Hinweis zu Channel Binding) | `psql --version` im Server-Image: `psql (PostgreSQL) 17.11`; Proben P1, P4, P7 | `docker run --rm "$PGI" psql --version`, dann P7 |
| keine weitere Version, kein weiterer Treiber, kein Betriebssystem | nicht gefahren (Handbuch sagt „nicht geprüft“, mehr nicht) | — |

**Aufbau und Hilfsfunktionen** (`S` ist ein Scratch-Verzeichnis außerhalb des Repos, `S/aw` für 777 beschreibbar, `S/pki` das Zertifikat, `S/py` die Quellen unten):

```bash
PGI='postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24'
docker network create impl-anmd-net
docker run -d --name impl-anmd-pg --network impl-anmd-net -e POSTGRES_HOST_AUTH_METHOD=trust "$PGI"
# pg_hba.conf (vor der Zeile fuer alle) und Rollen, wie tools/test/run-integration-tests.sh:
#   host all pw_scram all scram-sha-256 | host all pw_md5 all md5 | host all pw_klar all password | host all pw_reject all reject
#   pw_scram GEHEIMscram (scram), pw_md5 GEHEIMmd5 (md5), pw_klar GEHEIMklar, pw_reject, app geheim (scram, Datenbank app), Datenbank anmd; postgres: trust
rec_start() { f=$1; shift; docker rm -f impl-anmd-rec >/dev/null 2>&1
  docker run -d --name impl-anmd-rec --network impl-anmd-net -v $S/aw:/w -v $S/pki:/pki pgwire-recorder:dev record \
    --listen 0.0.0.0:15432 --upstream ${UPSTREAM:-impl-anmd-pg:5432} --output /w/$f "$@"; sleep 1; }
rec_stop()  { docker stop impl-anmd-rec; echo exit=$(docker wait impl-anmd-rec); docker logs impl-anmd-rec; docker rm impl-anmd-rec; }
psqlc()     { u=$1; p=$2; x=$3; shift 3; docker run --rm --network impl-anmd-net -e PGPASSWORD="$p" "$PGI" \
    psql "host=impl-anmd-rec port=15432 user=$u dbname=anmd sslmode=disable $x" -At -c "$1"; }
fake_start(){ docker run -d --name impl-anmd-fake --network impl-anmd-net -v $S/py:/py -e MODE=$1 python:3.13-alpine python /py/fake_server.py; sleep 1; }
pyc()       { docker run --rm --network impl-anmd-net -v $S/py:/py python:3.13-alpine python /py/client.py "$@"; }
```

**Proben** (jede mit Aufruf in Kurzform):

| Nr. | Aufruf | Ergebnis |
|---|---|---|
| P1 | `rec_start p1.yaml`; `psqlc` als `pw_scram`, `pw_md5`, `pw_klar` (je mit dem Passwort) und `postgres` (ohne); `rec_stop`; `grep -c -E 'GEHEIM\|md5[0-9a-f]{32}\|SCRAM-SHA\|n,,n=\|biws' p1.yaml` | alle vier Anmeldungen liefern ihre Zeile, 4 Sitzungen mit `sql: select current_user`, Exit 0, **0 Treffer** |
| P2 | wie P1 mit `--log-level debug` (ohne `postgres`), `docker logs` gegen dasselbe Muster; Kontrolle: `printf 'GEHEIMx\nmd5%032d\nSCRAM-SHA-256\nn,,n=a\nbiws\n' 0 \| grep -c -E '<Muster>'` | Log 0 Treffer, Aufzeichnung 0 Treffer, **Kontrolle 5 Treffer** (das Muster fängt jede der fünf Spuren). Das Log enthält bei diesem Ablauf nur die Zeile „record gestartet“ vor dem Beenden; die Zusage „auch nicht im Debug-Log“ trägt zusätzlich der Test des Code-Slice (`TestE2ERecordAnmeldung`, Spuren auch in `rec.stderr`, Mutant W18real, `done/slice-v1-abschluss-anmeldung` §7) |
| P3 | `rec_start r.yaml`, eine Anfrage als `pw_scram`, `rec_stop`; `replay --input` auf diese Datei, `psql` ohne Passwort als `wer`/`egal`; zweiter Lauf mit frischem `replay` und `PGPASSWORD=anders` | erster Lauf: ohne Passwort Zeile `1`, die zweite Verbindung (`anders`) kam durch die Anmeldung und endete mit `PGR-E5003` (nur eine Sitzung aufgezeichnet); zweiter Lauf, frisches `replay`: mit `anders` Zeile `1` |
| P4 | `rec_start p4.yaml --record-empty-sessions`; `psqlc` als `pw_scram`, `pw_md5`, `pw_klar` mit `falsch`, als `pw_reject`, als `nobody`; `rec_stop`; dieselben Aufrufe direkt am Server (`host=impl-anmd-pg port=5432`) | je `FATAL:  password authentication failed for user "<name>"`, `pg_hba.conf rejects connection …`, `role "nobody" does not exist`, **Text gleich dem direkten Lauf**; Log nur „gestartet“, „Herunterfahren“, „beendet“; Exit 0; Datei `empty_sessions: true`, `sessions: []` |
| P4b | wie P4 ohne `--record-empty-sessions`, ein Versuch mit `falsch` | gleicher Text, Exit 0, Log wie oben, Datei `sessions: []` (das Werkzeug schreibt die Datei beim Beenden; das Handbuch sagt dazu nichts) |
| P5 | Datei `connections: lokal: "postgresql://anderer:${FALSCH}@impl-anmd-pg:5432/andere"`, `record: upstream: lokal`; `-e PGWIRE_RECORDER_PASSWORD=falsch -e FALSCH=auchfalsch`; `psqlc` als `pw_scram` (Datenbank `anmd`) und `pw_md5` | beide Anmeldungen gelingen; Aufzeichnung trägt `user: pw_scram`, `database: anmd`, nicht `anderer`/`andere`; Exit 0 |
| P6 | `fake_start code:N` mit N = 2, 6, 7, 9, 42, `UPSTREAM=impl-anmd-fake:5432 rec_start`, `psqlc pw_klar …`, `rec_stop`; ohne und mit `--record-empty-sessions` (N = 7) | je `FATAL:  nicht unterstützt [PGR-E6001]: der Upstream verlangt ein Anmeldeverfahren (Code N), das record nicht vermittelt`, Log-Zeile `code=PGR-E6001`, **Exit 6**, `sessions: []` |
| P6b | `fake_start order11` (Code 11 als erste Nachricht), `fake_start ok-then-r` (`AuthenticationOk`, dann Code 3), `fake_start close3` (Code 3, dann nach der Antwort Verbindung schließen); mit `--record-empty-sessions` | `PGR-E4002` („Nachricht der Anmeldung mit Code 11“ / „Code 3 außerhalb der Reihenfolge“), `PGR-E4002`, `PGR-E4003` („Verbindung zum Upstream im Verbindungsaufbau beendet“); Client erhält die Meldung als `FATAL`, Log-Zeile mit dem Code, **Exit 4**, `sessions: []` |
| P7 | a: `rec_start`; `psqlc pw_scram … "channel_binding=require"` und `channel_binding=prefer`. b: `rec_start --tls-cert /pki/server.pem --tls-key /pki/server.key`; `psql` mit `sslmode=require` und `channel_binding=require` / `prefer` (als `pw_scram`), dazu `pw_md5` mit `sslmode=require` | a: `require` → `channel binding required, but SSL not in use`; `prefer` → `1`. b: `require` → `channel binding is required, but server did not offer an authentication method that supports channel binding`; `prefer` → `1`; `pw_md5` → `1`; 2 Sitzungen aufgezeichnet |
| P8 | `rec_start p8.yaml --record-empty-sessions`; `pyc startup-close impl-anmd-rec`, `pyc terminate …`, `pyc reset …` (Startnachricht, `R` Code 3 lesen, dann schließen / `X` / Reset mit `SO_LINGER` 0); danach `psqlc pw_scram … "select 'danach'"` | je „erhalten R 3“, Log ohne Zeile, **Exit 0**, Datei: nur `- id: 1` mit `sql: select 'danach'` |
| P10 | `rec_start --record-timing`; `pyc slow-klar` (Startnachricht, `R` Code 3, **1,5 s** warten, `p`, `Q`), und mit `DELAY=0` | `offset_ms: 1505` und `offset_ms: 4` (Gegenrichtung: ohne Warten klein) |
| P11 | die Blöcke des neuen Abschnitts per `awk` aus dem Handbuch gezogen: Block 1 mit `pgwire-recorder` als Funktion (`docker run -d … pgwire-recorder:dev "$@"`), `postgres:5432` → `impl-anmd-pg:5432`, `./recordings/` → `/recordings/`; Block 2 mit `host=localhost` → `host=impl-anmd-rec`, `psql` als Funktion mit `PGPASSWORD=geheim` (der Wert des Blocks) | `?column?` / `1`, Exit 0; Aufzeichnung `login.yaml` trägt `user: app`, `database: app`, `sql: select 1`, **0 Treffer** für `geheim\|SCRAM-SHA\|n,,n=\|biws` |
| P12 | `pyc klar-falsch impl-anmd-rec` (Startnachricht, `R` Code 3, `p` mit `falsch`, die Fehlerantwort als Hex) und `PORT=5432 … klar-falsch impl-anmd-pg` | die Fehlerantwort (Typ `E`, alle Felder) ist **byte-gleich** (`cmp`) |

**Stand-in-Server und Rohclient** (die Quellen der Proben; gekürzt auf das Nötige):

```python
# fake_server.py (MODE aus der Umgebung; Port 5432)
import socket, struct, os
mode = os.environ["MODE"]
s = socket.socket(); s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1); s.bind(("0.0.0.0", 5432)); s.listen(5)
def rd(c, n):
    b = b""
    while len(b) < n:
        d = c.recv(n - len(b))
        if not d: raise EOFError
        b += d
    return b
def auth(c, code, rest=b""): c.sendall(b"R" + struct.pack("!ii", 8 + len(rest), code) + rest)
while True:
    c, _ = s.accept()
    try:
        n = struct.unpack("!i", rd(c, 4))[0]; body = rd(c, n - 4)
        if struct.unpack("!i", body[:4])[0] == 80877103:        # SSLRequest -> N
            c.sendall(b"N"); n = struct.unpack("!i", rd(c, 4))[0]; body = rd(c, n - 4)
        if mode.startswith("code:"): auth(c, int(mode[5:])); c.recv(100)
        elif mode == "order11": auth(c, 11, b"r=x,s=y,i=1"); c.recv(100)
        elif mode == "ok-then-r": auth(c, 0); auth(c, 3); c.recv(100)
        elif mode == "close3":
            auth(c, 3); rd(c, 1); n = struct.unpack("!i", rd(c, 4))[0]; rd(c, n - 4); c.close()
    except (EOFError, OSError): pass
    finally: c.close()
```

```python
# client.py MODE HOST (Port 15432, PORT aus der Umgebung); startup(), msg(), send() sind die Nachrichten des Protokolls 3.0
#   startup-close / terminate / reset: Startnachricht {user: pw_klar, database: anmd}; erste Nachricht lesen (R, Code 3); schliessen / b"X" senden / SO_LINGER (1, 0) und schliessen
#   slow-klar: wie oben, dann time.sleep(DELAY, Standard 1.5), send(b"p", b"GEHEIMklar\0"), bis "Z" lesen, send(b"Q", b"select 1\0"), bis "Z" lesen, b"X"
#   klar-falsch: Startnachricht, R lesen, send(b"p", b"falsch\0"), die naechste Nachricht (E) als Hex ausgeben
```

**Probe, die nicht trug:** Ein Passwort als Startup-Parameter (`password=…` in der Startnachricht, `startup-param`): PostgreSQL 17 lehnt den Parameter ab (`unrecognized configuration parameter "password"`), die Verbindung wird keine Sitzung. Die Spezifikation (`LH-FA-05.b` *Aufzeichnung*, vom Architect auf das Verhalten des Recorders gekürzt: Startup-Parameter stehen unverändert in der Aufzeichnung) macht dazu keine Aussage über den Server; das Handbuch sagt dazu nichts.

**Aussage im Handbuch · Probe · Ergebnis · Gegenbeispiel** (jeder Satz mit Verhaltensaussage hat eine Zeile; Abschnitt *Eine Anwendung mit Passwort-Anmeldung aufzeichnen*, sofern nicht anders genannt):

| Aussage | Probe | Ergebnis | Gegenbeispiel (was der Satz ausschließt) |
|---|---|---|---|
| `record` reicht die Anmeldung durch; Datenbank meldet mit Klartext-Passwort, MD5, SCRAM-SHA-256 oder ohne Passwort an; das Verfahren bestimmt die Datenbank (*Voraussetzung*, *Hinweise*) | P1 | vier Benutzer, vier Verfahren, derselbe Aufruf von `record`, vier Zeilen | falsches Passwort desselben Benutzers: Fehler (P4) |
| Die Anfrage läuft, die Datei enthält sie mit Benutzer und Datenbank (*Ergebnis*) | P1, P5, P11 | `sql:`, `user:`, `database:` in der Aufzeichnung | — |
| Passwort und Nachrichten der Anmeldung stehen weder in der Datei noch im Log, auch mit `--log-level debug` (*Ergebnis*) | P2, P11, P1 | 0 Treffer in Datei und Log | Kontrolle: das Muster trifft die fünf Spuren (P2) |
| `replay` gibt die Datei wieder und nimmt dabei jedes Passwort an (*Ergebnis*; §6 *wie bisher*) | P3 | ohne Passwort, mit `anders`, anderer Benutzer und Datenbank: Zeile `1` | dieselbe Anmeldung direkt an der Datenbank scheitert (P4) |
| `record` hat keine Option für Passwort; liest weder Benutzer, Passwort, Datenbank der Verbindung aus `--upstream` noch `PGWIRE_RECORDER_PASSWORD`; sie gelten nur bei `play` (*Hinweise*; §5, unverändert) | P5; Befehlsfolge | Verbindung mit anderem Benutzer/Datenbank und gesetzte Variablen ändern nichts; `--help` von `record` nennt keine Option dazu | bei `play` wirkt die Variable (Bestand, `done/slice-doku-ist-stand` §7: `PGR-E4005`) |
| Lehnt die Datenbank ab, erhält die Anwendung deren Fehlerantwort unverändert (*Hinweise*) | P4, P12 | Text gleich dem direkten Lauf (drei Verfahren, `reject`, unbekannter Benutzer); Fehlerantwort byte-gleich (Klartext) | — |
| … das Werkzeug schreibt dazu keine Meldung, der Exit-Code ändert sich nicht (*Hinweise*) | P4, P4b | Log ohne Zeile, Exit 0 | `PGR-E6001` bei anderem Verfahren: Meldung, Exit 6 (P6) |
| Beendet die Anwendung die Verbindung in der Anmeldung (Schließen, Rücksetzen, `Terminate`), gilt dasselbe, ohne Fehlerantwort (*Hinweise*) | P8 | drei Wege, Log ohne Zeile, Exit 0 | — |
| Beides ist keine Sitzung, auch nicht mit `--record-empty-sessions` (*Hinweise*; §4 *Hinweise zum Aufzeichnen*, „Ebenso keine Sitzung …“) | P4, P8, P6b | `sessions: []` bzw. nur die nachfolgende Sitzung mit `id: 1` | eine gelungene Anmeldung ist Sitzung (P1, P8 `danach`) |
| Channel Binding gelingt nicht: `channel_binding=require` bricht ab, mit und ohne `--tls-cert`; mit `prefer` gelingt die Anmeldung auch mit TLS (*Hinweise*) | P7 | `require` beide Male Fehler, `prefer` beide Male `1` | beide Richtungen in einer Probe |
| Trifft das Werkzeug auf die Anforderung eines anderen Verfahrens, `PGR-E6001`, die Anwendung erhält die Meldung, nicht aufgezeichnet, Exit 6 (*Hinweise*; §7 Zeile `PGR-E6001`, *Die Anwendung kann sich nicht verbinden*) | P6 | Codes 2, 6, 7, 9, 42 (Stand-in-Server, kein echter Kerberos-Server); Meldung nennt „(Code N)“ | Code 3, 5, 10 werden vermittelt (P1) |
| Lesbare Nachricht außerhalb der Reihenfolge: `PGR-E4002`; Verbindung beendet in der Anmeldung: `PGR-E4003`; Exit 4 (*Hinweise*; §7 beide Zeilen) | P6b | drei Fälle, Meldung nennt den Code, Log-Zeile, Exit 4 | — |
| Die Anmeldung zählt zur Zeit von `offset_ms`, auch die Zeit, die die Anwendung für das Passwort braucht (§4 *Zeitangaben*) | P10 | 1505 mit 1,5 s Wartezeit | 4 ohne Wartezeit |
| Geprüft ist das mit PostgreSQL 17 und `psql` aus libpq 17; andere Versionen und Treiber sind nicht geprüft (*Voraussetzung*, README) | Tabelle der Versionen oben | Läufe P1 bis P12 | — |
| §6: „vermittelt die Anmeldung, ohne das Passwort aufzuzeichnen“, „nimmt jede Anmeldung an“ | P1, P3 | wie oben | — |
| Beispiele als Datei (zwei Blöcke) | P11 | Block 1 und 2 unverändert bis auf Host, Pfad, Funktion | — |

**Quantoren-Durchsicht des Bestands** (Merkmal: die Allgemeinregeln, die der neue Abschnitt einschränkt: „record meldet ohne Passwort an“, „Verbindungen ohne Anfrage sind keine Sitzung“, „Zeit der Anmeldung“; Suche `grep -n -i -E "anmeld|passwort|zugangsdaten"` mit den Quantoren `jede/alle/immer/stets/nie/keine/jeweils/nur/ohne` über das Handbuch, dazu dieselben Quantoren über den Diff):

| Fundstelle | Entscheidung |
|---|---|
| §1 *Voraussetzungen*, Zeile „…die den Benutzer ohne Passwort anmeldet“ | ersetzt (neue Fassung ohne Einschränkung auf `trust`) |
| §4 *Eine Anwendung aufzeichnen*, Hinweis „Verlangt die Datenbank ein Passwort, beendet das Werkzeug …“ | ersetzt durch Verweis |
| §4 *Mehrere Verbindungen*, „Keine Sitzung ist, auch mit der Option, …“ | verbunden: ergänzt um die Verbindung, deren Anmeldung nicht zu Ende geht (P4, P8, P6b) |
| §4 *Mehrere Verbindungen*, „Eine Verbindung zählt erst mit ihrer Startnachricht“ | unverändert: Aussage über die Verbindung vor ihrer Startnachricht; die Anmeldung liegt danach und steht in der ergänzten Zeile (P8) |
| §4 *Zeitangaben aufzeichnen*, „Die Zeit … bis zum Senden der Startup-Nachricht … zählt nicht mit“ | unverändert, mit Gegenprobe: die Zeit nach der Startnachricht zählt (P10) |
| §5 Einleitung, „Unbeachtet bleiben: `PGWIRE_RECORDER_PASSWORD` … bei `record` und `replay`“ | unverändert, mit Probe (P5) |
| §5 *Konfigurationsdatei*, „Bei `record` zählen nur Host und Port der URL; Benutzer, Passwort und Datenbank vermittelt die Anwendung selbst“ | unverändert, mit Probe (P5) |
| §6 *Rollen und Rechte*, „prüft keine Zugangsdaten“, „nimmt jede Anmeldung an“ | der erste Teil unverändert (die Datenbank prüft, P4); der zweite mit Probe (P3) |
| §7 `PGR-E4005`, `PGR-E2006` (Passwort bei `play`) | unverändert: Aussagen über `play` |
| §7 *Fehlercodes* `PGR-E4002`, `PGR-E4003`, `PGR-E6001` | ergänzt (P6, P6b) |
| §1 *Wichtig: Aufzeichnungen können vertrauliche Daten enthalten* | unverändert: Anfragen und Datensätze, nicht das Passwort der Anmeldung; mit Probe P1 (0 Treffer) |
| §10 *Support*, „(ohne Zugangsdaten)“ | unverändert: Hinweis an den Melder |

**Befehlsfolge aus Plan §3** (im Wurzelverzeichnis nach `make build`; die Befehle aus `done/slice-doku-ist-stand` §7, Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`, dazu die Gegenrichtungen Hilfe → Handbuch und Katalog → Handbuch):

| Stand | Ergebnis |
|---|---|
| Handbuch `e12ab97` | Optionen des Handbuchs, die `--help` nicht kennt: nur `--h`, `--help`, `--name`, `--rm`, `--version`; Variablen mit der benannten Ausnahme, Codes (`PGR-[EW]…` gegen `internal/hexagon/model/fehler.go`) und beide Gegenrichtungen: leer |
| Mutant (Kopie des Handbuchs mit angehängter Zeile `--tls-ca` `PGWIRE_RECORDER_TLS_CA` `PGR-E2099`) | je rot: `--tls-ca` unter den Optionen, `PGWIRE_RECORDER_TLS_CA` unter den Variablen, `PGR-E2099` unter den Codes |

`grep -n -i -E "slice|welle|adr-|LH-|SPEC-|ARC-|spec/|lastenheft|spezifikation|review|noch nicht|kommt|geplant|später|bisher|künftig"` über die Zeilen des Diffs von Handbuch und README: kein Treffer; im Handbuch insgesamt kein Treffer auf Slice, Welle, ADR, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review. README: der Bestandssatz zu `make doc-trace` (Zeile 99) nennt Slices als Gegenstand des Berichts, nicht als Stand des Produkts. `make docs-check`: 744 Dateien, 0 Befunde.

**Abweichungen und Funde für den Planner:**

- **Verhalten neben der Spezifikation:** keines gefunden. Der Fund zum Startup-Parameter `password` (Probe oben) ist mit `6a44529` und `0aadcac` in `LH-FA-05.b` bereinigt (Architect); das Handbuch nennt ihn nicht.
- **Bewusste Lücke (Review F-7):** Das Handbuch nennt nicht die nicht lesbare Nachricht des Servers und die Clientnachricht außer `p` und `Terminate` in der Anmeldung (`PGR-E6001`); keine Probe des Implementers dazu (§1). Verifikation V-1: Die Sätze zu `PGR-E4002` sind auf „lesbare Nachricht außerhalb der Reihenfolge“ gefasst (P6b: Code 11 zuerst, `AuthenticationOk` dann Code 3); die nicht lesbare Servernachricht (`PGR-E6001`, Probe des Verifiers mit Ersatz-Server) bleibt ungenannt.
- **Nachgezogen nach Review F-3:** Die Sätze zur Anmeldung in §1, §4 und §6 nennen „geprüft mit PostgreSQL 17“; der Satz zu anderen Verfahren bezieht sich auf das Werkzeug und nennt, dass die Codes 2, 6, 7, 9 und 42 von einem Ersatz-Server kamen (P6); die Zeile `PGR-E6001` und *Die Anwendung kann sich nicht verbinden* sprechen vom Werkzeug, das auf die Anforderung trifft.
- **Nicht gefahren, nicht genannt:** Windows und macOS; andere Treiber (pgx, JDBC, `psycopg`) und andere libpq-Versionen; PostgreSQL außer 17 (die Matrix liefert `slice-v1-abschluss-postgres-versionen`); ein echter Server mit Kerberos, GSSAPI oder SSPI (die Codes 2, 7, 9 kamen von einem Stand-in-Server); ein Client, der `SCRAM-SHA-256-PLUS` anbietet, ohne dass der Server es tut; gleichzeitige Verbindungen in der Anmeldung; Herunterfahren während der Anmeldung (das Handbuch sagt dazu nichts).
- **Randformen:** Keine nicht genannte oder offene Randform entschieden. Der Satz „Das Verfahren bestimmt die Datenbank“, die Aussagen zu Channel Binding und zum Ende des Servers (`PGR-E4003`) gibt `LH-FA-05.b` wieder; das Ende des Servers steht im Handbuch, weil die Spezifikation es entscheidet (*Ende des Servers*) und P6b es am Binary zeigt.
- **Abdeckungstabellen:** unberührt.

**Läufe und Größe:**

- `make docs-check`: 744 Dateien, 0 Befunde (nach jeder Änderung am Handbuch, zuletzt auf `e12ab97` mit dem Plan); `make kopf-check`: Exit 0.
- `make gates-schnell` auf dem Baum mit diesem Absatz (letzter Commit des Implementers): grün; `echo '{}' | bash .claude/hooks/stop-require-gates.sh` antwortet `{"decision":"approve"}` (Handbuch, README und Plan liegen auf der Positivliste, `make gates` voll ist nicht gelaufen; ein voller Lauf ist nicht nötig, weil kein Code, Test, Gate oder Skript im Diff steht).
- Größe (`AGENTS.md` §3.14, Summe aus Einfügungen und Löschungen von `git diff --numstat ea5019b..<letzter Commit> -- docs/user README.md`, Stand `e12ab97`, danach ändern nur Pläne): Handbuch 96 (81 + 15), README 7 (4 + 3), zusammen **103**; Marke 300 und Meldeschwelle 220 nicht erreicht; Schätzung des Architect 120 bis 200 mit Review-Runde (ohne sie 103, im Band darunter, weil kein Satz der Schätzung ein zweites Mal gefasst werden musste). Der Plan zählt nicht in die Metrik (rund 200 Zeilen §7).
- **Nach dem Review** (F-1, F-3, F-5, F-7; Handbuch `1291bcd`, Plan `30b9b49`): `make gates` (voll, ein Lauf, kein pkill, kein paralleler Lauf) auf `30b9b49`: Exit 0, darin `d-check: 745 Datei(en) geprüft, 0 Befund(e)`, `run-integration-tests: gruen`, `lint-gegenprobe: gruen`, `stop-hook-gegenprobe: gruen`; deckt die Spezifikationsänderung `6a44529` und `0aadcac` ab. Der Stop-Hook antwortet `approve`. `make kopf-check`: Exit 0. Größe auf `1291bcd` (letzter Commit an Handbuch und README): Handbuch 98 (83 + 15), README 7, zusammen **105**; Marke 300 und Meldeschwelle 220 nicht erreicht. Dieser Absatz ist der letzte Commit und ändert nur §7 (Positivliste).
- Liefer-Punkte 2, Schichten eine (Dokumentation), unverändert gegenüber §8.

### Closure

**Review und Verifikation.** Review (`docs/reviews/2026-10-14-review-slice-v1-abschluss-anmeldung-doku.md`: 0 HIGH, 1 MEDIUM, 2 LOW, 4 INFO): F-1 (MEDIUM, Plan folgte der Spezifikationsänderung nicht) behoben in `30b9b49`; F-2 (LOW, Satz zum Server in der Spezifikation) behoben durch den Architect in `0aadcac`; F-3 (LOW, Version des Servers nicht genannt) behoben in `1291bcd`; F-4 bis F-6 ohne Handlung (Mechanik von `make slice-mv`, DoD-Wortlaut, Zeilenumbruch ohne Konvention); F-7 (INFO) bleibt als benannte Lücke in §1. Verifikation (`docs/reviews/2026-10-14-verifikation-slice-v1-abschluss-anmeldung-doku.md`): DoD-Punkte 1 und 2 bestätigt, Punkt 3 mit Einschränkung (der volle Lauf ist nicht neu gefahren, Beleg `30b9b49` und Freigabe des Hooks). V-1 (LOW, der Satz zu `PGR-E4002` war für eine nicht lesbare Nachricht des Servers falsch, dort gilt `PGR-E6001`) behoben in `67cc189`: Hinweis und Tabellenzeile sprechen von einer „lesbaren Nachricht außerhalb der Reihenfolge“, die P6b belegt. V-2 (Kerberos, GSSAPI, SSPI nur mit Ersatz-Server) verlangt keine Änderung, das Handbuch sagt „Ersatz-Server“. V-3 bis V-5 INFO.

**Größe** (`AGENTS.md` §3.14, Summe aus Einfügungen und Löschungen von `git diff --numstat ea5019b..67cc189 -- docs/user README.md`, `67cc189` ist der letzte Commit an Handbuch und README): Handbuch 98 (83 + 15), README 7 (4 + 3), zusammen **105**; Marke 300 und Meldeschwelle 220 nicht erreicht; Schätzung des Architect 120 bis 200 unterschritten. Die Bedingung „nicht in einer Review-Sitzung prüfbar“ trat nicht ein (ein Review-Lauf). Liefer-Punkte 2, Schichten eine (Dokumentation), wie §8.

**Entscheidungen zu den Übergaben.**

- **F-4, F-5, F-6, V-2 (Rest), V-3, V-4:** ohne Handlung. F-6 und die langen Zeilen: Bestand bleibt bewusst stehen, keine Konvention verlangt einen Umbruch. V-3: Die Befehlsfolge war für diesen Slice nie rot, weil er keinen Namen einführt, den Hilfe und Katalog nicht kennen; der Plan behauptet kein „rot vorher“, nur der Mutant ist rot.
- **F-7 (nicht lesbare Nachricht des Servers im Aufbau, Clientnachricht außer `p` und `Terminate`, beide `PGR-E6001`):** bleibt eine benannte Lücke ohne Probe des Implementers (Verifier hat beide gefahren, `PGR-E6001`, Exit 6); das Handbuch sagt dazu nichts, kein Satz widerspricht ihr. Kein Nehmer: eine Adresse dafür gibt es nicht, ein Slice dafür wäre ohne beobachtbaren Anlass geplant.
- **Die Mitteilung des Verifiers zu den Risiko-Ausgängen** war ein Vorschlag; die Ausgänge in §6 sind die Entscheidung dieser Closure (Linux: weiter offen; Handbuch-Teil: entfallen; Version des Servers: weiter offen, mit dem Grund zum Nehmer in §6).

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge): siehe §6; Beispiel nur unter Linux: *weiter offen* (Register); Teil von Handbuch und README nicht geliefert: *entfallen*; Version des Servers: *weiter offen* (Register).

**Register** (Zähler = Dateien unter `evidence/`, gelesen am Stand dieser Closure):

- `BEO-REPO/plan-folgt-korrektur-nicht` — neue Datei `evidence/slice-v1-abschluss-anmeldung-doku.md` (F-1: der Plan nannte die Änderung der Spezifikation durch den Architect nicht): 44 → 45. Verkörpert (`AGENTS.md` §3.9, `make kopf-check`); die Regel gilt, sie wurde erst nach dem Review angewandt. `make kopf-check` prüft nur Kennungen aus §1 und §2, nicht eine Änderung an der Spezifikation während des Slice.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — neue Datei `evidence/slice-v1-abschluss-anmeldung-doku.md` (F-2, F-3, V-1): 54 → 55. Verkörpert (`AGENTS.md` §3.11).
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` — kein neues Auftreten; vermerkt in `state.md`. Der Zähler bleibt bei 3, der Ausgang ist weiter nicht zugewiesen (Lese-Schritt der Closure von welle-v1-abschluss).
- `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` — kein Auftreten (das Handbuch nennt die Version des Servers und die Grenze, der Gegenbeleg); vermerkt in `state.md`, bleibt 1×, offen.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` — kein Auftreten; vermerkt „Teil geliefert“ in `state.md`, bleibt 1×, offen.
- Kein nicht verkörperter Eintrag erreicht mit diesem Slice 3× neu (`BEO-REPO/verhalten-nur-unter-linux-geprueft` steht seit `slice-v1-abschluss-sqlite-format-lesen` auf 3×).

**Lerneintrag — geschärfte Regel (Vorschlag des Planners, noch nicht verkörpert).** Der Handbuchsatz zu `PGR-E4002` sagte „eine Nachricht, die der Ablauf nicht vorsieht“ und schloss damit die nicht lesbare Nachricht des Servers ein, für die die Spezifikation `PGR-E6001` regelt (V-1); die Durchsicht des Bestands (Quantoren) hatte den *vorhandenen* Text geprüft, nicht den *neuen* Satz gegen die Fallunterscheidung der Spezifikation. Die Zeile F-7 des Plans hielt die Serverseite für eine Lücke ohne Aussage, das Handbuch sagte aber etwas. Vorgeschlagen ist eine Zeile in `.claude/commands/implement-slice.md` Schritt 19: Nennt ein neuer Handbuchsatz einen Meldungscode, vergleicht der Implementer die Ursachen des Codes im Satz mit der Fallunterscheidung der Spezifikation (der Katalog in `internal/hexagon/model/fehler.go` und die Randformen-Tabelle); der Satz nennt genau die Ursachen, die eine Probe zeigt, oder grenzt die übrigen ab („lesbare“, „im Aufbau“), und eine bewusste Lücke in §1 gilt nur, wenn das Handbuch zu ihr schweigt. Ob sie aufgenommen wird, entscheidet der Nutzer; bis dahin trägt sie kein Feld `liegt in`. **Kein neuer Sensor:** Welche Ursachen ein Satz meint, sieht kein Skript. Zweiter Befund derselben Klasse und derselben Folge: Die Zeile zur Version des Gegenübers aus `slice-v1-abschluss-tls-client-doku` (F-3 hier: „PostgreSQL 17“ fehlte an vier Stellen) ist ebenfalls noch nicht verkörpert.

**Paarungen.** (a) Anker: Diese Closure trägt kein Feld `liegt in`; die Regel ist nicht verkörpert, also kein Zielort zu prüfen. (b) Folge-Slice: Diese Closure weist keinem Slice etwas zu; `slice-v1-abschluss-postgres-versionen` bleibt Adresse der Beobachtung zum Server (Eintrag dort unterblieben, siehe §6). (c) Register: die fünf genannten Kennungen existieren als Verzeichnis mit nicht leerem `evidence/`.

**Nachzählen.** Keine Sendung in einen Nehmer, daher kein Nachzählen von Liefer-Punkten und Schichten. `slice-v1-abschluss-postgres-versionen` ist unverändert (drei Liefer-Punkte, vier Schichten laut dessen §8).

**Stand der Gates.** `make gates` (voll) auf `30b9b49` grün (Belege des Implementers oben); seither ändern nur Handbuch (`67cc189`), Pläne, Reviews und Register; auf dem Stand der Closure bestätigt `make gates-schnell` den Baum, der Stop-Hook gibt auf Positivlisten-Unterschieden frei.

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
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-anmeldung`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-anmeldung` (Planner, 2026-10-14, `AGENTS.md` §3.13): Das Risiko zur Version des Servers in §6 gehört zu den beiden Liefer-Punkten (Handbuch, README) und fügt keinen hinzu. Liefer-Punkte: 2. Schichten: eine (Dokumentation), kein Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
