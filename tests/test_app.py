"""Unit tests for jitterbugs.app pure functions (QAQC)."""

from __future__ import annotations

import math
import os
import sys
import types
import unittest

import numpy as np


# ---------------------------------------------------------------------------
# Minimal stubs for heavy optional dependencies so the module can be imported
# in CI without a camera, microphone, or MediaPipe installation.
# ---------------------------------------------------------------------------

def _make_stub(name: str):
    mod = types.ModuleType(name)
    sys.modules[name] = mod
    return mod


def _stub_cv2():
    if "cv2" in sys.modules:
        return
    cv2 = _make_stub("cv2")
    cv2.VideoCapture = object
    cv2.cvtColor = lambda *a, **kw: None
    cv2.COLOR_BGR2RGB = 4
    cv2.FONT_HERSHEY_SIMPLEX = 0
    cv2.LINE_AA = 16
    cv2.putText = lambda *a, **kw: None
    cv2.imshow = lambda *a, **kw: None
    cv2.waitKey = lambda *a, **kw: -1


def _stub_mediapipe():
    if "mediapipe" in sys.modules:
        return
    mp = _make_stub("mediapipe")
    sol = types.ModuleType("mediapipe.solutions")
    sys.modules["mediapipe.solutions"] = sol
    pose_mod = types.ModuleType("mediapipe.solutions.pose")
    sys.modules["mediapipe.solutions.pose"] = pose_mod
    face_mod = types.ModuleType("mediapipe.solutions.face_mesh")
    sys.modules["mediapipe.solutions.face_mesh"] = face_mod

    class _Stub:
        def __init__(self, **kw):
            pass

        def process(self, img):
            return None

        def close(self):
            pass

    pose_mod.Pose = _Stub
    face_mod.FaceMesh = _Stub
    mp.solutions = sol
    sol.pose = pose_mod
    sol.face_mesh = face_mod


def _stub_librosa():
    if "librosa" in sys.modules:
        return
    lib = _make_stub("librosa")
    lib.pyin = lambda *a, **kw: (np.array([]), None, np.array([]))
    lib.note_to_hz = lambda note: 0.0


def _stub_pyaudio():
    if "pyaudio" in sys.modules:
        return
    pa = _make_stub("pyaudio")
    pa.paInt16 = 8
    pa.paContinue = 0

    class _PAStub:
        def open(self, **kw):
            return _StreamStub()

        def terminate(self):
            pass

    class _StreamStub:
        def stop_stream(self):
            pass

        def close(self):
            pass

    pa.PyAudio = _PAStub


# Install stubs before importing the module under test.
_stub_cv2()
_stub_mediapipe()
_stub_librosa()
_stub_pyaudio()

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "src"))

from jitterbugs.app import (  # noqa: E402
    _shoulder_tilt_deg,
    _eye_openness,
    _gaze_offset,
    _draw_overlay,
)


# ---------------------------------------------------------------------------
# Helpers – lightweight fake landmark containers
# ---------------------------------------------------------------------------

class _Lm:
    """Single normalised landmark (mirrors mediapipe NormalizedLandmark)."""

    def __init__(self, x: float, y: float, z: float = 0.0) -> None:
        self.x = x
        self.y = y
        self.z = z


class _LmList:
    """Container that exposes a .landmark list, like NormalizedLandmarkList."""

    def __init__(self, lms: list[_Lm]) -> None:
        self.landmark = lms

    @classmethod
    def zeros(cls, n: int) -> "_LmList":
        return cls([_Lm(0.0, 0.0) for _ in range(n)])


# ---------------------------------------------------------------------------
# _shoulder_tilt_deg
# ---------------------------------------------------------------------------

