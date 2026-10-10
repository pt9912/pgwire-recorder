# Review-Report: slice-v1-abschluss-einspielen-anmeldung — 2026-10-10

**Review-Art:** Code-Review gegen Plan §1, §3, §6 und §8, gegen `LH-FA-20.a` (*Passwort*, *Verfahren*, *SCRAM-Austausch*, *Anmelde-Nachrichten*, *Aufbau*) und `LH-FA-17.a`, gegen [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und gegen `AGENTS.md` §3.7, §3.9 bis §3.13.

**Gegenstand:** Diff `f9bc4d0..265b184`: Code-Commit `c8b7f44`, Belege `238c78b`, Nacharbeit zu den Randform-Fragen `e74e2fe`, Stand der Gates `265b184` (19 Dateien, +2225 −81). Neu `internal/adapters/driven/postgres/anmeldung.go`; geändert `einspielen.go`, CLI (`cli.go`, `upstream.go`, `verbindung.go`), `internal/bootstrap/bootstrap.go`, `tools/test/run-integration-tests.sh`, Tests, Abdeckungstabellen, Spezifikation (Architect-Commit `482951f`), Plan §3, §6, §7.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-anmeldung.md`, ganz gelesen mit §7 *Belege des Implementers*; den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-20.a` *Passwort*, *Verfahren*, *SCRAM-Austausch*, *Aufbau*; `SPEC-038`; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten).
- [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md) ganz, [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) und [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) (Abschnitt *Entscheidung*).
- `AGENTS.md` §3.7, §3.9 bis §3.13 (dort *Nachzählen beim Eintragen* und die Bootstrap-Regel).
- Der Nehmer `slice-v1-abschluss-einspielen-anmeldung-doku` (Kopf, §1, DoD): der Diff trägt keine neue Sendung hinein.
- Unverändert, aber vom Diff vorausgesetzt: `Verbinde` und `aufbau` in `einspielen.go`, `model.Meldungen`, Port `Einspielziel`.
- Nummern dieses Laufs ab F-591 (vorheriger Report: `slice-v1-abschluss-einspielen-laufsteuerung`, bis F-590).

**Ausgeführte Läufe:**

