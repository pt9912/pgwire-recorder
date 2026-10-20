# Verifikationsbericht: slice-v1-abschluss-homebrew — 2026-10-20

**Rolle:** Verifier (Modul 11). **Frage:** Bauen wir es richtig, gegen Plan, DoD und Entscheidungen.

**Gegenstand:** Diff `8dbe6a14..impl-brew2` (Kopf `fd5ef0ef`), neun Dateien; Review-Report `2026-10-20-review-slice-v1-abschluss-homebrew.md` (F-1 bis F-7).

**Entscheidungen und Spezifikation:** [ADR-0053](../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) (Accepted), [ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md) (Accepted), `SPEC-067` Punkt 1c und 19, `SPEC-055` Punkt 3, `LH-FA-19`.

**Umgebung:** Dev-Container `pgr-dev2`, eigenes Worktree `wt-ver` (Zweig `ver-brew`), `TMPDIR` im Worktree, Platz vor den Läufen 33 GB frei, nach den Läufen keine Reste (`tmpver` leer, danach entfernt). Kein voller Lauf von `make release-gegenprobe` (Nutzerentscheid vom 2026-10-20), kein `make gates`.

## Ergebnis

**Konform, ohne HIGH und ohne MEDIUM.** Jedes Finding des Reviews ist im Code oder im Plan behoben, nicht nur behauptet (Tabelle unten). Die gefahrenen Mutationen färben genau die Fälle, die §7 nennt. Offen sind zwei LOW/INFO-Punkte und die Closure-Pflichten.

## Belege und eigene Läufe

| Prüfung | Ergebnis |
|---|---|
| §7 trägt die Belege des Implementers (Befund-Tabelle, Läufe, Zusage · Mutation · roter Test, Grenzen) | vorhanden (F-5 behoben) |
| Unveränderte Kopie, Gruppe `p19-*` (`GP_MUT=^ok$`) | grün (140 s unter Last) |
| Vier Vorlage-Mutationen `m-p19-vorlage-(pruefung-entfaellt|stichwort-tap-formel|vor-summen|nach-klon)` und `m-p19-token-in-url` (`GP_ENTDECKE=1`, 62 s) | grün; rot genau wie in §7: Prüfung entfällt, Stichwort und nach Klon je `p19-vorlage-fehlt`, `-leer`, `-verzeichnis`, `-vor-tap`; vor Summen `p19-vorlage-nach-summen`; token-in-url `p19-token-nicht-im-argument`, `p19-token-nicht-in-url-und-config` |
| Stichprobe je Zusagen-Art (95 s): `m-p19-vorabversion-nur-rc`, `push-force`, `aktuell-immer-neu`, `hoeher-ueberschrieben`, `summen-altes-schema-angenommen`, `vorlage-name-mit-version-linux-amd64`, `vorlage-mit-windows`, `summen-form-63`, `int-falle-fehlt`, `token-im-argument-git` | grün; je Mutation dieselben roten Fälle wie in der Tabelle von §7 |
| `m-p1-homebrew-*` (alle vier) | grün; je ein roter Fall `p1-homebrew-…` wie in §7 |
| `make docs-check` am Worktree | 1299 Dateien, 0 Befunde |
| Vorlage gegen Text in `SPEC-067` Punkt 19 (4) | `cmp` gleich |
| Größenmarke §4 (`git diff --numstat 8dbe6a14..HEAD -- tools .github harness/README.md docs/maintainer`) | 860 + 49 = 909, Marke 1300 nicht erreicht (wie §7) |

Bewusstes Brechen: Für die DoD-Behauptungen mit Test (Zustände der Vorlage, Vorabversion, Token, Namen ohne Version in der Adresse, altes Schema) ist die Prüfung im Code per Mutation entfernt bzw. verändert und der benannte Fall rot gesehen; der Grund ist der richtige (Exit/Stichwort bzw. Aufruf weicht ab, nicht ein Folgefehler).

## Review-Findings: wirklich behoben?

| ID | Behauptung | Prüfung | Ergebnis |
|---|---|---|---|
| F-1 | Vorlage-Zustände mit Fall und Mutation, Stichwort im Kopfkommentar, Randform in Spec und §6 (h) | Code `release.sh` (`[ -n "$f" ]` nach den Summen, vor dem Klon), Kopfkommentar nennt `vorlage`, `SPEC-067` (4) und (9) tragen Zustände, Reihenfolge und Fälle, Läufe oben | behoben |
| F-2 | sechs Anker entfernt | `grep` in `tools/`: keiner der sechs mehr vorhanden; die 22 verbleibenden `# P19…`-Anker haben je mindestens einen Leser in der Gegenprobe | behoben |
| F-3 | `m-p19-token-in-url` | Lauf oben, rot genau die zwei Fälle der Spezifikation | behoben |
| F-4 | `SPEC-067` folgt der Annahme der ADR | Kopf und Punkt 19 (6) nennen die ADR `Accepted`; die Empfehlung zu Art und Rechten bleibt als Empfehlung | behoben |
| F-5 | Belege in §7 | vorhanden | behoben |
| F-6 | (INFO) Zusagen über Verhalten außerhalb des Klons in `releasing.md` | nicht geändert: siehe V-1 | offen |
| F-7 | Basis der Marke | §4 nennt `8dbe6a14`, Messung 909 bestätigt | behoben |

## Konformität

