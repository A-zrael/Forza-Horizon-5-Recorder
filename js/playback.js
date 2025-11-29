// playback.js - playback loop and scrubber handling

/**
 * Create a playback controller to advance cars and render frames.
 */
export function createPlaybackController({
  playBtn,
  pauseBtn,
  scrubber,
  cars,
  updateCarTrackIndex,
  onFrame
}) {
  let carsRef = cars;
  const state = {playing: false};

  function step() {
    if (!state.playing) return;
    if (!carsRef || !carsRef.length) {
      state.playing = false;
      if (scrubber) scrubber.value = 0;
      onFrame({smoothLens: false});
      return;
    }

    carsRef.forEach(car => {
      if (car.dataIndex < car.data.length - 1) car.dataIndex++;
      updateCarTrackIndex(car);
    });

    const maxIdx = Math.max(...carsRef.map(c => c.dataIndex));
    if (Number.isFinite(maxIdx) && scrubber) scrubber.value = maxIdx;

    onFrame({smoothLens: true});
    requestAnimationFrame(step);
  }

  function play() {
    state.playing = true;
    step();
  }
  function pause() {
    state.playing = false;
  }

  if (playBtn) playBtn.addEventListener("click", play);
  if (pauseBtn) pauseBtn.addEventListener("click", pause);

  if (scrubber) {
    scrubber.addEventListener("input", e => {
      const v = parseInt(e.target.value);
      if (!carsRef || !carsRef.length) {
        return;
      }
      carsRef.forEach(car => {
        car.dataIndex = Math.min(v, car.data.length - 1);
        updateCarTrackIndex(car);
      });
      onFrame({smoothLens: false});
    });
  }

  function setCars(newCars) {
    carsRef = newCars;
  }

  return {play, pause, state, setCars};
}
