// track.js - track rendering and lens helpers

/**
 * Convert world coordinates into canvas space using the current bounds.
 */
export function worldToCanvas(trackCanvas, bounds, x, y) {
  const pad = 40;
  const sx = (trackCanvas.width - pad * 2) / (bounds.maxX - bounds.minX || 1);
  const sy = (trackCanvas.height - pad * 2) / (bounds.maxY - bounds.minY || 1);
  const s = Math.min(sx, sy);
  return {
    x: pad + (x - bounds.minX) * s,
    y: trackCanvas.height - (pad + (y - bounds.minY) * s)
  };
}

/**
 * Draw the master track, sector highlights, start/end markers, and lens overlay.
 */
export function drawMasterTrack({
  trackCtx,
  trackCanvas,
  masterTrack,
  raceType,
  sectorCount,
  showTrackSectorsEl,
  minX,
  maxX,
  minY,
  maxY,
  lens,
  cars,
  speedToColor
}) {
  trackCtx.setTransform(1, 0, 0, 1, 0, 0);
  trackCtx.clearRect(0, 0, trackCanvas.width, trackCanvas.height);
  if (!masterTrack.length) return;

  const bounds = {minX, maxX, minY, maxY};
  const toCanvas = (x, y) => worldToCanvas(trackCanvas, bounds, x, y);
  const n = masterTrack.length;
  const totalDist = masterTrack[n - 1].dist || 1;

  // Sector highlighting along the track
  if (showTrackSectorsEl && showTrackSectorsEl.checked && sectorCount > 1) {
    const sectorColors = ['#4C6FFF', '#22C55E', '#FACC15', '#F97373', '#A855F7', '#06B6D4'];
    for (let s = 0; s < sectorCount; s++) {
      trackCtx.save();
      trackCtx.beginPath();
      let started = false;
      for (let i = 1; i < n; i++) {
        const midDist = (masterTrack[i - 1].dist + masterTrack[i].dist) * 0.5;
        let frac = midDist / totalDist;
        if (!isFinite(frac)) frac = 0;
        let si = Math.floor(frac * sectorCount);
        if (si < 0) si = 0;
        if (si >= sectorCount) si = sectorCount - 1;
        if (si !== s) continue;
        const a = toCanvas(masterTrack[i - 1].x, masterTrack[i - 1].y);
        const b = toCanvas(masterTrack[i].x, masterTrack[i].y);
        if (!started) {
          trackCtx.moveTo(a.x, a.y);
          started = true;
        }
        trackCtx.lineTo(b.x, b.y);
      }
      if (started) {
        const base = sectorColors[s % sectorColors.length];
        trackCtx.strokeStyle = base + '55';
        trackCtx.lineWidth = 8;
        trackCtx.stroke();
      }
      trackCtx.restore();
    }
  }

  // Speed heatmap polyline
  for (let i = 1; i < n; i++) {
    const a = toCanvas(masterTrack[i - 1].x, masterTrack[i - 1].y);
    const b = toCanvas(masterTrack[i].x, masterTrack[i].y);
    let color = "#9097A0";

    const v = masterTrack[i].heatSpeed;
    if (v != null && window.__heatRange) {
      let norm = (v - window.__heatMin) / window.__heatRange;
      norm = Math.max(0, Math.min(1, norm));
      color = speedToColor(norm);
    }

    trackCtx.beginPath();
    trackCtx.moveTo(a.x, a.y);
    trackCtx.lineTo(b.x, b.y);
    trackCtx.strokeStyle = color;
    trackCtx.lineWidth = 3;
    trackCtx.globalAlpha = 0.95;
    trackCtx.stroke();
  }
  trackCtx.globalAlpha = 1.0;

  // Start / end markers
  const start = toCanvas(masterTrack[0].x, masterTrack[0].y);
  const end = toCanvas(masterTrack[n - 1].x, masterTrack[n - 1].y);

  if (raceType === 'lapped') {
    trackCtx.beginPath();
    trackCtx.arc(start.x, start.y, 7, 0, Math.PI * 2);
    trackCtx.fillStyle = "#FFFFFF";
    trackCtx.shadowColor = "#FFFFFF";
    trackCtx.shadowBlur = 10;
    trackCtx.fill();
    trackCtx.shadowBlur = 0;
  } else {
    trackCtx.beginPath();
    trackCtx.arc(start.x, start.y, 7, 0, Math.PI * 2);
    trackCtx.fillStyle = "#3EE98A";
    trackCtx.shadowColor = "#3EE98A";
    trackCtx.shadowBlur = 10;
    trackCtx.fill();
    trackCtx.shadowBlur = 0;

    trackCtx.beginPath();
    trackCtx.arc(end.x, end.y, 7, 0, Math.PI * 2);
    trackCtx.fillStyle = "#FF5566";
    trackCtx.shadowColor = "#FF5566";
    trackCtx.shadowBlur = 10;
    trackCtx.fill();
    trackCtx.shadowBlur = 0;
  }

  // Magnifier lens overlay
  if (lens?.active) {
    const zoom = 2.0;

    // Clean the lens area first so unzoomed track does not bleed through
    trackCtx.save();
    trackCtx.beginPath();
    trackCtx.arc(lens.smoothX, lens.smoothY, lens.radius, 0, Math.PI * 2);
    trackCtx.clip();
    trackCtx.clearRect(
      lens.smoothX - lens.radius,
      lens.smoothY - lens.radius,
      lens.radius * 2,
      lens.radius * 2
    );
    trackCtx.restore();

    trackCtx.save(); // save before clipping

    // Clip to circular lens
    trackCtx.beginPath();
    trackCtx.arc(lens.smoothX, lens.smoothY, lens.radius, 0, Math.PI * 2);
    trackCtx.clip();

    // Draw zoomed track
    trackCtx.save();
    trackCtx.translate(lens.smoothX, lens.smoothY);
    trackCtx.scale(zoom, zoom);
    trackCtx.translate(-lens.smoothX, -lens.smoothY);

    for (let i = 1; i < n; i++) {
      const a = toCanvas(masterTrack[i - 1].x, masterTrack[i - 1].y);
      const b = toCanvas(masterTrack[i].x, masterTrack[i].y);

      let color = "#9097A0";
      const v = masterTrack[i].heatSpeed;
      if (v != null && window.__heatRange) {
        let norm = (v - window.__heatMin) / window.__heatRange;
        norm = Math.max(0, Math.min(1, norm));
        color = speedToColor(norm);
      }

      trackCtx.beginPath();
      trackCtx.moveTo(a.x, a.y);
      trackCtx.lineTo(b.x, b.y);
      trackCtx.strokeStyle = color;
      trackCtx.lineWidth = 3;
      trackCtx.stroke();
    }

    trackCtx.restore(); // restore zoom transform only (still clipped)

    // Draw car dots (still inside the clipped region)
    cars.forEach(car => {
      const idx = Math.max(0, Math.min(n - 1, car.index));
      const p = masterTrack[idx];
      const pos = toCanvas(p.x, p.y);

      const dx = pos.x - lens.smoothX;
      const dy = pos.y - lens.smoothY;

      const zx = lens.smoothX + dx * zoom;
      const zy = lens.smoothY + dy * zoom;

      trackCtx.beginPath();
      trackCtx.arc(zx, zy, 6, 0, Math.PI * 2);
      trackCtx.fillStyle = car.color;
      trackCtx.fill();

      trackCtx.beginPath();
      trackCtx.arc(zx, zy, 7, 0, Math.PI * 2);
      trackCtx.strokeStyle = "rgba(255,255,255,0.9)";
      trackCtx.lineWidth = 1;
      trackCtx.stroke();
    });

    trackCtx.restore(); // restore clip

    // Lens border
    trackCtx.beginPath();
    trackCtx.arc(lens.smoothX, lens.smoothY, lens.radius, 0, Math.PI * 2);
    trackCtx.strokeStyle = 'rgba(187,134,255,0.9)';
    trackCtx.lineWidth = 1.5;
    trackCtx.setLineDash([4, 3]);
    trackCtx.stroke();
    trackCtx.setLineDash([]);
  }
}

