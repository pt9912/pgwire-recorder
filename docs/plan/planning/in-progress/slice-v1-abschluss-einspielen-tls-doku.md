# Slice slice-v1-abschluss-einspielen-tls-doku: Handbuch und README zu: TLS zum Server beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-einspielen-tls`)

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-einspielen-tls` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-einspielen-tls`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-einspielen-tls` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-einspielen-tls` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-einspielen-tls` liefert, im Ist-Zustand des gebauten Binaries: Aus `slice-doku-ist-stand` (dort §6, *Teilweise geliefert*): Der Slice ersetzt im Handbuch die Grenz-Sätze zu TLS und `sslmode=require` bei `play` — §4 *Eine Aufzeichnung in eine Datenbank einspielen* (Hinweis „verbindet sich unverschlüsselt … `sslmode=require` ist bei `play` ungültig (`PGR-E2004`)“, dazu der Satz zum Klartext-Passwort „… also ohne TLS unverschlüsselt“, der mit TLS nicht mehr allgemein gilt: ein Klartext-Passwort ist über TLS verschlüsselt, ohne TLS nicht; Sendung aus `slice-v1-abschluss-einspielen-tls`), §5 *Konfigurationsdatei* (Absatz zu `sslmode`, soweit er `play` betrifft, dort mit DoD-Punkt 3) und §1 *Voraussetzungen* (Zeile zur Verbindung zur Datenbank) — und im README den Satz „Alle Verbindungen laufen unverschlüsselt“, soweit er `play` betrifft. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] Das Benutzerhandbuch beschreibt in §5 *Konfigurationsdatei* TLS einer Verbindung bei `play` wie geliefert (`sslmode=require` mit TLS und Prüfung des Zertifikats, ein gesetztes `--upstream-tls` vor `sslmode`; Befund F-533), dazu `--upstream-ca` und die Grenze *Zertifikatsspeicher nicht ladbar*; die Abdeckungstabellen sind über `make abdeckung` nachgezogen.
- [ ] `README.md` nennt, was `slice-v1-abschluss-einspielen-tls` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
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
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-einspielen-tls` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-einspielen-tls` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo. Für TLS gilt die Probe-Matrix in §6 (echte PostgreSQL-Server mit eigener CA, abgelaufenem Zertifikat und `ssl=off`, dazu ein Fake-Server); jeder Satz zu TLS im Handbuch trägt eine Zeile davon in §7.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen-tls` liegt in `done/` (Voraussetzung; Schritt 15 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

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

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob `slice-v1-abschluss-einspielen-tls` eine weitere Randform liefert, die dort nicht steht.

**Randformen aus `slice-v1-abschluss-einspielen-tls`** (Prüfung des Architect vom 2026-10-11, vor dem ersten Commit; Ort ist dieser Plan, der Slice legt keinen Vertrag an). Eine Aussage steht im Handbuch nur, wenn die Probe in §7 sie am gebauten Binary zeigt; wo die Probe nicht steht, bleibt der Satz ungenannt. Das Verhalten ist in `slice-v1-abschluss-einspielen-tls` §6 und `LH-FA-20.a` entschieden; hier steht je Randform, ob und wie das Handbuch sie nennt.

*Genannt* (jede mit einer Probe aus der Matrix unten):

- **TLS-Wunsch** — `--upstream-tls` (Option, `PGWIRE_RECORDER_UPSTREAM_TLS`, Schlüssel `upstream_tls` im Abschnitt `play:`) und `sslmode=require` der benutzten Verbindung; ein gesetztes `--upstream-tls`, auch `false`, geht `sslmode=require` vor, aus jeder Quelle; ohne beides keine TLS-Verbindung. Handbuch §5 Optionstabelle, Absatz zu `sslmode` in §5 *Konfigurationsdatei* und Hinweis in §4. `record` und `replay`: `sslmode=require` bleibt `PGR-E2004` (Satz bleibt, nennt jetzt nur `record`; bei `replay` gibt es keine Verbindung), Client-TLS wird nicht genannt (gibt es nicht).
- **Prüfung des Zertifikats** — kein Überspringen; Speicher des Systems plus `--upstream-ca`; unbekannte CA, falscher Name, abgelaufen: `PGR-E4005`; Server antwortet `N` auf die Anfrage nach Verschlüsselung: `PGR-E4005`, kein Rückfall auf eine unverschlüsselte Verbindung; Server lehnt unverschlüsselte Verbindung ab: `PGR-E4005`. Die Meldung nennt die Sitzung und den Grund, **weder Pfad noch Inhalt** der CA-Datei (an der Probe belegt, nicht angenommen).
- **Name** — gegen den Host der Verbindung, wie eingesetzt (nicht gegen den Namen aus der Aufzeichnung); IPv4 und IPv6 gegen die IP-Adressen des Zertifikats, IPv6 ohne Zone. Die Zone nennt das Handbuch nicht, wenn keine Probe sie zeigt.
- **`--upstream-ca`** — PEM-Datei mit einem oder mehreren Zertifikaten, **ergänzt** den Speicher des Systems; ohne TLS `PGR-E2001`; nicht lesbare, keine reguläre Datei, kein Zertifikat: `PGR-E2007` beim Start, Exit-Code 2; die Meldung nennt die Option, nicht den Pfad. Als Beispiel im Handbuch: Verbindung `sslmode=require` in der Konfigurationsdatei, `--upstream-ca` als Datei, gegen den Server mit eigener CA (Exit 0, Wirkung in der Tabelle); das Beispiel steht als Datei im Scratch und läuft gegen das Binary.
- **Speicher des Systems** — `SSL_CERT_FILE` (nur diese Variable, von F-603 Ende zu Ende geprüft) nennt dem Binary die Zertifikate des Systems; die Probe (13) zeigt es am Binary der Stufe `runtime`.
- **Klartext-Passwort** — Satz in §4 neu: Über TLS geht ein Klartext-Passwort verschlüsselt, ohne TLS unverschlüsselt (Probe: Server mit `password` in `pg_hba` und TLS einerseits, ohne TLS andererseits; der Satz sagt nur, was der Server sieht, nämlich die Verbindung, nicht, dass jemand mitliest). Der Satz bleibt bei Klartext; MD5 und SCRAM bekommen keinen neuen Satz.
- **Meldungscodes** — §7 des Handbuchs: neu `PGR-E2007`; `PGR-E2001` ergänzt um `--upstream-ca` ohne TLS; `PGR-E4005` ergänzt um TLS (abgelehnt, Aushandlung, Zertifikat) und Server ohne unverschlüsselte Verbindung; `PGR-E4002` ergänzt um die Antwort auf die Anfrage nach Verschlüsselung, die weder `S` noch `N` ist (Probe: Fake-Server mit einem anderen Byte und mit Ende ohne Antwort).
- **Voraussetzungen und Fehlerbehebung** — §1 *Voraussetzungen* (Zeile zur Datenbank beim Einspielen) nennt, dass `play` auf Wunsch TLS nutzt; der Abschnitt *Die Anwendung kann sich nicht verbinden* und die Zeile zur Verbindung zum Werkzeug betreffen die Verbindung der Anwendung und bleiben unberührt.

*Ungenannt* (kein Satz im Handbuch, jeweils mit Grund; akzeptierte Negative dieser Prüfung, hier festgehalten, damit die nächste Runde sie nicht als neu liest):

- **TLS-Version und Verfahren** — Voreinstellung der Standardbibliothek, ohne Zusage in der Spezifikation (`LH-FA-20.a` *TLS*); ein Satz darüber wäre eine Zusage ohne Test (`AGENTS.md` §3.11).
- **Speicher des Systems nicht ladbar, gilt als leer** — an der Außenseite nicht auslösbar und nicht von „kein passendes Zertifikat“ zu unterscheiden (beide: `PGR-E4005`, unbekannte CA); die Regel *Teilweise geliefert* (`slice-doku-ist-stand` §6) nennt eine Grenze nur mit beobachtbarer Reaktion. Das Handbuch sagt stattdessen den beobachtbaren Fall: Steht das Zertifikat des Servers weder im Speicher des Systems noch in `--upstream-ca`, endet `play` mit `PGR-E4005`. **Abweichung vom Wortlaut der DoD-Zeile 2** („die Grenze *Zertifikatsspeicher nicht ladbar*“): Das ist eine Änderung der DoD und gehört dem Nutzer (siehe Bericht des Architect); der Implementer hält die Zeile bis dahin so, wie die Probe sie trägt.
- **Bytes nach `S` in zwei Segmenten, Handshake-Frist, Dauer** — Segmentierung und Zeitverhalten sind nicht steuerbar und haben keinen Test (Grenze in `LH-FA-20.a` *TLS*); ein Server, der nach `S` schweigt, lässt `play` warten, bis ein Signal kommt: das Handbuch sagt es nicht, weil keine Frist zugesagt ist. Bytes, die mit dem `S` eintreffen, sind `PGR-E4002`; das ist ein Randfall des Protokolls, kein Fall für das Handbuch.
- **SNI, Punkt am Ende des Hosts, Größe der CA-Datei, Zertifikat mit Kopfzeilen im PEM-Block** — Verhalten der Standardbibliothek, kein Fall des Produkts (akzeptiert in `slice-v1-abschluss-einspielen-tls` §6).
- **Schutzziel in der README** — die README sagt nur Beobachtbares: dass `play` auf Wunsch (`--upstream-tls` oder `sslmode=require`) mit TLS verbindet, das Zertifikat prüft (auch gegen eine eigene CA) und TLS-Fehler als `PGR-E4005` meldet; sie sagt weder „sicher“ noch „geschützt vor …“ noch etwas über den Weg der Aufzeichnung. Der Satz „Alle Verbindungen laufen unverschlüsselt“ wird ersetzt durch: `record` und `replay` nehmen Verbindungen nur unverschlüsselt an und verbinden `record` unverschlüsselt zur Datenbank; `play` verbindet auf Wunsch mit TLS. Die Satzteile gelten je nur, wenn die Probe sie zeigt.

*Probe-Matrix* (wiederverwendbar; Setup wie in der Verifikation des Code-Slice, neu gefahren, im Scratch, nichts im Repo): das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz, mit Zertifikaten aus `openssl` in einem Hilfscontainer (alpine mit `openssl`): zwei CAs (`ca.pem`, `ca2.pem`), ein Blatt mit `DNS:pgssl, DNS:localhost, IP:<Adresse im Netz>` von `ca.pem`, ein Blatt für `pgexp` gültig vom 2020-01-01 bis 2020-01-02, `pgssl` mit `ssl=on`, `pgexp` mit dem abgelaufenen Blatt, `pgplain` mit `ssl=off`; `pg_hba.conf` des SSL-Servers mit `hostnossl all all all reject` (Server lehnt Unverschlüsseltes ab) und für die Passwort-Probe `password`; dazu ein Python-Fake (`socket`), der auf das `SSLRequest` mit `N`, einem anderen Byte, Ende ohne Antwort oder `S` und Schweigen antwortet. Aufruf je Zeile: Verbindung in der Datei, `--upstream-ca` als Datei, `PGWIRE_RECORDER_PASSWORD` für den Server mit Passwort. Je Zeile: Aufruf · Exit-Code · Meldung (stderr) · Wirkung (Tabelle in der Datenbank). Zeilen: (1) Erfolg `sslmode=require` + `--upstream-ca`; (2) Erfolg `--upstream-tls` mit `host:port`; (3) Host als IPv4-Adresse im Zertifikat; (4) unbekannte CA (`ca2.pem`); (5) falscher Name (Host ohne Eintrag im Blatt); (6) abgelaufen; (7) `N`; (8) anderes Byte und Ende ohne Antwort; (9) `pgplain` (`ssl=off`) mit `sslmode=require`: Antwort `N`, `PGR-E4005`, kein Rückfall; (10) `--upstream-ca` ohne TLS; (11) `--upstream-ca` nicht lesbar / Verzeichnis / ohne PEM; (12) `--upstream-tls=false` mit `sslmode=require` in der Datei: keine TLS-Verbindung; (13) `SSL_CERT_FILE` statt `--upstream-ca`; (14) Klartext-Passwort über TLS und ohne TLS; (15) `record` und `replay` mit `sslmode=require`: `PGR-E2004`. Die Meldungen aus (4) bis (8) und (11) werden auf Pfad, Passwort und Inhalt der Datei abgesucht (`grep`).

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

### Belege des Implementers

Wird bei Closure um die Closure-Notiz ergänzt (vor dem `git mv` nach `done/`). Kein Produkt-Code, kein Test, keine Mutation des Produkts: Der Slice liefert keinen neuen Vertrag; die Probe-Matrix aus §6 ersetzt die Mutation. Stand der Proben: Binary `pgwire-recorder:dev` (Stufe `runtime`, `make build` auf `15a789a`; `eea9569` ändert nur Handbuch und README), gefahren am 2026-10-11 unter Linux im Container.

**Aufbau (nach den Proben entfernt: alle Container, das Netz, das Hilfsimage und das Scratch-Verzeichnis, nur Namen mit Präfix `impl-tlsdoku-`).** Netz `impl-tlsdoku-net` (`172.30.77.0/24`, `fd00:77::/64`); Server `postgres:17-alpine` mit dem Digest aus `harness/mk/integration.mk`: `pgssl` (`ssl=on`, gültiges Blatt von `ca.pem` mit `DNS:pgssl, DNS:pgpw, DNS:localhost, IP:172.30.77.10, IP:172.30.77.13, IP:fd00:77::10, IP:fd00:77::13`; `pg_hba`: `hostnossl … reject`, `hostssl … trust`), `pgexp` (dasselbe Blatt, gültig 2020-01-01 bis 2020-01-02), `pgplain` (`ssl=off`), `pgpw` (`ssl=on`, `hostssl … password`, `hostnossl … reject`, Benutzer `tester`), `pgpwplain` (`ssl=off`, `host … password`), `pgip` (dasselbe Blatt auf einer Adresse, die nicht darin steht); Zertifikate mit `openssl` 3.5 im Hilfscontainer (`ca.pem`, `ca2.pem`); Fake-Server (`python:3.13-alpine`, `socket`) mit den Antworten `N`, `Z` und Ende ohne Antwort. Aufzeichnung für `play`: mit `record` gegen `pgplain` aufgenommen (`DROP TABLE IF EXISTS probe`, `CREATE TABLE probe(id int)`, `INSERT INTO probe VALUES (1)`); „Tabelle“ unten heißt: danach `select count(*) from probe` gleich 1 auf dem Server, nach `DROP TABLE` vor dem Lauf. Jeder Aufruf lief als `docker run` des Binaries mit den Dateien als Bind-Mount (nonroot).

| Aussage im Handbuch (§4 *Einspielen*, §5, §7) / README | Probe (Matrix aus §6) | Ergebnis |
|---|---|---|
| Beispiel `tls.yaml` mit `sslmode=require` plus `--upstream-ca ./ca.pem` (Handbuch §4, wörtlich als Dateien `tls.yaml`, `ca.pem`, `recordings/users.yaml`, Host `localhost` im Netzwerk-Namensraum des Containers von `pgssl`) | (1) | Exit 0, Tabelle; ohne `--upstream-ca`: Exit 4, `PGR-E4005`, „certificate signed by unknown authority“ |
| `sslmode=require` wünscht TLS; `--upstream-tls`, Option / Umgebung / Schlüssel `upstream_tls` | (1), (2), (12) | Exit 0, Tabelle (`pgssl` lehnt Unverschlüsseltes ab, die Verbindung hatte TLS); Schlüssel `upstream_tls`, `upstream_ca` und Umgebung `PGWIRE_RECORDER_UPSTREAM_TLS`/`_CA`: Exit 0, Tabelle |
| Ein gesetztes `--upstream-tls` geht `sslmode` vor, auch `false`, aus jeder Quelle | (12) | `sslmode=require` in der Datei gegen `pgplain` (lehnt TLS ab): mit `upstream_tls: false`, `--upstream-tls=false`, Umgebung `false` je Exit 0, Tabelle (kein TLS); Umgebung `false` und Option `=true`: Exit 4, „lehnt TLS ab“ (die Option geht der Umgebung vor); `--upstream-tls=false` mit `sslmode=require` gegen `pgssl`: Exit 4, `pg_hba … no encryption` (ohne TLS verbunden); `--upstream-tls` mit `sslmode=disable` in der Datei gegen `pgssl`: Exit 0, Tabelle (mit TLS) |
| Ohne beides kein TLS | (12) | `host:port` gegen `pgplain`: Exit 0, Tabelle |
| Jede Sitzung baut ihre Verbindung mit TLS auf | Aufzeichnung mit zwei Sitzungen gegen `pgssl` | Exit 0; im Log des Servers zwei Zeilen `SSL enabled` mehr als vorher |
| Kein Überspringen; unbekannte Zertifizierungsstelle → `PGR-E4005` | (4) | Exit 4, `PGR-E4005`, `failed to verify certificate … unknown authority`, keine Tabelle (mit `ca2.pem` und ohne CA) |
| Name gegen den Host, wie eingesetzt; Hostname als DNS-Name, IPv4/IPv6 als IP-Adresse | (3), (5) | `172.30.77.10` und `[fd00:77::10]`: Exit 0, Tabelle; Platzhalter `${DBH}`=`pgssl`: Exit 0; `${DBH}`=`pgother-pgssl`, `pgother-pgssl` als Host: Exit 4, „valid for pgssl, pgpw, localhost, not …“; `172.30.77.15` und `[fd00:77::15]` (nicht im Blatt): Exit 4, `PGR-E4005`. Die Zone einer IPv6-Adresse: nicht gefahren, nicht genannt |
| Abgelaufenes Zertifikat → `PGR-E4005` | (6) | Exit 4, „certificate has expired“, keine Tabelle |
| Antwort `N` → `PGR-E4005`, kein Rückfall; Server lehnt Unverschlüsseltes ab → `PGR-E4005` | (7), (9) | Fake `N`: Exit 4, „lehnt TLS ab“; `pgplain` mit `sslmode=require` und mit `--upstream-tls`: Exit 4, „lehnt TLS ab“, keine Tabelle (ein Rückfall hätte die Tabelle angelegt: `pgplain` nimmt Unverschlüsseltes an); ohne TLS gegen `pgssl`: Exit 4, `PGR-E4005`, `pg_hba … no encryption` |
| Antwort weder `S` noch `N` / Ende davor → `PGR-E4002` | (8) | Fake `Z`: Exit 4, `PGR-E4002`, „weder S noch N“; Fake Ende ohne Antwort: Exit 4, `PGR-E4002`. Fake `S` und Schweigen: nicht gefahren, nicht genannt |
| Meldungen nennen weder Pfad noch Inhalt der CA-Datei; Passwort nie | `grep -n -E` über die Ausgaben aller 50 Läufe auf Verzeichnis, Dateinamen und Kopfzeile der CA-Dateien sowie die Passwörter, dazu `grep -n -F` mit je zwei Zeilen Inhalt von `ca.pem` und `ca2.pem` | kein Treffer; die Gegenprobe des Suchmusters auf einen Text mit Pfad: 1 Treffer |
| `--upstream-ca` ergänzt den Speicher; mehrere Zertifikate in einer Datei | (13), (11) | `SSL_CERT_FILE`=`ca.pem` und `--upstream-ca ca2.pem`: Exit 0; `SSL_CERT_FILE`=`ca2.pem` und `--upstream-ca ca.pem`: Exit 0; beide `ca2.pem`: Exit 4; Bündel `ca2.pem`+`ca.pem` als `--upstream-ca`: Exit 0 |
| `--upstream-ca` ohne TLS → `PGR-E2001` | (10) | `host:port` ohne TLS und Datei mit `sslmode=disable`: Exit 2, `PGR-E2001`; auch mit `--upstream-tls=false` und `sslmode=require` |
| Datei nicht lesbar, keine reguläre Datei, kein Zertifikat → `PGR-E2007`, Exit 2, ohne Pfad | (11) | fehlende Datei, Datei mit Rechten 000, Verzeichnis, Textdatei, leere Datei, PEM nur mit Schlüssel: je Exit 2, `PGR-E2007`, Meldung nennt `--upstream-ca` und Grund |
| `SSL_CERT_FILE` zählt zum Speicher des Systems | (13) am Binary der Stufe `runtime` | `SSL_CERT_FILE`=`ca.pem` ohne `--upstream-ca`: Exit 0, Tabelle; mit Datei mit `sslmode=require` ebenso; ohne die Variable: Exit 4 (4). Bestehender Test: `TestE2EPlayTLSZertifikatsspeicherDesSystems` (`test/integration/play_tls_e2e_test.go`) |
| Klartext-Passwort mit TLS auf der verschlüsselten, ohne TLS auf der unverschlüsselten Verbindung | (14) | `pgpw` mit `--upstream-tls`: Exit 0, Tabelle, Server-Log `method=password` und `SSL enabled (protocol=TLSv1.3 …)`; `pgpwplain` ohne TLS: Exit 0, Tabelle, `method=password` ohne `SSL enabled`; `pgpw` ohne TLS: Exit 4, `pg_hba … no encryption`; falsches Passwort: Exit 4, `PGR-E4005`, SQLSTATE `28P01`, das Passwort steht in keiner Ausgabe |
| `record`: `sslmode=require` ungültig (`PGR-E2004`); ungenutzte Verbindung mit `require` gültig; `record` und `replay` nehmen TLS nicht an; `record` verbindet ohne TLS | (15) | `record` mit `require`: Exit 2, `PGR-E2004`; `record` und `replay` starten mit einer ungenutzten `require`-Verbindung in der Datei; `psql sslmode=require` gegen `replay` und `record`: „server does not support SSL“; `record` gegen `pgssl` mit `sslmode=prefer` des Clients: `pg_hba … no encryption`. `replay` hat keine Option `--upstream`: `PGR-E2001` |

**Befehlsfolge aus Plan §3** (Optionen, Variablen, Codes des Handbuchs gegen `--help` und Katalog `internal/hexagon/model/fehler.go`, Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`):

| Stand | Ergebnis |
|---|---|
| `15a789a` (vorher), Richtung Hilfe → Handbuch (`comm -13`) | rot: `--upstream-ca`, `--upstream-tls` fehlen; Katalog → Handbuch: `PGR-E2007` fehlt |
| `eea9569` (nachher) | Optionen des Handbuchs, die `--help` nicht kennt: nur `--h`, `--help`, `--version` (Nicht-Optionen im Text) und `--rm`, `--name` (`docker run`); Variablen mit der benannten Ausnahme: leer; Codes: leer; Richtung Hilfe → Handbuch und Katalog → Handbuch: leer |
| Mutanten an einer Kopie des Handbuchs: Zeile `--upstream-cert`, `PGWIRE_RECORDER_UPSTREAM_CERT`, `PGR-E2008` angehängt | je rot (1 Treffer in der jeweiligen Zeile); ohne die benannte Ausnahme: `PGWIRE_RECORDER_PASSWORD` in der Ausgabe |

`grep -n -i` im Handbuch nach Slice, Welle, ADR-, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review, „noch nicht“, „kommt“, „geplant“, „später“, „bisher“: Treffer nur im Bestand (Zeile 14 „später“, 153 „existiert noch nicht“ in der Voraussetzung einer Aufgabe, 594 „kommt nur zum Einsatz“), keiner in den neuen Sätzen. README: kein Treffer auf Slice, Welle, „noch nicht“, „geplant“, „bisher“, „sicher“, „geschützt“ in den neuen Sätzen. `make docs-check` (d-check): 528 Dateien, 0 Befunde. `make gates` auf dem sauberen Baum `eea9569`: grün (Ausgang 0, 3 min 5 s). Der Commit mit diesem Satz ändert nur §7 dieses Plans; `make gates` auf dessen sauberem Baum ist die Übergabe-Messung (Bericht).

**Abweichungen und Funde für den Planner:**

- **DoD-Zeile 2** nennt die Grenze *Zertifikatsspeicher nicht ladbar*; das Handbuch nennt sie nicht (Entscheidung des Nutzers, Option A, siehe §6). Es sagt stattdessen: Steht das Zertifikat der Zertifizierungsstelle weder im Speicher des Systems noch in `--upstream-ca`, endet `play` mit `PGR-E4005`. Der Planner gleicht die Zeile bei der Closure an.
- **DoD-Zeile 2, Halbsatz Abdeckungstabellen:** Kein Test ist geändert; `make abdeckung` ist nicht gefahren, `make abdeckung-check` (Teil von `make gates`) ist grün.
- **Nicht gefahren, nicht genannt:** Fake `S` und Schweigen, Zone einer IPv6-Adresse, MD5 und SCRAM mit TLS, Zertifikatsketten mit Zwischenzertifikat, der Standardpfad des Systems ohne `SSL_CERT_FILE` (das Zertifikat einer öffentlichen Zertifizierungsstelle wurde nicht gegen das Image geprobt), Windows und macOS (`BEO-REPO/verhalten-nur-unter-linux-geprueft`: alle Proben unter Linux im Container).
- **Bestand, nicht geändert:** Die Optionstabelle in Handbuch §5 führt `--config` vor `--log-level`, `play --help` führt `--log-level` vor `--config`; die neuen Zeilen stehen in der Reihenfolge von `play --help`.
- Die Meldung von `PGR-E4005` mit TLS nennt neben Sitzung und Grund auch `host:port` des Servers; das Handbuch sagt „Sitzung und Grund“, was zutrifft, ohne die Adresse auszuschließen.
- Ausgänge der Risiken aus §6 und die Closure-Notiz: bei Closure.

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
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-einspielen-tls`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-einspielen-tls` (Closure, `AGENTS.md` §3.13): Der Satz zum Klartext-Passwort im Handbuch §4 steht in derselben Datei und demselben Abschnitt wie der Hinweis zu `sslmode=require`, den DoD-Punkt 1 schon führt; er ist ein Satz mehr in diesem Liefer-Punkt, kein vierter. Liefer-Punkte: 3, Schichten: eine (Dokumentation); der Plan liegt innerhalb der Grenze. Den Ausschluss *Produkt-Code* in §1 trifft die Sendung nicht: Sie ändert nur Handbuchtext.

Liefer-Punkte: 3. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
