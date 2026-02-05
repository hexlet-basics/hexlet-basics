export default {
  currency: {
    format: {
      delimiter: ' ',
      format: '%n %u',
      precision: 2,
      separator: ',',
      significant: false,
      strip_insignificant_zeros: false,
      unit: 'руб.',
    },
  },
  format: {
    delimiter: ' ',
    precision: 3,
    round_mode: 'default',
    separator: ',',
    significant: false,
    strip_insignificant_zeros: false,
  },
  human: {
    decimal_units: {
      format: '%n %u',
      units: {
        billion: {
          few: 'миллиардов',
          many: 'миллиардов',
          one: 'миллиард',
          other: 'миллиардов',
        },
        million: {
          few: 'миллионов',
          many: 'миллионов',
          one: 'миллион',
          other: 'миллионов',
        },
        quadrillion: {
          few: 'квадриллионов',
          many: 'квадриллионов',
          one: 'квадриллион',
          other: 'квадриллионов',
        },
        thousand: {
          few: 'тысяч',
          many: 'тысяч',
          one: 'тысяча',
          other: 'тысяч',
        },
        trillion: {
          few: 'триллионов',
          many: 'триллионов',
          one: 'триллион',
          other: 'триллионов',
        },
        unit: '',
      },
    },
    format: {
      delimiter: '',
      precision: 1,
      significant: false,
      strip_insignificant_zeros: false,
    },
    storage_units: {
      format: '%n %u',
      units: {
        byte: {
          few: 'байта',
          many: 'байт',
          one: 'байт',
          other: 'байта',
        },
        eb: 'ЭБ',
        gb: 'ГБ',
        kb: 'КБ',
        mb: 'МБ',
        pb: 'ПБ',
        tb: 'ТБ',
      },
    },
  },
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
