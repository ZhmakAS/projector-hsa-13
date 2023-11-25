# projector-hsa-13

## Benchmarking Results

```
Benchmarking Beanstalkd PUT -- 1000000(iterations)
============================================
Mean ----->  23787.029306 nanos
(90 Percentile) ----->  26449.000000 nanos
(99 Percentile) ----->  36719.000000 nanos


Benchmarking Beanstalkd RESERVE -- 1000000(iterations)
============================================
Mean ----->  25372.686266 nanos
(90 Percentile) ----->  28316.000000 nanos
(99 Percentile) ----->  48796.000000 nanos


Benchmarking Redis SET with RDB -- 1000000(iterations)
============================================
Mean ----->  24411.790478 nanos
(90 Percentile) ----->  26877.000000 nanos
(99 Percentile) ----->  38107.000000 nanos


Benchmarking Redis SET with AOF -- 1000000(iterations)
============================================
Mean ----->  26528.672133 nanos
(90 Percentile) ----->  29686.000000 nanos
(99 Percentile) ----->  41481.000000 nanos
```

The difference between Redis SET operation is because RDB doesn’t impact the performance of the server.
Since the main process only has to fork its process and the child process will take care of all the writing on the disk,
the performance of the parent process is preserved.

Because AOF, Redis will log every write operation received by the
server so that AOF can not be very performant depending on the fsync setting.