- **Entscheidung zum Tap** ([ADR-0053](../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md)): Eigener Job `homebrew` mit `needs: release`, nur bei `push`, `contents: read`, Geheimnis nur im Schritt `Homebrew`, kein `packages: write`; Zeilen des Jobs gleich Punkt 1c (Review und `p1-homebrew-*`). Der Token geht nur über `GIT_ASKPASS`, nicht in Argument, URL, Datei oder Ausgabe (Mutationen oben rot).
- **Namen ohne Version** ([ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md)): Vorlage und Schritt rechnen mit `pgwire-recorder-<os>-<arch>` ohne Version; altes Schema wird mit Exit 1 `summen` abgelehnt (`p19-summen-altes-schema` mit Mutation rot); Adresse je Plattform mit Mutation geprüft; die Zeile `version` in der Formel ist als Abweichung in §6 (a) begründet.
- **SPEC-067 Punkt 1c und 19:** Reihenfolge der Gründe (Ereignis, Tag-Form, Vorabversion, Tokens, `summen-laden`, `summen`, `vorlage`, Tap) im Code wie in der Spezifikation; Exit 2 für Eingaben, Exit 1 für Summen und Tap; Vorabversion ohne Aufruf von `gh` und `git`.
- **Handbuch und README:** unverändert (`git diff` über `docs/user` und `README.md` leer, kein Treffer zu `homebrew` dort); Handbuch bleibt Ist-Zustand.
- **Nichts Echtes im Netz:** Die Gegenprobe klont ein Bare-Repository im Temp-Verzeichnis (`file://`), `gh` ist eine Attrappe; die Vorgabe-URL des echten Taps wird nur über `url.<…>.insteadOf` geprüft. Kein Push in den echten Tap, kein Token-Wert im Repo.
- **Plan-vs-Code:** §3 nennt die sechs Dateien des Diffs, dazu `release.yml`; §8 zählt drei Liefer-Punkte und zwei Schichten, der Diff berührt Werkzeug und Wartungs-Doku (`docs/reviews` und der Plan sind Planung, keine Schicht). Kein Abweichen gefunden.

## Findings

| ID | Schwere | Befund | Pfad |
|---|---|---|---|
| V-1 | LOW | `releasing.md` sagt Zustand und Verhalten außerhalb des Klons ohne Beleg zu (§3.11): „im Tap liegt keine Formel“ (Zustand des echten Taps, aus dem Klon nicht prüfbar) und „„Re-run failed jobs“ wiederholt nur ihn“ (Verhalten von GitHub; `SPEC-067` Punkt 15 führt es als Grenze). Das ist F-6 des Reviews, weder behoben noch als akzeptiert eingetragen. Vorschlag: enger fassen („nach Angabe des Betreibers“ bzw. als Grenze nach Punkt 15 kennzeichnen) oder in §7 als akzeptiertes Negativ mit Grund führen. | `docs/maintainer/releasing.md` Zeile 25 und 202 |
| V-2 | INFO | `m-p19-vorlage-nach-klon` und `-stichwort-tap-formel` färben die vier Fälle `p19-vorlage-fehlt`, `-leer`, `-verzeichnis`, `-vor-tap`; `SPEC-067` (9) nennt für nach-Klon nur `-vor-tap`. Die größere Menge ist folgerichtig (nach dem Klon ist „Tap nicht geklont“ verletzt), die Spezifikation nennt sie enger. Kein Befund gegen den Code; Spezifikation könnte die Menge nennen (Architect). | `harness/vertraege/SPEC-067-release.md` Punkt 19 (9) |
| V-3 | INFO | Kopf des Plans trägt `Verantwortlich: pt9912 (Zweig impl-brew)`; die Nacharbeit lief auf `impl-brew2`. Deklaration, kein Sensor. | `docs/plan/planning/in-progress/slice-v1-abschluss-homebrew.md` |
| V-4 | INFO | DoD-Häkchen, Closure-Notiz mit Lerneintrag, Register und Risiko-Ausgänge sind noch offen („Wird bei Closure gefüllt“); die vier Risiken aus §6 tragen schon den Ausgang *weiter offen* mit Adresse `slice-erster-release-homebrew-nachweis`. Vor dem `git mv` nach `done/` zu schließen. Der Review-Report liegt vor; diese Verifikation ist der noch offene Beleg. | Plan §2, §7 |

## Nicht gefahren

- voller Lauf von `make release-gegenprobe` (Nutzerentscheid);
- Gruppen `p1-*` außer `m-p1-homebrew-*`, `p2-*` bis `p16-*` und ihre Mutationen;
- aus `p19-*` die Mutationen, die nicht in der Stichprobe stehen (rund 95 von rund 110, darunter Reihenfolge der Gründe, Tap-Randformen, Askpass, Commit-Identität, Push-Fehler, SIGTERM); §7 nennt sie als gefahren (Läufe 1 bis 10), von mir nicht wiederholt;
- `make gates` und `make test-integration` (der Host fährt sie am Endstand; der Beleg in §7 „grün, 303 s“ ist ungeprüft);
- Echt-Lauf: Rechte des Tokens, `gh release download` mit dem echten `gh`, Annahme der Formel durch Homebrew, Binary auf macOS — Grenzen nach `SPEC-067` Punkt 14 und 15, übergeben an `slice-erster-release-homebrew-nachweis`.

## Verdikt

Verifikation bestanden, Closure kann folgen, wenn V-1 (enger fassen oder als Negativ führen) erledigt oder bewusst akzeptiert ist. Kein Blocker.
