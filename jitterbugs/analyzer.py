"""
Analyzer module for jitterbugs.

Aggregates detection results over a session and produces a summary report.
"""

from __future__ import annotations

from collections import Counter
from dataclasses import dataclass, field
from typing import Sequence

from .detector import Detector, DetectionResult


@dataclass
class SessionReport:
    """Summary of signals found during an analysis session."""

    total_items: int = 0
    flagged_items: int = 0
    hit_counts: Counter = field(default_factory=Counter)

    # ------------------------------------------------------------------

    @property
    def flagged_ratio(self) -> float:
        """Fraction of items that contained at least one signal."""
        if self.total_items == 0:
            return 0.0
        return self.flagged_items / self.total_items

    def top_patterns(self, n: int = 5) -> list[tuple[str, int]]:
        """Return the *n* most-triggered patterns and their counts."""
        return self.hit_counts.most_common(n)

    def __repr__(self) -> str:
        return (
            f"SessionReport(total={self.total_items}, "
            f"flagged={self.flagged_items}, "
            f"ratio={self.flagged_ratio:.2%})"
        )


class Analyzer:
    """
    Runs a :class:`~jitterbugs.Detector` over a collection of text items
    and accumulates a :class:`SessionReport`.

    Usage::

        a = Analyzer()
        report = a.analyze(["hello world", "buy cheap spam here"])
        print(report)
        print(report.top_patterns())
    """

    def __init__(self, detector: Detector | None = None) -> None:
        self._detector = detector if detector is not None else Detector()

    # ------------------------------------------------------------------
    # Public API
    # ------------------------------------------------------------------

    def analyze(self, items: Sequence[str]) -> SessionReport:
        """Scan every item and return a consolidated :class:`SessionReport`."""
        report = SessionReport()
        for text in items:
            report.total_items += 1
            hits: list[DetectionResult] = self._detector.scan(text)
            if hits:
                report.flagged_items += 1
                for hit in hits:
                    report.hit_counts[hit.pattern] += 1
        return report

    def analyze_one(self, text: str) -> list[DetectionResult]:
        """Convenience wrapper — scan a single item and return its hits."""
        return self._detector.scan(text)
