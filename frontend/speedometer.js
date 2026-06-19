class Speedometer extends HTMLElement {
  constructor() {
    super();
    this.GAUGE_MAX = 1000;
    this.ARC_LENGTH = 125.6; // PI * 40
    this.currentValue = 0;
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
}

customElements.define("speedometer-gauge", Speedometer);
