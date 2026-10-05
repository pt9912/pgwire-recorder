**Vorgang:** slice-extended-query-replay
**Fund:** Dass die Verbindung nach einer Extended-Abweichung endet, fing kein Test; ein Adapter, der weiterliest, blieb grün, weil der Test an der Frist des Clients endete (Review F-322, Mutation A7).
**Fund (Folge-Reviews):** Die Ordnung von Wächter, Lesefrist und `geweckt` beim Herunterfahren hatte keinen fangenden Test (F-337, M1 grün); die Herleitung des SQL in der Diagnose war an mehreren Stellen ungeschützt (F-340, D1 und D2 grün), zuletzt die Bedingung für eine einfache Anfrage am Cursor (F-349, M1 grün).
**Fund (Verifikation):** Auf den neuen Pfaden der Diagnose hielt kein Test die Zusage „nie Parameterwerte“ (V-31, Mutation NG2 grün). Gemeinsame Form: eine Zusage der Diagnose ohne Test, der ihre Bedingung bricht.
