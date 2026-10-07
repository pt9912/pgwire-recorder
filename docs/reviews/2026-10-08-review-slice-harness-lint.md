# Review-Report: slice-harness-lint — 2026-10-08

**Review-Art:** Code-Review (Gate-Anschluss, Gegenprobe, Sensor-Datei, Kommentare, Belege in §7). Geprüft wird gegen Plan §1, §3 und §6 (mit den Randform-Rückgaben aus `7a32d57`), `SPEC-049` (Punkte 1 bis 10 und Grenze), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und die Hard Rules (Modul 10). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff fdbcfaf..3bcf000`:

- `221845d`: `harness/mk/lint.mk` (beide Ziele an `GATE_CHECKS`), `tools/harness/lint-gegenprobe.sh` (neu), `harness/sensors/lint.md` (neu), `harness/README.md`, `AGENTS.md` §3.2, Kommentare in `Dockerfile`, `tools/harness/lint.sh` und `internal/hexagon/model/fehler.go`, Plan §3 und §7.
- `7a32d57`: Architect, `SPEC-049` Punkt 8 (Leerraum nach dem Doppelpunkt, `-` allein), Historien-Zeile, Plan §6 *Randform-Rückgaben aus `221845d`*.
- `3bcf000`: vier Fälle zu den Rückgaben, Kopf der Gegenprobe, Plan §7 *Nachtrag zu den Rückgaben*.

**Skill:** `.harness/skills/reviewer.md` am Stand `3bcf000`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-lint.md`: ganz gelesen am Stand `3bcf000`, dazu der Plan-Diff je Commit. §7 *Belege des Implementers* habe ich gelesen, den Bericht des Implementers nicht.
- `spec/spezifikation.md` `SPEC-049` ganz, Historie.
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md).
- `AGENTS.md` §3.2, §3.3, §3.7, §3.9 bis §3.13.
- `tools/harness/lint.sh` und `.golangci.yml` ganz (Stand `3bcf000`), `Makefile`, `harness/mk/build.mk`.
- Register: `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (`observation.md`, Abgrenzung zu `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`).
- Vorherige Reports: `slice-harness-lint-werkzeug` (Review und Verifikation), `slice-lint-bestand-driving` (F-465 bis F-470). Die höchste vergebene Nummer vor diesem Lauf war F-470.

**Ausgeführte Läufe:**

- Fünf frische Kopien des Arbeitsbaums an `3bcf000` unter dem Scratch-Pfad, je mit `cp -R` ohne `-p` und ohne `.git`. Mutationen per `sed` mit geprüfter Trefferzahl, danach `touch` auf die Datei. Nichts lief im Repo.
- In vier Kopien `make lint-gegenprobe` gleichzeitig (`LINT_GEGENPROBE_PARALLEL=4`). Die unveränderte Kopie ist grün: `lint-gegenprobe: gruen …`, Ausgang 0. Unter dieser Last dauerte jeder Lauf 349 bis 400 s; das ist kein Maß für die Laufzeit allein.
- In der fünften Kopie `make -k gates` mit einem `//nolint` an `func Meldungen` (G1, Tabelle unten).
- Eigene Mutationen: Tabelle unten. Die 35 Mutationen aus §7 habe ich nicht nachgefahren, ihre Fälle habe ich am Skript gelesen.
- `make docs-check` vor dem Commit.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-471 | MEDIUM | `SPEC-049` Punkt 6 sagt zu, dass jede Zeile einer Datei `*.go` unter `cmd/`, `internal/` und `test/` geprüft wird, auch in Testdateien. Kein Fall der Gegenprobe legt ein `nolint` in eine Datei `*_test.go`; alle zehn roten Fälle von Punkt 6 liegen in Dateien ohne `_test`. Mein Mutant M1 nimmt Testdateien aus der Prüfung (`go_dateien '*.go' \| grep -v '_test\.go$'`) und `make lint-gegenprobe` bleibt grün. Der Kopf führt ihn nicht unter OFFEN. *Failure-Szenario:* Eine Änderung an `go_dateien` oder an der Schleife nimmt Testdateien heraus. Danach geht `//nolint` in jedem Test durch, und `AGENTS.md` §3.2 sagt weiter, die Direktive breche das Gate. | `AGENTS.md` §3.10; `SPEC-049` Punkt 6; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | `tools/harness/lint-gegenprobe.sh` · `NF="$G/nolint/fall.go"`; `tools/harness/lint.sh` · `done < <(go_dateien '*.go')` | ja (M1 grün) | Zusage eines neuen Vertrags ohne unterscheidende Mutation |
| F-472 | MEDIUM | `SPEC-049` Punkt 5 sagt zu, dass `testpackage` nur eine Datei mit dem Namen `export_test.go` überspringt, `skip-regexp` `(^\|/)export_test\.go$`. Kein Fall legt eine Testdatei an, deren Name auf `export_test.go` endet, aber anders heißt. Mein Mutant M2 (`skip-regexp: export_test\.go$`) bleibt grün. Den Mutanten aus §7 (`(export\|internal)`) fängt `p7-internal-test`, M2 nicht. *Failure-Szenario:* Ist der Anker weg, ist eine White-Box-Datei `fooexport_test.go` im Paket des Codes kein Befund. Messmethode 3 von `LH-QA-07` hält dann nicht mehr, und das Gate bleibt grün. | `AGENTS.md` §3.10; `SPEC-049` Punkt 5 und 7; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | `.golangci.yml` · `skip-regexp: (^\|/)export_test\.go$`; `tools/harness/lint-gegenprobe.sh` · `kein_befund p7-bruecke-kein-testpackage` | ja (M2 grün) | Zusage eines neuen Vertrags ohne unterscheidende Mutation |
| F-473 | MEDIUM | Die DoD nennt die untere Grenze von `dupl` ausdrücklich: Sie braucht einen Fall, der sie unterscheidet, oder einen Eintrag als offen im Kopf. Mein Mutant M3 (`threshold: 149`) bleibt grün, `p4-dupl-gruen` liegt also nicht an der Schwelle. Unter OFFEN steht `dupl` nicht. Der Kopf sagt trotzdem „je Schwelle ein grüner und ein roter Fall an der Grenze“ und nennt `p4-dupl-*`. *Failure-Szenario:* Eine Senkung der Schwelle unter 150 (gegen `AGENTS.md` §3.6 nur mit ADR) läuft durch die Gegenprobe. Der Kopf behauptet, sie würde gefangen. | `AGENTS.md` §3.10, §3.11; `SPEC-049` Punkt 4; Plan §2 (DoD-Satz „die untere Grenze von `dupl`“) | `tools/harness/lint-gegenprobe.sh` · „je Schwelle ein grüner und ein roter Fall an der Grenze“; · `doppelt doppeltKurzA 45` | ja (M3 grün) | Zusage eines neuen Vertrags ohne unterscheidende Mutation |
| F-474 | MEDIUM | Zwei Kommentare sagen mehr zu, als ein Fall prüft. (a) `lint.sh` ordnet Punkt 2 „p0-grundlauf (Module aus deps ohne Netz)“ zu. Der Kopf der Gegenprobe sagt, der Grundlauf „zeigt die Module aus deps ohne Netz“. Derselbe Kopf führt `--network=none` aber unter OFFEN, und §7 nennt diesen Mutanten grün. „ohne Netz“ zeigt also kein Fall. (b) `lint.mk` sagt „Beide schreiben nichts in den Arbeitsbaum“, darunter stehen als `GEPRUEFT DURCH` nur die drei `p10-*`. Für `lint` steht die Zusage im Kopf der Gegenprobe als offen. Für die Gegenprobe selbst gibt es nur die einmalige Messung mit `git status` aus §7. *Failure-Szenario:* Jemand nimmt `--network=none` aus der Stufe, liest „p0-grundlauf … ohne Netz“ als geprüft, und das Gate bleibt grün. Dasselbe Muster war schon zweimal LOW (F-401, F-460) und einmal MEDIUM (F-465), deshalb MEDIUM nach Reviewer-Skill. | `AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `tools/harness/lint.sh` · „p0-grundlauf (Module aus deps ohne Netz)“; `tools/harness/lint-gegenprobe.sh` · „zeigt die Module aus deps ohne Netz“; `harness/mk/lint.mk` · „Beide schreiben nichts in den“ | ja (Mutant `--network=none` grün laut §7) | Kommentar-Zusage weiter als die Prüfung |
| F-475 | LOW | Die Gegenprobe hat mit `LINT_GEGENPROBE_PARALLEL` eine Einstellung (Default 6), deren Randformen weder §6 noch die Spezifikation nennen. Entschieden hat sie der Code: Mit `0` läuft die Warteschleife endlos (nachgefahren, `timeout 3` endet mit 124). Ein Wert, der keine Zahl ist, meldet `Ganzzahliger Ausdruck erwartet` und hebt die Grenze auf. Nicht MEDIUM: Die Variable steht nicht in `harness/README.md` und nicht im Hilfetext. Zugesagt ist sie nur im Kopf und in §3 („bis zu sechs gleichzeitig“). | `AGENTS.md` §3.12 („Ein Diff, der eine Randform entscheidet, die §6 nicht nennt“) | `tools/harness/lint-gegenprobe.sh` · `while [ "$(jobs -rp \| wc -l)" -ge "$parallel" ]` | ja (Lauf mit `0` hängt) | Randform einer Einstellung vom Code entschieden |
| F-476 | LOW | `7a32d57` ändert `SPEC-049` Punkt 8 nach dem ersten Code-Commit. §3 (Zeile `spec/spezifikation.md`) und der Kopf (*Berührte Spec-Stellen*) nennen weiter nur die Fortschreibung der Grenze am 2026-10-07. §1 sagt weiter, `SPEC-049` habe der Architect „vor dem Code geschrieben“. Die Entscheidung selbst steht vollständig in §6 *Randform-Rückgaben*, deshalb nur LOW: Ein Leser von §1 und §3 sieht die Änderung an Punkt 8 nicht. | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | `docs/plan/planning/in-progress/slice-harness-lint.md` · „die Grenze am 2026-10-07 um ungültiges YAML und den Cache des Builds fortgeschrieben“; · „`SPEC-049` und die Historien-Zeilen hat der Architect vor dem Code geschrieben“ | nein (Lesen) | Plan folgt der Spezifikations-Änderung nicht |
| F-477 | INFO | Zur Größe und zur nicht gezogenen Rückführung aus §4. Der Diff passt in eine Review-Sitzung: Ich habe die 1115 Zeilen der Gegenprobe ganz gelesen, jeden Fall gegen `SPEC-049` gehalten und vier Läufe gefahren. Das Skript besteht aus etwa 150 Zeilen Hilfsfunktionen und Tabellen aus Fall und Erwartung, ein Fall braucht eine bis drei Zeilen. Die Bedingung der Rückführung („nicht in einer Review-Sitzung prüfbar“) ist nicht eingetreten. Einen Schnitt halte ich nicht für nötig. F-471 bis F-473 sind Lücken in der Tabelle; das ist keine Frage der Größe. Liefer-Punkte: drei, Schichten: Harness und Doku, dazu nur ein Kommentar im Produkt-Code. | Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice; Plan §4 | `docs/plan/planning/in-progress/slice-harness-lint.md` · „Die Rückführung `in-progress → next` ist nicht gezogen“ | nein (Urteil) | — (Negativbefund) |
| F-478 | INFO | Zum Lesen eines roten Laufs. Hat der Bestand einen Befund (G1), ist neben `make lint` auch `make lint-gegenprobe` rot. Es meldet dann etwa zwanzig Fälle, die mit der Gegenprobe nichts zu tun haben, darunter `p8-form-eintrag-8` mit der ganzen erwarteten Liste als Ausgabe. Die Ursache steht im ersten Fall (`p0-grundlauf`). Weder die Sensor-Datei noch die Zeile in `harness/README.md` sagt, dass ein roter Bestand die Gegenprobe mit rot macht. Zweitens: Bricht das Skript nach dem ersten `starte` über `set -e` ab, löscht `trap` das Arbeitsverzeichnis, während Builds im Hintergrund weiterlaufen, und es erscheint keine Zeile `ROT`. Für den Implementer bzw. Architect. | Maintainability | `tools/harness/lint-gegenprobe.sh` · `trap 'rm -rf "$arbeit"' EXIT`; · `melde "Fall 'p8-form-eintrag-8' erwartet:"` | ja (G1) | Folgerot eines Gates ohne Lesehinweis |
| F-479 | INFO | Für den Verifier, zu den Schwerpunkten. **Verdrahtung:** `GATE_CHECKS += lint lint-gegenprobe`, keine Ordnungskante, `record-gates` hängt an allen. **Rotes Lint macht `make gates` rot:** Bei G1 meldet die Stufe `lint: internal/hexagon/model/fehler.go:175: Direktive nolint`, `make` meldet `[harness/mk/lint.mk:15: lint] Fehler 1`, `make -k gates` endet mit 2. Rot war in meiner Kopie auch `hook-gegenprobe`, weil ihr `.git` fehlt; das ist ein Artefakt der Kopie. **Eigene Zeile je rotem Fall:** Jeder Fall mit `befund`, `lint_zeile`, `form` oder `allein` prüft Pfad und Zeile bzw. den Wortlaut. Nur am Ausgang hängen `p9-laden` (mit `can't load config` und ohne `lint:`-Zeile, wie die Grenze es sagt), `p10-gates-rot` (am selben Log prüft `p9-lint-zeile-allein` die Zeile), `p10-gegenprobe-scharf` (Ersatzskript) und `p0-sammel`. **Kopien:** `mktemp -d`, `cp -R` ohne `-p`, `touch` nach jeder Änderung an einer Datei im Build-Kontext. Gelöscht wird nur `"$arbeit"` sowie per `find … -delete` in `"$k/harness/mk"` der eigenen Kopie. **Grundlauf:** grün. **OFFEN:** Gegen `SPEC-049` und die Rückgaben vollständig bis auf F-471 bis F-473. Pin, Plattform und `Felder nicht erkannt` stehen dort mit dem Grund aus §6. **`lint.sh` in `3bcf000`:** unverändert, im ganzen Diff nur der Kopfkommentar. Die Form nach Punkt 8 neu (`sub(/[ \t]+$/, "", rest)`, `/^-( \|$)/`) entspricht der Fassung aus `7a32d57`. **`fehler.go`:** Der Kommentar sagt die Tiefensuche zu. Der Fall `Tiefensuche` in `fehler_test.go` prüft, dass ein tiefer Fehler der ersten Ursache vor einem flachen der zweiten kommt. | Plan §6 *Prüfung vor dem Code*; `SPEC-049` Punkt 8 bis 10 | `harness/mk/lint.mk` · `GATE_CHECKS += lint lint-gegenprobe` | ja (G1, Lesen) | — (Negativbefund) |

