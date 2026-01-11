import { APP_VERSION } from './version.js';

class Header extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    // Determine active page
    const path = window.location.pathname;
    let dashboardActive = '';
    let tableActive = '';
    let settingsActive = '';
    if (path === '/' || path.startsWith('/dashboard')) {
      dashboardActive = 'active';
    } else if (path.startsWith('/table')) {
      tableActive = 'active';
    } else if (path.startsWith('/settings')) {
      settingsActive = 'active';
    }
    this.innerHTML = `
        <header>
            <div class="nav-bar">
                <a href="/">
                    <img src="/pwa/icon-512.png" alt="ExpenseOwl Logo" height="85" style="vertical-align: middle; margin-right: 20px;">
                </a>
                <a href="/" class="view-button ${dashboardActive}" data-tooltip="Dashboard">
                    <i class="fa-solid fa-chart-pie"></i>
                </a>
                <a href="/table" class="view-button ${tableActive}" data-tooltip="Table View">
                    <i class="fa-solid fa-table"></i>
                </a>
                <a href="/settings" class="view-button ${settingsActive}" data-tooltip="Settings">
                    <i class="fa-solid fa-gear"></i>
                </a>
                <span style="margin-left:auto;font-size:0.95em;color:#888;align-self:center;padding-left:1.5rem;">v${APP_VERSION}</span>
            </div>
        </header>
    `;
  }
}

customElements.define('header-component', Header);