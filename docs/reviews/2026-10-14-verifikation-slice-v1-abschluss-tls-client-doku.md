# Verifikation: slice-v1-abschluss-tls-client-doku

**Rolle:** Verifier. **Datum:** 2026-10-14. **Spanne:** `59e14db` (Move nach `in-progress/`) bis `174640b` (HEAD). **Gegenstand:** [slice-v1-abschluss-tls-client-doku](../plan/planning/in-progress/slice-v1-abschluss-tls-client-doku.md) gegen DoD (§2), den geschlossenen Code-Slice [slice-v1-abschluss-tls-client](../plan/planning/done/slice-v1-abschluss-tls-client.md), [`LH-FA-23`](../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client) und in der Spezifikation [`LH-FA-23.a`](../../spec/spezifikation.md#lh-fa-23a--tls-zum-client), `AGENTS.md` §3.9, §3.11, §3.13, §3.14 und Modul 11. Review: [Review-Report](2026-10-14-review-slice-v1-abschluss-tls-client-doku.md). Vorgänger-Verifikation des Code-Slice: [Verifikation](2026-10-13-verifikation-slice-v1-abschluss-tls-client.md).

Alle Belege unten sind in diesem Lauf selbst erzeugt, nicht aus §7 übernommen. Aufbau: Binary `pgwire-recorder:dev` (`make build`; `git diff 59e14db..HEAD -- internal cmd test spec harness .claude` ist leer, der Code ist seit dem Code-Slice unverändert), gepinntes `postgres:17-alpine` aus `harness/mk/integration.mk`, `psql` (libpq 17) und `python:3.13-alpine` (`ssl`, OpenSSL 3.5) als Clients im Netzwerk-Namensraum des Recorders, eigenes Docker-Netz `vrf-tlscd-net`, Container als Benutzer 1000. Läufe nacheinander, nie parallel; danach Container und Netz entfernt (`docker ps -a` und `docker network ls` zeigen nur den Bestand `pgr-dev`), kein `pkill`. Scratch-Dateien (Test-PKI, Aufzeichnungen) liegen außerhalb des Repos. `git status` im Repo war vor und nach den Läufen sauber.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| 1 Benutzerhandbuch beschreibt, was der Code-Slice liefert (Abschnitt, Grenz-Sätze §1/§4/§7, Einstellungen §5, Fehlercode-Tabelle, Sitzungs-Sätze, `config show`), jedes Beispiel als Datei gegen das Binary, keine Verweise, Optionen gegen `--help`, Codes gegen Katalog | bestätigt | Siehe Abschnitte *Beispiele als Datei*, *Verhaltensaussagen*, *Befehlsfolge*, *Ist-Zustand und Verweise* unten. Alle vier Beispiele laufen; jede neue oder geänderte Verhaltensaussage ist am Binary bestätigt. |
| 2 `README.md` nennt TLS ab 1.2 mit `--tls-cert`/`--tls-key`, `PGR-E6003`, `--allow-plaintext`; Verweise nach oben ohne Stand | bestätigt | Diff `README.md` (4 + 1): ein Satz im Ist-Zustand. „mindestens TLS 1.2“ (TLS-1.2-Client verbindet, `maxtls11` scheitert), „ohne ein Zertifikat der Anwendung zu verlangen“ (`psql` ohne Client-Zertifikat verbindet), `PGR-E6003` und `--allow-plaintext` am Binary. Kein Treffer auf Slice, Welle, ADR, „noch nicht“, „geplant“, „kommt“ im Diff; die Verweise auf `spec/` und `docs/plan/` (Zeilen 13, 95, 97) stehen im Bestand ohne Stand-Aussage. Der DoD-Wortlaut ist seit `174640b` an den Bestand angepasst (Review F-2). |
| 3 `make gates` grün | bestätigt (mit Einschränkung der Kette) | Der Stop-Hook `.claude/hooks/stop-require-gates.sh` antwortet auf HEAD `{"decision":"approve"}`; er verlangt keinen vollen Lauf: Die Unterschiede zum Stand des vollen Laufs (`dbe6232`) liegen auf der Positivliste (`docs/plan/planning/`, `docs/reviews/`). Statt der vollen Kette selbst gefahren: `make gates-schnell` auf HEAD, Exit 0 (8 s, Stufen zum Teil `CACHED`, Eingabe byte-gleich); `make docs-check` Exit 0, `d-check: 733 Datei(en) geprüft, 0 Befund(e)` (732 des Implementers plus der Review-Report). Ein voller Lauf ist nicht neu gefahren; der Beleg „grün, 7 min 32 s auf `dbe6232`“ in §7 ist deshalb nicht von mir bestätigt, sondern nur durch die Hook-Freigabe gestützt. |
| 4 Review durchgeführt | bestätigt | Report liegt vor (0 HIGH, 0 MEDIUM, 1 LOW, 4 INFO); F-1 bis F-3 im Plan behoben (unten). |
| 5 bis 8 Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen, DoD-Häkchen | offen (erwartet) | Alle Häkchen leer; §7 trägt bisher die Belege des Implementers; die drei Risiken in §6 stehen auf „offen bis Closure“. Closure-Schritte vor dem `git mv` nach `done/`. Zu den Risiken siehe *Risiko-Ausgänge* unten (Vorschläge, kein Eintrag). |

## Beispiele als Datei gegen das gebaute Binary (AGENTS.md §3.11)

Die vier Blöcke des Abschnitts *Verschlüsselte Verbindungen der Anwendung annehmen* (Zeilen 1142 bis 1177 von `docs/user/benutzerhandbuch.md`) per `sed` unverändert in Dateien gezogen und ausgeführt (nur `PGUSER=postgres` für die Anmeldung und `docker run` als Funktion `pgwire-recorder`):

| Block | Ergebnis |
|---|---|
| `openssl req -x509 …` | läuft mit OpenSSL 3.0.13 (Host) und 3.5.7 (`golang:1.27`): `server.pem` beginnt mit `BEGIN CERTIFICATE`, `server.key` mit `BEGIN PRIVATE KEY`; mit den 3.5.7-Dateien startet `replay` ohne Fehler. Die Angabe „der Aufruf lief mit OpenSSL 3.3“ ist von mir nicht nachfahrbar (kein Image mit 3.3 vorhanden, netzlos); sie stützt sich auf §7 (OpenSSL 3.3.7 in `alpine:3.20`). Sie nennt nur, was gelaufen ist, und ist nicht falsch, wenn §7 stimmt; siehe V-2. |
| `record --listen 127.0.0.1:15432 --upstream postgres:5432 --output ./recordings/users.yaml --tls-cert … --tls-key …` | läuft; `psql sslmode=require` liefert `1`, Aufzeichnung enthält `sql: select 1 as eins`. `sslmode=disable` wird mit `PGR-E6003` abgewiesen, Stopp danach Exit 6; `prefer` verbindet (TLS). |
| `psql … sslmode=verify-full sslrootcert=./server.pem` (`host=localhost`) | liefert `1`. |
| `replay`-YAML (`replay:` mit `tls_cert`, `tls_key`) mit `replay --config` | läuft; `psql verify-full` liefert `1`; `config show --config` mit derselben Datei Exit 0 ohne Meldung. Beim Beenden `PGR-W2001` (nicht verbrauchte Sitzungen der Aufzeichnung), Exit 0. |

Gegenprobe zur Aussage „dieselbe Aufzeichnung über beide Verbindungsarten“: `record --allow-plaintext`, einmal `sslmode=require`, einmal `disable`, jeweils `select 1 as eins`: `diff` der beiden Aufzeichnungen leer.

## Verhaltensaussagen der neuen und geänderten Handbuch-Stellen

Alle am Binary bestätigt (Exit-Codes und Meldungen aus dem Lauf):

| Aussage | Befund |
|---|---|
| `--tls-cert` allein, `--tls-key` allein: `PGR-E2001`; `--allow-plaintext=true` (auch als bloße Option und als Umgebungsvariable) ohne Dateien `PGR-E2001`; `=false` ohne Dateien gültig | bestätigt, Exit 2 bzw. lief |
| Quellen mischen (Zertifikat aus Umgebung, Schlüssel aus Datei) | lief |
| Relativer Pfad ab dem aktuellen Verzeichnis | im Arbeitsverzeichnis `/p`: lief; im Arbeitsverzeichnis `/`: Exit 2, „Datei nicht vorhanden“ (beide Richtungen) |
| `play`: Optionen unbekannt, Variablen `PGR-E2001`, `play.tls_cert` in der Datei `PGR-E2004` | bestätigt, für alle drei Variablen und die Option |
| `allow_plaintext: yes` in der Datei | `PGR-E2004`, „erlaubt sind true und false“ |
| Zertifikatsdatei nur `CERTIFICATE`-Blöcke; Schlüssel als Zertifikat, `beides.pem` für beide Optionen | `PGR-E2007`, „ein PEM-Block hat nicht den Typ CERTIFICATE“, Exit 2 |
| Schlüsseldatei: `PRIVATE KEY` (EC, RSA, Ed25519), `RSA PRIVATE KEY` (PKCS #1), `EC PRIVATE KEY` (SEC1) | lief; zwei Blöcke „mehr als ein PEM-Block“; X25519 „weder RSA noch ECDSA noch Ed25519“; verschlüsselt (`ENCRYPTED PRIVATE KEY`) „der Schlüssel ist verschlüsselt“, ohne Rückfrage, Exit 2 |
| Schlüssel passt nicht zum Zertifikat | „der Schlüssel passt nicht zum Zertifikat“, Exit 2 |
| Zertifikat abgelaufen / noch nicht gültig | je Exit 2, „abgelaufen“ bzw. „noch nicht gültig“ |
| Erster Fehler zuerst: Zertifikat, dann Schlüssel | abgelaufenes Zertifikat plus verschlüsselter Schlüssel: `--tls-cert`; gültiges Zertifikat plus verschlüsselter Schlüssel: `--tls-key` |
| Dateiarten: Verzeichnis, FIFO (beide Optionen), Gerät (`/dev/null`) | je „keine reguläre Datei“, Exit 2, sofort; ein Socket habe ich nicht neu gefahren (Test `TestTLSDateiKeineRegulaere`, vom Code-Slice verifiziert) |
| Meldung nennt Option und Ursache, nicht den Pfad | in allen Läufen so; fehlende Datei: „Datei nicht vorhanden“ ohne Pfad |
| Text außerhalb der Blöcke unbeachtet, in beiden Dateien | lief (Text vor, zwischen, nach den Blöcken der Zertifikatsdatei; vor und nach dem Block der Schlüsseldatei) |
| Kette: erster Block Zertifikat des Servers, Kette wird ausgeliefert; Gültigkeit nur am ersten Zertifikat | `verify-full` gegen die Wurzel: mit Kette `1`, nur Blatt „certificate verify failed“ (beide Richtungen); zweites Zertifikat abgelaufen: lief |
| Prüfung vor der Ausgabedatei, vorhandene Ausgabedatei unberührt, nicht `PGR-E2002` | bei fehlender Zertifikatsdatei: `PGR-E2007`, Ausgabedatei unverändert („alt“); ohne die Option und mit derselben Datei: `PGR-E2002` |
| Läuft das Zertifikat nach dem Start ab, werden neue Verbindungen angenommen | Zertifikat mit 16 s Gültigkeit: Verbindung vor und nach Ablauf mit `require` je `1`; `verify-full` danach: „certificate verify failed“; Neustart mit demselben Zertifikat: Exit 2 „abgelaufen“ |
| `config show` liest die Dateien nicht | Pfade, die nicht existieren: Exit 0, Ausgabe nennt sie |
| Antwort `S`, TLS mindestens 1.2; Client mit höchstens 1.1 scheitert; kein Client-Zertifikat angefordert | TLS 1.2 und 1.3 verbinden; höchstens TLS 1.1: Log `PGR-W3002`, „protocol version not supported“, Client „NO_PROTOCOLS_AVAILABLE“; `psql` ohne Client-Zertifikat verbindet |
| Scheitert der TLS-Aufbau (nicht vertraut, abbricht, Klartext nach `S`): `PGR-W3002` im Log mit `code`, `remote`, Ursache; Exit-Code unverändert | bestätigt (Abbruch nach dem `S`: `error=EOF`; Klartext nach dem `S`: „first record does not look like a TLS handshake“); Stopp danach Exit 0 |
| Direktes TLS: `PGR-W3003`, auch mit `--allow-plaintext` | `psql sslnegotiation=direct`: Client „EOF detected“, Log `PGR-W3003`, mit und ohne die Option |
| Weitere Anfrage nach Verschlüsselung auf verschlüsselter Verbindung `PGR-E6001` | Antwort `SSLRequest` und `GSSENCRequest` je `PGR-E6001` |
| Ohne `--allow-plaintext`: Klartext `PGR-E6003`, nicht aufgezeichnet, Exit 6 beim Beenden; auch nach `GSSENCRequest` (Antwort `N`) | bestätigt bei `record` und `replay`; mit `--allow-plaintext` wird der Klartext-Client nach `GSSENCRequest` angenommen |
| `PGR-E6002` bei anderer Protokollversion | nicht neu gefahren; Test `TestTLSRangfolge/andere_Version` laut §7, die Aussage ist die Bestandsregel des Handbuchs |
| Kein `close_notify`; Anwendung sieht das Ende der TCP-Verbindung | Python mit `suppress_ragged_eofs=False`: nach `Terminate` „EOF occurred in violation of protocol“ |
| Schlüssel in keiner Log-Zeile, auch `debug` | `record` und `replay` mit `--log-level debug` und drei Verbindungsarten: 0 Treffer für die zweite Schlüsselzeile (Gegenprobe des Suchmusters auf die Schlüsseldatei: 1 Treffer); `TestRunTLSLogOhneSchluessel` läuft in beiden Kommandos (`replay`, `record`), V-1 der Vorgänger-Verifikation ist im Code-Slice behoben |
| Sitzungen, `record`: TLS-Aufbau gescheitert und `PGR-E6003` sind keine Sitzung | `record --record-empty-sessions` mit Abbruch nach dem `S`, Handshake und Ende, `disable`, Klartext nach dem `S`, `GSSENCRequest`: danach genau eine Sitzung (die gültige) in der Aufzeichnung |
| Sitzungen, `replay first-request`: ebenso keine Verbindung | `PGR-W2001` nennt nach scheiternden Verbindungen unverändert „4 Sessions, die erste mit Kennung 1“, nach einer gültigen „3 … Kennung 2“ |
| Fehlerbehebung (`PGR-E6003`, `PGR-W3002`, `PGR-W3003`) | wie oben |
| `PGR-W3003`-Zeile der Tabelle (direktes TLS), `PGR-E6001`-Zeile, `PGR-E2004`-Zeile | wie oben |

## Befehlsfolge des Plans §3 (rot vorher, leer nachher)

Selbst nachgefahren (Hilfe aus dem gebauten Binary für `--help` und `record`, `replay`, `play`, `config show`; Katalog `internal/hexagon/model/fehler.go`; Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`), dazu die beiden Gegenrichtungen:

| Stand | Ergebnis (selbst gefahren) |
|---|---|
| Handbuch `59e14db` | rot in den Gegenrichtungen: Hilfe → Handbuch `--allow-plaintext`, `--tls-cert`, `--tls-key`, `PGWIRE_RECORDER_ALLOW_PLAINTEXT`, `_TLS_CERT`, `_TLS_KEY`; Katalog → Handbuch `PGR-E6003`, `PGR-W3002`. Die Vorwärtsrichtung (Handbuch → Hilfe) ist vorher leer, weil das Handbuch nichts Neues nannte. |
| Handbuch HEAD | Optionen nur `--h`, `--help`, `--name`, `--rm`, `--version`; Variablen, Codes und alle drei Gegenrichtungen leer |
| Mutant (angehängt: `--tls-ca`, `PGWIRE_RECORDER_TLS_CA`, `PGR-E2099`) | je rot: `--tls-ca`, `PGWIRE_RECORDER_TLS_CA`, `PGR-E2099` |

Das deckt sich mit §7. Die Tabelle in §5 trägt Standardwerte und Kommandos; sie stimmen mit `--help` (Zeilen `--tls-cert`, `--tls-key`, `--allow-plaintext[=true|false]` in `record` und `replay`, Standard `false`, `true` nur mit `--tls-cert`).

## Ist-Zustand und Verweise (grep)

- Handbuch, `grep -n -i` auf Slice, Welle, ADR-, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review, „noch nicht“, „kommt“, „geplant“, „später“, „bisher“, „künftig“, „vorerst“: 11 Treffer, alle im Sinn eines Zustands oder einer Reihenfolge, keiner ein Verweis oder eine Aussage über Künftiges. Neu im Diff: Zeile 1214 („noch nicht gültig“, Zustand des Zertifikats), 1219 („kommt vor der der Ausgabedatei“, Reihenfolge), 1595 (dasselbe in der Tabelle); die übrigen sind Bestand (14, 167, 533, 613, 712, 1109, 1375).
- Der neue Abschnitt, die Tabellenzeilen und die Hinweise verweisen nur auf Handbuch-Anker (`#5-einstellungen`, `#container`, `#verschlüsselte-verbindungen-der-anwendung-annehmen`).
- README: keine Chronik, keine Aussage über einen Stand; die Verweise auf `spec/` und `docs/plan/` (Verweis nach oben) stehen im Bestand unverändert.
- Versionszeile: `Version: 0.7`, `Stand: 14.10.2026` (vorher 0.6 und 13.10.2026), `Software-Version: dev` unverändert; `pgwire-recorder version` meldet `dev` (Binary aus `make build`).

## Quantoren-Durchsicht

Gegen das Binary gegengeprüft, soweit der neue Abschnitt einen Quantor aus dem Bestand einschränkt:

- „Das Werkzeug fragt nie nach Eingaben“: ein verschlüsselter Schlüssel endet mit Exit 2 ohne Rückfrage (Containerlauf ohne Standardeingabe). Bleibt wahr.
- „nimmt jede Anmeldung an“ (Rollen und Rechte): TLS fordert kein Zertifikat der Anwendung an (`psql` ohne Client-Zertifikat verbindet). Bleibt wahr.
- „Eine Verbindung zählt erst ab ihrer Startnachricht“ / „Keine Sitzung ist …“: bleiben wahr, jetzt um `PGR-W3002` und `PGR-E6003` ergänzt (oben).
- „Jede Datei muss … regulär sein“, „jede allein ist ein ungültiger Aufruf“, „in beiden Dateien“, „in keiner Log-Zeile“, „genau einen Block“: je am Binary bestätigt (oben).
- §1 *Voraussetzungen* und §4 *Mit einem Datenbanktreiber arbeiten*: „Außer an der Verschlüsselung“, „solange … ohne `--tls-cert`“, „Startet das Werkzeug ohne `--tls-cert`“ sind enge Fassungen; „Mit `--tls-cert` und `--tls-key` nimmt das Werkzeug die Anfrage nach Verschlüsselung an“ bestätigt.
- Bestandsaussagen über die Seite der Datenbank („`record` verbindet ohne TLS zur Datenbank“, §5): `record` mit TLS zum Client und ohne TLS zur Datenbank lief (Beispiel 2; die Upstream-Seite von `pg_stat_ssl` zeigte `ssl = f`). Unverändert wahr.
- Nicht neu gefahren: „auch nicht nach einer Anfrage nach Verschlüsselung“ in *Zeitangaben aufzeichnen* (Bestandssatz); `record --record-timing` über TLS lieferte `offset_ms: 3` für die erste Anfrage, die Wartezeit zwischen Handshake und Startnachricht habe ich nicht neu vermessen (P5 `timing` des Implementers, §7).

## Review-Findings

- **F-1 (LOW, Plan nennt Handbuch-Stellen nicht):** behoben in `174640b`. §1 nennt die Sitzungs-Sätze und den Satz zu `config show` unter *Lieferumfang im Handbuch*, §3 trägt sie in der Zeile `docs/user/benutzerhandbuch.md`, DoD-Punkt 1 nennt sie, §8 zählt nach (Liefer-Punkte 2, eine Schicht). Der Diff von `174640b` ändert nur den Plan.
- **F-2 (INFO, DoD-Wortlaut README):** behoben. DoD 2 nennt den ersetzten Satz „`record` und `replay` nehmen Verbindungen nur ohne TLS an …“ und vermerkt, dass „Alle Verbindungen laufen unverschlüsselt“ schon vorher im README fehlte.
- **F-3 (INFO, Größenzahl):** behoben. §7 nennt den Messbereich `59e14db..dbe6232`. Selbst gemessen (`git diff --numstat`, Summe aus Einfügungen und Löschungen): `dbe6232` 337 (Handbuch 196 + 30, README 4 + 1, Plan 105 + 1; ohne Plan 231), `77efd06` 344 (Plan 112 + 1), HEAD `174640b` 438 (Plan 119 + 4, Review-Report 84 + 0). Die Zahlen in §7 stimmen; eine Rückführungs-Marke in Zeilen nennt §4 nicht, die Metrik aus `AGENTS.md` §3.14 ist dort nicht verlangt.
- **F-4 (INFO, nur libpq 17 belegt):** Das Handbuch nennt keine ungeprüfte Umgebung als geprüft. Es steht kein Windows, macOS, JDBC, pgx als geprüft im neuen Abschnitt (`grep` auf Windows, macOS, JDBC, pgx, psycopg im Diff: kein Treffer). Offen bleibt die Allgemeinheit der Sätze „Mit libpq gelingt `sslmode=verify-full` … nur mit dieser Kette“ und „bei libpq `sslnegotiation=direct`“: belegt ist libpq 17 mit OpenSSL des Postgres-Images (eigener Lauf, gleiches Ergebnis wie Implementer und Review); `sslnegotiation` kennt libpq erst ab 17. Die Sätze nennen die Version nicht, was die Aussage für ältere libpq-Versionen ungeprüft lässt (V-1).
- **F-5 (INFO, lange Zeilen):** nicht behoben, keine Konvention verlangt es; in *Hinweise* stehen weiter zwei Zeilen über 160 Zeichen (Zeilen 1215 und 1238 des Handbuchs). Keine Aktion nötig.

## Plan-vs-Code-Diff

`git diff 59e14db..HEAD --name-status`: vier Dateien, `README.md`, `docs/user/benutzerhandbuch.md`, der Plan, der Review-Report (alle `M` bzw. `A`). Nichts unter `internal/`, `cmd/`, `test/`, `spec/`, `harness/`, `.claude/`, `docs/user/abdeckung-*.md`, `docs/user/benutzerhandbuch-standard.md`: die Ausschlüsse aus §1 sind eingehalten (kein Produkt-Code, kein Sensor, Bestand des Standards unberührt). Der Move `59e14db` ist rein (0 Einfügungen, 0 Löschungen, `next/` → `in-progress/`); der Inhalt steht in eigenen Commits (`AGENTS.md` §3.3); alle fünf Commits nennen die Slice-Kennung. Liefer-Punkte 2 (Handbuch, README), Schichten eine (Dokumentation). Jede Stelle, die der Handbuch-Diff ändert, steht im Plan (§1 Lieferumfang, §3, DoD 1): Abschnitt, Grenz-Sätze §1/§4/§7, Einstellungen §5 (Tabelle und Satz zu den Variablen bei `play`), Fehlercode-Tabelle (`PGR-E2001`, `E2004`, `E2007`, `E6001`, `E6003`, `W3002`, `W3003`), Sitzungs-Sätze, `config show`, Versionszeile. Die Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando` (Variablen bei `play` `PGR-E2001`) ist geliefert und belegt. Der Code-Slice liegt in `done/`; §4 ist erfüllt.

## Risiko-Ausgänge (Vorschläge für die Closure, kein Eintrag)

- *Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut* (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×): Das Beispiel lief in diesem Lauf mit zwei weiteren OpenSSL-Versionen (3.0.13, 3.5.7) unter Linux. Linux im Container bleibt die einzige belegte Umgebung; der Handbuch-Satz nennt keine andere. Ausgang *weiter offen* ins Register (zweiter Beleg, kein Zähler setzen) ist belegbar; *entfallen* wäre falsch, weil Windows und macOS nicht gefahren sind.
- *Slice liefert Handbuch-/README-Teil nicht vollständig* (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×): nicht eingetreten; die Lücke des Reviews (F-1) ist im Plan nachgezogen, und die Befehlsfolge (beide Gegenrichtungen) ist leer. Ausgang *entfallen* mit dieser Begründung trägt.
- *Clients mit eigener Zertifikatsprüfung scheitern am selbst signierten Zertifikat*: Das Handbuch beschreibt den Weg (Zertifikatsdatei als Wurzel; `verify-full` gegen `localhost` bestätigt; Kette). Ausgang *entfallen* trägt.

## Befunde

- **HIGH:** keine.
- **MEDIUM:** keine.
- **LOW (V-1, F-4):** Die Sätze zur Kette und zum direkten TLS sagen „mit libpq“, belegt ist libpq 17. Das ist keine ungeprüfte Umgebung im Satz, aber eine Allgemeinheit über alle libpq-Versionen. Abhilfe in der Hand des Planners: Version nennen („libpq ab 17“ bei `sslnegotiation`) oder in der Closure-Notiz als belegte Umgebung führen. Kein Blocker.
- **LOW (V-2):** Die Angabe „der Aufruf lief mit OpenSSL 3.3“ stützt sich auf §7; in diesem Lauf nicht nachfahrbar (kein Image mit OpenSSL 3.3 netzlos). Der Aufruf läuft mit 3.0.13 und 3.5.7; die Angabe ist die des Implementers. Entweder die Version streichen oder als belegt stehen lassen; kein Blocker.
- **INFO (V-3):** `make gates` (voll) habe ich nicht neu gefahren: Der Stop-Hook gibt auf HEAD frei (Positivliste), `make gates-schnell` und `make docs-check` sind grün. Der volle Lauf auf `dbe6232` (7 min 32 s) bleibt Aussage des Implementers; die Freigabe des Hooks stützt sie, bestätigt sie aber nicht selbst. Eine Pflicht zum vollen Lauf besteht laut Hook nicht.
- **INFO (V-4):** Nicht neu gefahren: Socket als Dateiart, `PGR-E6002` bei anderer Protokollversion, Wartezeit zwischen Handshake und Startnachricht gegen `offset_ms`, Treiber außer `psql` (libpq 17) und Python `ssl`. Alle sind durch Tests des Code-Slice belegt (dessen Verifikation Mutationen gefahren hat) oder Bestandssätze.
- **INFO (V-5):** Offen für die Closure wie in §2: DoD-Häkchen, Closure-Notiz mit Lerneintrag, Register (Kandidat: `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` bekommt keine zweite Datei, wenn kein Handbuch-Teil gefehlt hätte; die Lücke F-1 war eine Planlücke, keine Handbuchlücke), Ausgänge der drei Risiken, drei Paarungen.

## Fazit

Die zwei Liefer-Punkte sind bestätigt: Alle vier Handbuch-Beispiele laufen als Datei gegen das gebaute Binary, jede neue oder geänderte Verhaltensaussage (Sitzungs-Sätze, `config show`, Fehlercode-Tabelle, Dateiarten, Kette, kein `close_notify`, direktes TLS) stimmt am Binary, die Befehlsfolge ist vorher rot und nachher leer (Mutant rot), das Handbuch enthält keine Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein Künftiges, das README nennt TLS im Ist-Zustand mit Verweisen nach oben, die Versionszeile steht auf 0.7 und 14.10.2026. Die Review-Findings F-1 bis F-3 sind im Plan behoben, F-4 und F-5 sind INFO ohne Handbuch-Verstoß (V-1 als LOW-Hinweis). `make docs-check` (733 Dateien, 0 Befunde) und `make gates-schnell` sind auf HEAD grün; der volle Lauf ist nicht neu gefahren, der Hook verlangt ihn nicht. Freigabe zur Closure; vorher DoD-Häkchen, Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge und Paarungen (V-5), und optional V-1 und V-2.
