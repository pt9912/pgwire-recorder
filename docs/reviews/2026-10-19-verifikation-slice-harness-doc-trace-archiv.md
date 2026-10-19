# Verifikationsbericht: slice-harness-doc-trace-archiv — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“ gegen Plan, DoD und `SPEC-068`; kein Review, kein Umbau).

**Gegenstand:** Zweig `ver-dt`, Stand `9b57139d` (= `impl-dt` nach Review-Runde 1), Basis `main` `5f68ebef`, Basis der Metrik `dc7f0a69`. Eingang: Plan `docs/plan/planning/in-progress/slice-harness-doc-trace-archiv.md` (DoD, §1, §3, §6, §7), `harness/vertraege/SPEC-068-archiv-bezug.md`, Review-Report `docs/reviews/2026-10-19-review-slice-harness-doc-trace-archiv.md`.

**Anlass der Schärfe:** Der Reviewer hat die Behebungen von F-1 bis F-10 (Commit `4dec8155`) nicht erneut gesehen (`nacharbeit-nach-review-ohne-erneuten-lauf`, Zähler 3). Diese Verifikation ist der erneute Lauf: jede Behebung wurde von Hand gebrochen und der Fall auf den richtigen Grund geprüft.

## Urteil

**Konform, ohne HIGH und ohne MEDIUM.** Alle drei Liefer-Punkte der DoD sind am Stand geliefert und belegt, die Behebungen F-1 bis F-10 tragen (je Bruch rot aus dem richtigen Grund). Zwei LOW-Befunde und drei INFO stehen unten; keiner blockiert. Offen sind allein die Closure-Pflichten, die dem Planner gehören.

## 1. Behebungen F-1 bis F-10: Bruch von Hand

Methode (Modul 11, Bewusstes Brechen): den Fix im Klon entfernen, `GP_MUT=ok bash tools/harness/archiv-bezug-gegenprobe.sh` (unveränderte Kopie liest die Datei im Repo) fahren, die roten Fälle lesen, `git checkout` zurück. Die Brüche sind eigene Textersetzungen, nicht die Mutationen der Gegenprobe; danach lief zusätzlich die Gegenprobe der Gruppen `m-p4-art|m-p3-schreiben|m-p3-temp|m-p3-grund|m-p3-erster|m-p3-mv` (19 Läufe, grün).

