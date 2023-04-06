import 'select2';
// TODO: use setup webpack loader for eslint
// eslint-disable-next-line import/no-unresolved
import $ from 'jquery';

// eslint-disable-next-line func-names
$('.select2').each(function () {
  const el = $(this);
  el.select2({
    theme: 'bootstrap-5',
    ajax: {
      url: el.data('ajax--url'),
      data(params) {
        return { term: params.term };
      },
      delay: 250,
      cache: true,
      processResults(data) {
        const results = data.map((item) => ({
          id: item.id,
          text: `${item.text}`,
        }));
        return { results };
      },
    },
  });
});
