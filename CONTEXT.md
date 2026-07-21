# Hexlet Basics

A platform of free introductory programming courses. Each course teaches one
programming language through a sequence of in-browser coding lessons, with
per-learner progress tracking, localization, and a versioned content pipeline.

## Language

> This glossary records the ubiquitous language of the project. Where the
> spoken term differs from the code name, both are given.

### Learning content

**Course**:
A body of lessons that teaches one programming language. The central thing a
learner studies.
_Code name_: `Language`.
_Avoid_: language, programming language.

**Locale**:
The human language a piece of content is written in. The axis of translation.
_Avoid_: language.

**Version** (Course Version):
A built, immutable snapshot of a Course's whole content tree. A Course has one
live snapshot at a time.
_Code name_: `Language::Version`.
_Avoid_: build, release.

**Module**:
A named grouping of Lessons within a Course — the chapter level.
_Code name_: `Language::Module`.

**Lesson**:
The stable identity of one unit of content — the thing a learner opens and
solves in the browser. Keeps its identity across Versions.
_Code name_: `Language::Lesson`.

**Lesson Version / Module Version**:
A Lesson or Module as it exists inside one Course Version. The same Lesson has a
distinct Lesson Version in each Course Version.
_Code names_: `Language::Lesson::Version`, `Language::Module::Version`.

**Info**:
The localized text of a Course, Module Version, or Lesson Version (its name,
theory, instructions, tips) for one Locale.
_Code names_: `…::Version::Info`.
_Avoid_: translation, content.

**Category**:
A catalog grouping of Courses, used to organize the public catalog and landing
pages. A Course joins a Category through its Landing Page, not directly.
_Code name_: `Language::Category`.

> The Category-to-Course membership is a plain join record
> (`Language::Category::Item`, linking a Category to a Landing Page) — a
> mechanism, not a domain term.

**Landing Page**:
A localized marketing page for a Course.
_Code name_: `Language::LandingPage`.

### Readiness vs progress

**Readiness** (of a Course):
The authoring maturity of a Course — how finished its content is. A property of
the content, set by authors.
_Code name_: `Language#progress`.
_Avoid_: progress (for this meaning).

**Progress** (of a learner):
How far a learner has advanced through a Course. A property of the learner.
_Avoid_: readiness.

### Enrollment & learner state

**Enrollment**:
A user's participation in a Course and their Course-level Progress. A record
about a user's relationship to a Course — **not** a person. A user who could
begin a Course but has not is "ready to start."
_Code name_: `Language::Member`.
_Avoid_: member, membership.

**Lesson Progress**:
A learner's state on a single Lesson within an Enrollment. The unit that drives
Course-level Progress, and what a Lesson Assistant attaches to.
_Code name_: `Language::Lesson::Member`.
_Avoid_: lesson member, lesson attempt.

**First language / Second language**:
Course positioning — whether the Course is aimed at people new to programming or
at developers learning an additional language.
_Code name_: `Language#learn_as`.

### Identity & access

**User**:
A person on the platform. Reachable by at least one of: email, phone, or a
Linked Account.
_Avoid_: account, member, learner (in code).

**Linked Account**:
An external OAuth identity bound to a User (e.g. GitHub, Facebook). A sign-in
route — not the User and not billing.
_Code name_: `User::Account`.
_Avoid_: account, social identity.

**Passkey**:
A WebAuthn credential a User registers for passwordless sign-in.
_Code name_: `User::Credential`.
_Avoid_: credential.

**Admin**:
A User with platform-wide superuser access that bypasses resource permission
checks.
_Code name_: `User#admin`.

**Staff**:
A User with elevated back-office access — either an Admin or a Staff Member.
_Code name_: `User#staff?`.

**Staff Member**:
The grant that makes a User staff: a Role plus the content Locales they may
manage. Unrelated to Enrollment — a different "member."
_Code name_: `StaffMember`.
_Avoid_: member.

**Role**:
A named set of Permissions assigned to Staff Members.
_Code name_: `StaffMember::Role`.