| F | Bruch (Fix entfernt) | rote Fälle | Urteil |
|---|---|---|---|
| F-1 | Rückgabewert von `schreibe` nicht gelesen (`|| { rm -f …; exit 2 }` entfernt) | `p3-schreiben-scheitert`, `p3-schreiben-temp-aufgeraeumt` | trägt |
| F-1 von Hand | Attrappe `mv` scheitert für `slice-b.md`: Exit 2, Meldung `…/slice-b.md: Schreiben fehlgeschlagen`, Verzeichnis ohne Temp-Datei (`archiv.zip slice-a.md slice-b.md slice-c.md`), `slice-a.md` geschrieben, `slice-b.md` und `slice-c.md` nicht; zweiter Aufruf ohne Attrappe Exit 0 setzt fort, dritter Exit 0 ohne Änderung (`md5sum`) | — | trägt |
| F-2 (b)/(c): CR vor Kopf | `p3-reihenfolge-grund-b-c` | trägt |
| F-2 (c)/(d): `Bezug`-Anzahl vor CR | `p3-reihenfolge-grund-c-d` | trägt |
| F-2 innerhalb (c): Hervorgegangen vor CR | `p3-reihenfolge-grund-c-c` | trägt (Vertragslücke, siehe V-2) |
| F-2 (d)/(e): Beschreibbarkeit vor Abweichung | `p3-nicht-beschreibbar-gleich`, `p3-reihenfolge-grund-d-e` | trägt |
| F-3 | Datei `extra.log` in `done/`; in der Wurzel; `touch archiv.zip`; `touch` aller Stubs; Nebendatei nur beim Schreiben | je `p3-nur-stubs` (alle fünf) | trägt |
| F-4 | Attrappe `id -u` = 0: alle Mutationen `m-p3\|m-p2\|m-stichwort\|m-name`: grün (60 Läufe); Attrappe `locale` ohne `en_US`: `m-p2\|m-stichwort\|m-name\|m-p8`: grün (46 Läufe); beide Attrappen, volle Gegenprobe: grün (120 Läufe) | — | trägt |
| F-4 Bruch | Filter der erwarteten Liste (`grep -vxF -f …ENTFALLEN`) entfernt, Attrappe root | `m-p3-schreiben-vor-pruefung` und `m-p3-kein-abbruch-bei-befund` melden ROT mit erwartet ≠ rot (Symptom des Reviews reproduziert) | trägt |
| F-5 | Namen aus den Meldungen entfernt: `archiv.zip beschädigt`, `Form der Welle`, `Verzeichnis fehlt`, `Werkzeug fehlt`, `archiv.zip fehlt`, Wurzel | je zwei bis sieben benannte Fälle rot, nie keiner | trägt |
| F-6 | Kennungsart ohne Mutation: Regex ohne `QA` bzw. `RB` bzw. `FA` (Mutationen `m-p4-art-*` in der Gegenprobe) | 10, 2, 34 Fälle rot, wie in §7 | trägt |
| F-10 Wurzel | Prüfung der Wurzel gestrichen | `p2-wurzel-fehlt`, `p2-wurzel-leer`, `p2-wurzel-kein-verzeichnis`, `p2-reihenfolge-0-1` | trägt |
| F-10 falsche Stelle | `Zeile an falscher Stelle` zu `\|\| true` | `p6-gleich-falsche-stelle` | trägt |
| F-10 Welle ohne Stub | `compgen -G`-Filter entfernt (jede Welle zählt) | `p1-eine-ohne-stub`, `p1-leere-menge`, `p1-ohne-stub`, `p2-welle-ohne-stub-ohne-archiv` | trägt |

F-7, F-8, F-9 (Sätze): geprüft gegen den Text. `harness/README.md` nennt „die Fälle mit Dateirechten (Liste ENTFALLEN im Skript)“ und die Locale; das Fragment sagt „Punkt 1 bis 12“; `SPEC-068` Plattform sagt „Gegenprobe ≥ 4.3 wegen `wait -n`“ und führt `zip` in der Gegenprobe; der Plan-Kopf nennt die Messung (sieben vorher, nachher keine); §7 nennt „3 Einfügungen, 1 Löschung“ für die README (stimmt mit `git diff --numstat`).

## 2. Die 97 Stubs, unabhängig geprüft

- `git diff 5f68ebef..HEAD -- docs/plan/planning/done`: 97 Dateien, alle `M`, je `1 0` (numstat), alle 97 Pfade `done/welle-*/slice-*.md`, einzige hinzugekommene Zeile in jeder Datei `**Bezug:** …`, keine Löschung, keine andere Änderung, kein „No newline“. 97 von 97 `slice-*.md` im Bestand tragen genau eine Zeile `**Bezug:**`.
- Neuableitung mit einem eigenen Perl-Skript (nicht das Werkzeug; Kopf bis zur ersten Zeile `---`, erste Zeile `**Bezug:**` mit Fortsetzung bis zur nächsten Zeile `**`, Token `LH-(FA|QA|RB)-n` ohne Vorzeichen-Zeichen, Ordnung Gruppe, Zahl, Länge, ohne Duplikat; Quelle `unzip -p` aus den vier `archiv.zip`): **97 von 97 gleich**, jede Zeile steht unmittelbar hinter der einzigen Zeile `**Hervorgegangen:**`. Verteilung der Kennungen je Stub: 0 → 12, 1 → 27, 2 → 9, 3 → 24, 4 → 15, 5 → 4, 6 → 4, 7 → 2 (wie in §7: 12 Leerwerte). Stichprobe der Leerwerte (`slice-harness-upgrade-v6-18`: Bezug `MR-000`; `slice-harness-kopf-sensor`: „— (Harness-Arbeit; keine Produkt-Anforderung)“) gegen den Volltext gelesen.
- Idempotenz: `make archiv-bezug` am Stand: Exit 0, `git status` leer. `make -n archiv-bezug` druckt das Rezept (`bash tools/harness/archiv-bezug.sh "<Wurzel>" ""`).

