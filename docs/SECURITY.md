# Security Policy

## Supported Versions

The following versions of this project are currently supported with security updates:

| Version | Supported | Notes |
| ------- | --------- | ----- |
| 1.0.14   | ✅         | Latest - includes retry logic & circuit breaker with complete exports |
| 1.0.13   | ⚠️         | Missing JavaScript exports - upgrade to 1.0.14 |
| 1.0.12   | ⚠️         | Missing TypeScript definitions - upgrade to 1.0.14 |
| 1.0.11   | ⚠️         | No retry/circuit breaker features |
| 1.0.10   | ❌         | No longer supported |
| < 1.0.10 | ❌         | No longer supported |

> **Note:** Only versions 1.0.11 and above receive security fixes. Users are strongly encouraged to upgrade to version 1.0.14 for the latest features and complete API exports.

> **Feature Availability:** 
> - Retry logic and circuit breaker features are only available in version 1.0.12+
> - Complete TypeScript definitions for `setRetryOptions` are only available in version 1.0.13+
> - Complete JavaScript exports for `setRetryOptions` are only available in version 1.0.14+

---

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly.

### How to Report

* **GitHub:** Open a **private security advisory** via GitHub

Please include as much information as possible:

* A clear description of the issue
* Steps to reproduce (proof-of-concept if available)
* Affected versions
* Potential impact (e.g. DoS, data leak, RCE)

### What to Expect

* **Acknowledgement:** Within **48 hours** of receiving your report
* **Initial Assessment:** Within **5 business days**
* **Fix & Disclosure:** Timing depends on severity and complexity

If the vulnerability is accepted:

* We will work on a fix as quickly as possible
* A patched release will be published
* You will be credited in the security advisory (unless you prefer to remain anonymous)

If the vulnerability is declined:

* We will provide a clear explanation for the decision

### Scope

The following are **in scope**:

* Security issues in the project source code
* Vulnerabilities affecting consumers of the library

The following are **out of scope**:

* Issues in third-party dependencies (please report them upstream)
* Denial of service via unrealistic or extreme misuse
* Social engineering attacks

---

## Security Updates

Security fixes are released as soon as possible and documented in the release notes. No backports are provided for unsupported versions.

---

Thank you
