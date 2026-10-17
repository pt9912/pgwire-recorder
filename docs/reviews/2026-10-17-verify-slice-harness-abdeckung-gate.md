# Verifikation: slice-harness-abdeckung-gate — 2026-10-17

**Rolle:** Verifier (Modul 8). **Gegenstand:** `slice-harness-abdeckung-gate`, Diff `cbeb3381..HEAD` (HEAD `d7fede16`; Commits `01523fa0`, `cacaf26a` Architect, `d7fede16`).
**Maßstab:** DoD des Plans, `SPEC-048` Punkte 6 bis 8, [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md), `AGENTS.md` §3.9 bis §3.14.
Der Review-Report des Slice wurde nicht gelesen.

## Verdikt

**Die Lieferung stimmt mit Plan und Vertrag überein; Verdikt: bestanden mit zwei kleinen Befunden (F-1, F-2).** Die Hauptwirkung ist am Repo und an einem Gegenbeweis in einer Kopie gemessen. Die Closure-Voraussetzungen sind noch nicht erfüllt (Review, Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen offen; erwartet, siehe unten).

## Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make abdeckung-check` | grün (Exit 0) |
| `make abdeckung-gegenprobe` | grün (Exit 0), rund 6 s |
| `make docs-check` | grün (992 Dateien, 0 Befunde) |
| `make kopf-check` | grün (Exit 0) |
| `make doc-trace` | siehe Hauptwirkung |
| `make gates` | nicht gefahren (vereinbart); Beleg in §7 des Plans, nicht von mir bestätigt |

## Hauptwirkung am Repo

- `docs/user/abdeckung-gesamt.md`: `| LH-QA-07 | n/a | n/a | n/a | 1/3 Gate; 2/3 —; 3/3 Gate | teilweise |`. In `abdeckung-vollstaendig.md` kommt `LH-QA-07` nicht vor. `abdeckung-gate.md` führt genau die Zeilen `Messung-1-von-3` und `Messung-3-von-3` mit Nachweis `tools/harness/lint-gegenprobe.sh`.
- `make doc-trace`: `LH-QA-07` Coverage `—`; `LH-QA-02`, `LH-QA-06`, `LH-RB-01` `Abdeckung`.
- `git diff cbeb3381..HEAD -- docs/user/abdeckung-e2e.md docs/user/abdeckung-unit.md` ist leer (0 Byte): beide Tabellen byte-gleich.
- `.d-check.yml`: nur Kommentar und Label `Tests` → `Abdeckung`, Quelle unverändert.

## Gegenbeweis in einer Kopie (`git archive HEAD` im Scratchpad, nie im Arbeitsbaum)

- Mit Teil 2 zusätzlich im Kopf der Lint-Gegenprobe (alle drei Teile): `abdeckung-gesamt.md` zeigt `1/3 Gate; 2/3 Gate; 3/3 Gate | vollständig`, `abdeckung-vollstaendig.md` führt `LH-QA-07 | Gate`, `make doc-trace` zeigt `LH-QA-07` mit `Abdeckung`, `--check` grün.
- Mit Teil 1 entfernt (Teile 2 und 3): `teilweise` (`1/3 —; 2/3 Gate; 3/3 Gate`), nicht in `abdeckung-vollstaendig.md`, `doc-trace` zeigt `—`.
- Verdrahtung: veraltete Tabelle und fehlende Tabelle `abdeckung-gate.md` mit `--check` ergeben Exit 1 und nennen die Datei; fehlendes `docs/user/` ergibt Exit 2 mit Meldung.

## Mutationen von `tools/test/abdeckung.sh` (Kopie je Mutant, genau ein Treffer, unveränderte Gegenprobe dagegen)

