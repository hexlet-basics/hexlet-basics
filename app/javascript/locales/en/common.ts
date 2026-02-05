export default {
  check: {
    error: {
      headline: 'Oops!',
      message: 'Something went wrong, try one more time, please',
    },
    failed: {
      headline: 'Tests Failed',
      message:
        'Your code contains errors. Read carefully tests output and try to find errors yourself',
    },
    'failed-infinity': {
      headline: 'Tests Failed (Infinity Loop!)',
      message:
        'Maybe the test ran too long or you have an infinite loop in your solution. Make sure your loop has the correct breaking condition and you consider edge cases. If your loop works correctly, send it for checking again in 5 minutes. ',
    },
    passed: {
      headline: 'Tests passed',
      message:
        'Yippee! Well done! Compare your answer with the teacher’s one and move to the next lesson',
    },
  },
  confirm:
    'You want to reset the exercise progress. The current code version will not be saved. We hope you’ve already copied it. Continue resetting?',
  current_state: 'Current state',
  discuss: 'AI Assistent',
  editor: 'Editor',
  empty: 'Empty',
  errors: {
    network:
      'There was a network problem. Please try again. If it doesn’t work, make sure you have good internet and no blockers.',
    server:
      'Error on server. Maybe it’ll let go soon, but maybe not. Try to find out what happened in https://slack.hexlet.io/',
  },
  help: {
    controls: {
      body: "Reset progress. If something breaks and cannot be fixed click this button to return to the beginning of the exercises. The current code won’t be saved. If you'll need it, copy it before resetting. For example, to the notepads",
      header: 'Help',
    },
  },
  hours: {
    one: '%{count} hour',
    other: '%{count} hours',
  },
  instructions: 'Instructions',
  language_icon: '%{language} icon',
  languages: {
    en: 'English',
    ru: 'Russian',
  },
  lesson: 'Lesson',
  lessons: {
    one: '%{count} lesson',
    other: '%{count} lessons',
  },
  loading: 'Loading...',
  nextLesson: 'Next',
  organization: {
    address: 'The Republic of Kazakhstan, Almaty, Auezova St., 14A',
    email: 'support@hexlet.io',
    legal_name: 'TOO "Hexlet"',
  },
  output: 'Output',
  phones: ['+7 717 272 76 70'],
  prevLesson: 'Previous',
  reset: 'Reset',
  resetCode: 'Reset code',
  run: 'Run',
  showSolution: 'Show solution',
  signInSuggestion:
    '<a href="/users/new">Create an account</a> to save your progress',
  solution: 'Solution',
  solutionInstructions: "Teacher's solution will be available in:",
  solutionNotice:
    "It's best to solve the problem yourself, but if you're stuck for a long time, feel free to check out the solution. But make sure to study it thoroughly to truly understand it.",
  state_events: 'State Events',
  students: {
    one: '%{count} student',
    other: '%{count} students',
  },
  teacherSolution: "Teacher's solution:",
  testForExercise: 'Tests',
  testInstructions: 'Your exercise will be checked with these tests:',
  tos: 'Table Of Content',
  userCode: 'Your solution:',
  userCodeInstructions:
    "(start writing in Editor, your code will appear here and you'll be able to compare it to the teacher's solution)",
} as const;
