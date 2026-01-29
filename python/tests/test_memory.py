import os
import gc
import time
import platform
import statistics
import psutil
import enry
import pytest


def _mem_bytes(process: psutil.Process) -> int:
    """
    Return a memory metric suitable for leak detection.

    On Linux, prefer USS (unique set size) because RSS can jump in steps due to
    allocator arenas/tcache and may not return to the OS even after frees.
    Elsewhere, fall back to RSS.
    """
    try:
        if platform.system() == "Linux":
            return process.memory_full_info().uss
    except Exception:
        pass
    return process.memory_info().rss


def _stabilize(process: psutil.Process) -> int:
    """Force GC and take a stable memory sample."""
    gc.collect()
    time.sleep(0.01)  # helps stabilize readings on CI runners
    return _mem_bytes(process)


def _measure_growth(iterations: int, print_every: int = 0, warmup: int = 200):
    """
    Returns:
      total_growth_pct: final vs initial percent change
      checkpoints: list[(iter, mem_bytes)]
      initial: baseline mem_bytes
      final: final mem_bytes
    """
    process = psutil.Process(os.getpid())

    # Warm up to trigger one-time init allocations/caches (Go runtime, cgo, CFFI, etc.)
    for _ in range(warmup):
        enry.get_language("test.py", b"import os\n")

    initial = _stabilize(process)
    start = time.time()

    checkpoints = []
    for i in range(1, iterations + 1):
        enry.get_language("test.py", b"import os\nprint('Hello')")

        if print_every and i % print_every == 0:
            current = _stabilize(process)
            checkpoints.append((i, current))
            growth = (current / initial - 1) * 100
            elapsed = time.time() - start
            print(f"Iter {i:<6} MEM: {current}  Growth: {growth:.2f}%  Elapsed: {elapsed:.1f}s")

    final = _stabilize(process)
    total_growth_pct = (final / initial - 1) * 100
    return total_growth_pct, checkpoints, initial, final


def _tail_median_growth_pct(checkpoints, tail_points: int = 5) -> float:
    """
    Compute growth over the tail of the run using the median of the last N
    checkpoints. This is robust to a single allocator "step" at the very end.
    """
    if len(checkpoints) < 2:
        return 0.0

    tail = checkpoints[-tail_points:] if len(checkpoints) >= tail_points else checkpoints
    mems = [m for _, m in tail]
    start = mems[0]
    med = int(statistics.median(mems))
    if start == 0:
        return 0.0
    return (med / start - 1) * 100


def test_no_memory_leak_short():
    """
    Fast regression test (runs in the normal CI matrix).

    Designed to detect *continued* growth, not one-off allocator arena steps.
    """
    iterations = int(os.getenv("ENRY_MEM_ITERS_SHORT", "2000"))
    print(f"running short memory test with {iterations} iterations")

    total_growth, checkpoints, initial, final = _measure_growth(
        iterations=iterations,
        print_every=max(1, iterations // 10),
        warmup=int(os.getenv("ENRY_MEM_WARMUP_SHORT", "200")),
    )

    # Primary: tail of run should be essentially flat (robust to one final-step jump).
    tail_growth = _tail_median_growth_pct(checkpoints, tail_points=5)
    assert tail_growth < 1.0, f"Tail median memory growth too large: {tail_growth:.2f}%"

    # Secondary: guard against runaway growth (should never happen in 2k calls).
    # Allow allocator steps; this is just a safety net.
    assert total_growth < 10.0, (
        f"Total memory growth too large: {total_growth:.2f}% (initial={initial}, final={final})"
    )


@pytest.mark.skipif(
    os.getenv("ENRY_MEM_RUN_SMOKE", "1") != "1",
    reason="smoke test disabled via ENRY_MEM_RUN_SMOKE",
)
def test_no_memory_leak_smoke():
    """
    Very small smoke check for CI/dev. Kept loose; its job is to catch egregious leaks.
    """
    iterations = int(os.getenv("ENRY_MEM_ITERS_SMOKE", "200"))
    total_growth, _, _, _ = _measure_growth(
        iterations=iterations,
        print_every=0,
        warmup=int(os.getenv("ENRY_MEM_WARMUP_SMOKE", "50")),
    )
    assert total_growth < 5.0, f"Memory growth too large: {total_growth:.2f}%"


@pytest.mark.skipif(
    os.getenv("ENRY_MEM_RUN_LONG", "0") != "1",
    reason="long memory test disabled by default (set ENRY_MEM_RUN_LONG=1)",
)
def test_no_memory_leak_long():
    """
    Long regression test intended to catch true C-level leaks reliably.

    Enable in CI for a single job only (e.g. ubuntu + one python version).
    """
    iterations = int(os.getenv("ENRY_MEM_ITERS_LONG", "100000"))
    print(f"running long memory test with {iterations} iterations")

    total_growth, checkpoints, initial, final = _measure_growth(
        iterations=iterations,
        print_every=max(1, iterations // 20),  # ~5% checkpoints
        warmup=int(os.getenv("ENRY_MEM_WARMUP_LONG", "500")),
    )

    # Long test can use stricter "end segment" notion by just comparing early/late medians.
    tail_growth = _tail_median_growth_pct(checkpoints, tail_points=7)
    assert tail_growth < 1.0, f"Tail median memory growth too large: {tail_growth:.2f}%"

    # Allow some total allocator growth, especially on Linux, but keep it bounded.
    assert total_growth < 5.0, (
        f"Total memory growth too large: {total_growth:.2f}% "
        f"(initial={initial}, final={final})"
    )
