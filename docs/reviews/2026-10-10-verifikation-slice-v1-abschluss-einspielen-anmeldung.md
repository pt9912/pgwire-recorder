# Verifikation: slice-v1-abschluss-einspielen-anmeldung — 2026-10-10

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-10-review-slice-v1-abschluss-einspielen-anmeldung.md`, F-591 bis F-596); die Übergaben prüfe ich nach.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-anmeldung.md`, Kopf und §1 bis §8, gegen den Gesamt-Diff `f9bc4d0..556f8aa` (20 Dateien, +2510 −82). Architect: `cc58b41`, `482951f`, `0a063cd` (`LH-FA-20.a` *Passwort*, *Verfahren*, *SCRAM-Austausch*). Implementer: `c8b7f44`, `238c78b`, `e74e2fe`, `265b184`, `b263135`, `556f8aa`. Review: `3cc01e4`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `556f8aa`, besonders §1 (*Übernimmt*, *Ausdrücklich NICHT*), §2 (DoD), §3, §6 (*Randformen*, *Risiken*) und §7 (*Belege des Implementers*, beide Tabellen der Nacharbeit)
- `spec/spezifikation.md` `LH-FA-20.a` *Anmeldung*, *Passwort*, *Verfahren*, *SCRAM-Austausch*, *Aufbau*, *Anmelde-Nachrichten*, *Abbruch im Aufbau*; `LH-FA-17.a` (Passwort, Wirkung einer URL, Geheimnisse); `SPEC-033`
- [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md) ganz; [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) (Abschnitt *Entscheidung*); [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)
- am Stand `556f8aa`: `anmeldung.go` ganz, der Diff von `einspielen.go`, CLI (`cli.go`, `upstream.go`, `verbindung.go`) und `bootstrap.go`
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- `AGENTS.md` §3.9 bis §3.13

Beim Start war der Arbeitsbaum sauber, HEAD `556f8aa`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Alles lief unter dem Präfix `ver-anm-` im Scratchpad, nie im Repo.