/**
 * Draw car trails and dots on the master track.
 */
export function drawCars({trackCtx, masterTrack, cars, lens, worldToCanvas}) {
  cars.forEach(car => {
    const idx = Math.max(0, Math.min(masterTrack.length - 1, car.index));
    const p = masterTrack[idx] || masterTrack[0] || {x: 0, y: 0};
    const c = worldToCanvas(p.x, p.y);

    // Skip if inside the lens – zoom overlay draws it
    if (lens?.active) {
      const dxL = c.x - lens.smoothX;
      const dyL = c.y - lens.smoothY;
      const distL2 = dxL * dxL + dyL * dyL;
      if (distL2 <= lens.radius * lens.radius) {
        return;
      }
    }

    // Tail / trail
    trackCtx.beginPath();
    for (let t = Math.max(0, idx - 8); t <= idx; t++) {
      const pt = masterTrack[t] || p;
      const cp = worldToCanvas(pt.x, pt.y);
      if (t === Math.max(0, idx - 8)) trackCtx.moveTo(cp.x, cp.y);
      else trackCtx.lineTo(cp.x, cp.y);
    }
    trackCtx.strokeStyle = car.color;
    trackCtx.lineWidth = 3;
    trackCtx.shadowColor = car.color;
    trackCtx.shadowBlur = 8;
    trackCtx.stroke();
    trackCtx.shadowBlur = 0;

    // Main car dot
    trackCtx.beginPath();
    trackCtx.arc(c.x, c.y, 7, 0, Math.PI * 2);
    trackCtx.fillStyle = car.color;
    trackCtx.shadowColor = car.color;
    trackCtx.shadowBlur = 12;
    trackCtx.fill();
    trackCtx.shadowBlur = 0;
  });
}

