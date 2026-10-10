# Review-Report: slice-doku-ist-stand — 2026-10-10

**Review-Art:** Code-Review (Dokumentations-Diff). Geprüft wird gegen Plan §1, §3 und §6 von `slice-doku-ist-stand`, gegen `AGENTS.md` §3.9, §3.10, §3.11 (auch *Handbuch und README beschreiben den Ist-Zustand*), §3.12 und §3.13 und gegen die Entscheidungen des Architect in den Randformen (§6). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers. Maßstab für „wahr“ ist das gebaute Binary, nicht die Spezifikation.

**Gegenstand:** Diff `c8df409..9f8f48d`: `248a248` (Architect, Randformen in §6 und DoD), `94f26b0` (Implementer, `docs/user/benutzerhandbuch.md` +102 −214, `README.md` +26 −33), `9f8f48d` (Plan §7, Belege des Implementers). Rahmen: [MR-000](../../harness/conventions.md#mr-000--baseline-aussage).

**Skill:** `.harness/skills/reviewer.md` @ `6db272b`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-doku-ist-stand` am Stand `9f8f48d`, ganz gelesen, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen; die Zusagen dort habe ich nicht übernommen, sondern am Binary nachgespielt.
- `AGENTS.md` §3.9 bis §3.13; Rahmen [MR-000](../../harness/conventions.md#mr-000--baseline-aussage).
- Katalog `internal/hexagon/model/fehler.go`; `spec/spezifikation.md` an `LH-FA-05.b`, `LH-FA-17.a` (Fund SCRAM bei `record`).
- Die Nehmer der Plan-Zuweisungen, §1 und DoD: `slice-v1-abschluss-anmeldung`, `slice-v1-abschluss-einspielen-anmeldung`, `slice-v1-abschluss-einspielen-extended`, `slice-v1-abschluss-einspielen-tls`, `slice-v1-abschluss-tls-client`, `slice-v1-abschluss-homebrew`, `slice-v1-abschluss-container`, `slice-erster-release-veroeffentlichung`, `slice-erster-release-homebrew-nachweis`.
- Vorherige Findings am Gegenstand: Report zu `slice-v1-abschluss-einspielen-laufsteuerung` (F-564 bis F-569). Die Nummern dieses Laufs beginnen bei F-570.

**Ausgeführte Läufe:**

- **Binary.** `make build` am Stand `9f8f48d`; Image `sha256:b45c7f918488…`, `linux/amd64`, Benutzer `nonroot`, Einstiegspunkt `/pgwire-recorder`. Alle Proben in einem Scratch-Verzeichnis außerhalb des Repos (Präfix `rev-doku-`); Netz `rev-doku-net` mit dem gepinnten PostgreSQL-Image aus `harness/mk/integration.mk`, einmal mit `trust` und einmal mit Passwort (SCRAM). Container, Netz und Dateien danach entfernt; `docker ps -a` und `docker network ls` mit `rev-doku` leer.
- **Stichproben, 30 Aussagen quer, alle bestätigt:**
  - Optionen: `--help` für allgemein, `record`, `replay`, `play`, `config show` gegen die Tabelle in Handbuch §5 (14 Zeilen, Betriebsart-Spalte und Standardwerte); `-help`, `--h`, `version x` (`PGR-E2001`), `--version` (`PGR-E2001`).
  - Entfernte Optionen: `play --compare-responses` → `PGR-E2001`, `config show` mit `record.format` und `play.password` → `PGR-E2004` „unbekannter Schlüssel“.
  - Umgebungsvariablen: `PGWIRE_RECORDER_FORCE=1` → `PGR-E2001`, `=true` ersetzt die Datei; `--shutdown-timeout=5` → `PGR-E2001`; `--fail-on-unconsumed=1` → `PGR-E2001`, `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=yes` → `PGR-E2001`; `config show` listet gesetzte Variablen ohne Wert, auch ohne zugehörige Option.
  - Codes und Exit-Codes: `PGR-E2002` (Exit 2), `PGR-E2004` (`sslmode=require` bei `play` und `record`), `PGR-E2005`, `PGR-E2006`, `PGR-E3001` (Verzeichnis als `--output`, fehlendes Verzeichnis), `PGR-E3002`, `PGR-E3003`, `PGR-E3004` (`replay` Exit 3, `play` Exit 0), `PGR-E4002`, `PGR-E4004` (mit und ohne `--allow-recorded-errors`), `PGR-E4005` (Passwort verlangt; `--user niemand`), `PGR-E5001`, `PGR-E5002`, `PGR-E5003`, `PGR-E6001` (`record` gegen SCRAM, Client erhält FATAL, `sessions: []`, Exit 6), `PGR-E6002`, `PGR-W3003`. Exit-Code 5 bei `--fail-on-unconsumed`.
  - Verschlüsselung: `psql sslmode=require` an `replay` → „server does not support SSL“; `prefer` verbindet.
  - Installation und Log: `docker image inspect` (Benutzer, Einstiegspunkt), Binary aus dem Image kopiert (`version` → `pgwire-recorder dev`, statisch gebunden, x86-64), Log-Zeitstempel mit `TZ=Europe/Berlin` (`+02:00`), neue Zieldatei `0644` unter umask `022`.
  - Schlüssel der Datei: Platzhalter-Passwort für `play` (`PGR-E4005` gegen SCRAM, Exit 4), Klartext-Passwort → `PGR-E2006`.
- **Katalog und Optionen mechanisch.** Eigenes Skript (Scratch, nicht im Repo): jede `--option` des Handbuchs gegen die Vereinigung der `--help`-Ausgaben, jede `PGWIRE_RECORDER_*`-Variable gegen eine Option, jeder `PGR-…`-Code gegen den Katalog. Handbuch `9f8f48d`: nur Treffer, die das Handbuch selbst als Nicht-Option nennt (`--h`, `--help`, `--version`) oder die zu `docker run` gehören (`--rm`, `--name`); keine Variable ohne Option, kein Code außerhalb des Katalogs. Handbuch `c8df409`: 39 verschiedene Meldungen. Mutanten am Handbuch von `9f8f48d`, je eine angehängte Zeile: `--compare-responses`, `PGWIRE_RECORDER_PASSWORD`, `PGR-E5004`, `--keep-timing`, `PGWIRE_RECORDER_FORCEE` — alle fünf rot, mit dem genannten Namen.
- **`grep`.** Handbuch `-i` auf `spezifikation|lastenheft|\bADR-|slice|welle|review`: 0 Treffer. `LH-`, `SPEC-`, `MR-`, `AGENTS`: 0 Treffer. README auf `noch|erste Version|Beginn der Umsetzung|auf Wunsch|TLS`: 0 Treffer. Zeitwörter im Handbuch (`derzeit`, `später`, `künftig`, `bisher`, `vorerst`, `dieser Stand`, `kommt`, `folgt`): die Treffer sind „spielt sie später ohne die Datenbank wieder ab“ (Ablauf, kein Zielstand) und „existiert noch nicht“ (Zustand der Datei).
- **Gates.** `make docs-check` vor dem Commit dieses Reports: grün.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-570 | MEDIUM | Der Plan weist in §6 zwei Sendungen an andere Slices zu, ohne dass einer der Nehmer sie mit der Kennung des Gebers trägt. (a) *Installation*: „Die Abschnitte liefern `slice-v1-abschluss-homebrew`, `slice-v1-abschluss-container`, `slice-erster-release-veroeffentlichung` und `slice-erster-release-homebrew-nachweis`, je mit ihrer Zeile zum Ist-Zustand in DoD und §3“. (b) *Grenz-Sätze*: „der liefernde Slice ersetzt den Grenz-Satz“ — im Handbuch an `record` mit Passwort (§1, §4, §6, §7), `play` mit Passwort, `sslmode=require` und Extended. `grep slice-doku-ist-stand` in den neun Nehmern: 0 Treffer. Ihre DoD trägt nur die allgemeine Zeile „beschreiben, was dieser Slice liefert“ mit der Kennung von `slice-v1-abschluss-einspielen-laufsteuerung`; sie nennt weder die Installationsabschnitte noch die Grenz-Sätze, die zu ersetzen sind. *Failure-Szenario:* `slice-v1-abschluss-anmeldung` ergänzt einen Absatz zu SCRAM bei `record`; §1 *Voraussetzungen* („Datenbank … ohne Passwort“), der Hinweis in §4 und die Zeile `PGR-E6001` bleiben stehen und widersprechen dem Binary — der Fall, den §3.11 und die neue Regel verhindern sollen, und kein Nehmer weist ihn aus. §8 der Nehmer nennt die Sendung nicht. | `AGENTS.md` §3.13 (Nehmer trägt die Sendung mit der Kennung des Gebers; Nachzählen beim Eintragen); Reviewer-Skill MEDIUM *Adresse nimmt nicht an* | Plan · „Die Abschnitte liefern“; Plan · „der liefernde Slice ersetzt den Grenz-Satz“; `docs/plan/planning/next/slice-v1-abschluss-anmeldung.md` · „beschreiben, was dieser Slice liefert, im Ist-Zustand“ | ja (`grep -rl slice-doku-ist-stand docs/plan/planning/next` findet keinen Nehmer) | Zuweisung an Folge-Slice ohne Eintrag beim Nehmer |
| F-571 | LOW | Handbuch §7, Zeile `PGR-E2001`: „oder passt nicht zu einer anderen Option“. Das Binary kennt keine Kombination von Optionen mehr, die so abgelehnt wird: Die Optionspaare der entfernten Optionen (`--tls-cert`/`--tls-key`, `--keep-timing`/`--timing-mode`) sind weg, und `--compare-responses` mit `--allow-recorded-errors` gibt es nicht mehr. Ich habe keine solche Probe gefunden; der Satz sagt einen Fall zu, den weder Probe noch Test hält. | `AGENTS.md` §3.11 (Handbuch sagt nur zu, was ein Test prüft) | `docs/user/benutzerhandbuch.md` · „oder passt nicht zu einer anderen Option“ | ja (Probe je Optionspaar; kein Gate) | Zusage im Handbuch weiter als die Prüfung |
| F-572 | LOW | Plan §6 *Verweise im README* begründet den Wegfall der Aufzählung mit „seine Liste ist heute unvollständig“. `harness/README.md` §Sensors führt jedes der 14 Ziele aus `GATE_CHECKS` (`make -pn gates`). Die Zeile ist ohne Beleg; gemeint kann nur die alte Aufzählung im README gewesen sein. Das Urteil „unvollständig“ ist weder in §7 belegt noch als Fund an den Planner gegeben. | `AGENTS.md` §3.11 (Plan-Zeile sagt nur zu, was geprüft ist) | Plan · „seine Liste ist heute unvollständig“ | ja (`GATE_CHECKS` gegen die Zeilen in `harness/README.md` §Sensors) | Zusage im Plan weiter als die Prüfung |
| F-573 | INFO | Der Fund für den Planner (SCRAM bei `record` gegen `LH-FA-05.b` und `LH-FA-17.a`) ist als *Verhalten neben der Spezifikation* eingeordnet. Die Probe bestätigt das Verhalten (Exit 6, `sessions: []`). Es ist aber der bekannte Teilstand: Die Roadmap führt es seit 2026-10-04 im Drift-Log („Der Record-Modus vermittelt im Walking Skeleton keine Passwort-Anmeldung; LH-FA-05.b verlangt sie“), und `slice-v1-abschluss-anmeldung` hat es als Ziel. Der Fund nennt diesen Nehmer nicht; `LH-FA-17.a` regelt die Anmeldung des Clients nicht (die Spezifikation sagt dort nur, `record` vermittle sie, in der Aussage zu `sslmode`). Die Handbuch-Sätze selbst sind als Grenze im Präsens richtig. | `AGENTS.md` §3.13 (Adresse); Plan §6 *Teilweise geliefert* | Plan §7 · „`spec/spezifikation.md` LH-FA-05.b (*Record*: vermittelt die Authentifizierung transparent) und“ | nein | — (Hinweis an den Planner: Teilstand mit Adresse, keine neue Abweichung) |
| F-574 | INFO | Die Zeile `PGR-E1000` fiel aus der Tabelle *Fehlercodes*, weil keine Probe sie auslöst. Der Code steht im Katalog (`CodeInternal`), die Klasse `sonstiger Fehler` steht im Handbuch weiter, und Exit-Code 1 bleibt in der Tabelle. Wer `sonstiger Fehler [PGR-E1000]` liest, findet keine Zeile. Die Regel aus Plan §6 trägt die Streichung; der entfallene Hinweis („Starten Sie mit `--log-level debug` neu, und melden Sie das Problem“) war für den Ist-Zustand wahr. | Plan §6 *Wie geprüft wird* | `docs/user/benutzerhandbuch.md` · „| 1 | sonstiger Fehler |“ | nein | — (Hinweis an den Verifier und an den Planner) |
| F-575 | INFO | Der Beleg „Prüfskript als Gegenprobe“ in Plan §7 verweist auf ein Skript im Scratch, das nicht im Repo liegt. Ich habe die Prüfung selbst nachgebaut (siehe *Ausgeführte Läufe*); die drei Mutanten aus §7 und zwei eigene werden rot. Das Skript prüft Namen (Option, Variable, Code), nicht die Zuordnung zu Kommandos, nicht Standardwerte und nicht Exit-Codes; die Zeilen der Tabelle in §5 habe ich von Hand gegen `--help` gehalten (14 Zeilen, Spalten stimmen). In §7 steht außerdem „Tabelle mit 16 Zeilen wie `--help`“, die Tabelle hat 14 Datenzeilen (16 mit Kopf und Trennzeile). | `AGENTS.md` §3.9; Plan §1 (kein Sensor, Prüfung von Hand) | Plan §7 · „Tabelle mit 16 Zeilen wie `--help`“ | nein | — (Hinweis an den Verifier) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/user/benutzerhandbuch.md`, Kopf und §11 | geprüft, ohne Befund. `pgwire-recorder version` gibt `pgwire-recorder dev`; Kopfzeilen und der Satz „Es gibt keine veröffentlichte Version.“ nennen den Ist-Zustand ohne Chronik; kein Verweis auf Spezifikation, Lastenheft, ADR, Slice, Welle oder Review (0 Treffer). |
| `docs/user/benutzerhandbuch.md`, §1 und §2 | geprüft, ohne Befund. Kein macOS, kein Windows, keine Registry, kein Homebrew; `make build`, `version`, das Binary aus dem Image und der Benutzer `nonroot` hielten in den Proben; die Compose-Beispiele nennen `pgwire-recorder:dev`. |
| `docs/user/benutzerhandbuch.md`, §3 und §4 | geprüft, ohne Befund (F-571 betrifft §7). Aufzeichnen, Wiedergeben, Einspielen, `config show`, Rechte der Zieldatei, Herunterfahren, `--fail-on-unconsumed`: alle Aussagen der Stichprobe halten am Binary. Die Grenz-Sätze (Passwort bei `record` und `play`, `sslmode=require`, Extended bei `play`) zeigen die genannte Reaktion; die Folgen für die Nehmer stehen in F-570. |
| `docs/user/benutzerhandbuch.md`, §5 | geprüft, ohne Befund. 14 Optionen, 14 Variablen, Betriebsarten, Standardwerte und Prioritäten stimmen mit `--help`; Schlüssel der Datei (`config show`, `PGR-E2004` bei Entferntem); Exit-Codes 0 bis 6. |
| `docs/user/benutzerhandbuch.md`, §6 bis §10 | geprüft, ohne Befund außer F-571 und F-574. Codes der Tabellen stehen alle im Katalog und wurden bis auf `PGR-E4000`, `PGR-E4001`, `PGR-E4003`, `PGR-E4006`, `PGR-W2001` und `PGR-W3001` im Binary ausgelöst; diese sechs trägt der Beleg des Implementers mit Probe bzw. E2E-Test, ich habe sie nicht erneut ausgelöst. |
| Streichungen (`git diff c8df409..9f8f48d`, alle `-`-Zeilen) | geprüft, ohne Befund außer F-574. Entfernt sind Abschnitte und Optionen, die das Binary mit `PGR-E2001` oder `PGR-E2004` ablehnt (TLS zum Client, SQLite, Zeitangaben, Vergleich, Sitzungszuordnung, `--upstream-*`, `PGWIRE_RECORDER_PASSWORD`, Release-Binaries, Homebrew, Registry), dazu die Codes ohne Katalogeintrag. Keine Streichung nahm eine für den Ist-Zustand wahre Aussage mit, außer der Zeile `PGR-E1000`. |
| `README.md` | geprüft, ohne Befund. Leitsatz, Einleitung, *Was kann ich heute tun?* und *Kerngedanke* nennen nur Geliefertes; „lehnt `play` ab“, „unverschlüsselt“, „höchstens `--shutdown-timeout` (Standard 5 Sekunden)“ halten am Binary; `make build` baut `pgwire-recorder:dev`. Verweise auf `spec/`, `docs/plan/`, `docs/reviews/` ohne Aussage über einen Stand. |
| Plan §1, §3, §6 gegen den Diff (`AGENTS.md` §3.9) | geprüft, ohne Befund außer F-572. Ziel, Abgrenzung, Tabelle und Randformen folgen dem Diff; Größe in §7 (+102 −214, README +26 −33) stimmt mit `git diff --stat`. Kein Randform-Eintrag im Code-Commit: Die Entscheidungen stehen in `248a248` vor dem Handbuch-Commit `94f26b0` (`AGENTS.md` §3.12). |
| Mutations-Regel (`AGENTS.md` §3.10) | nicht anwendbar, geprüft: Der Slice liefert keinen neuen Vertrag (Plan §1, §6, Randformen); die Gegenprobe der Prüfung ist als Beleg geführt und von mir wiederholt (F-575). |
| Commit-Struktur (`AGENTS.md` §3.3) | geprüft, ohne Befund. Kein Move im Diff, drei Commits, jeder nennt `slice-doku-ist-stand` und `MR-000`. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zuweisung an Folge-Slice ohne Eintrag beim Nehmer · Zusage im Handbuch weiter als die Prüfung · Zusage im Plan weiter als die Prüfung

## Verdikt

**Merge-blockierend:** ja — wegen F-570 (MEDIUM). Der Handbuch- und README-Teil des Diffs hält in den Stichproben am Binary: keine Chronik, kein Zielstand, kein Verweis, 30 von 30 Aussagen bestätigt, 5 von 5 Mutanten des Prüfskripts rot. F-570 betrifft die Zuweisung an die Folge-Slices, nicht den Text der beiden Dokumente; F-571 und F-572 sind kleine Zusagen ohne Prüfung.

**Übergabe:** F-570 geht an den Implementer (Einträge bei den Nehmern, mit Nachzählen in deren §8); F-571 und F-572 an den Implementer; F-573 und F-574 als Hinweis an den Planner; F-575 an den Verifier. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier.
