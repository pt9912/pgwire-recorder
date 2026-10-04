# Review-Report: slice-walking-skeleton-build-gates — 2026-10-04

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD ist nicht Gegenstand; sie prüft der Verifier.

**Gegenstand:** `git diff 759fc4f HEAD` (HEAD `8143a47`): `1703db5` (Verantwortlich gesetzt), `554bd0b` und `d04e43c` (Moves open → next → in-progress, von Hand per `git mv`), `8143a47` (`go.mod`, `cmd/pgwire-recorder/main.go`, `doc.go` je Package unter `internal/`, `harness/mk/build.mk`, `harness/mk/arch-negativ.mk`, `tools/arch/a-check-negativ.sh`, `harness/README.md` §Sensors, `README.md`, Slice-Plan §3).

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- Slice-Plan `slice-walking-skeleton-build-gates` (§1 bis §6, §8)
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0009](../plan/adr/0009-implementierungssprache-go.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md)
- `LH-FA-01`, `LH-QA-03`, `LH-QA-04`
- `spec/architecture.md` §1, §2, §2.1, §3 (`ARC-012`), §6; `.a-check.yml`
- `AGENTS.md` (Hard Rules 3.1 bis 3.8), `harness/README.md`, `harness/conventions.md`
- `tools/harness/commit-msg-traceability.sh`, `tools/harness/slice-mv.sh` (für die Frage nach dem Lifecycle-Werkzeug)

