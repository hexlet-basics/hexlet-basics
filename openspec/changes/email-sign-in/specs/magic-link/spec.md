# Spec Delta

## Purpose

Lets an existing User sign in without their password by following a
short-lived link emailed to them.

## ADDED Requirements

### Requirement: Requesting a Magic Link reveals nothing about the email
Requesting a Magic Link SHALL always succeed with no content, whether or not the
email belongs to a User. Only an existing User SHALL be emailed, and no account
SHALL ever be created by the request.

#### Scenario: Registered email
- **WHEN** a visitor requests a Magic Link for a registered email
- **THEN** the response is 204
- **AND** that User receives a Magic Link email

#### Scenario: Unknown email
- **WHEN** a visitor requests a Magic Link for an email no User has
- **THEN** the response is 204, identical to the registered case
- **AND** no email is sent and no User is created

#### Scenario: Email differs only in case or surrounding spaces
- **WHEN** a visitor requests a Magic Link for ` Learner@Example.com `
- **THEN** the User registered as `learner@example.com` receives the email

### Requirement: A Magic Link signs its User in
Following a valid Magic Link SHALL sign its User in exactly as a password
sign-in does: the session cookie is set, the visitor's guest progress is
credited to the account, and a sign-in is recorded. The response SHALL be the
signed-in User.

#### Scenario: Guest with progress follows the link
- **WHEN** a guest who finished two Lessons follows their Magic Link
- **THEN** they are signed in as the link's User
- **AND** those two Lessons are credited to the account and the guest cookie is cleared

#### Scenario: Link followed twice
- **WHEN** a User follows the same Magic Link a second time within its lifetime
- **THEN** they are signed in again

### Requirement: A Magic Link dies after fifteen minutes or an email change
A Magic Link SHALL be refused once fifteen minutes have passed since it was
issued, or once its User's email has changed. A refused, forged or malformed
link SHALL answer not found and SHALL NOT sign anyone in.

#### Scenario: Expired link
- **WHEN** a Magic Link is followed sixteen minutes after it was issued
- **THEN** the response is 404 and no session cookie is set

#### Scenario: Email changed after the link was sent
- **WHEN** a User changes their email and then follows a Magic Link sent to the old address
- **THEN** the response is 404

#### Scenario: Link issued by the legacy app
- **WHEN** a link emailed by the Rails app before cutover is followed
- **THEN** the response is 404

### Requirement: Magic Link pages keep legacy URLs
The request form SHALL live at `/magic_links/new` and the emailed link SHALL
open `/magic_links/:token`, each under the locale prefix. A refused link SHALL
send the visitor back to the request form with an explanation.

#### Scenario: Visitor follows a refused link
- **WHEN** a visitor opens `/magic_links/<expired token>`
- **THEN** they land on `/magic_links/new` with a message that the link is no longer valid
