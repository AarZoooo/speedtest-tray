class Speedometer extends HTMLElement {
  constructor() {
    super();
    this.GAUGE_MAX = 1000;
    this.ARC_LENGTH = 125.6; // PI * 40
    this.currentValue = 0;
    this._sweepTimeout = null;
    this.sweeping = false;
  }

  connectedCallback() {
    this.innerHTML = `
            <div class="gauge-container">
                <svg id="gauge" viewBox="0 0 200 180" xmlns="http://www.w3.org/2000/svg">
                    <defs>
                        <linearGradient id="gauge-gradient" gradientUnits="userSpaceOnUse" x1="20" y1="100" x2="180" y2="100">
                            <stop offset="0%" style="stop-color: var(--accent-end); stop-opacity: 1" />
                            <stop offset="100%" style="stop-color: var(--accent-start); stop-opacity: 1" />
                        </linearGradient>
                    </defs>
                    <!-- Background Semicircle -->
                    <path class="gauge-bg" d="M 60 100 A 40 40 0 0 1 140 100" fill="none" stroke-width="80" />

                    <!-- Bloom Layer: Same gradient, blurred via CSS -->
                    <path id="gauge-bloom" class="gauge-bloom" d="M 60 100 A 40 40 0 0 1 140 100" fill="none" stroke-width="80" stroke-dasharray="0 ${this.ARC_LENGTH}" />

                    <!-- Main Fill Layer -->
                    <path id="gauge-fill" class="gauge-fill" d="M 60 100 A 40 40 0 0 1 140 100" fill="none" stroke-width="80" stroke-dasharray="0 ${this.ARC_LENGTH}" />

                    <!-- Labels positioned closer to the base -->
                    <text x="20" y="115" text-anchor="middle" class="gauge-label">0</text>
                    <text id="gauge-max-label" x="180" y="115" text-anchor="middle" class="gauge-label">${this.GAUGE_MAX}</text>

                    <!-- Balanced needle -->
                    <path id="needle" d="M 35 100 L 100 103 L 106 100 L 100 97 Z" class="needle" />

                    <!-- Text values below -->
                    <text id="gauge-value" x="100" y="145" text-anchor="middle" class="gauge-text">0</text>
                    <text x="100" y="165" text-anchor="middle" class="gauge-unit">Mbps</text>
                </svg>
            </div>
        `;

    this.needle = this.querySelector("#needle");
    this.gaugeFill = this.querySelector("#gauge-fill");
    this.gaugeBloom = this.querySelector("#gauge-bloom");
    this.gaugeValue = this.querySelector("#gauge-value");
    this.gaugeMaxLabel = this.querySelector("#gauge-max-label");
  }

  setMax(max) {
    this.GAUGE_MAX = parseFloat(max) || 1000;
    if (this.gaugeMaxLabel) {
      this.gaugeMaxLabel.textContent = this.GAUGE_MAX;
    }
    this.setValue(this.currentValue);
  }

  setValue(speed) {
    if (!this.gaugeValue) return;

    const val = parseFloat(speed) || 0;

    if (val === 0 && this.currentValue > 0) {
      this.needle.classList.add("easing");
      this.gaugeFill.classList.add("easing");
      if (this.gaugeBloom) {
        this.gaugeBloom.classList.add("easing");
      }
    } else if (val > 0) {
      this.needle.classList.remove("easing");
      this.gaugeFill.classList.remove("easing");
      if (this.gaugeBloom) {
        this.gaugeBloom.classList.remove("easing");
      }
    }

    this.currentValue = val;
    const clampedValue = Math.min(this.currentValue, this.GAUGE_MAX);

    const angle = (clampedValue / this.GAUGE_MAX) * 180;
    this.needle.style.transform = `rotate(${angle}deg)`;

    const fillLength = (clampedValue / this.GAUGE_MAX) * this.ARC_LENGTH;
    this.gaugeFill.style.strokeDasharray = `${fillLength} ${this.ARC_LENGTH}`;
    if (this.gaugeBloom) {
      this.gaugeBloom.style.strokeDasharray = `${fillLength} ${this.ARC_LENGTH}`;
    }

    this.gaugeValue.textContent = Math.round(this.currentValue);
  }

  // Startup sweep: needle swings to max, holds, swings back — once.
  // Returns a Promise that resolves when the full animation completes.
  playStartupSweep() {
    if (this.sweeping) return Promise.resolve();
    this.sweeping = true;
    // Engage sweep transitions
    this.needle?.classList.add("sweeping");
    this.gaugeFill?.classList.add("sweeping");
    this.gaugeBloom?.classList.add("sweeping");
    return new Promise((resolve) => this._runSweepCycle(resolve));
  }

  _runSweepCycle(resolve) {
    if (!this.sweeping) {
      resolve?.();
      return;
    }

    // Sweep to max smoothly (500ms transition)
    this._setSweepAngle(180);
    // Hold at max for 500ms, then sweep back down
    this._sweepTimeout = setTimeout(() => {
      if (!this.sweeping) {
        resolve?.();
        return;
      }
      this._setSweepAngle(0);
      // After sweep-down (500ms), mark as done and resolve
      this._sweepTimeout = setTimeout(() => {
        this.sweeping = false;
        this._sweepTimeout = null;
        this.needle?.classList.remove("sweeping");
        this.gaugeFill?.classList.remove("sweeping");
        this.gaugeBloom?.classList.remove("sweeping");
        resolve?.();
      }, 550);
    }, 1050);
  }

  _setSweepAngle(angleDeg) {
    if (!this.needle) return;
    this.needle.style.transform = `rotate(${angleDeg}deg)`;
    const fillLength = (angleDeg / 180) * this.ARC_LENGTH;
    if (this.gaugeFill) this.gaugeFill.style.strokeDasharray = `${fillLength} ${this.ARC_LENGTH}`;
    if (this.gaugeBloom) this.gaugeBloom.style.strokeDasharray = `${fillLength} ${this.ARC_LENGTH}`;
  }

  stopSweep() {
    this.sweeping = false;
    if (this._sweepTimeout !== null) {
      clearTimeout(this._sweepTimeout);
      this._sweepTimeout = null;
    }
    // Remove sweep transitions so normal setValue() takes over cleanly
    this.needle?.classList.remove("sweeping");
    this.gaugeFill?.classList.remove("sweeping");
    this.gaugeBloom?.classList.remove("sweeping");
  }
}

customElements.define("speedometer-gauge", Speedometer);
