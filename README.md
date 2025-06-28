<p align="center">
<img src="/assets/logo.png" alt="ExpenseOwl Logo" width="250" height="250" /><br>
</p>

<h1 align="center">ExpenseOwl</h1><br>

<p align="center">
<a href="https://github.com/kwokgordon/expenseowl/actions/workflows/release.yml"><img src="https://github.com/kwokgordon/expenseowl/actions/workflows/release.yml/badge.svg" alt="Release"></a>&nbsp;<a href="https://github.com/kwokgordon/expenseowl/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/kwokgordon/expenseowl"></a>&nbsp;<a href="https://hub.docker.com/r/kwokgordon/expenseowl"><img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/kwokgordon/expenseowl"></a>
</p>

<p align="center">
<a href="#why-create-this">Why Create This?</a>&nbsp;&bull;&nbsp;<a href="#features">Features</a>&nbsp;&bull;&nbsp;<a href="#screenshots">Screenshots</a><br><a href="#installation">Installation</a>&nbsp;&bull;&nbsp;<a href="#usage">Usage</a>&nbsp;&bull;&nbsp;<a href="#contributing">Contributing</a>&nbsp;&bull;&nbsp;<a href="#technology-stack">Tech Stack</a>
</p>

<br>

<p align="center">
<b>ExpenseOwl</b> is an extremely simple self-hosted expense tracking system with a modern monthly pie-chart visualization and cashflow showcase.
</p>

<br>

# Why Create This?

