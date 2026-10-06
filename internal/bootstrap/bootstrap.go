package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/postgres"
	"github.com/pt9912/pgwire-recorder/internal/adapters/driven/recording"
	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/cli"
	"github.com/pt9912/pgwire-recorder/internal/adapters/driving/pgwire"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/model"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/ports/driving"
	"github.com/pt9912/pgwire-recorder/internal/hexagon/services"
)

var (
	_ driving.Recorder = (*services.RecordService)(nil)
	_ driving.Replayer = (*services.ReplayService)(nil)
)

// Run führt einen Aufruf aus und liefert den Exit-Code (SPEC-013 bis SPEC-019).
// `record` und `replay` laufen, bis ctx endet; danach endet jede Verbindung nach
// ihrer laufenden Interaktion. Die Signalbehandlung liegt beim Aufrufer.
func Run(ctx context.Context, args []string, version string, stdout, stderr io.Writer) int {
	cmd, err := cli.Parse(args, stdout)
	if errors.Is(err, cli.ErrHelp) {
		return 0
	}
	if err != nil {
		return fail(stderr, err)
	}

	switch cmd.Name {
	case "version":
		fmt.Fprintln(stdout, "pgwire-recorder", version)
		return 0
	case "record":
		return record(ctx, cmd.Record, logger(stderr, cmd.Record.LogLevel), stderr)
	case "replay":
		return replay(ctx, cmd.Replay, logger(stderr, cmd.Replay.LogLevel), stderr)
	default:
		return fail(stderr, model.Errorf(model.CodeUsage, nil, "unbekanntes Kommando %q", cmd.Name))
	}
}

func record(ctx context.Context, o cli.RecordOptions, log *slog.Logger, stderr io.Writer) int {
	service, err := services.NewRecordService(ctx, &postgres.Upstream{Address: o.Upstream}, recording.YAML{}, o.Output, o.Force)
	if err != nil {
		return fail(stderr, err)
	}
	l, err := pgwire.Listen(o.Listen)
	if err != nil {
		return fail(stderr, err)
	}
	log.Info("record gestartet", "listen", l.Addr().String(), "upstream", o.Upstream)

	server := pgwire.NewRecordServer(service, log)
	done := make(chan struct{})
	go func() {
		server.Serve(ctx, l)
		close(done)
	}()

	<-ctx.Done()
	l.Close()
	<-done

	if err := service.Finish(context.Background()); err != nil {
		return fail(stderr, err)
	}
	log.Info("record beendet", "output", o.Output)
	return exitCode(server.FirstErrorCode())
}

func replay(ctx context.Context, o cli.ReplayOptions, log *slog.Logger, stderr io.Writer) int {
	var opts []services.ReplayOption
	if o.FailOnUnconsumed {
		opts = append(opts, services.FailOnUnconsumed)
	}
	service, err := services.NewReplayService(ctx, recording.YAML{}, o.Input, opts...)
	if err != nil {
		return fail(stderr, err)
	}
	l, err := pgwire.Listen(o.Listen)
	if err != nil {
		return fail(stderr, err)
	}
	log.Info("replay gestartet", "listen", l.Addr().String(), "input", o.Input)

	server := pgwire.NewReplayServer(service, log)
	done := make(chan struct{})
	go func() {
		server.Serve(ctx, l)
		close(done)
	}()

	<-ctx.Done()
	l.Close()
	<-done

	// Nie zugeordnete Sessions werden nach allen Verbindungsfehlern gemerkt und
	// bestimmen den Exit-Code nur ohne einen solchen (LH-FA-03.b, LH-FA-13.b).
	code := server.FirstErrorCode()
	w, err := service.Unassigned()
	if w != nil {
		log.Warn(w.Msg, "code", w.Code)
	}
	for _, m := range model.Meldungen(err) {
		log.Error("Fehler", "code", m.Code, "error", m.Text)
		if code == "" {
			code = m.Code
		}
	}
	log.Info("replay beendet")
	return exitCode(code)
}

// stufen bildet die Werte von --log-level auf die Stufen des Loggers ab
// (LH-FA-14.a).
var stufen = map[string]slog.Level{
	cli.LogError: slog.LevelError,
	cli.LogWarn:  slog.LevelWarn,
	cli.LogInfo:  slog.LevelInfo,
	cli.LogDebug: slog.LevelDebug,
}

// logger schreibt Log-Zeilen im Format logfmt nach stderr: time (RFC 3339 mit
// Millisekunden und Zonenversatz, Ortszeit), level, msg, dann die Attribute;
// er zeigt die Zeilen der Stufe und der strengeren (LH-FA-14.a).
func logger(stderr io.Writer, stufe string) *slog.Logger {
	return slog.New(slog.NewTextHandler(stderr, &slog.HandlerOptions{Level: stufen[stufe]}))
}

// exitCode ist 0 ohne Verbindungsfehler, sonst der Exit-Code der Klasse des
// ersten Verbindungsfehlers (LH-FA-13.b).
func exitCode(code string) int {
	if code == "" {
		return 0
	}
	return (&model.Error{Code: code}).ExitCode()
}

// fail schreibt je Meldung des Fehlers ihren Fehlertext als Zeile beim
// Prozessende nach stderr, unabhängig vom Log-Level, und liefert den Exit-Code
// der ersten (SPEC-034 §Ausgabe, LH-FA-14.a).
func fail(stderr io.Writer, err error) int {
	ms := model.Meldungen(err)
	for _, m := range ms {
		fmt.Fprintln(stderr, m.Text)
	}
	return ms[0].ExitCode()
}
