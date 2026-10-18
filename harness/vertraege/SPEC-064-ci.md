# SPEC-064 — Die CI auf GitHub: `ci.yml`, Berechnung der Range, Konfiguration der Range-Prüfungen (`ci-range`, `ci-gegenprobe`)

Jede Integration (ein Pull-Request, ein Push auf einen Zweig) läuft auf GitHub durch `make gates` und durch die Prüfung ihrer Commit-Range,
unabhängig vom Klon und vom lokalen Hook: `.github/workflows/ci.yml` hat genau den Text von Punkt 1; die Range berechnet das Skript
`tools/harness/ci-range.sh` aus dem Ereignis (Teil 2); `make doc-commits` und `make doc-immutable` prüfen sie mit der Konfiguration von
Teil 3. `make ci-gegenprobe` hält alle drei: den Text der Datei, das Skript und die Konfiguration, die ohne die Gegenprobe still nichts
prüfte. Die CI ist **kein** Gate der Kette; sie ist ein Lauf auf GitHub, den dieses Repo nicht fährt. Warum die CI auf GitHub Actions läuft und
den lokalen Hook ergänzt, steht in [ADR-0050](../../docs/plan/adr/0050-ci-auf-github-actions.md); diese Kennung wiederholt es nicht.

**Stand der Entscheidung.** Der Nutzer hat [ADR-0050](../../docs/plan/adr/0050-ci-auf-github-actions.md) am 2026-10-18 angenommen und bestätigt: `runs-on: ubuntu-24.04` statt `ubuntu-latest` (Punkt 2, Grund dort) sowie den Verzicht auf `concurrency` und `workflow_dispatch` (Punkt 2).

**Schnittstelle.** Das Fragment `harness/mk/ci.mk` definiert `ci-gegenprobe` (`.PHONY`), das `tools/harness/ci-gegenprobe.sh` ruft und in
`GEGENPROBEN_CONTAINER` steht, nicht in `GATE_CHECKS` ([`SPEC-060`](SPEC-060-gegenproben.md) Punkte 1 bis 3). Es übergibt `DCHECK_REF` (aus `d-check.mk`)
und `DOCKER` als Umgebung. `ci-range.sh` liest fünf Umgebungsvariablen (Punkt 6) und schreibt die Range auf stdout; es ist kein Make-Ziel.

