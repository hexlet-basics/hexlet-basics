export default {
  connection_refused:
    'Oops! Failed to connect to the Web Console middleware.\nPlease make sure a rails development server is running.\n',
  format: '%{attribute} %{message}',
  messages: {
    accepted: 'must be accepted',
    blank: "can't be blank",
    confirmation: "doesn't match %{attribute}",
    empty: "can't be empty",
    equal_to: 'must be equal to %{count}',
    even: 'must be even',
    exclusion: 'is reserved',
    greater_than: 'must be greater than %{count}',
    greater_than_or_equal_to: 'must be greater than or equal to %{count}',
    in: 'must be in %{count}',
    inclusion: 'is not included in the list',
    invalid: 'is invalid',
    less_than: 'must be less than %{count}',
    less_than_or_equal_to: 'must be less than or equal to %{count}',
    model_invalid: 'Validation failed: %{errors}',
    not_a_number: 'is not a number',
    not_an_integer: 'must be an integer',
    odd: 'must be odd',
    other_than: 'must be other than %{count}',
    password_too_long: 'is too long',
    present: 'must be blank',
    required: 'must exist',
    taken: 'has already been taken',
    too_long: {
      one: 'is too long (maximum is %{count} character)',
      other: 'is too long (maximum is %{count} characters)',
    },
    too_short: {
      one: 'is too short (minimum is %{count} character)',
      other: 'is too short (minimum is %{count} characters)',
    },
    wrong_length: {
      one: 'is the wrong length (should be %{count} character)',
      other: 'is the wrong length (should be %{count} characters)',
    },
  },
  template: {
    body: 'There were problems with the following fields:',
    header: {
      one: '%{count} error prohibited this %{model} from being saved',
      other: '%{count} errors prohibited this %{model} from being saved',
    },
  },
  unacceptable_request:
    'A supported version is expected in the Accept header.\n',
  unavailable_session:
    "Session %{id} is no longer available in memory.\n\nIf you happen to run on a multi-process server (like Unicorn or Puma) the process\nthis request hit doesn't store %{id} in memory. Consider turning the number of\nprocesses/workers to one (1) or using a different server in development.\n",
} as const;
