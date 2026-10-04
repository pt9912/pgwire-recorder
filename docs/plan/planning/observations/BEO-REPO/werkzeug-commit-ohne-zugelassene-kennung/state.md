**Stand:** offen

Für Slices ist die Ursache behoben: Der Träger `.githooks/commit-msg` nimmt seit
[ADR-0025](../../../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) die
Kennung jedes Slice im Index an. Für Wellen nicht: eine Message, die nur eine
Welle-Kennung trägt, lehnt der Träger weiter ab.