| Zusage | Mutation | Ergebnis |
|---|---|---|
| geteilte Messung erst mit allen Teilen vollständig | `|| return 1` → `&& return 0` | rot (12 Fälle, u. a. `teile-gesamt`) |
| i > n abgelehnt | `a > b` → `a > b + 1` | rot (9, u. a. `gate-teil-bereich`) |
| Gate-Deklaration: nur `LH-QA`/`LH-RB` (Zweig allein) | Zweig `LH-FA` abgeschaltet | grün; vom Plan als äquivalent geführt (das Muster lehnt `LH-FA` mit anderer Meldung ab), Wortlaut nicht zugesagt |
| Zeile = Anfang der Deklaration im Kopf | `dz = NR - 1` | rot (7, `gate-zeile`) |
| fünf Tabellen geschrieben und geprüft | Gate-Tabelle aus der Schleife | rot (17, `schreiben-fehlend`) |
| reguläre Dateien, Symlink nicht gelesen | `-type f` entfernt | rot (1, `gate-nicht-gelesen-symlink`) |
| verschiedene n sind Fehlform | Prüfung `fn[id] != n` abgeschaltet | rot (20, `widerspruch-n`) |
| Ausgang 2 ohne `docs/user/` | Prüfung abgeschaltet | rot (2, `ohne-docs-user`) |
| Ausgang 1 bei Fehlform | `exit 1` → `exit 0` | rot (96) |
| eingerückt mit Tab | nur Leerzeichen | rot (3, `gate-eingerueckt-tab`) |

Neun von zehn Mutanten rot, der zehnte ist die im Plan geführte Äquivalenz. Kein ungültiger Mutant.

## DoD gegen Diff

1. **Werkzeug und Vertrag:** erfüllt. Skript liest den Kopf (`gate_awk`), Pfad `Messung-<i>-von-<n>` (`pform`), Tabelle Gate, Vollständigkeit über alle Teile 1 bis n. Randformen stehen in `SPEC-048` Punkte 6 bis 8 (Architect `cacaf26a`); `harness/README.md` §Sensors nennt Nachweisart und geteilte Messung mit [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md) und `· seit slice-harness-abdeckung-gate`; jede Aussage hat einen Fall in der Gegenprobe (§3.11).
2. **Gegenprobe:** erfüllt. 56 Aufrufe gegen 13 vor dem Slice (740 Zeilen gegen 262); die vorhandenen Fälle bleiben, nur der Fall `lastenheft-ohne-anforderung` wurde auf Exit 2 gedreht (Vertrag Punkt 6, vom Plan so verlangt). Fall `gegenprobe-ohne-deklaration` (Datei nach dem Vorbild von `integration-volumes-gegenprobe.sh`) ist vorhanden (Sendung eingelöst). Die Zeile in `harness/README.md` nennt die neuen Fälle.
3. **Erste Gate-Deklarationen und RTM:** erfüllt, wie oben gemessen. Die Namen im Kopf der Lint-Gegenprobe tragen: jedes Präfix hat Treffer in der Datei (`p0-grundlauf` 4, `p3-` 22, `p4-` 23, `p5-` 48, `p6-` 18, `p8-` 73; `p7-whitebox-{kern,driven,pgwire,einstieg}` je 2, `p7-internal-test` 4, `p7-endet-auf-export-test` 4, `p7-test-in-bruecke` 4). Teil 2 ist bewusst nicht deklariert.
4. **`make gates` grün:** nicht von mir gefahren; Beleg in §7 (Lauf am Endstand, Hinweis auf einen ersten roten Hintergrundlauf mit Erklärung) liegt vor.
5. bis 9. (Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen): offen, `docs/plan/planning/observations/` im Diff unberührt.

## Plan gegen Code