**Permission**:
A Role's create/read/update/destroy rights over one admin Resource.
_Code name_: `StaffMember::Role::Permission`.

**Resource** (admin):
A type of back-office–managed entity a Permission applies to.
_Code name_: `StaffMember::Role::Permission::Resource`.

### AI tutor

**Lesson Assistant**:
The AI conversation that helps one learner with one Lesson — the in-lesson help
channel.
_Code name_: `AiChat`.
_Avoid_: AI chat, AI tutor.

> The supporting `AiMessage`, `AiModel`, and `AiToolCall` are standard RubyLLM
> plumbing, not domain terms.

### Growth & engagement

**Lead**:
A captured sales contact — a User's contact details, survey answers, and
Courses of interest, attributed to the Visit it came from.
_Code name_: `Lead`.

**Book Request**:
A User's request to download the free book (a lead magnet).
_Code name_: `BookRequest`.

**Banner**:
A scheduled, localized site-wide promo strip.
_Code name_: `Banner`.

**Blog Post**:
A localized blog article, authored by a User and optionally tied to a Course.
_Code name_: `BlogPost`.

**Like** (of a Blog Post):
A User's like on a Blog Post. May be anonymous — recorded without a User.
_Code name_: `BlogPost::Like`.

**Related Course** (of a Blog Post):
An ordered link from a Blog Post to a Course, shown as further reading.
_Code name_: `BlogPost::RelatedLanguageItem`.
_Avoid_: related language, related item.

**Review**:
A published learner testimonial about a Course. Public social proof.
_Code name_: `Review`.
_Avoid_: testimonial, отзыв (in code), feedback report.

**Lesson Feedback Report**:
An AI-generated report that analyzes learners' questions to the Lesson Assistant
and suggests how to improve a Lesson, per Locale. A diagnostics artifact for
authors — **not** a learner testimonial.
_Code name_: `Language::Lesson::Review`.
_Avoid_: lesson review, lesson summary, review.

### Surveys

**Survey**:
A single localized question with a set of Options. May branch — a Survey can be
a follow-up to a specific Option of another Survey.
_Code name_: `Survey`.
_Avoid_: question, poll, form (a Survey is one question, not a multi-question form).

**Survey Option**:
One selectable answer to a Survey's question.
_Code name_: `Survey::Item`.
_Avoid_: item, choice.

**Survey Answer**:
A User's response to a Survey — the Option they picked.
_Code name_: `Survey::Answer`.
_Avoid_: response.

**Scenario**:
An ordered sequence of Surveys delivered to a User as one campaign, started by a
Trigger.
_Code name_: `Survey::Scenario`.
_Avoid_: survey flow, questionnaire.

**Scenario Step**:
A Survey's position within a Scenario.
_Code name_: `Survey::Scenario::Item`.
_Avoid_: item.

**Scenario Trigger**:
An event rule that starts a Scenario when a given Event occurs enough times.
_Code name_: `Survey::Scenario::Trigger`.

**Scenario Participation**:
A User's progress through a Scenario.
_Code name_: `Survey::Scenario::Member`.
_Avoid_: member, scenario run.

> **Note on "member":** the code uses `…Member` for four unrelated
> user-relationship records — Enrollment (`Language::Member`), Lesson Progress
> (`Language::Lesson::Member`), Staff Member (`StaffMember`), and Scenario
> Participation (`Survey::Scenario::Member`). In speech, prefer the specific
> term; never say a bare "member."

### Catalog & analytics

**Q&A Item**:
A question/answer pair (FAQ entry) shown on a Category or Landing Page.
_Code names_: `Language::Category::QnaItem`, `Language::LandingPage::QnaItem`.
_Avoid_: FAQ, qna item.

**Visit / Event**:
Web analytics records of a browsing session and the named Events fired within
it. Domain-relevant because Event names are the hooks a Scenario Trigger fires
on and the signal a Lead is attributed to.
_Code names_: `Ahoy::Visit`, `Ahoy::Event`.
