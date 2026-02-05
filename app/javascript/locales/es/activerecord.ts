export default {
  errors: {
    messages: {
      record_invalid: 'La validación falló: %{errors}',
      restrict_dependent_destroy: {
        has_many:
          'No se puede eliminar el registro porque existen %{record} dependientes',
        has_one:
          'No se puede eliminar el registro porque existe un %{record} dependiente',
      },
    },
  },
} as const;