**Plattform** (`AGENTS.md` §3.12, Plattform von [`LH-QA-03`](../../spec/lastenheft.md#lh-qa-03--portabilität)). Die CI läuft auf `linux/amd64`
(`ubuntu-24.04`, Docker vorhanden); mehr ist nicht zugesagt. Andere Plattformen der Binaries baut `make test-container` per Cross-Kompilierung
innerhalb von `make gates`; ein Lauf auf `arm64`, macOS oder Windows als Runner ist nicht zugesagt (akzeptiertes Negativ; Grund: das Gate
`make gates` setzt Linux-Docker voraus, und die Veröffentlichung prüft `linux/arm64` als Image in ihrer eigenen Pipeline). `ci-range.sh` braucht bash und
`git` (≥ 2.30); `ci-gegenprobe` zusätzlich `make`, Docker, `mktemp` und GNU-Werkzeuge.

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes. Je Zusage und je Element einer Aufzählung gibt es einen Fall und eine Mutation, die nur die
Fälle ihrer Gruppe färbt (`AGENTS.md` §3.10); die unveränderte Kopie ist grün.

**Randform → Punkt** (Randformen aus §6 des Plans):

| Randform | Punkt |
|---|---|
| 1 Auslöser, doppelte Läufe | 1, 2 |
| 2 Pull-Requests aus Forks, Geheimnisse | 3 |
| 3 Abbruch älterer Läufe | 2 |
| 4 Zeitgrenze | 1 |
| 5 Cache der Docker-Layer | 2 |
| 6 Matrix, Runner | 1, 2 |
| 7 Rechte, Pins | 1, 3 |
| 8 „Grün“ für den Merge, Jobname | 4 |
| 9 Nachweis und Stop-Hook | 5 |
| 10 Ressourcen, Aufruf, Wiederholen | 2, 5 |
| 11 Range und Kennung | 6 bis 9, 11, 15 |
| 12 Gebündelte Prüfung | 1, 2 |
| 13 Netz, Docker Hub | 5 |
| Reihenfolge der Gründe von `ci-range.sh` | 9 |
| Konfiguration der Range-Prüfungen | 10 bis 13 |

## Teil 1: Die Datei `ci.yml`

1. **Der Text.** `.github/workflows/ci.yml` besteht, außer Leerzeilen und Kommentarzeilen (erstes Zeichen außer Leerzeichen `#`), aus genau diesen
   Zeilen in dieser Reihenfolge; Einzug und Anführungszeichen sind Teil des Textes. Für die Zeile mit `<SHA>` und `<Tag>` gilt: `<SHA>` sind 40 Hexziffern
   in Kleinschrift, `<Tag>` ein nicht leerer Text; beide stehen nicht in dieser Kennung (sie altern; die Prüfung der Pins gehört
   [`SPEC-063`](SPEC-063-workflows-check.md)).

   ```
   name: ci
   on:
     pull_request:
     push:
       branches:
         - '**'
       tags-ignore:
         - '**'
   permissions: {}
   jobs:
     gates:
       name: ci-gates
       runs-on: ubuntu-24.04
       timeout-minutes: 30
       permissions:
         contents: read
       env:
         CI_EVENT: ${{ github.event_name }}
         CI_HEAD: ${{ github.sha }}
         CI_PR_BASE: ${{ github.event.pull_request.base.sha }}
         CI_PUSH_BEFORE: ${{ github.event.before }}
         CI_DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}
       steps:
         - name: Checkout
           uses: actions/checkout@<SHA> # <Tag>
           with:
             fetch-depth: 0
             persist-credentials: false
         - name: Commit-Messages der Range
           run: |
             range="$(bash tools/harness/ci-range.sh)"
             make doc-commits RANGE="$range"
         - name: ADR-Immutabilität der Range
           if: ${{ !cancelled() }}
           run: |
             range="$(bash tools/harness/ci-range.sh)"
             make doc-immutable RANGE="$range"
         - name: Gates
           run: make gates
   ```

   Jede Änderung der Datei ist eine Änderung dieser Kennung. Fälle: `p1-gleich` (die Datei des Arbeitsbaums entspricht dem Text, grün); je Zeile des Textes ein
   Fall `p1-fehlt-<n>` an einer Kopie ohne diese Zeile und ein Fall `p1-geaendert-<n>` an einer Kopie, in der die Zeile um ein Zeichen verändert ist (rot);
   `p1-zusatz-<Wort>` an Kopien mit je einer zusätzlichen Zeile, die `concurrency:`, `workflow_dispatch:`, `pull_request_target:`, `workflow_run:`, `strategy:`,
   `continue-on-error: true` oder `${{ secrets.X }}` (Element je Wort) trägt (rot); `p1-zusatz-schritt` (ein zusätzlicher Schritt `uses: actions/cache@…` bzw.
   `actions/upload-artifact@…`, rot); `p1-pin-kurz` (SHA mit 39 Ziffern), `p1-pin-gross` (Großbuchstaben), `p1-pin-ohne-tag` (kein Kommentar), `p1-pin-andere-action`
   (rot); `p1-reihenfolge-schritte` (die Schritte `Commit-Messages` und `ADR-Immutabilität` vertauscht), `p1-reihenfolge-gates` (der Schritt `Gates` vor den beiden anderen).
   Mutationen: die Prüfung ignoriert eine Zeile (je Zeile eine Mutation, erzeugt in einer Schleife über den Text); die Prüfung lässt zusätzliche Zeilen zu;
   die Prüfung ignoriert die Reihenfolge; der Pin akzeptiert Großbuchstaben, 39 Ziffern oder einen fehlenden Kommentar (je eine). Die Gegenprobe trägt den Text
   wörtlich als eigene Kopie; weicht sie von diesem Punkt ab, ist das ein Fehler der Gegenprobe, den der Fall `p1-gleich` am Arbeitsbaum zeigt.
2. **Was der Text entscheidet** (Gründe kurz; die Abwägung steht in der ADR). *Auslöser:* `pull_request` und `push` auf Zweigen, Tags ausgenommen (Tags gehören der
   Veröffentlichung). Ein Pull-Request aus demselben Repo löst beide aus, zwei Läufe je Commit; das ist hingenommen (akzeptiertes Negativ; Grund: ein Pull-Request prüft
   die Range gegen die Basis, ein Push auf einen Zweig ohne Pull-Request prüft den Zweig, und der Lauf kostet auf einem öffentlichen Repo nichts). `workflow_dispatch`
   fehlt: der Wiederanlauf geschieht über die Oberfläche von GitHub, und ein Ereignis ohne Basis und ohne Vorgänger hätte keine Range (Punkt 6, Ereignis
   unbekannt). *Kein `concurrency`:* Ein abgebrochener Lauf eines Push verwirft die Range seines Vorgängers (der nächste Lauf prüft nur `Vorgänger..Stand` seines eigenen
   Push), und ein wartender Lauf ersetzt den älteren wartenden; beides ließe Commits ungeprüft. *Runner:* `ubuntu-24.04` fest statt `ubuntu-latest`, damit ein
   Wechsel des Images nicht still die Gates ändert (Docker-Version, Werkzeuge); der Wechsel ist eine Änderung dieser Kennung. *Zeitgrenze:* 30 Minuten, ohne Grenze je
   Schritt; die Überschreitung ist rot, ein Wiederholen im Workflow gibt es nicht (ein Wiederholen maskiert Flackern; Wiederanlauf von Hand). *Kein Cache, keine Matrix:*
   `make` fährt `docker build` über das `Dockerfile` ([ADR-0026](../../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md)); jeder Lauf ist kalt, das Ziehen der Images
   zählt in die Zeit; ein Cache wäre eine Messung und eine Änderung dieser Kennung. *Reihenfolge:* die billigen Range-Prüfungen zuerst; `ADR-Immutabilität` läuft nach einem
   roten ersten Schritt (`!cancelled()`), damit ein Lauf beide Befunde nennt; `Gates` läuft nur, wenn beide grün sind (kein Warten auf fünf Minuten, wenn die Range schon rot ist).
   `make gates` läuft ohne `-j` und unverändert (kein zweites Gate-Ziel nur für die CI); die gemessene Dauer steht im Beleg des Slice, nicht hier. Die Zuweisung
   `range="$(…)"` steht allein in einer Zeile, damit ein Fehler von `ci-range.sh` den Schritt beendet (in `echo`, `make … RANGE="$(…)"` oder `export` verschluckte ihn die
   Ersetzung, und die leere Range erreichte `make`).
3. **Fork-Pull-Requests und Geheimnisse.** Der Text kennt weder `pull_request_target` noch `secrets.` noch ein Token außer dem, das `actions/checkout` selbst nutzt; Rechte
   `contents: read`; `persist-credentials: false` lässt kein Token im `.git`, das ein Gate-Skript lesen könnte. Für einen Pull-Request aus einem Fork läuft der Text mit
   dem schreibgeschützten Token ohne Geheimnisse (das ist das Verhalten von GitHub für `pull_request`, nicht eine Zusage dieses Repos; ob ein Lauf eines Erstbeitragenden
   eine Freigabe braucht, ist eine Einstellung des Repos auf GitHub). Belegt durch Punkt 1 (Fälle `p1-zusatz-*`, `p1-geaendert-*` der Zeilen `permissions` und `persist-credentials`).
4. **Was „grün“ für den Merge heißt.** Der Status ist der Job `ci-gates` (Zeile `name: ci-gates`); sein Name ist Teil des Textes (Fälle `p1-fehlt-<n>` und
   `p1-geaendert-<n>` für diese Zeile). Ändert er sich, verliert ein Branch-Schutz seine Bedingung still. **Der Branch-Schutz ist eine Einstellung auf GitHub**, aus dem Klon
   nicht auditierbar und keine Zusage dieses Repos (Handlung des Betreibers); ohne ihn ist die CI ein Hinweis, kein Sperrriegel. `harness/README.md` nennt den Namen und die Handlung.
5. **Verhältnis zum lokalen Hook; Nachweis; Ressourcen.** Die CI **ergänzt** den Stop-Hook ([ADR-0038](../../docs/plan/adr/0038-gate-kette-voll-und-schnell.md)), sie ersetzt ihn nicht:
   sie läuft klonunabhängig und über die Range und fängt `--no-verify`; sie ersetzt keinen Nachweis des lokalen Hooks, und der Hook liest nichts aus ihr. Das `make gates`
   der CI schreibt `.harness/state/gates-passed.dateiliste` in den Arbeitsbaum des Runners; das ist folgenlos: der Text hat keinen Schritt, der den Ordner hochlädt, zwischenspeichert oder
   committet (Punkt 1, Fälle `p1-zusatz-schritt`). Die CI meldet sich nicht beim Docker Hub an (keine Geheimnisse); ein Fehler durch das Limit anonymer Pulls ist rot und bleibt
   es; er wird als solcher im Beleg des Slice festgehalten, und ein Wiederanlauf ist von Hand. Ob `make gates` auf dem Runner läuft (Docker-Netz der Integrationstests, `setsid`,
   Speicher, Zeit), zeigt nur ein Lauf; kein Gate dieses Repos sagt es voraus (akzeptiertes Negativ; Grund: nur dort messbar, wie bei [`SPEC-063`](SPEC-063-workflows-check.md) Punkt 11).

## Teil 2: Die Range (`tools/harness/ci-range.sh`)

6. **Eingaben.** Das Skript liest `CI_EVENT` (E1: `pull_request` oder `push`), `CI_HEAD` (E2: 40 Hexziffern, Kleinschrift), `CI_PR_BASE` (E3: 40 Hexziffern, nur bei
   `pull_request`), `CI_PUSH_BEFORE` (E4: 40 Hexziffern, bei `push`; 40 Nullen für einen neuen Zweig) und `CI_DEFAULT_BRANCH` (E5: ein Zweigname nach
   `git check-ref-format --branch`, nur bei `push` mit neuem Zweig oder unbekanntem Vorgänger). Eine Variable, die das Ereignis nicht liest, wird nicht geprüft. Fehlt eine
   gelesene Variable, ist sie leer oder ungültig, endet das Skript mit Exit 2 und `ci-range: <Name> fehlt oder ist ungültig`; ein anderes Ereignis (`workflow_dispatch`,
   `schedule`, leer) ebenso. Fälle: `p6-event-fehlt`, `p6-event-unbekannt`, `p6-head-fehlt`, `p6-head-ungueltig`, `p6-prbase-fehlt`, `p6-prbase-ungueltig`, `p6-before-fehlt`,
   `p6-before-ungueltig`, `p6-default-fehlt`, `p6-default-ungueltig`, `p6-ungenutzt` (Pull-Request mit ungültigem `CI_PUSH_BEFORE`: Range gegeben). Mutation je Element:
   die Prüfung der Variable entfällt.
7. **Formen.** Das Skript gibt die Range als `<Basis>..<Kopf>` mit aufgelösten Commit-SHAs (je 40 Hexziffern) aus. *Pull-Request:* Basis ist `CI_PR_BASE`, Kopf `CI_HEAD`; der Kopf darf
   ein Merge-Commit sein (der Lauf prüft den Merge-Stand, und die Prüfung der Commit-Messages nimmt Merge-Commits aus); liegt die Basis nicht im Klon, Exit 2. *Push:* Basis ist
   `CI_PUSH_BEFORE`, Kopf `CI_HEAD`. *Neuer Zweig* (`CI_PUSH_BEFORE` nur Nullen): Basis ist der Stand von `refs/remotes/origin/<CI_DEFAULT_BRANCH>`; fehlt er, Exit 2. *Vorgänger
   nicht mehr im Klon* (Force-Push): wie ein neuer Zweig, mit der Zeile `ci-range: Vorgänger <sha> nicht im Klon, Basis ist der Standard-Zweig` auf stderr (geprüft wird mehr, nicht weniger).
   *Leere Range* (null Commits nach `git rev-list --count`, etwa ein neuer Zweig auf dem Stand des Standard-Zweigs oder das erste Schieben des Standard-Zweigs): Exit 1 und
   `ci-range: Range <Basis>..<Kopf> ist leer`; eine leere Range prüft nichts und ist rot, nicht grün. *Flacher Klon* (`git rev-parse --is-shallow-repository` ist `true`): Exit 2
   und `ci-range: flacher Klon, der Checkout braucht fetch-depth: 0`. Fälle: `p7-pr`, `p7-pr-merge-kopf`, `p7-pr-basis-fehlt-im-klon`, `p7-push`, `p7-push-neu`,
   `p7-push-neu-ohne-standard`, `p7-push-vorgaenger-fehlt`, `p7-leer` (neuer Zweig auf dem Stand des Standard-Zweigs), `p7-flach`. Mutationen: Pull-Request nimmt `CI_PUSH_BEFORE`;
   Push nimmt `CI_PR_BASE`; Nullen gelten nicht als neuer Zweig (die Range `000…` ist unauflösbar); der Standard-Zweig fehlt ohne Exit 2; ein fehlender Vorgänger endet mit Exit 2
   statt auf den Standard-Zweig zu fallen; die leere Range endet mit Exit 0; der flache Klon wird nicht erkannt; die Ausgabe trägt Namen statt SHAs.
8. **Ausgabe und Ausgang.** stdout trägt bei Erfolg genau eine Zeile, die Range; alles andere (Hinweise, Fehler) steht auf stderr mit dem Präfix `ci-range:`; bei Exit ungleich 0 bleibt stdout leer.
   Exit 0 heißt eine nicht leere Range; 1 eine leere; 2 eine Eingabe, einen Klon oder eine Auflösung, die nicht möglich ist (auch ein `git`-Fehler). Fälle:
   `p8-stdout-eine-zeile` (auch beim Hinweis des Force-Push), `p8-fehler-stdout-leer` (jede Form mit Exit ungleich 0 ohne Ausgabe auf stdout). Mutationen: der Hinweis steht auf stdout;
   ein Fehler schreibt eine Teil-Range auf stdout.
9. **Reihenfolge der Gründe.** (1) Ereignis, (2) die vom Ereignis gelesenen Variablen in der Folge `CI_HEAD`, `CI_PR_BASE` bzw. `CI_PUSH_BEFORE`, `CI_DEFAULT_BRANCH`, (3) flacher Klon, (4) Auflösung von Basis und
   Kopf, (5) leere Range. Fälle: `p9-eingabe-vor-flach` (flacher Klon und ungültige Variable: Meldung der Variable), `p9-flach-vor-aufloesung` (flacher Klon, dessen Basis nicht aufzulösen ist:
   Meldung zum Klon), `p9-aufloesung-vor-leer`. Mutation je Paar: die Reihenfolge vertauscht.

## Teil 3: Konfiguration der Range-Prüfungen

10. **Ausgangslage und Zusage.** `make doc-commits` und `make doc-immutable` lesen ihre Regeln aus `.d-check.yml`. Ohne die Blöcke `commits` und `vcs` sind beide Module
    inert: sie melden „0 Befund(e)“, Exit 0, über jeder Range, und der Vorlauf `history-range-guard` fängt nur eine leere oder unauflösbare Range. Die CI wäre grün, ohne etwas zu
    prüfen. Diese Kennung legt deshalb die Blöcke fest (Punkte 11 und 12) und die Gegenprobe zeigt, dass sie wirken (Punkt 14).
11. **`commits`.** `.d-check.yml` führt `commits.id-patterns` mit genau diesen vier Mustern: `ADR-[0-9]{4}`, `LH-[A-Z]{2}-[0-9]{2}`, `MR-[0-9]{3}`, `(slice|welle)-[a-z0-9]+(-[a-z0-9]+)*`,
    und `commits.exempt-pattern` mit `^(Merge |Revert )` (der Betreff ohne Kennung ist erlaubt, wie im lokalen Hook). Geprüft wird die Anwesenheit einer Kennung nach diesen Mustern
    in jeder Nicht-Merge-Message der Range; nicht geprüft werden die Existenz der genannten Slice- oder Welle-Datei (der lokale Hook prüft sie am Index; akzeptiertes Negativ, Grund:
    ein Commit der Range kann einen Slice nennen, der in einem späteren Commit derselben Range umbenannt wurde, und eine Prüfung am Index der Spitze lehnte alte Commits zu Unrecht ab),
    und nicht die Struktur-Kennungen (`SPEC-`, `ARC-`), die der lokale Hook ablehnt ([`SPEC-050`](SPEC-050-struktur-kennungen-im-commit.md)); die CI lehnt sie weder ab noch zählt sie als Kennung.
    Fälle (Punkt 14): `p11-ohne-kennung` (rot, `commit-untraceable`); je Muster ein Fall mit einer Message, die nur diese Kennung trägt (`p11-adr`, `p11-lh`, `p11-mr`, `p11-slice`, `p11-welle`: grün);
    `p11-revert` und `p11-merge-betreff` (Betreff `Revert "x"`, `Merge x` ohne Kennung: grün); `p11-kennung-im-rumpf` (Kennung nur in Zeile 3: grün); `p11-kennung-im-kommentar`
    (Kennung nur in einer `#`-Zeile: rot). Mutationen: der Block `commits` fehlt (färbt `p11-ohne-kennung`, `p11-kennung-im-kommentar`); ein Muster fehlt (je Muster, färbt den Fall dieses Musters),
    die Ausnahme fehlt (färbt `p11-revert`, `p11-merge-betreff`).
12. **`vcs`.** `.d-check.yml` führt `vcs.paths` mit `docs/plan/adr/[0-9]*.md`, `vcs.immutable-when` mit `^\*\*Status:\*\* Accepted`, `vcs.exclude-sections` mit `[Geschichte]`, `vcs.status-line` mit
    `^\*\*Status:\*\*`, `vcs.head-allow` mit `^\*\*Status:\*\* (Accepted|Superseded by )` und `vcs.ignore-link-targets` mit `true`. Das heißt: eine ADR, die in der Basis `Accepted` war, ändert sich in der
    Range nur im Abschnitt `Geschichte`, in einem nachgezogenen Link-Ziel und im Übergang zu `Superseded by`; eine ADR in anderem Status ändert sich frei; der Index `docs/plan/adr/README.md` trägt keine Status-Zeile und ist damit nie unveränderlich. Das ist [`AGENTS.md`](../../AGENTS.md) §3.5, mechanisch. Fälle (Punkt 14): `p12-kern` (Text einer `Accepted`-ADR geändert: rot, `core-drift-vcs`), `p12-geschichte` (nur eine Zeile
    im Abschnitt `Geschichte`: grün), `p12-ersetzt` (Status wird `Superseded by [ADR-…]`: grün), `p12-ersetzt-anders` (Status wird `Rejected`: rot), `p12-vorschlag` (`Proposed`, Text geändert: grün),
    `p12-geloescht` (`Accepted`-ADR gelöscht: rot), `p12-link-ziel` (Link-Ziel einer umgezogenen Datei nachgezogen: grün). Mutationen: `vcs` fehlt (färbt `p12-kern`,
    `p12-ersetzt-anders`, `p12-geloescht`); `exclude-sections` fehlt (färbt `p12-geschichte`); `head-allow` fehlt (färbt `p12-ersetzt-anders`) bzw. ist zu eng (färbt `p12-ersetzt`);
    `ignore-link-targets` fehlt (färbt `p12-link-ziel`); `paths` leer (färbt wie `vcs` fehlt).
13. **`workflows`.** Die Bindung der Pin-Prüfung (`modules:` und `workflows.dir`) steht in `.d-check.yml`, sobald `ci.yml` existiert; sie ist Gegenstand von [`SPEC-063`](SPEC-063-workflows-check.md)
    Punkt 8, und `make workflows-check` hält sie am Arbeitsbaum. Diese Kennung wiederholt sie nicht. Das Modul meldet ohne Datei und ohne `uses:` einen Befund (leere Prüfmenge): der Eintrag kommt mit
    der ersten Workflow-Datei, nicht früher.
14. **Aufbau der Gegenprobe.** `make ci-gegenprobe` legt je Fall einen Temp-Baum an und räumt ihn weg; der Arbeitsbaum bleibt unberührt. *Punkt 1:* Kopien der Datei in Temp-Verzeichnissen, geprüft
    von der Funktion, die den Text hält. *Teil 2:* kleine Git-Repositories in Temp-Verzeichnissen mit Commits, Zweigen, einem Remote-Zweig `refs/remotes/origin/main` (per `git update-ref`) und einem flachen Klon
    (`git clone --depth 1`); `ci-range.sh` läuft mit den Variablen des Falls. *Teil 3:* je Fall ein Git-Repository mit Kopien von `.d-check.yml`, `d-check.mk`, `harness/mk/doc-gate.mk` und
    `tools/harness/history-range-guard.sh` und einem Makefile, das beide Fragmente einbindet; `make doc-commits RANGE=…` und `make doc-immutable RANGE=…` laufen mit dem echten Image von d-check (`DCHECK_REF`)
    gegen die Range des Falls; erwartet ist Exit 0 ohne Befund oder Exit 1 mit dem Code des Falls. Mutationen von Teil 3 ändern die Kopie von `.d-check.yml`. Voraussetzung: Docker und das Image von d-check
    (der erste Lauf zieht es, danach läuft die Gegenprobe ohne Netz). Die Gegenprobe trägt keine Abdeckungs-Deklaration (keine Anforderung des Lastenhefts; [`SPEC-048`](SPEC-048-abdeckung.md)). Aufrufarten von
    `make` (`-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`) haben hier keinen Fall: kein Hook und kein Stand liest das Ergebnis, und die Sammelziele behandeln die Flags ([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 10);
    `ci-range.sh` läuft nie unter `make`.
15. **Grenzen** (akzeptierte Negative, jeweils mit Grund). Die Basis eines Pull-Requests kommt aus der Nutzlast des Ereignisses; wandert die Basis zwischen Ereignis und Lauf, prüft der Lauf mehr Commits,
    die inzwischen in die Basis gelangt sind, nie weniger. Hat ein Pull-Request einen Merge-Konflikt, startet GitHub keinen Lauf für `pull_request` (Verhalten von GitHub); der Push-Lauf des Zweigs bleibt. Die erste Probe auf
    GitHub prüft die Range des ersten Push; enthält sie ältere Commits, die gegen Teil 3 verstoßen, ist der Lauf rot, ohne dass die CI falsch läuft (Bestand wird nicht geändert). Ein Löschen der Blöcke aus
    Teil 3 färbt die Gegenprobe, nicht die CI; Grund: die CI kann ihre eigene Konfiguration nicht von ihrer Abwesenheit unterscheiden, und das Review sieht die Änderung von `.d-check.yml`. Ob ein Lauf auf GitHub grün
    endet, belegt ein Lauf, kein Gate; der Beleg steht im Slice, und diese Kennung sagt es nicht zu.
