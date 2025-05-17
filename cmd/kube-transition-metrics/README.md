# kube-transition-metrics

The `kube-transition-metrics` controller exports Pod life-cycle metrics in JSON
format to stdout.

## Usage

```txt
Usage of kube-transition-metrics:
      --emit-partial                       Emit partial statistics for pods that have not yet become Ready and image pulls that have not yet completed. When set to false, pods that never become Ready and image pulls that never complete will not be included in the statistics. Partial statistics will always be emitted for pods that are deleted before they become Ready. When set to true, multiple statistics will be emitted for the same pod/image pull. (ADVANCED)
      --image-pull-cancel-delay float      The delay before canceling an image pull routine to ensure events are flushed (ADVANCED) (default 3)
      --kube-watch-max-events int          The Kubernetes Watch maximum events per response (ADVANCED) (default 100)
      --kube-watch-timeout int             The Kubernetes Watch API timeout (ADVANCED) (default 60)
      --kubeconfig-path $KUBECONFIG        The path to the kube configuration file, if it's not set the value of $KUBECONFIG will be used, if that's not set `$HOME/.kube/config` will be used.
      --listen-address /metrics            The host and port for HTTP server delivering prometheus metrics over /metrics and pprof profiling over `/debug/pprof` endpoints. (default "127.0.0.1:8080")
      --log-level string                   The global logging level, one of "trace", "debug", "info", "warn", "error", "fatal", "panic", "disabled", or "" (empty string). This option'svalues are case-insensitive. Setting a value of "disabled" will result inno metrics being emitted. (default "INFO")
      --statistic-event-queue-length int   The maximum number of queued statistic events (ADVANCED) (default 1000)
```
