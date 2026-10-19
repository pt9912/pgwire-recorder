# Verifikationsbericht: slice-erster-release-veroeffentlichung — 2026-10-19

**Rolle:** Verifier (Modul 8) — „Bauen wir es richtig?“ gegen Plan, DoD und Vertrag. Kein Review, keine Reparatur.

**Gegenstand:** `slice-erster-release-veroeffentlichung`, Fassung A (Tag und Binaries). Zweig `ver-rel-a` im Klon des Dev-Containers pgr-dev3, Stand `824e4a32` (= `impl-rel-a` nach Review-Runde 1), Basis `17a8c618`. Der letzte Commit (`824e4a32`) ändert nur den Plan; der Kopf des Codes ist `39d964e9`, auf den sich die Belege in §7 beziehen.

**Eingang:** DoD-Bestätigung und Belege des Implementers (Plan §7, *Belege des Implementers*); Plan §1, §3, §6, §7; `SPEC-067` Fassung 1a; [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md); Review-Report `2026-10-19-review-slice-erster-release-veroeffentlichung`.

**Grenzen des Laufs:** kein Push nach GitHub, kein Tag, kein Netz, keine Registry. Die volle Gegenprobe und `make gates` sind nicht wiederholt (engster Sensor; der Beleg des Implementers steht in §7 und gilt für `39d964e9`, der Unterschied zu `824e4a32` ist eine Plan-Datei).

---

## 1. Eigene Läufe (nicht die Behauptung, die Ausgabe)

| Lauf | Ergebnis |
|---|---|
| `make docs-check` | Exit 0, `d-check: 1209 Datei(en) geprüft, 0 Befund(e)` |
| `make workflows-check` | Exit 0, `2 Datei(en) geprüft, kein Befund` |
| `make kopf-check` | Exit 0 |
| `git diff --numstat 17a8c618..HEAD -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` (Summe) | **1425** (wie in §7 behauptet; `release-gegenprobe.sh` 1037, `release.sh` 231, `releasing.md` 83, `release.yml` 38, `release.mk` 12, `SPEC-067` 22, `harness/README.md` 1, `VERSION` 1) |
| Gegenprobe, unveränderte Kopie und 8 Mutationen (`GP_MUT`, `GP_ENTDECKE=1`) | `release-gegenprobe: grün (8 Läufe)`; die roten Mengen stimmen mit der Tabelle in §7 überein (unten) |
| Text von `release.yml` gegen `SPEC-067` Punkt 1a | Zeile für Zeile gelesen: gleich (38 Zeilen, Pin `actions/checkout@3d3c42e5… # v7.0.1`, `group: release-${{ github.event_name }}`, `cancel-in-progress: false`, `permissions: {}` oben, Job nur `contents: write`, `GH_TOKEN` nur im Schritt `GitHub-Release`, der Schritt nur bei `push`) |

Stichproben der Mutationen (Soll aus §7, Ist gemessen):

| Mutation | Ist (rot) | Soll |
|---|---|---|
| `m-p9-view-muster-geweitet` (F-1) | `p9-view-anderer-fehler` | gleich |
| `m-p3-pipe-abgeschnitten-angenommen` (F-5) | `p3-notes-pipe` | gleich |
| `m-p1-zeile-9-optional` (F-2, Zeile `group:`) | `p1-fehlt-9` | gleich |
| `m-p9-name-fehlt-windows_arm64` (F-4) | 20 Fälle, darunter `p9-anhaenge`, `p4-aufruf`, `p11-alte-marke` | „20 Fälle“ — gleich |
| `m-p9-vergleich-entfernt` | `p10-nicht-im-log`, `p9-release-existiert-abweichend` | gleich |
| `m-p3-zelle-leer-angenommen-lauf` | `p3-notes-zelle-leer` | gleich |
| `m-p4-summe-entfernt` | `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende`, `p4-summe-falsch` | gleich |

## 2. Bewusstes Brechen (Modul 11) — Fix selbst entfernt, nicht die Mutationsnamen des Implementers

Je DoD-Punkt, der sich auf einen Test beruft, die Stelle im Code von Hand gebrochen (Datei im Klon verändert, `make`-freier Lauf der Gegenprobe `GP_MUT='^ok$'`, danach `git checkout`; Arbeitsbaum wieder sauber):

