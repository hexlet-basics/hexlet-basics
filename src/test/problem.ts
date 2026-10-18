import { HttpResponse } from "msw";

// What the API answers a refused write with (ADR-0028): an RFC 9457 problem
// whose `errors` extension lists the messages by field. The client only reads
// `errors` off a body that is a whole ProblemDetails, so a mock that sends the
// extension alone is not what the server sends and is not shown at the field.
export const validationProblem = (errors: Record<string, string[]>) =>
  HttpResponse.json(
    {
      type: "about:blank",
      title: "Unprocessable Entity",
      status: 422,
      errors,
    },
    { status: 422, headers: { "Content-Type": "application/problem+json" } },
  );
