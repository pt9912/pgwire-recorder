# Review-Report: slice-v1-abschluss-protokollrand-doku — 2026-10-15

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). Gegenstand ist Dokumentation (Handbuch, README); geprüft wird gegen §1 und §6 des Plans, `AGENTS.md` §3.7, §3.9, §3.11, §3.13, §3.14 und das gebaute Binary.

**Gegenstand:** slice-v1-abschluss-protokollrand-doku, Diff `cf08408..2073225` (Architect `9ce025c`, Implementer `7a0be02` Handbuch und README, `2073225` Plan §3/§7)

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec0
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

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

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan slice-v1-abschluss-protokollrand-doku (§1 Sendung, §6 R1 bis R8, P1 bis P7, §7 Belege des Implementers, §8 Zählung); Geber slice-v1-abschluss-protokollrand (in `done/`)
- aktive ADRs: keine im Diff genannt; Spezifikationsstellen `spec/spezifikation.md` LH-FA-05.e (Protokollrand, Randformen, Zählung der Verbindungen), LH-FA-05.b, LH-FA-23.a
- berührte `LH-*`-IDs: LH-FA-05, LH-FA-12, LH-FA-23, LH-RB-02
- `AGENTS.md` (Hard Rules, insbesondere §3.7, §3.9, §3.11, §3.13, §3.14)

**Eigene Läufe dieses Reviews** (Binary aus `make build`, Image `pgwire-recorder:dev`, **als Container** gegen `postgres:17-alpine` (Digest aus `harness/mk/integration.mk`) in eigenem Docker-Netz; Clients: Python-Rohnachrichten, `psql` 17 und 18; Container und Netz danach entfernt, Arbeitsbaum sauber):

- `record` mit TLS-Pflicht: `_pq_.x` ergibt `PGR-E6001`, dieselbe Startnachricht ohne `_pq_.` `PGR-E6003`; Exit-Code 6. `replay` mit TLS-Pflicht: dasselbe Paar; Rohverbindung über TLS mit Benutzer `gibtesnicht`, Datenbank `nix` erhält `AuthenticationOk`, `Z`.
- `record --record-empty-sessions`: Version 3.2 (`PGR-E6002`), `_pq_.x`, Rumpf Länge 9, Nachricht `p`, `CopyData`, `COPY t FROM STDIN`, `COPY t TO STDOUT` (je `PGR-E6001`), danach reguläre Verbindung: Datei hat genau 1 Sitzung, Exit-Code 6. Rumpf ohne abschließendes Nullbyte, Länge 8 und Länge 9 mit `x` je `PGR-E6001` ohne Upstream.
- `record` `LISTEN`/`NOTIFY`: allein und nach der letzten Anfrage Exit-Code 0, beide Sitzungen in der Datei; in einer Anfrage, in zwei Anfragen derselben Verbindung, auf zwei Verbindungen `PGR-E6001`, Exit-Code 6, zuhörende Verbindung nicht in der Datei, sendende schon.
- `replay` mit `first-request` und mit `connection`: Version 3.2, `_pq_.x`, Rumpf Länge 9, danach zwei reguläre Verbindungen erhalten Sitzung 1 und 2, die dritte `PGR-E5003` („Verbindung 3“); `p`, `CopyData` (und `CopyDone`) nach dem Aufbau `PGR-E6001`.
- `psql` 17 gegen `replay`, mit und ohne TLS: `require_auth=password|md5|scram-sha-256` bricht beim Client ab; `replay` Exit-Code 0, im Log keine Meldung `PGR-`; `psql` 18 `max_protocol_version=3.0` verbindet mit `record` und `replay`, `latest` und `3.2` ergeben `PGR-E6002`.
- Meldungscodes: jede Kennung des Katalogs in `internal/hexagon/model/fehler.go` steht in der Tabelle, jede Kennung der Tabelle im Katalog. `grep` auf Spezifikation, Lastenheft, ADR, Slice, Welle, Review, `SPEC-`, `LH-` im Handbuch: 0.
- Größe (Metrik `AGENTS.md` §3.14): `git diff --numstat cf08408..HEAD -- docs/user/benutzerhandbuch.md README.md` = 78+23 und 9+1, Summe 111; Marke 200 nicht erreicht, stimmt mit §7 überein.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | „`replay` beginnt jede Verbindung mit `AuthenticationOk`“ ist nach eigenem Lauf falsch für die Verbindungen, die dieselben Texte als abgelehnt beschreiben (Version 3.2, `_pq_.`, nicht lesbarer Rumpf, unverschlüsselt trotz TLS-Pflicht): sie erhalten sofort die Fehlerantwort, keine Nachricht `R`. Die Quantoren-Durchsicht in §7 erfasste nur den Bestand, nicht diesen neu geschriebenen Satz; im README steht er direkt hinter dem Satz zu den abgelehnten Startnachrichten. | `AGENTS.md` §3.11 (Satz sagt nur zu, was ein Test prüft) | `README.md` · „`replay` beginnt jede Verbindung mit `AuthenticationOk`“; `docs/user/benutzerhandbuch.md` · „`replay` beginnt jede Verbindung mit `AuthenticationOk`, ohne ein Passwort zu“ | ja — Probe (Rohverbindung mit Version 3.2 gegen `replay`); `TestE2EReplayAuswahl` prüft nur die angenommene Verbindung | Quantor „jede“ in neu geschriebenem Satz nicht gegen die Ausnahmen derselben Seite geprüft |
| F-2 | INFO | Der Satz „unter Linux im Container geprüft“ trägt durch Läufe, die dieser Slice nicht fuhr: §7 nennt für ihn „Läufe der Tests im Container (`make test-integration`)“ ohne Lauf, Datum oder Ergebnis, und die Proben des Implementers liefen als Prozess unter Linux (von §7 selbst als nicht tragend benannt). In diesem Review lief jede Probe oben mit dem Binary im Container; der Satz hat damit einen Beleg, aber nicht den, den §7 zeigt. Der Satz nennt „das Verhalten am Protokollrand“ und umfasst damit auch die nur per Probe belegten Client-Aussagen (`psql` 17/18). | `AGENTS.md` §3.11; `BEO-REPO/verhalten-nur-unter-linux-geprueft` | `docs/user/benutzerhandbuch.md` · „Das Verhalten am Protokollrand ist unter Linux im Container geprüft“; Plan · „meine Proben liefen zusätzlich als Prozess unter Linux, nicht im Container“ | ja — `make test-integration` (Geber-Tests) bestätigt die Testseite | Plattform-Satz ohne Lauf-Beleg im Plan |