| Fix / Zusage | Bruch von Hand | Ergebnis |
|---|---|---|
| F-1: nur `release not found` heißt „nicht vorhanden“ | `grep -q 'release not found'` → `grep -qi 'not found'` in `release.sh` | rot: `p9-view-anderer-fehler` — aus dem richtigen Grund (das zweite Bein `HTTP 404: Not Found` führt zu `create`) |
| F-5: Zeile mit weiterem `\|` ist Fehler `notes` | die Zeile mit der Marke `# P3pipe` aus `pruefe_notes` entfernt | rot: `p3-notes-pipe` |
| F-2: Gruppe `release-${{ github.event_name }}` (DoD 3, `SPEC-067` Punkt 1 und 13) | `group: release` in `release.yml` | rot: `p1-gleich`; zusätzlich `p13-gegenprobe-ende` (Nebenwirkung: die Kopie der Gegenprobe prüft den Arbeitsbaum, siehe V-5) |

Die Tests sind ohne den Fix aus dem richtigen Grund rot. Die Behauptung „je Zusage ein roter Fall“ trägt für die Stichprobe.

## 3. DoD-Konformität (Plan §2)

| DoD-Punkt | Befund | Urteil |
|---|---|---|
| Liefer-Punkt 1, Eingaben (`SPEC-067` Punkt 3 (1) bis (6)) | `release.sh eingaben` prüft in der Reihenfolge des Vertrags (Ereignis/Ref, Form, `VERSION`, Gleichheit, Zeile der Änderungshistorie, Standard-Zweig); `VERSION` = `0.1.0-dev`, `make binaries VERSION=` über `harness/mk/container.mk` (`VERSION ?= dev`). Fälle `p3-*` mit Mutationen je Prüfung, je Paar, je Form-Element belegt (Stichprobe rot wie behauptet). Keine Meldung gibt ein Geheimnis aus. | erfüllt (Code und Test); Closure-Haken offen |
| Liefer-Punkt 2, Binaries am GitHub-Release (Punkte 4 und 9) | Dateiliste = sechs Binaries + `SHA256SUMS` (`erwartete_namen`), `--verify-tag`, `--prerelease` genau bei Vorabversion, `--repo` an jedem `gh`, fehlende Anhänge ergänzt, abweichende Summe rot. Handbuch und `README.md` unverändert (kein Diff), im Plan begründet. **Anlegen des Releases auf GitHub nicht belegt** — der Plan sagt das selbst (nur der erste echte Tag zeigt es). | Code/Test erfüllt; Lauf auf GitHub offen (siehe V-1) |
| Liefer-Punkt 3, Trockenlauf, Wiederanlauf, Fehlertabelle (Punkte 1a, 2, 11, 12) | Text 1a gleich; nur `v*` und `workflow_dispatch`; Trockenlauf ruft weder `gh` noch legt er einen Release an (`p2-trockenlauf-ohne-aussenwirkung`, Schritt `GitHub-Release` mit `if: push`). `releasing.md` §11 trägt die drei Zeilen `eingaben`, `binaries`, `github-release` und deckt sich mit der Tabelle des Vertrags; §3 sagt jetzt, dass „Probe-Tag“ nur das Verfahren der Formel betrifft (F-3 behoben). `make workflows-check` und `make docs-check` grün. **Zwei Belege der DoD stehen aus:** der Trockenlauf auf GitHub (Lauf-Kennung, Dauer) und der Beleg „ein roter Lauf ist nach der Tabelle wieder angestoßen“. | Code/Test erfüllt; Belege auf GitHub offen (V-1) |
| `make gates` grün | §7: Exit 0, 285 s, am Kopf des Codes `39d964e9`; nicht von mir wiederholt, durch den Plan-Commit `824e4a32` nicht berührt (nur Datei unter `docs/plan/planning/`) | Beleg vorhanden, nicht nachgefahren |
| Review durchgeführt | Report liegt vor (Runde 1: 2 MEDIUM, 3 LOW, 2 INFO); die Behebungen (`108c006b`, `39d964e9`) sind nicht erneut begutachtet, sondern hier in Stichproben (Abschnitt 2) verifiziert | Report vorhanden; kein zweiter Durchgang (V-3) |
| Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen | noch nicht geschrieben; §6 Risiko 1 steht „offen bis Closure“. §7 trägt Belege, aber keine Closure-Notiz. Erwartbar vor der Closure, kein Befund gegen die Lieferung. | offen (Closure) |