- **Unabhängige Vektoren.** Python (`hashlib`, `hmac`, `base64`) rechnet zwei vollständige SCRAM-Austausche (Passwort `pässwörd $%41`, Salz `00 01 "salt" ff`, 4096 Iterationen, Nonce 0..17; Passwort `x`, leeres Salz, 1 Iteration, Nonce 100..117): erste Nachricht, Antwort mit Beweis, Serversignatur; dazu zwei MD5-Antworten (`pw`/`bob`, `ä p`/`u` mit Salz `ff 00 10 20`). Ein eigener Test in einer Kopie ruft `neuerScramAustausch`, `schritt` und `md5Antwort` und stimmt byteweise überein. Der RFC-7677-Vektor des Implementers (`TestScramBeweisRFC7677`) deckt den Rechenkern zusätzlich.
- **Mutanten.** Je Mutant eine frische Kopie von `internal/` per `cp -r` ohne `-p` unter dem Scratch-Pfad `rev-anm-m`, genau eine Ersetzung, gofmt und vet im Image `pgwire-recorder:test` mit Bind-Mount, dann `go test` des Pakets. 38 Mutanten, alle rot in einem Test der Unit-Pakete (Liste unten); drei Mutanten, die nicht baute (`passwort-trim-env`, erste Fassungen von `sig-nicht-geprueft` und `nonce-fest`), sind nicht gezählt oder gofmt-sauber wiederholt. Bei `salz-leer-ablehnen` war die Kopie nicht gofmt-sauber; der Rot-Befund kommt aus `TestAnmeldungScramLeeresSalz`, nicht vom Format. Dazu ein E2E-Mutant in einer Kopie des Baums (`git archive`): `AuthenticationOk` sendet ein `PasswordMessage` mit dem Passwort, wenn eines da ist, `make test-integration` rot in `TestE2EPlayAnmeldung/*` (6 Fälle) und `TestE2EPlayAnmeldungOhneVerlangen`, Exit-Code 4.
- **Mutanten im Einzelnen (alle rot).** MD5: Reihenfolge `benutzer+passwort`, Präfix `md5` weg, Salz nicht angehängt, `len(rest) < 4`. Klartext: `TrimSpace`. SCRAM: Beweis nur `ClientSignature`, Label `Server Key` durch `Client Key`, Iterationen +1, Signatur nicht geprüft, `bytes.Equal` mit `&& false`, Nonce 16 Byte, Nonce fest statt `crypto/rand`, Servernonce gleich lang angenommen, Präfixprüfung weg, Attribute `< 3` statt `!= 3`, `c=eSws`, Kopf `y,,`, Leersalz abgelehnt, Zeilenumbruch in Base64 zugelassen, `v=` mit weiterem Attribut gekürzt, Austausch läuft nie. Iterationen: Obergrenze doppelt, Obergrenze exklusiv, führende Null zugelassen. Anforderungen: `beantwortet` ohne Wirkung, `beantwortet` vor dem Verfahren, SASL-Liste ohne Ende lesbar, `EqualFold`, `fehlendes Passwort` mit `CodeUpstream`, Prüfung des Passworts in `sasl` entfernt, unverlangtes Passwort bei `AuthenticationOk`. Geheimnis: `Format` gibt den Wert aus, Serversignatur in der Meldung. Quellen: Variable vor Platzhalter, Platzhalter-Teil immer, Bootstrap reicht das Passwort nicht, `Verbinde` reicht `zugang{}`.
- **Integration.** `make test-integration` zweimal am unveränderten Baum: beide grün (`TestE2EPlayAnmeldung` 6 Fälle, `TestE2EPlayAnmeldungOhneVerlangen`, `TestE2EPlayAnmeldungFehler` 9 Fälle), kein Flackern; `GEHEIM` kommt in keiner Ausgabe beider Läufe vor (`grep -c` 0); danach `docker ps -a`, `docker network ls` und `docker volume ls` mit `pgr-it` leer.
- **Gates.** `make a-check`: `gesamt: 0 Befund(e)` (der Hinweis auf 15 Dateien unter `test/integration` ohne Schicht ist Bestand). `go.mod`, `go.sum`, `.a-check.yml`, `.golangci.yml`, `Makefile` und `harness/` unverändert (`git diff --stat f9bc4d0..265b184` leer). `make abdeckung-check` und `make lint`: grün.
- **`-race`.** Das deps-Image (alpine) hat kein `gcc`; der Lauf mit `CGO_ENABLED=1` dort scheitert am Build. Stattdessen `golang:1.27` (Debian, vorhanden), Quellbaum schreibgeschützt, Modulcache aus dem Test-Image, `--network none`: `go test -race -count=1` über den Postgres-Adapter, den CLI-Adapter und den Bootstrap, alle drei `ok`, ohne Befund (Hinweis, kein Gate).
- **Eigene Probe in einer Kopie.** `%+v` auf `zugang`, `anmeldung` und `scramAustausch` gibt das Passwort aus (F-591); `dekodiereBase64("AB==")` ist gültig (F-593).

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v<X.Y.Z>` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
— Tag einsetzen, denn diese Zeile wandert in den eingefrorenen Report.

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund — je ein Finding pro
     Zeile, 1-2 Sätze. Absichtlich (noch) nicht gate-geprüft: d-check
     `structure` (`table.column[].cell-max-chars`) könnte die Spalten
     `Befund`/`Klasse` zellenlängen-prüfen, sobald genug reale Reports
     zeigen, welche Grenze die gelebte Praxis trägt — verfrüht gesetzt,
     bricht sie am ersten gründlichen Befund. -->

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-591 | MEDIUM | Der Kommentar sagt, jede Formatierung eines Werts, der ein Passwort trägt, verrate es nicht. Das gilt nur für den Typ `Passwort` in einem exportierten Feld; `%+v` auf `zugang`, `anmeldung` und `scramAustausch` (dort `passwort string`) gibt `GEHEIM` aus (in einer Kopie ausgeführt). Ein Test formatiert nur `Einspielziel`. Szenario: Wer den Zustand der Anmeldung für ein Debug-Log mit `%+v` formatiert, gibt das Passwort aus im Glauben an die Zusage. | `AGENTS.md` §3.11 (Kommentar sagt nur zu, was ein Test prüft); `LH-FA-20.a` *Passwort* | `internal/adapters/driven/postgres/anmeldung.go` · „%#v, %s oder %q auf einem Wert, der ein Passwort trägt, verrät es nicht" | ja — ein Test, der `zugang{}`, `anmeldung` und `scramAustausch` mit `%+v` formatiert | Zusage im Kommentar weiter als Prüfung |
| F-592 | MEDIUM | §6 nennt unter *Randformen aus der Rückgabe des Implementers* weiter „**Der Code weicht ab** (`leseIterationen` nimmt sie an)" und „**Auftrag:** Code ändern, Test in `TestScramServerErste` …", obwohl `e74e2fe` genau das umgesetzt hat; der Commit zieht nur §7 nach, nicht §6. Der Plan beschreibt einen Code-Stand, den es nicht mehr gibt. | `AGENTS.md` §3.9 (der Slice-Plan folgt jeder Korrektur) | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-anmeldung.md` · „**Der Code weicht ab** (`leseIterationen` nimmt sie an)" | nein — Urteil; `make kopf-check` prüft nur den Kopf | Plan folgt der Korrektur nicht |
| F-593 | LOW | Zwei Randformen des Austauschs entscheidet der Code ohne Eintrag in §6 und ohne Satz der Spezifikation: nicht-kanonisches Base64 (`s=` und `v=` mit gesetzten Restbits, etwa `AB==`) ist „gültig", und die Servernonce darf jedes Zeichen außer `,` tragen (RFC 5802 verlangt druckbare ASCII-Zeichen). Ein Schaden ist nicht zu erzählen: Die Signatur bleibt bytegenau geprüft. | `AGENTS.md` §3.12 (Randformen vor dem Code entschieden) | `internal/adapters/driven/postgres/anmeldung.go` · „b, err := base64.StdEncoding.DecodeString(text)" | ja — Test mit `v=` aus nicht-kanonischem Base64 desselben Werts | Randform im Code entschieden, nicht in §6 |
| F-594 | INFO | Der Kommentar „Ein unverlangtes Passwort wäre für den Server eine ungültige Nachricht, und das Einspielen gelänge nicht" steht im Konjunktiv über den Fehlerfall. Er trägt die Kopplung des Tests an das Verhalten des Servers; ich lese ihn als Kopplung und nicht als verworfene Alternative, der Satz beschreibt aber, was nicht eintritt. Der E2E-Mutant (Passwort bei `AuthenticationOk`) ist rot, die Kopplung hält. | `AGENTS.md` §3.7 | `test/integration/play_anmeldung_e2e_test.go` · „unverlangtes Passwort wäre für den Server eine ungültige Nachricht" | nein | Konjunktiv im Kommentar |
| F-595 | INFO | Der Runner legt die Zeilen von `pg_hba.conf` an und ruft `pg_reload_conf()`, das asynchron wirkt; der Testlauf startet danach in einem neuen Container. Bis zur Wirkung gilt `trust` für die drei Benutzer, ein Fehlerfall, der `28P01` erwartet, schlüge in diesem Fenster fehl, weil der Server ohne Passwort anmeldet. In zwei Läufen trat es nicht auf. | Maintainability | `tools/test/run-integration-tests.sh` · „SELECT pg_reload_conf();" | nein | Asynchrones Neuladen ohne Wartepunkt |
| F-596 | INFO | Rot ist der Mutant „Bootstrap reicht das Passwort nicht" erst nach 120 s (viermal die 30-s-Frist des Fake-Servers, Meldung „der Fake-Server liest binnen 30 s kein Passwort"); die Meldung ist die richtige, die Rückmeldung langsam. | Maintainability | `internal/bootstrap/play_anmeldung_test.go` · „der Fake-Server liest binnen 30 s kein Passwort" | ja | Rot erst nach Frist |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| SCRAM-Rechnung (Client-first `n,,n=,r=`, Auth-Message, PBKDF2-SHA-256, ClientKey/StoredKey/ClientProof, ServerSignature) | geprüft gegen eigene Python-Vektoren (zwei Austausche, ein leeres Salz) und den RFC-7677-Vektor, ohne Befund |
| Serversignatur-Vergleich, Zufall, Nonce | `hmac.Equal` (konstante Zeit; ein Mutant `bytes.Equal` ist im Ergebnis äquivalent, die Eigenschaft hat keinen Test und die Spezifikation sagt sie nicht zu); `crypto/rand`, 18 Byte, Base64; geprüft, ohne Befund |
| MD5-Formel und Klartext | geprüft gegen Python-Vektoren (`md5` + MD5(MD5(Passwort+Benutzer)+Salz), Benutzer aus den Startup-Daten), ohne Befund |
| Obergrenze, führende Null, Attributreihenfolge, Servernonce (Präfix, Länge), Zeilenumbruch | geprüft mit Mutanten, ohne Befund (offen nur F-593) |
| Operationen gegen §6 (Lesen, Senden, Schließen, Fehlercode, Signal, Frist) | `PGR-E4002` für Senden und Verbindungsende, `PGR-E4005` für Anmeldefehler, kein `Terminate`, `context.AfterFunc` schließt; keine weitere still entschiedene Randform außer F-593 |
| Passwort in Ausgaben | Optionen von `play` und `Einspielziel`: Formatierung maskiert (Tests, Mutant `Format` rot); `Error()` und Ursachenketten der Anmeldung nennen weder Passwort noch Signatur (Mutant rot); E2E: `GEHEIM` in keiner Ausgabe, auch mit `--log-level debug`; offen nur F-591 |
| `internal/adapters/driving/cli` | Quellen und Vorrang (Platzhalter vor Variable, leere Variable gilt als nicht gesetzt, kein Kürzen), Hilfetext, Abdeckungszeilen geprüft; die Zusage „keine Option für das Passwort" ([ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md)) halten die bestehenden Tests der Optionstabelle (Probe: eine Option `password` in der Tabelle macht `TestLeserAlleOptionen` und `TestLeserOptionen` rot), ohne Befund |
| `internal/bootstrap` | nur Reichen des Passworts und der Kommentar zur Kopplung (+4 −2), keine Logik, keine Rückführung nach §4 |
| `go.mod`, `go.sum`, `.a-check.yml`, `.golangci.yml`, `harness/` | unverändert; nur Standardbibliothek, keine neue Bibliothek, keine Ausnahme im Lint-Profil, kein `//nolint` |
| `tools/test/run-integration-tests.sh` | die drei Benutzer (`scram-sha-256`, `md5`, `password`), Passwörter nur als Umgebung des Testcontainers, nichts außerhalb des Containers und des benannten Netzes berührt; zwei Läufe grün (offen nur F-595) |
| `docs/user/abdeckung-*.md` | `make abdeckung-check` grün, Zeilen entsprechen den Deklarationen |
| `spec/spezifikation.md` | Änderung `482951f` ist Architect-Arbeit, im Diff nur gelesen: kein ADR-, Slice- oder Wellenbezug, ADR-Schärfung unberührt |
| Nehmer `slice-v1-abschluss-einspielen-anmeldung-doku` | keine neue Sendung im Diff; Adresse und Nachzählen unverändert, ohne Befund |
| Akzeptierte Negative (PBKDF2-Fehler nur im FIPS-Modus, nicht unterbrechbare Berechnung) | nachvollziehbar: Spezifikation nennt beide als Grenze und sagt keinen Test zu; die Berechnung ist durch 10 000 000 begrenzt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zusage im Kommentar weiter als Prüfung · Plan folgt der Korrektur nicht · Randform im Code entschieden, nicht in §6 · Konjunktiv im Kommentar · Asynchrones Neuladen ohne Wartepunkt · Rot erst nach Frist

## Verdikt

**Merge-blockierend:** ja — wegen F-591 und F-592 (beide MEDIUM, beide kleine Änderungen an Kommentar, Test und Plan-Text; kein Befund am Protokoll, an der Rechnung oder am Geheimnis im Betrieb). Kein HIGH.

**Übergabe je Finding:** F-591 an den Implementer (Kommentar oder Test; Artefakt: dieser Report, Zeile F-591). F-592 an den Implementer (§6 des Slice-Plans nachziehen, `AGENTS.md` §3.9). F-593 an den Architect (Entscheidung der beiden Randformen in `LH-FA-20.a` *SCRAM-Austausch*; der Implementer entscheidet sie nicht, `AGENTS.md` §3.12). F-594 bis F-596 ohne erwartete Aktion; Annahme oder Begründung genügt, die Sequenz über den Architect wäre Overkill (isolierte INFO). Kein Rollen-Widerspruch bisher.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