## Negativbefunde

<!--
Eine Zeile pro betrachtetem Bereich. Ohne diesen Block ist "keine
Findings" nicht von "nicht geprüft" unterscheidbar (Modul 10
§Reviewer berichtet auch, was er nicht gefunden hat).
-->

| Bereich | Ergebnis |
|---|---|
| `docs/user/benutzerhandbuch.md`, Sätze zu `_pq_.` (record, replay, mit und ohne TLS-Pflicht, `PGR-E6001` statt `PGR-E6003`), Version 3.2 (`PGR-E6002`), `libpq` 18 `max_protocol_version` | geprüft, ohne Befund (Lauf, siehe oben) |
| Handbuch, Zählung abgelehnter Verbindungen im `replay` (`first-request`, `connection`) und im `record` mit `--record-empty-sessions` | geprüft, ohne Befund (Lauf: Sitzung 1 und 2 für die regulären Verbindungen, genau 1 Sitzung in der Datei) |
| Handbuch/README, `p` und `CopyData` nach dem Aufbau, COPY-Antworten (`TO STDOUT`, `FROM STDIN`) im `record` | geprüft, ohne Befund (Lauf, beide Modi bzw. `record`) |
| Handbuch, `LISTEN`/`NOTIFY` im `record` (allein, eine Anfrage, zwei Anfragen, zwei Verbindungen, nach der letzten Anfrage) | geprüft, ohne Befund (Lauf; „spätestens auf die nächste Anfrage“ trägt, die Fehlerantwort kam teils schon auf die Anfrage mit `NOTIFY`) |
| Handbuch, `require_auth` und feste Auswahl im `replay` (Benutzer, Datenbank, Optionen ungeprüft; Exit-Code 0, keine Meldung) | geprüft, ohne Befund außer F-1 (Lauf mit `psql` 17, mit und ohne TLS; Rohverbindung über TLS) |
| Weggelassene Sätze: R4 `replay` mit `LISTEN`/`NOTIFY`, P6 `CancelRequest` „wertet nicht aus“ | geprüft, ohne Befund: weder Handbuch noch README sagen sie zu; `PGR-W3001` trägt den Wortlaut ohne P6, wie §6 verlangt; §7 meldet P6 als nicht baubar und den Fund für den Planner |
| Meldungscodes-Tabelle gegen `internal/hexagon/model/fehler.go` | geprüft, ohne Befund (beide Richtungen deckungsgleich; Zeilen E6001 bis E6003 und W3001 stimmen mit §6 überein) |
| Ist-Zustand `AGENTS.md` §3.11: Chronik, Zielstand, Verweise im Handbuch; README | geprüft, ohne Befund (0 Treffer für Spezifikation/Lastenheft/ADR/Slice/Welle/Review; die 4 Treffer für „noch nicht“ u. ä. sind Bestand; README ohne neuen Verweis und ohne Chronik) |
| Plan folgt dem Diff (`AGENTS.md` §3.9), Kopf, §1, §3, §6, §8-Zählung (2 Liefer-Punkte, eine Schicht), Größe nach Metrik §4 (111 von 200) | geprüft, ohne Befund |
| Randformen im Code-Commit (§3.12), Adressen (§3.13): Diff ändert nur Handbuch, README und diesen Plan; keine neue Randform in §6 im Implementer-Commit, keine neue Adresse | geprüft, ohne Befund |
| Beispiele und Codeblöcke als Datei gegen das Binary | geprüft, ohne Befund: der Diff enthält keinen neuen oder geänderten Codeblock (nur Prosa und Tabellenzeilen); Aufrufe stammen aus dem Bestand |
| `internal/`, `cmd/`, `test/`, `harness/`, `spec/` | nicht im Diff, unberührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Quantor „jede“ in neu geschriebenem Satz nicht gegen die Ausnahmen derselben Seite geprüft · Plattform-Satz ohne Lauf-Beleg im Plan

## Verdikt

**Merge-blockierend:** nein — kein HIGH und kein MEDIUM. F-1 ist ein Satz, den ein Randfall der eigenen Seite widerlegt, ohne dass er einen Lauf oder ein Gate bricht; F-2 ist ein Hinweis ohne erwartete Aktion am Code.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