- §3 (Plan) stimmt mit den 13 Pfaden des Diffs überein; `harness/mk/abdeckung.mk` ändert nur Kommentar und Hilfetexte, nicht das Rezept (Plan: „keine Änderung am Rezept“). `harness/sensors/abdeckung-check.md` und die Tabellen sind nachgezogen.
- §1 Ausschlüsse eingehalten: keine anderen Gegenproben angefasst (nur 3 Kopfzeilen der Lint-Gegenprobe), kein Produkt-Code, kein Lastenheft, kein Eingriff in `GATE_CHECKS`-Prüfung (Grenze steht im Skriptkopf), Teil 2 nicht deklariert, ungeteilter Pfad `Messung` gültig geblieben.
- §6 Randformen: die Entscheidungen von `cacaf26a` (Lastenheft nur `### LH-XX-01`, erste festlegende Deklaration, eingerückt an jeder Stelle, Zeile des Anfangs, Symlink) sind im Skript umgesetzt und in `d7fede16` mit Fällen und Mutationen versehen (`lastenheft-fremde-kennung`, `widerspruch-erste-*`, `gate-eingerueckt-nach-kopf`, `gate-zeile-mehrzeilig`, `gate-nicht-gelesen-symlink`; Symlink-Mutant von mir rot gesehen).
- Akzeptierte Negative und äquivalente Mutanten in §6 und im Kopf der Gegenprobe stimmen überein (`LC_ALL=C`, Reihenfolge von i > n und n = 1, Kopfende, `LH-FA`-Zweig); der Zweig `LH-FA` habe ich selbst grün gesehen.
- **Marke nach §3.14:** `git diff --numstat cbeb3381..HEAD` über die Pfade des Plans (`tools/test tools/harness/lint-gegenprobe.sh harness/README.md harness/sensors/abdeckung-check.md harness/vertraege/SPEC-048-abdeckung.md docs/user .d-check.yml`) ergibt +784/−121, Summe **905** (davon `SPEC-048` 32). Basis `cbeb3381` ist der Elternteil des Anspruchs-Commits `1ab4f9d6` (geprüft). Über der Meldeschwelle 800, unter der Marke 1000; die Schätzung 450 bis 700 ist überschritten (Gegenprobe 524 von 905, rund 58 Prozent, wie vom Architect-Faktor „Doppeltes“ erwartet, aber die Gegenprobe wuchs durch die Nacharbeit nach `cacaf26a`). Die Rückführung (>1000) tritt nicht ein.
- **Schichten und Liefer-Punkte:** Harness-Werkzeug und Nutzer-Doku, zwei; drei Liefer-Punkte. Ich stimme zu.

## Befunde

- **F-1 (niedrig):** §7 *Größe* nennt 823 (Stand vor `cacaf26a` und `d7fede16`); die Messung am Endstand ist 905. §7 ist als Beleg des Implementers damit veraltet (§3.9: der Plan folgt jeder Korrektur); die Zahl und die Zerlegung (Skript 271, Gegenprobe 474 …) sind bei Closure nachzuziehen. Dass 905 über der Meldeschwelle 800 liegt, hat der Implementer nach Auskunft des Auftrags gemeldet; der Plan nennt es nicht für 905.
- **F-2 (niedrig, Hinweis):** Der Hinweis in §6 (Architect) nennt für Teil 1 auch `p9-*`; der Kopf deklariert `p6-*` statt `p9-*`. Das ist enger als vom Architect bedacht, damit keine Überzusage, aber §7 „Abgleich der Namen“ führt `p6-*` und lässt `p9-*` weg, ohne es zu begründen. Bei Closure als bewusste Auswahl vermerken oder `p9-*` prüfen.

Keine Befunde gegen Vertrag, Gegenprobe oder Hauptwirkung.

## Closure-Voraussetzungen und vorgeschlagene Risiko-Ausgänge (Urteil bleibt beim Implementer)

Offen sind Review-Report, Closure-Notiz mit Lerneintrag (die sieben Platzhalter in §7), Register (nichts im Diff; „keine Beobachtung“ wäre ebenfalls eine Antwort), die Ausgänge der fünf Risiken in §6 und die drei Paarungen. Vorschlag nach meiner Messung:

- *Falsche Gesamtzahl n*: **weiter offen** (Register), das Skript fängt nur verschiedene n; die Zusage „einheitlich falsche Zahl“ bleibt Review ([ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md)).
- *Deklaration sagt mehr zu*: **entfallen**, Namen gegen die Datei abgeglichen (siehe DoD 3), mit F-2 als Randnotiz.
- *Bestehende Prüfung fällt weg*: **entfallen**, Fälle bleiben (13 → 56 Aufrufe, Mutant „Exit 0“ rot).
- *RTM zählt Teilabdeckung als belegt*: **entfallen**, am Repo (`—`) und im Gegenbeweis (`teilweise` bei fehlendem Teil) gemessen.
- *Andere Gegenproben werden gelesen*: **entfallen**, Fall `gegenprobe-ohne-deklaration`, und `abdeckung-check` ist am Repo grün.

## Grenzen dieser Verifikation

Nicht gefahren: `make gates` (der Hauptbaum fährt ihn einmal am Ende), `make lint-gegenprobe`. Ob Teil 1 und Teil 3 inhaltlich die Messmethoden 1 und 3 von `LH-QA-07` decken, ist Urteil ([ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md)); ich habe nur geprüft, dass die genannten Fälle in der Lint-Gegenprobe existieren.