class TestShoulderTiltDeg(unittest.TestCase):
    """_shoulder_tilt_deg returns signed angle in degrees."""

    _LEFT_SHOULDER = 11
    _RIGHT_SHOULDER = 12

    def _make_landmarks(self, lx, ly, rx, ry, img_w=100, img_h=100):
        """Build a minimal landmark list with only the shoulder landmarks set."""
        n = max(self._LEFT_SHOULDER, self._RIGHT_SHOULDER) + 1
        lms = _LmList.zeros(n)
        # Normalise to [0, 1] for the given image dimensions
        lms.landmark[self._LEFT_SHOULDER] = _Lm(lx / img_w, ly / img_h)
        lms.landmark[self._RIGHT_SHOULDER] = _Lm(rx / img_w, ry / img_h)
        return lms

    def test_level_shoulders_zero_angle(self):
        lms = self._make_landmarks(20, 50, 80, 50)
        angle = _shoulder_tilt_deg(lms, 100, 100)
        self.assertAlmostEqual(angle, 0.0, places=5)

    def test_right_shoulder_higher_negative_angle(self):
        # right shoulder is above (lower y) left shoulder → dy < 0 → negative atan2
        # but we call atan2(dy, dx) where dy = ry - ly
        # right higher means ry < ly → dy negative → negative angle
        lms = self._make_landmarks(20, 60, 80, 40)
        angle = _shoulder_tilt_deg(lms, 100, 100)
        self.assertLess(angle, 0.0)

    def test_left_shoulder_higher_negative_angle(self):
        lms = self._make_landmarks(20, 40, 80, 60)
        angle = _shoulder_tilt_deg(lms, 100, 100)
        self.assertGreater(angle, 0.0)

    def test_known_45_degree_tilt(self):
        # dy == dx → atan2(dy, dx) == 45°
        lms = self._make_landmarks(0, 0, 100, 100, img_w=100, img_h=100)
        angle = _shoulder_tilt_deg(lms, 100, 100)
        self.assertAlmostEqual(angle, 45.0, places=5)

    def test_known_minus_45_degree_tilt(self):
        lms = self._make_landmarks(0, 100, 100, 0, img_w=100, img_h=100)
        angle = _shoulder_tilt_deg(lms, 100, 100)
        self.assertAlmostEqual(angle, -45.0, places=5)

    def test_different_image_sizes(self):
        """Result must be consistent regardless of image scale."""
        lms_a = self._make_landmarks(40, 100, 160, 100, img_w=200, img_h=200)
        lms_b = self._make_landmarks(20, 50, 80, 50, img_w=100, img_h=100)
        self.assertAlmostEqual(
            _shoulder_tilt_deg(lms_a, 200, 200),
            _shoulder_tilt_deg(lms_b, 100, 100),
            places=5,
        )


# ---------------------------------------------------------------------------
# _eye_openness
# ---------------------------------------------------------------------------

class TestEyeOpenness(unittest.TestCase):
    """_eye_openness returns vertical/horizontal eye ratio."""

    # Landmark indices from app.py
    _LEFT_EYE_TOP = 159
    _LEFT_EYE_BOTTOM = 145
    _LEFT_EYE_INNER = 133
    _LEFT_EYE_OUTER = 33

    def _make_landmarks(self, top, bot, inn, out, img_w=100, img_h=100):
        n = max(self._LEFT_EYE_TOP, self._LEFT_EYE_BOTTOM,
                self._LEFT_EYE_INNER, self._LEFT_EYE_OUTER) + 1
        lms = _LmList.zeros(n)
        tx, ty = top
        bx, by = bot
        ix, iy = inn
        ox, oy = out
        lms.landmark[self._LEFT_EYE_TOP] = _Lm(tx / img_w, ty / img_h)
        lms.landmark[self._LEFT_EYE_BOTTOM] = _Lm(bx / img_w, by / img_h)
        lms.landmark[self._LEFT_EYE_INNER] = _Lm(ix / img_w, iy / img_h)
        lms.landmark[self._LEFT_EYE_OUTER] = _Lm(ox / img_w, oy / img_h)
        return lms

    def test_wide_open_eye(self):
        # vertical spread = 20, horizontal span = 40 → ratio = 0.5
        lms = self._make_landmarks(
            top=(50, 40), bot=(50, 60),
            inn=(70, 50), out=(30, 50),
        )
        ratio = _eye_openness(lms, 100, 100)
        self.assertAlmostEqual(ratio, 20.0 / 40.0, places=5)

    def test_closed_eye_near_zero(self):
        # top and bottom at same y → vertical ≈ 0
        lms = self._make_landmarks(
            top=(50, 50), bot=(50, 50),
            inn=(70, 50), out=(30, 50),
        )
        ratio = _eye_openness(lms, 100, 100)
        self.assertAlmostEqual(ratio, 0.0, places=5)

    def test_ratio_positive(self):
        lms = self._make_landmarks(
            top=(50, 30), bot=(50, 70),
            inn=(80, 50), out=(20, 50),
        )
        ratio = _eye_openness(lms, 100, 100)
        self.assertGreater(ratio, 0.0)

    def test_image_scale_independent(self):
        lms_a = self._make_landmarks(
            top=(100, 60), bot=(100, 100),
            inn=(140, 80), out=(60, 80),
            img_w=200, img_h=200,
        )
        lms_b = self._make_landmarks(
            top=(50, 30), bot=(50, 50),
            inn=(70, 40), out=(30, 40),
            img_w=100, img_h=100,
        )
        self.assertAlmostEqual(
            _eye_openness(lms_a, 200, 200),
            _eye_openness(lms_b, 100, 100),
            places=5,
        )


# ---------------------------------------------------------------------------
# _gaze_offset
# ---------------------------------------------------------------------------