## 3. Matrix und bestehende Gates

| Messung | Ergebnis |
|---|---|
| `make doc-trace` im Worktree von `5f68ebef` | 34 Anforderung(en), **7** Waise(n) |
| `make doc-trace` am Stand `9b57139d` | 34 Anforderung(en), **0** Waise(n) (LH-FA-19 zählt über `slice-erster-release-veroeffentlichung`, wie §7) |
| `make docs-check` | 1222 Datei(en), 0 Befunde (der Implementer nannte 1221 am Stand `4dec8155`; die Differenz ist die Datei des Review-Reports bzw. dieses Berichts im Scan, kein Befund) |
| `make kopf-check`, `make abdeckung-check` | Exit 0 |
| `bash -n` beider Skripte | grün |
| Stop-Hook (`{"stop_hook_active": false}`) | Stand `9b57139d`: `approve`; Zeile an `tools/harness/archiv-bezug.sh` angehängt: `block` („A changed path is not on the inert list“); zurückgenommen: `approve` |
| volle Gegenprobe `make archiv-bezug-gegenprobe` (mitgefahren, 1 min 5 s) | grün, 127 Läufe (126 Mutationen und die unveränderte Kopie); 91 Fälle (`^f_p`), 123 `mut`-Zeilen plus 3 aus der Schleife von Punkt 11 |
| §7-Tabelle gegen das Skript | 126 Mutationsnamen in der Tabelle, 126 im Skript (123 + 3), kein Name nur auf einer Seite |
| Metrik `git diff --numstat dc7f0a69..HEAD -- tools harness` | **625** (gleich am Stand `4dec8155`; Marke 800, Meldeschwelle 650 nicht erreicht) |

`make gates` und `make gates-schnell-gegenprobe` am Stand `4dec8155` nicht wiederholt (Beleg des Implementers, §7: grün, 4 min 33 s bzw. 1 min 6 s); seit `4dec8155` hat sich nur der Plan geändert (`git diff --stat 4dec8155..HEAD`: eine Datei, `slice-harness-doc-trace-archiv.md`), und der Stop-Hook bestätigt den Nachweis (`approve`).

## 4. DoD gegen Beleg

| DoD-Punkt | Befund |
|---|---|
| Liefer-Punkt 1: Werkzeug, Fragment, `SPEC-068`, Index/Historie | `tools/harness/archiv-bezug.sh`, `harness/mk/archiv-bezug.mk` (`.PHONY`, Rezept ohne `$(MAKE)` und ohne `+`, nicht an `GATE_CHECKS`), `SPEC-068`; `harness/vertraege/README.md` führt Index (Zeile 45) und Historie (Zeile 80) seit dem Vertragscommit `19799c99`, der vor der Basis `5f68ebef` liegt. Erfüllt. |
| Liefer-Punkt 2: Gegenprobe | Fragment in `GEGENPROBEN_CONTAINER` (Fall `p11-keine-kette` prüft, dass es weder in `GATE_CHECKS` noch in einer der beiden Gegenproben-Listen für das Werkzeug selbst steht), 91 Fälle, 126 Mutationen, Zeile im Gate-Index von `harness/README.md`. Erfüllt. |
| Liefer-Punkt 3: Anwendung und Berichtigung | vier Inhalts-Commits (3, 4, 5, 85 Stubs = 97), Waisen 7 → 0, Zeile `make doc-trace` berichtigt (nennt die Regel an den Schreiber und zeigt auf `trace:` in `.d-check.yml`). Erfüllt. |
| Mutation je Zusage, vom Implementer rot gesehen | Tabelle in §7 vollständig (siehe oben); Stichproben (Abschnitt 1) rot aus dem richtigen Grund. Ausnahme: V-1. |
| Wirkung auf bestehende Gates | gemessen (Abschnitt 3), wie in `SPEC-068` Punkt 11 (c) vorhergesagt. |
| `make gates`, `make archiv-bezug-gegenprobe`, `make gates-schnell-gegenprobe` | Beleg in §7; Gegenprobe selbst nachgefahren (grün). |
| Review, Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | **offen**, gehören dem Planner (V-3). |

