export default {
  currency: {
    format: {
      delimiter: ',',
      format: '%u%n',
      negative_format: '-%u%n',
      precision: 2,
      separator: '.',
      significant: false,
      strip_insignificant_zeros: false,
      unit: '$',
    },
  },
  format: {
    delimiter: ',',
    precision: 3,
    round_mode: 'default',
    separator: '.',
    significant: false,
    strip_insignificant_zeros: false,
  },
  human: {
    decimal_units: {
      format: '%n %u',
      units: {
        billion: 'Billion',
        million: 'Million',
        quadrillion: 'Quadrillion',
        thousand: 'Thousand',
        trillion: 'Trillion',
        unit: '',
      },
    },
    format: {
      delimiter: '',
      precision: 3,
      significant: true,
      strip_insignificant_zeros: true,
    },
    storage_units: {
      format: '%n %u',
      units: {
        byte: {
          one: 'Byte',
          other: 'Bytes',
        },
        eb: 'EB',
        gb: 'GB',
        kb: 'KB',
        mb: 'MB',
        pb: 'PB',
        tb: 'TB',
        zb: 'ZB',
      },
    },
  },
  nth: {},
  percentage: {
    format: {
      delimiter: '',
      format: '%n%',
    },
  },
  precision: {
    format: {
      delimiter: '',
    },
  },
} as const;
