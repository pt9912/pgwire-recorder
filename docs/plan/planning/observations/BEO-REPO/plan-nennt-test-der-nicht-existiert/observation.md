# Der Plan nennt einen Test, den es im Repo nicht gibt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

§3, §6 oder §7 eines Slice-Plans nennen einen Test mit Namen als Beleg einer Zusage, und
der Test steht nicht im Code, oder die Zeile behauptet für ihn eine rote Mutation, die nie
gefahren wurde. Die Regeln dafür gelten (`AGENTS.md` §3.10: je Zusage ein Test und die
Mutation selbst gefahren; §3.11: eine Plan-Zeile sagt nur zu, was ein Test prüft); keine
Stelle verlangt aber, die Namen gegen den Code zu halten. Die Verifikation fand es, weil
sie die Mutation selbst fuhr und den Namen suchte.

Abgrenzung zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: Dort ist der Wortlaut
weiter als die Prüfung; hier ist der Wortlaut genau, aber sein Gegenstand fehlt. Abgrenzung
zu `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht`: Der Beleg stand hier im Plan und
erreichte die Prüfer, er war falsch.
