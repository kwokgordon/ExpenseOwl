class Transaction extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.innerHTML = `
        <header>
            <!-- <h1 align="center"><img src="../pwa/icon-512.png" alt="ExpenseOwl Logo" height="75" style="vertical-align: middle; margin-right: 10px;">ExpenseOwl</h1> -->
            <div class="nav-bar">
                <a href="/">
                    <img src="/pwa/icon-512.png" alt="ExpenseOwl Logo" height="85" style="vertical-align: middle; margin-right: 20px;">
                </a>
                <a href="/" class="view-button" data-tooltip="Dashboard">
                    <i class="fa-solid fa-chart-pie"></i>
                </a>
                <a href="/table" class="view-button active" data-tooltip="Table View">
                    <i class="fa-solid fa-table"></i>
                </a>
                <a href="/settings" class="view-button" data-tooltip="Settings">
                    <i class="fa-solid fa-gear"></i>
                </a>
            </div>
        </header>
    `;
  }
}

customElements.define('transaction-component', Transaction);