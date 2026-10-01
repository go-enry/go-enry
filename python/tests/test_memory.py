import os
import gc
import time
import platform
import statistics
import sys
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


def _assert_stable_growth(iterations, print_every, warmup, tail_points, total_limit):
    total_growth, checkpoints, initial, final = _measure_growth(
        iterations, print_every, warmup)
    tail_growth = _tail_median_growth_pct(checkpoints, tail_points)
    if tail_growth >= 1.0 and total_growth < total_limit:
        # A retained allocator arena can produce one step followed by a plateau.
        # Confirm stability over another full segment; keep the original baseline
        # so this cannot reset away cumulative growth from a real leak.
        _, checkpoints, _, final = _measure_growth(iterations, print_every, warmup=0)
        total_growth = (final / initial - 1) * 100
        tail_growth = _tail_median_growth_pct(checkpoints, tail_points)
    assert tail_growth < 1.0, f"Tail median memory growth too large: {tail_growth:.2f}%"
    assert total_growth < total_limit, (
        f"Total memory growth too large: {total_growth:.2f}% "
        f"(initial={initial}, final={final})"
    )


def test_allocator_step_requires_a_stable_confirmation(monkeypatch):
    samples = iter([
        (3.0, [(i, m) for i, m in enumerate([100, 100, 103, 103, 103])], 100, 103),
        (0.0, [(i, 103) for i in range(5)], 103, 103),
    ])
    monkeypatch.setattr(sys.modules[__name__], "_measure_growth",
                        lambda *args, **kwargs: next(samples))
    _assert_stable_growth(100, 20, 0, 5, 5.0)


def test_continuing_growth_fails_confirmation(monkeypatch):
    samples = iter([
        (3.0, [(i, m) for i, m in enumerate([100, 100, 103, 103, 103])], 100, 103),
        (3.0, [(i, m) for i, m in enumerate([103, 104, 105, 106, 107])], 103, 107),
    ])
    monkeypatch.setattr(sys.modules[__name__], "_measure_growth",
                        lambda *args, **kwargs: next(samples))
    with pytest.raises(AssertionError):
        _assert_stable_growth(100, 20, 0, 5, 5.0)


def test_cumulative_growth_cannot_reset_at_confirmation(monkeypatch):
    samples = iter([
        (3.0, [(i, m) for i, m in enumerate([100, 100, 103, 103, 103])], 100, 103),
        (3.0, [(i, 106) for i in range(5)], 103, 106),
    ])
    monkeypatch.setattr(sys.modules[__name__], "_measure_growth",
                        lambda *args, **kwargs: next(samples))
    with pytest.raises(AssertionError, match="Total memory growth"):
        _assert_stable_growth(100, 20, 0, 5, 5.0)


def test_no_memory_leak_short():
    """
    Fast regression test (runs in the normal CI matrix).

    Designed to detect *continued* growth, not one-off allocator arena steps.
    """
    iterations = int(os.getenv("ENRY_MEM_ITERS_SHORT", "2000"))
    print(f"running short memory test with {iterations} iterations")

    _assert_stable_growth(
        iterations=iterations,
        print_every=max(1, iterations // 10),
        warmup=int(os.getenv("ENRY_MEM_WARMUP_SHORT", "200")),
        tail_points=5,
        total_limit=10.0,
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

    _assert_stable_growth(
        iterations=iterations,
        print_every=max(1, iterations // 20),
        warmup=int(os.getenv("ENRY_MEM_WARMUP_LONG", "500")),
        tail_points=7,
        total_limit=5.0,
    )
