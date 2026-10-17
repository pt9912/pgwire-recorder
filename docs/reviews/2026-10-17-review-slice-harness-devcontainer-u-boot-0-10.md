# Review-Report: slice-harness-devcontainer-u-boot-0-10 — 2026-10-17

**Review-Art:** Code (Doku) — gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-harness-devcontainer-u-boot-0-10`, Diff `8070a14d..HEAD` (`bf1d170e` Dokument, `641a5b14` Belege in §7), Zweig `dev-container-u-boot`.

**Skill:** `.harness/skills/reviewer.md` @ `641a5b14`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form.** Kennung, nicht Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext:**

- Slice-Plan `slice-harness-devcontainer-u-boot-0-10` (§1, §2, §3, §4, §6, §7 *Belege des Implementers*); Nehmer `slice-harness-devcontainer-einrichten` (§1, §3)
- `AGENTS.md` (§3.7, §3.9, §3.11, §3.12, §3.13, §3.14)
- keine `LH-*`-Kennung, keine ADR

**Eigene Läufe:** `make docs-check` (1029 Dateien, 0 Befunde), `make kopf-check` (Exit 0). `git diff --numstat 05fd1403..HEAD -- docs/developer`: 297 + 301 = 598. Suche im Dokument nach Name, E-Mail, Konto-Kennung des Nutzers und Schlüsselmustern: keine Treffer; die Werte der Variablen stehen nur als Länge (23, 28) und als `<Schlüssel>`/`<Token>`. Docker- und Gate-Läufe nicht gefahren (Auftrag); das Scratch-Gerüst ist entfernt, die Proben ließen sich nicht wiederholen, geprüft wurde das Dokument gegen die Tabelle in §7.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Summe nach der Metrik aus §4 ist 598 und liegt über der Marke 450; die Rückführung `in-progress` → `next` ist nach dem Wortlaut von §4 eingetreten („Die Summe nach der Metrik liegt über 450“). §7 nennt das und überlässt die Entscheidung dem Planner. Der Schnitt, den §4 vorsieht (Gerüst und Dokument zuerst, Messung als Folge-Slice), ist am Diff nicht mehr zu trennen, weil Dokument und Messung in denselben Abschnitten stehen. | `AGENTS.md` §3.14 | `docs/plan/planning/in-progress/slice-harness-devcontainer-u-boot-0-10.md` · „Marke 450 (nur Dokument)“ und „Summe 598“ | ja — `git diff --numstat 05fd1403..HEAD -- docs/developer` | Rückführungs-Marke überschritten, Entscheidung offen |
| F-2 | MEDIUM | Das Gerüst weicht vom Plan ab: `devcontainer`-CLI und `mk-override.sh` statt `compose.yaml` und `start.sh`, `appPort` im Override statt Port-Mapping. §1, §2 (Punkt 1), §3 (Tabelle der Dateien), §4 und §6 (Randformen 6, 9) führen weiter `compose.yaml`, `start.sh` und `u-boot down`/`up`, das Dokument sagt dagegen „`u-boot up` startet diesen Container nicht“ (`no service selected`, Exit 12). Der Plan folgt dem Diff nicht (der Planner zieht nach, laut Auftrag). | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-harness-devcontainer-u-boot-0-10.md` · „`$PGR_DEV/.devcontainer/start.sh`“ und „`u-boot down`, `u-boot up`“ | nein — kein Gate liest den Abgleich | Plan folgt der Korrektur nicht |
| F-3 | MEDIUM | Der Nehmer `slice-harness-devcontainer-einrichten` nimmt die Sendung („Quellenauflösung“, Neustart) in einer Form an, die das Dokument widerlegt: sein §1 setzt den Neustart mit `u-boot down`/`up` und ein „`u-boot up`“ voraus, seine Remote nennt `/workspaces/austausch.git` (das Dokument: `/workspaces/pgr-dev/austausch.git`), und die Zugangsvariablen kommen nach dem Dokument nur in Prozesse von `devcontainer exec`, nicht „beim Start“ und nicht in `ssh`-Sitzungen. Ein Lauf des Nehmers nach seinem jetzigen Text scheitert an Pfad und Aufruf. | `AGENTS.md` §3.13 | `docs/plan/planning/open/slice-harness-devcontainer-einrichten.md` · „ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/austausch.git“ und „`u-boot down`, `u-boot up`“ | ja — `grep` nach Pfad und Aufruf in beiden Dateien | Adresse nimmt nicht (mehr) an |
| F-4 | MEDIUM | Viele Aussagen des Dokuments haben in der Tabelle der Proben in §7 keine Zeile, obwohl §13 sie nicht als ungeprüft führt und der Plan „je Probe einen Lauf“ verlangt: Inhalt von `devcontainer.json`, `Dockerfile`, `docker-init.sh` (§3.4; §7 verweist für die Dateien auf das Dokument selbst, ein Zirkel); Bare-Repo, Eigentümer von `/workspaces`, `git`-Version 2.55.0, Klon, `make hooks-install`, Push und `fetch` (§6, §9); `fetch` aus zwei Containern mit zwei Schlüsseln (§10); `Docker Engine ≥ 24` und Node.js 22.18.0 am Host (§2); `sshd`-Konfiguration, `/etc/shadow`, `sudo -n -l`, `ss -ltn`, `AcceptEnv`, `cgroup memory.max`, Pseudo-Terminal-Lauf (§5, §7, §11). Ob sie gelaufen sind, ist aus dem Beleg nicht zu sehen. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „Gemessen am frischen Volume: `/workspaces` gehört `root`“ und „`ssh/id_ed25519` … `Docker Engine ≥ 24`“ | nein — der Verifier braucht die Rohausgaben der Proben | Gemessen-Aussage ohne Zeile im Beleg |
| F-5 | LOW | Die zwei Zeilen von `make gates` in §8 nennen beide 1701 s; §7 führt dagegen drei Läufe der Hostliste (Lauf 3 Exit 0), den Lauf auf „frischem Volume“ und einen „zuvor“ mit Cache aus Fehlläufen (5,424 GiB). Die Reihenfolge in §8 (frisch zuerst) und in §7 (Cache „zuvor“) widerspricht sich, der Satz „die Gleichheit der Sekunden ist Zufall“ ist eine Behauptung ohne Beleg, und die 5,424 GiB des zweiten Laufs stehen nur in §7. Ob „frisches Volume, Daemon ohne Images“ nach dem Erfolgslauf 3 der Hostliste noch zutraf, ist offen. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „die Gleichheit der Sekunden ist Zufall“ | nein | Messwert doppelt ohne unterscheidbare Bedingung |
| F-6 | LOW | §8 begründet den Neustart je Slice damit, das Zugangstoken sei „beim Start frisch“. Nach §7 liest `devcontainer exec` die Variable bei jedem Aufruf aus der Shell des Hosts, und „ein Neustart des Containers ist dafür nicht nötig“; der Neustart ändert am Token nichts. Die beiden Abschnitte tragen verschiedene Gründe. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „das Zugangstoken ist beim Start frisch“ | nein | Begründung widerspricht Messung im selben Dokument |
| F-7 | LOW | §11 schreibt „Die Anmeldung von Claude Code erreichte `api.anthropic.com`“, §7 „das Token wurde von `api.anthropic.com` geprüft“. Gelaufen sind nur Platzhalter (`Invalid API key`, `401 Invalid bearer token`); eine Anmeldung gab es nicht, und dass die Antwort von diesem Host kam, ist ein Schluss aus der Meldung, kein Mitschnitt (§7 nennt das DNS-Log nur für die Liste insgesamt). Der Satz kann als Beleg des Logins gelesen werden, den §13 als ungeprüft führt. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „Die Anmeldung von Claude Code erreichte `api.anthropic.com`“ | nein | Zusage weiter als die Probe |
| F-8 | LOW | §3.5 nennt die Zusätze „`cpus` und `mem_limit`“; das Skript setzt `--cpus` und `--memory` in `runArgs`, `mem_limit` ist ein Begriff von Compose. Die Zeile über dem Skript in §3.5 hat ein Leerzeichen vor dem Semikolon („nennt ; fehlt die Datei“). | Maintainability | `docs/developer/devcontainer.md` · „Name, `cpus` und `mem_limit`“ | nein | Begriff des Vorgänger-Gerüsts übrig |
| F-9 | INFO | Der echte Login (API-Schlüssel, Abo-Token) steht im Dokument nicht als geprüft: §7 (Absatz „lief nicht“), §13 (erster Punkt) und die Befehle für den Nutzer sind einheitlich; die Vorbehalte der Vorfassung (kein Sandbox, `sudo`, Privilegien, Hostschlüssel, Verschachtelung) sind in §1 und §11 erhalten. Das Dokument nennt weder Name noch E-Mail des Nutzers und keinen Wert eines Zugangs. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „Der echte Login mit einem gültigen Schlüssel oder Token **lief nicht**“ | nein | — |
| F-10 | INFO | Die Größe der Gerüst-Dateien (108 Zeilen alt gegen neu, dazu 134 erzeugt) liegt außerhalb der Metrik von §4; die Summe 598 betrifft nur das Dokument. Die Hälfte der Löschungen (301) ist der Austausch ganzer Abschnitte, nicht Umbruch. Zuständig für die Entscheidung über den Schnitt ist der Planner. | `AGENTS.md` §3.14 | `docs/plan/planning/in-progress/slice-harness-devcontainer-u-boot-0-10.md` · „Gerüst-Umfang“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/developer/devcontainer.md` §11, §12, §13: Vorbehalte der Vorfassung gegen den Diff (Risiko „Umbau verliert Vorbehalte“) | geprüft, ohne Befund |
| `docs/developer/devcontainer.md`: Kommentarklassen, Konjunktiv über verworfene Alternativen, Chronik (§3.7) | geprüft, ohne Befund |
| `docs/developer/devcontainer.md`: personenbezogene Angaben und Zugangsdaten | geprüft, ohne Befund |
| `docs/developer/devcontainer.md` §4, §5, §7, §10 gegen die Zeilen von §7 (Zeiten, Fehlermeldungen, Längen, Sichtbarkeit des Werts, Hostliste) | geprüft, ohne Befund außer F-5 und F-7 |
| Schicht-Abgrenzung (kein Produkt-Code, kein Harness-Werkzeug) | geprüft, ohne Befund: der Diff berührt nur `docs/developer/` und den Plan |
| `make docs-check`, `make kopf-check` | grün |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Rückführungs-Marke überschritten, Entscheidung offen · Plan folgt der Korrektur nicht · Adresse nimmt nicht (mehr) an · Gemessen-Aussage ohne Zeile im Beleg · Messwert doppelt ohne unterscheidbare Bedingung · Begründung widerspricht Messung im selben Dokument · Zusage weiter als die Probe · Begriff des Vorgänger-Gerüsts übrig

## Verdikt

**Merge-blockierend:** ja — vier MEDIUM. F-1 bis F-3 liegen beim Planner (Schnitt, Plan und Nehmer nachziehen, vom Auftrag so vorgesehen); F-4 betrifft den Implementer (Zeilen in §7 oder Aussage in §13). Der Inhalt des Dokuments ist ohne HIGH; der nicht gemessene Login steht an keiner Stelle als geprüft.