- **Binary:** eine frische Kopie per `git archive HEAD | tar -x`, daraus die Stufe `runtime` des `Dockerfile` mit eigenem Tag. Ein internes Docker-Netz (`--internal`, festes Subnetz) und ein PostgreSQL-Container mit dem gepinnten Image aus `harness/mk/integration.mk`. `pg_hba.conf`: die Adresse des Recorders `trust`, `play_scram` mit `scram-sha-256`, `play_md5` mit `md5`, `play_pw` mit `password`, `play_trust` mit `trust`; Passwörter `GEHEIM…` (das SCRAM-Passwort mit Leerzeichen, `$` und `%`). Je Benutzer eine Datenbank, die er besitzt. Vier Aufzeichnungen mit `record` gegen `trust` (psql durch den Recorder: `CREATE TABLE`, `INSERT`, `SELECT`). Jeder `play`-Lauf als eigener Container mit `timeout 120`; die Wirkung las ich mit `psql` als Administrator und löschte die Tabelle danach.
- **Fake-Server:** ein selbst geschriebenes Python-Skript (`socket`, `hashlib`, `hmac`, `base64`), nur auf `127.0.0.1`, `play` mit `--network host`. Es rechnet SCRAM-Server-Seite, MD5 und Klartext selbst und protokolliert jede Nachricht des Clients. Das ist die unabhängige Rechnung (Aufgabe 3).
- **Unit-Mutanten:** je Mutant eine frische Kopie von `internal/` per `cp -r` ohne `-p`, genau eine Ersetzung (das Skript bricht ab, wenn der Suchtext nicht genau einmal trifft), im Image der Stufe `test` mit der Kopie als Bind-Mount, `--network none`: `gofmt -l` leer, `go vet` ohne Befund, dann `go test -count=1` der betroffenen Pakete. Der Grundlauf ohne Mutation war grün.
- **E2E-Mutanten:** eine Kopie des Baums, `make test-integration` mit eigenem Tag des Integrations-Images (das Image des Arbeitsbaums blieb unberührt).
- **`-race`:** `golang:1.27` (Debian, mit `gcc`), Quellbaum schreibgeschützt, Modulcache aus dem Test-Image, `--network none`.
- **Danach:** Container, Netz, Kopien und eigene Images/Tags entfernt (siehe Ende des Berichts). `make gates` lief im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `play` meldet sich an einem Server mit `password`, `md5` und `scram-sha-256` an und spielt eine Aufzeichnung ein; das Passwort aus dem Platzhalter der Verbindung, sonst aus `PGWIRE_RECORDER_PASSWORD`, leer gilt als nicht gesetzt; verlangt der Server keines, sendet `play` keines (Integrationstest). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Einspielen gegen die drei Verfahren, Passwort aus `PGWIRE_RECORDER_PASSWORD` bei `host:port` | Binary gegen PostgreSQL 17: Aufzeichnung gegen `trust` (`record`), danach `play --upstream ver-anm-pg:5432` ohne `--user` mit dem Benutzer der Aufzeichnung, je `play_scram`, `play_md5`, `play_pw`: Exit 0, die Tabelle trägt den Wert der Aufzeichnung (`eins-scram`, `eins-md5`, `eins-pw`) | bestätigt |
| Passwort aus dem Platzhalter (`postgresql://play_x:${PWV}@…`, Konfigurationsdatei) | Binary: je Verfahren Exit 0, Tabelle gefüllt | bestätigt |
| Platzhalter geht der Variable vor | Binary: Platzhalter richtig, `PGWIRE_RECORDER_PASSWORD=GEHEIMfalsch`: je Verfahren Exit 0. Umgekehrt (Platzhalter falsch, Variable richtig): `PGR-E4005`, Exit 4 (die Variable springt nicht ein); Platzhalter-Variable leer gesetzt neben richtiger `PGWIRE_RECORDER_PASSWORD`: `PGR-E2005`, Exit 2 (die leere Variable des Platzhalters gilt als nicht gesetzt, die andere Quelle springt nicht ein) | bestätigt |
| Variable bei Verbindung ohne Passwortteil | Binary: `postgresql://play_scram@…` und `play_md5@…` mit `PGWIRE_RECORDER_PASSWORD`: Exit 0 | bestätigt |
| Leere Variable gilt als nicht gesetzt | Binary: `PGWIRE_RECORDER_PASSWORD=` (leer) gegen scram, md5, password: `PGR-E4005` „der Server verlangt ein Passwort, play hat keines“, Exit 4 (wie ohne Variable) | bestätigt |
| Passwort unverändert (kein Kürzen, Nicht-ASCII) | Fake-Server: Passwort `pw GEHEIM ä$%41 ` (Leerzeichen am Ende, `ä`): die Klartext-Antwort ist byteweise das Passwort, MD5 und SCRAM stimmen mit der eigenen Rechnung überein | bestätigt |
| `--user` wirkt auf MD5 | Binary: Aufzeichnung des Benutzers `play_scram`, `play --user play_md5 --database d_md5`: Exit 0 (MD5 rechnet mit dem gesendeten `user`); Fake-Server: `md5-check` gleich, `user` = Benutzer der Startup-Daten | bestätigt |
| Server verlangt kein Passwort, `play` sendet keines | Binary: `play_trust` mit `PGWIRE_RECORDER_PASSWORD` (host:port) und mit Platzhalter: Exit 0; ohne Passwort ebenfalls Exit 0. Fake-Server `trust` mit gesetztem Passwort: der Server sah nur `Q`, `Q`, `Q`, `X`, keine Nachricht `p` | bestätigt |
| Integrationstest | `TestE2EPlayAnmeldung` (6 Fälle), `TestE2EPlayAnmeldungOhneVerlangen`, `TestE2EPlayAnmeldungFehler` (9 Fälle) grün im eigenen `make gates`. E2E-Mutanten, je eine Kopie, `make test-integration`: **E1** Bootstrap reicht das Passwort nicht: rot, 6 Fälle in `TestE2EPlayAnmeldung` und 6 in `TestE2EPlayAnmeldungFehler`; **E2** MD5 mit `benutzer + passwort`: rot, `TestE2EPlayAnmeldung/md5_Variable` und `/md5_Platzhalter_vor_der_Variable`; **E3** Passwort bei `AuthenticationOk` gesendet: rot, 6 Fälle und `TestE2EPlayAnmeldungOhneVerlangen` | bestätigt |

