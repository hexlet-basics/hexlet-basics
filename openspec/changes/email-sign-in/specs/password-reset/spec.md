# Spec Delta

## Purpose

Lets a User who forgot their password set a new one through a short-lived
emailed link, and fixes the rule every newly set password must meet.

## ADDED Requirements

### Requirement: Requesting a Password Reset reveals nothing about the email
Requesting a Password Reset SHALL always succeed with no content, whether or not
the email belongs to a User. Only an existing User SHALL be emailed.

#### Scenario: Unknown email
- **WHEN** a visitor requests a Password Reset for an email no User has
- **THEN** the response is 204
- **AND** no email is sent

#### Scenario: Registered email
- **WHEN** a visitor requests a Password Reset for a registered email
- **THEN** the response is 204 and that User receives a Password Reset email

### Requirement: A reset link can be checked before the form is shown
Checking a reset link SHALL answer no content while the link is valid and not
found otherwise, without changing anything.

#### Scenario: Valid link
- **WHEN** the reset page checks a link issued five minutes ago
- **THEN** the response is 204

#### Scenario: Forged link
- **WHEN** the reset page checks a token that was not issued by the system
- **THEN** the response is 404

### Requirement: Setting a new password signs the User in
Submitting a new password through a valid reset link SHALL replace the User's
password and sign them in exactly as a password sign-in does (session cookie,
guest progress credited, sign-in recorded). The response SHALL be the signed-in
User.

#### Scenario: Successful reset
- **WHEN** a User submits `new-secret` through a valid reset link
- **THEN** they are signed in
- **AND** signing in afterwards works with `new-secret` and not with the old password

### Requirement: A reset link dies after fifteen minutes or one use
A reset link SHALL be refused once fifteen minutes have passed since it was
issued, or once the User's password has changed by any means. A refused link
SHALL answer not found on both check and submit.

#### Scenario: Link reused after a successful reset
- **WHEN** a User resets their password and then submits the same link again
- **THEN** the response is 404 and the password is unchanged

#### Scenario: Expired link
- **WHEN** a reset link is submitted sixteen minutes after it was issued
- **THEN** the response is 404

### Requirement: Passwords are at least six characters
Every password a User sets — at sign-up or through a reset — SHALL be at least
six characters long. The form SHALL show the error on the password field before
submitting; the API SHALL reject a shorter one as a bad request, like any other
contract constraint, and nothing SHALL change.

#### Scenario: Short password on reset
- **WHEN** `12345` is submitted through a valid reset link
- **THEN** the response is 400
- **AND** the link is still valid and the password is unchanged

#### Scenario: Short password on sign-up
- **WHEN** a visitor signs up with the password `12345`
- **THEN** the response is 400 and no User is created

#### Scenario: Short password typed into a form
- **WHEN** a visitor types `12345` into the sign-up or reset form and submits
- **THEN** the password field shows an error and no request is sent

### Requirement: Password Reset pages keep legacy URLs
The request form SHALL live at `/remind_password/new` and the emailed link SHALL
open `/password/:token/edit`, each under the locale prefix. A refused link SHALL
send the visitor to the request form with an explanation.

#### Scenario: Visitor opens an expired reset link
- **WHEN** a visitor opens `/password/<expired token>/edit`
- **THEN** they land on `/remind_password/new` with a message that the link is no longer valid
