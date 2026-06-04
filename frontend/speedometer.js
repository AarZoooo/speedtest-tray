class Speedometer extends HTMLElement {
  constructor() {
    super();
    this.GAUGE_MAX = 1000;
    this.ARC_LENGTH = 251.3; // PI * radius (PI * 80)
    this.currentValue = 0;
  }

  connectedCallback() {
    this.innerHTML = `
            <div class="gauge-container">
                <svg id="gauge" viewBox="0 0 200 180" xmlns="http://www.w3.org/2000/svg">
                    <defs>
                        <linearGradient id="gauge-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                            <stop offset="0%" style="stop-color: var(--accent-start); stop-opacity: 1" />
                            <stop offset="100%" style="stop-color: var(--accent-end); stop-opacity: 1" />
                        </linearGradient>
                    </defs>
                    <!-- 180-degree arc at Y=100 with stroke-width 8 -->
                    <path class="gauge-bg" d="M 20 100 A 80 80 0 0 1 180 100" fill="none" stroke-width="8" />
                    <path id="gauge-fill" class="gauge-fill" d="M 20 100 A 80 80 0 0 1 180 100" fill="none" stroke-width="8" stroke-dasharray="0 ${this.ARC_LENGTH}" />

                    <!-- Labels for the scale -->
                    <text x="20" y="120" text-anchor="middle" class="gauge-label">0</text>
                    <text id="gauge-max-label" x="180" y="120" text-anchor="middle" class="gauge-label">${this.GAUGE_MAX}</text>

                    <!-- Kite-shaped needle with shorter back corner (108 instead of 115) -->
                    <path id="needle" d="M 35 100 L 100 104 L 108 100 L 100 96 Z" class="needle" />

                    <!-- Text values below -->
                    <text id="gauge-value" x="100" y="145" text-anchor="middle" class="gauge-text">0</text>
                    <text x="100" y="165" text-anchor="middle" class="gauge-unit">Mbps</text>
                </svg>
            </div>
        `;

    this.needle = this.querySelector("#needle");
    this.gaugeFill = this.querySelector("#gauge-fill");
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

    this.currentValue = parseFloat(speed) || 0;
    const clampedValue = Math.min(this.currentValue, this.GAUGE_MAX);

    // Rotation around the center (100, 100)
    const angle = (clampedValue / this.GAUGE_MAX) * 180;
    this.needle.style.transform = `rotate(${angle}deg)`;

    const fillLength = (clampedValue / this.GAUGE_MAX) * this.ARC_LENGTH;
    this.gaugeFill.style.strokeDasharray = `${fillLength} ${this.ARC_LENGTH}`;

    this.gaugeValue.textContent = Math.round(this.currentValue);
  }
}

customElements.define("speedometer-gauge", Speedometer);
