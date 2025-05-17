# `lib` — AWS Profiles Library

This package provides a programmatic interface for managing custom AWS profiles stored outside of `~/.aws`. It is designed to power the `aws-profiles v2` project, enabling secure and flexible profile management with first-class support for native encryption and CLI/TUI integration.

---

## 🔧 Features

- **Add, edit, remove** AWS profiles
- Store profiles in a **separate dotfile** (not `~/.aws`)
- Support the following fields per profile:
  - `name`
  - `region`
  - `role_arn`
  - `access_key_id`
  - `secret_access_key`

---

## 🔐 Security

- Profiles can be **optionally encrypted**:
  - Uses **system-native credential storage** (e.g., macOS Keychain, Windows Credential Manager) when available
  - Falls back to **password-based encryption** if system keyring is not supported

---

## 🔁 AWS Migration Support

- One-time import of profiles from `~/.aws/config` and `~/.aws/credentials`
    - A backup will be made and will require manual deletion.
- Imported profiles are moved into the new config file
- Leaves behind an **empty default profile** in `.aws`

---

## ✅ Profile Activation

- Set a profile as "active" to:
  - Recreate and write a minimal `.aws/config` and `.aws/credentials`
  - Overwrite the files cleanly without concern for legacy format

---

## 📦 Usage

TODO LATER

---

## 📁 Storage

Profiles are saved in:

```
~/.aws-profiles/config.json
```

Or, if encryption is enabled:

```
~/.aws-profiles/config-encrypted.bin
```