### Punkt 2: Falsches Passwort, fehlendes Passwort, nicht unterstütztes Verfahren und jeder SCRAM-Fehler sind `PGR-E4005`, Exit 4; nichts gesendet bei fehlendem Passwort und nicht unterstütztem Verfahren; Verbindungsende `PGR-E4002`; Signal in der Anmeldung; kein Geheimnis in Ausgabe (Test); Beleg in §7 je Zusage Zusage · Mutation · roter Test. **Zum Teil bestätigt: Verhalten bestätigt, Beleg für eine Zusage falsch (V-162).**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Falsches Passwort | Binary: je Verfahren `PGR-E4005` mit „Fehlerantwort im Aufbau 28P01“, Exit 4, auch mit Platzhalter; Tabelle nicht angelegt | bestätigt |
| Fehlendes Passwort, nichts gesendet | Binary (echter Server): ohne Variable `PGR-E4005`, Exit 4, drei Verfahren. Fake-Server: Klartext, MD5, SASL ohne Passwort: `PGR-E4005`, der Server sah nach dem Startup kein Byte und kein `Terminate`, Ende durch `EOF` | bestätigt |
| Nicht unterstützte Verfahren | Fake-Server: Anforderungen der Codes 2, 6, 7, 9 und 99 (ohne Passwort): je `PGR-E4005` „Anmeldeverfahren (Code n), das play nicht unterstützt“, nichts gesendet. Dasselbe mit Passwort: Code 7 nach einer Antwort bleibt `PGR-E4005` | bestätigt |
| SASL ohne `SCRAM-SHA-256` | Fake-Server: `SCRAM-SHA-256-PLUS` allein: `PGR-E4005` „bietet kein SCRAM-SHA-256 an“, nichts gesendet; leere Liste (`10` + `00`): ebenso; Liste ohne Ende (`10`, kein Byte): `PGR-E4005` „nicht lesbar“ | bestätigt |
| Nicht lesbare Anforderung eines unterstützten Verfahrens | Fake-Server: Klartext mit einem Byte Rest, MD5 mit drei und mit fünf Byte Salz: je `PGR-E4005` „nicht lesbar“ | bestätigt |
| SCRAM-Fehlerfälle | Fake-Server (Tabelle unten): Servernonce ohne Präfix, gleich lang, mit Leerzeichen; Iterationen 0, 10 000 001, `04096`; zusätzliches Attribut in der ersten Nachricht; `e=`; falsche Signatur; Abschluss ohne `v=`; Zusatzattribut hinter `v=`: je `PGR-E4005`, Exit 4, kein `Terminate`. Die Fortsetzung (Code 11) nach dem Abschluss: `PGR-E4002` | bestätigt |
| Erfolgsfälle des Austauschs | Fake-Server: nicht-kanonisches Base64 in `s=` (`YR==` für `a`) gelingt; nicht-kanonisches `v=` (letztes Zeichen mit gesetzten Restbits) gelingt; leeres Salz gelingt; genau 10 000 000 Iterationen gelingen. Jedes Mal `proof-check` gleich, die Nonce hat 18 Byte | bestätigt |
| Zweite Anforderung nach Antwort | Fake-Server, nach Klartext und nach MD5: Codes 3, 5, 10: je `PGR-E4002` „weitere Anforderung … nach der Antwort“; Code 5 mit drei Byte Salz, Code 3 mit Rest, Code 10 ohne Ende: ebenfalls `PGR-E4002`; Codes 2, 6, 7, 9, 99: `PGR-E4005`. Nach SCRAM-Abschluss: Codes 3, 5, 10 `PGR-E4002`, Code 7 `PGR-E4005`, Codes 8, 11, 12 `PGR-E4002` | bestätigt |
| Verbindungsende in der Anmeldung | siehe Tabelle unten; Spec-Zusage `PGR-E4002` hält an `TestAnmeldungSendenScheitert` und `TestAnmeldungAbgelehnt` (Mutation unten) | bestätigt |
| Signal in der Anmeldung | Binary gegen den haltenden Fake-Server (je Klartext, MD5, SCRAM): erstes `SIGINT` (und `SIGTERM`): `play` läuft weiter, der Aufbau endet (Server liefert nach vier Sekunden `AuthenticationOk`), `play` sendet `Terminate` ohne eine Interaktion zu spielen, Zeile „Abbruchsignal“, Exit 0. Zwei Signale, Server hält nach der Passwort-Antwort (bei SCRAM nach dem Beweis): Ende binnen 0,1 s nach dem zweiten Signal (Zeitstempel der Log-Zeilen), Exit 0, der Server sah kein `Terminate` | bestätigt |
| Kein Passwort in einer Ausgabe | Binary: alle Läufe oben (Erfolg, falsches, fehlendes, leeres Passwort, Platzhalter, Variable, `--log-level debug` bei Erfolg und Fehler, `config show` mit gesetzten `PWV` und `PGWIRE_RECORDER_PASSWORD`, unbekannte Datei, nicht erreichbarer Host, falscher Benutzer, `--password GEHEIMopt`): `grep GEHEIM` über alle Dateien von stdout und stderr: 0 Treffer. `config show` nennt die Variable nur mit Namen (`PGWIRE_RECORDER_PASSWORD`) und zeigt `${PWV}` unaufgelöst. `--password` ist `PGR-E2001` („flag provided but not defined“), ohne den Wert. Die vier Aufzeichnungen tragen kein `GEHEIM` | bestätigt |
| Test je Zusage mit Mutation, roter Test (§3.10) | **61 Unit-Mutanten** und 3 E2E-Mutanten gefahren (Tabelle unten): 60 Unit-Mutanten rot aus dem richtigen Grund. **Ein Mutant bleibt grün:** `.Strict()` in `dekodiereBase64`. §6 und §7 nennen dafür `TestAnmeldungScramSalzNichtKanonisch` und `TestAnmeldungScramSignaturNichtKanonisch`; keiner der beiden Tests ist im Repo (`grep -rn` über `internal/` und `test/`: nur der Hilfsrumpf `anmScramRundeSalzText`, den niemand mit einem Salztext aufruft). Das Verhalten stimmt am Binary (Erfolgsfälle oben), die Zusage ist ungeprüft | **abgelehnt für diese Zusage, V-162** |

