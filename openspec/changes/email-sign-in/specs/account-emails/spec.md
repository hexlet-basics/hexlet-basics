# Spec Delta

## Purpose

Defines how the emails that let a User back into their account are localized,
throttled and delivered, independently of which flow asked for them.

## ADDED Requirements

### Requirement: Account emails are sent in the request's locale
An account email SHALL be written in the locale of the request that asked for
it, and the link it carries SHALL open the site in that same locale. A request
with no supported locale SHALL get `ru`.

#### Scenario: Spanish request
- **WHEN** a visitor on the `es` site requests a Magic Link
- **THEN** the email subject and body are in Spanish
- **AND** the link points at the `/es/` URL of the page it opens

#### Scenario: Russian request
- **WHEN** a visitor on the `ru` site requests a Password Reset
- **THEN** the email is in Russian and its link opens the `ru` site

#### Scenario: Request with no supported locale
- **WHEN** a Password Reset is requested with no locale, or with `de`
- **THEN** the email is in Russian

### Requirement: Repeated requests within a minute send one email
A second request for the same kind of email to the same address within one
minute of the first SHALL NOT send another email, and SHALL answer exactly as
the first request did.

#### Scenario: Double click on "send"
- **WHEN** a Magic Link is requested twice for one address thirty seconds apart
- **THEN** exactly one email is delivered
- **AND** both requests receive the same success response

#### Scenario: After the cooldown
- **WHEN** a Magic Link is requested again more than one minute later
- **THEN** a new email is delivered

### Requirement: Delivery does not hold up the request
Asking for an account email SHALL answer before the email is delivered. A
delivery that fails SHALL be retried, and a delivery that keeps failing SHALL be
reported to error monitoring; the visitor SHALL NOT be told.

#### Scenario: Mail provider is down
- **WHEN** the mail provider rejects every delivery attempt
- **THEN** the request that asked for the email still succeeded
- **AND** after the retries run out, the failure appears in error monitoring

### Requirement: Non-production environments never send real email
Outside production an account email SHALL NOT reach a mail provider. In
development its recipient, subject and link SHALL be written to the log; in
tests it SHALL be observable by the test.

#### Scenario: Developer requests a Magic Link
- **WHEN** a Magic Link is requested on a development server
- **THEN** the log contains the recipient and the link, and no email is sent