**Ausgeführte Läufe (alle schreibgeschützt):** `make help`, `make a-check` (0 Befunde), `make build` (grün), `make test` (grün, alle Packages „no test files"), `make a-check-negativ` (grün). Dazu Mutationen von a-check an Kopien des Arbeitsbaums im Scratchpad, nicht im Repo: `pgproto3` in `internal/hexagon/model` (core-impurity, rot); dasselbe ohne die `pgproto3`-Zeilen der `tech`-Liste (0 Befunde, grün — die Aussage des Implementers ist bestätigt); interne Verstöße services → ports/driving, cli → driven/recording, driven/postgres → ports/driving, ports/driven → driven/postgres (alle vier gemeldet); `os` in `internal/hexagon/model` (0 Befunde); `pgproto3` bzw. `crypto/tls` in `internal/adapters/driven/postgres` (tech-leak, rot, siehe F-204). `git status` nach allen Läufen sauber. `make gates` nicht gelaufen.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-204 | HIGH | Die `tech`-Regel des Architektur-Gates lässt `github.com/jackc/pgx/v5/pgproto3` und `crypto/tls` nur in `internal/adapters/driving/pgwire` zu: a-check wertet bei zwei Einträgen mit gleichem `pattern` nur den ersten, ein `pgproto3`- oder `crypto/tls`-Import in `internal/adapters/driven/postgres` ergibt „tech-leak: Tech … außerhalb internal/adapters/driving/pgwire". [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) und `ARC-012` erlauben die Bibliothek in beiden PGWire-Adaptern; das Gate setzt die Entscheidung damit nicht um und wird den Upstream-Adapter ablehnen, sobald er Code trägt. Der Plan führt `.a-check.yml`/`a-check.mk` als „vorhanden, keine Änderung nötig"; das Risiko in §6 („ob a-check die Go-Importe wie modelliert erkennt, erst mit Code belegbar") ist mit diesem Lauf belegt, und §4 nennt dafür eine Rückführung. Der `doc.go` des PGWire-Adapters sagt zugleich zu, die Bibliothek stehe „hier und im Upstream-Adapter". | [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md); [ADR-0001](../plan/adr/0001-hexagonale-architektur.md); `ARC-012` | `.a-check.yml` · `tech`; Slice-Plan §3 (Zeile `a-check.mk`), §6 Risiko 1; `internal/adapters/driving/pgwire/doc.go` | ja — `make a-check` an einer Kopie mit `pgproto3`-Import in `internal/adapters/driven/postgres` | Gate-Konfiguration bildet ADR-Entscheidung nicht ab |
| F-205 | MEDIUM | Der aktive Commit-Träger (`core.hooksPath=.githooks`) akzeptiert als Slice-Kennung nur `slice-[0-9]+`; dieses Repo vergibt benannte Kennungen (`slice-<Kennung>`, `harness/conventions.md` MR-000), und `harness/README.md` §Traceability rules sagt `slice-*` als zulässige Kennung zu. Die Messages von `make slice-mv` tragen nur den Slice-Namen und fallen deshalb durch; die drei Lifecycle-Commits dieses Slice entstanden darum von Hand, ohne den Verweis-Nachzug des Werkzeugs. Das tritt bei jedem Lifecycle-Wechsel jedes benannten Slice erneut ein. Nicht im Diff, vom Slice aufgedeckt. | `harness/README.md` §Traceability rules; `harness/conventions.md` MR-000 (ID-Schema); Maintainability | `tools/harness/commit-msg-traceability.sh` · `patterns=`; `tools/harness/slice-mv.sh` · Commit-Messages | ja — `make slice-mv` mit benanntem Slice bei aktivem Hook | Kennungs-Muster des Commit-Wächters enger als das deklarierte ID-Schema |
| F-206 | LOW | Die Gegenprobe ist in der Sensors-Tabelle an [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) gebunden und der Commit `8143a47` nennt [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) und [ADR-0009](../plan/adr/0009-implementierungssprache-go.md). Was sie färbt, ist die `pgproto3`-Zeile der `tech`-Liste — die Fitness Function von [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md); ohne diese Zeile ist die Gegenprobe rot, ohne Änderung an Schichten oder Kanten nicht. [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) steht weder in der Bindung noch in der Commit-Message. | `harness/README.md` §Sensors (Bindung); [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) | `harness/README.md` · Sensors, Zeile `make a-check-negativ`; Commit `8143a47` | ja (Lesen; Mutation der `tech`-Liste) | Sensor-Bindung zeigt auf benachbarte statt prüfende ADR |
| F-207 | LOW | Die Gegenprobe gilt als grün, wenn a-check mit Exit ungleich 0 endet und die Ausgabe irgendwo `pgproto3` enthält. Welcher Befund (Datei, Regelklasse) es ist, wird nicht geprüft; ein Scheitern aus anderem Grund, dessen Meldung die Zeichenkette enthält (etwa ein Konfigurationsfehler, der die `tech`-Zeile zitiert), färbt sie ebenso grün. | Maintainability | `tools/arch/a-check-negativ.sh` · zweite Prüfung (`grep -q "pgproto3"`) | nein | Gegenprobe prüft Meldungstext unscharf |
| F-208 | LOW | Die Log-Datei `"$kopie.log"` liegt neben, nicht in dem temporären Verzeichnis; der `trap … EXIT` räumt nur das Verzeichnis. Die drei expliziten `rm -f` decken die regulären Ausgänge; bei Abbruch durch ein Signal zwischen `docker run` und dem `rm` bleibt die Datei im Temp-Verzeichnis liegen. Auf den Arbeitsbaum wirkt das Skript nicht (nach dem Lauf `git status` sauber). | Maintainability | `tools/arch/a-check-negativ.sh` · `trap`, Umleitung von `docker run` | nein | Aufräumen deckt nicht alle temporären Artefakte |
| F-209 | LOW | Neun `doc.go` tragen „Das Package trägt noch keinen Code.", `main.go` „Das Programm kennt noch keine Kommandos und endet mit Exit-Code 0." Die Sätze beschreiben den Zustand (Hard Rule 3.7 erfüllt), werden aber mit dem ersten Code in jedem Package falsch, ohne dass ein Sensor es meldet; die Last liegt bei jedem späteren Slice, der ein Package füllt. | AGENTS.md §3.7; Maintainability | `internal/**/doc.go`; `cmd/pgwire-recorder/main.go` · Package-Kommentar | nein | Zustandsaussage im Kommentar ohne Sensor |
| F-210 | INFO | Reichweite der Gegenprobe: Sie belegt, dass a-check einen in der `tech`-Liste bekannten Import im Domain Model als core-impurity meldet. Sie belegt weder die Kantenregeln (die meldet a-check laut Mutation korrekt, aber kein Gate hält das fest) noch den tech-leak-Pfad in den Adaptern (F-204). Die Kopfzeile von `arch-negativ.mk` („belegt, dass das Gate einen verbotenen Import meldet") ist allgemeiner als das, was geprüft wird; die Sensors-Zeile ist dagegen genau. Die DoD verlangt nur diesen einen Fall — Bewertung Verifier. | [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) Fitness Function; Verweis Verifier | `harness/mk/arch-negativ.mk` · Kopfkommentar | ja — Mutation einer Kante | Reichweite einer Gegenprobe allgemeiner beschrieben als geprüft |
| F-211 | INFO | `build`/`test` laufen mit `--network none`, `GOFLAGS=-mod=readonly` und einem leeren `GOPATH` im Container. Das trägt, solange das Modul keine Abhängigkeiten hat (der Kopfkommentar von `build.mk` nennt diese Kopplung). Der nächste Slice der Welle bringt mit [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) die erste Abhängigkeit; woher sie netzlos kommt (vendorte Module, vorbefüllter Cache o. ä.), sagt kein Plan, und `slice-walking-skeleton-record` nennt es nicht unter Risiken. Hinweis an Planner/Architect; kein Defekt dieses Diffs. | AGENTS.md §3.1, §3.6; Verweis Planner | `harness/mk/build.mk` · `GO_RUN`, Kopfkommentar | ja — `make build` nach Hinzufügen einer Abhängigkeit | Netzlos-Zusage an abhängigkeitsfreien Zustand gekoppelt |
| F-212 | INFO | `spec/architecture.md` §2 untersagt dem Domain Model den Zugriff auf das Dateisystem; ein `import "os"` in `internal/hexagon/model` ergibt in a-check 0 Befunde. Das Gate deckt diese Einschränkung nicht; die Sensors-Zeile von `make a-check` behauptet sie auch nicht. Für den Reviewer-Skill (HIGH-Liste „Core-Reinheit … Dateisystem") bleibt damit der Reviewer der einzige Wächter dieses Teils. Nicht im Diff. | `spec/architecture.md` §2; [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | `.a-check.yml` | ja — `make a-check` an einer Kopie mit `os`-Import im Model | Spec-Constraint ohne Gate-Abdeckung |
| F-213 | INFO | `test/integration/` aus `spec/architecture.md` §2.1 ist nicht angelegt; ein leeres Verzeichnis trägt git nicht, und ein Package ohne Test wäre nur ein Platzhalter. Der Plan §3 („leere Packages gemäß Package-Struktur") nennt die Auslassung nicht; zehn offene Slices führen `test/integration` als „update", obwohl der erste es neu anlegt. Alle übrigen Pfade aus §2.1 stimmen mit dem Diff überein. | `spec/architecture.md` §2.1; Maintainability | Slice-Plan §3; `test/integration/` (fehlt) | ja (Lesen) | Auslassung gegenüber Package-Struktur nicht im Plan notiert |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Package-Struktur gegen `spec/architecture.md` §2.1 | geprüft; `cmd/pgwire-recorder`, `internal/hexagon/{model,services,ports/driving,ports/driven}`, `internal/adapters/{driving/{cli,pgwire},driven/{postgres,recording}}`, `internal/bootstrap` vollständig und gleich benannt; Package-Namen passen zum Verzeichnis; `ARC-*`-Zuordnung der `doc.go` stimmt mit §1/§2.1; Ausnahme `test/integration/`, F-213 |
| Core-Reinheit (`internal/hexagon/**`) | geprüft, ohne Befund; keine Imports, nur Package-Klauseln |
| `go.mod` | geprüft, ohne Befund; Modulpfad, `go 1.27` passt zum Image `golang:1.27-alpine`, keine `require`-Zeilen |
| `harness/mk/build.mk` — netzlos, schreibgeschützt, gepinnt | geprüft, ohne Befund; `--network none`, Mount `:ro`, `-u` des Aufrufers, `GOTOOLCHAIN=local`, `-mod=readonly`, Cache/GOPATH im Container unter `/tmp`, Image per Digest, `-trimpath -buildvcs=false`; `make build`/`make test` grün, Arbeitsbaum danach sauber; Kopplung an Abhängigkeitsfreiheit siehe F-211; die Digest-Aussage „Manifestliste für amd64 und arm64" ist netzlos nicht nachgeprüft |
| `harness/mk/arch-negativ.mk` | geprüft; Variablen aus `a-check.mk` werden zur Rezeptzeit expandiert, die Include-Reihenfolge ist damit unerheblich; hängt an `GATE_CHECKS`; Reichweite siehe F-210 |
| `tools/arch/a-check-negativ.sh` — Seiteneffekte auf den Arbeitsbaum | geprüft, ohne Befund; Kopie per `tar` nach `mktemp -d`, `.git` ausgenommen, Mount `:ro`, `--network none`; prüft die `.a-check.yml` des Arbeitsbaums; Lecks und Schärfe siehe F-207, F-208 |
| Gegenprobe als Wächter | geprüft; Mutation „`pgproto3`-Zeilen der `tech`-Liste entfernt" färbt sie rot (bestätigt); interne Kantenverstöße meldet a-check korrekt (vier Mutationen); Bindung und Reichweite siehe F-206, F-210 |
| Hard Rule 3.1 — Docker-only | geprüft, ohne Befund; alle neuen Ziele laufen über `make` in Docker; das Skript braucht auf dem Host nur `bash`, `tar`, `mktemp`, `grep` und die Runtime |
| Hard Rule 3.2 — Suppression | geprüft, ohne Befund; keine Inline-Suppression, keine Änderung an `.a-check.yml` |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund; `554bd0b` und `d04e43c` sind reine Renames (0 Einfügungen, 0 Löschungen), `1703db5` ändert nur Inhalt in `open/`, die Plan-Änderung §3 steht in `8143a47` nach dem Move; keine Verweise auf den Lifecycle-Pfad des Slice blieben zurück (Welle und Folge-Slice zitieren die Kennung) |
| Hard Rule 3.5 — ADRs | geprüft, ohne Befund; keine ADR im Diff |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund; keine Schwelle gesenkt, drei Gates hinzugekommen; „Nicht behauptet" nennt Lint und Testabdeckung |
| Hard Rule 3.7 — Kommentare in `doc.go`, `main.go`, Make-Fragmenten, Skript | geprüft; Indikativ, Zusage/Kopplung/Grenze, keine verworfene Alternative, kein abwesender Text, kein abgebrochener Satz; Zustandsaussagen siehe F-209, Reichweite F-210 |
| Hard Rule 3.8 | geprüft, ohne Befund; keine ADR im Diff |
| Sensors-Tabelle und `README.md` gegen `make help` | geprüft, ohne Befund; `build`, `test`, `a-check`, `a-check-negativ`, `docs-check`, `baseline-verify`, `gates` existieren; die Aufzählung in `README.md` deckt sich mit den an `GATE_CHECKS` gehängten Zielen dieses Diffs; Bindung siehe F-206 |
| Slice-Plan §3 | geprüft; die Zeilen geben den Diff wieder (kein `Dockerfile`, Gegenprobe als Gate); Zeile `a-check.mk` siehe F-204, Auslassung siehe F-213 |
| Commit-Messages | geprüft; jede trägt `LH-*`, `8143a47` zusätzlich ADR-Kennungen; keine Struktur-IDs; fehlende [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) siehe F-206; Werkzeug-Konflikt siehe F-205 |
| `spec/` | geprüft, ohne Befund; nicht im Diff |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Gate-Konfiguration bildet ADR-Entscheidung nicht ab · Kennungs-Muster des Commit-Wächters enger als das deklarierte ID-Schema · Sensor-Bindung zeigt auf benachbarte statt prüfende ADR · Gegenprobe prüft Meldungstext unscharf · Aufräumen deckt nicht alle temporären Artefakte · Zustandsaussage im Kommentar ohne Sensor · Reichweite einer Gegenprobe allgemeiner beschrieben als geprüft · Netzlos-Zusage an abhängigkeitsfreien Zustand gekoppelt · Spec-Constraint ohne Gate-Abdeckung · Auslassung gegenüber Package-Struktur nicht im Plan notiert

## Verdikt

**Merge-blockierend:** ja — F-204 (HIGH). Das Gate des Slice setzt [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) nicht um, und der Plan nennt für genau diesen Fall eine Rückführung (§4). Widerspricht der Implementer der Einstufung, läuft der Konflikt-Pfad über den Architect (Modul 8). F-205 (MEDIUM) liegt außerhalb des Diffs und blockiert diesen Slice nicht; es betrifft jeden weiteren Lifecycle-Wechsel und gehört an den Harness (`harness/conventions.md`).

**Übergabe:** Findings an den Implementer; F-211 zusätzlich an den Planner (Folge-Slice `slice-walking-skeleton-record`), F-212 an den Architect. Die Finding-Klassen gehen in die Slice-Closure §7. DoD-Konformität prüft der Verifier separat.