**Tabelle der Unit-Mutanten** (alle in frischer Kopie, gofmt-sauber, vet ohne Befund; „rot“ heißt: ein genannter Test schlug mit seiner Meldung fehl, nicht Build, Vet oder gofmt).

| Gruppe | Mutanten (Ersetzung) | roter Test |
|---|---|---|
| Geheimnis | `Format` an `zugang` gibt `z.passwort` aus; `Format` an `zugang` entfernt; `Format` an `scramAustausch` gibt `s.passwort` aus; `Format` an `Passwort` (Postgres-Adapter) gibt den Wert aus; `Format` an `Passwort` (CLI) gibt den Wert aus; erwartete Signatur in der Meldung; Nachrichtentext des Abschlusses in der Meldung | `TestAnmeldungStrukturenOhnePasswort/{zugang,anmeldung_Zeiger,anmeldung_Wert,scramAustausch_*}`; `TestAnmeldungOhneGeheimnis` (letzte Schleife, `/SCRAM_falsche_Signatur`, `/SCRAM_Abschluss_ohne_v=`); `TestPlayOptionenOhnePasswort` |
| F-593 | Prüfung der Servernonce entfernt; untere Grenze 0x20; obere 0x7F; obere 0x7D; führende Null erlaubt; Prüfung nur auf `00`; leeres Salz abgelehnt; Zeilenumbruch in Base64 erlaubt; `RawStdEncoding` | `TestScramServerErste/{r_mit_Leerzeichen,…_Byte_0x7F,Iterationen_…_führende_Null,s_mit_Zeilenumbruch,s_ohne_Auffüllung}`, `TestAnmeldungScramLeeresSalz` |
| **F-593 Restbits** | **`.Strict()` in `dekodiereBase64`** | **keiner: grün (V-162)** |
| MD5 | Reihenfolge `benutzer + passwort`; Salz vor dem Hash; Präfix `md5` weg; Benutzer fest `postgres`; `len(rest) < 4` | `TestAnmeldungMD5` (Vektoren aus `hashlib`), `/ohne_Benutzer`, `TestAnmeldungAnforderungNichtLesbar/MD5_mit_fünf_Byte` |
| Quellen und Reichen | Variable nicht gelesen; Platzhalter ignoriert; Variable vor Platzhalter; Variable nicht bei `host:port`; Bootstrap reicht das Passwort nicht (rot nach 0,01 s); `Verbinde` reicht `zugang{}` | `TestPlayPasswort/*`, `TestPlayVerbindung`, `TestRunPlayPasswort/*`, `TestRunPlayErstesSignalInAnmeldung`, `TestRunPlayZweitesSignalInAnmeldung`, `TestAnmeldungVerbinde`, `TestAnmeldungNonce`, `TestAnmeldungSignal` |
| Fehlendes Passwort | Prüfung entfernt; `CodeUpstream` statt `CodeLogin` | `TestAnmeldungOhnePasswort/{Klartext,MD5,SCRAM}`, `TestAnmeldungVerbinde`, `TestEinspielAufbauFehler` |
| SCRAM-Rechnung | Label `Server Key` → `Client Key`; Beweis nur `ClientSignature`; Iterationen + 1; Nonce 16 Byte; Nonce fest; `c=eSws`; Kopf `y,,`; Benutzer `n=u` | `TestScramBeweisRFC7677`, `TestAnmeldungScram`, `TestAnmeldungScramLeeresSalz`, `TestAnmeldungNonce` |
| SCRAM-Nachrichten | Signatur nicht geprüft; Servernonce-Präfix nicht geprüft; gleich lang erlaubt; Obergrenze + 1; Iteration 0 erlaubt; Überlaufprüfung nach der Schleife; `len(teile) < 3`; Zusatzattribut hinter `v=` gekürzt; Code 12 vor 11 erlaubt; Code 11 nach der Antwort erlaubt; Abschluss beendet den Austausch nicht; Austausch gilt nie als laufend | `TestAnmeldungScramFehler/*`, `TestScramServerErste/*` (u. a. `Iterationen_2_hoch_64_plus_5`), `TestAnmeldungScramUnvorgesehen`, `TestAnmeldungNachScramAbschluss`, `TestAnmeldungOhneGeheimnis/SCRAM_falsche_Signatur` |
| Verfahren | Name als Präfix; `EqualFold`; nur das erste Verfahren der Liste; Klartext mit Rest erlaubt; Liste ohne Ende lesbar; `TrimSpace` auf das Klartext-Passwort | `TestAnmeldungNichtUnterstuetzt/*`, `TestAnmeldungScram/SCRAM_nach_PLUS`, `TestAnmeldungAnforderungNichtLesbar/*`, `TestAnmeldungKlartext` |
| `beantwortet`-Reihenfolgen | `beantwortet` vor dem Verfahren (Code 7, 9, unbekannt nach einer Antwort); Lesbarkeit vor `beantwortet` (Klartext mit Rest nach einer Antwort); `beantwortet` ohne Wirkung; Fortsetzung 8 nicht erkannt; Nachricht nach `AuthenticationOk` nicht erkannt | `TestAnmeldungWeitereAnforderungArt/*` (sechs und zwei Fälle), `TestAnmeldungWeitereAnforderung`, `TestAnmeldungNachScramAbschluss`, `TestAnmeldungFortsetzungOhneAustausch`, `TestEinspielAufbauFehler` |
| Senden und Signal | gescheitertes Senden verworfen; `context.AfterFunc` durch Attrappe | `TestAnmeldungSendenScheitert`, `TestAnmeldungSignal/*`, `TestEinspielAufbauAbgebrochen` |