There are a ton of amazing projects for expense tracking across GitHub ([Actual](https://github.com/actualbudget/actual), [Firefly III](https://github.com/firefly-iii/firefly-iii), etc.). They're all really incredible! I just find that they aren't the *fastest* or *simplest* to add expenses. They also offer too many features I never use (like varying data formats or complex budgeting). *Don't get me wrong*, they're amazing when complexity is needed, but I wanted something ***dead simple*** that only gives me a quick monthly pie chart and a tabular representation. NOTHING else!

That's why I created this project and I use it in my home lab for my expense tracking. The intention is to track spending across your categories in a simplistic manner. No complicated searching or editing - just `add`, `delete`, and `view`! This intention will not change throughout the project's lifecycle. This is *not* an app for budgeting; it's for straightforward tracking.

# Features

### Core Functionality

- Expense tracking with essential details only (optional name, date, amount, and category)
- Flat file storage system (`data/expenses.json`)
- REST API for expense management
- Single-user focused (mainly for a home lab deployment)
- CSV and JSON export and import of all expense data from the UI
- Custom categories, currency symbol, and start date via app settings
- Beautiful interface that automatically adapts to system for light/dark theme
- UUID-based expense identification in the backend
- Self-contained binary and container image to ensure no internet interaction
- Multi-architecture Docker container with support for persistent storage

### Visualization

1. Main dashboard - category breakdown (pie chart)
    - Click on a category to exclude it from the graph and total; click again to add it back
    - This helps visualize the breakdown without considering some categories like Rent
    - The legend shows categories that make up the the total expenditure of the month
2. Main dashboard - cashflow indicator
    - The default settings have an `Income` category, items in which are not considered expenses
    - If a month has an item in `Income`, ExpenseOwl automatically shows cashflow below the graph
    - Cashflow shows total income, total expenses, and balance (red or green based on +ve or -ve)
3. Table view for detailed expense listing
    - This is where you can view individual expenses chronologically and delete them
    - You can use the browser's search to find a name if needed
4. Month-by-month navigation in both dashboard and table views
5. Settings page for configuring the application
    - Reorder, add, or remove custom categories
    - Select a custom currency to display
    - Select a custom start date to show expenses for a different period
    - Exporting data as CSV or JSON and import data from JSON or CSV

### Progressive Web App (PWA)

The front end of ExpenseOwl can be installed as a Progressive Web App on desktop and mobile devices (i.e., the back end still needs to be self-hosted). To install:

- Desktop: Click the install icon in your browser's address bar
- iOS: Use Safari's "Add to Home Screen" option in the share menu
- Android: Use Chrome's "Install" option in the menu

# Screenshots

Dashboard Showcase:

| | Desktop View | Mobile View |
| --- | --- | --- |
| Dark | <img src="/assets/desktop-dark-main.png" alt="Dashboard Dark" /> | <img src="/assets/mobile-dark-main.png" alt="Mobile Dashboard Dark" /> |
| Light | <img src="/assets/desktop-light-main.png" alt="Dashboard Light" /> | <img src="/assets/mobile-light-main.png" alt="Mobile Dashboard Light" /> |

<details>
<summary>Expand this to see screenshots of other pages</summary>

| | Desktop View | Mobile View |
| --- | --- | --- |
| Table Dark | <img src="/assets/desktop-dark-table.png" alt="Dashboard Dark" /> | <img src="/assets/mobile-dark-table.png" alt="Mobile Dashboard Dark" /> |
| Table Light | <img src="/assets/desktop-light-table.png" alt="Dashboard Light" /> | <img src="/assets/mobile-light-table.png" alt="Mobile Dashboard Light" /> |
| Settings Dark | <img src="/assets/desktop-dark-settings.png" alt="Table Dark" /> | <img src="/assets/mobile-dark-settings.png" alt="Mobile Table Dark" /> |
| Settings Light | <img src="/assets/desktop-light-settings.png" alt="Table Light" /> | <img src="/assets/mobile-light-settings.png" alt="Mobile Table Light" /> |

</details>

# Installation

### Docker Installation (Recommended)

Create a volume or a directory for the project:

```bash
mkdir $HOME/expenseowl
```

```bash
docker run --rm -d \
--name expenseowl \
-p 8080:8080 \
-v $HOME/expenseowl:/app/data \
kwokgordon/expenseowl:main
```

To use it with Docker compose or a container-management system like Portainer or Dockge, use this YAML definition:

```yaml
services:
  budgetlord:
    image: kwokgordon/expenseowl:main
    restart: unless-stopped
    ports:
      - 5006:8080
    volumes:
      - /home/tanq/expenseowl:/app/data # CHANGE DIR
```

### Using the Binary

Download the appropriate binary from the project releases. Running the binary automatically sets up a `data` directory in your CWD. You can visit the frontend at `http://localhost:8080`.

### Building from Source

To directly install the binary from source into your GOBIN, use:

```bash
go install github.com/kwokgordon/expenseowl/cmd/expenseowl@latest
```

Otherwise, to build it yourself:

```bash
git clone https://github.com/kwokgordon/expenseowl.git && \
cd expenseowl && \
go build ./cmd/expenseowl
```

### Kubernetes Deployment

The project also has a community-contributed Kubernetes spec. The spec is a sample and you should review it before deploying in your cluster. Read the [associated readme](./kubernetes/README.md) for more information.

# Usage

Once deployed, use the web interface to do everything. Access it through your browser:

- Dashboard: `http://localhost:8080/`
- Table View: `http://localhost:8080/table`

> [!NOTE]
> This app has no authentication, so deploy carefully. It works very well with a reverse proxy like Nginx Proxy Manager and is mainly intended for homelab use. The app has not undergone a pentest to allow for any production deployment. It should strictly be deployed in a home lab setting, behind authentication, and for only one (or a few non-destructive) user(s).

If command-line automations are required for use with the REST API, read on!

### Executable

The application binary can be run directly within CLI for any common OS and architecture:

```bash
./expenseowl
# or from a custom directory
./expenseowl -data /custom/path
```

### REST API

ExpenseOwl provides an API to allow adding expenses via automations or simply via cURL, Siri Shortcuts, or other automations.

Add Expense:

```bash
curl -X PUT http://localhost:8080/expense \
-H "Content-Type: application/json" \
-d '{
 "name": "Groceries",
 "category": "Food",
 "amount": 75.50,
 "date": "2024-03-15T14:30:00Z"
}'
```

Get All Expenses:

```bash
curl http://localhost:8080/expenses
```

### Config Options

The primary config is stored in the data directory in the `config.json` file. A pre-defined configuration is automatically initialized. The currency in use and the categories can be customized from the `/settings` endpoint within the UI.

ExpenseOwl supports multiple currencies through the CURRENCY environment variable. If not specified, it defaults to USD ($). All available options are shown in the UI settings page.

If setting up for the first time, an environment variable can be used for ease. For example, to use Euro:

```bash
CURRENCY=eur ./expenseowl
```

ExpenseOwl also supports custom categories. A default set is pre-loaded in the config for ease of use and can be easily changed within the UI.

Like currency, if setting up for the first time, categories can be specified in an environment variable like so:

```bash
EXPENSE_CATEGORIES="Rent,Food,Transport,Fun,Bills" ./expenseowl
```

> [!TIP]
> The environment variables can be set in a compose stack or using `-e` in the command line with a Docker command. However, remember that they are only effective in setting up the configuration for first start. Otherwise, use the settings UI.

Similarly, the start date can also be set via the settings UI or the `START_DATE` environment variable.

### Data Import/Export

ExpenseOwl contains a sophisticated method for importing an exporting expenses. The settings page provides the options for exporting all expense data as JSON or CSV. The same page also allows importing data in both JSON and CSV formats.

**Importing CSV**

ExpenseOwl is meant to make things simple, and importing CSV abides by the same philosophy. ExpenseOwl will accept any CSV file as long as it contains the columns - `name`, `category`, `amount`, and `date`. This is case-insensitive so `name` or `Name` doesn't matter.

> [!TIP]
> This feature allows ExpenseOwl to use exported data from any tool as long as the required categories are present, making it insanely easy to shift from any provider.

**Importing JSON**

Primarily, ExpenseOwl maintains a JSON-backend for storing both expenses and config data. If you backed up a Docker volume containing the `config.json` and `expenses.json` files, the recommended way to restore is by mounting the same volume (or directory) to your new container. All data will be instantly usable.

However, in case you need to import JSON formatted data from elsewhere (this is generally rare), you can use the import JSON feature.

> [!WARNING]
> If the time field is not a proper date string (i.e., including time and zone), ExpenseOwl will do a best guess match to set the time to midnight UTC equivalent. This is because time zones are a ... thing.

> [!NOTE]
> ExpenseOwl goes through every row in the imported data, and will intelligently fail on rows that have invalid or absent data. There is a 10 millisecond delay per record to reduce disk overhead, so please allow appropriate time for ingestion (eg. 10 seconds for 1000 records).

# Contributing

Contributions are welcome; please ensure they align with the project's philosophy of maintaining simplicity by strictly using the current [tech stack](#technology-stack). It is intended for home lab use, i.e., a self-hosted first approach (containerized use). Consider the following:

- Additions should have sensible defaults without breaking foundations
- Environment variables can be used for user configuration in containers
- Found a typo or need to ask a question? Please open an issue instead of a PR

# Technology Stack

- Backend: Go
- Storage: JSON file system
- Frontend: Chart.js and vanilla web stack (HTML, JS, CSS)
- Interface: CLI + Web UI
