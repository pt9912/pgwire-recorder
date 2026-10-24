# Review-Report: slice-release-v0-2-0 und slice-go-toolchain-1-27-3 — 2026-10-24

**Review-Art:** Code — Code-Review gegen Plan, Entscheidungen und Hard Rules (`AGENTS.md`), ein Report für zwei Slices.

**Gegenstand:** `main` bei `511428c9`. (1) `slice-go-toolchain-1-27-3`: Basis `7cfad5ab`, Pin-Commit `143e9994`, Belege `32d0d84a`. (2) `slice-release-v0-2-0`: Basis `193f0e0b`, Release-Commit `d720e44f`, Tag `v0.2.0` auf `32d0d84a`, Eintrag `docs/maintainer/freigaben/0.2.0.md`, Bezug-Commit `da5fe8b6`, Commit zu `releasing.md` `2b0869ff`, Carveout `CO-002`.

**Skill:** `.harness/skills/reviewer.md` @ `4155023d`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24

> **Zitier-Form.** Kennung, nicht Adresse: Slices als `slice-<Kennung>`, Sensoren als `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext:**

- Slice-Pläne `slice-go-toolchain-1-27-3` und `slice-release-v0-2-0` (§1 bis §4, §6, §7 Belege des Implementers, §8)
- `CO-002`, `CO-001` (gelesen), Eintrag `freigaben/0.2.0.md`, `docs/maintainer/releasing.md` §10
- `SPEC-067` Punkte 16 bis 20, `SPEC-052` Punkt 12, `SPEC-060` Punkt 14; `LH-FA-19.a` (Spezifikation)
- `LH-FA-19`, `LH-FA-16`, `LH-FA-01`, `LH-FA-24`, `LH-QA-03`
- aktive ADRs der Pläne: [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md), [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md), [ADR-0052](../plan/adr/0052-scan-des-images-vor-dem-push.md), [ADR-0053](../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md), [ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md), [ADR-0055](../plan/adr/0055-gegenproben-haben-ein-zeitbudget.md), [ADR-0056](../plan/adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)
- `AGENTS.md` §3.9, §3.11, §3.13, §3.14
- vorherige Reports am Gegenstand: `slice-play-feldvergleich-doku` (Klasse: Zusage weiter als Prüfung), `slice-harness-doc-trace-archiv` F-8 (Plan folgt dem Befund nicht)

Gelesen und gegen die Originale nachgefahren (nur lesend, kein `make gates`): `git diff --numstat` der Größenmetrik und der sechs Pfade aus `CO-002`; `gh run view 40074219261` (Jobs, Zeilen, `skopeo`-Wortlaut, kein `release: hinweis:`); `gh release view v0.2.0` (sieben Anhänge, kein Entwurf); `git ls-remote` für Tag und Tap; `docker buildx imagetools inspect` des Image (Index-Digest `sha256:1d501812…`, `linux/amd64` und `linux/arm64`); das veröffentlichte Binary `pgwire-recorder-linux-amd64` (Summe gegen `SHA256SUMS`, `version`, `--version`, `-v`, `--help`, `play --help`, `sha256sum -c` ohne `--ignore-missing`). Nicht nachgefahren: `make homebrew-nachweis`, `docker pull` der vier Referenzen, die Läufe von `make gates` und der Gegenproben.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Handbuch, `README.md` und `releasing.md` §1 sagen für `0.2.0`, die Installation aus dem Tap sei unter Linux auf `amd64` gefahren. Eine Beleg-Datei nach `SPEC-067` Punkt 20 (6) (`docs/reviews/<Datum>-nachweis-<Kennung>.md` mit Kopf, `uname`, Tag und Digest des Image, den Befehlen (a) bis (e) mit vollständiger Ausgabe und Exit-Code, §7 nennt Pfad und Summe) gibt es nicht; Punkt 11 des Eintrags fasst die Läufe in Prosa zusammen, und Punkt 20 (6) lässt den Eintrag der Freigabe von Tap-Belegen unberührt. Szenario: `LH-FA-19.a` sagt „Eine Angabe ohne die Ausgaben ist kein Nachweis“; der Verifier und jeder spätere Lauf können die Zeilen (`level=INFO`, die erwartete `PGR-W2001`, Exit-Code 0, Image-Digest) nicht prüfen, und die Aussage zu Linux/`amd64` steht im Handbuch als belegt. Die unveränderten Sätze „baut nichts aus dem Quelltext und installiert kein weiteres Paket“ und „die Formel trägt die höchste stabile Version“ nennt der Eintrag für `0.2.0` nicht. | `SPEC-067` Punkt 20 (6); `LH-FA-19.a`; `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „Gefahren ist diese Installation unter Linux auf `amd64`“; `docs/maintainer/freigaben/0.2.0.md` · „Homebrew: `make homebrew-nachweis` (Beginn 2026-10-24 07:42:09Z“ | nein — `ls docs/reviews` zeigt keine Datei `…-nachweis-…` zu `0.2.0`; kein Gate liest sie | Nachweis ohne Beleg-Datei, Aussage als belegt geführt |
| F-2 | MEDIUM | Der Plan folgt der Korrektur „`make homebrew-nachweis` nach dem Tag gefahren“ nur in §1 (Abgrenzung), DoD 2 und §7: Die Rollen-Tabelle in §1 (Implementer), die Tabelle und die Zeile zu Punkt 11 in §3, die Reihenfolge Schritt 5 und Randform (h) in §6 nennen ihn nicht; §7 sagt „§1 schloss einen neuen Homebrew-Nachweis aus“, was §1 nicht mehr sagt. Das Muster war in früheren Reports schon zweimal LOW (Wiederholungsregel des Skills). Szenario: Wer §3 als Plan liest, kennt den Schritt und seine Grenze (Linux/`amd64`, Beleg-Form) nicht. | `AGENTS.md` §3.9; Wiederholung nach Reviewer-Skill (MEDIUM) | `docs/plan/planning/in-progress/slice-release-v0-2-0.md` · „DoD 2; jeder Befehl läuft nach dem Tag gegen die **veröffentlichten** Artefakte“ | nein — Urteil beim Lesen; `make kopf-check` prüft nur den Kopf | `plan-folgt-korrektur-nicht` |
| F-3 | LOW | §7 des Release-Plans sagt, `slice-erster-release-homebrew-nachweis` „bleibt für die übrigen Plattformen zuständig“, und §1 weist die Belege für macOS und `arm64` ihm zu. Der Slice liegt in `done/`, sein §7 führt macOS und Linux auf `arm64` als Grenze ohne Lieferung; kein offener Slice nimmt die Belege an. | `AGENTS.md` §3.13 | `docs/plan/planning/in-progress/slice-release-v0-2-0.md` · „bleibt für die übrigen Plattformen zuständig“ | nein — Verzeichnisposition `done/slice-erster-release-homebrew-nachweis.md` | `folge-slice-adresse-nimmt-nicht-an` |
| F-4 | LOW | §1 des Pin-Slice nennt unter den Ausschlüssen „Ein Sensor, der beide Pin-Stellen abgleicht, und der Nachtlauf“ als Sache von `slice-harness-ci-freshness-images`; dessen §1 und Randform 4 schließen den Abgleich der zwei Stellen ausdrücklich aus. Der Satz „der Abgleich ist hier eine Handlung“ nimmt den Rest zurück, die Adresse nimmt den Abgleich aber nicht an. | `AGENTS.md` §3.13 | `docs/plan/planning/in-progress/slice-go-toolchain-1-27-3.md` · „Ein Sensor, der beide Pin-Stellen abgleicht“ | nein — `grep -n 'Abgleich' open/slice-harness-ci-freshness-images.md` | `folge-slice-adresse-nimmt-nicht-an` |
| F-5 | LOW | `releasing.md` §1 sagt, das Werkzeug habe die Installation „für Linux/amd64 mit den Formeln `0.1.1` und `0.2.0`“ belegt; §7 nennt weiter „belegt für Linux/amd64 mit der Formel `0.1.1`“. Die zwei Stellen desselben Dokuments widersprechen sich (und stützen F-1 nicht gleich). | Maintainability | `docs/maintainer/releasing.md` · „belegt für **Linux/amd64** mit der Formel `0.1.1`: Homebrew nimmt die Formel an“ | nein — `grep -n 'Formel .0\.1\.1' docs/maintainer/releasing.md` | Stand-Aussage an zwei Stellen verschieden |
| F-6 | LOW | Der Eintrag schreibt „Die Probe belegt den Satz am Binary“ zu `play --diff-all`; am Binary gelesen ist nur der Hilfetext (`play --help`). Das Verhalten stützt der Eintrag auf Tests früherer Slices. Das Gegenmittel in §6 des Plans nannte einen Lauf des gebauten Binaries `play --diff-all` mit und ohne `--compare-responses`; die Abweichung davon nennt der Eintrag nicht. Die Grenze „keine Mutation gefahren“ steht dort. | `AGENTS.md` §3.11 | `docs/maintainer/freigaben/0.2.0.md` · „Die Probe belegt den Satz am Binary, nicht die Wirksamkeit der Tests.“ | nein — Lauf `play` gegen eine reale Instanz fehlt | `zusage-im-kommentar-weiter-als-pruefung` |
| F-7 | LOW | `CO-002` nennt für die Ersatzbelege 2 und 3 den Hash `R` und definiert `R` als „getaggten Commit“; der Tag liegt auf `32d0d84a`, und die Ersatzbelege laufen dort. Die Korrektur steht in §3 des Plans, nicht im Text des Carveouts; `releasing.md` §10 nennt als getaggten Commit den, der `VERSION` und die Zeile trägt (`d720e44f`). | Maintainability | `docs/plan/carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md` · „Hash `R` (der getaggte Commit“ | nein — `git ls-remote origin 'refs/tags/v0.2.0^{}'` gibt `32d0d84a` | Carveout-Text nennt anderen Hash als der Tag |
| F-8 | INFO | Punkt 1 (b) des Eintrags nennt als Fundstelle ein Log `gates-main4.log`; es liegt in keinem Repo-Pfad (nur im Scratchpad einer Sitzung). Der Eintrag nennt Dauer und Last von (b) als nicht erfasst; der Hash des Laufs ist über (a) und den Diff `143e9994..32d0d84a` (nur Plandatei) gestützt. Rolle: Verifier. | Maintainability | `docs/maintainer/freigaben/0.2.0.md` · „Log `gates-main4.log`“ | nein | Fundstelle eines Belegs nicht im Repo |
| F-9 | INFO | Pin-Slice: DoD 2 verlangt Lauf-Kennung, Hash und Zeile des Scans in §7; sie stehen in `freigaben/0.2.0.md` (Punkte 5 und 6, Lauf `40072517299` auf `32d0d84a`), nicht in §7 des Pin-Slice. DoD „`make gates` am Endstand“ ist an `143e9994` belegt; der Endstand `32d0d84a` unterscheidet sich nur durch die Plandatei (inert nach `SPEC-052` Punkt 12). Beides ist bis zur Closure offen, nicht falsch. | Maintainability | `docs/plan/planning/in-progress/slice-go-toolchain-1-27-3.md` · „DoD 2 (Scan) belegt die Koordination über den Trockenlauf; nicht abgehakt.“ | nein | — |
| F-10 | INFO | Dieser Report deckt zwei Slices; sein Dateiname trägt nur `slice-release-v0-2-0`. Eine Deckungsprüfung, die den Report über die volle Kennung zuordnet, findet für `slice-go-toolchain-1-27-3` keinen. Rolle: Orchestrator. | `.harness/skills/reviewer.md` §Output-Schema | `docs/reviews/2026-10-24-review-slice-release-v0-2-0.md` · Dateiname | nein | Ein Report für zwei Slices |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `Dockerfile` Zeile 9 und `harness/mk/build.mk` Zeile 15 | geprüft, ohne Befund: gleicher Text `golang:1.27.3-alpine@sha256:b8c9b82a…505e9c` an beiden Stellen; `grep` nach `golang:1`, `f92b6ef800e4`, `1.27-alpine` außerhalb `docs/` und `baseline/` ohne weiteren Treffer; Größenmetrik des Pin-Slice 4 Zeilen (Marke 1500) |
| Lifecycle-Moves beider Slices (`fed48ec2`, `7cfad5ab`, `1b51308c`, `193f0e0b`) | geprüft, ohne Befund: reine Renames (100 %), Inhalt (`Verantwortlich`, Basis-Hash) in eigenen Commits |
| `VERSION`, Zeile `0.2.0` der Änderungshistorie (`d720e44f`) | geprüft, ohne Befund: drei Zellen, dritte nicht leer, kein weiteres Pipe-Zeichen; jeder Satz am veröffentlichten Binary gelaufen (`--version`, `-v`, `version` je `pgwire-recorder 0.2.0`, Exit 0; `--help` nennt beide); Verhalten von `--diff-all` siehe F-6 |
| Eintrag `freigaben/0.2.0.md`, Punkte 4 bis 10 | geprüft, ohne Befund: Lauf `40074219261` `success` auf `32d0d84a`, beide Jobs, Zeilen wörtlich im Log, `skopeo`-Wortlaut, kein `release: hinweis:` im Lauf des Tags, Tag-Objekt und aufgelöster Hash, sieben Anhänge, Tap-Commit `d96e2968`; der Wiederanstoß ist als nicht belegt genannt |
| Eintrag, Punkt 3 und `CO-002` Ersatz 3 | geprüft, ohne Befund: `git diff -M --numstat 35778a2f 32d0d84a` über die sechs Pfade gleich dem Eintrag (1/1, 4/4, 1/1), die sechs Zeilen sind Pfade, ein Kommentar und die Berechnung der Wurzel (`../..` zu `../../..`, im Eintrag als einzige Codezeile genannt und vom Nutzer gedeckt, Wiedergabe der Koordination); Grenze „Lauf am Host steht aus“ steht im Eintrag |
| Punkte 1 und 2 des Eintrags | geprüft, ohne Befund außer F-8: Ersatzbelege nach `CO-002` mit Befehl, Hash, Datum, Dauer, Last, Ausgang; „keine Mutation gefahren“ genannt, die Schlusszeile von `test-container-gegenprobe` nicht als Beleg zitiert |
| Kein Wert eines Geheimnisses | geprüft, ohne Befund: Eintrag nennt `HOMEBREW_TAP_GITHUB_TOKEN` und den Zeitstempel von `gh secret list`, keinen Wert; `skopeo`-Zeilen tragen `***` |
| Bezug-Commit `da5fe8b6`, `docs/user/benutzerhandbuch.md` und `README.md` | geprüft, ohne Befund außer F-1: nur diese zwei Dateien, kein `0.1.1` außerhalb der Änderungshistorie; Digest, Tags, beide Plattformen im Index, sechs Namen, `sha256sum -c --ignore-missing` (OK) und ohne die Option (Exit 1, Warnung über fünf Dateien) nachgefahren; Grenzen für macOS, Windows und `linux/arm64` genannt |
| Versionszählung | geprüft, ohne Befund: Handbuch `0.13` (eine Hebung für `d720e44f` und `da5fe8b6`), Stand 24.10.2026; `releasing.md` `0.16` mit Stand 2026-10-24, die nachgetragene Zeile `0.15` stimmt mit `cd88df3e` (2026-10-21, Kopf `0.14` zu `0.15`) überein |
| `docs/plan/carveouts/README.md`, `CO-002` | geprüft, ohne Befund außer F-7: Index-Zeile und Trigger stimmen; Geltungsbereich allein `v0.2.0`, Folge-Tag nein in Plan §1, §4, §6 und im Carveout gleich |
| Größenmetrik des Release-Plans (`AGENTS.md` §3.14) | geprüft, ohne Befund: Basis `193f0e0b`, 75 Zeilen am Stand `da5fe8b6`, 93 am Stand `511428c9` (Marke 1500) |
| Plan §8 beider Slices (Zählung) | geprüft, ohne Befund: Release 3 Liefer-Punkte und 2 Schichten, Pin-Slice 2 Punkte und 1 Schicht; die Nachzählung der Korrektur steht im selben Plan |
| Reihenfolge der Commits und inert/nicht inert (`SPEC-052` Punkt 12) | geprüft, ohne Befund: `docs/maintainer/` (Eintrag, `releasing.md`) ist nicht inert, `make gates` am Endstand `2b0869ff` in §7 belegt; `511428c9` ändert nur den Plan |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 3 |

Je Slice: `slice-go-toolchain-1-27-3` F-4 (LOW), F-9 (INFO); `slice-release-v0-2-0` F-1 bis F-3, F-5 bis F-8, dazu F-10 für beide.

**Finding-Klassen dieses Laufs:** Nachweis ohne Beleg-Datei, Aussage als belegt geführt · plan-folgt-korrektur-nicht · folge-slice-adresse-nimmt-nicht-an · zusage-im-kommentar-weiter-als-pruefung · Stand-Aussage an zwei Stellen verschieden · Carveout-Text nennt anderen Hash als der Tag

## Verdikt

**Merge-blockierend:** `slice-go-toolchain-1-27-3`: nein — nur LOW und INFO. `slice-release-v0-2-0`: ja für die Closure — F-1 und F-2 (MEDIUM) bleiben vor ihr zu beantworten; der Tag und der Lauf sind davon nicht berührt (der Lauf ist `success`, jede nachgefahrene Zeile des Eintrags zu Lauf, Anhängen, Digest und Binary stimmt).

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier separat.
