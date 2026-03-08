"""
Detector module for jitterbugs.

Identifies patterns of problematic behaviour in text or event streams.
"""

import re
from typing import Iterable


# Default set of signal patterns used during detection.
DEFAULT_PATTERNS: list[str] = [
    r"\b(spam|scam|phish)\b",
    r"\b(abuse|harass|threaten)\w*\b",
    r"(https?://\S+\.ru\b)",          # example domain-based signal
]


class DetectionResult:
    """Container for a single detection hit."""

    def __init__(self, pattern: str, match: str, position: int) -> None:
        self.pattern = pattern
        self.match = match
        self.position = position

    def __repr__(self) -> str:
        return (
            f"DetectionResult(pattern={self.pattern!r}, "
            f"match={self.match!r}, position={self.position})"
        )


class Detector:
    """
    Scans input text for signals that indicate problematic content.

    Usage::

        d = Detector()
        results = d.scan("Buy cheap meds now! spam link http://example.ru/x")
        for r in results:
            print(r)
    """

    def __init__(self, patterns: Iterable[str] | None = None) -> None:
        raw_patterns = list(patterns) if patterns is not None else DEFAULT_PATTERNS
        self._compiled = [(p, re.compile(p, re.IGNORECASE)) for p in raw_patterns]

    # ------------------------------------------------------------------
    # Public API
    # ------------------------------------------------------------------

    def scan(self, text: str) -> list[DetectionResult]:
        """Return all detection hits found in *text*."""
        hits: list[DetectionResult] = []
        for raw, regex in self._compiled:
            for m in regex.finditer(text):
                hits.append(DetectionResult(raw, m.group(), m.start()))
        return hits

    def is_flagged(self, text: str) -> bool:
        """Return ``True`` if *text* contains at least one signal."""
        return bool(self.scan(text))

    def add_pattern(self, pattern: str) -> None:
        """Append a new regex *pattern* to the detector at runtime."""
        self._compiled.append((pattern, re.compile(pattern, re.IGNORECASE)))
