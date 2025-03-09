# DHIS2 Validation Results Cleaner 🛠️

This Go script connects to a **DHIS2** instance, fetches validation results, and deletes them automatically. It prompts for **username and password** securely (hides password input) and processes validation results page by page.

## 🚀 Features
- Fetches and deletes validation results from DHIS2.
- **Secure authentication**: Prompts for **username and password** (password input is hidden).
- Uses **Base64 encoding** for authentication.
- **Handles pagination**: Fetches validation results **only from the first page**.
- **Error handling**: Reports failed requests and deletions.

## 📌 Prerequisites
- Install **Go 1.18+**
- DHIS2 instance access with API permissions.

## 📥 Installation
Clone this repository:
```sh
git clone https://github.com/didate/dhis2-validation-cleaner.git
cd dhis2-validation-cleaner
```

## 🛠️ Usage
Run the script:

```sh
go run main.go
```

## 🔐 Authentication
You will be prompted for:
- DHIS2 Base URL (e.g., https://play.dhis2.org/dev)
- Username
- Password (input is hidden for security)

## 🏃 Example Execution

```
Enter DHIS2 base URL: https://play.dhis2.org/dev
Enter your DHIS2 username: admin
Enter your DHIS2 password:
ℹ️ Fetching: https://play.dhis2.org/dev/api/validationResults?page=1&pageSize=50
✅ No more validation results found. Stopping.
```