DoD-Behauptungen, die sich auf einen Test berufen: jede Stichprobe oben ist „ohne Fix rot“ gezeigt (Modul 11).

## 5. Entscheidungs-Konformität (Plan §1, §6, `SPEC-068`)

- §1 Abgrenzung gehalten: keine Änderung an `.d-check.yml`, kein Gate, keine Aufnahme in `GATE_CHECKS`, keine ADR-, `SPEC-`- oder `ARC-`-Kennung in den Stubs (Fälle `p12-*`), Wellen-Dateien und Review-Reports unberührt (`p1-welle-flach`, `p1-welle-datei`). Das Diff berührt außerhalb von `tools/`, `harness/`, `docs/plan/planning/` nichts.
- §6 Entscheidungen des Architect umgesetzt: Modus „prüfen“ nein; Matrix-Wirkung als Messung in §7; eigene Zeile `**Bezug:**` hinter `**Hervorgegangen:**`, keine `MR`; Randformen je Punkt im Code wie im Vertrag (Reihenfolge der Gründe 0 bis 5 über den ganzen Aufruf, Prüfung aller Stubs vor dem ersten Schreiben, Schreibfehler Exit 2, Wurzel, leere Menge, falsche Stelle, `LC_ALL=C`, CR).
- §3.3: reine Inhalts-Commits, kein `git mv`, jeder Commit des Bereichs nennt die Slice-Kennung (0 Commits ohne).
- §3.13 Nachzählung in §8: drei Liefer-Punkte, eine Schicht; hält.

## Findings

