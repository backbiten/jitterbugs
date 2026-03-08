"""Tests for the Analyzer class."""

import pytest
from jitterbugs.analyzer import Analyzer, SessionReport
from jitterbugs.detector import Detector


class TestSessionReport:
    def test_flagged_ratio_empty(self):
        r = SessionReport()
        assert r.flagged_ratio == 0.0

    def test_flagged_ratio_calculation(self):
        r = SessionReport(total_items=10, flagged_items=3)
        assert abs(r.flagged_ratio - 0.3) < 1e-9

    def test_top_patterns_empty(self):
        r = SessionReport()
        assert r.top_patterns() == []

    def test_repr(self):
        r = SessionReport(total_items=4, flagged_items=2)
        assert "50.00%" in repr(r)


class TestAnalyzer:
    def test_analyze_all_clean(self):
        a = Analyzer()
        report = a.analyze(["hello world", "good morning"])
        assert report.total_items == 2
        assert report.flagged_items == 0

    def test_analyze_some_flagged(self):
        a = Analyzer()
        items = ["totally fine", "buy cheap spam now", "scam alert"]
        report = a.analyze(items)
        assert report.total_items == 3
        assert report.flagged_items == 2

    def test_analyze_empty_list(self):
        a = Analyzer()
        report = a.analyze([])
        assert report.total_items == 0
        assert report.flagged_items == 0
        assert report.flagged_ratio == 0.0

    def test_analyze_one(self):
        a = Analyzer()
        hits = a.analyze_one("this is spam")
        assert len(hits) > 0

    def test_analyze_one_clean(self):
        a = Analyzer()
        hits = a.analyze_one("nothing suspicious")
        assert hits == []

    def test_custom_detector_injection(self):
        custom = Detector(patterns=[r"\bcustom\b"])
        a = Analyzer(detector=custom)
        report = a.analyze(["normal text", "custom trigger"])
        assert report.flagged_items == 1

    def test_top_patterns_reflects_counts(self):
        a = Analyzer()
        items = ["spam here", "more spam", "clean text"]
        report = a.analyze(items)
        top = report.top_patterns(1)
        assert top
        pattern, count = top[0]
        assert count >= 2
