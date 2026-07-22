# A successful Course build auto-promotes to live

When a Course Version finishes building successfully it becomes the live Version
in the same step — there is no manual publish, staging, or preview gate between
"built" and "shown to learners." We chose this for operational simplicity: a
rebuild is the only editing path (content is never edited in place), so a
separate approval step would add friction to every content fix without a clear
owner. The trade-off is that we cannot preview a built Version before learners
see it, and a bad build that still *succeeds* goes live immediately.

The one safety property we rely on: a build that **fails** leaves the previously
live Version in place, so a broken rebuild never takes a Course down — it just
fails to advance.