| ID | Schwere | Befund | Beleg | Pfad |
|---|---|---|---|---|
| V-1 | LOW | Zusage „Stub hat keinen Eintrag **gleichen Pfads** im Archiv“ (`SPEC-068` Punkt 3 (a)) hat keinen Fall für einen fast gleichen Pfad. Der Mutant `grep -Fxq --` → `grep -Fq --` (Teilzeichenkette statt ganzer Zeile) bleibt grün; §7 führt ihn weder als Mutation noch als äquivalent mit Invariante, Ort und Sonde (`AGENTS.md` §3.10, grüner Mutant). Gemessen: Archiv mit nur `…/slice-a.md.orig`, Stub `slice-a.md`: Original `kein Eintrag im Archiv`, Mutant `Eintrag im Archiv ohne Kopf` (Stichwort falsch, Exit 1 in beiden, nichts geschrieben); zusätzlich Meldung `caution: filename not matched` von `unzip` auf stderr. | Handmessung im Temp-Baum; `GP_MUT=ok` mit dem Mutanten: „grün (1 Läufe)“ | `tools/harness/archiv-bezug.sh` · `grep -Fxq -- "$pfad" "$tmp/liste"`; Gegenprobe `f_p3_stub_ohne_eintrag` |
| V-2 | LOW | Die Reihenfolge **innerhalb** von (c) (CR vor „nicht genau eine Zeile `**Hervorgegangen:**`“) steht in keinem Vertragstext: Punkt 3 nennt (a) bis (e), Punkt 7 und 8 nennen beide als (c) ohne Rang. Die Gegenprobe legt sie mit `p3-reihenfolge-grund-c-c` und `m-p3-grund-hervorgegangen-vor-cr` fest (Randform im Test, nicht im Vertrag, `AGENTS.md` §3.12). Der Bruch von Hand färbt genau diesen Fall; der Vertrag deckt ihn nicht. Für den Architect: im Vertrag entscheiden oder den Fall entfernen. | Bruch „Hervorgegangen vor CR“ in Abschnitt 1; `grep` im Vertrag auf „innerhalb (c)“ ohne Treffer | `harness/vertraege/SPEC-068-archiv-bezug.md` Punkt 3 (c); `tools/harness/archiv-bezug-gegenprobe.sh` |
| V-3 | INFO | Closure-Pflichten offen (erwartet vor der Closure): Häkchen der DoD, `§7` Felder (`Was hat funktioniert`, `Steering-Loop-Eintrag`, `Beobachtungs-Register`, `Folge-Slices`, `Risiken aus §6`, `Drei Paarungen`) stehen auf `—`; die vier Risiken aus §6 tragen noch „**Ausgang:** —“. Der Plan schlägt vor: Risiko 1 *weiter offen* → `BEO-REPO/archiv-stub-verliert-anforderungs-kennung` (der Vertrag Punkt 12 nennt sie), Risiko 2 mit der Messung in §7 *entfallen*, Risiko 3 *entfallen* (dieselbe Lesart wie bei lebenden Slices), Risiko 4 *entfallen* (Punkt 3 und Grenze im Vertrag; der Fall `p3-schreiben-scheitert` zeigt den halben Stand und die Fortsetzung). Die Finding-Klassen des Reviews gehen in den Lerneintrag. Die Beobachtung `nacharbeit-nach-review-ohne-erneuten-lauf` ist mit diesem Lauf nicht „wieder aufgetreten“, sondern bedient: die Nacharbeit wurde erneut gesehen. | Plan §6, §7 | `docs/plan/planning/in-progress/slice-harness-doc-trace-archiv.md` |
| V-4 | INFO | Fallnamen im Vertrag und in der Gegenprobe weichen ab: der Vertrag nennt `p3-nicht-beschreibbar` und `p2-form-welle`, die Gegenprobe führt `p3-datei-nicht-beschreibbar`, `p3-verzeichnis-nicht-beschreibbar` und fünf `p2-form-welle-*`. Die Zusagen sind gedeckt (Datei und Verzeichnis je ein Fall), die Namen im Vertrag sind Beispiele. Keine Wirkung auf Gates. | Vergleich Vertrag Punkt 2 und 3 gegen `^f_p` | `SPEC-068` · Gegenprobe |
| V-5 | INFO | `set -o pipefail` im Werkzeug hat keine beobachtbare Wirkung (Mutant `:` statt `set -o pipefail` bleibt grün): Die einzige Pipeline (`kennungen`) läuft in einer Prozess-Ersetzung, deren Status niemand liest, und kein weiterer Schritt verlässt sich auf `pipefail`. Kein Vertrags-Satz berührt; als Kommentar-Zusage steht nichts darüber. Für den Implementer ohne Eile: entfernen oder begründen. | Handmessung `GP_MUT=ok` | `tools/harness/archiv-bezug.sh` · `set -o pipefail` |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Stubs: Inhalt, Stelle, nur eine Zeile, keine andere Änderung | 97 von 97 unabhängig bestätigt, ohne Befund |
| Matrix: 7 Waisen → 0 | nachgefahren (Worktree `5f68ebef` gegen `9b57139d`), ohne Befund |
| Plan §7 gegen Skript: Tabelle, Zählung (91 Fälle, 126 Mutationen, 127 Läufe), Metrik 625 | ohne Befund |
| Wirkung auf `docs-check`, `kopf-check`, `abdeckung-check`, Stop-Hook | ohne Befund |
| Entfallen-Pfad der Gegenprobe (root, ohne Locale) | trägt, ohne Befund |
| Temp-Verzeichnis des Werkzeugs bei Exit 2 (Wurzel fehlt) | aufgeräumt (Zahl der `tmp.*` in `/tmp` vor und nach dem Lauf gleich) |
| Idempotenz: zweiter und dritter Lauf nach Schreibfehler und Fortsetzung | Exit 0, Dateien byte-gleich (`md5sum`), ohne Befund |

## Summary

| Schwere | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 (V-1, V-2) |
| INFO | 3 (V-3, V-4, V-5) |

**Übergabe an den Planner:** Der Slice ist DoD-konform in den Liefer-Punkten und trägt die Behebungen aus Review-Runde 1. V-1 (Fall und Mutation für den fast gleichen Pfad) und V-5 gehen an den Implementer, V-2 an den Architect (Vertrag). Sie blockieren die Closure nicht, wenn der Planner sie als Folgepunkt führt oder vorher behoben haben will. Die Closure-Pflichten (V-3) gehören dem Planner. Dieser Bericht ist ein Lauf-Beleg.
