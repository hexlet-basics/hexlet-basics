/* eslint-disable max-len */
export default {
  translation: {
    help: {
      controls: {
        header: 'Ayuda',
        body: 'Restablecer el progreso. Puede regresar al inicio del ejercicio presionando este botón si algo está roto y no se arregla. El código actual no se guardará. Si aún lo necesita, cópielo antes de restablecerlo, por ejemplo, en un bloc de notas',
      },
    },
    errors: {
      network: 'Hubo un problema de red. Intente nuevamente. Si no funciona, asegúrese de que su conexión a Internet esté funcionando correctamente y desactive los bloqueadores de anuncios',
      server: 'Error en el servidor. Puede que se resuelva pronto o puede que no.',
    },
    signInSuggestion: '<a href="/users/new">Crea una cuenta</a> para guardar tu progreso',
    run: 'Verificar',
    resetCode: 'Restablecer código',
    confirm: '¿Desea restablecer el progreso del ejercicio? La versión actual del código no se guardará, esperamos que ya lo haya copiado. ¿Continuar con el restablecimiento?',
    editor: 'Editor',
    output: 'Salida',
    testForExercise: 'Pruebas',
    solution: 'Solución',
    teacherSolution: 'Solución del profesor:',
    userCode: 'Tu código:',
    userCodeInstructions: 'Cuando comiences a escribir una solución en el Editor, aparecerá aquí para compararla con la del profesor',
    solutionInstructions: 'La solución del profesor se abrirá en:',
    testInstructions: 'Tu ejercicio se evalúa con estas pruebas',
    solutionNotice: 'Es recomendable resolver el problema por ti mismo, pero si te quedas atascado durante mucho tiempo, consulta la solución del profesor. Sin embargo, asegúrate de entenderla y repasarla por tu cuenta',
    showSolution: 'Mostrar solución',
    lesson: 'Lección',
    discuss: 'Discusión',
    instructions: 'Instrucciones',
    nextLesson: 'Siguiente',
    prevLesson: 'Anterior',
    loading: 'Cargando...',
    check: {
      error: {
        message: 'Algo salió mal. Intente nuevamente',
        headline: '¡Ups!',
      },
      passed: {
        message: '¡Hurra! ¡Todo salió bien! Compara tu solución con la del profesor y pasa a la siguiente lección',
        headline: 'Pruebas superadas',
      },
      failed: {
        /*message: 'Tu código tiene errores. Lee detenidamente la salida de las pruebas, encuentra los errores y corrígelos. <a href="http://help.hexlet.io/ru/articles/111500-kak-nayti-oshibki-v-kode" target="_blank">Cómo encontrar errores en el código</a>. Si te resulta difícil aprender por ti mismo, considera <a href="https://hexlet.la/" target="_blank">aprender con un tutor</a>',*/
		message: 'Tu código tiene errores. Lee detenidamente la salida de las pruebas, encuentra los errores y corrígelos. Si te resulta difícil aprender por ti mismo, considera <a href="https://hexlet.la/" target="_blank">aprender con un tutor</a>',
        headline: 'Pruebas no superadas',
      },
      'failed-infinity': {
        message: 'Es posible que las pruebas se hayan ejecutado durante demasiado tiempo o que tu solución tenga un bucle infinito. Asegúrate de que haya una condición de salida adecuada en el bucle y de que estés considerando casos límite. Si el bucle está bien, intenta enviar el código para su evaluación nuevamente en 5 minutos',
        headline: 'Ejecución prolongada',
      },
    },
  },
};