/**
 * Resize the track and timeline canvases to fit the viewport.
 */
export function resizeTrackLayout({
  trackCanvas,
  dashColumn,
  eventCanvas,
  setupPanel,
  playbackBar,
  drawMasterTrack,
  drawCars,
  drawEventTimeline
}) {
  const headerH = document.querySelector('header').offsetHeight;
  const setupH = setupPanel.classList.contains('open') ? setupPanel.offsetHeight : 0;
  const playbackH = playbackBar.offsetHeight;
  const margin = 24;
  const available = window.innerHeight - headerH - setupH - playbackH - margin;
  const trackHeight = Math.max(300, available);

  trackCanvas.height = trackHeight;
  trackCanvas.width = trackCanvas.clientWidth || trackCanvas.width;
  dashColumn.style.maxHeight = trackHeight + "px";

  const timelineHeight = Math.max(40, Math.min(80, playbackH - 32));
  eventCanvas.height = timelineHeight;
  eventCanvas.width = eventCanvas.clientWidth || eventCanvas.width;

  drawMasterTrack();
  drawCars();
  drawEventTimeline();
}

/**
 * Wire up lens tracking to the canvas.
 */
export function attachLensHandlers(trackCanvas, lensState) {
  const onMove = (e) => {
    const rect = trackCanvas.getBoundingClientRect();
    lensState.x = e.clientX - rect.left;
    lensState.y = e.clientY - rect.top;
    lensState.active = true;
  };

  const onLeave = () => {
    lensState.active = false;
  };

  trackCanvas.addEventListener("mousemove", onMove);
  trackCanvas.addEventListener("mouseleave", onLeave);

  return () => {
    trackCanvas.removeEventListener("mousemove", onMove);
    trackCanvas.removeEventListener("mouseleave", onLeave);
  };
}