Die drei E2E-Mutanten stehen bei Punkt 1. Zwei meiner Mutanten (`nonce-fest`, `scram-nur-erstes-verfahren`) bauten im ersten Anlauf wegen eines unbenutzten Imports nicht und sind wiederholt (jetzt rot in `TestAnmeldungNonce` und `TestAnmeldungScram/SCRAM_nach_PLUS`).

### Punkt 3: `make gates` grün. **Bestätigt durch eigenen Lauf.**

Eigener Lauf im Arbeitsbaum am Stand `556f8aa`, sauber vor und nach dem Lauf: Exit 0. Darin `baseline-verify: v6.18.0 OK — 54 Dateien`, `d-check: 494 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen`, `abdeckung-gegenprobe: gruen`, `commit-msg-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `run-integration-tests: gruen` mit `--- PASS` für `TestE2EPlayAnmeldung` (6 Fälle), `TestE2EPlayAnmeldungOhneVerlangen` und `TestE2EPlayAnmeldungFehler` (9 Fälle); kein `FAIL` im Log.

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`3cc01e4`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** noch offen; §6 trägt bei beiden Risiken „offen bis Closure“. Das ist die Reihenfolge, kein Befund. Für die Closure aus meiner Sicht: Risiko 1 (SCRAM nur gegen einen Testserver geprüft, der denselben Fehler rechnen könnte) ist nicht eingetreten: Der echte Server (SCRAM, MD5, Klartext) hat den Beweis in allen Läufen angenommen, und der Fake-Server rechnet mit `hashlib` unabhängig vom Code, `proof-check` war jedes Mal gleich. Risiko 2 (nur Version 17) bleibt weiter offen: Ich habe ebenfalls nur gegen Version 17 gefahren.
- **Lerneintrag:** nach V-162 gehört die Klasse *Test im Plan genannt, im Repo nicht vorhanden* hinein (siehe Befunde).

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-591 `Format` an `zugang`, `anmeldung`, `scramAustausch` | `TestAnmeldungStrukturenOhnePasswort` formatiert Wert, Zeiger, Liste und Abbildung mit sieben Verben. Mutanten (Format entfernt, Format gibt den Wert aus, an `zugang` und an `scramAustausch`): **rot**. Die Klasse „Passwort in einem unexportierten Feld“ ist im Postgres-Adapter abgedeckt; die CLI hält eine ähnliche Struktur (`ziel.passwort`), die niemand formatiert: V-163 | behoben; V-163 |
| F-592 §6 „Der Code weicht ab“ | `grep` nach „weicht ab“ im Plan: nur die Erledigungszeile in §7. §6 sagt zur führenden Null „Umgesetzt mit `e74e2fe`“; Code und Plan stimmen (`i=0004096`, `i=01`, `i=010` sind am Binary und im Mutanten-Lauf `PGR-E4005`) | behoben |
| F-593 Servernonce `printable`, Restbits | Servernonce: **behoben**, Spezifikation (*SCRAM-Austausch*) und Code stimmen, fünf Mutanten rot, Binary lehnt Leerzeichen ab. Restbits: Spezifikation und Code stimmen (Binary nimmt `YR==` und ein `v=` mit gesetzten Restbits an), **die zwei genannten Tests fehlen** | Servernonce behoben; Restbits nicht behoben, V-162 |
| F-594 Konjunktiv im E2E-Kommentar | Der Kommentar ist unverändert, der Plan nennt ihn als Kopplung an das Verhalten des Servers. E2E-Mutant E3 (Passwort bei `AuthenticationOk`) ist rot, die Kopplung hält | eingeordnet, ohne Aktion |
| F-595 `pg_reload_conf()` asynchron | Unverändert; in allen meinen Läufen nicht aufgetreten (ich habe vor der ersten Anmeldung zwei Sekunden gewartet, der Runner nicht). Der Lauf von `make test-integration` im Gate war beim ersten Mal grün; in rund 120 Anmeldungen meiner Binary-Läufe gegen den Server trat es nicht auf | eingeordnet, ohne Aktion |
| F-596 Rot erst nach 120 s | Mutant „Bootstrap reicht das Passwort nicht“: `TestRunPlayPasswort/*`, `TestRunPlayErstesSignalInAnmeldung` und `TestRunPlayZweitesSignalInAnmeldung` sind nach **0,00 bis 0,01 s** rot, das Paket in 3,1 s. Die Frist von 30 s bleibt als Literal | behoben |

---

## 3. Hexagon und Architektur-Gate

- `make a-check` im eigenen `make gates`: `gesamt: 0 Befund(e)`; `a-check-negativ: gruen`.
- `git diff --stat f9bc4d0..556f8aa -- go.mod go.sum .a-check.yml .golangci.yml Makefile harness cmd` ist leer; der Diff berührt nur Upstream-Adapter, CLI-Adapter, Bootstrap (`bootstrap.go` +4 −2, `play_test.go` +1 −1: Verdrahtung und Kommentar), Tests, den Runner der Integrationstests, Abdeckungstabellen, Spezifikation und Plan. Kein Code im Kern (`internal/hexagon`) und im PGWire-Adapter.
- `anmeldung.go` importiert die Standardbibliothek und `pgproto3` (Nachrichten), das [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) im Upstream-Adapter erlaubt; keine neue Bibliothek.
- [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md): Verfahren Klartext, MD5, SCRAM-SHA-256 am Binary gesehen; Quellen Platzhalter, dann Variable; keine Option für das Passwort (`--password` ist am Binary `PGR-E2001`); kein Passwort in der Aufzeichnung (vier Aufzeichnungen ohne `GEHEIM`, gegen `trust` aufgenommen).
- `SPEC-033` (Logs erzeugen keine zusätzlichen Geheimnisse): `--log-level debug` bei Erfolg und Fehler zeigt zwei Zeilen und kein Geheimnis.
- `-race` über `./internal/adapters/driven/postgres`, `./internal/adapters/driving/cli`, `./internal/bootstrap` (Hinweis, kein Gate; Go 1.27.1 im Image `golang:1.27`): alle drei `ok`, kein Befund.

---

## 4. Plan gegen Code

- **§1:** Geliefert sind die Passwortquellen, die drei Verfahren, `PGR-E4005` der Anmeldung, die Bindung an den Vertrag der Ports (jeder Fehler der Anmeldung trägt seinen Code, am Binary gesehen: `PGR-E4005` und `PGR-E4002`, kein Fehler ohne Code, Exit-Code 4 in beiden Fällen) und der Kommentar zur Kopplung im Bootstrap. Die Abgrenzungen halten: TLS und `--upstream-tls` kommen weiter nicht vor; `record` ist unberührt; Kerberos, GSSAPI, SSPI sind `PGR-E4005` ohne Senden. Benutzerhandbuch und `README.md` sind nicht im Diff (Doku-Folge-Slice).
- **§3:** Jede geänderte Datei steht dort. `internal/bootstrap/play_test.go` (+2 −2) und `internal/adapters/driving/cli/export_test.go` (+1) decken die Zeilen *Tests* und *Ablösung des Zwischenstands*; `go.mod` bleibt unverändert, wie §6 *Bibliothek für SCRAM* sagt. Der Bootstrap ist reine Verdrahtung (`+6 −3`), keine Rückführung nach §4.
- **§6:** Jede Randform, die ich am Binary gefahren habe, folgt ihrem Ort in `LH-FA-20.a`: *Passwort* (Quellen, Vorrang, leer, unverändert, kein Geheimnis), *Verfahren* (Wahl `SCRAM-SHA-256`, Rest nicht lesbar, zweite Anforderung), *SCRAM-Austausch* (Reihenfolge der Attribute, Nonce, Salz leer, Iterationen 1 bis 10 000 000 ohne führende Null, `e=`, Signatur, Restbits), *Aufbau*, *Anmelde-Nachrichten*, *Abbruch im Aufbau* (kein `Terminate` nach einem Fehler und nach dem zweiten Signal). Die akzeptierten Negative (Fehler der Schlüsselableitung, kein `bytes.Equal`-Test, keine Frist) habe ich nicht gefahren; die Spezifikation sagt sie als Grenze, nicht als geprüft.
- **§7:** Die Belege nennen Größe, Schichten, Weg der Mutanten, die Tabelle, die Läufe und die Nacharbeit. Falsch ist in der Nacharbeit zu F-593 die Zeile zu `TestAnmeldungScramSalzNichtKanonisch` und `TestAnmeldungScramSignaturNichtKanonisch` (V-162). Die Kopfzahl „84 Unit-Mutanten … 83 rot, 1 grün“ ist nicht nachgezogen: Die Nacharbeit zum Review hat Mutanten (F-591, F-593) ergänzt, deren Zahl in §7 nicht steht; die Tabellen der Nacharbeit nennen sie einzeln.

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-162 | HIGH | Die Zusage „`s=` und `v=` mit gesetzten Restbits (`AB==`, `YR==`) werden angenommen“ (`LH-FA-20.a` *SCRAM-Austausch*, F-593) hat keinen Test. §6 und die Tabelle der Nacharbeit in §7 nennen `TestAnmeldungScramSalzNichtKanonisch` und `TestAnmeldungScramSignaturNichtKanonisch` und behaupten, `.Strict()` in `dekodiereBase64` mache beide rot; keiner der beiden Tests existiert (`grep` über `internal/` und `test/`: 0 Treffer; im Commit `b263135` nur der Hilfsrumpf `anmScramRundeSalzText` ohne Aufrufer mit einem Salztext). Der Mutant bleibt grün (`go test` der drei Pakete `ok`). Das Verhalten stimmt am Binary (Fake-Server: Salz `YR==` und ein `v=` mit gesetzten Restbits gelingen). Damit ist DoD-Punkt 2 für diese Zusage nicht erfüllt (`AGENTS.md` §3.10) und der Plan sagt zu, was kein Test prüft (§3.11). | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-anmeldung.md` · „`TestAnmeldungScramSalzNichtKanonisch` und“; `internal/adapters/driven/postgres/einspielen_anmeldung_test.go` · „func anmScramRundeSalzText“ | Implementer: die beiden Tests schreiben (ein Austausch mit `s=YR==` für Salz `a` und einer mit `v=` im nicht-kanonischen Base64 desselben Werts, je mit dem Beweis aus `hashlib`), den Mutanten `.Strict()` rot sehen und die Tabelle in §7 belegt nachziehen; bis dahin den Satz in §6 und §7 nicht als geprüft führen |
| V-163 | LOW | `cli.ziel` hält das eingesetzte Passwort in einem unexportierten `string`-Feld (`passwort`), also derselben Klasse wie F-591: `%+v` auf einem `ziel` gäbe es aus. Nichts im Diff formatiert einen Wert dieser Struktur (`grep` nach Formatierungs-Verben in `cli` und `bootstrap`: nur Kommentar), und der Plan begrenzt die Klassenaussage auf den Upstream-Adapter. Latent, kein Verhalten am Binary. | `internal/adapters/driving/cli/verbindung.go` · „passwort                     string“ | Planner: ohne Aktion annehmen oder in die Klasse von F-591 aufnehmen; der Implementer entscheidet bei der nächsten Berührung der Datei, ob `ziel` ein `Format` bekommt |
| V-164 | INFO | Mehrere Mutanten sind im Paket Postgres-Adapter erst nach 24 bis 84 s rot, weil Tests mit dem Fake-Server (`TestAnmeldungVerbinde`, `TestEinspielAufbauFehler`, `TestAnmeldungScramFehler/…`, `TestAnmeldungOhneGeheimnis/…`) bis zur Frist von 20 s warten; die direkten Einheitstests desselben Mutanten (`TestAnmeldungOhnePasswort`, `TestScramServerErste`) sind nach 0,00 s rot. Aus dem richtigen Grund rot ist jeder; nur die Rückmeldung ist langsam. Dieselbe Klasse wie F-596 im Bootstrap. | `internal/adapters/driven/postgres/einspielen_anmeldung_test.go` · „TestAnmeldungVerbinde“ | ohne erwartete Aktion; Planner nimmt die Klasse *Rot erst nach Frist* in die Closure auf |

Die Falsifikationsversuche, die nichts fanden: SCRAM-Rechnung (Client-first, Beweis, Serversignatur) gegen die Rechnung des Fake-Servers mit `hashlib`/`hmac` in 19 Austauschen gleich (Salz leer, nicht-kanonisch, 10 000 000 Iterationen, Fehlerfälle nach dem Beweis eingeschlossen); MD5 gegen `hashlib` mit Salz `ff 00 10 20`, Benutzer der Startup-Daten, gleich; Klartext byteweise gleich; Nonce 18 Byte in jedem Austausch; kein `GEHEIM` in 296 Ausgabedateien; kein `Terminate` nach einem Fehler der Anmeldung; erstes Signal lässt den Aufbau zu Ende laufen, zweites schließt ohne `Terminate`.

Befunde: 1 HIGH, 0 MEDIUM, 1 LOW, 1 INFO. Die Verhaltens-DoD (Punkte 1 und 3, Punkt 2 im Verhalten) sind bestätigt; Punkt 2 ist im Beleg einer Zusage (V-162) nicht erfüllt.

**Finding-Klassen dieses Laufs:** Test im Plan genannt, im Repo nicht vorhanden · Passwort in unexportiertem Feld einer Struktur · Rot erst nach Frist

**Merge-blockierend / Closure-blockierend:** ja, wegen V-162. Die Korrektur ist klein (zwei Tests, ein Lauf des Mutanten, Nachzug von §7).

**Übergabe an den Planner:** dieser Bericht mit V-162 (an den Implementer zurück), V-163 und V-164; dazu für die Closure die Ausgänge der zwei Risiken aus Abschnitt 1 (*Übrige DoD-Punkte*) und die Klassen oben.

**Aufgeräumt:** Container `ver-anm-pg` und alle `ver-anm-play-*`, `ver-anm-fk-*`, `ver-anm-sg-*`; Netz `ver-anm-net`; Images `ver-anm-runtime:1`, `ver-anm-test:1` und `ver-anm-int-*`; die Verzeichnisse `ver-anm-*` im Scratchpad. `docker ps -a`, `docker network ls` und `docker volume ls` mit `ver-anm` und `pgr-it` danach leer; das Image `pgwire-recorder:test` des Arbeitsbaums blieb stehen.
