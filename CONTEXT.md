# Domain Glossary

## Language

### Accounts

**User sign-up**:
The first successful creation of a User account through public registration.
Later sign-ins, profile changes, and administrator-created users are not
sign-ups.

### Learning

**Course Membership**:
A learner's relationship with one Course, created the first time they open any
of its Lessons and never created twice for the same pair.
_Avoid_: Language Member, enrollment, subscription

**Lesson Membership**:
A learner's relationship with one Lesson of a Course they are a member of.
_Avoid_: Language Lesson Member, attempt, submission

**Not started**:
The condition of a learner who has no Course Membership for a Course. It is the
absence of a membership, not a state a membership can be in.
_Avoid_: ready to start

**Started**:
The state of a Course or Lesson Membership from creation until it is finished.

**Finished**:
The state of a Lesson Membership whose exercise the learner has passed, and of a
Course Membership all of whose current Lessons are finished. A membership never
returns to Started.
_Avoid_: completed, done, passed

**Current Lessons**:
The Lessons belonging to a Course's current Version. A Course's Lesson set
changes when a new Version is promoted, so what counts as finishing the Course
is defined against the current Version and not against Lessons a learner may
have finished under an older one.

**Progress**:
The share of a Course's Current Lessons a learner has finished, as a percentage.
Lessons finished under a retired Version do not count towards it.

**Next Lesson**:
The first Lesson, in course order, among a Course's Current Lessons that the
learner has not finished.

**Guest progress**:
Lessons a visitor passed before having an account. It is claimed once, when that
visitor signs up, and becomes ordinary Course and Lesson Memberships in the
Finished state.
_Avoid_: anonymous progress, trial progress