Haken in der DoD sind alle leer; das ist vor der Closure der Sollzustand (Haken sind Bedingung für `done/` und kommen mit der Closure).

## 4. Entscheidungs-Konformität

- **[ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md):** Träger GitHub Actions, nur Tags `v*` und `workflow_dispatch`, **kein Probe-Tag** (Trockenlauf statt dessen), nur `GITHUB_TOKEN` in Fassung A, kein `environment:`: eingehalten (Datei, `releasing.md` §3 und §6, Plan §1).
- **`SPEC-067` Fassung 1a:** Text, Reihenfolge der Gründe (1) bis (6), Schlüsselwörter, Marken in `$RUNNER_TEMP/release`, Ausgangscodes (2 Eingabe/Vorbedingung, 1 Schritt), `gh` mit `--repo` und ohne `GH_TOKEN` Exit 2, Abbruch (`SIGINT`/`SIGTERM`) und Gleichzeitigkeit (Gruppe, Grenze als akzeptiertes Negativ): im Code wiedergefunden. Die Entscheidungen des Architect `5f9fddb5` und `eec9c839` stammen aus dessen Commits; der Implementer hat den Vertragstext nicht selbst geändert.
- **`AGENTS.md`:** §3.10 (Mutationstabelle, zwei grüne Mutanten mit Invariante/Ort/Sonde; Sonden im Review nachgefahren), §3.12 (Randformen vorab entschieden, Rückgabe an den Architect erfolgt), §3.13 (keine neue Sendung), §3.3 (kein Move im Diff), §3.9 (Plan §1, §3, §6 folgen dem Diff; §3 listet `SPEC-067` nicht als Implementer-Datei, §7 erklärt die Merges des Architect), §3.11 (Handbuch/README unverändert, Ist-Zustand): kein Verstoß gefunden.

## 5. Plan-vs-Code-Diff

