import 'select2';
import $ from 'jquery';

// eslint-disable-next-line func-names
$('.select2').each(function () {
  const el = $(this);
  el.select2({
    theme: 'bootstrap-5',
    // ajax: {
    //   url: el.data('ajax--url'),
    //   dataType: 'json',
    // },
  });
});
