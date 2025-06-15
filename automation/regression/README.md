# Regression Testing for Loan Service

## Pengantar (Introduction)

Repositori ini berisi tes regresi otomatis untuk Loan Service menggunakan Cypress. Tes-tes ini memastikan bahwa semua fungsionalitas utama aplikasi bekerja dengan baik setelah perubahan kode.

## Overview

This repository contains automated regression tests for the Loan Service application using Cypress. These tests ensure that all core functionalities of the application continue to work correctly after code changes.

The regression test suite covers the complete loan lifecycle, including:
- Creating a new loan
- Validating duplicate borrower prevention
- Approving loans
- Adding investments to loans
- Disbursing approved loans

## Setup

### Prerequisites

- Node.js (latest LTS version recommended)
- npm or yarn
- The Loan Service application running locally on port 7002

### Installation

```bash
# Install dependencies
npm install
```

## Cara Menjalankan Tes (Running Tests)

Untuk menjalankan tes regresi, pastikan aplikasi Loan Service berjalan di `http://localhost:7002`, kemudian jalankan perintah berikut:

```bash
npx cypress run
```

Untuk menjalankan tes dengan UI Cypress:

```bash
npx cypress open
```

## Running Tests

To run the regression tests, ensure the Loan Service application is running at `http://localhost:7002`, then execute:

```bash
# Run tests in headless mode
npx cypress run

# Run tests with Cypress UI
npx cypress open
```

## Test Description

The test suite (`loan-lifecycle.cy.js`) verifies the complete loan process:

1. **Loan Creation**: Tests the API endpoint for creating a new loan with borrower details, principal amount, rate, and ROI.

2. **Duplicate Validation**: Ensures the system prevents creating loans with duplicate borrower IDs.

3. **Loan Approval**: Verifies the loan approval process with validator information and proof documentation.

4. **Investment Addition**: Tests adding investor funds to an approved loan.

5. **Loan Disbursement**: Confirms the final step of disbursing funds to the borrower after all requirements are met.

## Troubleshooting

If you encounter issues:

1. Verify the Loan Service is running on port 7002
2. Check network connectivity between Cypress and the application
3. Review the Cypress logs for detailed error information

## Pemecahan Masalah (Troubleshooting)

Jika mengalami masalah:

1. Pastikan Loan Service berjalan di port 7002
2. Periksa konektivitas jaringan antara Cypress dan aplikasi
3. Tinjau log Cypress untuk informasi error yang lebih detail