| Plan §3 | Diff | Abweichung |
|---|---|---|
| `release.yml` (neu) | vorhanden, 38 Zeilen | keine |
| `release.sh` (neu) | vorhanden, Schritte `eingaben`, `binaries`, `github-release`, Marken wie beschrieben | keine |
| `VERSION` (neu) | `0.1.0-dev` | keine |
| `release-gegenprobe.sh`, `release.mk`, Eintrag `GEGENPROBEN_CONTAINER` | vorhanden (`release.mk` hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`) | keine |
| `harness/README.md` §Sensors | eine Zeile `make release-gegenprobe` | keine |
| `docs/maintainer/releasing.md` §1, 2, 3, 5, 6, 9, 11 | geändert (59/24) | keine |
| Handbuch, `README.md` | geprüft, unverändert | keine (kein Diff) |
| (nicht in §3) `SPEC-067` | 22 Zeilen aus den Commits des Architect | in §7 erklärt, keine Abweichung des Implementers |

Im Diff liegt nichts, was §3 nicht nennt, außer den Architect-Merges.

## 6. Größe (nur benannt)

Gemessen **1425** (Metrik nach `AGENTS.md` §3.14, Basis `17a8c618` und `11a2ba91` liefern dieselbe Summe) gegen Schätzung des Architect 820 bis 1070, Meldeschwelle 1000, **Marke 1300**: Marke und Meldeschwelle sind überschritten, die Rückführungs-Bedingung `in-progress → next` ist nach dem Wortlaut von Plan §4 eingetreten. §7 nennt das selbst; die datierte Zeile „weiter“ oder „schneiden“ steht aus (Entscheidung der Koordination, nicht des Verifiers). Die Gegenprobe trägt 1037 von 1425 Zeilen (73 %), die Faktor-2-Annahme der Schätzung liegt unter der Messung (gemessen rund 4,5 : 1).

## 7. Findings

| ID | Schwere | Befund | Quelle | Ort |
|---|---|---|---|---|
| V-1 | MEDIUM (closure-blockierend, kein Defekt des Codes) | Zwei Belege der DoD Punkt 3 stehen aus und sind nur aus dem Klon nicht zu liefern: (a) der Trockenlauf auf GitHub auf dem Stand der Lieferung (Lauf-Kennung, Dauer, `SPEC-067` Punkt 12); (b) „ein roter Lauf ist nach der Tabelle wieder angestoßen (Beleg in §7)“. (a) ist in §7 ausdrücklich als offen benannt („nicht gefahren … Beleg steht aus“) — korrekt, nicht als erledigt dargestellt. (b) benennt §7 **nicht** als offen; im Plan steht dafür nur der Attrappen-Beleg (`p11-wiederanlauf-release-*`), der kein Lauf auf GitHub ist. Auch das Anlegen des GitHub-Release ist nicht belegt (bewusst, erster echter Tag, `SPEC-067` Punkt 12). Dazu gehört das Risiko „Rechte des Trockenlaufs auf GitHub“ (§6, „offen bis Closure“) und der nicht gemessene Wortlaut `release not found` von `gh` (akzeptierte Grenze, `SPEC-067` Punkt 9). | Plan §2 Punkt 3; §6; §7 | Plan §7 *Trockenlauf auf GitHub* |
| V-2 | LOW | Größe: Marke 1300 gemessen 1425 überschritten, Entscheidung „weiter“/„schneiden“ ohne datierte Zeile (siehe Abschnitt 6). Nur benannt. | Plan §4 | Plan §7 *Messung der Größe* |
| V-3 | LOW | Die Behebungen von F-1 bis F-5 (`108c006b`, `39d964e9`) hat der Reviewer nicht erneut gesehen; sie sind hier nur durch Stichproben abgedeckt (Abschnitt 1 und 2), nicht durch einen zweiten Review-Durchgang. F-3 (`releasing.md` §3) ist am Text geprüft: behoben. F-2 ist durch Entscheidung des Architect gelöst (Gruppe getrennt, Grenze in `SPEC-067` Punkt 13), das GitHub-Verhalten selbst bleibt ungeprüft (akzeptiertes Negativ). | Review-Report F-1 bis F-5 | — |
| V-4 | INFO | `schritt_binaries` entfernt ein vorhandenes `$RUNNER_TEMP/release/dist` nicht; ein zweiter Aufruf im selben Verzeichnis scheitert an `make binaries` (`DIST` muss fehlen oder leer sein). Gelesen, nicht gefahren; in der Pipeline nicht erreichbar (jeder Re-run läuft auf einem neuen Runner), nur beim Aufruf von Hand. Nicht im Vertrag zugesagt. | `SPEC-055` Punkt 3 | `tools/harness/release.sh` · `schritt_binaries` |
| V-5 | INFO | Unter einem veränderten Arbeitsbaum (Fix von Hand gebrochen) färbt sich zusätzlich `p13-gegenprobe-ende`, weil die Kopie der Gegenprobe im Fall `ok` den Arbeitsbaum prüft. Kein Fehler, nur beim Bruch von Hand zu beachten. | — | `tools/harness/release-gegenprobe.sh` |
| V-6 | INFO | Zählung `docs-check` 1209 Dateien (§7: 1208 am Kopf `39d964e9`; der Unterschied ist der Review-Report bzw. spätere Planänderung). Beide ohne Befund. | — | — |

Kein HIGH. Keine DoD-Verletzung im Code oder in den Tests gefunden.

## 8. Urteil

**Code- und Test-Teil der DoD (Liefer-Punkte 1 bis 3): bestätigt** — die Stichproben stimmen mit den Belegen des Implementers überein, der Fix von F-1, F-2 und F-5 hängt an einem Test, der ohne ihn aus dem richtigen Grund rot ist, `docs-check`, `workflows-check` und `kopf-check` sind grün, Plan und Diff stimmen überein, Entscheidungen ([ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md), `SPEC-067` 1a) sind eingehalten.

**Slice noch nicht schließbar:** die Belege auf GitHub (V-1) stehen aus und sind als offen benannt bzw. hier benannt; die Größenentscheidung (V-2) und die Closure (Notiz, Lerneintrag, Register, Risiko-Ausgänge, Haken, Paarungen) stehen aus. Der Trockenlauf-Beleg gehört in §7 vor das `git mv` nach `done/`; er ist Handlung des Betreibers oder des Nutzers (Push nach GitHub ist dem Verifier untersagt).

**Übergabe an den Planner:** V-1 (Belege auf GitHub einholen oder als Ausgang festlegen; (b) in §7 als offen benennen), V-2 (Entscheidung „weiter“/„schneiden“ als datierte Zeile), V-3 (ob ein zweiter Review-Durchgang gewünscht ist).
