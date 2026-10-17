# Verifikation: slice-harness-integration-volumes — 2026-10-17

**Rolle:** Verifier (Modul 8). **Gegenstand:** `slice-harness-integration-volumes`, Diff `de908591..HEAD` (HEAD `c9400ec4`).
**Maßstab:** DoD des Plans, `SPEC-057`, `SPEC-055` Punkt 9, `SPEC-052` Punkt 13, `AGENTS.md` §3.9 bis §3.14.
Der Review-Report des Slice lag beim Urteil nicht vor und wurde nicht gelesen.

## Verdikt

**Die Lieferung stimmt mit Plan und Vertrag überein; Verdikt: bestanden mit zwei kleinen Befunden (F-1, F-2).** Die Closure-Voraussetzungen sind noch nicht erfüllt, weil Review, Closure-Notiz, Lerneintrag und Risiko-Ausgänge fehlen (siehe Abschnitt Closure).

## Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make docs-check` | grün (980 Dateien, 0 Befunde) |
| `make kopf-check` | grün (Exit 0) |
| `make integration-volumes-gegenprobe` | grün, 59 s (Plan nennt 42 s; das Docker war zeitgleich durch den Reviewer belastet) |
| `make test-integration` (Handlauf) | grün, 2 min 9 s; ein zweiter Lauf rot an `TestE2EPlayVergleichAbweichung` (Zeitüberschreitung bei `DROP DATABASE`, kein Volume-Bezug, unter Last durch den parallelen Reviewer); ein dritter Lauf grün |
| `make test-container` (Handlauf) | grün |

## Hauptzusage am echten Docker (Differenz nach Namen)

Der Host war nicht ruhig: der Reviewer legte zeitgleich Volumes an (`pgr-ivg-…` und kurzlebige anonyme Hex-Volumes). Eine reine Vorher-nachher-Differenz des Hosts ist deshalb nicht aussagekräftig (sie zeigte 3 Hex-Namen und acht `pgr-ivg-`-Namen, die nach Minuten teils schon wieder fehlten). Ich habe darum die Volumes **meines Laufs** über `PGR_LAUF` und das Mitlesen von `docker container inspect` auf `<Lauf>-pg` und `<Lauf>-tests` erfasst und nach dem Lauf nach Namen geprüft:

- `make test-integration`, grün: gesehen `68d05faf…` (anonym, vom PostgreSQL-Container) und `pgr-ver-…-daten`; nach dem Lauf **kein** Volume und kein Container mehr vorhanden. Beim roten Lauf ebenso (`53a617b9…`, nach dem Lauf weg).
- **Alter Runner** (`git show de908591:tools/test/run-integration-tests.sh`, Kopie im Scratchpad, ohne Eingriff in den Arbeitsbaum): hinterließ `d16e26ca…`. Das habe ich per Namen entfernt (`docker volume rm d16e26ca…`); sonst habe ich nichts entfernt, kein `prune`.
- `make test-container`: das Volume `f9b980f1…` des Containers `pgr-ct-…-pg` war nach dem Lauf weg, kein Container übrig.

Die Hauptzusage ist damit am echten Docker belegt, einschließlich des Gegenstücks vor der Behebung.

## Mutationen (Kopien im Scratchpad, nie im Arbeitsbaum)

| Zusage | Mutation | Ergebnis |
|---|---|---|
| `SPEC-057` Punkt 2 (`rm -fv` im Runner) | `rm -f` in `integration-aufraeumen.sh` | Gegenprobe rot, unveränderte Kopie `p2-grundfall`, `p2-ohne-v` |
| `SPEC-057` Punkt 4 (Rest macht den Lauf rot) | Prüfung und Aufräumen enden immer mit 0 | Gegenprobe rot (`p3-angelegt`, `p4-rest-rot` in der unveränderten Kopie) |
| `SPEC-055` Punkt 9 (`rm -fv` für `$pg` im Aufräumen) | `rm -f` in `run-container.sh` Zeile 37 | `container-gegenprobe` rot, `p9-sigINT-psql-volume` und `p9-sigTERM-psql-volume` |

Zwei der Mutationen nutzen einen groben Textersatz; die Gegenprobe lehnt die Kopie ab, bevor die feinere Zuordnung zählt. Die feine Zuordnung (genaue Menge roter Fälle je Mutant) liegt im Lauf der Gegenprobe selbst und war in der unveränderten Fassung grün.

## DoD gegen Diff

