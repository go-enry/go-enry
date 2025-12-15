import os
import gc
import time
import os
import psutil
import enry
import pytest


def _measure_growth(iterations: int, print_every: int = 0) -> float:
    process = psutil.Process(os.getpid())

    # warm up the enry library to trigger one-time init allocations
    for _ in range(10):
        enry.get_language("test.py", b"import os\n")

    initial = process.memory_info().rss
    start = time.time()

    for i in range(1, iterations + 1):
        enry.get_language("test.py", b"import os\nprint('Hello')")
        if print_every and i % print_every == 0:
            gc.collect()
            current = process.memory_info().rss
            growth = (current / initial - 1) * 100
            elapsed = time.time() - start
            print(f"Iter {i:<6} RSS: {current}  Growth: {growth:.2f}%  Elapsed: {elapsed:.1f}s")

    final = process.memory_info().rss
    final_growth = (final / initial - 1) * 100
    return final_growth


def test_no_memory_leak_short():
    """Quick test. Iteration count configurable via `ENRY_MEM_ITERS_SHORT`.

    Default is reduced for fast test runs (2000). Set the env var to
    a larger number for thorough testing (e.g. 20000).
    """
    iterations = int(os.getenv("ENRY_MEM_ITERS_SHORT", "2000"))
    print(f"running short memory test with {iterations} iterations")
    growth = _measure_growth(iterations=iterations, print_every=max(1, iterations // 10))
    assert growth < 1.0, f"Memory growth too large: {growth:.2f}%"


@pytest.mark.skipif(os.getenv("ENRY_MEM_RUN_SMOKE", "1") != "1",
                    reason="smoke test disabled via ENRY_MEM_RUN_SMOKE")
def test_no_memory_leak_smoke():
    """Smoke test. Iterations configurable via `ENRY_MEM_ITERS_SMOKE`.

    Default is reduced for CI/dev convenience (200).
    Set `ENRY_MEM_RUN_SMOKE=0` to skip this test, or override iterations
    with `ENRY_MEM_ITERS_SMOKE`.
    """
    iterations = int(os.getenv("ENRY_MEM_ITERS_SMOKE", "200"))
    growth = _measure_growth(iterations=iterations, print_every=0)
    assert growth < 2.0, f"Memory growth too large: {growth:.2f}%"
