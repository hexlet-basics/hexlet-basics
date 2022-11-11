/* eslint-disable max-len */
export default {
  translation: {
    help: {
      controls: {
        header: 'Help',
        body: "Reset progress. If something breaks and cannot be fixed click this button to return to the beginning of the exercises. The current code won’t be saved. If you'll need it, copy it before resetting. For example, to the notepads",
      },
    },
    errors: {
      network: 'There was a network problem. Please try again. If it doesn’t work, make sure you have good internet and no blockers.',
      server: 'Error on server. Maybe it’ll let go soon, but maybe not. Try to find out what happened in https://slack.hexlet.io/',
    },
    signInSuggestion: '<a href="/users/new">Create an account</a> to save your progress',
    run: 'Run',
    resetCode: 'Reset code',
    confirm: 'You want to reset the exercise progress. The current code version will not be saved. We hope you’ve already copied it. Continue resetting?',
    editor: 'Editor',
    output: 'Output',
    solution: 'Solution',
    teacherSolution: "Teacher's solution:",
    testForExercise: 'Tests',
    userCode: 'Your solution:',
    userCodeInstructions: "(start writing in Editor, your code will appear here and you'll be able to compare it to the teacher's solution)",
    solutionInstructions: "Teacher's solution will be available in:",
    testInstructions: 'Your exercise will be checked with these tests:',
    solutionNotice: "It's best to solve the problem yourself, but if you're stuck for a long time, feel free to check out the solution. But make sure to study it thoroughly to truly understand it.",
    showSolution: 'Show solution',
    lesson: 'Lesson',
    discuss: 'Discussion',
    instructions: 'Instructions',
    nextLesson: 'Next',
    prevLesson: 'Previous',
    loading: 'Loading...',
    check: {
      error: {
        message: 'Something went wrong, try one more time, please',
        headline: 'Oops!',
      },
      passed: {
        message: 'Yippee! Well done! Compare your answer with the teacher’s one and move to the next lesson',
        headline: 'Tests passed',
      },
      failed: {
        message: 'Your code contains errors. Read carefully tests output and try to find errors yourself',
        headline: 'Tests Failed',
      },
      'failed-infinity': {
        message: 'Maybe the test ran too long or you have an infinite loop in your solution. Make sure your loop has the correct breaking condition and you consider edge cases. If your loop works correctly, send it for checking again in 5 minutes. ',
        headline: 'Tests Failed (Infinity Loop!)',
      },
    },
  },
};