**Eigene Mutationen** (jede in einer frischen Kopie an `3bcf000`; die unveränderte Kopie ist grün):

| ID | Mutation | Lauf | Ergebnis |
|---|---|---|---|
| M1 | `tools/harness/lint.sh`: Prüfung nach Punkt 6 ohne Testdateien, `done < <(go_dateien '*.go' \| grep -v '_test\.go$')` | `make lint-gegenprobe` | grün (F-471) |
| M2 | `.golangci.yml`: `skip-regexp: export_test\.go$` (ohne `(^\|/)`) | `make lint-gegenprobe` | grün (F-472) |
| M3 | `.golangci.yml`: `dupl` `threshold: 149` | `make lint-gegenprobe` | grün (F-473) |
| G1 | `internal/hexagon/model/fehler.go`: `//nolint` am Ende von `func Meldungen(…) {` | `make -k gates` | rot: `lint` mit der Zeile der Direktive, `lint-gegenprobe` mit `p0-grundlauf` und Folgefällen (F-478, F-479) |

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/lint-gegenprobe.sh` | geprüft. Befunde F-471 bis F-475, F-478. Arbeit nur in Kopien, Aufräumen nur der eigenen Pfade, Grundlauf grün. |
| `tools/harness/lint.sh` | geprüft. Nur der Kopf geändert; Befund F-474 (a). Kein Verhalten geändert, wie §3 sagt. |
| `harness/mk/lint.mk` | geprüft. Verdrahtung richtig (F-479); Befund F-474 (b). |
| `Dockerfile` | geprüft, ohne Befund. Nur der Kommentar der Stufe `lint`, der Pin ist unverändert. |
| `harness/sensors/lint.md` | geprüft, ohne Befund. Vertrag als Link auf `SPEC-049`. Die Grenze gibt die der Spezifikation wieder, die Sensor-Datei entscheidet nichts neu. Zum Folgerot siehe F-478. |
| `harness/README.md` | geprüft, ohne Befund. `make lint` steht unter §Sensors, `make lint-gegenprobe` ist neu, „Nicht behauptet“ nennt nur noch Testabdeckung, die Zeile `make lint` ist aus den Werkzeugen entfernt. |
| `AGENTS.md` §3.2 | geprüft, ohne Befund. Der Träger ist genannt, Falsch/Richtig sind aus diesem Repo, und es steht kein Target darin, das es nicht gibt. |
| `internal/hexagon/model/fehler.go` | geprüft, ohne Befund. Nur der Kommentar ist geändert; die Zusage hält ein Test (F-479). |
| `spec/spezifikation.md` | geprüft, ohne Befund. Die Änderung in `7a32d57` steht in Punkt 8 und in der Historie. Die Historie nennt weder ADR noch Slice. |
| `docs/plan/planning/in-progress/slice-harness-lint.md` | geprüft. Befunde F-476, F-477. |
| Hard Rule 3.2 / Suppression | geprüft, ohne Befund. Der Diff enthält kein `//nolint`, und `.golangci.yml` ist unverändert. |
| Hard Rule 3.3 | geprüft, ohne Befund. Im Diff gibt es keinen Move, nur zwei neue Dateien. |
| Hard Rule 3.6 | geprüft, ohne Befund. Keine Schwelle ist gesenkt, das Profil ist unverändert. |
| Hard Rule 3.7 (Kommentar-Klassen) | geprüft, ohne Befund. Die neuen Kommentare beschreiben den Ist-Zustand. Zu §3.11 siehe F-474. |
| Hard Rule 3.12 / MEDIUM *Randform im Code-Commit* | geprüft, ohne Befund für das MEDIUM. Die Rückgaben aus `221845d` hat der Architect in `7a32d57` entschieden, einem Commit ohne Code. `3bcf000` ändert §6 nicht. Zur Einstellung `LINT_GEGENPROBE_PARALLEL` siehe F-475. |
| Hard Rule 3.13 / MEDIUM *Adresse nimmt nicht an* | geprüft, ohne Befund. Der Diff nennt keinen Slice neu als Adresse. |
| Core-Reinheit, [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | geprüft, ohne Befund. In `internal/hexagon/` hat sich nur ein Kommentar geändert. |
| Commit-Messages `221845d`, `7a32d57`, `3bcf000` | geprüft, ohne Befund. Jede nennt `slice-harness-lint` und [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), keine nennt eine Struktur-ID. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 2 |
| INFO | 3 |

**Summary:** 0 HIGH · 4 MEDIUM · 2 LOW · 3 INFO (F-471 `nolint` in Testdateien ohne Fall, M1 grün; F-472 Anker von `skip-regexp` ohne Fall, M2 grün; F-473 untere Grenze von `dupl` weder Fall noch offen, M3 grün; F-474 Kopf von `lint.sh` und `lint.mk` sagen „ohne Netz“ bzw. „nichts in den Arbeitsbaum“ als geprüft zu; F-475 Randformen von `LINT_GEGENPROBE_PARALLEL` vom Code entschieden; F-476 §1, §3 und Kopf folgen der Änderung an Punkt 8 nicht; F-477 Größe trägt, keine Rückführung nötig; F-478 Folgerot der Gegenprobe ohne Lesehinweis; F-479 Verdrahtung und Hygiene der Gegenprobe bestätigt). Wiederkehrende Klassen: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-471 bis F-473), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-474), `BEO-REPO/plan-folgt-korrektur-nicht` (F-476).