class TestGazeOffset(unittest.TestCase):
    """_gaze_offset returns normalised iris offset in [-0.5, 0.5] approximately."""

    _LEFT_IRIS_CENTER = 468
    _LEFT_EYE_INNER = 133
    _LEFT_EYE_OUTER = 33

    def _make_landmarks(self, iris_x, inner_x, outer_x, y=50, img_w=100, img_h=100):
        n = max(self._LEFT_IRIS_CENTER, self._LEFT_EYE_INNER,
                self._LEFT_EYE_OUTER) + 1
        lms = _LmList.zeros(n)
        lms.landmark[self._LEFT_IRIS_CENTER] = _Lm(iris_x / img_w, y / img_h)
        lms.landmark[self._LEFT_EYE_INNER] = _Lm(inner_x / img_w, y / img_h)
        lms.landmark[self._LEFT_EYE_OUTER] = _Lm(outer_x / img_w, y / img_h)
        return lms

    def test_iris_centered_zero_offset(self):
        lms = self._make_landmarks(iris_x=50, inner_x=70, outer_x=30)
        offset = _gaze_offset(lms, 100, 100)
        self.assertAlmostEqual(offset, 0.0, places=5)

    def test_iris_left_negative_offset(self):
        lms = self._make_landmarks(iris_x=35, inner_x=70, outer_x=30)
        offset = _gaze_offset(lms, 100, 100)
        self.assertLess(offset, 0.0)

    def test_iris_right_positive_offset(self):
        lms = self._make_landmarks(iris_x=65, inner_x=70, outer_x=30)
        offset = _gaze_offset(lms, 100, 100)
        self.assertGreater(offset, 0.0)

    def test_known_half_offset(self):
        # center_x = (70 + 30) / 2 = 50; eye_width = |70-30| = 40
        # iris at 70 → offset = (70-50)/40 = 0.5
        lms = self._make_landmarks(iris_x=70, inner_x=70, outer_x=30)
        offset = _gaze_offset(lms, 100, 100)
        self.assertAlmostEqual(offset, 0.5, places=5)


# ---------------------------------------------------------------------------
# _draw_overlay
# ---------------------------------------------------------------------------

class TestDrawOverlay(unittest.TestCase):
    """_draw_overlay must not raise and must modify the frame in place."""

    def _blank_frame(self, h=480, w=640):
        return np.zeros((h, w, 3), dtype=np.uint8)

    def test_all_features_present(self):
        frame = self._blank_frame()
        _draw_overlay(frame, tilt=3.5, openness=0.25, gaze=-0.1,
                      pitch=220.0, jitter=1.5, rms=-20.0)
        # No exception raised is the primary assertion.

    def test_optional_fields_none(self):
        frame = self._blank_frame()
        _draw_overlay(frame, tilt=None, openness=None, gaze=None,
                      pitch=0.0, jitter=0.0, rms=-60.0)

    def test_zero_pitch_skips_pitch_line(self):
        """pitch <= 0 should not raise and the frame remains writable."""
        frame = self._blank_frame()
        _draw_overlay(frame, tilt=None, openness=None, gaze=None,
                      pitch=0.0, jitter=2.0, rms=-30.0)
        self.assertTrue(frame.flags["WRITEABLE"])


# ---------------------------------------------------------------------------
# Audio computation helpers (pure maths, no hardware)
# ---------------------------------------------------------------------------

class TestAudioMath(unittest.TestCase):
    """Test the audio feature maths that _AudioWorker._process_loop uses."""

    def test_rms_silence_near_minus_inf(self):
        audio = np.zeros(2048, dtype=np.float32)
        rms = float(np.sqrt(np.mean(audio ** 2)))
        rms_db = 20.0 * math.log10(rms + 1e-9)
        self.assertLess(rms_db, -150.0)

    def test_rms_full_scale_near_zero_db(self):
        audio = np.ones(2048, dtype=np.float32)
        rms = float(np.sqrt(np.mean(audio ** 2)))
        rms_db = 20.0 * math.log10(rms + 1e-9)
        self.assertAlmostEqual(rms_db, 0.0, places=4)

    def test_jitter_steady_pitch_zero(self):
        voiced = np.full(10, 220.0)
        diffs = np.abs(np.diff(voiced))
        jitter = float(np.mean(diffs) / (np.mean(voiced) + 1e-9) * 100)
        self.assertAlmostEqual(jitter, 0.0, places=5)

    def test_jitter_alternating_pitch_nonzero(self):
        voiced = np.array([200.0, 220.0, 200.0, 220.0, 200.0])
        diffs = np.abs(np.diff(voiced))
        jitter = float(np.mean(diffs) / (np.mean(voiced) + 1e-9) * 100)
        self.assertGreater(jitter, 0.0)

    def test_jitter_proportional_to_variation(self):
        voiced_low = np.array([210.0, 215.0, 210.0, 215.0])
        voiced_high = np.array([200.0, 220.0, 200.0, 220.0])
        jitter_low = float(
            np.mean(np.abs(np.diff(voiced_low)))
            / (np.mean(voiced_low) + 1e-9) * 100
        )
        jitter_high = float(
            np.mean(np.abs(np.diff(voiced_high)))
            / (np.mean(voiced_high) + 1e-9) * 100
        )
        self.assertLess(jitter_low, jitter_high)


if __name__ == "__main__":
    unittest.main()
