# Verifikation: slice-harness-upgrade-v6-18 — 2026-10-10

**Rolle:** Verifier (Modul 11). Geprüft wird, ob der Slice Plan und DoD erfüllt (bauen wir es richtig?). Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft; ob das Richtige gebaut wird, prüft der Validator.

**Gegenstand:** Slice-Plan `slice-harness-upgrade-v6-18` (Kopf, §1 bis §8) gegen den Gesamt-Diff `5e25195..48a6b00`. Architect: `954fc1b`. Implementer: `420f955`, `35503be`, `0e61bed`, `ddb288d`. Review: `48a6b00`. Rahmen: [MR-000](../../harness/conventions.md#mr-000--baseline-aussage).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `48a6b00`, besonders §1 (*Ausdrücklich NICHT*), §2 (DoD), §6 (*Entscheidungen des Architect*) und §7 (*Belege des Implementers*)
- der Review-Report zu `48a6b00` (F-576 bis F-579); den Bericht des Implementers habe ich nicht gelesen
- `AGENTS.md` §3, `harness/conventions.md`, `.harness/skills/reviewer.md`, `.claude/commands/implement-slice.md`
- `v6.18.0` · `regelwerk/modul-05-planning-harness.md`, `regelwerk/modul-10-review-harness.md`, `regelwerk/modul-13-quality-gates.md`, `regelwerk/grundlagen-harness-dateien.md`

Beim Start war der Arbeitsbaum sauber, HEAD `48a6b00`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Scratch außerhalb des Repos, Präfix `ver-upg-`: Release-Asset per `gh release download v6.18.0 -R pt9912/ai-harness-course` in ein eigenes Verzeichnis; für die Gegenproben eine Kopie des Arbeitsbaums (`tar` ohne `.git`, ohne Eigentümer und Zeitstempel) und je Fall eine frische `cp -r`-Kopie davon. Im Repo wurde nichts verändert.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Baseline `v6.18.0` vendored

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Asset sha256 vor dem Entpacken gegen die `SHA256SUMS` des Releases geprüft, Wert und Quelle in §7 | §7 nennt `18e5443b…b8d9`, Rang 1 (Release-`SHA256SUMS`) und Rang 2 (API-Digest). Eigener Abruf: `sha256sum -c` der Release-`SHA256SUMS` über das frisch geladene Zip meldet OK, derselbe Wert | bestätigt |
| Zip enthält nur `regelwerk/` und `templates/`, keine absoluten Pfade, kein `..`, keine Symlinks | `unzip -l` ohne Treffer auf `^/` und `..`; `find -type l` im entpackten Baum leer; 54 Dateien | bestätigt |
| Baum `.harness/baseline/v6.18.0/{regelwerk,templates}/` entspricht dem Asset | `diff -r` entpacktes `regelwerk/` und `templates/` gegen den vendored Baum: leer | bestätigt |
| vendored `SHA256SUMS` vollständig und richtig | `make baseline-verify`: `v6.18.0 OK — 54 Dateien`; Datei-Anzahl im Asset und in der Liste 54 | bestätigt |
| `v6.16.0` im selben Commit entfernt | `git diff -M 5e25195..48a6b00`: alle Pfade als Rename `v6.16.0 => v6.18.0` in `420f955`; `ls .harness/baseline` nennt nur `v6.18.0` | bestätigt |
| Pins in `harness/conventions.md` (Stand, Datum, Asset-URL, Stand-Zeile „Kurs-Welle 161 · 2026-10-10“, Messzeile) und `AGENTS.md` §1 | `git diff`: alle fünf Stellen nennen `v6.18.0`; die Stand-Zeile stimmt mit `regelwerk/README.md` des neuen Baums | bestätigt |
| `grep "v6\.16\.0"` über `AGENTS.md`, `README.md`, `harness/` (ohne `harness/conventions/done/`), `.claude/`, `.harness/skills/`, `tools/`, `Makefile`, `*.mk`, `spec/` leer außer der Zeile `MR-001` | eigener Lauf: ein Treffer, `harness/conventions.md` Zeile 206 (`MR-001`, Stand der Auflösung, laut DoD ausgenommen) | bestätigt |
| beide Symlinks unter `.claude/rules/` lösen auf `v6.18.0` auf; `find -xtype l` leer | `ls -l` zeigt beide Ziele `.harness/baseline/v6.18.0/…`; `find . -path ./.git -prune -o -xtype l -print` leer | bestätigt |

Treffer auf `v6.16.0` außerhalb der DoD-Grenzen (`done/`-Pläne, `docs/reviews/`, Register-Belege, `open/slice-harness-gate-index-werkzeug-teil`, Drift-Log der Roadmap) sind Chronik und nach §6 bewusst stehengelassen.

**Urteil Punkt 1: bestätigt.**

### Punkt 2: Abgleich `v6.16.0` → `v6.18.0`

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Jeder Regelblock des neuen Baums steht mit genau einem Wert in der Tabelle | `ls` des neuen `regelwerk/` gegen `harness/conventions.md`: alle 25 `modul-*`/`grundlagen-*`-Dateien und `README.md` stehen je mindestens einmal; drei Dateien (`grundlagen-harness-dateien`, `modul-02`, `modul-15`) je zweimal nach Abschnitt, wie vor dem Sprung | bestätigt |
| neue, umbenannte, entfallene Blöcke einzeln genannt | `git diff -M` der Baumpfade: keine neue, keine entfallene, keine umbenannte Datei; §7 sagt das ausdrücklich | bestätigt |
| Delta gegen verkörperte Regeln gehalten | Inhaltliches Delta nach `git diff -M`: `modul-10` (Absatz zum Deckungs-Sensor, Dateiname mit voller Slice-Kennung), `modul-13` (ein Satz zur Disjunktheit), `grundlagen-harness-dateien` (ein Satz zur Disjunktheit), drei Vorlagen. Die übrigen Regelwerk-Dateien tragen nur die Versionszeile. Keine Aussage in `AGENTS.md` §3, den Commands, Rollen-Typen oder Skills widerspricht dem Delta | bestätigt |
| Adaptions-Durchgang (`MR-000`, aufgelöste `MR-001`) | `MR-000` unverändert („keine inhaltlichen Adaptionen“); das Delta berührt `grundlagen-referenz-richtung.md` nur in der Versionszeile | bestätigt |
| Die drei Punkte der Vorprüfung je mit Ausgang in §7 | §7: Disjunktheit (kein Gate-Index in Teilen, kein Befund), Modul `reviews` (keine verkörperte Regel berührt, Modul bleibt aus), Dateiname (akzeptierte Abweichung, `-r2` nachgezogen) | bestätigt |
| Dateiname gegen Skill und Bestand | Wortlaut von `modul-10` verlangt „die volle Slice-Kennung im Dateinamen“, verbietet keinen Rollen-Infix; der Bestand unter `docs/reviews/` trägt die volle Kennung. Skill und `implement-slice` nennen die Folgelauf-Form `-r2`, `-r3` | bestätigt |
| jeder Befund hat genau einen Ausgang | Der Abgleich fand keinen Befund außer der akzeptierten Abweichung; die Ausgänge der Review-Findings F-576 und F-578 stehen noch aus (siehe V-154) | bestätigt, Ausgänge der Review-Findings offen bis Closure |
| Stichprobe gegen den Bestand | `modul-07` §Ziel-Form: Carveout im Komplement des Deltas (nur Versionszeile); `docs/plan/carveouts/` enthält nur `.gitkeep`, `harness/README.md` nennt „bisher keiner“. Plausibel und nachvollziehbar | bestätigt |

**Urteil Punkt 2: bestätigt.**

### Punkt 3: `make gates` grün

Eigener Lauf am Stand `48a6b00` (Arbeitsbaum sauber): Exit 0; darin `baseline-verify: v6.18.0 OK — 54 Dateien`, `d-check: 474 Datei(en) geprüft, 0 Befund(e)`, `a-check-negativ`, `lint-gegenprobe`, Unit- und Integrationstests grün.

**Urteil Punkt 3: bestätigt.**

### Punkte 4 bis 8: Closure-Pflichten

Review-Report liegt vor (`48a6b00`, 0 HIGH/MEDIUM/LOW, 4 INFO): bestätigt. Die Häkchen in §2, die Closure-Notiz mit Lerneintrag, die Register-Fortschreibung, die Risiko-Ausgänge in §6 (alle drei stehen auf „offen bis Closure“) und die drei Paarungen sind zum Prüfstand noch nicht gesetzt. Das ist vor der Verifikation der erwartete Stand (Closure folgt nach `AGENTS.md` §3.3 danach, siehe V-154), kein Mangel.

**Urteil Punkte 4 bis 8: Review bestätigt; der Rest ist Closure und noch zu tun.**

---

## 2. Bewusstes Brechen: `make baseline-verify`

Je Fall eine frische Kopie des Arbeitsbaums im Scratchpad, `make baseline-verify`. Grundprobe in der unveränderten Kopie: grün (`v6.18.0 OK — 54 Dateien`).

| Fall | Eingriff in der Kopie | Ergebnis |
|---|---|---|
| (a) geänderte Datei | eine Zeile an `regelwerk/modul-05-planning-harness.md` angehängt | rot, Exit 2, „weicht von SHA256SUMS ab (geänderte oder fehlende Datei)“ |
| (b) fehlende Datei | `regelwerk/modul-07-carveouts.md` gelöscht | rot, Exit 2, dieselbe Meldung |
| (c) zusätzliche Datei | `regelwerk/extra.md` angelegt | rot, Exit 2, Diff-Ausgabe nennt `> regelwerk/extra.md` |
| (d1) SHA256SUMS-Zeile fehlt | Zeile 1 (`regelwerk/README.md`) entfernt | rot, Exit 2, Ausgabe nennt `> regelwerk/README.md` |
| (d2) SHA256SUMS-Zeile abweichend | erstes Zeichen der ersten Summe geändert | rot, Exit 2, „weicht von SHA256SUMS ab“ |
| (d3) SHA256SUMS fehlt | Datei gelöscht | rot, Exit 2, „SHA256SUMS fehlt“ |
| zusätzlich: zweites Tag-Verzeichnis | Kopie nach `v6.16.0` | rot, Exit 2, „mehr als ein <tag>-Verzeichnis“ |

Gegenprobe aus dem Plan (`find . -path ./.git -prune -o -xtype l -print`) im Repo: leer. Die Aussage in §7, die Gegenprobe der Baseline-Verifikation sei der Tausch selbst, ist durch diese Proben gedeckt; der Implementer hat sie nicht selbst gebrochen (siehe V-155).

## 3. Unabhängige Prüfung des Assets

`gh release download v6.18.0 -R pt9912/ai-harness-course` (Zip und `SHA256SUMS`): sha256 `18e5443b4fca32ce452ec5dcf7cf011d7e04b0d29741b777a8dceaa452dcb8d9`, `sha256sum -c` OK; entpackt 54 Dateien, `diff -r` gegen den vendored Baum leer. Der Beleg in §7 stimmt.

## 4. Stichprobe der verkörperten Regeln gegen das neue Regelwerk

| Abschnitt des neuen Regelwerks | Verkörperte Stelle | Ergebnis |
|---|---|---|
| `modul-05` §Closure und Risiko-Ausgänge | `AGENTS.md` §3.3, §3.9, §3.13; Plan-Form | unverändert gegenüber `v6.16.0` (nur Versionszeile), kein Widerspruch |
| `modul-09` Hard Rules | `AGENTS.md` §3, `implement-slice` | unverändert (nur Versionszeile), kein Widerspruch |
| `grundlagen-harness-dateien` §Kommentar und §Gate-Index | `AGENTS.md` §3.7, `harness/README.md` §Sensors | Satz zur Disjunktheit neu; dieses Repo führt keinen Gate-Index in Teilen, Zelle *kommt nicht mit* bleibt wahr |
| `modul-10` §Harness-Einordnung | `.harness/skills/reviewer.md`, `implement-slice` Schritt 21 | Dateiname und Folgelauf: Skill nennt ihn mit Rollen-Infix (Abweichung begründet, F-576); Deckungs-Sensor-Absatz ohne Träger in diesem Repo (F-578) |
| `modul-13` §Gate-Aggregator | `harness/mk/` | Satz zur Disjunktheit setzt Teile voraus, die dieses Repo nicht führt |

Kein Widerspruch zu `AGENTS.md`, `harness/conventions.md` oder dem Reviewer-Skill, der nicht schon als Review-Finding oder akzeptierte Abweichung benannt ist.

## 5. Review-Übergaben

- **F-576** (Abweichung im Dateinamen, Adresse fehlt): bleibt für die Closure offen. Die Zusage trägt der Wortlaut. Ausgang: Register-Eintrag, am besten im Eintrag zum Dateinamen/Infix angehängt, als Zeiger auf den künftigen Slice zum Modul `reviews` (nicht angelegt, deshalb keine Zuweisung nach `AGENTS.md` §3.13). Übergabe an den Closure-Lauf (V-154).
- **F-578** (Deckungs-Sensor-Absatz ohne Adresse): derselbe Weg, Register-Eintrag; Übergabe an den Closure-Lauf (V-154).
- **F-577** (zwei Folgelauf-Formen nebeneinander): bleibt als Bestand bewusst stehen (§1 *eingefrorene Zeitdokumente*), keine Aktion. Eingeordnet als INFO.
- **F-579** (Move und Inhalt in einem Commit): vom Architect ausdrücklich entschieden (§6). Eigene Messung: alle Pfade als Rename erkannt (`git diff -M`), Rename-Erkennung hält. Keine Aktion.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-154 | INFO | Zum Prüfstand sind die Closure-Pflichten des Plans offen: DoD-Häkchen in §2 ungesetzt, §7 trägt noch keine Closure-Notiz und keinen Lerneintrag, die drei Risiken in §6 stehen auf „offen bis Closure“, das Register trägt weder F-576 noch F-578. | `modul-05` Closure (Plan §2, letzte vier Punkte); Plan §6 | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-18.md` · „Ausgang:** offen bis Closure“ | ja — der Closure-Lauf zeigt es, `make gates` prüft die Ablage | Closure-Pflicht vor Verifikation noch offen (erwarteter Stand) |
| V-155 | INFO | §7 nennt als Gegenprobe der Baseline-Verifikation „der Tausch selbst“; die Fälle geänderte, fehlende und zusätzliche Datei sowie fehlende oder abweichende Summenzeile hat der Implementer nicht gebrochen (nur der Reviewer den ersten). Der Verifier hat sie gefahren (Abschnitt 2), alle rot. | `modul-11` Bewusstes Brechen für DoD-Testbehauptungen | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-18.md` · „die Gegenprobe der Baseline-Verifikation ist der Tausch selbst“ | ja — Abschnitt 2 dieses Berichts | Gegenprobe im Beleg des Implementers nur durch den Tausch begründet |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Baum und Asset (`.harness/baseline/v6.18.0/`) | geprüft, ohne Befund (sha256, `diff -r`, 54 Dateien, keine Symlinks) |
| Verweise und Symlinks (`grep`, `find -xtype l`, `.claude/rules/`) | geprüft, ohne Befund |
| Pins in `AGENTS.md` und `harness/conventions.md` | geprüft, ohne Befund |
| Regelblock-Tabelle, `MR-000`, `MR-001` | geprüft, ohne Befund |
| Abgrenzung §1 (kein `.go`-File, `.d-check.yml`, `d-check.mk`, `tools/`, eingefrorene Dokumente) | geprüft, ohne Befund (`git diff --name-status` außerhalb des Baums: nur Commands, Symlinks, Skill, `AGENTS.md`, Konventionen, Plan, Review-Report) |
| `make gates` am Stand `48a6b00` | geprüft, ohne Befund (Exit 0) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Closure-Pflicht vor Verifikation noch offen (erwarteter Stand) · Gegenprobe im Beleg des Implementers nur durch den Tausch begründet

## Verdikt

**Merge-blockierend:** nein — alle drei Liefer-Punkte der DoD (Baum und Verweise, Abgleich, `make gates`) sind bestätigt; das Review liegt vor. Kein Befund blockiert.

**Übergabe:** Der Planner (Closure-Lauf) trägt F-576 und F-578 als Register-Einträge ein (Ausgang „Dateiname“ und „Deckungs-Sensor ohne Träger“, Zeiger auf den künftigen Slice zum Modul `reviews`, ohne Zähler), schreibt Closure-Notiz mit Lerneintrag, gibt den drei Risiken aus §6 je einen Ausgang (voraussichtlich *entfallen* mit Begründung: kein Befund trat ein) und setzt die Häkchen vor dem reinen `git mv` nach `done/`. V-155 geht als Klasse in den Zähler. Dieser Report ist ein **Lauf-Beleg**.