1. **Aufräumen im Runner:** erfüllt im Code (`rm -fv "$tests" "$pg"`); Belege in §7 (Fundstellen, `.Config.Volumes` auch für `a-check` und `d-migrate`, Handlauf) vorhanden und in meiner Messung bestätigt. Häkchen noch nicht gesetzt.
2. **`test-container`:** `rm -fv` an allen drei Stellen von `run-container.sh` (`:37`, `:145`, `:198`), Fälle `p9-*-volume` je Ende und Mutationen vorhanden; `p7-pg-vor-replay` nachgezogen. Beleg in §7 vorhanden; bestätigt.
3. **Prüfung im Runner:** `integration-aufraeumen.sh` liest Namen vor dem Entfernen (ohne `$daten`), prüft danach, meldet Namen, entfernt per Namen, Exit 1 mit Vorrang des Laufcodes; Runner und Gegenprobe laden dieselbe Datei. Erfüllt.
4. **Gegenprobe:** Ziel hängt an `GATE_CHECKS` (`integration.mk`), steht in `GATE_SCHNELL_AUSNAHME`, Zeile in `harness/README.md` §Sensors, `LANGSAM` der Gegenprobe der schnellen Kette nachgezogen. Fälle, Mutanten und Tabelle *Zusage · Mutation · roter Fall* in §7 vorhanden. `SPEC-052` Punkt 13: die Gegenprobe liest nur Kopien unter `tools/test/`, nichts aus `docs/user/` oder der Wurzel-`README.md` (geprüft im Skript). Erfüllt.
5. **`make gates` grün:** nicht von mir gefahren (Wanduhr 1018 s, Reviewer teilt das Docker); der Beleg in §7 nennt einen grünen dritten Lauf am Endstand `dc6da5fa`. Seither änderte sich nur Dokumentation (`git diff dc6da5fa..HEAD` berührt nur Plan und `SPEC-057`; `SPEC-057` liegt auf der Positivliste nicht, daher bleibt es dem Closure-Lauf überlassen, ihn frisch zu fahren; siehe Closure).
6. bis 10. (Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen): offen.

## Plan gegen Code

- §3 (Plan) stimmt mit den 13 Dateien des Diffs überein, einschließlich der Dateien, die der Plan nach dem Code nachgezogen hat (`gates-schnell-gegenprobe.sh`, `postgres-versionen-gegenprobe.sh`).
- §6 und die akzeptierten Negative stimmen mit `SPEC-057` Punkt 5 und 7 überein: `p5-term-entfaellt` und `p5-int-entfaellt` ohne roten Fall; die Falle für `INT` ist im Code (`trap 'trap "" INT; exit 130' INT`), aber nicht zugesagt.
- §1 Ausschlüsse eingehalten: kein `prune`, kein Eingriff in fremde Volumes (Fall `p4-fremd`), keine Änderung der übrigen `docker run`, kein Produkt-Code.
- Schichten: zwei (Harness-Werkzeug; Gate-Index und Verträge), Liefer-Punkte drei. Ich stimme zu.
- **Marke nach §3.14:** `git diff --numstat de908591..HEAD` über die Pfade des Plans ergibt +463/−42, Summe **505**. Das liegt unter Meldeschwelle 750 und Marke 900. Der Implementer meldete 457 bei `5e15247b`; die Differenz sind die späteren Commits (Nachzug von `SPEC-057`/`SPEC-055`, Gegenprobe-Nacharbeit).

## Befunde

- **F-1 (niedrig):** §4 *Größe* nennt „Basis: der Commit des Anspruchs, beim Anspruch einzutragen“, nennt aber keinen Hash; die Basis steht nur indirekt in §7 (`de908591`). Die Metrik ist damit nach `AGENTS.md` §3.14 nicht ohne Rückfrage nachrechenbar. Mit dem Hash `de908591` ist sie nachgerechnet (505).
- **F-2 (niedrig):** §7 *Offene Punkte für den Architect/Planner* führt Punkte (zweites `SIGTERM`, Falle für `INT`), die §6 *Prüfung des Architect* bereits entschieden hat; vor der Closure bereinigen, sonst liest die Closure-Notiz offene Fragen, die nicht mehr offen sind (`AGENTS.md` §3.9).
- Hinweis, kein Befund: Der Test `TestE2EPlayVergleichAbweichung` lief unter der Last des parallelen Reviewers einmal in eine Zeitüberschreitung. Mit dem Slice hat das nichts zu tun, spricht aber für die bekannte Flake-Klasse unter Last.

## Closure: stehen die Voraussetzungen?

Noch nicht, und das ist erwartbar vor dem Review:

- DoD-Häkchen: keine gesetzt. Review-Report: fehlt. Closure-Notiz und Lerneintrag: fehlen (§7 trägt Belege, aber keine Closure-Notiz).
- **Register** `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`: zwei Belege unter `evidence/` (`slice-walking-skeleton-record`, `slice-v1-abschluss-container-doku`), also 2×. Dieser Slice hat keine eigene Datei dazugelegt; DoD verlangt die Fortschreibung oder die Aussage „keine Beobachtung“. Die „Beobachtung für das Register“ in §7 (Plan nannte Mutationen so, als bräche jede nur den eigenen Fall) ist als Text da, aber nicht als Datei eingetragen.
- **Risiken §6:** alle drei stehen mit „offen bis Closure“, das ist keiner der drei Ausgänge. Vorschlag: Laufzeit (`gate-laufzeit-waechst-mit-gegenprobe`) mit Zahl aus §7 → *weiter offen* ins Register; volle Platte → *entfallen* mit Begründung (Leck behoben, gemessen); neues `docker run` ohne `--rm` → *weiter offen* (`SPEC-057` Punkt 11, Grenze). Urteil bleibt beim Planner.
- Vor dem `git mv` nach `done/` muss `make gates` am Endstand frisch grün sein.

## Beleg-Prüfung

Alle Belege, die die DoD verlangt (Namen vor und nach, Fundstellen, `.Config.Volumes`, Zusage·Mutation·roter Fall, Marke, Laufzeit), stehen in §7; keiner fehlt.
