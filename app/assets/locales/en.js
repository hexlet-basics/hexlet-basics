/* eslint-disable max-len */
export default {
  translation: {
    help: {
      controls: {
        header: 'Help',
        body: "Квадратная кнопка справа нужна для сброса прогресса, с её помощью можно восстановить начальное состояние упражнения. Это полезно, если во время экспериментов что-то сломалось и не чинится. Не забудьте, что сброс удалит ваш код, поэтому сохраните его куда-нибудь перед сбросом",
      },
    },
    errors: {
      network: 'There was a network problem. Please try again. If it doesn’t work, make sure you have good internet and no blockers.',
      server: 'Error on server. Maybe it’ll let go soon, but maybe not. Try to find out what happened in https://slack-ru.hexlet.io/',
    },
    run: 'Run',
    resetCode: 'Reset code',
    confirm: 'ATTENTION! Your code will be deleted. Are you sure you want to reset the current progress and return the exercise to its initial state?',
    editor: 'Editor',
    output: 'Output',
    solution: 'Solution',
    teacherSolution: "Teacher's soltuion:",
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
    nextLesson: 'Next Lesson',
    prevLesson: 'Previous Lesson',
    loading: 'Loading...',
    check: {
      error: {
        message: 'Something went wrong, try one more time, please',
        headline: 'Oops!',
      },
      passed: {
        message: 'Whoa! You did it! Proceed.',
        headline: 'Tests passed',
      },
      failed: {
        message: 'Fix errros and run again',
        headline: 'Tests Failed',
      },
      'failed-infinity': {
        message: 'Code was running too long. Check it for infinity loops.',
        headline: 'Tests Failed (Infinity Loop!)',
      },
    },
  },
};
