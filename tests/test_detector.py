"""Tests for the Detector class."""

import pytest
from jitterbugs.detector import Detector, DetectionResult, DEFAULT_PATTERNS


class TestDetector:
    def test_default_patterns_present(self):
        d = Detector()
        assert len(d._compiled) == len(DEFAULT_PATTERNS)

    def test_scan_returns_hits_for_flagged_text(self):
        d = Detector()
        results = d.scan("this is spam content")
        assert results
        assert all(isinstance(r, DetectionResult) for r in results)

    def test_scan_returns_empty_for_clean_text(self):
        d = Detector()
        results = d.scan("hello, how are you today?")
        assert results == []

    def test_is_flagged_true(self):
        d = Detector()
        assert d.is_flagged("buy cheap spam")

    def test_is_flagged_false(self):
        d = Detector()
        assert not d.is_flagged("completely ordinary sentence")

    def test_scan_case_insensitive(self):
        d = Detector()
        assert d.is_flagged("SPAM everywhere")
        assert d.is_flagged("Spam here")

    def test_add_pattern(self):
        d = Detector()
        before = len(d._compiled)
        d.add_pattern(r"\bfoo\b")
        assert len(d._compiled) == before + 1
        assert d.is_flagged("foo bar")

    def test_custom_patterns_only(self):
        d = Detector(patterns=[r"\btest\b"])
        assert d.is_flagged("this is a test")
        assert not d.is_flagged("spam scam abuse")

    def test_detection_result_fields(self):
        d = Detector(patterns=[r"\bspam\b"])
        results = d.scan("hello spam world")
        assert len(results) == 1
        r = results[0]
        assert r.match.lower() == "spam"
        assert r.position == 6

    def test_detection_result_repr(self):
        r = DetectionResult(r"\bspam\b", "spam", 0)
        assert "spam" in repr(r)