**Finding-Klassen dieses Laufs:** Zusage eines neuen Vertrags ohne unterscheidende Mutation · Kommentar-Zusage weiter als die Prüfung · Randform einer Einstellung vom Code entschieden · Plan folgt der Spezifikations-Änderung nicht · Folgerot eines Gates ohne Lesehinweis

## Verdikt

**Merge-blockierend:** nein, aber F-471 bis F-474 brauchen vor der Closure einen Ausgang. Das Gate hängt richtig an der Kette und wird bei rotem Lint rot. Die Gegenprobe arbeitet sauber in Kopien, ihr Grundlauf ist grün. Drei Zusagen von `SPEC-049` fängt sie aber nicht, und keine davon steht im Kopf als offen. Das ist die Klasse, die die DoD mit „je mit einem Fall, der sie unterscheidet, oder benannt als offen im Kopf“ ausschließen will.

**Übergabe:** F-471 bis F-474 gehen an den Implementer: je ein Fall oder ein Eintrag unter OFFEN, die Kommentare enger. Ob ein Eintrag als offen trägt, entscheidet bei Bedarf der Architect. F-475 geht an den Architect (Randform), F-476 an den Implementer bzw. Architect (Plan nachziehen). F-477 bis F-479 gehen an den Verifier und an den Planner (Größe). Dieser Report ersetzt keine Verifikation.
